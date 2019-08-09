import RecordCrateWorker from 'worker-loader!~/assets/recordcrate.worker.js';
import {encoder} from '~/common/textencoding';
import {sendTrackBlobToEos, sendTrackToEos, buildPlaylist} from '~/assets/eos';
import _builders from '~/builders/index';
export const builders = Object.freeze(_builders);

const cacheVersion = 'v1';
const baseTracksUrl = '/tracks.json';
const processedTracksUrl = '/processed-tracks.json';
let initializeCalled = false;

export async function initialize(countryCode) {
    if (initializeCalled) {
        return;
    }
    initializeCalled = true;
    let data = await getProcessedTracks(countryCode);
    await sendTrackBlobToEos(data);
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
        cache = await caches.open(cacheVersion);
        let response = await getFromCache(cache, baseTracksUrl);
        if (response) {
            return response;
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
    return response;
}

async function getProcessedTracks(countryCode) {
    let req = getProcessedTracksRequest(countryCode);
    let cache;
    if ('caches' in window) {
        cache = await caches.open(cacheVersion);
        let response = await getFromCache(cache, req);
        if (response) {
            let data = await response.arrayBuffer();
            return data;
        }
    } else {
        console.log('Browser does not support cache');
    }
    let tracksResponse = await getTracks();
    let expires = tracksResponse.headers.get('expires');
    let tracks = await tracksResponse.arrayBuffer();
    let recordCrateWorker = new RecordCrateWorker();
    let data = await new Promise((resolve, reject) => {
        recordCrateWorker.addEventListener('message', async ({data}) => {
            if (data.type === 'sendProcessedTracks') {
                resolve(data.data);
            }
            else {
                reject();
            }
        });
        recordCrateWorker.postMessage({type: 'sendTracks', tracks, countryCode});
    });
    recordCrateWorker.terminate();
    if (cache && expires) {
        console.log('Caching good processed tracks JSON for next time');
        await cache.put(req, new Response(data, {
            status: 200,
            statusText: 'OK',
            headers: {
                'Vary': 'X-Phosphor-Accept-Region',
                'X-Phosphor-Content-Region': countryCode,
                'Expires': expires,
                'Content-Type': 'application/json',
                'Content-Length': data.length
            }
        }));
    }
    return data;
}

async function getFromCache(cache, request) {
    let url = request.url || request;
    let response = await cache.match(request);
    if (!response) {
        console.log(`${url} not found in cache`);
        return;
    }
    console.log(`${url} found in cache`);
    let expiresHeader = response.headers.get('expires');
    if (!expiresHeader) {
        console.log(`${url} has no expiration header, ignoring`)
        return
    }
    console.log(`Cached ${url} has expiration header:`, expiresHeader);
    let expires = new Date(expiresHeader);
    let now = new Date();
    if (expires > now) {
        console.log(`Cached ${url} expiration header is in the future, using cache:`, expires);
        return response;
    } else {
        console.log(`Cached ${url} has expired, ignoring`);
    }
}

function getProcessedTracksRequest(countryCode) {
    return new Request(processedTracksUrl, {
        headers: new Headers({
            'X-Phosphor-Accept-Region': countryCode
        })
    });
}
