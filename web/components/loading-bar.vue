<template>
    <div>
        <progress v-if="!notchMode" v-show="loadingProgressSticky" :class="{hidden: hide}" max="100" :value="loadingProgressSticky">{{loadingProgressSticky}}%</progress>
        <canvas v-if="notchMode" v-show="loadingProgressSticky" ref="canvas" :class="{hidden: hide}">{{loadingProgressSticky}}%</canvas>
    </div>
</template>

<style scoped>
    div {
        position: absolute;
    }

    progress {
        width: 100%;
        position: fixed;
        top: env(safe-area-inset-top, 0px);
        height: 2px;
        appearance: none;
        border: 0px;
        background: none;
        color: cyan;
        transition: opacity 1s linear;
        z-index: 999999;
    }

    progress.hidden {
        opacity: 0;
    }

    canvas {
        position: fixed;
        top: 0px;
        transition: opacity 1s linear;
        z-index: 999999;
    }

    canvas.hidden {
        opacity: 0;
    }

    progress::-webkit-progress-bar {
        background: none;
    }

    progress::-webkit-progress-value {
        background-color: cyan;
    }

    progress::-moz-progress-bar {
        background-color: cyan;
    }
</style>

<script>
    const iOSNotchHeight = 30;

    function buildNotchPath(ctx) {
        ctx.arc(40, 40, 39, Math.PI, Math.PI * 1.5);
        ctx.lineTo(76, 0);
        ctx.arc(76, 6, 6, Math.PI * 1.5, 0);
        ctx.arc(103, 10, 20, Math.PI, Math.PI * 0.5, true);
        ctx.lineTo(272, 30);
        ctx.arc(272, 10, 20, Math.PI * 0.5, 0, true);
        ctx.arc(299, 6, 6, Math.PI, Math.PI * 1.5);
        ctx.lineTo(299, 0);
        ctx.arc(335, 40, 39, Math.PI * 1.5, 0);
    }

    export default {
        data() {
            return {
                loadingProgressSticky: 0,
                hide: false,
                interval: null,
                notchMode: false
            };
        },
        computed: {
            loadingProgress() {
                return this.$store.getters['loading/progress'];
            }
        },
        watch: {
            loadingProgress(newVal) {
                this.interval && clearInterval(this.interval);
                if (!newVal) {
                    this.loadingProgressSticky = 100;
                    this.hide = true;
                    setTimeout(() => {
                        this.loadingProgressSticky = 0;
                        this.hide = false;
                    }, 1000);
                } else if (newVal > this.loadingProgressSticky) {
                    this.hide = false;
                    this.loadingProgressSticky = newVal;
                    this.interval = setInterval(() => {
                        if (this.loadingProgressSticky >= 95) {
                            clearInterval(this.interval);
                            return;
                        }
                        this.loadingProgressSticky += 0.05;
                    }, 10);
                }
            },
            loadingProgressSticky() {
                this.redrawCanvas();
            }
        },
        mounted() {
            addEventListener('orientationchange', this.resetNotchMode);
            this.resetNotchMode();
        },
        beforeDestroy() {
            removeEventListener('orientationchange', this.resetNotchMode);
        },
        methods: {
            resetNotchMode() {
                let oldNotchMode = this.notchMode;
                if (this.isNotchedFullscreenIphonePortrait()) {
                    this.notchMode = true;
                } else {
                    this.notchMode = false;
                }
                if (this.notchMode && !oldNotchMode) {
                    this.$nextTick(() => {
                        this.initializeCanvas();
                    });
                }
            },
            isNotchedFullscreenIphonePortrait() {
                return /\b(iPhone)\b/.test(navigator.userAgent) && !orientation && navigator.standalone && screen.width * devicePixelRatio === 1125 && screen.height * devicePixelRatio === 2436;
            },
            initializeCanvas() {
                if (!this.notchMode) {
                    return;
                }
                let canvas = this.$refs.canvas;
                let ctx = canvas.getContext('2d');
                let width = innerWidth;
                let height = iOSNotchHeight + 4;
                if (navigator.standalone) {
                    width = outerWidth;
                }
                canvas.width = width * devicePixelRatio;
                canvas.height = height * devicePixelRatio;
                canvas.style.width = `${width}px`;
                canvas.style.height = `${height}px`;
                ctx.scale(devicePixelRatio, devicePixelRatio);
                this.redrawCanvas();
            },
            redrawCanvas() {
                if (!this.notchMode) {
                    return;
                }
                requestAnimationFrame(() => {
                    let canvas = this.$refs.canvas;
                    let ctx = canvas.getContext('2d');
                    let width = innerWidth;
                    let height = iOSNotchHeight + 4;
                    if (navigator.standalone) {
                        width = outerWidth;
                    }
                    ctx.save();
                    ctx.clearRect(0, 0, width, height);
                    ctx.translate(0, 1);
                    ctx.beginPath();
                    buildNotchPath(ctx);
                    ctx.strokeStyle = 'cyan';
                    ctx.lineWidth = 2;
                    ctx.stroke();
                    ctx.globalCompositeOperation = 'destination-in';
                    ctx.fillStyle = 'white';
                    ctx.translate(0, -1);
                    ctx.fillRect(0, 0, (this.loadingProgressSticky / 100) * width, height);
                    ctx.restore();
                });
            }
        }
    };
</script>
