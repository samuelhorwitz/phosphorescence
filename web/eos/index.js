import RunnerWorker from 'worker-loader?inline=true&fallback=false!./runner.worker.js';
import {encoder, decoder} from '../common/textencoding';

let additionalTracks = {};
let additionalTracksUrl;
let tracksReadyResolver;
let tracksReady = new Promise(resolve => tracksReadyResolver = resolve);

async function destroyIDB() {
    // Since we cannot completely sandbox this origin, user's may create Indexed DBs.
    // We don't want them to do that and will destroy them.
    (await indexedDB.databases()).forEach(db => indexedDB.deleteDatabase(db.name));
}

addEventListener('message', async ({origin, data}) => {
    if (origin === process.env.PHOSPHOR_ORIGIN) {
        let {responsePort} = data;
        if (data.type === 'buildPlaylist') {
            await destroyIDB();
            let finished = false;
            let runner = new RunnerWorker();
            let secret = crypto.getRandomValues(new Uint8Array(32)).reduce((str, byte) => str + byte.toString(16).padStart(2, '0'), '');
            runner.addEventListener('message', async e => {
                if (finished) {
                    return;
                }
                if (e.data.secret === secret) {
                    if (e.data.type === 'playlist') {
                        responsePort.postMessage({type: 'playlist', playlist: e.data.playlist, dimensions: e.data.dimensions});
                    }
                    else {
                        responsePort.postMessage({type: 'playlistError', error: e.data.error});
                    }
                }
                else {
                    responsePort.postMessage({type: 'playlistError', error: 'Blocked potentially-malicious playlist builder hijacking'});
                }
                runner.terminate();
                finished = true;
                await destroyIDB();
            }, {once: true});
            addEventListener('message', async ({origin, data}) => {
                if (finished) {
                    return;
                }
                if (origin === process.env.PHOSPHOR_ORIGIN && data.type === 'terminateAll') {
                    responsePort.postMessage({type: 'playlistError', error: 'User killed builder'});
                    runner.terminate();
                    finished = true;
                    await destroyIDB();
                }
            }, {once: true});
            let tracksUrl = await tracksReady;
            let trackData = await new Promise(async resolve => {
                let response = await fetch(tracksUrl);
                resolve(await response.arrayBuffer());
            });
            let additionalTrackData = new ArrayBuffer(0);
            if (additionalTracksUrl) {
                additionalTrackData = await new Promise(async resolve => {
                    let response = await fetch(additionalTracksUrl);
                    resolve(await response.arrayBuffer());
                });
            }
            runner.postMessage({type: 'buildPlaylist', trackData, additionalTrackData, trackCount: data.trackCount, firstTrackOnly: data.firstTrackOnly, firstTrack: data.firstTrack, script: encoder.encode(`(function(){${decoder.decode(data.script)}})()`), secret}, [trackData, additionalTrackData]);
        }
        else if (data.type === 'loadTracks') {
            tracksReadyResolver(URL.createObjectURL(new Blob([data.tracks], {type: 'application/json'})));
            responsePort.postMessage({type: 'acknowledge'});
        }
        else if (data.type === 'loadAdditionalTrack') {
            additionalTracks[data.track.track.id] = data.track;
            additionalTracksUrl = URL.createObjectURL(new Blob([JSON.stringify(additionalTracks)], {type: 'application/json'}));
            responsePort.postMessage({type: 'acknowledge'});
        }
    }
});
