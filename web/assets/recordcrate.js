import RecordCrateWorker from 'worker-loader!~/assets/recordcrate.worker.js';
import {encoder} from '~/common/textencoding';
import {getAccessToken} from '~/assets/session';
import {sendTrackBlobToEos, buildPlaylist} from '~/assets/eos';
import _builders from '~/builders/index';
export const builders = Object.freeze(_builders);

let initializeCalled = false;

export async function initialize(countryCode) {
    if (initializeCalled) {
        return;
    }
    initializeCalled = true;
    let response = await fetch(`/api/spotify/tracks?token=${await getAccessToken()}`);
    let encodedTracks = await response.arrayBuffer();
    let recordCrateWorker = new RecordCrateWorker();
    await new Promise((resolve, reject) => {
        recordCrateWorker.addEventListener('message', async ({data}) => {
            if (data.type === 'sendProcessedTracks') {
                await sendTrackBlobToEos(data.data);
                resolve();
            }
            else {
                reject();
            }
        });
        recordCrateWorker.postMessage({type: 'sendTracks', tracks: encodedTracks, countryCode});
    });
    recordCrateWorker.terminate();
}

export async function loadNewPlaylist(count, builder, firstTrackBuilder) {
    if (!builder) {
        builder = builders.randomwalk;
    }
    if (!count) {
        count = 10;
    }
    let playlist;
    let dimensions;
    try {
        let response = await buildPlaylist(count, builder, firstTrackBuilder);
        playlist = response.playlist;
        dimensions = response.dimensions;
    }
    catch (e) {
        console.error('Could not build playlist:', e);
        return {error: e};
    }
    return {playlist, dimensions};
}
