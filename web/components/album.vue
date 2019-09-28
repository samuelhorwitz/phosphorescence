<template>
    <aside>
        <a v-if="currentTrack" target="_blank" rel="external noopener" :href="currentAlbumUrl" class="art-wrapper-link">
            <div class="art-wrapper" :title="currentTrackImageAltText" :style="{'background-image': 'url(' + currentTrackImage + ')'}"></div>
        </a>
    </aside>
</template>

<style scoped>
    body:not(.playerConnected) aside {
        display: none;
    }

    aside {
        max-height: 100%;
        display: flex;
        flex: 1;
        align-items: center;
        justify-content: center;
        flex-direction: column;
        height: 100%;
        grid-column: 2 / 3;
        margin-right: 2em;
    }

    .art-wrapper {
        border: 7px outset magenta;
        background-color: aquamarine;
        display: flex;
        align-items: center;
        justify-content: center;
        background-repeat: no-repeat;
        background-size: contain;
        background-position: center center;
        height: 100%;
        width: 60vh;
        max-height: 60vh;
        box-sizing: border-box;
    }

    .art-wrapper-link {
        height: 100%;
        display: flex;
        align-items: flex-start;
    }

    @media only screen and (max-width: 1099px) and (min-height: 450px) {
        aside {
            flex-basis: 25%;
            flex-grow: 0;
            grid-column: 1 / 2;
            margin: 0px 2em;
        }

        .art-wrapper-link {
            width: 100%;
        }

        .art-wrapper {
            width: 100%;
        }
    }

    @media only screen and (min-width: 1499px) and (min-height: 450px) {
        aside {
            flex: 0;
            margin-left: 0.5em;
        }
    }

    @media only screen and (max-height: 449px) {
        aside {
            margin-right: 1em;
        }
    }

    @media only screen and (max-height: 449px) {
        .art-wrapper-link {
            margin-top: 0.5em;
        }

        .art-wrapper {
            width: 90vh;
            max-height: 90vh;
        }

        aside {
            grid-column: 3 / 4;
            grid-row: 1 / 5;
        }
    }

    @media only screen and (max-height: 374px) {
        .art-wrapper {
            width: 85vh;
            max-height: 85vh;
        }
    }

    @media only screen and (max-height: 249px) {
        .art-wrapper {
            width: 80vh;
            max-height: 80vh;
        }
    }

    @media only screen and (min-height: 1199px) {
        .art-wrapper-link {
            align-items: center;
        }
    }

    @media only screen and (max-width: 949px) and (min-height: 250px) and (max-height: 449px) {
        aside {
            grid-column: 1 / 4;
            grid-row: 2 / 3;
            margin: 0;
        }
    }

    @media only screen and (max-width: 449px) and (max-height: 249px) {
        aside {
            margin: 0 1em;
        }
    }
</style>

<script>
    import {getSpotifyAlbumUrl} from '~/assets/spotify';

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
            currentAlbumUrl() {
                return getSpotifyAlbumUrl(this.currentTrack.track.album.id);
            }
        }
    };
</script>