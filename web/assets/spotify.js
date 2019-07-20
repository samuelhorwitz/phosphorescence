import {getAccessToken, quickReturnAccessTokenWithoutGuarantee} from '~/assets/session';

export async function initializePlayer(store, storePrefix) {
    let script, player;
    await new Promise((resolve, reject) => {
        window.onSpotifyWebPlaybackSDKReady = async () => {
            let deviceId;
            try {
                let res = await initializePlayerListeners(store, storePrefix);
                player = res.player;
                deviceId = res.deviceId;
            }
            catch ({message, reauth}) {
                console.error(message);
                reject(message);
                if (reauth) {
                    await getAccessToken();
                }
                return;
            }
            await fetch(`${process.env.API_ORIGIN}/device/${deviceId}?playState=pause`, {
                method: 'PUT',
                credentials: 'include'
            });
            store.commit(`${storePrefix}/deviceId`, deviceId);
            resolve();
        };
        script = document.createElement('script');
        script.src = 'https://sdk.scdn.co/spotify-player.js';
        document.head.appendChild(script);
    });
    return {
        destroyer: () => {
            player.disconnect();
            document.head.removeChild(script);
        },
        player
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
            console.log('Player state changed');
            if (!state) {
                console.log('Device disconnected');
                store.commit(`${storePrefix}/disconnected`);
                return;
            }
            console.log(state);
            store.commit(`${storePrefix}/connected`);
            let ourCurrentTrack = store.getters[`${storePrefix}/currentTrack`];
            let trackOkay = true;
            if (ourCurrentTrack && state.track_window.current_track.id != ourCurrentTrack.track.id) {
                let newCursor = store.getters[`${storePrefix}/getPlaylistCursorById`](state.track_window.current_track.id);
                if (newCursor) {
                    if (newCursor != -1) {
                        store.commit(`${storePrefix}/seekTrack`, newCursor);
                    } else {
                        // TODO alert user that this song cannot be played in Phosphorescence as it's not in the current lineup.
                        store.dispatch(`${storePrefix}/stop`);
                        trackOkay = false;
                    }
                }
            }
            if (trackOkay) {
                if (state.paused && !(store.getters[`${storePrefix}/paused`] || store.getters[`${storePrefix}/stopped`])) {
                    store.commit(`${storePrefix}/pause`);
                }
                else if (!state.paused && (store.getters[`${storePrefix}/paused`] || store.getters[`${storePrefix}/stopped`])) {
                    store.commit(`${storePrefix}/resume`);
                }
            }
            store.commit(`${storePrefix}/lastKnownSpotifyState`, state);
        });

        // Ready
        player.addListener('ready', ({ device_id }) => {
            console.log('Ready with Device ID', device_id);
            resolve({player, deviceId: device_id});
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
