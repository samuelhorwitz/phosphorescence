<template>
    <div class="mainContainer" ref="mainContainer">
        <div class="dropzone" :class="{dropzoneReady: dragStarted, dropHoverActive: isDropHovering, dropHoverBlockActive: isAlreadyGeneratingHover, invalidDragObject: isBadDropHovering}" @dragenter="handleDragenter" @dragleave="handleDragleave" @drop="handleDrop">
            {{dropText}}
        </div>
        <elastic v-if="isIOS"></elastic>
        <logo></logo>
        <toolbar></toolbar>
        <main :class="{playlistLoaded: $store.getters['tracks/playlistLoaded']}">
            <nuxt/>
        </main>
        <album></album>
        <player></player>
        <foot></foot>
    </div>
</template>

<style scoped>
    .mainContainer {
        display: grid;
        grid-template-columns: 10fr 1fr;
        height: 100vh;
    }

    @media only screen and (display-mode: fullscreen) and (orientation: portrait) {
        .mainContainer {
            height: calc(100vh - 70px);
        }
    }

    @media only screen and (display-mode: fullscreen) and (orientation: landscape) {
        .mainContainer {
            height: calc(100vh - 20px);
        }
    }

    body.playerConnected .mainContainer {
        grid-template-rows: min-content min-content minmax(100px, 1fr) min-content min-content;
    }

    body:not(.playerConnected) .mainContainer {
        grid-template-rows: min-content min-content minmax(100px, 1fr) min-content;
    }

    .dropzone {
        display: none;
        z-index: 1000000000;
        color: white;
        font-family: 'Caveat';
        font-size: 5em;
        outline: none;
        cursor: pointer;
        text-shadow: -1px -1px 0 midnightblue, 1px -1px 0 midnightblue, -1px 1px 0 midnightblue, 1px 1px 0 midnightblue;
        align-items: center;
        justify-content: center;
        box-sizing: border-box;
    }

    .dropzone * {
        pointer-events: none;
    }

    .dropzone.dropzoneReady {
        display: flex;
        position: absolute;
        width: 100%;
        height: 100%;
    }

    .dropzone.dropHoverActive {
        background-color: rgba(255, 0, 255, 0.75);
        border: 20px dashed aqua;
        cursor: grabbing;
    }

    .dropzone.dropHoverBlockActive {
        background-color: rgba(255, 0, 0, 0.75);
        border: 20px dashed red;
        cursor: no-drop;
    }

    .dropzone.invalidDragObject {
        cursor: no-drop;
        color: transparent;
    }

    body.playerConnected main.playlistLoaded {
        grid-column: 1 / 2;
    }

    body:not(.playerConnected) main,
    main:not(.playlistLoaded) {
        grid-column: 1 / 3;
    }

    @media only screen and (min-height: 450px) and (max-width: 1099px) {
        .mainContainer {
            grid-template-columns: 100%;
        }

        body.playerConnected .mainContainer {
            grid-template-rows: min-content min-content minmax(100px, 1fr) 100px min-content min-content;
        }

        body:not(.playerConnected) .mainContainer {
            grid-template-rows: min-content min-content minmax(100px, 1fr) min-content;
        }
    }

    @media only screen and (max-height: 449px) {
        .mainContainer {
            grid-template-columns: 28vw minmax(40vw, 10fr) min-content;
        }

        body.playerConnected .mainContainer {
            grid-template-rows: max-content minmax(100px, 1fr) max-content max-content;
        }

        body:not(.playerConnected) .mainContainer {
            grid-template-rows: max-content minmax(100px, 1fr) max-content;
        }

        main {
            grid-row: 2 / 3;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 1em 0;
        }

        body.playerConnected main,
        body.playerConnected main.playlistLoaded,
        body:not(.playerConnected) main {
            grid-column: 1 / 3;
        }
    }

    @media only screen and (max-height: 449px) and (max-width: 1199px) {
        .mainContainer {
            grid-template-columns: 28vw minmax(30vw, 10fr) min-content;
        }
    }

    @media only screen and (max-width: 949px) and (min-height: 275px) and (max-height: 449px) {
        body.playerConnected main.playlistLoaded {
            display: none;
        }

        body:not(.playerConnected) .mainContainer {
            grid-template-rows: max-content minmax(100px, 1fr);
        }
    }

    @media only screen and (max-height: 274px) {
        body.playerConnected main.playlistLoaded {
            display: none;
        }

        body:not(.playerConnected) .mainContainer {
            grid-template-rows: max-content minmax(100px, 1fr);
        }
    }

    @media only screen and (max-width: 500px) and (max-height: 249px) {
        .mainContainer {
            grid-template-columns: 0 minmax(40vw, 10fr) min-content;
        }
    }
</style>

