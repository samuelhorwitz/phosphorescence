<template>
    <div class="dropzone" ref="dropzone" :class="{dropzoneReady: dragStarted, dropHoverActive: isDropHovering, dropHoverBlockActive: isAlreadyGeneratingHover, invalidDragObject: isBadDropHovering, dropHoverImmediate: isDropHoveringOverTarget}" @dragenter.capture="handleDragenter" @dragleave="handleDragleave" @drop="handleDragCancel">
        <div class="dropzoneTarget" ref="dropzoneTarget" :class="{dropzoneReady: dragStarted, miniDropzone: isDropFromPhosphor}" @dragenter.capture="handleDragenterTarget" @dragleave="handleDragleaveTarget" @drop="handleDrop">
            {{dropText}}
        </div>
    </div>
</template>

<style scoped>
    .dropzone {
        display: none;
        z-index: 1000000000;
        cursor: no-drop;
    }

    .dropzone.dropzoneReady {
        display: flex;
        position: absolute;
        width: 100%;
        height: 100%;
        top: 0px;
        left: 0px;
    }

    .dropzoneTarget.dropzoneReady {
        display: flex;
        position: absolute;
        width: 100%;
        height: 100%;
        color: white;
        font-family: 'Caveat';
        font-size: 5em;
        outline: none;
        text-shadow: -1px -1px 0 midnightblue, 1px -1px 0 midnightblue, -1px 1px 0 midnightblue, 1px 1px 0 midnightblue;
        align-items: center;
        justify-content: center;
        box-sizing: border-box;
    }

    .dropzoneTarget.dropzoneReady.miniDropzone {
        width: 10ex;
        height: 10ex;
        right: 0px;
        bottom: 0px;
    }

    .dropzone.dropHoverActive .dropzoneTarget {
        background-color: rgba(255, 0, 255, 0.75);
        border: 20px dashed aqua;
    }

    .dropzone.dropHoverImmediate.dropHoverActive {
        cursor: grabbing;
    }

    .dropzone.dropHoverBlockActive .dropzoneTarget {
        background-color: rgba(255, 0, 0, 0.75);
        border: 20px dashed red;
    }

    .dropzone.dropHoverImmediate.dropHoverBlockActive {
        cursor: no-drop;
    }

    .dropzone.invalidDragObject .dropzoneTarget {
        color: transparent;
    }

    .dropzone.dropHoverImmediate.invalidDragObject {
        cursor: no-drop;
    }
</style>

