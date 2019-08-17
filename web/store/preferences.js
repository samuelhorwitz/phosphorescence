export const state = () => ({
    tracksPerPlaylist: 10,
    seedStyle: null,
    onlyTheHits: false
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
    updateOnlyTheHits(state, onlyTheHits) {
        localStorage.setItem('onlyTheHits', onlyTheHits);
        state.onlyTheHits = onlyTheHits;
    },
    restore(state) {
        let tracksPerPlaylist = localStorage.getItem('tracksPerPlaylist');
        let seedStyle = localStorage.getItem('seedStyle');
        let onlyTheHits = localStorage.getItem('onlyTheHits');
        if (tracksPerPlaylist) {
            state.tracksPerPlaylist = parseInt(tracksPerPlaylist, 10);
        }
        if (seedStyle && seedStyle !== 'null') {
            state.seedStyle = seedStyle;
        }
        if (onlyTheHits) {
            state.onlyTheHits = onlyTheHits === 'true';
        }
    }
};
