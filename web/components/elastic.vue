<template>
    <aside ref="elastic" :class="{real: useRealHeader, ready: isReadyToRelease, noTrack: noTrackCurrentlyPlaying}">
        <loadingBar></loadingBar>
        <div ref="container" class="container">
            <div class="svgContainer">
                <svg ref="pullArrow" id="pullArrow" v-show="isPulling || noTrackCurrentlyPlaying || isReadyToRelease" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" version="1.1" style="shape-rendering:geometricPrecision;text-rendering:geometricPrecision;image-rendering:optimizeQuality;" viewBox="0 0 847 847" x="0px" y="0px" fill-rule="evenodd" clip-rule="evenodd"><defs></defs><g><path class="fil0" d="M579 492l-59 0 0 -18c0,-8 -6,-14 -14,-14l-165 0c-8,0 -14,6 -14,14l0 17 -60 0 155 151 157 -150zm-71 -287l-170 0c-14,0 -14,25 0,25l170 0c15,0 15,-25 0,-25zm0 51l-170 0c-14,0 -14,25 0,25l170 0c15,0 15,-25 0,-25zm0 50l-170 0c-14,0 -14,26 0,26l170 0c15,0 15,-26 0,-26zm0 51l-170 0c-14,0 -14,26 0,26l170 0c15,0 15,-26 0,-26zm0 51l-170 0c-14,0 -14,26 0,26l170 0c15,0 15,-26 0,-26z"></path></g></svg>
                <svg id="loadingCancel" v-show="isReleased" @click="cancelGeneration()" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0px" y="0px" viewBox="0 0 500 500" enable-background="new 0 0 500 500" xml:space="preserve"><g><g><g><path d="M424.3,250c0,5.1-0.2,10.3-0.6,15.4c-0.2,2.3-0.4,4.6-0.6,6.9c-0.2,1.8-1.2,6.8,0.2-1.2     c-0.3,1.8-0.5,3.6-0.8,5.4c-1.6,9.8-4,19.4-7.1,28.9c-1.5,4.4-3.2,8.8-4.8,13.2c-1.1,3.1,2.7-6.1,0.7-1.7     c-0.5,1.2-1.1,2.4-1.6,3.6c-1.1,2.4-2.3,4.8-3.5,7.2c-4.4,8.6-9.4,16.9-15,24.8c-1.2,1.7-7.1,8.9-1.6,2.3     c-1.4,1.7-2.7,3.4-4.1,5.1c-3.4,4-6.9,7.9-10.5,11.7c-3,3.1-6.2,6.2-9.4,9.1c-1.9,1.8-3.9,3.5-5.9,5.1c-1.7,1.5-9.1,6.5-2.1,1.8     c-8.1,5.4-15.7,11-24.2,15.7c-4.2,2.3-8.6,4.5-13,6.5c-1.9,0.9-9.3,3.4-1.3,0.7c-2.5,0.9-5,2-7.5,2.9c-9.7,3.5-19.7,6.3-29.8,8.3     c-2.2,0.4-4.5,0.8-6.7,1.2c-1.3,0.2-2.7,0.4-4.1,0.6c0.2,0,6.7-0.8,2.6-0.4c-5.5,0.6-11.1,1.1-16.6,1.3     c-10.3,0.4-20.5-0.2-30.7-1.3c-5.2-0.5,6.3,1,1.2,0.2c-0.9-0.1-1.8-0.3-2.7-0.4c-2.7-0.4-5.4-0.9-8.1-1.5     c-4.9-1-9.7-2.2-14.5-3.5c-4.7-1.3-9.4-2.8-14-4.5c-2.1-0.8-4.1-1.6-6.2-2.4c-4.7-1.8,5.7,2.5,1.1,0.5c-1.6-0.7-3.2-1.4-4.8-2.2     c-8.8-4.1-17.2-8.8-25.3-14.2c-1.8-1.2-3.6-2.4-5.4-3.7c-1.1-0.8-2.1-1.5-3.2-2.3c-0.6-0.5-4.2-3.2-0.6-0.4c3.5,2.8,0,0-0.6-0.5     c-1-0.8-2-1.6-3-2.5c-2-1.7-4-3.4-5.9-5.1c-7.1-6.4-13.8-13.4-20-20.7c-1.4-1.7-6.5-9.2-1.8-2.1c-1.2-1.8-2.6-3.5-3.9-5.3     c-2.9-4.1-5.7-8.3-8.3-12.7c-2.4-4-4.6-8-6.7-12.1c-1.1-2.2-2.2-4.3-3.2-6.6c-0.5-1-0.9-2-1.4-3c-2.7-6.2,1.5,4-0.1-0.1     c-3.6-9.2-6.7-18.5-9-28.1c-1.2-4.9-2.1-9.8-2.9-14.7c-0.1-0.9-0.3-1.8-0.4-2.7c-0.8-5.1,0.7,6.4,0.2,1.2     c-0.3-2.8-0.6-5.5-0.8-8.3c-0.8-10.5-0.8-21,0.1-31.5c0.2-2.5,0.5-5,0.7-7.6c0.5-5-0.8,5.3-0.3,1.9c0.2-1.1,0.3-2.3,0.5-3.4     c0.9-5.4,1.9-10.7,3.2-16c2.4-9.6,5.9-18.7,9.1-28c-2.7,7.9-0.2,0.6,0.7-1.3c1.1-2.4,2.3-4.8,3.5-7.2c2.2-4.3,4.5-8.5,7-12.7     c2.5-4.1,5.1-8.2,7.9-12.1c1.3-1.8,2.7-3.5,3.9-5.3c-4.7,7.1,0.8-0.9,2.2-2.6c6.3-7.5,13.2-14.6,20.5-21.1     c3.3-2.9,6.7-5.6,10.1-8.4c-6.6,5.5,0.6-0.4,2.3-1.6c2.1-1.5,4.3-3,6.5-4.4c8.1-5.3,16.6-10,25.4-14c1.2-0.6,2.4-1.1,3.6-1.6     c4.6-2.1-5.8,2.3-1.1,0.5c2.1-0.8,4.1-1.6,6.2-2.4c4.8-1.8,9.7-3.3,14.7-4.7c4.6-1.3,9.2-2.4,13.8-3.3c2.5-0.5,4.9-0.9,7.4-1.4     c1.1-0.2,2.3-0.3,3.4-0.5c4.5-0.7-6.8,0.8,0.2,0c10.2-1.1,20.5-1.5,30.7-1.1c5.1,0.2,10.2,0.7,15.2,1.2c1.7,0.2,4.6,0.9-2.6-0.4     c1.3,0.2,2.7,0.4,4.1,0.6c2.7,0.4,5.4,0.9,8.1,1.5c9.7,2,19.2,4.7,28.5,8c2.5,0.9,4.9,2,7.5,2.9c-8-2.7-0.5-0.2,1.3,0.7     c4.8,2.2,9.5,4.6,14.1,7.2c3.8,2.1,7.6,4.4,11.3,6.8c2.2,1.4,4.4,2.9,6.5,4.4c1.8,1.3,8.8,7,2.3,1.6c7.4,6.2,14.7,12.3,21.4,19.3     c3.3,3.5,6.5,7,9.6,10.7c1.4,1.7,2.7,3.4,4.1,5.1c-5.4-6.5,0.3,0.5,1.6,2.3c5.8,8.2,11,16.9,15.6,26c1,2,1.9,4,2.9,6     c0.6,1.2,1.1,2.4,1.6,3.6c1.6,3.7-0.9-2.1-1-2.3c1.6,5.1,3.9,10,5.5,15.1c3,9.5,5.3,19.1,6.9,28.9c0.2,1.4,0.4,2.7,0.6,4.1     c-1.3-7.4-0.5-3.7-0.3-1.9c0.3,2.5,0.5,5,0.7,7.6C424.1,239.7,424.3,244.9,424.3,250c0,13.1,11.5,25.6,25,25     c13.5-0.6,25-11,25-25c-0.2-45.6-14-91.6-40.7-128.7C406.3,83.1,369,54.5,324.8,38.4c-86.3-31.5-188.7-3.6-247,67.4     c-30.3,36.9-48.3,81.1-52,128.8c-3.5,45.3,7.9,92.1,31.4,131c22.9,37.9,57.5,69.7,98,88c44.5,20,93.2,25.8,141.1,16.1     c88.4-17.9,161.1-93.2,174.8-182.5c1.9-12.3,3.2-24.6,3.2-37.1c0-13.1-11.5-25.6-25-25C435.7,225.6,424.4,236,424.3,250z"></path></g></g><g><g><path d="M148.3,184.5c18.8,18.8,37.6,37.6,56.4,56.4c29.8,29.8,59.6,59.6,89.4,89.4c6.8,6.8,13.7,13.7,20.5,20.5     c9.3,9.3,26.2,9.9,35.4,0c9.2-10,9.9-25.5,0-35.4c-18.8-18.8-37.6-37.6-56.4-56.4c-29.8-29.8-59.6-59.6-89.4-89.4     c-6.8-6.8-13.7-13.7-20.5-20.5c-9.3-9.3-26.2-9.9-35.4,0C139.2,159.1,138.4,174.6,148.3,184.5L148.3,184.5z"></path></g></g><g><g><path d="M183.7,350.9c18.8-18.8,37.6-37.6,56.4-56.4c29.8-29.8,59.6-59.6,89.4-89.4c6.8-6.8,13.7-13.7,20.5-20.5     c9.3-9.3,9.9-26.2,0-35.4c-10-9.2-25.5-9.9-35.4,0c-18.8,18.8-37.6,37.6-56.4,56.4c-29.8,29.8-59.6,59.6-89.4,89.4     c-6.8,6.8-13.7,13.7-20.5,20.5c-9.3,9.3-9.9,26.2,0,35.4C158.3,360,173.8,360.8,183.7,350.9L183.7,350.9z"></path></g></g></g></svg>
            </div>
            <div class="messageContainer">
                <p v-show="failed">somethin went wrong :'(</p>
                <p v-show="noTrackCurrentlyPlaying">no track currently playing</p>
                <p v-show="isPulling">create playlist from current track</p>
                <p v-show="isReadyToRelease">release to create playlist</p>
                <p v-show="isReleased">creating playlist...</p>
                <p v-if="track && showTrackData" class="track">{{track.name}} - {{trackArtists}}</p>
            </div>
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
        z-index: 99999999;
    }

    aside.real {
        color: white !important;
        height: 60px !important;
    }

    .container {
        display: flex;
        align-items: center;
        justify-content: center;
        flex: 1;
        box-sizing: border-box;
        width: 100%;
    }

    .messageContainer {
        display: flex;
        flex: 10;
        flex-direction: column;
        align-items: flex-start;
        justify-content: center;
        overflow-x: hidden;
    }

    aside.noTrack .messageContainer {
        align-items: center;
    }

    .svgContainer {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    aside.real .svgContainer {
        min-width: 10em;
        min-height: 10em;
    }

    svg {
        fill: white;
    }

    svg#loadingCancel {
        width: 5em;
    }

    svg#pullArrow {
        width: 10em;
        transition: transform 0.2s linear 0s, opacity 0.2s linear 0s;
    }

    aside.ready svg#pullArrow {
        transform: rotate(-180deg);
    }

    aside.noTrack svg#pullArrow {
        opacity: 0.5;
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
        max-width: 100%;
        color: cyan;
    }
