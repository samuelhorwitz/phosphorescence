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
    const iPhoneXNotchHeight = 30;
    const iPhoneXMaxNotchHeight = 30;
    const iPhoneXRNotchHeight = 32;

    function buildNotchPathIphoneX(ctx) {
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

    function buildNotchPathIphoneXSMax(ctx) {
        ctx.arc(40, 40, 39, Math.PI, Math.PI * 1.5);
        ctx.lineTo(96, 0);
        ctx.arc(96, 6, 6, Math.PI * 1.5, 0);
        ctx.arc(122, 10, 20, Math.PI, Math.PI * 0.5, true);
        ctx.lineTo(292, 30);
        ctx.arc(292, 10, 20, Math.PI * 0.5, 0, true);
        ctx.arc(319, 6, 6, Math.PI, Math.PI * 1.5);
        ctx.lineTo(319, 0);
        ctx.arc(374, 40, 39, Math.PI * 1.5, 0);
    }

    function buildNotchPathIphoneXR(ctx) {
        ctx.arc(40, 40, 39, Math.PI, Math.PI * 1.5);
        ctx.lineTo(85, 0);
        ctx.arc(85, 6, 6, Math.PI * 1.5, 0);
        ctx.arc(112, 12, 20, Math.PI, Math.PI * 0.5, true);
        ctx.lineTo(302, 32);
        ctx.arc(302, 12, 20, Math.PI * 0.5, 0, true);
        ctx.arc(330, 6, 6, Math.PI, Math.PI * 1.5);
        ctx.lineTo(330, 0);
        ctx.arc(374, 40, 39, Math.PI * 1.5, 0);
    }

    export default {
        data() {
            return {
                loadingProgressSticky: 0,
                hide: false,
                interval: null,
                notchMode: false,
                notchStyle: null
            };
        },
        computed: {
            loadingProgress() {
                return this.$store.getters['loading/progress'];
            },
            notchHeight() {
                if (this.notchStyle === 'X') {
                    return iPhoneXNotchHeight;
                } else if (this.notchStyle === 'Max') {
                    return iPhoneXMaxNotchHeight;
                } else if (this.notchStyle === 'R') {
                    return iPhoneXRNotchHeight;
                } else {
                    return null;
                }
            }
        },
        watch: {
            loadingProgress(newVal) {
                if (!newVal) {
                    this.finish();
                } else if (newVal > this.loadingProgressSticky) {
                    this.tick(newVal);
                } else {
                    this.clearSmoothLoadInterval();
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
            start() {
                this.tick(0.001);
            },
            finish() {
                this.clearSmoothLoadInterval();
                this.loadingProgressSticky = 100;
                this.hide = true;
                setTimeout(() => {
                    this.loadingProgressSticky = 0;
                    this.hide = false;
                }, 1000);
            },
            fail() {
                this.finish();
            },
            increase(newVal) {
                this.loadingProgressSticky = newVal;
            },
            tick(newVal) {
                this.clearSmoothLoadInterval();
                this.hide = false;
                this.loadingProgressSticky = newVal;
                this.interval = setInterval(() => {
                    if (this.loadingProgressSticky >= 95) {
                        clearInterval(this.interval);
                        return;
                    }
                    this.loadingProgressSticky += 0.05;
                }, 10);
            },
            clearSmoothLoadInterval() {
                this.interval && clearInterval(this.interval);
            },
            resetNotchMode() {
                let oldNotchMode = this.notchMode;
                if (this.isIphoneFullscreenPortrait()) {
                    if (this.isIphoneX()) {
                        this.notchMode = true;
                        this.notchStyle = 'X';
                    } else if (this.isIphoneXSMax()) {
                        this.notchMode = true;
                        this.notchStyle = 'Max';
                    } else if (this.isIphoneXR()) {
                        this.notchMode = true;
                        this.notchStyle = 'R';
                    } else {
                        this.notchMode = false;
                        this.notchStyle = null;
                    }
                } else {
                    this.notchMode = false;
                    this.notchStyle = null;
                }
                if (this.notchMode && !oldNotchMode) {
                    this.$nextTick(() => {
                        this.initializeCanvas();
                    });
                }
            },
            isIphoneFullscreenPortrait() {
                return /\b(iPhone)\b/.test(navigator.userAgent) && !orientation && navigator.standalone;
            },
            isIphoneX() {
                return screen.width === 375 && screen.height === 812 && devicePixelRatio === 3;
            },
            isIphoneXSMax() {
                return screen.width === 414 && screen.height === 896 && devicePixelRatio === 3;
            },
            isIphoneXR() {
                return screen.width === 414 && screen.height === 896 && devicePixelRatio === 2;
            },
            initializeCanvas() {
                if (!this.notchMode || !this.notchStyle) {
                    return;
                }
                let canvas = this.$refs.canvas;
                let ctx = canvas.getContext('2d');
                let width = innerWidth;
                let height = this.notchHeight + 4;
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
                if (!this.notchMode || !this.notchStyle) {
                    return;
                }
                requestAnimationFrame(() => {
                    let canvas = this.$refs.canvas;
                    let ctx = canvas.getContext('2d');
                    let width = innerWidth;
                    let height = this.notchHeight + 4;
                    if (navigator.standalone) {
                        width = outerWidth;
                    }
                    ctx.save();
                    ctx.clearRect(0, 0, width, height);
                    ctx.translate(0, 1);
                    ctx.beginPath();
                    if (this.notchStyle === 'X') {
                        buildNotchPathIphoneX(ctx);
                    } else if (this.notchStyle === 'Max') {
                        buildNotchPathIphoneXSMax(ctx);
                    } else if (this.notchStyle === 'R') {
                        buildNotchPathIphoneXR(ctx);
                    }
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
