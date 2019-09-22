<template>
    <section>
        <div class="container" ref="container">
            <span class="label vertical left" title="Chthonicness">chthonic</span>
            <div class="inner">
                <span class="label top" title="Transcendentalness">transcendental</span>
                <div class="canvasWrapper">
                    <canvas ref="canvas" :class="{pointer: hoverTrack !== null}" :title="canvasTitle" @click="handleCanvasClick" @dblclick="handleCanvasDoubleClick" tabindex="0" @keydown.arrow-up="moveCursorUp" @keydown.arrow-down="moveCursorDown" @keydown.enter="handleEnter" @blur="detailsTrack = null">
                        <ol>
                            <li v-for="(track, index) in tracks" @click="seekTrack(index)">
                                "{{track.track.name}}" by {{humanReadableArtists(track.track.artists)}} on album "{{track.track.album.name}}" falls on the chthonic-aethereal axis at {{Math.floor(track.evocativeness.aetherealness * 100)}}% and on the transcendental-primordial axis at {{Math.floor(track.evocativeness.primordialness * 100)}}%. The track is in key {{getHarmonics(track)}} and has {{track.features.tempo}} beats per minute.
                            </li>
                        </ol>
                    </canvas>
                    <div class="dummy" ref="dummy"></div>
                </div>
                <span class="label bottom" title="Primordialness">primordial</span>
            </div>
            <span class="label vertical right" title="Aetherealness">aethereal</span>
        </div>
        <div class="details" :style="{left: detailsOffsetX + 'px', top: detailsOffsetY + 'px'}" ref="details" v-if="detailsTrack">
            <dl>
                <dt>Track</dt>
                <dd><a target="_blank" rel="external noopener" :href="detailsTrack.track.external_urls.spotify">{{detailsTrack.track.name}}</a></dd>
                <dt>Artists</dt>
                <dd>
                    <ol class="artists">
                        <li class="artist" v-for="artist in detailsTrack.track.artists">
                            <a target="_blank" rel="external noopener" :href="artist.external_urls.spotify">{{artist.name}}</a>
                        </li>
                    </ol>
                </dd>
                <dt>Album</dt>
                <dd><a target="_blank" rel="external noopener" :href="detailsTrack.track.album.external_urls.spotify">{{detailsTrack.track.album.name}}</a></dd>
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

    .container {
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: rgb(40, 27, 61, 0.7);
        border: 2px solid aqua;
        height: 100%;
        box-sizing: border-box;
        width: 65vh;
        max-height: 65vh;
        margin-bottom: 1em;
    }

    .inner {
        display: flex;
        flex-direction: column;
        flex: 1;
        height: 100%;
    }

    canvas.pointer {
        cursor: pointer;
    }

    canvas:focus {
        outline: none;
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

    .canvasWrapper {
        position: relative;
        display: flex;
        align-items: center;
        justify-content: center;
        flex: 1;
    }

    .dummy {
        flex: 1;
        box-sizing: border-box;
        touch-action: none;
        height: 100%;
    }

    canvas {
        position: absolute;
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

    @media only screen and (max-width: 1099px) {
        .container {
            width: 58vh;
            max-height: 58vh;
        }
    }

    @media only screen and (max-width: 599px) {
        section {
            justify-content: flex-start;
            flex-direction: column;
        }

        .container {
            width: 45vh;
            max-height: 45vh;
        }

        dl {
            white-space: normal;
        }
    }

    @media only screen and (max-width: 414px) {
        .container {
            width: 40vh;
            max-height: 40vh;
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
    import {initializeCanvas} from '~/assets/constellation';

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
    const segmentBreaks = [[400, 10], [350, 8], [150, 6], [0, 4]];
    const outerSizeDifference = 50;
    const padding = outerSizeDifference / 2;
    const pointBoundingBox = 30;
    const pointBoundingRadius = pointBoundingBox / 2;
    const pulseTTL = 1000;
    const maxHarmonicDistance = 7;
    const dragThreshold = 20;

    export default {
        data() {
            return {
                hoverTrack: null,
                detailsTrack: null,
                innerSize: null,
                destroyCanvas: null,
                haloInterval: null,
                haloRotation: 0,
                beatPulses: [],
                beatPulseConsumer: [],
                beatIntervals: [],
                edges: [],
                canvasAbsoluteX: null,
                canvasAbsoluteY: null,
                ongoingTouch: null
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
            canvasTitle() {
                if (typeof this.hoverTrack == 'undefined' || this.hoverTrack === null) {
                    return 'Evocativeness constellation';
                }
                let track = this.tracks[this.hoverTrack];
                return `${track.track.name} - ${track.track.artists.map(artist => artist.name).join(', ')} - ${track.track.album.name} (${this.getHarmonics(track)}, ${track.features.tempo} BPM)`;
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
            }
        },
        watch: {
            innerSize(newVal) {
                this.destroyCanvas && this.destroyCanvas();
                if (!newVal) {
                    return;
                }
                this.destroyCanvas = initializeCanvas(this.$refs.canvas, this, fps);
            },
            tracks: {
                immediate: true,
                handler(newTracks) {
                    this.destroyBeatIntervals();
                    this.beatPulses = [];
                    this.beatPulseConsumer = [];
                    this.edges = [];
                    this.detailsTrack = null;
                    this.hoverTrack = null;
                    let lastTrack;
                    for (let track of newTracks) {
                        let beat = (1 / (track.features.tempo / 60)) * 1000;
                        this.beatIntervals.push(setInterval(() => {
                            this.beatPulseConsumer.push({start: Date.now(), x: track.evocativeness.aetherealness, y: track.evocativeness.primordialness, opacity: 1, radiusMultiplier: 1, trackId: track.track.id});
                        }, beat));
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
            this.$refs.canvas.addEventListener('mousemove', this.handleCanvasMouseMove);
            this.$refs.canvas.addEventListener('mouseleave', this.handleCanvasMouseLeave);
            this.$refs.canvas.addEventListener('touchstart', this.handleTouchStart, {passive: true});
            this.$refs.canvas.addEventListener('touchmove', this.handleTouchMove, {passive: true});
            this.$refs.canvas.addEventListener('touchend', this.handleTouchEnd, {passive: true});
            this.$refs.canvas.addEventListener('touchcancel', this.handleTouchEnd, {passive: true});
            addEventListener('resize', this.handleResize);
            addEventListener('orientationchange', this.handleResizeAfterTimeout);
            this.handleResize();
            this.$nextTick(this.resetCanvasBounds);
        },
        created() {
            this.haloInterval = setInterval(this.updateHalos, 1000 / (fps * 2));
        },
        beforeDestroy() {
            clearInterval(this.haloInterval);
            this.destroyBeatIntervals();
            this.$refs.canvas.removeEventListener('mousemove', this.handleCanvasMouseMove);
            this.$refs.canvas.removeEventListener('mouseleave', this.handleCanvasMouseLeave);
            this.$refs.canvas.removeEventListener('touchstart', this.handleTouchStart, {passive: true});
            this.$refs.canvas.removeEventListener('touchmove', this.handleTouchMove, {passive: true});
            this.$refs.canvas.removeEventListener('touchend', this.handleTouchEnd, {passive: true});
            this.$refs.canvas.removeEventListener('touchcancel', this.handleTouchEnd, {passive: true});
            removeEventListener('resize', this.handleResize);
            removeEventListener('orientationchange', this.handleResizeAfterTimeout);
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
                let now = Date.now();
                let newBeatPulses = [];
                for (let beatPulse of this.beatPulses) {
                    if (beatPulse.start + pulseTTL > now) {
                        let percentDone = (now - beatPulse.start) / pulseTTL;
                        beatPulse.opacity = 1 - percentDone;
                        beatPulse.radiusMultiplier = 1 + percentDone;
                        newBeatPulses.push(beatPulse);
                    }
                }
                for (let beatPulse of this.beatPulseConsumer) {
                    newBeatPulses.push(beatPulse);
                }
                this.beatPulseConsumer = [];
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
                this.detailsTrack = this.tracks[this.hoverTrack];
            },
            handleCanvasDoubleClick() {
                if (typeof this.hoverTrack == 'undefined' || this.hoverTrack === null) {
                    return;
                }
                this.$store.dispatch('tracks/seekTrack', this.hoverTrack);
            },
            moveCursorUp() {
                if (typeof this.hoverTrack == 'undefined' || this.hoverTrack === null) {
                    this.hoverTrack = 0;
                    return;
                }
                if (this.hoverTrack === 0) {
                    return;
                }
                this.hoverTrack = Math.max(0, this.hoverTrack - 1);
            },
            moveCursorDown() {
                if (typeof this.hoverTrack == 'undefined' || this.hoverTrack === null) {
                    this.hoverTrack = 0;
                    return;
                }
                if (this.hoverTrack === this.tracks.length - 1) {
                    return;
                }
                this.hoverTrack = Math.min(this.tracks.length - 1, this.hoverTrack + 1);
            },
            handleEnter() {
                this.handleCanvasDoubleClick();
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
            handleResizeAfterTimeout() {
                setTimeout(this.handleResize, 200);
            },
            handleResize() {
                let dummy = this.$refs.dummy;
                let size = Math.min(dummy.offsetWidth, dummy.offsetHeight) - outerSizeDifference;
                this.innerSize = Math.min(maxSize, Math.max(minSize, size));
                if (this.$refs.container.offsetHeight >= 300 && this.$refs.container.offsetWidth >= 300) {
                    this.$refs.container.style.fontSize = `${Math.max(8, Math.floor((size / maxSize) * 16))}px`;
                } else {
                    this.$refs.container.style.fontSize = '0px';
                }
                this.resetCanvasBounds();
            },
            resetCanvasBounds() {
                let {x, y} = this.$refs.canvas.getBoundingClientRect();
                this.canvasAbsoluteX = x;
                this.canvasAbsoluteY = y;
            },
            destroyBeatIntervals() {
                if (this.beatIntervals) {
                    for (let beatInterval of this.beatIntervals) {
                        clearInterval(beatInterval);
                    }
                }
                this.beatIntervals = [];
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
                    } else {
                        this.moveCursorUp();
                    }
                    this.detailsTrack = this.tracks[this.hoverTrack];
                    this.ongoingTouch = newTouch;
                }
            },
            handleTouchEnd({touches}) {
                this.ongoingTouch = null;
            }
        }
    };
</script>