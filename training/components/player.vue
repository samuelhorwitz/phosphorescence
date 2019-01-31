<template>
    <article>
        <button @click="play" :disabled="!playlistLoaded">Play!</button>
        <button @click="next" :disabled="!playlistLoaded">Next!</button>
        <div>{{currentTrackName}} - {{currentTrackArtists}}</div>
        <!-- <pre class="analysis">{{currentTrackAnalysis}}</pre> -->
    </article>
</template>

<style scoped>
    .analysis {
        color: white;
    }

    .analysis:hover {
        color: black;
    }
</style>

<script>
    export default {
        props: {
            deviceId: String
        },
        computed: {
            playlistLoaded() {
                return !!this.$store.state.tracks.tracks;
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
            currentTrackArtists() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.artists.map(artist => artist.name).join(', ');
            },
            currentTrackAnalysis() {
                return JSON.stringify(this.$store.getters['tracks/currentTrackAnalysis'], null, 4);
            }
        },
        created() {
            this.$store.subscribe(({type}) => {
                if (type == 'tracks/nextTrack') {
                    this.play();
                }
            });
        },
        methods: {
            play() {
                fetch(`https://api.spotify.com/v1/me/player/play?device_id=${this.deviceId}`, {
                    method: 'PUT',
                    headers: {
                        Authorization: `Bearer ${this.$store.state.session.accessToken}`
                    },
                    body: JSON.stringify({
                        uris: [this.currentTrack.track.uri]
                    })
                });
            },
            next() {
                this.$store.commit('tracks/nextTrack');
            }
        }
    };
</script>