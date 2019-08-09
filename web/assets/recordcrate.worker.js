import {loadModel, tensor2d, tidy} from '@tensorflow/tfjs';
import {getTrackTag} from '~/common/normalize';
import {encoder, decoder} from '~/common/textencoding';
import pako from 'pako';

const modelsReady = new Promise(async resolve => {
    let aModel = await loadModel('/models/aetherealness/model.json');
    let pModel = await loadModel('/models/primordialness/model.json');
    let meanstdRaw = await fetch('/models/meanstd.json');
    let meanstd = await meanstdRaw.json();
    resolve({aModel, pModel, meanstd});
});

addEventListener('message', async ({data}) => {
    if (data.type == 'sendTracks') {
        let tracks = JSON.parse(decoder.decode(data.tracks));
        let countryCode = data.countryCode;
        tracks = filterTracks(tracks, countryCode);
        console.log('Processing tracks...');
        tracks = await getEvocativeness(tracks);
        console.log('Tracks processed');
        let {tags, idsToTags} = buildTags(tracks);
        let responseData = encoder.encode(JSON.stringify({tracks, tags, idsToTags}));
        let gzipResponseData = pako.gzip(responseData, {level: 9});
        postMessage({type: 'sendProcessedTracks', gzipData: gzipResponseData.buffer});
    } else if (data.type == 'sendTrack') {
        let track = data.track;
        let countryCode = data.countryCode;
        console.log('Processing track...');
        track = await getEvocativenessOfSingleTrack(track);
        console.log('Track processed');
        let tag = getTrackTag(track.track);
        postMessage({type: 'sendProcessedTrack', data: {track: track.track, features: track.features, evocativeness: track.evocativeness, tag}});
    }
});

function buildTags(tracks) {
    let tags = {};
    let idsToTags = {}
    for (let [id, {track}] of Object.entries(tracks)) {
        let tag = getTrackTag(track);
        if (tags[tag]) {
            tags[tag].push(id);
        }
        else {
            tags[tag] = [id];
        }
        idsToTags[id] = tag;
    }
    return {tags, idsToTags};
}

async function getEvocativeness(tracks) {
    let {aModel, pModel, meanstd} = await modelsReady;
    for (let [id, track] of Object.entries(tracks)) {
        let {features} = track;
        tracks[id].evocativeness = predict(features, aModel, pModel, meanstd);
    }
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

function filterTracks(tracks, countryCode) {
    console.log('Tracks before region pruning:', Object.keys(tracks).length);
    for (let [id, track] of Object.entries(tracks)) {
        if (track.track.available_markets.indexOf(countryCode) == -1) {
            delete tracks[id];
        }
    }
    console.log('Tracks after region pruning:', Object.keys(tracks).length);
    return tracks;
}
