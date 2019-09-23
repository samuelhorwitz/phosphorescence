export const state = () => ({
    user: null,
    country: null // only for non-logged-in user
});

export const mutations = {
    user(state, user) {
        state.user = user;
    },
    country(state, country) {
        country = country.toUpperCase();
        localStorage.setItem('country', country);
        state.country = country;
    },
    restore(state) {
        let country = localStorage.getItem('country');
        if (country && country !== 'null') {
            state.country = country;
        }
    }
};

export const getters = {
    country(state) {
        if (state.user) {
            return state.user.country;
        } else if (state.country) {
            return state.country;
        }
        return 'US';
    }
}
