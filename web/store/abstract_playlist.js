import {getAccessToken} from '~/assets/session';

const STOPPED = 0;
const PLAYING = 1;
const PAUSED = 2;

export function initialize(storagePrefix) {
    return {
        state: getState(),
        mutations: getMutations(storagePrefix),
        actions: getActions(),
        getters: getGetters()
    };
};

const getState = () => () => ({
    currentTrackCursor: 0,
    playlist: null,
    playback: STOPPED,
    deviceId: null,
    spotifyState: null,
    spotifyFullyRestored: false,
    spotifyAppearsDown: false
});

const getMutations = storagePrefix => Object.assign({
    nextTrack(state) {
        if (!canSkipForward(state)) {
            return;
        }
        state.currentTrackCursor++;
    },
    previousTrack(state) {
        if (!canSkipBackward(state)) {
            return;
        }
        state.currentTrackCursor--;
    },
    seekTrack(state, cursor) {
        if (!isCursorInRange(state, cursor)) {
            return;
        }
        state.currentTrackCursor = cursor;
    },
    loadPlaylist(state, playlist) {
        sessionStorage.setItem(`${storagePrefix}/currentPlaylist`, JSON.stringify(playlist));
        state.playlist = playlist;
        state.currentTrackCursor = 0;
        state.playback = STOPPED;
    },
    restore(state) {
        let playlist = sessionStorage.getItem(`${storagePrefix}/currentPlaylist`);
        if (playlist) {
            state.playlist = JSON.parse(playlist);
            state.currentTrackCursor = 0;
            state.playback = STOPPED;
        }
        let spotifyState = sessionStorage.getItem(`${storagePrefix}/spotifyState`);
        if (spotifyState) {
            state.spotifyState = JSON.parse(spotifyState);
        }
    },
    lastKnownSpotifyState(state, spotifyState) {
        sessionStorage.setItem(`${storagePrefix}/spotifyState`, JSON.stringify(spotifyState));
        state.spotifyState = spotifyState;
    },
    play(state) {
        state.playback = PLAYING;
    },
    pause(state) {
        state.playback = PAUSED;
    },
    resume(state) {
        state.playback = PLAYING;
    },
    stop(state) {
        state.playback = STOPPED;
    },
    stopBecauseBroken(state) {
        state.playback = STOPPED;
        state.spotifyAppearsDown = true;
    },
    deviceId(state, deviceId) {
        state.deviceId = deviceId;
    },
    spotifyFullyRestored(state) {
        state.spotifyFullyRestored = true;
    }
});