</style>

<script>
    import {builders, loadNewPlaylist, processTrack} from '~/assets/recordcrate';
    import {terminatePlaylistBuilding} from '~/assets/eos';
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
        props: {
            container: HTMLElement
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
                } else if (this.useRealHeader) {
                    requestAnimationFrame(this.repaintReal);
                }
            },
            repaint() {
                let absScrollY = Math.abs(scrollY);
                this.$refs.elastic.style.height = `${absScrollY}px`;
                this.$refs.container.style.opacity = 1;
                this.$refs.elastic.style.color = `rgba(255, 255, 255, ${Math.min(absScrollY / barHeight, 1)})`;
                this.$refs.pullArrow.style.fill = `rgba(255, 255, 255, ${Math.min(absScrollY / barHeight, 1)})`;
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
            repaintReal() {
                this.$refs.container.style.opacity = Math.min(Math.max((barHeight - scrollY) / barHeight, 0), 1);
            },
            cancelGeneration() {
                terminatePlaylistBuilding();
            },
            handleTouchEnd() {
                if (this.state === readyToRelease) {
                    let absScrollY = Math.abs(scrollY);
                    this.state = released;
                    this.container.style.transform = `translateY(${absScrollY}px)`;
                    scrollTo(0, absScrollY);
                    this.container.addEventListener('webkitTransitionEnd', () => {
                        this.container.style.transition = '';
                    }, {once: true});
                    this.container.style.transition = 'transform 0.2s ease-out 0s';
                    this.container.style.transform = `translateY(${barHeight}px)`;
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
                        console.debug(now, this.trackState.fetchedAt, (now - this.trackState.fetchedAt) / 1000, this.trackState.progress, this.trackState.progress / 1000);
                        let offset = Math.max(0, (now - this.trackState.fetchedAt) + this.trackState.progress);
                        if (offset > track.duration_ms) {
                            this.$store.dispatch('tracks/next');
                            offset = 0;
                        }
                        this.$store.dispatch('tracks/play', offset);
                    }
                    this.state = complete;
                    this.$store.commit('loading/completeProgress', {id: 'generate'});
                    this.$store.commit('loading/resetProgress');
                } catch (e) {
                    this.state = failed;
                    this.$store.dispatch('loading/failProgress');
                    await new Promise(res => setTimeout(res, 1000));
                }
                this.state = complete;
                this.$store.commit('loading/playlistGenerationComplete');
                this.$store.dispatch('loading/endLoadAfterDelay');
                this.track = null;
                this.trackState = null;
                this.container.addEventListener('webkitTransitionEnd', () => {
                    this.container.style.transition = '';
                    this.container.style.transform = '';
                    this.state = notTouched;
                }, {once: true});
                this.container.style.transition = 'transform 0.5s ease 0s';
                this.container.style.transform = 'translateY(0)';
            }
        }
    };
</script>