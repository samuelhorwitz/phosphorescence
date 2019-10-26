<template>
    <section ref="container">
        <div class="canvasWrapper">
            <canvas
                ref="canvas"
                :class="{pointer: hoverTrack !== null}"
                :title="canvasTitle"
                @click="handleCanvasClick(); $ga.event('constellation', 'click', trackHovering ? 'track' : 'nothing')"
                @dblclick="handleCanvasDoubleClick(); $ga.event('constellation', 'double-click', trackHovering ? 'track' : 'nothing')"
                tabindex="0"
                @keydown.arrow-up="moveCursorUp(); $ga.event('constellation', 'key', 'up')"
                @keydown.arrow-down="moveCursorDown(); $ga.event('constellation', 'key', 'down')"
                @keydown.enter="handleEnter(); $ga.event('constellation', 'key', 'enter')"
                @keydown.esc="handleEscape(); $ga.event('constellation', 'key', 'escape')"
                @focusout="handleBlur"
                @mousemove="handleCanvasMouseMove"
                @mouseleave="handleCanvasMouseLeave"
                @touchstart.passive="handleTouchStart"
                @touchmove.prevent="handleTouchMove"
                @touchend.passive="handleTouchEnd"
                @touchcancel.passive="handleTouchEnd"
                v-spotify-uri:track="canvasSpotifyDragTrack"
                v-spotify-uri-title="canvasSpotifyDragTrackTitle">
                <ol>
                    <li v-for="(track, index) in tracks" @click="seekTrack(index); $ga.event('constellation-accessible', 'seek-track')">
                        "{{track.track.name}}" by {{humanReadableArtists(track.track.artists)}} on album "{{track.track.album.name}}" falls on the chthonic-aethereal axis at {{Math.floor(track.evocativeness.aetherealness * 100)}}% and on the transcendental-primordial axis at {{Math.floor(track.evocativeness.primordialness * 100)}}%. The track is in key {{getHarmonics(track)}} and has {{track.features.tempo}} beats per minute.
                    </li>
                </ol>
            </canvas>
            <div class="hover" :class="{show: showAxisLabels}">
                <span class="label vertical left" title="Chthonicness">chthonic</span>
                <div class="hoverInner">
                    <span class="label top" title="Transcendentalness">transcendental</span>
                    <span class="label bottom" title="Primordialness">primordial</span>
                </div>
                <span class="label vertical right" title="Aetherealness">aethereal</span>
            </div>
        </div>
        <div class="details" :style="{left: detailsOffsetX + 'px', top: detailsOffsetY + 'px'}" ref="details" v-if="detailsTrack">
            <dl>
                <dt>Track</dt>
                <dd><a tabindex="0" target="_blank" rel="external noopener" :href="detailsTrackUrl" v-spotify-uri:track="detailsTrack.track.id" v-spotify-uri-title="getSpotifyTrackDragTitle(detailsTrack)">{{detailsTrack.track.name}}</a></dd>
                <dt>Artists</dt>
                <dd>
                    <ol class="artists">
                        <li class="artist" v-for="artist in detailsTrackArtists">
                            <a tabindex="0" target="_blank" rel="external noopener" :href="artist.url" v-spotify-uri:artist="artist.id" v-spotify-uri-title="artist.name">{{artist.name}}</a>
                        </li>
                    </ol>
                </dd>
                <dt>Album</dt>
                <dd><a tabindex="0" target="_blank" rel="external noopener" :href="detailsTrackAlbumUrl" v-spotify-uri:album="detailsTrack.track.album.id" v-spotify-uri-title="getSpotifyAlbumDragTitle(detailsTrack.track.album)">{{detailsTrack.track.album.name}}</a></dd>
                <dt>Key</dt>
                <dd>{{getHarmonics(detailsTrack)}}</dd>
                <dt>BPM</dt>
                <dd>{{detailsTrack.features.tempo}}</dd>
            </dl>
        </div>
    </section>
</template>

