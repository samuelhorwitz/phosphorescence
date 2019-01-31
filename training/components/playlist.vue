<template>
    <article>
        <input type="text" v-model="id" placeholder="Playlist ID">
        <button @click="getPlaylistData">Load Playlist</button>
        <button @click="getNewTranceSongs">Load 100 Trance Songs</button>
        <button @click="getNewCTTranceSongs">Load 10 C/T Songs</button>
        <button @click="getNewCPTranceSongs">Load 10 C/P Songs</button>
        <button @click="getNewATTranceSongs">Load 10 A/T Songs</button>
        <button @click="getNewAPTranceSongs">Load 10 A/P Songs</button>
        <button @click="getNewTransitionTranceSongs">Load 10 Transitional Songs</button>
        <div class="tableWrapper" ref="tableWrapper">
            <table>
                <thead>
                    <tr>
                        <th>
                            No.
                        </th>
                        <th>
                            Track
                        </th>
                        <th>
                            Artists
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(track, index) in $store.state.tracks.tracks" :class="{playing: isPlaying(track.track.id)}">
                        <td>{{index + 1}}</td>
                        <td>{{track.track.name}}</td>
                        <td>{{concatArtistNames(track.track.artists)}}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </article>
</template>

<style scoped>
    article {
        border: 1px solid black;
    }

    .tableWrapper {
        border-top: 2px solid black;
        overflow: auto;
        height: 10em;
    }

    table {
        width: 100%;
    }

    th {
        font-weight: bold;
    }

    thead tr, tbody tr:nth-child(even) {
        background-color: #ccc;
    }

    tr.playing {
        background-color: blue !important;
        color: white;
    }

    thead tr:nth-child(1) th{
        background: gray;
        position: sticky;
        top: 0;
        z-index: 10;
    }
</style>

<script>
    import {getPlaylist, getSomeTrance, shuffle, predict} from '~/util/spotify';

    export default{
        data() {
            return {
                id: '5lSgExb6yTKfLusGag7bm7',
                isLoadingPlaylist: false
            };
        },
        computed: {
            currentTrack: function() {
                return this.$store.getters['tracks/currentTrack'];
            }
        },
        watch: {
            currentTrack: function() {
                let playingEl = this.$el.querySelector('.tableWrapper .playing');
                if (playingEl) {
                    this.$refs.tableWrapper.scrollTop = playingEl.offsetTop;
                }
            }
        },
        methods: {
            async getPlaylistData() {
                if (this.isLoadingPlaylist) {
                    return;
                }
                this.isLoadingPlaylist = true;
                let {tracks, allTracksAnalysis} = await getPlaylist(this.id);
                this.$store.commit('tracks/load', {tracks: shuffle(tracks), analysis: allTracksAnalysis.reduce((acc, cur) => {
                    if (!cur) return acc;
                    acc[cur.id] = cur;
                    return acc;
                }, {})});
                this.isLoadingPlaylist = false;
            },
            async getNewTranceSongs() {
                if (this.isLoadingPlaylist) {
                    return;
                }
                this.isLoadingPlaylist = true;
                let {tracks, allTracksAnalysis} = await getSomeTrance();
                this.$store.commit('tracks/load', {tracks: shuffle(tracks), analysis: allTracksAnalysis.reduce((acc, cur) => {
                    if (!cur) return acc;
                    acc[cur.id] = cur;
                    return acc;
                }, {})});
                this.isLoadingPlaylist = false;
            },
            async getNewCTTranceSongs() {
                this.getFilteredSongs((a, p) => a < 0.4 && p < 0.4);
            },
            async getNewCPTranceSongs() {
                this.getFilteredSongs((a, p) => a < 0.4 && p >= 0.6);
            },
            async getNewATTranceSongs() {
                this.getFilteredSongs((a, p) => a >= 0.6 && p < 0.4);
            },
            async getNewAPTranceSongs() {
                this.getFilteredSongs((a, p) => a >= 0.6 && p >= 0.6);
            },
            async getNewTransitionTranceSongs() {
                this.getFilteredSongs((a, p) => a >= 0.4 && a < 0.6 && p >= 0.4 && p < 0.6);
            },
            async getFilteredSongs(filter) {
                if (this.isLoadingPlaylist) {
                    return;
                }
                this.isLoadingPlaylist = true;
                let picked = [];
                let allAnalysis = {};
                let triedMap = {};
                while (picked.length < 10) {
                    let {tracks, allTracksAnalysis} = await getSomeTrance();
                    let reducedAnalysis = allTracksAnalysis.reduce((acc, cur) => {
                        if (!cur) return acc;
                        acc[cur.id] = cur;
                        return acc;
                    }, {});
                    tracks = shuffle(tracks);
                    for (let track of tracks) {
                        if (triedMap[track.track.id]) {
                            continue;
                        }
                        triedMap[track.track.id] = true;
                        let {aetherealness, primordialness} = await predict(reducedAnalysis[track.track.id]);
                        if (filter(aetherealness, primordialness)) {
                            picked.push(track);
                        }
                        if (picked.length >= 10) {
                            break;
                        }
                    }
                    allAnalysis = Object.assign({}, allAnalysis, reducedAnalysis)
                }
                this.$store.commit('tracks/load', {tracks: picked, analysis: allAnalysis});
                this.isLoadingPlaylist = false;
            },
            concatArtistNames(artists) {
                return artists.map(artist => artist.name).join(', ');
            },
            isPlaying(id) {
                if (!this.currentTrack) {
                    return false;
                }
                return this.currentTrack.track.id == id;
            }
        }
    };
</script>