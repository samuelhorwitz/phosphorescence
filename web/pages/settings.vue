<template>
    <div class="pageWrapper">
        <article class="authPage">
            <h2 class="pageHeader">Settings</h2>
            <p>You are logged in as the Spotify user <a target="_blank" :href="'https://open.spotify.com/user/' + $store.state.user.user.spotifyId">{{$store.state.user.user.name}}</a>.</p>
            <p>Currently, all settings are stored locally on the device and automatically updated as they are changed.</p>
            <p>You may disconnect your Spotify account from this application at any time by visiting the <a target="_blank" href="https://www.spotify.com/us/account/apps/">Spotify user account page</a>.</p>
        </article>
    </div>
</template>

<style scoped>
    h2 {
        margin-bottom: 0.3em;
    }
</style>

<script>
    export default {
        async fetch({store, error}) {
            let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
            if (!userResponse.ok) {
                return error({statusCode: userResponse.status, message: "Could not get user information"});
            }
            let {user} = await userResponse.json();
            store.commit('user/user', user);
        }
    };
</script>