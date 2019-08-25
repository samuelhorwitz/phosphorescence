export const state = () => ({
    query: ''
});

export const mutations = {
    setQuery(state, query) {
        state.query = query;
    },
    clearQuery(state) {
        state.query = '';
    }
};
