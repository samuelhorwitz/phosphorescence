<template>
    <article :class="{loading: !$store.getters['tracks/playlistLoaded']}">
        <div class="tableWrapper" :class="{loading: this.$store.state.loading.playlistGenerating}" ref="tableWrapper" v-if="$store.getters['tracks/playlistLoaded'] && !showConstellation">
            <table>
                <thead>
                    <tr>
                        <th class="playButton" v-if="playlistPlayButton">
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
                        <td class="playButton" v-if="playlistPlayButton">
                            <button @click.stop="seekTrack(index); play(); $ga.event('playlist', 'click', 'play', index)" class="playButton">
                                <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Play Track</title><path d="M3,0.25V31.71L30.25,16ZM5,3.71L26.25,16,5,28.24V3.71Z"></path></svg>
                            </button>
                        </td>
                        <td :title="humanReadableEvocativeness[index]" class="number" :class="{pullRight: !playlistPlayButton}">
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
        <constellation class="constellation" :class="{loading: this.$store.state.loading.playlistGenerating}" v-if="$store.getters['tracks/playlistLoaded'] && showConstellation"></constellation>
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

    .playButton button {
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
    }

    .playButton button svg {
        width: 1.5em;
        fill: white;
        stroke: white;
        stroke-linejoin: round;
    }

    .playButton button:hover svg {
        fill: magenta;
        stroke: magenta;
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

    td.number.pullRight {
        text-align: right;
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
</style>

<script>
    import {getSpotifyAlbumUrl, getSpotifyArtistUrl, getSpotifyTrackUrl, getSpotifyTrackDragTitle, getSpotifyAlbumDragTitle} from '~/assets/spotify';
    import constellation from '~/components/constellation';
    import loadingScreen from '~/components/loading-screen';

    export default {
        components: {constellation, loadingScreen},
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
            playlistPlayButton() {
                return !this.$store.getters['tracks/isPlayerDisconnected'];
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
                this.$store.dispatch('tracks/seekTrack', i);
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