<template>
    <article>
        <div class="container">
            <h2><b class="appname">Phosphorescence</b>
            builds coherent playlists for <b class="genre">trance</b> and <b class="genre">chill-out</b> listeners.
            </h2>
            <form @submit.prevent.stop>
                <label for="rememberMe">
                    <input type="checkbox" id="rememberMe" v-model="rememberMe">
                    <span>Remember Me <em>(do not check if using a public computer)</em></span>
                </label>
                <button @click="login">Login With Spotify</button>
            </form>
            <div class="boilerplate">
                <p>We are an <a target="_blank" rel="external noopener" href="https://github.com/samuelhorwitz/phosphorescence">open source project</a> and 100% free to use.</p>
                <p>Spotify provides the music, however we are <em>not</em> affiliated with Spotify in any way.</p>
            </div>
        </div>
    </article>
</template>

<style scoped>
    article {
        margin: 0px 1em;
        font-size: 16px;
        padding: 1em;
        flex: 1;
        display: flex;
        justify-content: center;
        color: white;
    }

    .container {
        max-width: 40em;
    }

    a {
        color: cyan;
    }

    a:hover {
        color: magenta;
    }

    b.appname {
        font-variant: small-caps;
        font-weight: bolder;
        font-family: 'Varela';
    }

    b.genre {
        font-weight: inherit;
        color: aqua;
    }

    button {
        border: 7px outset aqua;
        background-color: magenta;
        -webkit-appearance: none;
        width: 100%;
        font-size: 2em;
        color: white;
        padding: 0.5em 0;
        cursor: pointer;
        margin-top: 0.5em;
    }

    button:hover {
        border-style: inset;
    }

    h2 {
        margin-top: 0px;
        font-size: 2.2em;
    }

    form {
        display: flex;
        flex-direction: column;
    }

    label {
        display: flex;
        align-items: center;
    }

    input[type=checkbox] {
        appearance: none;
        background-color: white;
        border: 10px inset aqua;
        width: 5em;
        min-width: 5em;
        height: 5em;
        margin: 0px;
        margin-right: 0.5em;
        border-radius: 0px;
    }

    input[type=checkbox]:checked {
        background-color: magenta;
    }

    .boilerplate {
        font-size: 0.75em;
        color: #eee;
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