<template>
    <article :class="{loading: !$store.getters['tracks/playlistLoaded']}">
        <div class="tableWrapper" :class="{loading: $store.state.loading.playlistGenerating}" ref="tableWrapper" v-if="$store.getters['tracks/playlistLoaded'] && !showConstellation">
            <table>
                <thead>
                    <tr>
                        <th class="playButton">
                        </th>
                        <th class="number">
                        </th>
                        <th>
                            Title
                        </th>
                        <th>
                            Artist
                        </th>
                        <th class="album">
                            Album
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <tr
                        v-for="(track, index) in $store.state.tracks.playlist"
                        ref="playlistTrack"
                        :class="{currentTrack: isPlaying(track.id), selectedTrack: isSelected(track.id)}"
                        @click="selectTrack(index); $ga.event('playlist', 'click', 'track', index)"
                        @dblclick="seekTrack(index); $ga.event('playlist', 'double-click', 'track', index)"
                        tabindex="0"
                        @keydown.arrow-up="moveCursorUp(); $ga.event('playlist', 'key', 'up')"
                        @keydown.arrow-down="moveCursorDown(); $ga.event('playlist', 'key', 'down')"
                        @keydown.enter="seekTrack(index); $ga.event('playlist', 'key', 'enter')"
                        v-spotify-uri:track="track.id"
                        v-spotify-uri-title="getSpotifyTrackDragTitle(track)">
                        <td class="playButton">
                            <button @click.stop="seekTrack(index); $ga.event('playlist', 'click', 'play', index)" :disabled="$store.state.tracks.previewLocked || ($store.getters['tracks/isPlayerDisconnected'] && !previewUrls[index])" class="playButton" v-if="$store.state.tracks.currentPreview != track.id && !(isPlaying(track.id) && $store.getters['tracks/playing'])">
                                <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">{{playButtonText}}</title><path d="M3,0.25V31.71L30.25,16ZM5,3.71L26.25,16,5,28.24V3.71Z"></path></svg>
                            </button>
                            <button @click.stop="handlePreviewStop(); $ga.event('playlist', 'click', 'stop', index)" class="stopButton" v-if="$store.state.tracks.currentPreview == track.id" @mouseover="previewPlaybackShadowColor = 'cyan'" @mouseout="previewPlaybackShadowColor = 'magenta'">
                                <previewPlayback :percent="$store.state.tracks.currentPreviewPercent" :shadowColor="previewPlaybackShadowColor"></previewPlayback>
                                <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px"aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Stop Preview</title><path d="M1,1V31H31V1H1ZM29,29H3V3H29V29Z"></path></svg>
                            </button>
                            <button @click.stop="pause(); $ga.event('playlist', 'click', 'pause', index)" class="currentlyPlaying" v-if="isPlaying(track.id) && $store.getters['tracks/playing']">
                                <svg id="currentlyPlaying" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0px" y="0px" viewBox="0 0 128 128" style="enable-background:new 0 0 128 128;" xml:space="preserve" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Currently Playing Track</title><g><path d="M72.3,16.1L39.1,39.8H20.5c-1.3,0-2.4,1.1-2.4,2.4v43.5c0,1.3,1.1,2.4,2.4,2.4h18.6l33.2,23.7c0.4,0.3,0.9,0.4,1.4,0.4   c0.4,0,0.8-0.1,1.1-0.3c0.8-0.4,1.3-1.2,1.3-2.1V18.1c0-0.9-0.5-1.7-1.3-2.1S73,15.6,72.3,16.1z M22.9,44.7h14.5v38.7H22.9V44.7z    M71.2,105.2l-29-20.7v-41l29-20.7V105.2z"></path><path d="M97.6,90.3c0.6,0,1.2-0.2,1.7-0.7c6.8-6.9,10.6-16,10.6-25.6s-3.8-18.8-10.6-25.6c-0.9-0.9-2.5-0.9-3.4,0   c-0.9,0.9-0.9,2.5,0,3.4c5.9,6,9.2,13.8,9.2,22.2s-3.3,16.3-9.2,22.2c-0.9,0.9-0.9,2.5,0,3.4C96.3,90.1,97,90.3,97.6,90.3z"></path><path d="M87.3,81.1c0.5,0.5,1.1,0.7,1.7,0.7c0.6,0,1.2-0.2,1.7-0.7c4.6-4.5,7.1-10.6,7.1-17.1s-2.5-12.5-7.1-17.1   c-0.9-0.9-2.5-0.9-3.4,0c-0.9,0.9-0.9,2.5,0,3.4C91,54,93,58.8,93,64s-2,10-5.7,13.7C86.4,78.6,86.4,80.1,87.3,81.1z"></path></g></svg>
                                <svg id="pauseTrack" xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Pause Track</title><path d="M13.76,1H6.63V31h7.14V1Zm-2,28H8.63V3h3.14V29Z"></path><path d="M25.37,1H18.24V31h7.14V1Zm-2,28H20.24V3h3.14V29Z"></path></svg>
                            </button>
                        </td>
                        <td :title="humanReadableEvocativeness[index]" class="number">
                            <span>{{index + 1}}</span>
                        </td>
                        <td :title="track.track.name"><a target="_blank" rel="external noopener" :href="getSpotifyTrackUrl(track.id)" @click.stop v-spotify-uri:track="track.id" v-spotify-uri-title="getSpotifyTrackDragTitle(track)">{{track.track.name}}</a></td>
                        <td :title="track.track.artists.map(artist => artist.name).join(', ')">
                            <ol>
                                <li class="artist" v-for="artist in track.track.artists">
                                    <a target="_blank" rel="external noopener" :href="getSpotifyArtistUrl(artist.id)" @click.stop v-spotify-uri:artist="artist.id" v-spotify-uri-title="artist.name">{{artist.name}}</a>
                                </li>
                            </ol>
                        </td>
                        <td class="album" :title="track.track.album.name"><a target="_blank" rel="external noopener" :href="getSpotifyAlbumUrl(track.track.album.id)" @click.stop v-spotify-uri:album="track.track.album.id" v-spotify-uri-title="getSpotifyAlbumDragTitle(track.track.album)">{{track.track.album.name}}</a></td>
                    </tr>
                </tbody>
            </table>
        </div>
        <constellation class="constellation" :class="{loading: $store.state.loading.playlistGenerating}" v-if="$store.getters['tracks/playlistLoaded'] && showConstellation"></constellation>
        <loadingScreen v-if="!$store.getters['tracks/playlistLoaded']"></loadingScreen>
    </article>
