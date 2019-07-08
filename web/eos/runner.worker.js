(() => {
    const kdTree = require('./kdtree').default;
    const {getTrackTag} = require('../common/normalize');

    let tags;
    let idToTagMap;
    let tracks;
    let additionalTracks = {};
    let tree;
    let getTree = () => tree;
    let getIdToTagMap = () => idToTagMap;

    require('./api.js').injector({getTree, getIdToTagMap});

    addEventListener('message', async ({data}) => {
        // The secret is used so that the worker code cannot try and falsify a postMessage to `self`
        // The message handlers set up to listen to this worker pass in a secret they expect back
        // on response which lives inside this closure and cannot be accessed by user code.
        let secret = data.secret;
        if (data.type !== 'buildPlaylist') {
            console.error('Could not build playlist, unexpected request', data.type);
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
        tracks = trackData.tracks;
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
        let trackCount = data.trackCount;
        let firstTrack;
        if (data.firstTrackOnly) {
            trackCount = 1;
        }
        else if (data.firstTrack) {
            firstTrack = data.firstTrack;
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
        self.postMessage({type: 'playlist', playlist, dimensions: tree.getDimensions(), secret});
    });

    async function buildPlaylist(script, goalTracks, firstTrack) {
        if (!goalTracks) {
            goalTracks = 20;
        }
        let blobUrl = URL.createObjectURL(new Blob([script], {type: 'application/javascript'}));
        importScripts(blobUrl);
        let points = tracksToPoints(tracks);
        tree = await self.hooks.buildTree(kdTree, JSON.parse(JSON.stringify({tracks, idToTagMap, points})));
        let tags = {};
        let playlist = []
        if (!firstTrack) {
            firstTrack = await getTrack(function() {
                return self.hooks.getFirstTrack(JSON.parse(JSON.stringify({playlist, tags, goalTracks})), tree);
            });
            if (!firstTrack) {
                throw new Error('Builder was unable to get a first track');
            }
        }
        else {
            tree.removeById(firstTrack.track.id);
        }
        tags[getTrackTag(firstTrack.track)] = true;
        playlist.push(firstTrack);
        for (let i = 0; i < goalTracks - 1; i++) {
            let nextTrack = await getTrack(function() {
                return self.hooks.getNextTrack(JSON.parse(JSON.stringify({playlist, tags, goalTracks})), JSON.parse(JSON.stringify(playlist[playlist.length - 1])), tree);
            });
            if (!nextTrack) {
                throw new Error(`Builder was unable to get track ${i + 2}`);
            }
            tags[getTrackTag(nextTrack.track)] = true;
            playlist.push(nextTrack);
        }
        return JSON.parse(JSON.stringify(playlist));
    }

    async function getTrack(getTrackFn) {
        let track = await getTrackFn();
        if (!track) {
            return null;
        }
        let {point} = track;
        tree.removeById(point.id);
        return tracks[point.id];
    }
})();
