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

    function initializeCanvas(canvasEl, component) {
        let canvas = canvasEl;
        let ctx = canvas.getContext('2d');
        resizeCanvas(canvas, ctx, component);
        return beginLoop(ctx, component);
    }

    function resizeCanvas(canvas, ctx, component) {
        let outerSize = component.outerSize;
        canvas.width = outerSize * devicePixelRatio;
        canvas.height = outerSize * devicePixelRatio;
        canvas.style.width = `${outerSize}px`;
        canvas.style.height = `${outerSize}px`;
        ctx.scale(devicePixelRatio, devicePixelRatio);
    }

    function beginLoop(ctx, component) {
        let keepLooping = true;
        (async () => {
            while (keepLooping) {
                requestAnimationFrame(() => {
                    ctx.clearRect(0, 0, component.outerSize, component.outerSize);
                    ctx.save();
                    ctx.translate(padding, padding);
                    drawGrid(ctx, component);
                    drawTracks(ctx, component);
                    ctx.restore();
                });
                await new Promise(resolve => setTimeout(resolve, 1000 / fps));
            }
        })();
        return () => {
            keepLooping = false;
            ctx.clearRect(0, 0, component.outerSize, component.outerSize);
        };
    }

    function drawGrid(ctx, component) {
        let innerSize = component.innerSize;
        let totalGridlines = component.totalGridlines;
        let centerGridlineIndex = component.centerGridlineIndex;
        let gridInterval = component.gridInterval;
        let gridSegments = component.gridSegments;
        ctx.save();
        ctx.lineWidth = 4;
        ctx.strokeStyle = 'white';
        ctx.shadowBlur = 20;
        ctx.shadowColor = 'magenta';
        ctx.beginPath();
        ctx.moveTo(0, 0);
        ctx.lineTo(innerSize, 0);
        ctx.lineTo(innerSize, innerSize);
        ctx.lineTo(0, innerSize);
        ctx.closePath();
        ctx.stroke();
        ctx.beginPath();
        for (let i = 0; i < totalGridlines; i++) {
            if (i === centerGridlineIndex) {
                // we draw the middle ones separately
                continue;
            }
            ctx.moveTo(0, (i + 1) * gridInterval);
            ctx.lineTo(innerSize, (i + 1) * gridInterval);
            ctx.moveTo((i + 1) * gridInterval, 0);
            ctx.lineTo((i + 1) * gridInterval, innerSize);
        }
        ctx.closePath();
        ctx.stroke();
        ctx.lineWidth = 7;
        ctx.strokeStyle = 'white';
        ctx.shadowBlur = 20;
        ctx.shadowColor = 'cyan';
        ctx.beginPath();
        ctx.moveTo(0, (gridSegments / 2) * gridInterval);
        ctx.lineTo(innerSize, (gridSegments / 2) * gridInterval);
        ctx.moveTo((gridSegments / 2) * gridInterval, 0);
        ctx.lineTo((gridSegments / 2) * gridInterval, innerSize);
        ctx.closePath();
        ctx.stroke();
        ctx.restore();
    }

    function drawTracks(ctx, component) {
        for (let i = 0; i < component.tracks.length - 1; i++) {
            let track = component.tracks[i];
            let nextTrack = component.tracks[i+1];
            drawConnector(ctx, component, component.edges[i], track.evocativeness.aetherealness, track.evocativeness.primordialness,
                nextTrack.evocativeness.aetherealness, nextTrack.evocativeness.primordialness);
        }
        let haloX, haloY;
        let hoverHaloX, hoverHaloY;
        let detailsHaloX, detailsHaloY;
        let activeId, hoverId, detailsId;
        for (let track of component.tracks) {
            drawTrack(ctx, component, track.evocativeness.aetherealness, track.evocativeness.primordialness);
            if (component.currentTrack && component.currentTrack.track.id === track.track.id) {
                activeId = track.track.id;
                haloX = track.evocativeness.aetherealness;
                haloY = track.evocativeness.primordialness;
            }
            if (component.hoverTrack !== null && track.track.id === component.tracks[component.hoverTrack].track.id) {
                hoverId = track.track.id;
                hoverHaloX = track.evocativeness.aetherealness;
                hoverHaloY = track.evocativeness.primordialness;
                continue;
            }
            if (!!component.detailsTrack && track.track.id === component.detailsTrack.track.id) {
                detailsId = track.track.id;
                detailsHaloX = track.evocativeness.aetherealness;
                detailsHaloY = track.evocativeness.primordialness;
            }
        }
        drawBeatPulses(ctx, component, activeId, hoverId, detailsId, component.beatPulses);
        if (haloX && haloY) {
            drawActiveHalo(ctx, component, haloX, haloY);
        }
        if (hoverHaloX && hoverHaloY) {
            drawHoverHalo(ctx, component, hoverHaloX, hoverHaloY);
        }
        if (detailsHaloX && detailsHaloY) {
            drawHoverHalo(ctx, component, detailsHaloX, detailsHaloY);
        }
    }

    function drawConnector(ctx, component, edge, x1, y1, x2, y2) {
        let innerSize = component.innerSize;
        ctx.save();
        ctx.beginPath();
        ctx.moveTo(x1 * innerSize, y1 * innerSize);
        ctx.lineTo(x2 * innerSize, y2 * innerSize);
        ctx.globalCompositeOperation = 'screen';
        ctx.strokeStyle = 'lightcyan';
        ctx.lineWidth = Math.min(5, Math.max(3, 5 * (1 - (edge / (maxHarmonicDistance + 1)))));
        ctx.shadowBlur = 30;
        ctx.shadowColor = 'cyan';
        ctx.stroke();
        ctx.restore();
    }

    function drawTrack(ctx, component, x, y) {
        let innerSize = component.innerSize;
        ctx.save();
        ctx.beginPath();
        ctx.arc(x * innerSize, y * innerSize, 6, 0, 2 * Math.PI);
        ctx.closePath();
        ctx.fillStyle = 'white';
        ctx.strokeStyle = 'lightcyan';
        ctx.shadowColor = 'cyan';
        ctx.lineWidth = 6;
        ctx.shadowBlur = 30;
        ctx.globalCompositeOperation = 'screen';
        ctx.fill();
        ctx.stroke();
        ctx.restore();
    }

    function drawBeatPulses(ctx, component, activeId, hoverId, detailsId, beatPulses) {
        let innerSize = component.innerSize;
        ctx.save();
        ctx.globalCompositeOperation = 'screen';
        ctx.shadowColor = 'white';
        ctx.shadowBlur = 30;
        ctx.lineWidth = 5;
        for (let beatPulse of beatPulses) {
            if (beatPulse.trackId !== activeId && beatPulse.trackId !== hoverId && beatPulse.trackId !== detailsId) {
                continue;
            }
            ctx.beginPath();
            ctx.arc(beatPulse.x * innerSize, beatPulse.y * innerSize, 6 * beatPulse.radiusMultiplier * beatPulse.radiusMultiplier * beatPulse.radiusMultiplier, 0, 2 * Math.PI);
            ctx.closePath();
            ctx.strokeStyle = `rgba(255, 255, 255, ${beatPulse.opacity * 0.5})`;
            ctx.stroke();
        }
        ctx.restore();
    }

    function drawActiveHalo(ctx, component, x, y) {
        let innerSize = component.innerSize;
        ctx.save();
        ctx.translate(x * innerSize, y * innerSize);
        ctx.rotate((component.haloRotation * 0.5) * Math.PI / 180);
        ctx.translate(-x * innerSize, -y * innerSize);
        ctx.beginPath();
        ctx.arc(x * innerSize, y * innerSize, 22, 0, 2 * Math.PI);
        ctx.closePath();
        ctx.globalCompositeOperation = 'source-over';
        ctx.strokeStyle = 'magenta';
        ctx.lineWidth = 5;
        ctx.setLineDash([17, 5]);
        ctx.stroke();
        ctx.restore();
    }

    function drawHoverHalo(ctx, component, x, y) {
        let innerSize = component.innerSize;
        ctx.save();
        ctx.translate(x * innerSize, y * innerSize);
        ctx.rotate(-(component.haloRotation * 0.65) * Math.PI / 180);
        ctx.translate(-x * innerSize, -y * innerSize);
        ctx.beginPath();
        ctx.arc(x * innerSize, y * innerSize, 15, 0, 2 * Math.PI);
        ctx.closePath();
        ctx.globalCompositeOperation = 'source-over';
        ctx.strokeStyle = 'cyan';
        ctx.lineWidth = 5;
        ctx.setLineDash([7, 2]);
        ctx.stroke();
        ctx.restore();
    }

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
                canvasAbsoluteY: null
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
                if (this.hoverTrack === null) {
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
                this.destroyCanvas = initializeCanvas(this.$refs.canvas, this);
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
                            this.edges.push(this.calculateHarmonicDifference(lastTrack.features, track.features));
                        }
                        lastTrack = track;
                    }
                }
            }
        },
        mounted() {
            this.$refs.canvas.addEventListener('mousemove', this.handleCanvasMouseMove);
            this.$refs.canvas.addEventListener('mouseleave', this.handleCanvasMouseLeave);
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
                if (this.hoverTrack === null) {
                    return;
                }
                this.$store.dispatch('tracks/seekTrack', this.hoverTrack);
            },
            moveCursorUp() {
                if (this.hoverTrack === null) {
                    this.hoverTrack = 0;
                    return;
                }
                if (this.hoverTrack === 0) {
                    return;
                }
                this.hoverTrack--;
            },
            moveCursorDown() {
                if (this.hoverTrack === null) {
                    this.hoverTrack = 0;
                    return;
                }
                if (this.hoverTrack === this.tracks.length - 1) {
                    return;
                }
                this.hoverTrack++;
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
            }
        }
    };
</script>