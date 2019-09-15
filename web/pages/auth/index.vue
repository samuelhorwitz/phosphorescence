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
                <button @click="login" autofocus>Login With Spotify</button>
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

    button:focus,
    button:active {
        outline: none;
        box-shadow: 5px 5px teal;
    }

    h2 {
        margin-top: 0px;
        margin-bottom: 1em;
        font-size: 2.2em;
    }

    form {
        display: flex;
        flex-direction: column;
    }

    label {
        display: flex;
        align-items: center;
        cursor: pointer;
    }

    input[type=checkbox] {
        appearance: none;
        background-color: white;
        border: 10px inset aqua;
        width: 5em;
        min-width: 5em;
        height: 5em;
        margin: 0px;
        margin-right: 1em;
        border-radius: 0px;
        cursor: pointer;
    }

    input[type=checkbox]:checked::after {
        content: url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB3aWR0aD0iODUuNTg1OTM3NSIgaGVpZ2h0PSIxMzEuOTA2MjUiPjxwYXRoIGZpbGw9Im1hZ2VudGEiIHN0cm9rZT0ibGlnaHRjeWFuIiBkPSJNMzQuNTYgMTEwLjMwTDM0LjU2IDExMC4zMFEyOS40OCAxMDEuMTIgMjYuNzQgOTYuODBMMjYuNzQgOTYuODBMMjMuODQgOTIuMjZMMjEuNjYgODguOTJRMTQuNDkgNzcuNzcgNi43MyA3MC40OUw2LjczIDcwLjQ5UTExLjE2IDY2LjgzIDE0Ljk4IDY2LjgzTDE0Ljk4IDY2LjgzUTE5LjY5IDY2LjgzIDIzLjM1IDcwLjI3UTI3LjAyIDczLjcyIDMyLjU5IDgzLjQwTDMyLjU5IDgzLjQwUTM4Ljk0IDYyLjc4IDQ4LjczIDQ0LjMwTDQ4LjczIDQ0LjMwUTU0LjE0IDM0LjIzIDU4LjM4IDMwLjc2UTYyLjYyIDI3LjI5IDY5LjYyIDI3LjI5TDY5LjYyIDI3LjI5UTczLjM0IDI3LjI5IDc4LjgwIDI4LjQ0TDc4LjgwIDI4LjQ0UTY0Ljc1IDM5Ljc2IDU2LjMzIDU1LjM0UTQ3LjkxIDcwLjkzIDM0LjU2IDExMC4zMFoiLz48L3N2Zz4K');
        font-family: 'Caveat';
        font-size: 10em;
        position: relative;
        top: -0.5em;
        left: -0.1em;
        color: magenta;
        text-shadow: -1px -1px 0 lightcyan, 1px -1px 0 lightcyan, -1px 1px 0 lightcyan, 1px 1px 0 lightcyan;
    }

    input[type=checkbox]:focus,
    input[type=checkbox]:active {
        outline: none;
        box-shadow: 5px 5px teal;
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
        showLogo: false,
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