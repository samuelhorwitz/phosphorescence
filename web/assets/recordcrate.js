import RecordCrateWorker from 'worker-loader?inline!~/assets/recordcrate.worker.js';
import {sendTrackBlobToEos, sendTrackToEos, buildPlaylist} from '~/assets/eos';
import _builders from '~/builders/index';
import SecureMessenger from '~/secure-messenger/secure-messenger';
import {getCaptchaToken} from '~/assets/captcha';
import throttle from 'lodash/throttle';
export const builders = Object.freeze(_builders);

const cacheVersion = 'v2';
const cachesToDelete = ['v1']
const processedTracksUrl = '/processed-tracks.json';
let initializeCalled = false;

function getBaseTracksUrl(region) {
    return `tracks.${region.toLowerCase()}.json`;
}

export async function initialize(countryCode, isLoggedIn, loadingHandler) {
    if (initializeCalled) {
        return;
    }
    if ('caches' in window) {
        for (let toDelete of cachesToDelete) {
            await caches.delete(toDelete);
        }
    }
    initializeCalled = true;
    let data = await getProcessedTracks(countryCode, isLoggedIn, loadingHandler);
    await sendTrackBlobToEos(data);
}

export async function processTrack(track) {
    let recordCrateWorker = new RecordCrateWorker();
    let messenger = new SecureMessenger(location.origin);
    await messenger.listen(recordCrateWorker);
    let {data} = await messenger.postMessage({type: 'sendTrack', track});
    if (data.type === 'sendProcessedTrack') {
        await sendTrackToEos(data.data);
    }
    messenger.close();
    recordCrateWorker.terminate();
    return data.data;
}

export async function processTracks(tracks) {
    let recordCrateWorker = new RecordCrateWorker();
    let messenger = new SecureMessenger(location.origin);
    await messenger.listen(recordCrateWorker);
    let {data} = await messenger.postMessage({type: 'sendTracks', tracks, raw: true});
    if (data.type !== 'sendProcessedTracks') {
        throw new Error('Invalid response from track processing', data.type);
    }
    messenger.close();
    recordCrateWorker.terminate();
    return data.data;
}

export function loadNewPlaylist({count, builder, firstTrackBuilder, firstTrack, replacementTracks, pruners}, loadPercent) {
    if (!builder) {
        builder = builders.randomwalk;
    }
    if (!count) {
        count = 10;
    }
    return buildPlaylist({count, builder, firstTrackBuilder, firstTrack, replacementTracks, pruners}, loadPercent);
}

async function getTracks(isLoggedIn, region) {
    let baseTracksUrl = getBaseTracksUrl(region);
    let cache;
    if ('caches' in window) {
        cache = await caches.open(cacheVersion);
        let response = await getFromCache(cache, baseTracksUrl);
        if (response) {
            return response;
        }
    } else {
        console.warn('Browser does not support cache');
    }
    let tracksUrlResponse;
    if (isLoggedIn) {
        tracksUrlResponse = await fetch(`${process.env.API_ORIGIN}/spotify/tracks`, {credentials: 'include'});
    } else {
        let captcha = await getCaptchaToken('api/tracks');
        tracksUrlResponse = await fetch(`${process.env.API_ORIGIN}/spotify/unauthenticated/tracks/${region}?captcha=${captcha}`);
    }
    let {tracksUrl} = await tracksUrlResponse.json();
    let response = await fetch(tracksUrl);
    if (cache && response.ok && response.headers.get('expires')) {
        console.log('Caching good tracks JSON for next time');
        try {
            await cache.put(baseTracksUrl, response.clone());
        } catch (e) {
            console.warn('Could not cache tracks JSON', e);
        }
    }
    return response;
}

async function getProcessedTracks(countryCode, isLoggedIn, loadingHandler) {
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
        console.warn('Browser does not support cache');
    }
    let tracksResponse = await getTracks(isLoggedIn, countryCode);
    let expires = tracksResponse.headers.get('expires');
    let tracks = await getArrayBufferWithProgress(tracksResponse, percent => {
        loadingHandler(percent * 0.45);
    });
    let recordCrateWorker = new RecordCrateWorker();
    let messenger = new SecureMessenger(location.origin);
    await messenger.listen(recordCrateWorker);
    await messenger.openInterruptListenerPort({type: 'initializeLoadingNotificationChannel'}, ({type, value}) => {
        if (type == 'loadPercent') {
            loadingHandler(0.45 + (value * 0.45), true);
        }
    });
    let {data} = await messenger.postMessage({type: 'sendTracks', tracks});
    let gzipData;
    if (data.type === 'sendProcessedTracks') {
        gzipData = data.gzipData;
    }
    messenger.close();
    recordCrateWorker.terminate();
    if (cache && expires) {
        console.log('Caching good processed tracks JSON for next time');
        try {
            await cache.put(req, new Response(gzipData, {
                status: 200,
                statusText: 'OK',
                headers: {
                    'Vary': 'X-Phosphor-Accept-Region, Accept-Encoding',
                    'X-Phosphor-Content-Region': countryCode,
                    'Expires': expires,
                    'Content-Type': 'application/json',
                    'Content-Encoding': 'gzip',
                    'Content-Length': gzipData.byteLength
                }
            }));
        } catch (e) {
            console.warn('Could not cache processed tracks JSON', e);
        }
    }
    return gzipData;
}

// https://javascript.info/fetch-progress
async function getArrayBufferWithProgress(response, progressHandler) {
    progressHandler = throttle(progressHandler, 100);
    let reader = response.body.getReader();
    let contentLength = +(response.headers.get('x-amz-meta-uncompressed-length') || response.headers.get('Content-Length'));
    let receivedLength = 0;
    let chunks = [];
    while (true) {
        let {done, value} = await reader.read();
        if (done) {
            break;
        }
        chunks.push(value);
        receivedLength += value.length;
        progressHandler(receivedLength / contentLength);
    }
    let chunksAll = new Uint8Array(receivedLength);
    let position = 0;
    for(let chunk of chunks) {
        chunksAll.set(chunk, position);
        position += chunk.length;
    }
    return chunksAll.buffer;
}

async function getFromCache(cache, request) {
    let url = request.url || request;
    let response = await cache.match(request);
    if (!response) {
        console.info(`${url} not found in cache`);
        return;
    }
    console.info(`${url} found in cache`);
    let expiresHeader = response.headers.get('expires');
    if (!expiresHeader) {
        console.info(`${url} has no expiration header, ignoring`)
        return
    }
    console.info(`Cached ${url} has expiration header:`, expiresHeader);
    let expires = new Date(expiresHeader);
    let now = new Date();
    if (expires > now) {
        console.info(`Cached ${url} expiration header is in the future, using cache:`, expires);
        return response;
    } else {
        console.info(`Cached ${url} has expired, ignoring`);
    }
}

function getProcessedTracksRequest(countryCode) {
    return new Request(processedTracksUrl, {
        headers: new Headers({
            'Accept-Encoding': 'gzip',
            'X-Phosphor-Accept-Region': countryCode
        })
    });
}
