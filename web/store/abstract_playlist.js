import {getAccessToken} from '~/assets/session';

const STOPPED = 0;
const PLAYING = 1;
const PAUSED = 2;

const NOT_READY = 0;
const CONNECTED = 1;
const DISCONNECTED = 2;

let player;

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
    spotifyAppearsDown: false,
    playerState: NOT_READY,
    neighborSeekLocked: false,
    deviceName: 'Phosphorescence'
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
        if (/\b(iPhone|iPod)\b/.test(navigator.userAgent)) {
            let oldId = getCurrentPlaylistIdIos(storagePrefix);
            localStorage.removeItem(`${storagePrefix}/currentPlaylist-${oldId}`);
            let id = new Date().getTime();
            localStorage.setItem(`${storagePrefix}/currentPlaylist-${id}`, JSON.stringify(playlist));
            location.hash = `#${id}`;
            sessionStorage.setItem(`${storagePrefix}/currentPlaylistId`, id);
        } else {
            sessionStorage.setItem(`${storagePrefix}/currentPlaylist`, JSON.stringify(playlist));
        }
        state.playlist = playlist;
        state.currentTrackCursor = 0;
        state.playback = STOPPED;
    },
    clearPlaylist(state) {
        state.playlist = null;
        state.currentTrackCursor = 0;
        state.playback = STOPPED;
    },
    restore(state) {
        let playlist;
        if (/\b(iPhone|iPod)\b/.test(navigator.userAgent)) {
            let id = getCurrentPlaylistIdIos(storagePrefix);
            if (id) {
                playlist = localStorage.getItem(`${storagePrefix}/currentPlaylist-${id}`);
                if (!playlist) {
                    location.hash = '';
                    sessionStorage.removeItem(`${storagePrefix}/currentPlaylistId`);
                }
            }
        } else {
            playlist = sessionStorage.getItem(`${storagePrefix}/currentPlaylist`);
        }
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
    connected(state) {
        state.playerState = CONNECTED;
    },
    disconnected(state) {
        state.playerState = DISCONNECTED;
    },
    lockNeighborSeeking(state) {
        state.neighborSeekLocked = true;
    },
    unlockNeighborSeeking(state) {
        state.neighborSeekLocked = false;
    }
});

const getActions = () => Object.assign({
    registerPlayer({commit}, pl) {
        player = pl;
        commit('connected');
    },
    async play({commit, dispatch, getters, state}) {
        commit('loading/startLoad', null, {root: true});
        commit('play');
        let response = await fetch('https://api.spotify.com/v1/me/player/play', {
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
        if (player) {
            await player.resume();
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    async pause({commit, dispatch, getters}) {
        commit('loading/startLoad', null, {root: true});
        commit('pause');
        if (player) {
            await player.pause();
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    async stop({commit, dispatch, getters}) {
        commit('loading/startLoad', null, {root: true});
        commit('stop');
        if (player) {
            await player.pause();
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
            await seekToCurrentTrack(state, commit);
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
            await seekToCurrentTrack(state, commit);
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
            await seekToCurrentTrack(state, commit);
        }
        dispatch('loading/endLoadAfterDelay', null, {root: true});
    },
    loadPlaylist({commit, dispatch}, playlist) {
        commit('loadPlaylist', playlist);
        dispatch('stop');
    },
    clearPlaylist({commit, dispatch}) {
        commit('clearPlaylist');
        dispatch('stop');
    }
});

const getGetters = () => Object.assign({
    currentTrack(state) {
        if (!state.playlist) {
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
    isPlayerConnected(state) {
        return state.playerState == CONNECTED;
    },
    isPlayerDisconnected(state) {
        return state.playerState == NOT_READY || state.playerState == DISCONNECTED;
    },
    getPlaylistCursorById: state => id => state.playlist ? state.playlist.findIndex(track => track.track.id == id) : null
});

function canSkipBackward(state) {
    if (state.neighborSeekLocked) {
        return false;
    }
    if (!state.playlist) {
        return false;
    }
    return state.currentTrackCursor > 0;
}

function canSkipForward(state) {
    if (state.neighborSeekLocked) {
        return false;
    }
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

async function seekToCurrentTrack(state, commit) {
    commit('lockNeighborSeeking');
    let response = await fetch('https://api.spotify.com/v1/me/player/play', {
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
    // Seeking responds before seek is guaranteed, throttle seek requests
    // to prevent weird sync-y issues. We could reach out and check but it
    // is one more request to Spotify so let's keep it simple.
    await new Promise(resolve => setTimeout(resolve, 500));
    commit('unlockNeighborSeeking');
}

function getCurrentPlaylistIdIos(storagePrefix) {
    let id = location.hash.slice(1);
    if (!id) {
        id = sessionStorage.getItem(`${storagePrefix}/currentPlaylistId`);
    }
    return id;
}
