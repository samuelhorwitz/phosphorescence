import {getTokens, saveLocalTokens, getLocalTokens, sendToAuth} from '~/util/spotify';

export const state = () => ({
    accessToken: null
});

export const mutations = {
    tokens(state, {access}) {
        state.accessToken = access;
    }
};

export const actions = {
    tokens({commit, dispatch}, {access, refresh, expires}) {
        commit('tokens', {access, refresh});
        saveLocalTokens({access, refresh, expires});
        let {expiresAt} = getLocalTokens();
        return new Promise((resolve, reject) => {
            setTimeout(async () => {
                if (refresh) {
                    dispatch('tokens', await getTokens(refresh, true));
                } else {
                    sendToAuth();
                }
                resolve();
            }, expiresAt - Date.now());
        });
    }
};