<script>
    import {initialize, builders, loadNewPlaylist, processTrack} from '~/assets/recordcrate';
    import elastic from '~/components/elastic';
    import logo from '~/components/logo';
    import toolbar from '~/components/toolbar';
    import album from '~/components/album';
    import player from '~/components/player';
    import foot from '~/components/foot';

    export default {
        components: {
            elastic,
            logo,
            toolbar,
            album,
            player,
            foot
        },
        middleware: 'authenticated',
        head() {
            return {
                bodyAttrs: {
                    class: this.$store.getters['tracks/isPlayerConnected'] ? 'playerConnected' : ''
                }
            }
        },
        data() {
            return {
                destroyResizeListener: false,
                dragStarted: false,
                isDropHovering: false,
                isAlreadyGeneratingHover: false,
                isBadDropHovering: false
            };
        },
        computed: {
            isIOS() {
                return /\b(iPhone|iPod)\b/.test(navigator.userAgent);
            },
            dropText() {
                if (this.isAlreadyGeneratingHover) {
                    return 'Please wait';
                } else if (this.isBadDropHovering) {
                    return '';
                }
                return 'Drop your track!';
            },
            playlistGenerating() {
                return this.$store.state.loading.playlistGenerating;
            }
        },
        watch: {
            playlistGenerating() {
                if (this.isAlreadyGeneratingHover) {
                    this.isAlreadyGeneratingHover = false;
                    this.isDropHovering = true;
                }
            }
        },
        async created() {
            this.$store.commit('loading/startLoad');
            this.$store.commit('preferences/restore');
            this.$store.commit('tracks/restore');
            this.$store.commit('loading/playlistGenerating');
            let messageId = await this.$store.dispatch('loading/pushMessage', 'Downloading and processing track data');
            this.$store.dispatch('loading/initializeProgress', {id: 'tracks', weight: 60, ms: 300});
            await initialize(this.$store.state.user.user.country);
            this.$store.commit('loading/completeProgress', {id: 'tracks'});
            this.$store.commit('loading/clearMessage', messageId);
            if (!this.$store.getters['tracks/playlistLoaded']) {
                let loadingMessage = 'Generating playlist';
                if (this.$store.state.preferences.seedStyle) {
                    loadingMessage += ` (${this.$store.state.preferences.seedStyle})`;
                } else {
                    loadingMessage += ' (random)';
                }
                let messageId = await this.$store.dispatch('loading/pushMessage', loadingMessage);
                this.$store.dispatch('loading/initializeProgress', {id: 'generate', weight: 35, ms: 200, amount: 2});
                let pruners;
                if (this.$store.state.preferences.onlyTheHits) {
                    pruners = [builders.hits];
                }
                try {
                    let {playlist} = await loadNewPlaylist(this.$store.state.preferences.tracksPerPlaylist, builders.randomwalk, builders[this.$store.state.preferences.seedStyle], null, pruners);
                    this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                }
                catch (e) {
                    console.error('Playlist generation failed', e);
                    // TODO add some visual UI indication
                }
                this.$store.commit('loading/completeProgress', {id: 'generate'});
                this.$store.commit('loading/clearMessage', messageId);
            }
            this.$store.commit('loading/playlistGenerationComplete');
            this.$store.dispatch('loading/endLoadAfterDelay');
        },
        mounted() {
            document.body.addEventListener('dragenter', this.handleWindowDragenter);
            document.body.addEventListener('dragover', this.handleWindowDragover);
            document.body.addEventListener('drop', this.handleWindowDrop);
            if (/\b(iPhone|iPod)\b/.test(navigator.userAgent)) {
                document.body.addEventListener('resize', this.handleResize);
                window.addEventListener('orientationchange', this.handleResize);
                this.handleResize();
                this.destroyResizeListener = true;
            }
        },
        beforeDestroy() {
            document.body.removeEventListener('dragenter', this.handleWindowDragenter);
            document.body.removeEventListener('dragover', this.handleWindowDragover);
            document.body.removeEventListener('drop', this.handleWindowDrop);
            if (this.destroyResizeListener) {
                document.body.removeEventListener('resize', this.handleResize);
                window.removeEventListener('orientationchange', this.handleResize);
            }
        },
        methods: {
            handleResize(e) {
                if (matchMedia('(display-mode: fullscreen)').matches) {
                    return;
                }
                this.$refs.mainContainer.style.height = `${innerHeight}px`;
            },
            handleWindowDragenter(e) {
                this.dragStarted = true;
            },
            handleWindowDragover(e) {
                e.preventDefault();
            },
            handleWindowDrop(e) {
                e.preventDefault();
            },
            handleDragenter(e) {
                if (e.dataTransfer.types.includes('text/x-spotify-tracks')) {
                    if (this.$store.state.loading.playlistGenerating) {
                        this.isAlreadyGeneratingHover = true;
                    } else {
                        this.isDropHovering = true;
                    }
                } else {
                    this.isBadDropHovering = true;
                }
            },
            handleDragleave() {
                this.isDropHovering = false;
                this.isAlreadyGeneratingHover = false;
                this.isBadDropHovering = false;
                this.dragStarted = false;
            },
            async handleDrop(e) {
                let shouldHandleDrop = false;
                if (this.isDropHovering) {
                    shouldHandleDrop = true;
                }
                this.isDropHovering = false;
                this.isAlreadyGeneratingHover = false;
                this.isBadDropHovering = false;
                this.dragStarted = false;
                if (!shouldHandleDrop) {
                    return;
                }
                this.$store.commit('loading/startLoad');
                this.$store.commit('loading/playlistGenerating');
                e.preventDefault();
                let url = e.dataTransfer.getData('text/x-spotify-tracks');
                if (!url) {
                    return;
                }
                let trackParts = url.split('/');
                let trackId = trackParts[trackParts.length - 1];
                try {
                    let trackResponse = await fetch(`${process.env.API_ORIGIN}/track/${trackId}`, {credentials: 'include'});
                    let {track} = await trackResponse.json();
                    let processedTrack = await processTrack(this.$store.state.user.user.country, track);
                    let pruners;
                    if (this.$store.state.preferences.onlyTheHits) {
                        pruners = [builders.hits];
                    }
                    let {playlist} = await loadNewPlaylist(this.$store.state.preferences.tracksPerPlaylist, builders.randomwalk, null, processedTrack, pruners);
                    this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                }
                catch (e) {
                    console.error('Playlist generation failed', e);
                    // TODO add some visual UI indication
                }
                this.$store.commit('loading/playlistGenerationComplete');
                this.$store.dispatch('loading/endLoadAfterDelay');
            }
        }
    };
</script>
