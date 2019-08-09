export const state = () => ({
    loading: true,
    loadSemaphore: 0,
    descriptions: [],
    progresses: {},
    progressWeights: {},
    progressIntervals: {}
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
    initializeProgress(state, {id, weight, intervalId}) {
        if (state.progressIntervals[id]) {
            clearInterval(state.progressIntervals[id]);
        }
        if (!weight) {
            weight = 100;
        }
        weight = Math.min(weight, 100);
        state.progresses[id] = 0;
        state.progressWeights[id] = weight;
        state.progressIntervals[id] = intervalId;
    },
    tickProgress(state, {id, amount}) {
        if (!amount) {
            amount = 1;
        }
        let old = state.progresses[id];
        state.progresses = {...state.progresses, [id]: Math.min(old + amount, 99)};
    },
    completeProgress(state, {id}) {
        state.progresses = {...state.progresses, [id]: 100};
        clearInterval(state.progressIntervals[id]);
    },
    resetProgress(state) {
        for (let i in state.progressIntervals) {
            clearInterval(state.progressIntervals[i]);
        }
        state.progresses = {};
        state.progressWeights = {};
        state.progressIntervals = {};
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
    },
    initializeProgress({state, commit}, {id, weight, ms, amount}) {
        let intervalId = setInterval(() => {
            commit('tickProgress', {id, amount});
        }, ms);
        commit('initializeProgress', {id, weight, intervalId});
    }
}

export const getters = {
    progress(state) {
        let total = 0;
        for (let i in state.progresses) {
            total += state.progresses[i] * (state.progressWeights[i] / 100)
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
