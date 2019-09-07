import RunnerWorker from 'worker-loader!./runner.worker.js';
import SecureMessenger from '../secure-messenger/secure-messenger';
import {encoder, decoder} from '../common/textencoding';
import pako from 'pako';

let additionalTracks = {};
let tracksReadyResolver;
let tracksReady = new Promise(resolve => tracksReadyResolver = resolve);
let loadingInterruptPort;
let phosphorMessenger;
const terminationEvent = 'eosTerminateAll';

(async () => {
    phosphorMessenger = new SecureMessenger(process.env.PHOSPHOR_ORIGIN);
    await phosphorMessenger.knock([parent]);
    phosphorMessenger.messageHandlerLoop((data, interruptPort) => {
        if (data.type == 'buildPlaylist') {
            return handleBuildPlaylist(data);
        } else if (data.type == 'pruneTracks') {
            return handlePruneTracks(data);
        } else if (data.type == 'loadTracks') {
            return handleLoadTracks(data);
        } else if (data.type == 'loadAdditionalTrack') {
            return handleLoadAdditionalTrack(data);
        } else if (data.type == 'initializeLoadingNotificationChannel') {
            loadingInterruptPort = interruptPort;
            return {type: 'acknowledge'};
        } else if (data.type == 'requestTerminationChannel') {
            return phosphorMessenger.respondWithInterruptListenerPort('terminationChannel', ({type}) => {
                if (type === 'terminateAll') {
                    dispatchEvent(new Event(terminationEvent));
                }
            }, event => {
                console.error('Could not handle termination channel message', event);
            });
        }
        return {type: 'error', error: 'invalid request'};
    }, event => {
        console.error('Eos could not handle message from Phosphor', event);
    });
})();

async function handleBuildPlaylist(data) {
    let response = await (await buildRunner())({
        type: 'buildPlaylist',
        tracksUrl: await tracksReady,
        additionalTracksUrl: URL.createObjectURL(new Blob([JSON.stringify(additionalTracks)], {type: 'application/json'})),
        prunedTrackIds: data.prunedTrackIds,
        trackCount: data.trackCount,
        firstTrackOnly: data.firstTrackOnly,
        firstTrack: data.firstTrack,
        script: encoder.encode(`(function(){${decoder.decode(data.script)}})()`)
    });
    if (response.type === 'playlist') {
        return {type: 'playlist', playlist: response.playlist, dimensions: response.dimensions};
    }
    return {type: 'error', error: response.error || 'Unknown error'};
}

async function handlePruneTracks(data) {
    let response = await (await buildRunner())({
        type: 'pruneTracks',
        tracksUrl: await tracksReady,
        additionalTracksUrl: URL.createObjectURL(new Blob([JSON.stringify(additionalTracks)], {type: 'application/json'})),
        prunedTrackIds: data.prunedTrackIds,
        script: encoder.encode(`(function(){${decoder.decode(data.script)}})()`)
    });
    if (response.type === 'prunedTracks') {
        return {type: 'prunedTracks', prunedTrackIds: response.prunedTrackIds, dimensions: response.dimensions};
    }
    return {type: 'error', error: response.error || 'Unknown error'};
}

function handleLoadTracks(data) {
    let trackData = pako.ungzip(data.tracks);
    tracksReadyResolver(URL.createObjectURL(new Blob([trackData], {type: 'application/json'})));
    return {type: 'acknowledge'};
}

function handleLoadAdditionalTrack(data) {
    additionalTracks[data.track.track.id] = data.track;
    return {type: 'acknowledge'};
}

async function buildRunner() {
    let workerMessenger = new SecureMessenger(location.origin);
    let runner = new RunnerWorker();
    let terminationPromise = new Promise(resolve => {
        addEventListener(terminationEvent, async () => {
            runner.terminate();
            resolve({type: 'runnerError', error: 'User killed builder'});
        }, {once: true});
    });
    await workerMessenger.listen(runner);
    await workerMessenger.openInterruptListenerPort({type: 'initializeLoadingNotificationChannel'}, ({type, value}) => {
        if (type == 'loadPercent') {
            loadingInterruptPort.postMessage({type, value});
        }
    });
    return (postData) => {
        return Promise.race([terminationPromise, new Promise(async resolve => {
            let {data} = await workerMessenger.postMessage(postData);
            runner.terminate();
            resolve(data);
        })]);
    };
}