<script>
    import {getCaptchaToken} from '~/assets/captcha';
    import {builders, loadNewPlaylist, processTrack, processTracks} from '~/assets/recordcrate';

    export default {
        data() {
            return {
                dragStarted: false,
                isDropHovering: false,
                isAlreadyGeneratingHover: false,
                isBadDropHovering: false,
                isDropFromPhosphor: false,
                isDropHoveringOverTarget: false,
                count: 0 //safari sux and cant even handle related target
            };
        },
        computed: {
            dropText() {
                if (this.isAlreadyGeneratingHover) {
                    return 'Please wait';
                } else if (this.isBadDropHovering) {
                    return '';
                }
                return 'Drop!';
            }
        },
        mounted() {
            document.body.addEventListener('dragenter', this.handleWindowDragenter);
            document.body.addEventListener('dragover', this.handleWindowDragover);
            document.body.addEventListener('drop', this.handleWindowDrop);
        },
        beforeDestroy() {
            document.body.removeEventListener('dragenter', this.handleWindowDragenter);
            document.body.removeEventListener('dragover', this.handleWindowDragover);
            document.body.removeEventListener('drop', this.handleWindowDrop);
        },
        methods: {
            handleWindowDragenter(e) {
                this.dragStarted = true;
            },
            handleWindowDragover(e) {
                e.preventDefault();
            },
            handleWindowDrop(e) {
                e.preventDefault();
                this.count = 0;
            },
            handleDragenter(e) {
                this.count++;
                if (e.dataTransfer.types.includes('text/x-spotify-tracks') || e.dataTransfer.types.includes('text/x-spotify-playlists') || e.dataTransfer.types.includes('text/x-spotify-albums')) {
                    if (this.$store.state.loading.playlistGenerating) {
                        this.isAlreadyGeneratingHover = true;
                    } else {
                        if (e.dataTransfer.types.includes('text/x-phosphor-origin')) {
                            this.isDropFromPhosphor = true;
                        }
                        this.isDropHovering = true;
                    }
                } else {
                    this.isBadDropHovering = true;
                }
                if (e.target == this.$refs.dropzoneTarget) {
                    return;
                }
                e.dataTransfer.dropEffect = 'none';
            },
            handleDragenterTarget(e) {
                this.count++;
                e.dataTransfer.dropEffect = 'copy';
                this.isDropHoveringOverTarget = true;
            },
            handleDragleave(e) {
                this.count--;
                if (e.relatedTarget == this.$refs.dropzoneTarget || e.relatedTarget == this.$refs.dropzone || this.count != 0) {
                    return;
                }
                this.handleDragCancel();
            },
            handleDragleaveTarget(e) {
                this.count--;
                this.isDropHoveringOverTarget = false;
            },
            handleDragCancel() {
                this.isDropHovering = false;
                this.isDropHoveringOverTarget = false;
                this.isDropFromPhosphor = false;
                this.isAlreadyGeneratingHover = false;
                this.isBadDropHovering = false;
                this.dragStarted = false;
            },
            async handleDrop(e) {
                this.count = 0;
                let shouldHandleDrop = false;
                if (this.isDropHovering || (this.isBadDropHovering && (e.dataTransfer.types.includes('text/x-spotify-tracks') || e.dataTransfer.types.includes('text/x-spotify-playlists') || e.dataTransfer.types.includes('text/x-spotify-albums')))) {
                    shouldHandleDrop = true;
                }
                this.isDropHovering = false;
                this.isDropHoveringOverTarget = false;
                this.isDropFromPhosphor = false;
                this.isAlreadyGeneratingHover = false;
                this.isBadDropHovering = false;
                this.dragStarted = false;
                if (!shouldHandleDrop) {
                    return;
                }
                this.$store.commit('loading/startLoad');
                this.$store.commit('loading/playlistGenerating');
                e.preventDefault();
                let trackUrls = e.dataTransfer.getData('text/x-spotify-tracks');
                let playlistUrl = e.dataTransfer.getData('text/x-spotify-playlists');
                let albumUrl = e.dataTransfer.getData('text/x-spotify-albums');
                if (trackUrls) {
                    let trackUrlsArr = trackUrls.split('\n');
                    let trackIds = [];
                    for (let trackUrl of trackUrlsArr) {
                        let trackParts = trackUrl.split('/');
                        let trackId = trackParts[trackParts.length - 1];
                        trackIds.push(trackId);
                    }
                    this.handleTrackDrop(trackIds);
                } else if (playlistUrl) {
                    let playlistParts = playlistUrl.split('/');
                    let playlistId = playlistParts[playlistParts.length - 1];
                    this.handlePlaylistDrop(playlistId);
                } else if (albumUrl) {
                    let albumParts = albumUrl.split('/');
                    let albumId = albumParts[albumParts.length - 1];
                    this.handleAlbumDrop(albumId);
                }
            },
            async handleTrackDrop(trackIds) {
                this.$store.commit('loading/initializeProgress', {id: 'generate'});
                try {
                    let trackResponse;
                    let trackIdsStr = trackIds.join(',');
                    if (this.isLoggedInUser) {
                        trackResponse = await fetch(`${process.env.API_ORIGIN}/track/${trackIdsStr}`, {credentials: 'include'});
                    } else {
                        let captcha = await getCaptchaToken('api/track');
                        trackResponse = await fetch(`${process.env.API_ORIGIN}/track/unauthenticated/${this.$store.getters['user/country']}/${trackIdsStr}?captcha=${captcha}`);
                    }
                    let {tracks} = await trackResponse.json();                    
                    if (tracks.length === 1) {
                        let processedTrack = await processTrack(tracks[0]);
                        let pruners;
                        if (this.$store.state.preferences.onlyTheHits) {
                            pruners = [builders.hits];
                        }
                        let {playlist} = await loadNewPlaylist({
                            count: this.$store.state.preferences.tracksPerPlaylist,
                            builder: builders.randomwalk,
                            firstTrack: processedTrack,
                            pruners
                        }, percent => {
                            this.$store.commit('loading/tickProgress', {id: 'generate', percent});
                        });
                        this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                    } else {
                        let processedTracks = await processTracks(tracks);
                        let {playlist} = await loadNewPlaylist({
                            count: Object.keys(processedTracks.tracks).length,
                            builder: builders.randomwalk,
                            replacementTracks: processedTracks
                        }, percent => {
                            this.$store.commit('loading/tickProgress', {id: 'generate', percent});
                        });
                        this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                    }
                    this.$store.commit('loading/completeProgress', {id: 'generate'});
                    this.$store.commit('loading/resetProgress');
                }
                catch (e) {
                    this.$ga.exception(e);
                    console.error('Playlist generation failed', e);
                    this.$store.dispatch('loading/failProgress');
                }
                this.$store.commit('loading/playlistGenerationComplete');
                this.$store.dispatch('loading/endLoadAfterDelay');
            },
            async handlePlaylistDrop(playlistId) {
                this.$store.commit('loading/initializeProgress', {id: 'generate'});
                try {
                    let playlistResponse;
                    if (this.isLoggedInUser) {
                        playlistResponse = await fetch(`${process.env.API_ORIGIN}/playlist/${playlistId}`, {credentials: 'include'});
                    } else {
                        let captcha = await getCaptchaToken('api/playlist');
                        playlistResponse = await fetch(`${process.env.API_ORIGIN}/playlist/unauthenticated/${this.$store.getters['user/country']}/${playlistId}?captcha=${captcha}`);
                    }
                    let {playlist: {tracks}} = await playlistResponse.json();
                    let processedTracks = await processTracks(tracks);
                    let {playlist} = await loadNewPlaylist({
                        count: Object.keys(processedTracks.tracks).length,
                        builder: builders.randomwalk,
                        replacementTracks: processedTracks
                    }, percent => {
                        this.$store.commit('loading/tickProgress', {id: 'generate', percent});
                    });
                    this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                    this.$store.commit('loading/completeProgress', {id: 'generate'});
                    this.$store.commit('loading/resetProgress');
                }
                catch (e) {
                    this.$ga.exception(e);
                    console.error('Playlist generation failed', e);
                    this.$store.dispatch('loading/failProgress');
                }
                this.$store.commit('loading/playlistGenerationComplete');
                this.$store.dispatch('loading/endLoadAfterDelay');
            },
            async handleAlbumDrop(albumId) {
                this.$store.commit('loading/initializeProgress', {id: 'generate'});
                try {
                    let albumResponse;
                    if (this.isLoggedInUser) {
                        albumResponse = await fetch(`${process.env.API_ORIGIN}/album/${albumId}`, {credentials: 'include'});
                    } else {
                        let captcha = await getCaptchaToken('api/album');
                        albumResponse = await fetch(`${process.env.API_ORIGIN}/album/unauthenticated/${this.$store.getters['user/country']}/${albumId}?captcha=${captcha}`);
                    }
                    let {album: tracks} = await albumResponse.json();
                    let processedTracks = await processTracks(tracks);
                    let {playlist} = await loadNewPlaylist({
                        count: Object.keys(processedTracks.tracks).length,
                        builder: builders.randomwalk,
                        replacementTracks: processedTracks
                    }, percent => {
                        this.$store.commit('loading/tickProgress', {id: 'generate', percent});
                    });
                    this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                    this.$store.commit('loading/completeProgress', {id: 'generate'});
                    this.$store.commit('loading/resetProgress');
                }
                catch (e) {
                    this.$ga.exception(e);
                    console.error('Playlist generation failed', e);
                    this.$store.dispatch('loading/failProgress');
                }
                this.$store.commit('loading/playlistGenerationComplete');
                this.$store.dispatch('loading/endLoadAfterDelay');
            }
        }
    };
</script>
