import SecureMessenger from '../secure-messenger/secure-messenger';

// First lets lockdown and run sanity checks
async function runChecks() {
    ['indexedDB', 'caches', 'CacheStorage', 'Caches', 'postMessage', 'close'].forEach(o => {
        let t = self;
        while (!!Object.getOwnPropertyDescriptor(t, o) || !!t.__lookupGetter__(o)) {
            Object.defineProperty(t, o, {
                get() {
                    throw new Error(`Disallowed Exception: You cannot use global property ${o} in a track builder script`);
                },
                configurable: false
            });
            t = Object.getPrototypeOf(t);
        }
    });

    // Test data offloading CSP
    let failedChecks = {
        remoteConnectionCheck: true,
        idbCheck: true,
        caches: true,
        close: true,
        postMessage: true
    }

    try {
        let response = await fetch('http://example.com/this-fetch-should-fail');
    } catch (e) {
        failedChecks.remoteConnectionCheck = false;
    }

    try {
        indexedDB;
    } catch (e) {
        failedChecks.idbCheck = false;
    }

    try {
        caches;
    } catch (e) {
        failedChecks.caches = false;
    }

    try {
        close();
    } catch (e) {
        failedChecks.close = false;
    }

    try {
        postMessage('foo');
    } catch (e) {
        failedChecks.postMessage = false;
    }

    for (let key in failedChecks) {
        let failedCheck = failedChecks[key];
        if (failedCheck) {
            return true;
        }
    }

    return false;
}

