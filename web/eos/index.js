import RunnerWorker from 'worker-loader!./runner.worker.js';
import SecureMessenger from '../secure-messenger/secure-messenger';
import {TextEncoder, TextDecoder} from 'text-encoding-shim';
import pako from 'pako';

const encoder = new TextEncoder();
const decoder = new TextDecoder();

let additionalTracks = {};
let tracksReadyResolver;
let tracksReady = new Promise(resolve => tracksReadyResolver = resolve);
let phosphorMessenger;
const terminationEvent = 'eosTerminateAll';

(async () => {
    phosphorMessenger = new SecureMessenger(process.env.PHOSPHOR_ORIGIN);
    await phosphorMessenger.knock([parent]);
    let loadingInterruptPort;
    phosphorMessenger.messageHandlerLoop((data, interruptPort) => {
        if (data.type == 'buildPlaylist') {
            return handleBuildPlaylist(loadingInterruptPort, data);
        } else if (data.type == 'pruneTracks') {
            return handlePruneTracks(loadingInterruptPort, data);
        } else if (data.type == 'loadTracks') {
            return handleLoadTracks(data);
        } else if (data.type == 'loadAdditionalTrack') {
            return handleLoadAdditionalTrack(data);
        } else if (data.type == 'initializeLoadingNotificationChannel') {
            loadingInterruptPort = interruptPort;
            return {type: 'acknowledge'};
        } else if (data.type == 'requestTerminationChannel') {
            return phosphorMessenger.respondWithInterruptListenerPort({type: 'terminationChannel'}, ({type}) => {
                if (type === 'terminateAll') {
                    console.debug('Should dispatch termination event...');
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

async function handleBuildPlaylist(loadingInterruptPort, data) {
    let response = await buildRunner(loadingInterruptPort)({
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

async function handlePruneTracks(loadingInterruptPort, data) {
    let response = await buildRunner(loadingInterruptPort)({
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

function buildRunner(loadingInterruptPort) {
    let workerMessenger = new SecureMessenger(location.origin);
    let runner = new RunnerWorker();
    let terminationHandler;
    let terminationPromise = new Promise(resolve => {
        terminationHandler = () => {
            console.debug("Should terminate...");
            runner.terminate();
            resolve({type: 'runnerError', error: 'User killed builder'});
        };
        addEventListener(terminationEvent, terminationHandler);
    });
    return (postData) => {
        return Promise.race([terminationPromise, new Promise(async resolve => {
            await workerMessenger.listen(runner);
            await workerMessenger.openInterruptListenerPort({type: 'initializeLoadingNotificationChannel'}, ({type, value}) => {
                if (type == 'loadPercent') {
                    loadingInterruptPort.postMessage({type, value});
                }
            });
            let {data} = await workerMessenger.postMessage(postData);
            runner.terminate();
            resolve(data);
        })]).then(passThru => {
            removeEventListener(terminationEvent, terminationHandler);
            return passThru;
        });
    };
}
