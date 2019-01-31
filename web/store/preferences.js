export const state = () => ({
    tracksPerPlaylist: 10,
    seedStyle: null
});

export const mutations = {
    updateTracksPerPlaylist(state, count) {
        localStorage.setItem('tracksPerPlaylist', count);
        state.tracksPerPlaylist = count;
    },
    updateSeedStyle(state, seedStyle) {
        localStorage.setItem('seedStyle', seedStyle);
        state.seedStyle = seedStyle;
    },
    restore(state) {
        let tracksPerPlaylist = localStorage.getItem('tracksPerPlaylist');
        let seedStyle = localStorage.getItem('seedStyle');
        if (tracksPerPlaylist) {
            state.tracksPerPlaylist = parseInt(tracksPerPlaylist, 10);
        }
        if (seedStyle) {
            state.seedStyle = seedStyle;
        }
    }
};
