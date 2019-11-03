import {getAccessToken} from '~/assets/session';
import {getSpotifyTrackUri} from '~/assets/spotify';
import {getCaptchaToken} from '~/assets/captcha';

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
    selectedTrackCursor: 0,
    playlist: null,
    previews: {},
    currentPreview: null,
    currentPreviewPercent: 0,
    previewLocked: false,
    playback: STOPPED,
    deviceId: null,
    spotifyState: null,
    spotifyAppearsDown: false,
    playerState: NOT_READY,
    neighborSeekLocked: false,
    deviceName: 'Phosphorescence'
});

const getMutations = storagePrefix => Object.assign({
    lockPreview(state) {
        state.previewLocked = true;
    },
    unlockPreview(state) {
        state.previewLocked = false;
    },
    nextTrack(state) {
        if (!canSkipForward(state)) {
            return;
        }
        state.currentTrackCursor++;
        state.selectedTrackCursor = state.currentTrackCursor;
    },
    previousTrack(state) {
        if (!canSkipBackward(state)) {
            return;
        }
        state.currentTrackCursor--;
        state.selectedTrackCursor = state.currentTrackCursor;
    },
    seekTrack(state, cursor) {
        if (!isCursorInRange(state, cursor)) {
            return;
        }
        state.currentTrackCursor = cursor;
        state.selectedTrackCursor = state.currentTrackCursor;
    },
    selectTrack(state, cursor) {
        if (!isCursorInRange(state, cursor)) {
            return;
        }
        state.selectedTrackCursor = cursor;
    },
    selectNextTrack(state) {
        if (!canSelectNext(state)) {
            return;
        }
        state.selectedTrackCursor++;
    },
    selectPreviousTrack(state) {
        if (!canSelectPrevious(state)) {
            return;
        }
        state.selectedTrackCursor--;
    },
    loadPlaylist(state, playlist) {
        if (navigator.standalone) {
            localStorage.setItem(`${storagePrefix}/currentPlaylist`, JSON.stringify(playlist));
        } else if (/\b(iPhone|iPod)\b/.test(navigator.userAgent)) {
            let oldId = getCurrentPlaylistIdIos(storagePrefix);
            localStorage.removeItem(`${storagePrefix}/currentPlaylist-${oldId}`);
            let id = new Date().getTime();
            localStorage.setItem(`${storagePrefix}/currentPlaylist-${id}`, JSON.stringify(playlist));
            setCurrentPlaylistIdIos(storagePrefix, id);
        } else {
            sessionStorage.setItem(`${storagePrefix}/currentPlaylist`, JSON.stringify(playlist));
        }
        state.playlist = playlist;
        state.currentTrackCursor = 0;
        state.selectedTrackCursor = 0;
        state.playback = STOPPED;
    },
    clearPlaylist(state) {
        state.playlist = null;
        state.currentTrackCursor = 0;
        state.selectedTrackCursor = 0;
        state.playback = STOPPED;
    },
    loadPreviews(state, previews) {
        state.previews = {...state.previews, ...previews};
    },
    restore(state) {
        let playlist;
        if (navigator.standalone) {
            playlist = localStorage.getItem(`${storagePrefix}/currentPlaylist`);
        } else if (/\b(iPhone|iPod)\b/.test(navigator.userAgent)) {
            let id = getCurrentPlaylistIdIos(storagePrefix);
            if (id) {
                playlist = localStorage.getItem(`${storagePrefix}/currentPlaylist-${id}`);
                if (!playlist) {
                    clearCurrentPlaylistIdIos(storagePrefix);
                }
            }
        } else {
            playlist = sessionStorage.getItem(`${storagePrefix}/currentPlaylist`);
        }
        if (playlist) {
            state.playlist = JSON.parse(playlist);
            state.currentTrackCursor = 0;
            state.selectedTrackCursor = 0;
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
    playPreview(state, id) {
        if (state.currentPreview == id) {
            return;
        }
        if (!state.previews[id]) {
            state.currentPreview = null;
        } else {
            state.currentPreview = id;
        }
        state.currentPreviewPercent = 0;
    },
    playPreviewOfSelectedTrack(state) {
        if (!state.playlist || !state.playlist[state.selectedTrackCursor]) {
            return;
        }
        let id = state.playlist[state.selectedTrackCursor].id;
        if (state.currentPreview == id) {
            return;
        }
        if (!state.previews[id]) {
            state.currentPreview = null;
        } else {
            state.currentPreview = id;
        }
        state.currentPreviewPercent = 0;
    },
    stopPreview(state) {
        state.currentPreview = null;
        state.currentPreviewPercent = 0;
    },
    completePreview(state) {
        state.currentPreview = null;
        state.currentPreviewPercent = 0;
    },
    updatePreviewPercent(state, percent) {
        state.currentPreviewPercent = percent;
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
    async play({commit, dispatch, getters, state}, msOffset) {
        commit('loading/startLoad', null, {root: true});
        commit('play');
        let body = {
            uris: state.playlist.map(track => getSpotifyTrackUri(track.id)),
            offset: {position: state.currentTrackCursor}
        };
        if (msOffset) {
            body.position_ms = msOffset;
        }
        let response = await fetch('https://api.spotify.com/v1/me/player/play', {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${await getAccessToken()}`
            },
            body: JSON.stringify(body)
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
    seekSelectedTrack({dispatch, state}) {
        dispatch('seekTrack', state.selectedTrackCursor);
    },
    loadPlaylist({commit, dispatch, rootState, rootGetters}, playlist) {
        (async () => {
            if (!playlist) {
                return;
            }
            let trackPreviews = await loadPreviews(playlist.map(t => t.id), !!rootState.user.user, rootGetters['user/country']);
            commit('loadPreviews', trackPreviews);
        })();
        commit('loadPlaylist', playlist);
        dispatch('stop');
    },
    restore({commit, state, rootState, rootGetters}) {
        commit('restore');
        (async () => {
            if (!state.playlist) {
                return;
            }
            let trackPreviews = await loadPreviews(state.playlist.map(t => t.id), !!rootState.user.user, rootGetters['user/country']);
            commit('loadPreviews', trackPreviews);
        })();
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
    selectedTrack(state) {
        if (!state.playlist) {
            return null;
        }
        return state.playlist[state.selectedTrackCursor];
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
    getPlaylistCursorById: state => id => state.playlist ? state.playlist.findIndex(track => track.id == id) : null
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

function canSelectPrevious(state) {
    if (!state.playlist) {
        return false;
    }
    return state.selectedTrackCursor > 0;
}

function canSelectNext(state) {
    if (!state.playlist) {
        return false;
    }
    return state.selectedTrackCursor < state.playlist.length - 1;
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
            uris: state.playlist.map(track => getSpotifyTrackUri(track.id)),
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
    if (id) {
        setCurrentPlaylistIdIos(storagePrefix, id);
    }
    return id;
}

function setCurrentPlaylistIdIos(storagePrefix, id) {
    location.hash = `#${id}`;
    sessionStorage.setItem(`${storagePrefix}/currentPlaylistId`, id);
}

function clearCurrentPlaylistIdIos(storagePrefix) {
    location.hash = '';
    sessionStorage.removeItem(`${storagePrefix}/currentPlaylistId`);
}

async function loadPreviews(ids, isLoggedInUser, region) {
    let trackPreviewsResponse;
    let trackIdsStr = ids.join(',');
    if (isLoggedInUser) {
        trackPreviewsResponse = await fetch(`${process.env.API_ORIGIN}/track/preview/${trackIdsStr}`, {credentials: 'include'});
    } else {
        let captcha = await getCaptchaToken('api/track/preview');
        trackPreviewsResponse = await fetch(`${process.env.API_ORIGIN}/track/unauthenticated/preview/${region}/${trackIdsStr}?captcha=${captcha}`);
    }
    let {trackPreviews} = await trackPreviewsResponse.json();
    return trackPreviews;
}
