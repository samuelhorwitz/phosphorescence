import {getAccessToken, quickReturnAccessTokenWithoutGuarantee} from '~/assets/session';

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
