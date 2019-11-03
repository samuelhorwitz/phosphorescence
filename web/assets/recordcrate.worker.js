import {loadModel, tensor2d, tidy} from '@tensorflow/tfjs';
import {getTrackTag} from '~/common/normalize';
import {TextEncoder, TextDecoder} from 'text-encoding-shim';
import pako from 'pako';
import SecureMessenger from '~/secure-messenger/secure-messenger';
import throttle from 'lodash/throttle';

const encoder = new TextEncoder();
const decoder = new TextDecoder();

const modelsReady = new Promise(async resolve => {
    let aModel = await loadModel(`${process.env.STATIC_ORIGIN}/models/aetherealness/model.json`);
    let pModel = await loadModel(`${process.env.STATIC_ORIGIN}/models/primordialness/model.json`);
    let meanstdRaw = await fetch(`${process.env.STATIC_ORIGIN}/models/meanstd.json`);
    let meanstd = await meanstdRaw.json();
    resolve({aModel, pModel, meanstd});
});

let loadingInterruptPort;
let throttledLoadingMessage;
let loadingPercentBase = 0;

(async () => {
    let messenger = new SecureMessenger(location.origin);
    await messenger.knock([self]);
    messenger.messageHandlerLoop((data, interruptPort, finish) => {
        if (data.type == 'sendTracks') {
            finish();
            return handleSendTracks(data);
        } else if (data.type == 'sendTrack') {
            finish();
            return handleSendTrack(data);
        } else if (data.type == 'initializeLoadingNotificationChannel') {
            loadingInterruptPort = interruptPort;
            throttledLoadingMessage = throttle(loadingInterruptPort.postMessage.bind(loadingInterruptPort), 100);
            return {type: 'acknowledge'};
        }
        return {type: 'error', error: 'invalid request'};
    }, event => {
        console.error('Record crate worker could not handle Phosphor message', event);
    });
})();

async function handleSendTracks(data) {
    console.log('Will process tracks...');
    let tracksArr;
    if (data.raw) {
        tracksArr = data.tracks;
    } else {
        tracksArr = JSON.parse(decoder.decode(new Uint8Array(data.tracks)));
    }
    console.log('Processing tracks...');
    tracksArr = await getEvocativeness(tracksArr);
    console.log('Tracks processed');
    let {tags, idsToTags} = buildTags(tracksArr);
    let tracks = {};
    for (let track of tracksArr) {
        tracks[track.id] = track;
    }
    if (data.raw) {
        loadingInterruptPort && loadingInterruptPort.postMessage({type: 'loadPercent', value: 1});
        return {type: 'sendProcessedTracks', data: {tracks, tags, idsToTags}};
    }
    let responseData = encoder.encode(JSON.stringify({tracks, tags, idsToTags}));
    let gzipResponseData = pako.gzip(responseData, {level: 9});
    loadingInterruptPort && loadingInterruptPort.postMessage({type: 'loadPercent', value: 1});
    return {type: 'sendProcessedTracks', gzipData: gzipResponseData.buffer};
}

async function handleSendTrack(data) {
    let track = data.track;
    console.log('Processing track...');
    track = await getEvocativenessOfSingleTrack(track);
    console.log('Track processed');
    let tag = getTrackTag(track.track);
    return {type: 'sendProcessedTrack', data: {id: track.id, track: track.track, features: track.features, evocativeness: track.evocativeness, tag}};
}

function buildTags(tracks) {
    let tags = {};
    let idsToTags = {};
    let i = 0;
    for (let {id, track} of tracks) {
        i++;
        let tag = getTrackTag(track);
        if (tags[tag]) {
            tags[tag].push(id);
        }
        else {
            tags[tag] = [id];
        }
        idsToTags[id] = tag;
        updateLoadingPercent(i / tracks.length, 0.45);
    }
    loadingPercentBase = 0.9;
    return {tags, idsToTags};
}

async function getEvocativeness(tracks) {
    let {aModel, pModel, meanstd} = await modelsReady;
    let i = 0;
    for (let track of tracks) {
        i++;
        let {features} = track;
        track.evocativeness = predict(features, aModel, pModel, meanstd);
        updateLoadingPercent(i / tracks.length, 0.45);
    }
    loadingPercentBase = 0.45;
    return tracks;
}

async function getEvocativenessOfSingleTrack(track) {
    let {aModel, pModel, meanstd} = await modelsReady;
    track.evocativeness = predict(track.features, aModel, pModel, meanstd);
    return track;
}

function predict(a, aModel, pModel, meanstd) {
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
    return tidy(() => {
        let aTensor = tensor2d([tensor.slice()]);
        let pTensor = tensor2d([tensor.slice()]);
        let aetherealness = aModel.predict(aTensor).dataSync()[0];
        let primordialness = pModel.predict(pTensor).dataSync()[0];
        return {aetherealness, primordialness};
    });
}

function updateLoadingPercent(partialPercent, percentOfTotal) {
    throttledLoadingMessage && throttledLoadingMessage({type: 'loadPercent', value: loadingPercentBase + (partialPercent * percentOfTotal)});
}
