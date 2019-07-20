<template>
    <div class="wrapper" v-show="$store.getters['idetracks/playlistLoaded'] && $store.getters['idetracks/deviceLoaded']">
        <div class="container">
            <menu>
                <li><button @click="play" v-if="$store.getters['idetracks/stopped']">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title>Play Track</title><path d="M3,0.25V31.71L30.25,16ZM5,3.71L26.25,16,5,28.24V3.71Z"></path></svg>
                </button></li>
                <li><button @click="resume" v-if="$store.getters['idetracks/paused']">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title>Play Track</title><path d="M3,0.25V31.71L30.25,16ZM5,3.71L26.25,16,5,28.24V3.71Z"></path></svg>
                </button></li>
                <li><button @click="pause" v-if="$store.getters['idetracks/playing']">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title>Pause Track</title><path d="M13.76,1H6.63V31h7.14V1Zm-2,28H8.63V3h3.14V29Z"></path><path d="M25.37,1H18.24V31h7.14V1Zm-2,28H20.24V3h3.14V29Z"></path></svg>
                </button></li>
                <li><button @click="previous" :class="{disabled: !$store.getters['idetracks/canSkipBackward']}">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title>Previous Track</title><path d="M3.46,2h-2V30h2V16.08L30.54,31.71V0.25L3.46,15.89V2ZM28.54,3.72V28.25L7.3,16Z"></path></svg>
                </button></li>
                <li><button @click="next" :class="{disabled: !$store.getters['idetracks/canSkipForward']}">
                    <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Next Track</title><path d="M28.54,15.89L1.46,0.25V31.71L28.54,16.08V30h2V2h-2v13.9ZM3.46,28.25V3.72L24.7,16Z"></path></svg>
                </button></li>
            </menu>
            <a class="albumImgLink" target="_blank" :href="currentAlbumUrl"><img class="albumImg" :alt="currentTrackImageAltText" :src="currentTrackImage"></a>
            <div class="trackContainer">
                <div class="trackData">
                    <div class="marquee">
                        <span class="trackName">
                            <a target="_blank" :href="currentTrackUrl">{{currentTrackName}}</a>
                        </span>
                        <ol class="artistsNames">
                            <li v-for="artist in currentTrackArtists">
                                <a target="_blank" :href="artist.external_urls.spotify">{{artist.name}}</a>
                            </li>
                        </ol>
                        <span class="albumName">
                            <a target="_blank" :href="currentAlbumUrl">{{currentAlbumName}}</a>
                        </span>
                    </div>
                </div>
            </div>
            <div class="spotifyLogo">
                <a target="_blank" href="https://spotify.com"><img alt="Spotify" class="spotifyLogoImg" src="/images/spotify_small.png"></a>
            </div>
        </div>
    </div>
</template>

<style scoped>
    aside {
        width: 100%;
        text-align: center;
        position: fixed;
        bottom: 3.5em;
        display: flex;
        justify-content: center;
    }

    menu {
        margin: 0px;
        padding: 0px;
        cursor: pointer;
        padding-left: 1em;
    }

    menu li {
        display: inline;
    }

    .trackContainer {
        display: flex;
        flex: 1;
        overflow-y: hidden;
    }

    .wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: rgba(26, 17, 16, 1);
        border: 2px solid black;
        width: 100%;
        height: 3.2em;
        position: relative;
        flex-direction: column;
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
        flex: 1;
        white-space: nowrap;
        overflow-x: scroll;
        overflow-y: hidden;
        padding-right: 0.5em;
        padding-bottom: 2em;
        margin-left: .2em;
    }

    .trackData.stopped {
        overflow: hidden;
    }

    .trackData::-webkit-scrollbar { 
        display: none; 
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

    .container {
        display: flex;
        margin-top: auto;
        align-items: center;
        width: 100%;
    }

    .albumImg {
        max-height: 2.5em;
        border: 1px solid white;
    }

    .albumImgLink {
        padding-left: 1em;
    }

    .spotifyLogo {
        position: absolute;
        bottom: 0px;
        right: 0px;
    }

    .spotifyLogoImg {
        margin-left: auto;
        max-width: 70px;
    }

    .marquee {
        color: white;
        animation: marquee 15s linear infinite;
        padding-left: 100%;
        display: inline-block;
        cursor: pointer;
    }

    .marquee:hover {
        animation-play-state: paused;
    }

    @keyframes marquee {
        0%   { transform: translate(0, 0); }
        100% { transform: translate(-100%, 0); }
    }
</style>

<script>
    import {initializePlayer} from '~/assets/spotify';

    export default {
        data() {
            return {
                destroyer: null
            };
        },
        computed: {
            tracks() {
                return this.$store.state.idetracks.playlist;
            },
            currentTrack() {
                return this.$store.getters['idetracks/currentTrack'];
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
            currentTrackImage() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.album.images[track.track.album.images.length - 1].url;
            },
            currentTrackImageAltText() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return `${track.track.album.name} - ${track.track.album.artists.map(artist => artist.name).join(', ')}`;
            }
        },
        methods: {
            play() {
                this.$store.dispatch('idetracks/play');
            },
            resume() {
                this.$store.dispatch('idetracks/resume');
            },
            pause() {
                this.$store.dispatch('idetracks/pause');
            },
            next() {
                this.$store.dispatch('idetracks/next');
            },
            previous() {
                this.$store.dispatch('idetracks/previous');
            }
        },
        async created() {
            let playerWrapper = await initializePlayer(this.$store, 'idetracks');
            this.destroyer = playerWrapper.destroyer;
            this.$store.dispatch('idetracks/registerPlayer', playerWrapper.player);
        },
        beforeDestroy() {
            if (this.destroyer) {
                this.destroyer();
            }
        }
    };
</script>
