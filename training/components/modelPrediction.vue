<template>
    <div class="wrapper">
        <div>Transcendental</div>
        <div class="subwrapper">
            <div>Chthonic</div>
            <canvas ref="canvas" :width="retinaSize" :height="retinaSize" :style="canvasStyle"></canvas>
            <div>Aethereal</div>
        </div>
        <div>Primordial</div>
    </div>
</template>

<style scoped>
    canvas {
        border: 1px solid black;
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
    import {predict} from '~/util/spotify';

    export default {
        data() {
            return {
                aetherealness: null,
                primordialness: null,
                context: null,
                size: 500,
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
            }
        },
        mounted() {
            this.context = this.$refs.canvas.getContext('2d');
            this.drawLoop();
            this.$store.watch(() => this.$store.getters['tracks/currentTrackAnalysis'], async a => {
                if (a) {
                    let {aetherealness, primordialness} = await predict(a);
                    this.aetherealness = aetherealness;
                    this.primordialness = primordialness;
                    this.dirty = true;
                }
            });
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
                this.fillBackground();
                this.drawGrid();
                this.drawDot();
            },
            fillBackground() {
                this.context.save();
                this.context.fillStyle = 'black';
                this.context.fillRect(0, 0, this.retinaSize, this.retinaSize);
                this.context.restore();
            },
            drawGrid() {
                this.context.save();
                this.context.strokeStyle = 'white';
                this.context.beginPath();
                this.context.moveTo(this.retinaSize / 2, 0);
                this.context.lineTo(this.retinaSize / 2, this.retinaSize);
                this.context.moveTo(0, this.retinaSize / 2);
                this.context.lineTo(this.retinaSize, this.retinaSize / 2);
                this.context.stroke();
                this.context.restore();
            },
            drawDot() {
                if (!this.aetherealness || !this.primordialness) {
                    return;
                }
                this.context.save();
                this.context.beginPath();
                this.context.arc(this.aetherealness * this.retinaSize, this.primordialness * this.retinaSize, 20, 0, 2 * Math.PI);
                this.context.fillStyle = 'pink';
                this.context.fill();
                this.context.restore();
            }
        }
    };
</script>