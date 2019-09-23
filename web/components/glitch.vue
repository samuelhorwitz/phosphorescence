<template>
    <span>
        <span class="t" v-if="rawText" :class="{glitch: enable}" :data-text="rawGlyph" v-html="rawText" ref="raw"></span>
        <span class="t" v-if="!rawText" :class="{glitch: enable}" :data-text="text">{{text}}</span>
    </span>
</template>

<script>
    export default {
        props: {text: String, rawText: String, enable: Boolean},
        data() {
            return {
                rawGlyph: ''
            };
        },
        watch: {
            rawText: {
                immediate: true,
                handler() {
                    this.$nextTick(() => {
                        if (!this.$refs.raw) {
                            return;
                        }
                        this.rawGlyph = this.$refs.raw.innerText;
                    });
                }
            }
        }
    };
</script>

<style scoped>
    span.t {
        color: inherit;
        white-space: nowrap;
    }

    @keyframes noise-anim {
        0% {
            clip: rect(32px, 9999px, 9px, 0);
        }
        5% {
            clip: rect(77px, 9999px, 66px, 0);
        }
        10% {
            clip: rect(33px, 9999px, 20px, 0);
        }
        15% {
            clip: rect(42px, 9999px, 87px, 0);
        }
        20% {
            clip: rect(56px, 9999px, 7px, 0);
        }
        25% {
            clip: rect(92px, 9999px, 28px, 0);
        }
        30% {
            clip: rect(66px, 9999px, 34px, 0);
        }
        35% {
            clip: rect(57px, 9999px, 86px, 0);
        }
        40% {
            clip: rect(42px, 9999px, 25px, 0);
        }
        45% {
            clip: rect(50px, 9999px, 9px, 0);
        }
        50% {
            clip: rect(96px, 9999px, 12px, 0);
        }
        55% {
            clip: rect(72px, 9999px, 96px, 0);
        }
        60% {
            clip: rect(93px, 9999px, 43px, 0);
        }
        65% {
            clip: rect(66px, 9999px, 30px, 0);
        }
        70% {
            clip: rect(87px, 9999px, 33px, 0);
        }
        75% {
            clip: rect(24px, 9999px, 82px, 0);
        }
        80% {
            clip: rect(29px, 9999px, 19px, 0);
        }
        85% {
            clip: rect(44px, 9999px, 96px, 0);
        }
        90% {
            clip: rect(23px, 9999px, 37px, 0);
        }
        95% {
            clip: rect(40px, 9999px, 1px, 0);
        }
        100% {
            clip: rect(30px, 9999px, 41px, 0);
        }
    }
    .glitch:after {
        content: attr(data-text);
        position: absolute;
        left: 2px;
        text-shadow: -1px 0 aqua;
        top: 0;
        color: white;
        overflow: hidden;
        clip: rect(0, 900px, 0, 0);
        animation: noise-anim 2s infinite linear alternate-reverse;
    }

    @keyframes noise-anim-2 {
        0% {
            clip: rect(44px, 9999px, 98px, 0);
        }
        5% {
            clip: rect(99px, 9999px, 79px, 0);
        }
        10% {
            clip: rect(37px, 9999px, 91px, 0);
        }
        15% {
            clip: rect(88px, 9999px, 34px, 0);
        }
        20% {
            clip: rect(98px, 9999px, 86px, 0);
        }
        25% {
            clip: rect(2px, 9999px, 2px, 0);
        }
        30% {
            clip: rect(72px, 9999px, 80px, 0);
        }
        35% {
            clip: rect(39px, 9999px, 1px, 0);
        }
        40% {
            clip: rect(51px, 9999px, 62px, 0);
        }
        45% {
            clip: rect(57px, 9999px, 73px, 0);
        }
        50% {
            clip: rect(38px, 9999px, 77px, 0);
        }
        55% {
            clip: rect(11px, 9999px, 33px, 0);
        }
        60% {
            clip: rect(4px, 9999px, 84px, 0);
        }
        65% {
            clip: rect(63px, 9999px, 82px, 0);
        }
        70% {
            clip: rect(83px, 9999px, 65px, 0);
        }
        75% {
            clip: rect(66px, 9999px, 47px, 0);
        }
        80% {
            clip: rect(97px, 9999px, 17px, 0);
        }
        85% {
            clip: rect(90px, 9999px, 80px, 0);
        }
        90% {
            clip: rect(9px, 9999px, 100px, 0);
        }
        95% {
            clip: rect(59px, 9999px, 55px, 0);
        }
        100% {
            clip: rect(49px, 9999px, 70px, 0);
        }
    }
    .glitch:before {
        content: attr(data-text);
        position: absolute;
        left: -1px;
        text-shadow: 2px 0 magenta;
        top: 0;
        color: white;
        overflow: hidden;
        clip: rect(0, 900px, 0, 0);
        animation: noise-anim-2 3s infinite linear alternate-reverse;
    }
</style>