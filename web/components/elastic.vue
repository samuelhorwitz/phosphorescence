<template>
    <aside ref="elastic" :class="{real: useRealHeader}">
        <div class="container">
            <loadingBar></loadingBar>
            <p v-show="failed">somethin went wrong :'(</p>
            <p v-show="noTrackCurrentlyPlaying">no track currently playing</p>
            <p v-show="isPulling">create playlist from current track</p>
            <p v-show="isReadyToRelease">release to create playlist</p>
            <p v-show="isReleased">creating playlist...</p>
            <p v-if="track && showTrackData" class="track">{{track.name}} - {{trackArtists}}</p>
        </div>
    </aside>
</template>

<style scoped>
    aside {
        position: fixed;
        height: 0px;
        width: 100%;
        z-index: -1;
        display: flex;
        align-items: flex-start;
        justify-content: flex-start;
        color: transparent;
    }

    aside.real {
        color: white !important;
    }

    .container {
        display: flex;
        align-items: center;
        justify-content: center;
        flex: 1;
        padding: 1em;
        flex-direction: column;
        max-width: 100vw;
        box-sizing: border-box;
    }

    p {
        font-size: 16px;
        margin: 0px;
        padding: 0px;
        font-style: italic;
        flex: 1;
        text-align: center;
    }

    p.track {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-style: normal;
        font-size: 14px;
        max-width: 100vw;
        color: cyan;
    }
</style>

