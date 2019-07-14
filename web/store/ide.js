export const state = () => ({
    script: '',
    dimensions: null,
    numberOfTracks: 10,
    user: null,
    scriptData: null,
    scriptVersionData: null
});

export const mutations = {
    saveScript(state, script) {
        state.script = script;
    },
    dimensions(state, dimensions) {
        sessionStorage.setItem('ide/scriptDimensions', JSON.stringify(dimensions));
        state.dimensions = dimensions;
    },
    numberOfTracks(state, numberOfTracks) {
        sessionStorage.setItem('ide/numberOfTracks', numberOfTracks);
        state.numberOfTracks = numberOfTracks;
    },
    restore(state) {
        let oldDimensions = sessionStorage.getItem('ide/scriptDimensions');
        let oldNumberOfTracks = sessionStorage.getItem('ide/numberOfTracks');
        if (oldDimensions) {
            state.dimensions = JSON.parse(oldDimensions);
        }
        if (oldNumberOfTracks) {
            state.numberOfTracks = parseInt(oldNumberOfTracks, 10);
        }
    },
    user(state, user) {
        state.user = user;
    },
    scriptData(state, scriptData) {
        state.scriptData = scriptData;
    },
    mostRecentVersion(state, scriptVersionData) {
        state.scriptVersionData = scriptVersionData;
    }
};

export const getters = {
    isScriptOwnedByUser(state) {
        if (!state.user || !state.scriptData) {
            return false;
        }
        return state.user.spotifyId == state.scriptData.authorSpotifyId;
    },
}
