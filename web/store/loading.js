export const state = () => ({
    loading: true,
    loadSemaphore: 0,
    descriptions: [],
    progresses: {},
    progressWeights: {},
    playlistGenerating: false
});

export const mutations = {
    startLoad(state) {
        state.loading = true;
        state.loadSemaphore++;
    },
    endLoad(state) {
        state.loadSemaphore = safeSubtract(state.loadSemaphore);
        if (state.loadSemaphore == 0) {
            state.loading = false;
            state.loadSemaphore = 0;
        }
    },
    pushMessage(state, description) {
        state.descriptions.push(description);
    },
    clearMessage(state, messageId) {
        let newDescriptions = [];
        for (let i in state.descriptions) {
            if (state.descriptions[i].id !== messageId) {
                newDescriptions.push(state.descriptions[i]);
            }
        }
        state.descriptions = newDescriptions;
    },
    initializeProgress(state, {id, weight}) {
        if (!weight) {
            weight = 100;
        }
        weight = Math.min(weight, 100);
        state.progresses[id] = 0;
        state.progressWeights[id] = weight;
    },
    tickProgress(state, {id, percent}) {
        if (typeof state.progresses[id] === 'undefined') {
            return;
        }
        state.progresses = {...state.progresses, [id]: percent};
    },
    completeProgress(state, {id}) {
        if (typeof state.progresses[id] === 'undefined') {
            return;
        }
        state.progresses = {...state.progresses, [id]: 1};
    },
    resetProgress(state) {
        state.progresses = {};
        state.progressWeights = {};
    },
    playlistGenerating(state) {
        state.playlistGenerating = true;
    },
    playlistGenerationComplete(state) {
        state.playlistGenerating = false;
    }
};

let descriptionId = 0;

export const actions = {
    loadFlash({commit, dispatch, state}) {
        if (!state.loading) {
            commit('startLoad');
            dispatch('endLoadAfterDelay');
        }
    },
    endLoadAfterDelay({commit}) {
        setTimeout(() => {
            commit('endLoad');
        }, 1000);
    },
    pushMessage({commit, state}, description) {
        let id = descriptionId++;
        commit('pushMessage', {id, description});
        return id;
    }
}

export const getters = {
    progress(state) {
        let total = 0;
        for (let i in state.progresses) {
            total += state.progresses[i] * state.progressWeights[i];
        }
        return total;
    }
}

function safeSubtract(val) {
    if (val == 0) {
        return 0;
    }
    return val - 1;
}
