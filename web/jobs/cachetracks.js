if (process.env.NODE_ENV !== 'production') {
    require('dotenv').config({path: require('path').resolve(process.cwd(), '../.env')});
}
const axios = require('axios');
const querystring = require('querystring');
const fs = require('fs');
const {seeders, whitelists, blacklists} = require('./playlists.js');
const mixCutsMatcher = /([\[(](mix cut|mixed)[\])]|\bmix cut)/i;

async function consumeFromSpotify(url, token) {
    try {
        let {status, statusText, data} = await axios.get(url, {headers:
            {Authorization: `Bearer ${token}`}
        });
        return data;
    }
    catch ({response, request, message}) {
        console.error(response.data);
        console.error(response.status);
        console.error(response.headers);
        return null;
    }
}

async function getAppToken() {
    try {
        let {data} = await axios.post('https://accounts.spotify.com/api/token', querystring.stringify({
            grant_type: 'client_credentials',
            client_id: process.env.SPOTIFY_CLIENT_ID,
            client_secret: process.env.SPOTIFY_SECRET
        }));
        return data.access_token;
    }
    catch ({response, request, message}) {
        console.error(response.data);
        console.error(response.status);
        console.error(response.headers);
        return null;
    }
}

async function getFeatures(tracks, token) {
    let ids = Object.keys(tracks).join(',');
    let response = await consumeFromSpotify(`https://api.spotify.com/v1/audio-features?ids=${ids}`, token);
    if (!response) {
        throw new Error('Could not get features, failing');
    }
    let featuresById = {};
    for (let features of response.audio_features) {
        if (!features) {
            continue;
        }
        featuresById[features.id] = features;
    }
    let newTracks = {};
    for (let [id, track] of Object.entries(tracks)) {
        if (track.track && featuresById[id]) {
            newTracks[id] = track;
            newTracks[id].features = featuresById[id];
        }
    }
    return newTracks;
}

async function getFeaturesInBatches(tracks, token) {
    let cursor = 0;
    let trackIds = Object.keys(tracks);
    let tracksLength = trackIds.length;
    let newTracks = {};
    while (cursor < tracksLength) {
        let ids = trackIds.slice(cursor, cursor + 100);
        console.log(`Getting feature batch ${(cursor / 100) + 1} of ${Math.ceil(tracksLength / 100)}, ids ${ids[0]} to ${ids[ids.length - 1]}...`)
        let batchOfTracks = {};
        for (let id of ids) {
            batchOfTracks[id] = tracks[id];
        }
        let tracksWithFeatures = await getFeatures(batchOfTracks, token);
        for (let [id, trackWithFeatures] of Object.entries(tracksWithFeatures)) {
            newTracks[id] = trackWithFeatures;
        }
        cursor += 100;
    }
    return newTracks;
}

async function getPlaylist(playlistId, token) {
    let nextUrl = `https://api.spotify.com/v1/playlists/${playlistId}/tracks?limit=100`;
    let tracks = {};
    while (nextUrl) {
        console.log(`Handling ${nextUrl}...`);
        let response = await consumeFromSpotify(nextUrl, token);
        if (!response) {
            throw new Error('Could not get tracks, failing');
        }
        for (let track of response.items) {
            tracks[track.track.id] = {track: track.track};
        }
        nextUrl = response.next;
    }
    return tracks;
}

async function run(playlistId) {
    let token = await getAppToken();
    let blacklist = {};
    for (let blacklist of blacklists) {
        let blacklistTracks = await getPlaylist(blacklist, token);
        for (let [id, blacklistTrack] of Object.entries(blacklistTracks)) {
            blacklist[id] = true;
        }
    }
    let tracks = {};
    for (let whitelist of whitelists) {
        let whitelistTracks = await getPlaylist(whitelist, token);
        for (let [id, whitelistTrack] of Object.entries(whitelistTracks)) {
            tracks[id] = whitelistTrack;
        }
    }
    for (let seeder of seeders) {
        let playlistTracks = await getPlaylist(seeder, token);
        for (let [id, playlistTrack] of Object.entries(playlistTracks)) {
            if (mixCutsMatcher.test(playlistTrack.track.name) || blacklist[id]) {
                continue;
            }
            tracks[id] = playlistTrack;
        }
    }
    let newTracks = await getFeaturesInBatches(tracks, token);
    return newTracks;
}

(async function() {
    try {
        let tracks = await run();
        await new Promise((resolve, reject) => {
            fs.writeFile(`${__dirname}${process.env.TRACKS_JSON}`, JSON.stringify(tracks), function(err) {
                if (err) {
                    reject(err);
                }
                else {
                    resolve();
                }
            }); 
        });
    }
    catch (e) {
        console.error(e);
    }
})();
