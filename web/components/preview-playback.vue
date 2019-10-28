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

    function toMidnightRadians(deg) {
        return (deg - 90) * Math.PI / 180
    }

    export default {
        props: {percent: Number},
        data() {
            return {
                ctx: null
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
        },
        watch: {
            percent() {
                requestAnimationFrame(() => {
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
                    ctx.shadowColor = 'cyan';
                    ctx.shadowBlur = 10;
                    ctx.globalCompositeOperation = 'screen';
                    ctx.fill();
                    ctx.restore();
                });
            }
        }
    };
</script>