<script>
    import {builders, loadNewPlaylist, processTrack} from '~/assets/recordcrate';
    import loadingBar from '~/components/loading-bar';

    const failed = -3;
    const playlistGenerating = -2;
    const noTrack = -1;
    const notTouched = 0;
    const pulling = 1;
    const readyToRelease = 2;
    const released = 3;
    const complete = 4;
    const barHeight = 60;

    export default {
        components: {
            loadingBar
        },
        data() {
            return {
                state: notTouched,
                recheckCurrentlyPlaying: true,
                track: null,
                trackState: null
            }
        },
        computed: {
            trackArtists() {
                if (!this.track || !this.track.artists) {
                    return '';
                }
                return this.track.artists.map(artist => artist.name).join(', ');
            },
            useRealHeader() {
                return this.isReleased || this.failed;
            },
            noTrackCurrentlyPlaying() {
                return this.state === noTrack;
            },
            isPulling() {
                return this.state === pulling;
            },
            isReadyToRelease() {
                return this.state === readyToRelease;
            },
            isReleased() {
                return this.state === released;
            },
            showTrackData() {
                return this.isReadyToRelease || this.isReleased;
            },
            playlistGenerating() {
                return this.state === playlistGenerating;
            },
            failed() {
                return this.state === failed;
            },
            disabled() {
                return this.playlistGenerating;
            }
        },
        mounted() {
            addEventListener('scroll', this.handleScroll, {passive: true});
            addEventListener('touchend', this.handleTouchEnd, {passive: true});
        },
        beforeDestroy() {
            removeEventListener('scroll', this.handleScroll, {passive: true});
            removeEventListener('touchend', this.handleTouchEnd, {passive: true});
        },
        methods: {
            handleScroll() {
                if (matchMedia('(orientation: landscape)').matches) {
                    return;
                }
                if (scrollY <= 0) {
                    requestAnimationFrame(this.repaint);
                    if (scrollY <= -10 && this.recheckCurrentlyPlaying && this.isPulling && !this.$store.state.loading.playlistGenerating) {
                        this.recheckCurrentlyPlaying = false;
                        this.loadCurrentlyPlaying();
                    }
                    if (scrollY === 0) {
                        this.recheckCurrentlyPlaying = true;
                    }
                }
            },
            repaint() {
                let absScrollY = Math.abs(scrollY);
                this.$refs.elastic.style.height = `${absScrollY}px`;
                this.$refs.elastic.style.color = `rgba(255, 255, 255, ${Math.min(absScrollY / barHeight, 1)})`;
                if (this.$store.state.loading.playlistGenerating && this.state !== released && this.state !== failed) {
                    this.state = playlistGenerating;
                } else if (absScrollY > barHeight && ((this.state === pulling && this.track) || this.state === readyToRelease)) {
                    this.state = readyToRelease;
                } else if (scrollY < 0 && this.state === readyToRelease) {
                    this.state = pulling;
                } else if (scrollY === 0 && this.state <= pulling && this.state !== failed) {
                    this.state = notTouched;
                } else if (this.state === notTouched) {
                    this.state = pulling;
                } else if (this.state === playlistGenerating && !this.$store.state.loading.playlistGenerating) {
                    this.state = pulling;
                }
            },
            handleTouchEnd() {
                if (this.state === readyToRelease) {
                    let absScrollY = Math.abs(scrollY);
                    this.state = released;
                    this.$refs.elastic.style.transform = `translateY(-${absScrollY}px)`;
                    this.$refs.elastic.parentNode.style.transform = `translateY(${absScrollY}px)`;
                    scrollTo(0, absScrollY);
                    this.$refs.elastic.parentNode.addEventListener('webkitTransitionEnd', () => {
                        this.$refs.elastic.parentNode.style.transition = '';
                    }, {once: true});
                    this.$refs.elastic.addEventListener('webkitTransitionEnd', () => {
                        this.$refs.elastic.style.transition = '';
                    }, {once: true});
                    this.$refs.elastic.parentNode.style.transition = 'transform 0.2s ease-out 0s';
                    this.$refs.elastic.parentNode.style.transform = `translateY(${barHeight}px)`;
                    this.$refs.elastic.style.transition = 'transform 0.2s ease-out 0s';
                    this.$refs.elastic.style.transform = `translateY(-${barHeight}px)`;
                    this.generateFromTrack();
                }
            },
            async loadCurrentlyPlaying() {
                let currentlyPlayingResponse = await fetch(`${process.env.API_ORIGIN}/user/me/currently-playing`, {credentials: 'include'});
                if (currentlyPlayingResponse.ok) {
                    let {track, isPlaying, progress, fetchedAt} = await currentlyPlayingResponse.json();
                    this.track = track;
                    this.trackState = {isPlaying, progress, fetchedAt};
                } else {
                    this.state = noTrack;
                }
            },
            async generateFromTrack() {
                this.$store.commit('loading/startLoad');
                this.$store.commit('loading/playlistGenerating');
                this.$store.commit('loading/initializeProgress', {id: 'generate'});
                try {
                    let trackResponse = await fetch(`${process.env.API_ORIGIN}/track/${this.track.id}`, {credentials: 'include'});
                    let {track} = await trackResponse.json();
                    let processedTrack = await processTrack(this.$store.state.user.user.country, track);
                    let pruners;
                    if (this.$store.state.preferences.onlyTheHits) {
                        pruners = [builders.hits];
                    }
                    let {playlist} = await loadNewPlaylist(this.$store.state.preferences.tracksPerPlaylist, builders.randomwalk, null, processedTrack, pruners, percent => {
                        this.$store.commit('loading/tickProgress', {id: 'generate', percent});
                    });
                    this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                    if (this.trackState.isPlaying) {
                        let now = new Date().getTime();
                        console.log(now, this.trackState.fetchedAt, (now - this.trackState.fetchedAt) / 1000, this.trackState.progress, this.trackState.progress / 1000);
                        let offset = Math.max(0, (now - this.trackState.fetchedAt) + this.trackState.progress);
                        if (offset > track.duration_ms) {
                            this.$store.dispatch('tracks/next');
                            offset = 0;
                        }
                        this.$store.dispatch('tracks/play', offset);
                    }
                    this.state = complete;
                } catch (e) {
                    this.state = failed;
                    await new Promise(res => setTimeout(res, 1000));
                }
                this.state = complete;
                this.$store.commit('loading/completeProgress', {id: 'generate'});
                this.$store.commit('loading/resetProgress');
                this.$store.commit('loading/playlistGenerationComplete');
                this.$store.dispatch('loading/endLoadAfterDelay');
                this.track = null;
                this.trackState = null;
                this.$refs.elastic.parentNode.addEventListener('webkitTransitionEnd', () => {
                    this.$refs.elastic.parentNode.style.transition = '';
                    this.$refs.elastic.parentNode.style.transform = '';
                    this.state = notTouched;
                }, {once: true});
                this.$refs.elastic.parentNode.style.transition = 'transform 0.5s ease 0s';
                this.$refs.elastic.parentNode.style.transform = 'translateY(0)';
                this.$refs.elastic.style.transform = '';
            }
        }
    };
</script>