<style scoped>
    section {
        max-height: 100%;
        display: flex;
        flex: 1;
        align-items: center;
        justify-content: center;
        height: 100%;
    }

    canvas {
        -webkit-touch-callout: none;
        -webkit-user-select: none;
    }

    canvas.pointer {
        cursor: pointer;
    }

    canvas:focus {
        outline: none;
    }

    .canvasWrapper {
        background-color: rgb(40, 27, 61, 0.7);
        border: 2px solid aqua;
        box-sizing: border-box;
        position: relative;
    }

    .hover {
        display: flex;
        position: absolute;
        top: 0px;
        left: 0px;
        width: 100%;
        height: 100%;
        background-image: radial-gradient(circle, rgba(0, 0, 0 ,0) 51%, rgba(0, 0, 28, 1) 100%);
        pointer-events: none;
        transition: opacity 0.3s ease-in 0s;
        opacity: 0;
    }

    .hover.show {
        opacity: 1;
    }

    .hoverInner {
        display: flex;
        flex-direction: column;
        flex: 1;
        justify-content: space-between;
    }

    .label {
        text-align: center;
        color: white;
        font-family: VT323;
        font-size: 3.5em;
    }

    .label.vertical {
        writing-mode: vertical-rl;
        text-orientation: upright;
    }

    .label.left {
        margin-right: -0.3em;
    }

    .label.right {
        margin-left: -0.3em;
    }

    .label.top {
        margin-bottom: -0.3em;
    }

    .label.bottom {
        margin-top: -0.4em;
        margin-bottom: 0.1em;
    }

    .details {
        display: flex;
        font-size: 16px;
        background-color: rgb(40, 27, 61, 0.7);
        border: 2px solid aqua;
        box-sizing: border-box;
        color: white;
        padding: 1em;
        width: 100%;
        overflow-y: auto;
    }

    a {
        color: aqua;
        text-decoration: none;
    }

    a:hover {
        text-decoration: underline;
    }

    a:visited {
        color: aqua;
    }

    dl {
        margin: 0px;
        padding: 0px;
        width: 100%;
        white-space: nowrap;
    }

    dt {
        display: inline;
        float: left;
        clear: both;
        width: 25%
    }

    dd {
        display: inline;
        margin: 0px;
        float: left;
        width: 75%;
        text-overflow: ellipsis;
        overflow-x: hidden;
    }

    ol.artists {
        list-style: none;
        margin: 0px;
        padding: 0px;
        overflow-x: hidden;
        text-overflow: ellipsis;
    }

    ol.artists li {
        display: inline;
    }

    ol.artists li:not(:last-child):after {
        content: ',';
        padding-right: 1ex;
    }

    @media only screen and (min-width: 600px) and (pointer: fine) {
        .details {
            position: absolute;
            width: 20em;
            background-color: rgb(40, 27, 61);
            z-index: 9999999;
            border: 4px solid aqua;
            overflow-y: unset;
        }

        .details:after, .details:before {
            right: 100%;
            top: 50%;
            border: solid transparent;
            content: " ";
            height: 0;
            width: 0;
            position: absolute;
            pointer-events: none;
        }

        .details:after {
            border-color: rgba(136, 183, 213, 0);
            border-right-color: rgb(40, 27, 61);
            border-width: 30px;
            margin-top: -30px;
        }
        .details:before {
            border-color: rgba(194, 225, 245, 0);
            border-right-color: aqua;
            border-width: 36px;
            margin-top: -36px;
        }
    }

    @media only screen and (max-width: 599px) {
        section {
            justify-content: flex-start;
            flex-direction: column;
        }

        dl {
            white-space: normal;
        }

        .details {
            margin-top: 1ex;
        }
    }

    @media only screen and (orientation: landscape) and (pointer: coarse) {
        dl {
            white-space: normal;
        }

        .details {
            margin-left: 1em;
        }
    }
</style>