const getActions = () => Object.assign({
    async play({commit, dispatch, getters, state}) {
        commit('loading/startLoad', null, {root: true});
        commit('play');
        let response = await fetch(`https://api.spotify.com/v1/me/player/play${getters.deviceIdAsQueryString}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${await getAccessToken()}`
            },
            body: JSON.stringify({
                uris: state.playlist.map(track => track.track.uri),
                offset: {position: state.currentTrackCursor}
            })
        });
        if (response && response.status >= 500) {
            commit('stopBecauseBroken');
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    async resume({commit, dispatch, getters}) {
        commit('loading/startLoad', null, {root: true});
        commit('play');
        let response = await fetch(`https://api.spotify.com/v1/me/player/play${getters.deviceIdAsQueryString}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${await getAccessToken()}`
            }
        });
        if (response && response.status >= 500) {
            commit('stopBecauseBroken');
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    async pause({commit, dispatch, getters}) {
        commit('loading/startLoad', null, {root: true});
        commit('pause');
        let response = await fetch(`https://api.spotify.com/v1/me/player/pause${getters.deviceIdAsQueryString}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${await getAccessToken()}`
            }
        });
        if (response && response.status >= 500) {
            commit('stopBecauseBroken');
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    async stop({commit, dispatch, getters}) {
        commit('loading/startLoad', null, {root: true});
        commit('stop');
        let response = await fetch(`https://api.spotify.com/v1/me/player/pause${getters.deviceIdAsQueryString}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${await getAccessToken()}`
            }
        });
        if (response && response.status >= 500) {
            commit('stopBecauseBroken');
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    async next({commit, dispatch, getters, state}) {
        if (!canSkipForward(state)) {
            return;
        }
        commit('loading/startLoad', null, {root: true});
        commit('nextTrack');
        if (!getters.stopped) {
            let response = await fetch(`https://api.spotify.com/v1/me/player/next${getters.deviceIdAsQueryString}`, {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${await getAccessToken()}`
                }
            });
            if (response && response.status >= 500) {
                commit('stopBecauseBroken');
            }
            else {
                commit('play');
            }
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    async previous({commit, dispatch, getters, state}) {
        if (!canSkipBackward(state)) {
            return;
        }
        commit('loading/startLoad', null, {root: true});
        commit('previousTrack');
        if (!getters.stopped) {
            let response = await fetch(`https://api.spotify.com/v1/me/player/previous${getters.deviceIdAsQueryString}`, {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${await getAccessToken()}`
                }
            });
            if (response && response.status >= 500) {
                commit('stopBecauseBroken');
            }
            else {
                commit('play');
            }
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    async seekTrack({commit, dispatch, getters, state}, cursor) {
        if (!isCursorInRange(state, cursor)) {
            return;
        }
        commit('loading/startLoad', null, {root: true});
        commit('seekTrack', cursor);
        if (!getters.stopped) {
            let response = await fetch(`https://api.spotify.com/v1/me/player/play${getters.deviceIdAsQueryString}`, {
                method: 'PUT',
                headers: {
                    Authorization: `Bearer ${await getAccessToken()}`
                },
                body: JSON.stringify({
                    uris: state.playlist.map(track => track.track.uri),
                    offset: {position: state.currentTrackCursor}
                })
            });
            if (response && response.status >= 500) {
                commit('stopBecauseBroken');
            }
            else {
                commit('play');
            }
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    async restoreSpotifyState({commit, dispatch, getters, state, nextTick}) {
        commit('loading/startLoad', null, {root: true});
        let response = await fetch('https://api.spotify.com/v1/me/player', {
            method: 'GET',
            headers: {
                Authorization: `Bearer ${await getAccessToken()}`
            }
        });
        if (false && response && response.status == 200) {
            let currentSpotifyState = await response.json();
            if (currentSpotifyState.currently_playing_type == 'track') {
                commit('seekTrack', getters.getPlaylistCursorById(currentSpotifyState.item.id));
                let response = await fetch('https://api.spotify.com/v1/me/player', {
                    method: 'PUT',
                    headers: {
                        Authorization: `Bearer ${await getAccessToken()}`
                    },
                    body: JSON.stringify({
                        device_ids: [state.deviceId]
                    })
                });
            }
        }
        else if (false && state.spotifyState) {
            commit('seekTrack', getters.getPlaylistCursorById(state.spotifyState.track_window.current_track.id));
        }
        commit('spotifyFullyRestored');
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    loadPlaylist({commit, dispatch}, playlist) {
        commit('loadPlaylist', playlist);
        dispatch('stop');
    }
});

const getGetters = () => Object.assign({
    currentTrack(state) {
        if (!state.playlist || !state.spotifyFullyRestored) {
            return null;
        }
        return state.playlist[state.currentTrackCursor];
    },
    canSkipBackward(state) {
        return canSkipBackward(state);
    },
    canSkipForward(state) {
        return canSkipForward(state);
    },
    stopped(state) {
        return state.playback == STOPPED;
    },
    playing(state) {
        return state.playback == PLAYING;
    },
    paused(state) {
        return state.playback == PAUSED;
    },
    playlistLoaded(state) {
        return !!state.playlist;
    },
    deviceLoaded(state) {
        return !!state.deviceId;
    },
    deviceIdAsQueryString(state) {
        let queryString = '';
        if (state.deviceId) {
            queryString = `?device_id=${state.deviceId}`;
        }
        return queryString;
    },
    getPlaylistCursorById: state => id => state.playlist ? state.playlist.findIndex(track => track.track.id == id) : null
});

function canSkipBackward(state) {
    if (!state.playlist) {
        return false;
    }
    return state.currentTrackCursor > 0;
}

function canSkipForward(state) {
    if (!state.playlist) {
        return false;
    }
    return state.currentTrackCursor < state.playlist.length - 1;
}

function isCursorInRange(state, cursor) {
    if (!state.playlist) {
        return false;
    }
    return cursor >= 0 && cursor <= state.playlist.length - 1;
}
