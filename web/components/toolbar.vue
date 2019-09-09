<template>
    <div class="container">
        <div class="wrapper">
            <div class="bg"></div>
            <menu>
                <li class="menuItem regenerateOrCancel" @click="regenerateOrCancel()">
                    <button v-if="!$store.getters['tracks/playlistLoaded'] || !$store.state.loading.playlistGenerating" :disabled="advancedOpen || $store.state.loading.loading || $nuxt.$route.path !== '/'">Generate New</button>
                    <button v-if="$store.getters['tracks/playlistLoaded'] && $store.state.loading.playlistGenerating">Cancel</button>
                </li>
                <li class="menuItem toPage" :class="{active: advancedOpen}" @click="toggledAdvanced()">
                    <button>Advanced</button>
                </li>
                <li class="menuItem toPage" @click="flash()">
                    <button><nuxt-link to="/settings">Settings</nuxt-link></button>
                </li>
                <li class="menuItem logout" @click="logout()">
                    <button>Logout</button>
                </li>
            </menu>
        </div>
        <aside class="advancedWrapper" v-if="advancedOpen">
            <div class="bg bg-inverse"></div>
            <ul class="advanced">
                <li class="advancedMenuItem">
                    <label for="tracksPerPlaylist">Tracks</label>
                    <input name="tracksPerPlaylist" type="number" v-model="tracksPerPlaylist">
                </li>
                <li class="advancedMenuItem">
                    <label for="seedStyle">Style</label>
                    <select name="seedStyle" v-model="seedStyle">
                        <option :value="null">Random</option>
                        <option :value="'chillout'">Chillout</option>
                        <option :value="'primal'">Primal</option>
                        <option :value="'melancholy'">Melancholy</option>
                        <option :value="'emotional'">Emotional</option>
                        <option :value="'progressive'">Progressive</option>
                        <option :value="'hard'">Hard</option>
                        <option :value="'dark'">Dark</option>
                        <option :value="'trippy'">Trippy</option>
                    </select>
                </li>
                <li class="advancedMenuItem">
                    <label for="onlyTheHits">Only Hits</label>
                    <input name="onlyTheHits" type="checkbox" v-model="onlyTheHits">
                </li>
                <li class="advancedMenuItem" @click="regenerateOrCancel()">
                    <button v-if="!$store.state.loading.playlistGenerating" :disabled="$store.state.loading.loading || $nuxt.$route.path !== '/'">Generate!</button>
                    <button v-if="$store.state.loading.playlistGenerating">Cancel</button>
                </li>
            </ul>
        </aside>
    </div>
</template>

<style scoped>
    .container {
        display: flex;
        flex-direction: column;
        grid-column: 1 / 3;
    }

    .advanced {
        padding: 0px;
        margin: 0px;
        margin-left: 7vw;
        margin-right: 7vw;
        display: flex;
        align-items: center;
        justify-content: flex-start;
        flex: 1;
        z-index: 1;
    }

    .advancedWrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 2em;
        margin-top: 1em;
    }

    .advancedMenuItem {
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: bold;
        color: white;
        margin-right: 2em;
    }

    menu {
        padding: 0px;
        margin: 0px;
        margin-bottom: 1em;
        margin-left: 7vw;
        margin-right: 7vw;
        display: flex;
        align-items: center;
        justify-content: flex-start;
        flex: 1;
    }

    .menuItem {
        display: inline;
        cursor: pointer;
        margin-right: 3em;
        white-space: nowrap;
    }

    .menuItem:last-child {
        margin-left: auto;
    }

    button {
        appearance: none;
        border: 0px;
        background: transparent;
        color: white;
        font-family: 'Caveat';
        font-size: 2.5em;
        outline: none;
        transform: rotate(-10deg);
        cursor: pointer;
        text-shadow: -1px -1px 0 midnightblue, 1px -1px 0 midnightblue, -1px 1px 0 midnightblue, 1px 1px 0 midnightblue;
    }

    button[disabled] {
        cursor: not-allowed;
    }

    .menuItem:hover button, .menuItem.active button, .advancedMenuItem:hover button {
        color: magenta;
        text-shadow: -1px -1px 0 lightcyan, 1px -1px 0 lightcyan, -1px 1px 0 lightcyan, 1px 1px 0 lightcyan;
    }

    .menuItem:hover button[disabled] {
        text-decoration: line-through;
    }

    .bg {
        background-color: mediumvioletred;
        height: 2em;
        transform: skewX(-36deg);
        position: absolute;
        width: 90vw;
        border: 3px outset mediumturquoise;
    }

    .bg-inverse {
        background-color: mediumvioletred;
        border-color: mediumvioletred;
        transform: skewX(36deg);
    }

    .wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
    }

    a {
        color: inherit;
        text-decoration: none;
    }

    label {
        font-family: 'Montserrat';
        font-size: 1.3em;
        white-space: nowrap;
        margin-right: 0.5em;
    }

    input[type="number"] {
        -webkit-appearance: none;
        border: 3px inset gray;
        height: 1em;
        font-size: 1em;
        text-align: right;
        outline: 0;
        width: 3em;
    }

    select {
        height: 3.2em;
        border-radius: 0px;
        border: 3px outset gray;
        background-color: dimgray;
    }

    @media only screen and (max-width: 329px) {
        .container {
            font-size: 7.5px;
        }

        .menuItem {
            margin-right: 0.5em;
        }

        .advancedWrapper {
            font-size: 8px;
        }

        .advancedMenuItem label {
            font-size: 5.5px;
        }

        .advancedMenuItem select {
            height: 1.75em;
        }
    }

    @media only screen and (min-width: 330px) and (max-width: 419px) {
        .container {
            font-size: 7px;
        }

        .menuItem {
            margin-right: 0.5em;
        }

        .advancedWrapper {
            font-size: 8px;
        }

        .advancedMenuItem label {
            font-size: 6.5px;
        }

        .advancedMenuItem select {
            height: 2em;
        }
    }

    @media only screen and (max-height: 449px) {
        .container {
            grid-column: 2 / 3;
            grid-row: 1 / 2;
        }

        .wrapper {
            position: relative;
            width: 100%;
            margin-top: 0.5em;
        }

        .bg {
            width: 90%;
        }

        .menuItem {
            margin-right: 2em;
        }

        .menuItem.toPage {
            display: none;
        }

        .advancedWrapper {
            display: none;
        }
    }

    @media only screen and (max-width: 949px) and (max-height: 449px) {
        .bg {
            width: 80%;
        }
    }

    @media only screen and (max-width: 749px) and (max-height: 449px) {
        .container {
            grid-column: 1 / 3;
        }
    }

    @media only screen and (max-width: 524px) and (max-height: 449px) {
        .menuItem.logout {
            display: none;
        }

        .menuItem {
            margin: 0;
        }

        menu {
            justify-content: center;
        }

        .bg {
            width: 75%;
        }
    }

    @media only screen and (min-height: 450px) and (max-width: 1099px) {
        .container {
            grid-column: 1 / 2;
        }
    }

    @media only screen and (max-height: 449px) and (min-width: 950px) and (max-width: 1099px) {
        body.playerConnected .menuItem.logout {
            display: none;
        }
    }

    @media only screen and (min-width: 420px) and (max-width: 599px) and (min-height: 450px) {
        .container {
            font-size: 12px;
        }

        .menuItem {
            margin-right: 0.5em;
        }
    }

    @media only screen and (min-width: 600px) and (max-width: 899px) and (min-height: 450px) {
        .container {
            font-size: 16px;
        }

        .menuItem {
            margin-right: 0.5em;
        }
    }
