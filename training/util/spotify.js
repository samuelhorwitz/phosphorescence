import * as tf from '@tensorflow/tfjs';

let store, aModel, pModel, meanstd;
let storeReadyResolver;
let modelsReadyResolver;
let storeReady = new Promise(resolve => storeReadyResolver = resolve);
let modelsReady = new Promise(resolve => modelsReadyResolver = resolve);
if (process.browser) {
    window.onNuxtReady(async ({$store}) => {
        store = $store;
        storeReadyResolver();
        aModel = await tf.loadModel('/models/aetherealness/model.json');
        pModel = await tf.loadModel('/models/primordialness/model.json');
        let meanstdRaw = await fetch('/models/meanstd.json');
        meanstd = await meanstdRaw.json();
        modelsReadyResolver();
    });
}

export async function getTokens(code, refresh) {
    let type = 'type=auth';
    if (refresh) {
        type = 'type=refresh';
    }
    let response = await fetch(`/api/spotify/tokens?code=${code}&${type}`, {cache: 'no-cache'});
    let {status, statusText} = response;
    if (status != 200) {
        throw new Error(`Could not get tokens: ${statusText}`);
    }
    return response.json();
}

export function saveLocalTokens({access, refresh, expires}) {
    let now = new Date();
    sessionStorage.setItem('accessToken', access);
    sessionStorage.setItem('refreshToken', refresh);
    sessionStorage.setItem('expires', expires);
    sessionStorage.setItem('expiresAt', now.setSeconds(now.getSeconds() + expires));
}

export function getLocalTokens() {
    let access = sessionStorage.getItem('accessToken');
    let refresh = sessionStorage.getItem('refreshToken');
    let expires_ = sessionStorage.getItem('expires');
    let expiresAt_ = sessionStorage.getItem('expiresAt');
    let expires = parseInt(expires_, 10);
    let expiresAt = parseInt(expiresAt_, 10);
    return {access, refresh, expires, expiresAt, okay: areLocalTokensOkay({access, refresh, expires, expiresAt})};
}

export function areLocalTokensOkay({access, refresh, expires, expiresAt}) {
    return !!(access && refresh && expires && expiresAt && Date.now() < expiresAt);
}

export function sendToAuth() {
    sessionStorage.clear();
    location.href = '/api/spotify/authorize';
}

export function tryToRefresh(refresh) {
    if (refresh) {
        try {
            return getTokens(refresh, true);
        }
        catch (e) {
            sendToAuth();
        }
    } else {
        sendToAuth();
    }
}

export function shuffle(array) {
    let currentIndex = array.length, temporaryValue, randomIndex;
    // While there remain elements to shuffle...
    while (0 !== currentIndex) {
        // Pick a remaining element...
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex -= 1;

        // And swap it with the current element.
        temporaryValue = array[currentIndex];
        array[currentIndex] = array[randomIndex];
        array[randomIndex] = temporaryValue;
    }
    return array;
}

export async function getSomeTrance() {
    await storeReady;
    let response = await fetch('https://api.spotify.com/v1/recommendations?seed_genres=trance&limit=100', {
        headers: {
            Authorization: `Bearer ${store.state.session.accessToken}`
        }
    });
    let tracks = await response.json();
    let tracksInPlaylistStyle = tracks.tracks.map(t => {return {track: t}});
    let {audio_features} = await getAudioFeatures(tracksInPlaylistStyle);
    return {tracks: tracksInPlaylistStyle, allTracksAnalysis: audio_features};
}

export async function getPlaylist(id) {
    let nextUrl = id;
    let isNext = false;
    let tracks = [];
    let allTracksAnalysis = [];
    while (true) {
        let {items, next} = await getPlaylistPage(nextUrl, isNext);
        isNext = true;
        tracks = tracks.concat(items);
        nextUrl = next;
        let {audio_features} = await getAudioFeatures(items);
        allTracksAnalysis = allTracksAnalysis.concat(audio_features);
        if (!nextUrl) {
            break;
        }
    }
    return {tracks, allTracksAnalysis};
}

async function getPlaylistPage(idOrUrl, isNext) {
    await storeReady;
    let url;
    if (isNext) {
        url = idOrUrl;
    } else {
        url = `https://api.spotify.com/v1/playlists/${idOrUrl}/tracks`;
    }
    let response = await fetch(url, {
        headers: {
            Authorization: `Bearer ${store.state.session.accessToken}`
        }
    });
    return response.json();
}

async function getAudioFeatures(tracks) {
    await storeReady;
    let response = await fetch(`https://api.spotify.com/v1/audio-features?ids=${tracks.map(t => t.track.id).join(',')}`, {
        headers: {
            Authorization: `Bearer ${store.state.session.accessToken}`
        }
    });
    return response.json();
}

export async function predict(a) {
    let keyC = a.key == 0 ? 1 : 0;
    let keyCs = a.key == 1 ? 1 : 0;
    let keyD = a.key == 2 ? 1 : 0;
    let keyDs = a.key == 3 ? 1 : 0;
    let keyE = a.key == 4 ? 1 : 0;
    let keyF = a.key == 5 ? 1 : 0;
    let keyFs = a.key == 6 ? 1 : 0;
    let keyG = a.key == 7 ? 1 : 0;
    let keyGs = a.key == 8 ? 1 : 0;
    let keyA = a.key == 9 ? 1 : 0;
    let keyAs = a.key == 10 ? 1 : 0;
    let keyB = a.key == 11 ? 1 : 0;
    let metrics = [a.danceability, a.energy, a.loudness, a.acousticness, a.instrumentalness, a.valence, a.tempo];
    let newMetrics = [];
    for (let i = 0; i < metrics.length; i++) {
        let metric = metrics[i];
        newMetrics[i] = (metric - meanstd.mean[i]) / meanstd.std[i];
    }
    let tensor = [a.mode, keyC, keyCs, keyD, keyDs, keyE, keyF, keyFs, keyG, keyGs, keyA, keyAs, keyB, newMetrics[0], newMetrics[1], newMetrics[2], newMetrics[3], newMetrics[4], newMetrics[5], newMetrics[6]];
    let aTensor = tf.tensor2d([tensor.slice()]);
    let pTensor = tf.tensor2d([tensor.slice()]);
    await modelsReady;
    let originalA = aModel.predict(aTensor).dataSync()[0];
    let aetherealness = Math.pow(originalA, 1);
    let originalP = pModel.predict(pTensor).dataSync()[0];
    let primordialness = Math.pow(originalP, 1);
    console.log(`Aethereal: ${aetherealness}, Primordial: ${primordialness}`);
    return {aetherealness, primordialness};
}
