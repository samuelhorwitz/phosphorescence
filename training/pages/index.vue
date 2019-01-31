<template>
    <main>
        <h1>Welcome To Phosphorescence ML Training</h1>
        <playlist></playlist>
        <player :deviceId="deviceId"></player>
        <button @click="testModels" :disabled="testingModels">Test models</button>
        <axis v-if="!testingModels"></axis>
        <model-prediction v-if="testingModels"></model-prediction>
    </main>
</template>

<script>
    import playlist from '~/components/playlist';
    import player from '~/components/player';
    import axis from '~/components/axis';
    import modelPrediction from '~/components/modelPrediction';
    import {getLocalTokens, tryToRefresh} from '~/util/spotify';

    export default {
        components: {
            axis,
            playlist,
            player,
            modelPrediction
        },
        data() {
            return {
                deviceId: null,
                testingModels: false
            };
        },
        async created() {
            if (!process.browser) {
                return;
            }
            let {access, refresh, expires, okay} = getLocalTokens();
            if (!okay) {
                let response = await tryToRefresh(refresh);
                access = response.access;
                refresh = response.refresh;
                expires = response.expires;
            }
            this.$store.dispatch('session/tokens', {access, refresh, expires});
            this.setUpSpotifyPlayer(access);
        },
        methods: {
            setUpSpotifyPlayer(token) {
                window.onSpotifyWebPlaybackSDKReady = () => {
                    let player = new Spotify.Player({
                        name: 'Phosphorescence Training Player',
                        getOAuthToken: cb => cb(token)
                    });
                    player.addListener('initialization_error', async ({ message }) => {
                        console.error(message);
                        let {access, refresh, expires} = await tryToRefresh();
                        this.$store.dispatch('session/tokens', {access, refresh, expires});
                        this.$router.push('/');
                    });
                    player.addListener('authentication_error', async ({ message }) => {
                        console.error(message);
                        let {access, refresh, expires} = await tryToRefresh();
                        this.$store.dispatch('session/tokens', {access, refresh, expires});
                        this.$router.push('/');
                    });
                    player.addListener('account_error', async ({ message }) => {
                        console.error(message);
                        let {access, refresh, expires} = await tryToRefresh();
                        this.$store.dispatch('session/tokens', {access, refresh, expires});
                        this.$router.push('/');
                    });
                    player.addListener('playback_error', ({ message }) => { console.error(message); });

                    // Playback status updates
                    player.addListener('player_state_changed', state => { console.log(state); });

                    // Ready
                    player.addListener('ready', ({ device_id }) => {
                        console.log('Ready with Device ID', device_id);
                        this.deviceId = device_id;
                    });

                    // Not Ready
                    player.addListener('not_ready', ({ device_id }) => {
                        console.log('Device ID has gone offline', device_id);
                    });

                    // Connect to the player!
                    player.connect();
                };

                let script = document.createElement('script');
                script.src = 'https://sdk.scdn.co/spotify-player.js';
                document.head.appendChild(script);
            },
            testModels() {
                this.testingModels = true;
            }
        }
    };
</script>