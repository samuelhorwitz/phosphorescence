import RunnerWorker from 'worker-loader!./runner.worker.js';
import {encoder, decoder} from '../common/textencoding';

let additionalTracks = {};
let tracksReadyResolver;
let tracksReady = new Promise(resolve => tracksReadyResolver = resolve);

async function cleanup() {
    // Since we cannot completely sandbox this origin, users may create Indexed DBs or caches.
    // We don't want them to do that and will destroy them.
    let indexedDB = window.indexedDB || window.mozIndexedDB || window.webkitIndexedDB || window.msIndexedDB;
    if (indexedDB && indexedDB.databases) {
        (await indexedDB.databases()).forEach(db => indexedDB.deleteDatabase(db.name));
    }
    if (caches) {
        (await caches.keys()).forEach(cache => caches.delete(cache));
    }
}

addEventListener('message', async ({origin, data}) => {
    if (origin !== process.env.PHOSPHOR_ORIGIN) {
        return;
    }
    let {responsePort} = data;
    if (data.type === 'buildPlaylist') {
        await cleanup();
        let finished = false;
        let runner = new RunnerWorker();
        async function killRunnerAndCleanup() {
            runner.terminate();
            finished = true;
            await cleanup();
        }
        let secret = crypto.getRandomValues(new Uint8Array(32)).reduce((str, byte) => str + byte.toString(16).padStart(2, '0'), '');
        runner.addEventListener('message', async e => {
            if (finished) {
                return;
            }
            await killRunnerAndCleanup();
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
        }, {once: true});
        addEventListener('message', async ({origin, data}) => {
            if (finished) {
                return;
            }
            if (origin === process.env.PHOSPHOR_ORIGIN && data.type === 'terminateAll') {
                await killRunnerAndCleanup();
                responsePort.postMessage({type: 'playlistError', error: 'User killed builder'});
            }
        }, {once: true});
        runner.postMessage({
            type: 'buildPlaylist',
            tracksUrl: await tracksReady,
            additionalTracksUrl: URL.createObjectURL(new Blob([JSON.stringify(additionalTracks)], {type: 'application/json'})),
            trackCount: data.trackCount,
            firstTrackOnly: data.firstTrackOnly,
            firstTrack: data.firstTrack,
            script: encoder.encode(`(function(){${decoder.decode(data.script)}})()`),
            secret
        });
    }
    else if (data.type === 'loadTracks') {
        tracksReadyResolver(URL.createObjectURL(new Blob([data.tracks], {type: 'application/json'})));
        responsePort.postMessage({type: 'acknowledge'});
    }
    else if (data.type === 'loadAdditionalTrack') {
        additionalTracks[data.track.track.id] = data.track;
        responsePort.postMessage({type: 'acknowledge'});
    }
});
