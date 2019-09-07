import {encoder} from '~/common/textencoding';
import SecureMessenger from '~/secure-messenger/secure-messenger';

let messenger;
let terminationPort;

export function initialize() {
    return new Promise(async (resolve, reject) => {
        let iframe = document.createElement('iframe');
        iframe.src = process.env.EOS_ORIGIN;
        iframe.sandbox = 'allow-scripts allow-same-origin';
        iframe.style.display = 'none';
        iframe.addEventListener('load', () => {
            resolve();
        });
        messenger = new SecureMessenger(process.env.EOS_ORIGIN);
        let ready = messenger.listen(iframe);
        document.body.appendChild(iframe);
        await ready;
        let {type, interruptPort} = await messenger.postMessage({type: 'requestTerminationChannel'});
        if (type !== 'terminationChannel') {
            reject('Could not open termination channel');
        }
        terminationPort = interruptPort;
    });
}

export async function sendTrackBlobToEos(raw) {
    let {data} = await messenger.postMessage({type: 'loadTracks', tracks: raw});
    if (data.type === 'acknowledge') {
        return;
    }
    else if (data.type === 'error') {
        throw new Error(data.error);
    }
    throw new Error(`Unknown error when sending track data: ${data.type}`);
}

export async function sendTrackToEos(track) {
    let {data} = await messenger.postMessage({type: 'loadAdditionalTrack', track});
    if (data.type === 'acknowledge') {
        return;
    }
    else if (data.type === 'error') {
        throw new Error(data.error);
    }
    throw new Error(`Unknown error when sending track data: ${data.type}`);
}

export async function buildPlaylist(trackCount, builder, firstTrackBuilder, firstTrack, pruners, loadPercent) {
    let allDimensions = [];
    let prunedTrackIds;
    function appendDimensions(newDims) {
        if (!newDims) {
            return;
        }
        allDimensions = [...allDimensions, ...newDims].filter((val, index, arr) => arr.indexOf(val) === index);
    }
    let totalLoaders = pruners.length + 1;
    loadPercent(0.1);
    if (pruners) {
        let i = 0;
        for (let pruner of pruners) {
            let script = encoder.encode(pruner);
            let response = await callBuilder({prunedTrackIds, script}, 'pruneTracks', percent => {
                loadPercent(0.1 + (0.9 * ((percent / totalLoaders) + (i / totalLoaders))));
            });
            if (!prunedTrackIds) {
                prunedTrackIds = [];
            }
            prunedTrackIds = [...prunedTrackIds, ...response.prunedTrackIds];
            appendDimensions(response.dimensions);
            i++;
        }
    }
    let script = encoder.encode(builder);
    if (firstTrackBuilder) {
        let response = await callBuilder({firstTrackOnly: true, prunedTrackIds, script: encoder.encode(firstTrackBuilder)}, 'buildPlaylist', () => {
            loadPercent(0.1 + (0.9 * (((1 / trackCount) / totalLoaders) + ((totalLoaders - 1) / totalLoaders))));
        });
        let firstTrack = response.playlist[0];
        appendDimensions(response.dimensions);
        let {playlist, dimensions} = await callBuilder({firstTrack, trackCount, prunedTrackIds, script}, 'buildPlaylist', percent => {
            loadPercent(0.1 + (0.9 * ((percent / totalLoaders) + ((totalLoaders - 1) / totalLoaders))));
        });
        appendDimensions(dimensions);
        return {
            playlist,
            dimensions: allDimensions
        };
    }
    else if (firstTrack) {
        let {playlist, dimensions} = await callBuilder({firstTrack, trackCount, prunedTrackIds, script}, 'buildPlaylist', percent => {
            loadPercent(0.1 + (0.9 * ((percent / totalLoaders) + ((totalLoaders - 1) / totalLoaders))));
        });
        appendDimensions(dimensions);
        return {
            playlist,
            dimensions: allDimensions
        };
    }
    else {
        let {playlist, dimensions} = await callBuilder({trackCount, prunedTrackIds, script}, 'buildPlaylist', percent => {
            loadPercent(0.1 + (0.9 * ((percent / totalLoaders) + ((totalLoaders - 1) / totalLoaders))));
        });
        appendDimensions(dimensions);
        return {
            playlist,
            dimensions: allDimensions
        };
    }
}

export function terminatePlaylistBuilding() {
    terminationPort.postMessage({type: 'terminateAll'});
}

async function callBuilder(body, type, loadPercentHandler) {
    let {closer} = await messenger.openInterruptListenerPort({type: 'initializeLoadingNotificationChannel'}, ({type, value}) => {
        if (type == 'loadPercent') {
            loadPercentHandler && loadPercentHandler(value);
        }
    });
    let {data} = await messenger.postMessage(Object.assign({type}, body));
    closer();
    if (data.type === 'playlist') {
        return {playlist: data.playlist, dimensions: data.dimensions};
    }
    else if (data.type === 'prunedTracks') {
        return {prunedTrackIds: data.prunedTrackIds, dimensions: data.dimensions};
    }
    else if (data.type === 'error') {
        throw new Error(data.error);
    }
    throw new Error(`Unknown error: ${data.type}`);
}
