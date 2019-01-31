export const state = () => ({
    tracks: null,
    analysis: null,
    currentTrackIndex: 0
});

export const mutations = {
    load(state, {tracks, analysis}) {
        state.tracks = tracks;
        state.analysis = analysis;
        state.currentTrackIndex = 0;
    },
    nextTrack(state) {
        state.currentTrackIndex += 1;
    }
};

export const getters = {
    currentTrack(state) {
        if (!state.tracks) {
            return null;
        }
        return state.tracks[state.currentTrackIndex];
    },
    currentTrackAnalysis(state) {
        if (!state.tracks) {
            return null;
        }
        return state.analysis[state.tracks[state.currentTrackIndex].track.id];
    }
};
