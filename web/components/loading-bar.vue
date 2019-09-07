<template>
    <progress v-show="loadingProgressSticky" :class="{hidden: hide}" max="100" :value="loadingProgressSticky">{{loadingProgressSticky}}%</progress>
</template>

<style scoped>
    progress {
        width: 100%;
        position: absolute;
        top: 0px;
        height: 2px;
        appearance: none;
        border: 0px;
        background: none;
        color: cyan;
        z-index: 999999;
        transition: opacity 1s linear;
    }

    progress.hidden {
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
    export default {
        data() {
            return {
                loadingProgressSticky: 0,
                hide: false,
                interval: null
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
                        this.loadingProgressSticky += 0.1;
                    }, 10);
                }
            }
        }
    };
</script>