</style>

<script>
    import {logout} from '~/assets/session';
    import {builders, loadNewPlaylist} from '~/assets/recordcrate';
    import {terminatePlaylistBuilding} from '~/assets/eos';

    export default {
        data() {
            return {
                advancedOpen: false
            };
        },
        computed: {
            seedStyle: {
                get() {
                    return this.$store.state.preferences.seedStyle;
                },
                set(newValue) {
                    this.$store.commit('preferences/updateSeedStyle', newValue);
                }
            },
            tracksPerPlaylist: {
                get() {
                    return this.$store.state.preferences.tracksPerPlaylist;
                },
                set(newValue) {
                    this.$store.commit('preferences/updateTracksPerPlaylist', newValue);
                }
            },
            onlyTheHits: {
                get() {
                    return this.$store.state.preferences.onlyTheHits;
                },
                set(newValue) {
                    this.$store.commit('preferences/updateOnlyTheHits', newValue);
                }
             }
        },
        methods: {
            regenerateOrCancel() {
                if (this.$store.state.loading.playlistGenerating) {
                    terminatePlaylistBuilding();
                } else {
                    this.regenerate();
                }
            },
            async regenerate() {
                if (!this.$store.getters['tracks/stopped'] && !confirm('This will destroy the current playlist. Are you sure?')) {
                    return;
                }
                this.$store.commit('loading/startLoad');
                let loadingMessage = 'Regenerating playlist';
                if (this.$store.state.preferences.seedStyle) {
                    loadingMessage += ` (${this.$store.state.preferences.seedStyle})`;
                } else {
                    loadingMessage += ' (random)';
                }
                let messageId = await this.$store.dispatch('loading/pushMessage', loadingMessage);
                this.$store.commit('loading/resetProgress');
                this.$store.commit('loading/initializeProgress', {id: 'generate'});
                this.$store.commit('loading/playlistGenerating');
                let pruners;
                if (this.$store.state.preferences.onlyTheHits) {
                    pruners = [builders.hits];
                }
                try {
                    let {playlist} = await loadNewPlaylist(this.$store.state.preferences.tracksPerPlaylist, builders.randomwalk, builders[this.$store.state.preferences.seedStyle], null, pruners, percent => {
                        this.$store.commit('loading/tickProgress', {id: 'generate', percent});
                    });
                    this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                }
                catch (e) {
                    console.error('Playlist generation failed', e);
                    // TODO add some visual UI indication
                }
                this.$store.commit('loading/completeProgress', {id: 'generate'});
                this.$store.commit('loading/resetProgress');
                this.$store.commit('loading/clearMessage', messageId);
                this.$store.commit('loading/playlistGenerationComplete');
                this.$store.dispatch('loading/endLoadAfterDelay');
            },
            logout() {
                this.$store.dispatch('loading/loadFlash');
                logout();
                this.$router.push('/auth');
            },
            flash() {
                this.$store.dispatch('loading/loadFlash');
            },
            toggledAdvanced() {
                this.flash();
                this.advancedOpen = !this.advancedOpen;
            }
        }
    };
</script>