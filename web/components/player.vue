<template>
    <aside v-if="ready">
        <div class="wrapper">
            <menu>
                <li><button @click="play" v-if="$store.getters['tracks/stopped']" :class="{disabled: !ready}">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title>Play Track</title><path d="M3,0.25V31.71L30.25,16ZM5,3.71L26.25,16,5,28.24V3.71Z"></path></svg>
                </button></li>
                <li><button @click="resume" v-if="$store.getters['tracks/paused']" :class="{disabled: !ready}">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title>Play Track</title><path d="M3,0.25V31.71L30.25,16ZM5,3.71L26.25,16,5,28.24V3.71Z"></path></svg>
                </button></li>
                <li><button @click="pause" v-if="$store.getters['tracks/playing']" :class="{disabled: !ready}">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title>Pause Track</title><path d="M13.76,1H6.63V31h7.14V1Zm-2,28H8.63V3h3.14V29Z"></path><path d="M25.37,1H18.24V31h7.14V1Zm-2,28H20.24V3h3.14V29Z"></path></svg>
                </button></li>
                <li><button @click="previous" :class="{disabled: !ready || !$store.getters['tracks/canSkipBackward']}">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title>Previous Track</title><path d="M3.46,2h-2V30h2V16.08L30.54,31.71V0.25L3.46,15.89V2ZM28.54,3.72V28.25L7.3,16Z"></path></svg>
                </button></li>
                <li><button @click="next" :class="{disabled: !ready || !$store.getters['tracks/canSkipForward']}">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Next Track</title><path d="M28.54,15.89L1.46,0.25V31.71L28.54,16.08V30h2V2h-2v13.9ZM3.46,28.25V3.72L24.7,16Z"></path></svg>
                </button></li>
            </menu>
            <div class="trackData" :class="{stopped: $store.getters['tracks/stopped']}">
                <span class="trackName" v-if="!$store.getters['tracks/stopped']">
                    <a target="_blank" :href="currentTrackUrl">{{currentTrackName}}</a>
                </span>
                <ol class="artistsNames" v-if="!$store.getters['tracks/stopped']">
                    <li v-for="artist in currentTrackArtists">
                        <a target="_blank" :href="artist.external_urls.spotify">{{artist.name}}</a>
                    </li>
                </ol>
                <span class="albumName" v-if="!$store.getters['tracks/stopped']">
                    <a target="_blank" :href="currentAlbumUrl">{{currentAlbumName}}</a>
                </span>
                <span class="nothingPlaying" v-if="$store.getters['tracks/stopped'] && !$store.state.tracks.spotifyAppearsDown">
                    ... Welcome to Phosphorescence ... Please Click "Play" To Listen ... 💿 💻 ...
                </span>
                <span class="nothingPlaying" v-if="$store.getters['tracks/stopped'] && $store.state.tracks.spotifyAppearsDown">
                    ... 😢😢😢 Spotify's Playback API Appears To Be Down Right Now 😢 ...
                </span>
            </div>
        </div>
    </aside>
</template>

