export const state = () => ({
    query: '',
    searchResults: null
});

export const mutations = {
    setQuery(state, query) {
        state.query = query;
    },
    setSearchResults(state, searchResults) {
        state.searchResults = searchResults;
    }
};
