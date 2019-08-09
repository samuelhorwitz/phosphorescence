import RecordCrateWorker from 'worker-loader!~/assets/recordcrate.worker.js';
import {encoder} from '~/common/textencoding';
import {sendTrackBlobToEos, sendTrackToEos, buildPlaylist} from '~/assets/eos';
import _builders from '~/builders/index';
export const builders = Object.freeze(_builders);

const cacheVersion = 'v1';
const baseTracksUrl = 'https://phosphorescence.sfo2.digitaloceanspaces.com/tracks.json';
let initializeCalled = false;

export async function initialize(countryCode) {
    if (initializeCalled) {
        return;
    }
    initializeCalled = true;
    let encodedTracks = await getTracks();
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

export async function processTrack(countryCode, track) {
    let recordCrateWorker = new RecordCrateWorker();
    let processedTrack = await new Promise((resolve, reject) => {
        recordCrateWorker.addEventListener('message', async ({data}) => {
            if (data.type === 'sendProcessedTrack') {
                await sendTrackToEos(data.data);
                resolve(data.data);
            }
            else {
                reject();
            }
        });
        recordCrateWorker.postMessage({type: 'sendTrack', track, countryCode});
    });
    recordCrateWorker.terminate();
    return processedTrack;
}

export async function loadNewPlaylist(count, builder, firstTrackBuilder, firstTrack) {
    if (!builder) {
        builder = builders.randomwalk;
    }
    if (!count) {
        count = 10;
    }
    let playlist;
    let dimensions;
    try {
        let response = await buildPlaylist(count, builder, firstTrackBuilder, firstTrack);
        playlist = response.playlist;
        dimensions = response.dimensions;
    }
    catch (e) {
        console.error('Could not build playlist:', e);
        return {error: e};
    }
    return {playlist, dimensions};
}

async function getTracks() {
    let cache;
    if ('caches' in window) {
        console.log('Browser supports cache');
        cache = await caches.open(cacheVersion);
        let response = await cache.match(baseTracksUrl);
        if (response) {
            console.log('Tracks JSON found in cache');
            let expiresHeader = response.headers.get('expires');
            if (expiresHeader) {
                console.log('Cached tracks JSON has expiration header:', expiresHeader);
                let expires = new Date(expiresHeader);
                let now = new Date();
                if (expires > now) {
                    console.log('Cached tracks JSON expiration header is in the future, using cache:', expires);
                    return response.arrayBuffer();
                } else {
                    console.log('Cached tracks JSON has expired, ignoring');
                }
            } else {
                console.log('Tracks JSON has no expiration header, ignoring')
            }
        } else {
            console.log('Tracks JSON not found in cache');
        }
    } else {
        console.log('Browser does not support cache');
    }
    let tracksUrlResponse = await fetch(`${process.env.API_ORIGIN}/spotify/tracks`, {credentials: 'include'});
    let {tracksUrl} = await tracksUrlResponse.json();
    let response = await fetch(tracksUrl);
    if (cache && response.ok && response.headers.get('expires')) {
        console.log('Caching good tracks JSON for next time');
        await cache.put(baseTracksUrl, response.clone());
    }
    return response.arrayBuffer();
}
