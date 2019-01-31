export const state = () => ({
    script: '',
    dimensions: null,
    numberOfTracks: 10
});

export const mutations = {
    save(state, script) {
        sessionStorage.setItem('ide/editorBackup', script);
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
        let oldScript = sessionStorage.getItem('ide/editorBackup');
        let oldDimensions = sessionStorage.getItem('ide/scriptDimensions');
        let oldNumberOfTracks = sessionStorage.getItem('ide/numberOfTracks');
        if (oldScript) {
            state.script = oldScript;
        }
        if (oldDimensions) {
            state.dimensions = JSON.parse(oldDimensions);
        }
        if (oldNumberOfTracks) {
            state.numberOfTracks = parseInt(oldNumberOfTracks, 10);
        }
    }
};
