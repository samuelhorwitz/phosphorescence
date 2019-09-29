<template>
    <article>
        <div class="container">
            <form @submit.prevent.stop>
                <label for="rememberMe">
                    <input type="checkbox" id="rememberMe" v-model="rememberMe">
                    <span>Remember Me <em>(do not check if using a public computer)</em></span>
                </label>
                <button @click="login" autofocus>Login With Spotify</button>
            </form>
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

    button:focus {
        outline: none;
        box-shadow: 5px 5px teal;
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
        pointer-events: none;
    }

    input[type=checkbox]:focus {
        outline: none;
        box-shadow: 5px 5px teal;
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