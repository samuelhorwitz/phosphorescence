import {builders} from '~/assets/recordcrate';

export const state = () => ({
    tracksPerPlaylist: 10,
    seedStyle: null,
    onlyTheHits: false,
    lowEnergy: false
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
    updateLowEnergy(state, lowEnergy) {
        localStorage.setItem('lowEnergy', lowEnergy);
        state.lowEnergy = lowEnergy;
    },
    restore(state) {
        let tracksPerPlaylist = localStorage.getItem('tracksPerPlaylist');
        let seedStyle = localStorage.getItem('seedStyle');
        let onlyTheHits = localStorage.getItem('onlyTheHits');
        let lowEnergy = localStorage.getItem('lowEnergy');
        if (tracksPerPlaylist) {
            state.tracksPerPlaylist = parseInt(tracksPerPlaylist, 10);
        }
        if (seedStyle && seedStyle !== 'null') {
            state.seedStyle = seedStyle;
        }
        if (onlyTheHits) {
            state.onlyTheHits = onlyTheHits === 'true';
        }
        if (lowEnergy) {
            state.lowEnergy = lowEnergy === 'true';
        }
    }
};

export const getters = {
    pruners(state) {
        let pruners = [];
        if (state.onlyTheHits) {
            pruners.push(builders.hits);
        }
        if (state.lowEnergy) {
            pruners.push(builders.lowEnergy);
        }
        if (pruners.length) {
            return pruners;
        }
        return null;
    }
}