// Now lets run code
(async () => {
    const kdTree = require('./kdtree').default;
    const {getTrackTag} = require('../common/normalize');

    let tags;
    let idToTagMap;
    let tracks;
    let unprunedTracks;
    let additionalTracks = {};
    let extraDimensions = [];
    let tree;
    let getTree = () => tree;
    let getIdToTagMap = () => idToTagMap;
    let registerDimension = dim => extraDimensions.push(dim);

    require('./api.js').injector({getTree, getIdToTagMap, registerDimension});

    let messenger = new SecureMessenger(location.origin);
    let loadingInterruptPort;
    await messenger.knock([self]);
    let failedChecks = await runChecks();
    messenger.messageHandlerLoop((data, interruptPort, finish) => {
        if (failedChecks) {
            finish();
            console.error('Checks failed');
            return {type: 'error', error: 'checks failed'};
        }
        if (data.type == 'buildPlaylist') {
            finish();
            return handleBuildPlaylist(data);
        } else if (data.type == 'pruneTracks') {
            finish();
            return handlePruneTracks(data);
        } else if (data.type == 'initializeLoadingNotificationChannel') {
            loadingInterruptPort = interruptPort;
            return {type: 'acknowledge'};
        }
        return {type: 'error', error: 'invalid request'};
    }, event => {
        console.error('Eos worker could not handle Eos message', event);
    });

    async function handleTrackData(data) {
        let trackData = await (await fetch(data.tracksUrl)).json();
        let additionalTrackData = await (await fetch(data.additionalTracksUrl)).json();
        tags = trackData.tags;
        idToTagMap = trackData.idsToTags;
        unprunedTracks = trackData.tracks;
        tracks = {};
        if (data.prunedTrackIds) {
            for (let prunedTrackId of data.prunedTrackIds) {
                tracks[prunedTrackId] = unprunedTracks[prunedTrackId];
            }
        } else {
            tracks = unprunedTracks;
        }
        for (let [id, track] of Object.entries(additionalTrackData)) {
            if (tags[track.tag]) {
                tags[track.tag].push(id);
            }
            else {
                tags[track.tag] = [id];
            }
            idToTagMap[id] = track.tag;
            additionalTracks[id] = {
                track: track.track,
                features: track.features,
                evocativeness: track.evocativeness
            };
        }
        console.info(`${Object.keys(unprunedTracks).length} unpruned tracks, ${Object.keys(tracks).length} pruned tracks, ${Object.keys(additionalTracks).length} additional tracks`)
    }

    async function handleBuildPlaylist(data) {
        await handleTrackData(data);
        let trackCount = data.trackCount;
        let firstTrack;
        if (data.firstTrackOnly) {
            console.debug('Getting first track...');
            trackCount = 1;
        }
        else {
            if (data.firstTrack) {
                firstTrack = data.firstTrack;
            }
            console.debug('Getting tracks...');
        }
        let playlist;
        try {
            playlist = await buildPlaylist(data.script, trackCount, firstTrack);
        }
        catch (e) {
            console.error('Could not build playlist:', e);
            return {type: 'playlistError', error: e.message};
        }
        let dimensions = [...extraDimensions];
        if (tree) {
            dimensions = [...dimensions, ...tree.getDimensions()];
        }
        return {type: 'playlist', playlist, dimensions};
    }

    async function handlePruneTracks(data) {
        await handleTrackData(data);
        console.debug('Pruning...');
        let prunedTrackIds;
        try {
            prunedTrackIds = await prune(data.script);
        }
        catch (e) {
            console.error('Could not prune tracks:', e);
            return {type: 'playlistError', error: e.message};
        }
        console.debug(`Tracks pruned down to ${prunedTrackIds.length}`);
        loadingInterruptPort.postMessage({type: 'loadPercent', value: 1});
        return {type: 'prunedTracks', prunedTrackIds, dimensions: extraDimensions};
    }

    async function prune(script) {
        let blobUrl = URL.createObjectURL(new Blob([script], {type: 'application/javascript'}));
        importScripts(blobUrl);
        let {data: prunedTracksUnsafe} = await self.hooks.prune(JSON.parse(JSON.stringify({tracks, idToTagMap, unprunedTracks})));
        let safeTracks = validatePrunedTracks(prunedTracksUnsafe, tracks);
        return Object.keys(safeTracks);
    }

    async function buildPlaylist(script, goalTracks, firstTrack) {
        if (!goalTracks) {
            goalTracks = 20;
        }
        let blobUrl = URL.createObjectURL(new Blob([script], {type: 'application/javascript'}));
        importScripts(blobUrl);
        let {data: prunedTracksUnsafe} = await self.hooks.prune(JSON.parse(JSON.stringify({tracks, idToTagMap, unprunedTracks})));
        let prunedTracks = validatePrunedTracks(prunedTracksUnsafe, tracks);
        let points = tracksToPoints(prunedTracks);
        let {data} = await self.hooks.buildTree(kdTree, JSON.parse(JSON.stringify({tracks: prunedTracks, idToTagMap, points})));
        tree = data;
        let tags = {};
        let playlist = []
        if (!firstTrack) {
            firstTrack = await getTrack(function() {
                return self.hooks.getFirstTrack(JSON.parse(JSON.stringify({playlist, tags, goalTracks, points, tracks: prunedTracks})), tree);
            }, prunedTracks);
            if (!firstTrack) {
                throw new Error('Builder was unable to get a first track');
            }
        }
        else {
            tree && tree.removeById(firstTrack.id);
        }
        loadingInterruptPort.postMessage({type: 'loadPercent', value: 1 / goalTracks});
        tags[getTrackTag(firstTrack.track)] = true;
        playlist.push(firstTrack);
        for (let i = 0; i < goalTracks - 1; i++) {
            let nextTrack = await getTrack(function() {
                let previousTrack = playlist[playlist.length - 1];
                return self.hooks.getNextTrack(JSON.parse(JSON.stringify({playlist, tags, goalTracks, points, tracks: prunedTracks, previousTrack})), tree);
            }, prunedTracks);
            if (!nextTrack) {
                console.warn(`Builder was unable to get track ${i + 2}`);
                break;
            }
            tags[getTrackTag(nextTrack.track)] = true;
            playlist.push(nextTrack);
            loadingInterruptPort.postMessage({type: 'loadPercent', value: (i + 2) / goalTracks});
        }
        return JSON.parse(JSON.stringify(playlist));
    }

    async function getTrack(getTrackFn, allowedTracks) {
        let {data: unsafeTrack} = await getTrackFn();
        if (!unsafeTrack) {
            return null;
        }
        let {point} = unsafeTrack;
        let safeTrack = allowedTracks[point.id];
        if (!safeTrack) {
            throw new Error(`Builder returned invalid track ${point.id}`);
        }
        tree && tree.removeById(point.id);
        return safeTrack;
    }

    function validatePrunedTracks(unsafeTracks, allowedTracks) {
        let safeTracks = {};
        for (let trackId of Object.keys(unsafeTracks)) {
            let safeTrack = allowedTracks[trackId];
            if (!safeTrack) {
                throw new Error(`Builder pruning returned invalid track ${trackId}`);
            }
            safeTracks[trackId] = safeTrack;
        }
        return safeTracks;
    }
})();
