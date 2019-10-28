<template>
    <canvas ref="canvas">{{Math.round(percent * 100)}}% complete</canvas>
</template>

<style scoped>
    canvas {
        position: absolute;
    }
</style>

<script>
    const size = 45;
    const center = size / 2;
    const circleSize = 15;
    const fullDegrees = 360;
    const fps = 60;
    const previewLength = 30; // maybe we should pass this in?
    const totalTics = (fps * 2) * previewLength;

    function toMidnightRadians(deg) {
        return (deg - 90) * Math.PI / 180
    }

    export default {
        props: {
            percent: Number,
            shadowColor: String
        },
        data() {
            return {
                ctx: null,
                isDirty: true
            };
        },
        mounted() {
            let canvas = this.$refs.canvas;
            canvas.width = size * devicePixelRatio;
            canvas.height = size * devicePixelRatio;
            canvas.style.width = `${size}px`;
            canvas.style.height = `${size}px`;
            this.ctx = canvas.getContext('2d');
            this.ctx.scale(devicePixelRatio, devicePixelRatio);
            this.beginLoop();
        },
        watch: {
            percent() {
                this.isDirty = true;
            },
            shadowColor() {
                this.isDirty = true;
            }
        },
        methods: {
            async beginLoop() {
                while (this.percent <= 1) {
                    if (this.isDirty) {
                        requestAnimationFrame(this.paint);
                        this.isDirty = false;
                    }
                    if (this.percent === 1) {
                        break;
                    }
                    await new Promise(resolve => setTimeout(resolve, 1000 / fps));
                }
            },
            paint() {
                let ctx = this.ctx;
                let start = toMidnightRadians(0);
                let end = toMidnightRadians(this.percent * fullDegrees);
                ctx.clearRect(0, 0, size, size);
                ctx.save();
                ctx.beginPath();
                ctx.moveTo(center, center);
                ctx.arc(center, center, circleSize, start, end);
                ctx.lineTo(center, center);
                ctx.closePath();
                ctx.fillStyle = 'white';
                if (this.shadowColor) {
                    ctx.shadowColor = this.shadowColor;
                    ctx.shadowBlur = 10;
                }
                ctx.globalCompositeOperation = 'screen';
                ctx.fill();
                ctx.restore();
            }
        }
    };
</script>