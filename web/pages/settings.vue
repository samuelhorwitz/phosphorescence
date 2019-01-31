<template>
    <div class="pageWrapper">
        <article class="authPage">
            <h2 class="pageHeader">Settings</h2>
            <p>You are logged in as the Spotify user {{name}}.</p>
            <p>Currently, all settings are stored locally on the device and automatically updated as they are changed.</p>
        </article>
    </div>
</template>

<style scoped>
    h2 {
        margin-bottom: 0.3em;
    }
</style>

<script>
    import {getAccessToken} from '~/assets/session';

    export default {
        data() {
            return {
                name: '█████ █████'
            };
        },
        async created() {
            let response = await fetch('https://api.spotify.com/v1/me', {
                method: 'GET',
                headers: {
                    Authorization: `Bearer ${await getAccessToken()}`
                }
            });
            let {display_name} = await response.json();
            setTimeout(() => {
                this.name = display_name;
            }, 500);
        }
    };
</script>