<template>
    <article v-show="$store.getters['tracks/playlistLoaded'] && $store.state.tracks.spotifyFullyRestored">
        <div class="tableWrapper" ref="tableWrapper">
            <table>
                <thead>
                    <tr>
                        <th class="number">
                        </th>
                        <th>
                            Title
                        </th>
                        <th>
                            Artist
                        </th>
                        <th v-show="showAlbums">
                            Album
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(track, index) in $store.state.tracks.playlist" :class="{currentTrack: isPlaying(track.track.id)}" @click="seekTrack(index)">
                        <td class="number">{{index + 1}}</td>
                        <td :title="track.track.name"><a target="_blank" :href="track.track.external_urls.spotify">{{track.track.name}}</a></td>
                        <td :title="track.track.artists.map(artist => artist.name).join(', ')">
                            <ol>
                                <li v-for="artist in track.track.artists">
                                    <a target="_blank" :href="artist.external_urls.spotify">{{artist.name}}</a>
                                </li>
                            </ol>
                        </td>
                        <td v-show="showAlbums" :title="track.track.album.name"><a target="_blank" :href="track.track.album.external_urls.spotify">{{track.track.album.name}}</a></td>
                    </tr>
                </tbody>
            </table>
        </div>
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
        align-items: flex-start;
        height: 100%;
    }

    ol {
        list-style-type: none;
        margin: 0px;
        padding: 0px;
        display: inline;
    }

    li {
        display: inline;
    }

    li:not(:last-child)::after {
        content: ', ';
        display: inline;
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
    }

    table {
        width: 100%;
        height: 100%;
        border-spacing: 0px;
        border-collapse: unset;
        table-layout: fixed;
    }

    th {
        font-weight: bold;
    }

    th.number {
        width: 3em;
    }

    tr {
        color: white;
        background-color: rgb(15,10,54);
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

    td.number {
        text-align: right;
    }

    tr.currentTrack, tr:hover {
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

    .tableWrapper {
        max-width: 600px;
    }

    @media only screen and (min-width: 1499px), (max-width: 1099px) and (min-height: 450px) {
        article {
            margin: 0px 2em;
        }

        .tableWrapper {
            max-width: unset;
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
        .tableWrapper {
            max-height: 600px;
        }

        .tracks {
            align-items: center;
        }
    }
</style>

<script>
    import {accessTokenExists, refreshUser} from '~/assets/session';

    export default {
        data() {
            return {
                showAlbumsMql: null,
                showAlbums: false
            };
        },
        watch: {
            currentTrack() {
                let playingEl = this.$el.querySelector('.tableWrapper .playing');
                if (playingEl) {
                    this.$refs.tableWrapper.scrollTop = playingEl.offsetTop;
                }
            }
        },
        computed: {
            currentTrack() {
                return this.$store.getters['tracks/currentTrack'];
            },
            currentTrackImage() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.album.images[0].url;
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
            isPlaying(id) {
                if (!this.currentTrack) {
                    return false;
                }
                return this.currentTrack.track.id == id;
            },
            seekTrack(i) {
                this.$store.dispatch('tracks/seekTrack', i);
            },
            handleShowAlbumsChange({matches}) {
                this.showAlbums = matches;
            }
        },
        created() {
            this.showAlbumsMql = matchMedia('only screen and (max-width: 1099px), (min-width: 1499px)');
            this.handleShowAlbumsChange(this.showAlbumsMql);
            this.showAlbumsMql.addListener(this.handleShowAlbumsChange);
        },
        beforeDestroy() {
            if (this.showAlbumsMql) {
                this.showAlbumsMql.removeListener(this.handleShowAlbumsChange);
            }
        },
        async fetch({store, error}) {
            if (!accessTokenExists()) {
                await refreshUser();
            }
            let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
            if (!userResponse.ok) {
                return error({statusCode: userResponse.status, message: "Could not get user information"});
            }
            let {user} = await userResponse.json();
            store.commit('user/user', user);
        }
    };
</script>