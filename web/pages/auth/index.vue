<template>
    <article>
        <div>
            <p><span class="appname">Phosphorescence</span>
            is a playlist building tool for fans of trance music.
            </p>
            <p>We are an <a target="_blank" rel="external noopener" href="https://github.com/samuelhorwitz/phosphorescence">open source project</a> and 100% free to use.
            Spotify provides the music, however we are <em>not</em> affiliated with Spotify in any way.</p>
            <p><span class="appname">Phosphorescence</span> requires an active <a rel="external noopener" target="_blank" href="https://www.spotify.com">Spotify</a> account.
            Certain features, such as in-browser streaming, are only available to paid accounts.</p>
            <p>We will not store more information than necessary to provide you with this service, nor will we use or sell your information for advertising or other purposes.</p>
            <p>Please be sure to read our <nuxt-link to="/legal/tos">Terms of Service</nuxt-link> and <nuxt-link to="/legal/privacy">Privacy Policy</nuxt-link>.</p>
        </div>
        <div>
            <label for="rememberMe">
                <input type="checkbox" id="rememberMe" v-model="rememberMe">
                Remember Me <em>(do not check if using a public computer)</em>
            </label>
            <button @click="login">Login With Spotify</button>
        </div>
    </article>
</template>

<style scoped>
    article {
        margin: 1em;
        font-size: 16px;
        background-color: teal;
        padding: 1em;
        flex: 1;
        overflow-y: scroll;
        -webkit-overflow-scrolling: touch;
        border: 5px outset magenta;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }

    a {
        color: indigo;
    }

    a:hover {
        color: magenta;
    }

    .appname {
        font-variant: small-caps;
        font-weight: bolder;
        font-family: 'Varela';
    }

    button {
        border: 7px outset darkgray;
        background-color: gray;
        -webkit-appearance: none;
        width: 100%;
        font-size: 2em;
        color: black;
        padding: 0.5em 0;
        cursor: pointer;
        margin-top: 1em;
        margin-bottom: 1em;
    }

    button:hover {
        border-style: inset;
    }
</style>

<script>
    import {authorizeUserRedirect} from '~/assets/session';

    export default {
        layout: 'empty',
        middleware: 'unauthenticated',
        data() {
            return {
                rememberMe: false
            };
        },
        created() {
            this.$store.dispatch('loading/endLoadAfterDelay');
        },
        methods: {
            login() {
                this.$store.dispatch('loading/loadFlash');
                authorizeUserRedirect(this.rememberMe);
            }
        }
    };
</script>