import {encoder} from '~/common/textencoding';

let iframe;

export function initialize() {
    return new Promise(resolve => {
        iframe = document.createElement('iframe');
        iframe.src = process.env.EOS_ORIGIN;
        iframe.sandbox = 'allow-scripts allow-same-origin';
        iframe.style.display = 'none';
        iframe.addEventListener('load', () => {
            resolve();
        });
        document.body.appendChild(iframe);
    });
}

export function sendTrackBlobToEos(raw) {
    return new Promise((resolve, reject) => {
        let channel = new MessageChannel();
        channel.port1.onmessage = ({origin, data}) => {
            if (data.type === 'acknowledge') {
                resolve();
            }
            else if (data.type === 'error') {
                reject(data.error);
            }
            else {
                reject(`Unknown error when sending track data: ${data.type}`)
            }
            channel.port1.close();
        };
        iframe.contentWindow.postMessage({type: 'loadTracks', tracks: raw, responsePort: channel.port2}, process.env.EOS_ORIGIN, [channel.port2]);
    });
}

export function sendTrackToEos(track) {
    return new Promise((resolve, reject) => {
        let channel = new MessageChannel();
        channel.port1.onmessage = ({origin, data}) => {
            if (data.type === 'acknowledge') {
                resolve();
            }
            else if (data.type === 'error') {
                reject(data.error);
            }
            else {
                reject(`Unknown error when sending track data: ${data.type}`)
            }
            channel.port1.close();
        };
        iframe.contentWindow.postMessage({type: 'loadAdditionalTrack', track, responsePort: channel.port2}, process.env.EOS_ORIGIN, [channel.port2]);
    });
}

export async function buildPlaylist(trackCount, builder, firstTrackBuilder, firstTrack, pruners) {
    let allDimensions = [];
    let prunedTrackIds;
    function appendDimensions(newDims) {
        if (!newDims) {
            return;
        }
        allDimensions = [...allDimensions, ...newDims].filter((val, index, arr) => arr.indexOf(val) === index);
    }
    if (pruners) {
        for (let pruner of pruners) {
            let script = encoder.encode(pruner);
            let response = await callBuilder({prunedTrackIds, script}, 'pruneTracks');
            if (!prunedTrackIds) {
                prunedTrackIds = [];
            }
            prunedTrackIds = [...prunedTrackIds, ...response.prunedTrackIds];
            appendDimensions(response.dimensions);
        }
    }
    let script = encoder.encode(builder);
    if (firstTrackBuilder) {
        let response = await callBuilder({firstTrackOnly: true, prunedTrackIds, script: encoder.encode(firstTrackBuilder)}, 'buildPlaylist');
        let firstTrack = response.playlist[0];
        appendDimensions(response.dimensions);
        let {playlist, dimensions} = await callBuilder({firstTrack, trackCount, prunedTrackIds, script}, 'buildPlaylist');
        appendDimensions(dimensions);
        return {
            playlist,
            dimensions: allDimensions
        };
    }
    else if (firstTrack) {
        let {playlist, dimensions} = await callBuilder({firstTrack, trackCount, prunedTrackIds, script}, 'buildPlaylist');
        appendDimensions(dimensions);
        return {
            playlist,
            dimensions: allDimensions
        };
    }
    else {
        let {playlist, dimensions} = await callBuilder({trackCount, prunedTrackIds, script}, 'buildPlaylist');
        appendDimensions(dimensions);
        return {
            playlist,
            dimensions: allDimensions
        };
    }
}

export function terminatePlaylistBuilding() {
    iframe.contentWindow.postMessage({type: 'terminateAll'}, process.env.EOS_ORIGIN);
}

function callBuilder(body, type) {
    return new Promise((resolve, reject) => {
        let channel = new MessageChannel();
        channel.port1.onmessage = ({origin, data}) => {
            if (data.type === 'playlist') {
                resolve({playlist: data.playlist, dimensions: data.dimensions});
            }
            else if (data.type === 'prunedTracks') {
                resolve({prunedTrackIds: data.prunedTrackIds, dimensions: data.dimensions});
            }
            else if (data.type === 'playlistError') {
                reject(data.error);
            }
            else {
                reject(`Unknown error when building playlist: ${data.type}`)
            }
            channel.port1.close();
        };
        iframe.contentWindow.postMessage(Object.assign({type, responsePort: channel.port2}, body), process.env.EOS_ORIGIN, [channel.port2]);
    });
}
