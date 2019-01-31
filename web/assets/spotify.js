import {getAccessToken, quickReturnAccessTokenWithoutGuarantee} from '~/assets/session';

const playlists = [
    '5lSgExb6yTKfLusGag7bm7',
    '5QafFMGgQKGwqgV7k3qHy6',
    '68BiK8KG3otORDaYB3ZaO9',
    '0uVIGYfnUAkOT5REqtQICx',
    '5WdZ7SDAgBE8kK0EYkn3xj',
    '5ztDD0tx94xaHpIRpepP2r',
    '3bYwSveexd8YNXvS84Hj12'
];

export async function spider() {
    let allTracks = {};
    console.log('Getting all genre seed tracks...');
    while (Object.entries(allTracks).length < 850) {
        console.log(Object.entries(allTracks).length);
        allTracks = Object.assign(allTracks, await getGenreSeedTracks());
    }
    let popularTranceArtists = getMostPopularArtists(allTracks);
    console.log('Getting more tracks by artist...');
    while (Object.entries(allTracks).length < 10000) {
        console.log(Object.entries(allTracks).length);
        allTracks = Object.assign(allTracks, await getRandomSeedArtists(popularTranceArtists.splice(0, 4)));
    }
    return allTracks;
}

function getMostPopularArtists(tracks) {
    let artistCounts = {};
    for (let track of Object.values(tracks)) {
        for (let artist of track.track.artists) {
            if (!artistCounts[artist.id]) {
                artistCounts[artist.id] = {count: 1, name: artist.name};
            }
            else {
                artistCounts[artist.id].count++;
            }
        }
    }
    let orderedArtists = [];
    for (let [id, {count, name}] of Object.entries(artistCounts)) {
        orderedArtists.push({id, name, count});
    }
    orderedArtists.sort((a, b) => {
        if (a.count > b.count) return -1;
        else if (a.count < b.count) return 1;
        return 0;
    });
    console.log(orderedArtists);
    debugger
    return orderedArtists;
}

async function getGenreSeedTracks() {
    let response = await fetch('https://api.spotify.com/v1/recommendations?seed_genres=trance&limit=100', {
        headers: {
            Authorization: `Bearer ${await getAccessToken()}`
        }
    });
    let json = await response.json();
    let newTracks = json.tracks.filter(t => t != null);
    return indexTracks(newTracks, await getAudioFeatures(newTracks));
}

async function getRandomSeedArtists(artists) {
    let response = await fetch(`https://api.spotify.com/v1/recommendations?seed_genres=trance&limit=100&seed_artists=${artists.map(t => t.id).join(',')}`, {
        headers: {
            Authorization: `Bearer ${await getAccessToken()}`
        }
    });
    let json = await response.json();
    let newTracks = json.tracks.filter(t => t != null);
    return indexTracks(newTracks, await getAudioFeatures(newTracks));
}

async function getAudioFeatures(tracks) {
    let response = await fetch(`https://api.spotify.com/v1/audio-features?ids=${tracks.map(t => t.id).join(',')}`, {
        headers: {
            Authorization: `Bearer ${await getAccessToken()}`
        }
    });
    let {audio_features} = await response.json();
    return audio_features;
}

function indexTracks(tracks, features) {
    let allTracks = {};
    for (let track of tracks) {
        if (allTracks[track.id]) {
            continue;
        }
        allTracks[track.id] = {track};
    }
    for (let trackFeatures of features) {
        if (!allTracks[trackFeatures.id]) {
            continue;
        }
        allTracks[trackFeatures.id].features = trackFeatures;
    }
    return allTracks;
}

export async function initializePlayer(store, storePrefix) {
    let script;
    store.dispatch(`${storePrefix}/restoreSpotifyState`);
    await new Promise((resolve, reject) => {
        window.onSpotifyWebPlaybackSDKReady = async () => {
            let deviceId;
            try {
                deviceId = await initializePlayerListeners(store, storePrefix);
            }
            catch ({message, reauth}) {
                console.error(message);
                reject(message);
                if (reauth) {
                    await getAccessToken();
                }
                return;
            }
            store.commit(`${storePrefix}/deviceId`, deviceId);
            resolve();
        };
        script = document.createElement('script');
        script.src = 'https://sdk.scdn.co/spotify-player.js';
        document.head.appendChild(script);
    });
    return () => {
        document.head.removeChild(script);
    };
}

async function initializePlayerListeners(store, storePrefix) {
    return new Promise((resolve, reject) => {
        let player = new Spotify.Player({
            name: 'Phosphorescence',
            getOAuthToken: cb => cb(quickReturnAccessTokenWithoutGuarantee())
        });
        player.addListener('initialization_error', async ({ message }) => {
            reject({message});
        });
        player.addListener('authentication_error', async ({ message }) => {
            reject({message, reauth: true});
        });
        player.addListener('account_error', async ({ message }) => {
            reject({message});
        });
        player.addListener('playback_error', ({ message }) => { reject(message); });

        // Playback status updates
        player.addListener('player_state_changed', state => {
            if (!state) {
                return;
            }
            console.log(state);
            let ourCurrentTrack = store.getters[`${storePrefix}/currentTrack`];
            if (ourCurrentTrack && state.track_window.current_track.id != ourCurrentTrack.track.id) {
                let newCursor = store.getters[`${storePrefix}/getPlaylistCursorById`](state.track_window.current_track.id);
                if (newCursor) {
                    store.commit(`${storePrefix}/seekTrack`, newCursor);
                }
            }
            if (state.paused && !store.getters[`${storePrefix}/paused`]) {
                store.commit(`${storePrefix}/pause`);
            }
            else if (!state.paused && store.getters[`${storePrefix}/paused`]) {
                store.commit(`${storePrefix}/resume`);
            }
            store.commit(`${storePrefix}/lastKnownSpotifyState`, state);
        });

        // Ready
        player.addListener('ready', ({ device_id }) => {
            console.log('Ready with Device ID', device_id);
            resolve(device_id);
        });

        // Not Ready
        player.addListener('not_ready', ({ device_id }) => {
            // TODO
            console.log('Device ID has gone offline', device_id);
        });

        // Connect to the player!
        player.connect();
    });
}