<script>
    import {getSpotifyAlbumUrl, getSpotifyArtistUrl, getSpotifyTrackUrl, getSpotifyTrackDragTitle, getSpotifyAlbumDragTitle} from '~/assets/spotify';
    import {initializeCanvas} from '~/assets/constellation';
    import mainViewportEventBus from '~/assets/mainviewport';

    const A_FLAT = 8, E_FLAT = 3, B_FLAT = 10, F = 5, C = 0, G = 7, D = 2, A = 9, E = 4, B = 11, F_SHARP = 6, D_FLAT = 1;
    const key = ['C', 'C♯', 'D', 'D♯', 'E', 'F', 'F♯', 'G', 'G♯', 'A', 'A♯', 'B'];
    const mode = ['minor', 'major'];
    const minorsCircle = [A_FLAT, E_FLAT, B_FLAT, F, C, G, D, A, E, B, F_SHARP, D_FLAT];
    const majorsCircle = [B, F_SHARP, D_FLAT, A_FLAT, E_FLAT, B_FLAT, F, C, G, D, A, E];
    const minorsPitchToPositionMap = minorsCircle.reduce((acc, cur, i) => {acc[cur] = i; return acc;}, []);
    const majorsPitchToPositionMap = majorsCircle.reduce((acc, cur, i) => {acc[cur] = i; return acc;}, []);
    const fps = 60;
    const maxSize = 500;
    const minSize = 50;
    const segmentBreaks = [[600, 14], [500, 12], [400, 10], [250, 8], [150, 6], [0, 4]];
    const outerSizeDifference = 50;
    const padding = outerSizeDifference / 2;
    const pointBoundingBox = 30;
    const pointBoundingRadius = pointBoundingBox / 2;
    const pulseTTL = 1000;
    const maxHarmonicDistance = 7;
    const dragThreshold = 20;
    const canvasMargin = 50;

    export default {
        data() {
            return {
                hoverTrack: null,
                detailsTrack: null,
                innerSize: null,
                destroyCanvas: null,
                markConstellationAsDirty: null,
                haloInterval: null,
                haloRotation: 0,
                beatPulses: [],
                lastBeat: {},
                edges: [],
                canvasAbsoluteX: null,
                canvasAbsoluteY: null,
                ongoingTouch: null,
                showAxisLabels: false
            };
        },
        computed: {
            outerSize() {
                return this.innerSize + outerSizeDifference;
            },
            gridSegments() {
                for (let segmentBreak of segmentBreaks) {
                    if (this.innerSize >= segmentBreak[0]) {
                        return segmentBreak[1];
                    }
                }
                return segmentBreaks[segmentBreaks.length - 1][1];
            },
            gridInterval() {
                return this.innerSize / this.gridSegments;
            },
            totalGridlines() {
                return (this.innerSize / this.gridInterval) - 1;
            },
            centerGridlineIndex() {
                return (this.gridSegments / 2) - 1;
            },
            tracks() {
                return this.$store.state.tracks.playlist;
            },
            currentTrack() {
                if (this.$store.getters['tracks/isPlayerDisconnected']) {
                    return null;
                }
                return this.$store.getters['tracks/currentTrack'];
            },
            trackHovering() {
                return !(typeof this.hoverTrack == 'undefined' || this.hoverTrack === null);
            },
            canvasTitle() {
                if (!this.trackHovering) {
                    return 'Evocativeness constellation';
                }
                let track = this.tracks[this.hoverTrack];
                return `${track.track.name} - ${track.track.artists.map(artist => artist.name).join(', ')} - ${track.track.album.name} (${this.getHarmonics(track)}, ${track.features.tempo} BPM)`;
            },
            canvasSpotifyDragTrack() {
                if (!this.trackHovering) {
                    return this.tracks.map(track => track.id);
                }
                return this.tracks[this.hoverTrack].id;
            },
            canvasSpotifyDragTrackTitle() {
                if (!this.trackHovering) {
                    return `${this.tracks.length} track${this.tracks.length == 1 ? '' : 's'}`;
                }
                return getSpotifyTrackDragTitle(this.tracks[this.hoverTrack]);
            },
            detailsOffsetX() {
                if (!this.detailsTrack || !this.innerSize) {
                    return 0;
                }
                return this.canvasAbsoluteX + (this.detailsTrack.evocativeness.aetherealness * this.innerSize) + 80;
            },
            detailsOffsetY() {
                if (!this.detailsTrack || !this.innerSize) {
                    return 0;
                }
                return (this.canvasAbsoluteY + (this.detailsTrack.evocativeness.primordialness * this.innerSize)) - 42;
            },
            detailsTrackUrl() {
                if (!this.detailsTrack) {
                    return null;
                }
                return getSpotifyTrackUrl(this.detailsTrack.id);
            },
            detailsTrackAlbumUrl() {
                if (!this.detailsTrack) {
                    return null;
                }
                return getSpotifyAlbumUrl(this.detailsTrack.track.album.id);
            },
            detailsTrackArtists() {
                if (!this.detailsTrack) {
                    return null;
                }
                let artists = [];
                for (let artist of this.detailsTrack.track.artists) {
                    artists.push({...artist, url: getSpotifyArtistUrl(artist.id)});
                }
                return artists;
            }
        },
        watch: {
            innerSize(newVal) {
                this.destroyCanvas && this.destroyCanvas();
                if (!newVal) {
                    return;
                }
                let {destroyer, redrawConstellation} = initializeCanvas(this.$refs.canvas, this, fps);
                this.destroyCanvas = destroyer;
                this.markConstellationAsDirty = redrawConstellation;
            },
            tracks: {
                immediate: true,
                handler(newTracks) {
                    this.markConstellationAsDirty && this.markConstellationAsDirty();
                    this.beatPulses = [];
                    this.edges = [];
                    this.detailsTrack = null;
                    this.hoverTrack = null;
                    let lastTrack;
                    for (let track of newTracks) {
                        if (lastTrack) {
                            let diff = this.calculateHarmonicDifference(lastTrack.features, track.features);
                            this.edges.push(Math.min(5, Math.max(3, 5 * (1 - (diff / (maxHarmonicDistance + 1))))));
                        }
                        lastTrack = track;
                    }
                }
            }
        },
        mounted() {
            addEventListener('resize', this.handleResize);
            addEventListener('orientationchange', this.handleResizeAfterTick);
            mainViewportEventBus.$on('resize', this.handleResizeAfterTick);
            this.handleResize();
            this.$nextTick(this.resetCanvasBounds);
        },
        created() {
            this.haloInterval = setInterval(this.updateHalos, 1000 / (fps * 2));
        },
        beforeDestroy() {
            clearInterval(this.haloInterval);
            removeEventListener('resize', this.handleResize);
            removeEventListener('orientationchange', this.handleResizeAfterTick);
            mainViewportEventBus.$off('resize', this.handleResizeAfterTick);
            this.destroyCanvas && this.destroyCanvas();
        },
        methods: {
            updateHalos() {
                this.updateHaloRotation();
                this.updatePulseHalos();
            },
            updateHaloRotation() {
                this.haloRotation++;
            },
            updatePulseHalos() {
                let newBeatPulses = [];
                for (let beatPulse of this.beatPulses) {
                    if (beatPulse.start + beatPulse.ttl > performance.now()) {
                        newBeatPulses.push(beatPulse);
                    }
                }
                for (let track of this.tracks) {
                    let lastBeat = this.lastBeat[track.id];
                    let start = performance.now();
                    if (lastBeat) {
                        let beat = (1 / (track.features.tempo / 60)) * 1000;
                        let timeSince = start - lastBeat.start;
                        if (timeSince < beat) {
                            continue;
                        }
                        let diff = timeSince - beat;
                        start -= diff;
                    }
                    let beat = {start, x: track.evocativeness.aetherealness, y: track.evocativeness.primordialness, ttl: pulseTTL, trackId: track.id};
                    newBeatPulses.push(beat);
                    this.lastBeat[track.id] = beat;
                }
                this.beatPulses = newBeatPulses;
            },
            calculateHarmonicDifference(a, b) {
                if (!a.mode || !a.key || !b.mode || !b.key) {
                    return 1;
                }
                let diff;
                if (a.mode == b.mode) {
                    // We can use major here even if they are both minor because
                    // distance is the same regardless of offset on the wheel.
                    diff = Math.abs(majorsPitchToPositionMap[a.key] - majorsPitchToPositionMap[b.key]);
                    if (diff > 12 / 2) {
                        diff = 12 % diff;
                    }
                }
                else {
                    let diff;
                    if (a.mode == MINOR && b.mode == MAJOR) {
                        diff = Math.abs(minorsPitchToPositionMap[a.key] - majorsPitchToPositionMap[b.key]);
                    }
                    else {
                        diff = Math.abs(majorsPitchToPositionMap[a.key] - minorsPitchToPositionMap[b.key]);
                    }
                    if (diff > 12 / 2) {
                        diff = 12 % diff;
                    }
                    diff += 1;
                }
                return diff;
            },
            handleCanvasMouseMove(e) {
                let rect = e.target.getBoundingClientRect();
                let x = e.clientX - rect.left;
                let y = e.clientY - rect.top;
                let ctx = e.target.getContext('2d');
                this.hoverTrack = null;
                for (let i in this.tracks) {
                    let track = this.tracks[i];
                    let trackX = (track.evocativeness.aetherealness * this.innerSize) + padding;
                    let trackY = (track.evocativeness.primordialness * this.innerSize) + padding;
                    if ((y < trackY - pointBoundingRadius) ||
                       (y > trackY + pointBoundingRadius) ||
                       (x < trackX - pointBoundingRadius) ||
                       (x > trackX + pointBoundingRadius)) {
                        continue;
                    }
                    this.hoverTrack = i;
                    break;
                }
            },
            handleCanvasMouseLeave() {
                this.hoverTrack = null;
            },
            handleCanvasClick() {
                let oldDetailsTrack = this.detailsTrack;
                this.detailsTrack = this.tracks[this.hoverTrack];
                if (this.detailsTrack || (oldDetailsTrack && !this.detailsTrack)) {
                    this.showAxisLabels = false;
                } else {
                    this.showAxisLabels = !this.showAxisLabels;
                }
            },
            handleCanvasDoubleClick() {
                if (!this.trackHovering) {
                    return;
                }
                this.$store.dispatch('tracks/seekTrack', this.hoverTrack);
            },
            moveCursorUp() {
                if (!this.trackHovering) {
                    this.hoverTrack = 0;
                    return;
                }
                if (this.hoverTrack === 0) {
                    return;
                }
                this.hoverTrack = Math.max(0, this.hoverTrack - 1);
            },
            moveCursorDown() {
                if (!this.trackHovering) {
                    this.hoverTrack = 0;
                    return;
                }
                if (this.hoverTrack === this.tracks.length - 1) {
                    return;
                }
                this.hoverTrack = Math.min(this.tracks.length - 1, this.hoverTrack + 1);
            },
            handleEnter() {
                this.handleCanvasClick();
            },
            handleEscape() {
                this.$refs.container.blur();
            },
            handleBlur(e) {
                if (this.$refs.container.contains(e.relatedTarget)) {
                    e.preventDefault();
                    return;
                }
                this.detailsTrack = null;
                this.showAxisLabels = false;
            },
            humanReadableArtists(artists) {
                if (artists.length === 1) {
                    return `artist "${artists[0].name}"`;
                }
                let artistsStr = artists.slice(0, -1).map(artist => `"${artist.name}"`).join(', ');
                for (let artist of artists) {
                    artistsStr += ''
                }
                return `artists ${artistsStr} and "${artists[artists.length - 1].name}"`;
            },
            getHarmonics(track) {
                return `${key[track.features.key]} ${mode[track.features.mode]}`;
            },
            handleResizeAfterTick() {
                this.$nextTick(this.handleResize);
                setTimeout(this.resetCanvasBounds, 200);
            },
            handleResize() {
                this.innerSize = Math.min(this.$refs.container.parentNode.clientWidth - canvasMargin - (outerSizeDifference / 2), this.$refs.container.parentNode.clientHeight - canvasMargin - (outerSizeDifference / 2));
                this.$refs.container.style.fontSize = `${Math.max(8, Math.floor((this.innerSize / maxSize) * 16))}px`;
                this.resetCanvasBounds();
            },
            resetCanvasBounds() {
                let {x, y} = this.$refs.canvas.getBoundingClientRect();
                this.canvasAbsoluteX = x;
                this.canvasAbsoluteY = y;
            },
            handleTouchStart({touches}) {
                if (touches.length > 1 || this.ongoingTouch) {
                    return;
                }
                this.$refs.canvas.focus();
                this.ongoingTouch = touches[0];
            },
            handleTouchMove({touches}) {
                if (!this.ongoingTouch) {
                    return;
                }
                let newTouch;
                for (let touch of touches) {
                    if (touch.identifier === this.ongoingTouch.identifier) {
                        newTouch = touch;
                        break;
                    }
                }
                if (!newTouch) {
                    return;
                }
                if (Math.abs(newTouch.screenY - this.ongoingTouch.screenY) > dragThreshold) {
                    if (newTouch.screenY > this.ongoingTouch.screenY) {
                        this.moveCursorDown();
                        this.$ga.event('constellation', 'touch-scroll', 'down');
                    } else {
                        this.moveCursorUp();
                        this.$ga.event('constellation', 'touch-scroll', 'up');
                    }
                    this.detailsTrack = this.tracks[this.hoverTrack];
                    this.ongoingTouch = newTouch;
                }
            },
            handleTouchEnd({touches}) {
                this.ongoingTouch = null;
            },
            getSpotifyTrackDragTitle,
            getSpotifyAlbumDragTitle
        }
    };
</script>