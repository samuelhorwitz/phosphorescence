export const state = () => ({
    user: null
});

export const mutations = {
    user(state, user) {
        state.user = user;
    }
};
