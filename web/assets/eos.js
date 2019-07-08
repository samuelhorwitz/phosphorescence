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
        iframe.contentWindow.postMessage({type: 'loadTracks', tracks: raw, responsePort: channel.port2}, '*', [channel.port2]);
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
        iframe.contentWindow.postMessage({type: 'loadAdditionalTrack', track, responsePort: channel.port2}, '*', [channel.port2]);
    });
}

export async function buildPlaylist(trackCount, builder, firstTrackBuilder, firstTrack) {
    let script = encoder.encode(builder);
    if (firstTrackBuilder) {
        let firstTrack, firstTrackDimensions;
        try {
            let response = await callBuilder({firstTrackOnly: true, script: encoder.encode(firstTrackBuilder)});
            firstTrack = response.playlist[0];
            firstTrackDimensions = response.dimensions;
        }
        catch (e) {
            console.error(e);
            return null;
        }
        let {playlist, dimensions} = await callBuilder({firstTrack, trackCount, script});
        return {
            playlist,
            dimensions: [...firstTrackDimensions, ...dimensions].filter((val, index, arr) => arr.indexOf(val) === index)
        };
    }
    else if (firstTrack) {
        return callBuilder({firstTrack, trackCount, script});
    }
    else {
        return callBuilder({trackCount, script});
    }
}

export function terminatePlaylistBuilding() {
    iframe.contentWindow.postMessage({type: 'terminateAll'}, '*');
}

function callBuilder(body) {
    return new Promise((resolve, reject) => {
        let channel = new MessageChannel();
        channel.port1.onmessage = ({origin, data}) => {
            if (data.type === 'playlist') {
                resolve({playlist: data.playlist, dimensions: data.dimensions});
            }
            else if (data.type === 'playlistError') {
                reject(data.error);
            }
            else {
                reject(`Unknown error when building playlist: ${data.type}`)
            }
            channel.port1.close();
        };
        iframe.contentWindow.postMessage(Object.assign({type: 'buildPlaylist', responsePort: channel.port2}, body), '*', [channel.port2]);
    });
}
