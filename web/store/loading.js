export const state = () => ({
    loading: true,
    loadSemaphore: 0
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
    }
};

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
    }
}

function safeSubtract(val) {
    if (val == 0) {
        return 0;
    }
    return val - 1;
}
