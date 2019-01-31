<template>
    <div class="wrapper">
        <div>Transcendental</div>
        <div class="subwrapper">
            <div>Chthonic</div>
            <canvas ref="canvas" :width="retinaSize" :height="retinaSize" :style="canvasStyle" @click="rankTrack"></canvas>
            <div>Aethereal</div>
        </div>
        <div>Primordial</div>
    </div>
</template>

<style scoped>
    canvas {
        border: 1px solid black;
        cursor: none;
    }

    .wrapper {
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    .subwrapper {
        display: flex;
        align-items: center;
    }
</style>

<script>
    export default {
        data() {
            return {
                context: null,
                size: 500,
                cursorX: null,
                cursorY: null,
                dirty: true
            };
        },
        computed: {
            retinaSize() {
                if (!process.browser) {
                    return this.size;
                }
                return this.size * window.devicePixelRatio;
            },
            canvasStyle() {
                return `width: ${this.size}px; height: ${this.size}px;`;
            },
            cursor() {
                let rect = this.$refs.canvas.getBoundingClientRect();
                let x = Math.max(0, Math.min(this.retinaSize, (this.cursorX - rect.left) * window.devicePixelRatio));
                let y = Math.max(0, Math.min(this.retinaSize, (this.cursorY - rect.top) * window.devicePixelRatio));
                let floatX = x / this.retinaSize;
                let floatY = y / this.retinaSize;
                return {x, y, floatX, floatY};
            }
        },
        created() {
            if (!process.browser) {
                return;
            }
            addEventListener('mousemove', this.handleMouse, {passive: true});
        },
        mounted() {
            this.context = this.$refs.canvas.getContext('2d');
            this.drawLoop();
        },
        beforeDestroy() {
            window.removeEventListener('mousemove', this.handleMouse, {passive: true});
        },
        methods: {
            drawLoop() {
                requestAnimationFrame(this.draw);
            },
            draw() {
                this.drawLoop();
                if (!this.dirty) {
                    return;
                }
                this.dirty = false;
                this.context.clearRect(0, 0, this.retinaSize, this.retinaSize);
                this.drawGrid();
                this.drawCursorHover();
            },
            drawCursorHover() {
                if (!this.cursorX || !this.cursorY) {
                    return;
                }
                this.context.save();
                let {x, y} = this.cursor;
                this.context.beginPath();
                this.context.arc(x, y, 20, 0, 2 * Math.PI);
                this.context.fillStyle = 'blue';
                this.context.fill();
                this.context.restore();
            },
            drawGrid() {
                this.context.save();
                this.context.beginPath();
                this.context.moveTo(this.retinaSize / 2, 0);
                this.context.lineTo(this.retinaSize / 2, this.retinaSize);
                this.context.moveTo(0, this.retinaSize / 2);
                this.context.lineTo(this.retinaSize, this.retinaSize / 2);
                this.context.stroke();
                this.context.restore();
            },
            handleMouse(event) {
                this.cursorX = event.clientX;
                this.cursorY = event.clientY;
                this.dirty = true;
            },
            rankTrack() {
                let {floatX, floatY} = this.cursor;
                let currentTrack = this.$store.getters['tracks/currentTrack'];
                let data = {track: currentTrack, analysis: this.$store.getters['tracks/currentTrackAnalysis'], evocativeness: {x: floatX, y: floatY}};
                console.log(data);
                localStorage.setItem(currentTrack.track.id, JSON.stringify(data));
                this.$store.commit('tracks/nextTrack');
            }
        }
    };
</script>