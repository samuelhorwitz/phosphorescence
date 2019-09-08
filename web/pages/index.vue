<template>
    <article :class="{loading: !$store.getters['tracks/playlistLoaded']}">
        <div class="tableWrapper" ref="tableWrapper" v-show="$store.getters['tracks/playlistLoaded']">
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
                        <td :title="humanReadableEvocativeness[index]" class="number">{{index + 1}}</td>
                        <td :title="track.track.name"><a target="_blank" :href="track.track.external_urls.spotify">{{track.track.name}}</a></td>
                        <td :title="track.track.artists.map(artist => artist.name).join(', ')">
                            <ol>
                                <li class="artist" v-for="artist in track.track.artists">
                                    <a target="_blank" :href="artist.external_urls.spotify">{{artist.name}}</a>
                                </li>
                            </ol>
                        </td>
                        <td class="album" :title="track.track.album.name"><a target="_blank" :href="track.track.album.external_urls.spotify">{{track.track.album.name}}</a></td>
                    </tr>
                </tbody>
            </table>
        </div>
        <aside class="loading" v-show="!$store.getters['tracks/playlistLoaded']">
            <ul>
                <li class="loadingMessage" v-for="loadMessage in $store.state.loading.descriptions" :class="{done: loadMessage.done}">
                    {{loadMessage.description}}...
                </li>
            </ul>
        </aside>
    </article>
</template>

<style scoped>
    article.loading {
        align-items: flex-start;
    }

    aside.loading {
        overflow: auto;
        color: white;
        padding: 1em;
        flex: 1;
        z-index: 10000000;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    aside.loading ul {
        margin: 0px;
        padding: 0px;
        flex: 1;
    }

    li.loadingMessage {
        display: block;
        font-size: 3em;
    }

    li.loadingMessage:not(.done)::before {
        content: '‣';
        color: magenta;
        display: inline-block;
        width: 1em;
    }

    li.loadingMessage.done::before {
        content: '✓';
        color: cyan;
        display: inline-block;
        width: 1em;
    }

    progress {
        width: 100%;
    }

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

    li.artist {
        display: inline;
    }

    li.artist:not(:last-child)::after {
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

        aside.loading {
            flex: 1;
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
                return this.currentTrack.track.id == id;
            },
            seekTrack(i) {
                this.$store.dispatch('tracks/seekTrack', i);
            }
        }
    };
</script>