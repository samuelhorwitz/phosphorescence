<template>
    <article v-show="$store.getters['tracks/playlistLoaded']">
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
                        <th class="album">
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
                        <td class="album" :title="track.track.album.name"><a target="_blank" :href="track.track.album.external_urls.spotify">{{track.track.album.name}}</a></td>
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
        margin: 0px 2em;
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
        font-size: 16px;
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

    @media only screen and (min-width: 1199px), (max-width: 1099px) and (min-height: 450px) {
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
    import {accessTokenExists, refreshUser} from '~/assets/session';

    export default {
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
                if (this.$store.getters['tracks/isPlayerDisconnected']) {
                    return null;
                }
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