</template>

<style scoped>
    a {
        color: white;
        text-decoration: none;
    }

    a:hover {
        text-decoration: underline;
    }

    article {
        display: flex;
        align-items: center;
        flex-direction: column;
        height: 100%;
        margin: 0px 2em;
        width: 100%;
    }

    article.loading {
        margin: 0px;
    }

    ol {
        list-style-type: none;
        margin: 0px;
        padding: 0px;
        display: inline;
    }

    li.artist {
        display: inline;
    }

    li.artist:not(:last-child)::after {
        content: ', ';
        display: inline;
    }

    td.playButton button {
        visibility: hidden;
        appearance: none;
        background-color: transparent;
        border: 0px;
        margin: 0px;
        padding: 0px;
        display: inline;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        position: relative;
    }

    tr:hover td.playButton button,
    td.playButton button.stopButton,
    td.playButton button.currentlyPlaying {
        visibility: visible;
    }

    td.playButton button:focus {
        outline: none;
    }

    td.playButton button svg {
        width: 1.5em;
        fill: white;
        stroke: white;
        stroke-linejoin: round;
    }

    td.playButton button.currentlyPlaying svg {
        fill: aquamarine;
        stroke: aquamarine;
    }

    td.playButton button.currentlyPlaying svg#currentlyPlaying {
        width: 2em;
    }

    td.playButton button.currentlyPlaying svg#pauseTrack {
        display: none;
    }

    td.playButton button:disabled svg {
        fill: gray !important;
        stroke: gray !important;
        cursor: not-allowed;
    }

    td.playButton button.stopButton svg {
        position: absolute;
        fill: aquamarine;
        stroke: magenta;
    }

    td.playButton button:hover svg {
        fill: aquamarine;
        stroke: magenta;
    }

    td.playButton button.stopButton:hover svg {
        fill: aquamarine;
        stroke: aquamarine;
    }

    td.playButton button.currentlyPlaying:hover svg#currentlyPlaying {
        display: none;
    }

    td.playButton button.currentlyPlaying:hover svg#pauseTrack {
        display: block;
    }

    .tracks {
        display: flex;
        flex: 1;
        justify-content: center;
        height: 100%;
    }

    .tableWrapper {
        overflow: auto;
        border: 7px outset aquamarine;
        box-sizing: border-box;
        height: 100%;
        transition: transform 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275) 0s, opacity 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275) 0s;
    }

    .constellation {
        transition: transform 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275) 0s, opacity 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275) 0s;
    }

    .tableWrapper.loading,
    .constellation.loading {
        transform: scale(0.7);
        opacity: 0.3;
        pointer-events: none;
    }

    table {
        width: 100%;
        height: 100%;
        border-spacing: 0px;
        border-collapse: unset;
        table-layout: fixed;
        font-size: 16px;
    }

    th {
        font-weight: bold;
    }

    th.number {
        width: 3em;
    }

    th.playButton {
        width: 3em;
    }

    tr {
        color: white;
        background-color: rgb(15,10,54);
    }

    tr:focus {
        outline: none;
    }

    td {
        padding: 0.5em;
    }

    tr:not(:last-child) td {
        border-bottom: 1px solid rgba(245,188,251,.8);
    }

    td:not(.number) {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    td.playButton {
        padding: 0px;
    }

    tr.currentTrack, tr.selectedTrack, tr:hover {
        background-color: transparent !important;
    }

    tr:hover {
        cursor: pointer;
    }

    thead tr:nth-child(1) th{
        background-color: rgb(15,10,54);
        position: sticky;
        top: 0;
        z-index: 10;
        padding: 0.5em 0;
        border-bottom: 1px solid rgba(245,188,251,.8);
    }

    @media only screen and (min-width: 1199px), (max-width: 1099px) and (min-height: 550px) {
        article {
            margin: 0px 2em;
        }
    }

    @media only screen and (max-height: 449px) {
        article {
            margin: 0 1em;
        }
    }

    @media only screen and (max-height: 274px) {
        .tableWrapper {
            display: none;
        }
    }

    @media only screen and (max-height: 999px) {
        .tracks {
            align-items: center;
        }
    }

    @media only screen and (max-width: 1099px) {
        body:not(.playerConnected) th.album, body:not(.playerConnected) td.album {
            display: none;
        }
    }

    @media only screen and (max-width: 1499px) {
        body.playerConnected th.album, body.playerConnected td.album {
            display: none;
        }
    }

    @media only screen and (max-width: 399px) {
        th.playButton, td.playButton {
            display: none;
        }

        td.number {
            text-align: right;
        }
    }
</style>

<script>
    import {getSpotifyAlbumUrl, getSpotifyArtistUrl, getSpotifyTrackUrl, getSpotifyTrackDragTitle, getSpotifyAlbumDragTitle} from '~/assets/spotify';
    import constellation from '~/components/constellation';
    import loadingScreen from '~/components/loading-screen';
    import previewPlayback from '~/components/preview-playback';

    export default {
        components: {constellation, loadingScreen, previewPlayback},
        data() {
            return {
                previewPlaybackShadowColor: 'magenta'
            };
        },
        watch: {
            currentTrack() {
                let playingEl = this.$el.querySelector('.tableWrapper .playing');
                if (playingEl) {
                    this.$refs.tableWrapper.scrollTop = playingEl.offsetTop;
                }
            },
            showConstellation(newVal, oldVal) {
                if (oldVal && !newVal) {
                    setTimeout(() => {
                        // This is purely to fix a weird rendering bug in Desktop Safari
                        // MacOS 10.14.6, Safari 13.0.1
                        if (this.$refs.tableWrapper) {
                            this.$refs.tableWrapper.style.display = 'none';
                            this.$refs.tableWrapper.offsetHeight;
                            this.$refs.tableWrapper.style.display = '';
                        }

                        // This is for focusing the correct tr
                        if (this.$refs.playlistTrack && this.$refs.playlistTrack[this.$store.state.tracks.selectedTrackCursor]) {
                            this.$refs.playlistTrack[this.$store.state.tracks.selectedTrackCursor].focus();
                        }
                    }, 10);
                }
            }
        },
        computed: {
            previewUrls() {
                let previewUrls = [];
                for (let track of this.$store.state.tracks.playlist) {
                    previewUrls.push(this.$store.state.tracks.previews[track.id]);
                }
                return previewUrls;
            },
            playButtonText() {
                if (this.$store.getters['tracks/isPlayerDisconnected']) {
                    return 'Play Preview';
                }
                return 'Play Track';
            },
            showConstellation() {
                return this.$store.state.preferences.showConstellation;
            },
            currentTrack() {
                if (this.$store.getters['tracks/isPlayerDisconnected']) {
                    return null;
                }
                return this.$store.getters['tracks/currentTrack'];
            },
            selectedTrack() {
                return this.$store.getters['tracks/selectedTrack'];
            },
            humanReadableEvocativeness() {
                return this.$store.state.tracks.playlist.map(({evocativeness}) => {
                    let str = '';
                    if (evocativeness.aetherealness >= 0.5) {
                        str += `${Math.floor((evocativeness.aetherealness - 0.5) * 200)}% aethereal and `
                    } else {
                        str += `${Math.floor((0.5 - evocativeness.aetherealness) * 200)}% chthonic and `
                    }
                    if (evocativeness.primordialness >= 0.5) {
                        str += `${Math.floor((evocativeness.primordialness - 0.5) * 200)}% primordial`
                    } else {
                        str += `${Math.floor((0.5 - evocativeness.primordialness) * 200)}% transcendental`
                    }
                    return str;
                });
            }
        },
        methods: {
            isPlaying(id) {
                if (!this.currentTrack) {
                    return false;
                }
                return this.currentTrack.id == id;
            },
            isSelected(id) {
                if (!this.selectedTrack) {
                    return false;
                }
                return this.selectedTrack.id == id;
            },
            selectTrack(i) {
                this.$store.commit('tracks/selectTrack', i);
            },
            seekTrack(i) {
                if (!this.$store.getters['tracks/isPlayerDisconnected']) {
                    this.$store.dispatch('tracks/seekTrack', i);
                    this.play();
                } else {
                    this.$store.commit('tracks/seekTrack', i);
                    this.$store.commit('tracks/playPreview', this.$store.state.tracks.playlist[i].id);
                }
                this.$refs.playlistTrack[i].focus();
            },
            moveCursorUp() {
                this.$store.commit('tracks/selectPreviousTrack');
                this.$refs.playlistTrack[this.$store.state.tracks.selectedTrackCursor].focus();
            },
            moveCursorDown() {
                this.$store.commit('tracks/selectNextTrack');
                this.$refs.playlistTrack[this.$store.state.tracks.selectedTrackCursor].focus();
            },
            handleKeyPress(e) {
                if (e.keyCode === 67) {
                    this.$ga.event('playlist', 'toggle-constellation');
                    this.$store.commit('preferences/toggleConstellation');
                }
            },
            play() {
                this.$store.dispatch('tracks/play');
            },
            pause() {
                this.$store.dispatch('tracks/pause');
            },
            handlePreviewStop() {
                this.$store.commit('tracks/stopPreview');
            },
            getSpotifyAlbumUrl,
            getSpotifyArtistUrl,
            getSpotifyTrackUrl,
            getSpotifyTrackDragTitle,
            getSpotifyAlbumDragTitle
        },
        mounted() {
            document.addEventListener('keydown', this.handleKeyPress);
        },
        beforeDestroy() {
            document.removeEventListener('keydown', this.handleKeyPress);
        },
        head() {
            return {
                title: 'Phosphorescence | Build trance and chill-out playlists for Spotify',
                meta: [
                    {
                        hid: 'og:type',
                        name: 'og:type',
                        content: 'website'
                    },
                    {
                        hid: 'og:site_name',
                        name: 'og:site_name',
                        content: 'Phosphorescence'
                    },
                    {
                        hid: 'og:image',
                        name: 'og:image',
                        content: 'https://phosphor.me/og.png'
                    },
                    {
                        hid: 'og:description',
                        name: 'og:description',
                        content: 'Build trance and chill-out playlists for Spotify'
                    },
                    {
                        hid: 'og:url',
                        name: 'og:url',
                        content: 'https://phosphor.me'
                    },
                    {
                        hid: 'og:title',
                        name: 'og:title',
                        content: 'Phosphorescence'
                    }
                ]
            };
        }
    };
</script>