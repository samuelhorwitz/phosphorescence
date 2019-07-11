<template>
    <editor></editor>
</template>

<script>
    import editor from '~/components/ide/editor';
    import {accessTokenExists, refreshUser} from '~/assets/session';

    export default {
        layout: 'ide',
        components: {
            editor
        },
        async fetch({store, params, error}) {
            if (!accessTokenExists()) {
                await refreshUser();
            }
            let {id: scriptId} = params;
            let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
            if (!userResponse.ok) {
                return error({statusCode: userResponse.status, message: "Could not get user information"});
            }
            let {user} = await userResponse.json();
            store.commit('ide/user', user);
            if (scriptId != null) {
                let scriptResponse = await fetch(`${process.env.API_ORIGIN}/script/${scriptId}`, {credentials: 'include'});
                if (!scriptResponse.ok) {
                    return error({statusCode: scriptResponse.status, message: "Could not get script"});
                }
                let {script} = await scriptResponse.json();
                store.commit('ide/scriptData', script);
                store.commit('ide/mostRecentVersion', script.mostRecent);
            }
            if (store.getters['ide/isScriptOwnedByUser']) {
                let draftResponse = await fetch(`${process.env.API_ORIGIN}/script/${scriptId}/version/draft?count=1`, {credentials: 'include'});
                if (!draftResponse.ok) {
                    return error({statusCode: draftResponse.status, message: "Could not get most recent draft version of script"});
                }
                let {scriptVersions} = await draftResponse.json();
                if (scriptVersions.length > 0) {
                    store.commit('ide/mostRecentVersion', scriptVersions[0]);
                }
            }
            if (store.state.ide.scriptVersionData) {
                let scriptResponse = await fetch(`${process.env.SCRIPTS_ORIGIN}/${store.state.ide.scriptVersionData.fileId}`);
                if (!scriptResponse.ok) {
                    return error({statusCode: scriptResponse.status, message: "Could not get associated script file"});
                }
                let initialValue = await scriptResponse.text();
                store.commit('ide/saveScript', initialValue);
            }
        }
    };
</script>