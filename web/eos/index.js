import RunnerWorker from 'worker-loader!./runner.worker.js';
import {encoder, decoder} from '../common/textencoding';

let tracksReadyResolver;
let tracksReady = new Promise(resolve => tracksReadyResolver = resolve);

addEventListener('message', async ({origin, data}) => {
    if (origin === process.env.PHOSPHOR_ORIGIN) {
        let {responsePort} = data;
        if (data.type === 'buildPlaylist') {
            let finished = false;
            let runner = new RunnerWorker();
            let secret = crypto.getRandomValues(new Uint8Array(32)).reduce((str, byte) => str + byte.toString(16).padStart(2, '0'), '');
            runner.addEventListener('message', e => {
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
                runner.terminate();
                finished = true;
            }, {once: true});
            addEventListener('message', ({origin, data}) => {
                if (finished) {
                    return;
                }
                if (origin === process.env.PHOSPHOR_ORIGIN && data.type === 'terminateAll') {
                    responsePort.postMessage({type: 'playlistError', error: 'User killed builder'});
                    runner.terminate();
                    finished = true;
                }
            }, {once: true});
            runner.postMessage({type: 'buildPlaylist', tracksUrl: await tracksReady, trackCount: data.trackCount, firstTrackOnly: data.firstTrackOnly, firstTrack: data.firstTrack, script: encoder.encode(`(function(){${decoder.decode(data.script)}})()`), secret});
        }
        else if (data.type === 'loadTracks') {
            tracksReadyResolver(URL.createObjectURL(new Blob([data.tracks], {type: 'application/json'})));
            responsePort.postMessage({type: 'acknowledge'});
        }
    }
});
