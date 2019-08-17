(() => {
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

    addEventListener('message', async ({data}) => {
        // The secret is used so that the worker code cannot try and falsify a postMessage to `self`
        // The message handlers set up to listen to this worker pass in a secret they expect back
        // on response which lives inside this closure and cannot be accessed by user code.
        let secret = data.secret;
        if (data.type !== 'buildPlaylist' && data.type !== 'pruneTracks') {
            console.error('Unexpected request', data.type);
            self.postMessage({type: 'playlistError', error: 'unexpected request type', secret});
            return;
        }
        let trackData = await new Promise(async resolve => {
            let response = await fetch(data.tracksUrl);
            resolve(await response.json());
        });
        let additionalTrackData = await new Promise(async resolve => {
            let response = await fetch(data.additionalTracksUrl);
            resolve(await response.json());
        });
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
        console.log(`${Object.keys(unprunedTracks).length} unpruned tracks, ${Object.keys(tracks).length} pruned tracks, ${Object.keys(additionalTracks).length} additional tracks`)
        if (data.type === 'pruneTracks') {
            console.log('Pruning...');
            let prunedTrackIds;
            try {
                prunedTrackIds = await prune(data.script);
            }
            catch (e) {
                console.error('Could not prune tracks:', e);
                self.postMessage({type: 'playlistError', error: e.message, secret});
                return;
            }
            console.log(`Tracks pruned down to ${prunedTrackIds.length}`);
            self.postMessage({type: 'prunedTracks', prunedTrackIds, dimensions: extraDimensions, secret});
        } else {
            let trackCount = data.trackCount;
            let firstTrack;
            if (data.firstTrackOnly) {
                console.log('Getting first track...');
                trackCount = 1;
            }
            else {
                if (data.firstTrack) {
                    firstTrack = data.firstTrack;
                }
                console.log('Getting tracks...');
            }
            let playlist;
            try {
                playlist = await buildPlaylist(data.script, trackCount, firstTrack);
            }
            catch (e) {
                console.error('Could not build playlist:', e);
                self.postMessage({type: 'playlistError', error: e.message, secret});
                return;
            }
            let dimensions = [...extraDimensions];
            if (tree) {
                dimensions = [...dimensions, ...tree.getDimensions()];
            }
            self.postMessage({type: 'playlist', playlist, dimensions, secret});
        }
    });

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
            tree && tree.removeById(firstTrack.track.id);
        }
        tags[getTrackTag(firstTrack.track)] = true;
        playlist.push(firstTrack);
        for (let i = 0; i < goalTracks - 1; i++) {
            let nextTrack = await getTrack(function() {
                let previousTrack = playlist[playlist.length - 1];
                return self.hooks.getNextTrack(JSON.parse(JSON.stringify({playlist, tags, goalTracks, points, tracks: prunedTracks, previousTrack})), tree);
            }, prunedTracks);
            if (!nextTrack) {
                throw new Error(`Builder was unable to get track ${i + 2}`);
            }
            tags[getTrackTag(nextTrack.track)] = true;
            playlist.push(nextTrack);
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