<style scoped>
    aside {
        width: 100%;
        height: 3em;
        text-align: center;
        display: flex;
        justify-content: center;
        grid-column: 1 / 3;
        margin-top: 1em;
    }

    menu {
        margin: 0px;
        padding: 0px;
        cursor: pointer;
        padding-left: .5em;
    }

    menu li {
        display: inline;
    }

    .wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: rgba(26, 17, 16, 0.9);
        border: 5px inset teal;
        border-radius: 10px;
        width: 90%;
        position: relative;
    }

    .trackName::after, .artistsNames::after {
        content: ' - ';
        color: white;
    }

    .nothingPlaying {
        color: white;
        animation: marquee 15s linear infinite;
        padding-left: 100%;
        display: inline-block;
        cursor: pointer;
    }

    .nothingPlaying:hover {
        animation-play-state: paused;
    }

    @keyframes marquee {
        0%   { transform: translate(0, 0); }
        100% { transform: translate(-100%, 0); }
    }

    a {
        color: white;
        text-decoration: none;
    }

    a:hover {
        text-decoration: underline;
    }

    ol {
        list-style-type: none;
        margin: 0px;
        padding: 0px;
        display: inline;
    }

    ol li {
        display: inline;
    }

    ol li:not(:last-child)::after {
        content: ', ';
        display: inline;
        color: white;
    }

    .trackData {
        font-family: VT323;
        font-size: 1.3em;
        flex: 1;
        white-space: nowrap;
        overflow-x: scroll;
        overflow-y: hidden;
        padding-left: 0.5em;
        padding-right: 0.5em;
        padding-bottom: 0.1em;
    }

    .trackData.stopped {
        overflow: hidden;
    }

    .trackData::-webkit-scrollbar { 
        display: none; 
    }

    .trackData::before {
        content: '';
        background: linear-gradient(90deg, rgba(26,17,16,0.9) 46%, rgba(255,255,255,0) 100%);
        width: 5px;
        position: absolute;
        top: 5px;
        left: 82px;
        height: 70%;
    }

    .trackData::after {
        content: '';
        background: linear-gradient(270deg, rgba(26,17,16,0.9) 46%, rgba(255,255,255,0) 100%);
        width: 5px;
        position: absolute;
        top: 5px;
        right: 0px;
        height: 70%;
    }

    button {
        appearance: none;
        border: 0px;
        background-color: transparent;
        margin: 0px;
        padding: 0px;
        outline: none;
    }

    button svg {
        width: 2em;
        fill: white;
        stroke: white;
        cursor: pointer;
        stroke-linejoin: round;
    }

    button svg:hover {
        fill: aquamarine;
        stroke: magenta;
    }

    button.disabled svg {
        fill: gray;
        stroke: gray;
    }

    button.disabled svg:hover {
        cursor: not-allowed;
    }

    @media only screen and (max-height: 449px) {
        aside {
            grid-column: 1 / 3;
            grid-row: 3 / 4;
            justify-content: center;
            margin: 0;
        }

        .wrapper {
            margin: 0 1em;
            width: 90%;
            flex: 1;
        }
    }

    @media only screen and (max-height: 249px) {
        aside {
            grid-row: 2 / 4;
        }
    }

    @media only screen and (min-height: 450px) and (max-width: 1099px) {
        aside {
            grid-column: 1 / 2;
            height: unset;
        }
    }
</style>

<script>
    import {initializePlayer} from '~/assets/spotify';

    export default {
        data() {
            return {
                ready: false,
                destroyer: null
            };
        },
        computed: {
            tracks() {
                return this.$store.state.tracks.playlist;
            },
            currentTrack() {
                return this.$store.getters['tracks/currentTrack'];
            },
            currentTrackName() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.name;
            },
            currentTrackUrl() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.external_urls.spotify;
            },
            currentAlbumName() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.album.name;
            },
            currentAlbumUrl() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.album.external_urls.spotify;
            },
            currentTrackArtists() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.artists;
            },
        },
        methods: {
            play() {
                this.$store.dispatch('tracks/play');
            },
            resume() {
                this.$store.dispatch('tracks/resume');
            },
            pause() {
                this.$store.dispatch('tracks/pause');
            },
            next() {
                this.$store.dispatch('tracks/next');
            },
            previous() {
                this.$store.dispatch('tracks/previous');
            }
        },
        async created() {
            this.$store.commit('loading/startLoad');
            try {
                this.destroyer = await initializePlayer(this.$store, 'tracks');
                this.ready = true;
            }
            catch (e) {
                this.ready = false;
            }
            this.$store.dispatch('loading/endLoadAfterDelay');
        },
        beforeDestroy() {
            if (this.destroyer) {
                this.destroyer();
            }
        }
    };
</script>