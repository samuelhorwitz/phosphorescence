<template>
    <div>
        <elastic v-if="isLoggedInUser && isIOS" :container="$refs.mainContainer"></elastic>
        <div class="mainContainer" ref="mainContainer">
            <loadingBar v-if="!isIOS"></loadingBar>
            <dropzone></dropzone>
            <logo></logo>
            <toolbar></toolbar>
            <main :class="{playlistLoaded: $store.getters['tracks/playlistLoaded']}">
                <nuxt/>
            </main>
            <album></album>
            <player></player>
            <foot></foot>
        </div>
    </div>
</template>

<style scoped>
    .mainContainer {
        display: grid;
        grid-template-columns: 10fr 1fr;
        height: 100vh;
    }

    main {
        display: flex;
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
    import {initialize, builders, loadNewPlaylist} from '~/assets/recordcrate';
    import elastic from '~/components/elastic';
    import logo from '~/components/logo';
    import toolbar from '~/components/toolbar';
    import album from '~/components/album';
    import player from '~/components/player';
    import foot from '~/components/foot';
    import loadingBar from '~/components/loading-bar';
    import dropzone from '~/components/dropzone';

    const deviceNotSupportedError = 'Your device may not be supported, or your ad blocker may be blocking our playlist building engine. Please visit this page with an up-to-date version of the Edge, Chrome, Firefox, Opera or Safari browsers.';

    export default {
        components: {
            elastic,
            logo,
            toolbar,
            album,
            player,
            foot,
            loadingBar,
            dropzone
        },
        middleware: ['authenticated'],
        head() {
            let additionalClasses = ['appHome', 'vividSunrise'];
            if (this.$store.getters['tracks/isPlayerConnected']) {
                additionalClasses.push('playerConnected');
            } else {
                document.body.classList.remove('playerConnected');
            }
            let existingClasses = document.body.getAttribute('class') || '';
            if (additionalClasses.length === 0) {
                return {bodyAttrs: {class: existingClasses}};
            }
            let missingClasses = [];
            for (let additionalClass of additionalClasses) {
                if (!document.body.classList.contains(additionalClass)) {
                    missingClasses.push(additionalClass);
                }
            }
            if (missingClasses.length === 0) {
                return {bodyAttrs: {class: existingClasses}};
            }
            return {
                bodyAttrs: {
                    class: `${existingClasses} ${missingClasses.join(' ')}`
                }
            }
        },
        data() {
            return {
                destroyResizeListener: false
            };
        },
        computed: {
            isIOS() {
                return /\b(iPhone|iPod)\b/.test(navigator.userAgent);
            },
            playlistGenerating() {
                return this.$store.state.loading.playlistGenerating;
            },
            isLoggedInUser() {
                return !!this.$store.state.user.user;
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
            this.$store.dispatch('tracks/restore');
            this.$store.commit('loading/tracksDownloading');
            let messageId = await this.$store.dispatch('loading/pushMessage', 'Downloading track data');
            this.$store.commit('loading/initializeProgress', {id: 'tracks', weight: 60});
            this.$store.commit('loading/initializeProgress', {id: 'generate', weight: 35});
            let downloadingPhase = true;
            try {
                await initialize(this.$store.getters['user/country'], this.isLoggedInUser, async (percent, processingBegan) => {
                    if (downloadingPhase && processingBegan) {
                        this.$store.commit('loading/clearMessage', messageId);
                        messageId = await this.$store.dispatch('loading/pushMessage', 'Processing track data');
                        downloadingPhase = false;
                    }
                    this.$store.commit('loading/tickProgress', {id: 'tracks', percent});
                });
            }
            catch (e) {
                this.$ga.exception(e, true);
                console.error('Track initialization failed', e);
                this.$store.dispatch('loading/failProgress');
                this.$nuxt.error({message: deviceNotSupportedError});
                return;
            }
            this.$store.commit('loading/trackDownloadingComplete');
            this.$store.commit('loading/completeProgress', {id: 'tracks'});
            this.$store.commit('loading/clearMessage', messageId);
            if (!this.$store.getters['tracks/playlistLoaded']) {
                this.$store.commit('loading/playlistGenerating');
                let loadingMessage = 'Generating playlist';
                if (this.$store.state.preferences.seedStyle) {
                    loadingMessage += ` (${this.$store.state.preferences.seedStyle})`;
                } else {
                    loadingMessage += ' (random)';
                }
                let messageId = await this.$store.dispatch('loading/pushMessage', loadingMessage);
                let pruners;
                if (this.$store.state.preferences.onlyTheHits) {
                    pruners = [builders.hits];
                }
                try {
                    let {playlist} = await loadNewPlaylist({
                        count: this.$store.state.preferences.tracksPerPlaylist,
                        builder: builders.randomwalk,
                        firstTrackBuilder: builders[this.$store.state.preferences.seedStyle],
                        pruners
                    }, percent => {
                        this.$store.commit('loading/tickProgress', {id: 'generate', percent});
                    });
                    this.$store.commit('loading/completeProgress', {id: 'generate'});
                    await new Promise(resolve => setTimeout(resolve, 200));
                    this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                }
                catch (e) {
                    this.$ga.exception(e, true);
                    console.error('Playlist generation failed', e);
                    this.$store.dispatch('loading/failProgress');
                    this.$nuxt.error({message: deviceNotSupportedError});
                    return;
                }
                this.$store.commit('loading/clearMessage', messageId);
                this.$store.commit('loading/playlistGenerationComplete');
            } else {
                this.$store.commit('loading/completeProgress', {id: 'generate'});
            }
            this.$store.commit('loading/resetProgress');
            this.$store.dispatch('loading/endLoadAfterDelay');
            this.$ga.time('App Layout', 'load', Math.round(performance.now()));
        },
        mounted() {
            if (/\b(iPhone|iPod)\b/.test(navigator.userAgent)) {
                document.body.addEventListener('resize', this.handleResize);
                window.addEventListener('orientationchange', this.handleResize);
                this.handleResize();
                this.destroyResizeListener = true;
            }
        },
        beforeDestroy() {
            this.$store.dispatch('loading/resetAll');
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
                setTimeout(() => {
                    this.$refs.mainContainer.style.height = `${innerHeight}px`;
                }, 200);
            }
        }
    };
</script>
