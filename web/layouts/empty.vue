<template>
    <div>
        <div class="mainContainer" ref="mainContainer">
            <loadingBar></loadingBar>
            <logo></logo>
            <main>
                <nuxt/>
            </main>
            <foot></foot>
        </div>
    </div>
</template>

<style scoped>
    .mainContainer {
        display: grid;
        grid-template-columns: 10fr 1fr;
        height: 100vh;
        grid-template-rows: min-content minmax(100px, 1fr) min-content;
    }

    @media only screen and (display-mode: fullscreen) and (orientation: portrait) {
        .mainContainer {
            height: calc(100vh - 70px);
        }
    }

    @media only screen and (display-mode: fullscreen) and (orientation: landscape) {
        .mainContainer {
            height: calc(100vh - 20px);
        }
    }

    main {
        grid-column: 1 / 3;
        display: flex;
    }

    @media only screen and (min-height: 450px) and (max-width: 1099px) {
        .mainContainer {
            grid-template-columns: 100%;
            grid-template-rows: min-content minmax(100px, 1fr) min-content;
        }
    }

    @media only screen and (max-height: 449px) {
        .mainContainer {
            grid-template-columns: 28vw minmax(40vw, 10fr) min-content;
            grid-template-rows: max-content minmax(100px, 1fr) max-content;
        }

        main {
            grid-row: 2 / 3;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 1em 0;
            grid-column: 1 / 3;
        }
    }

    @media only screen and (max-height: 449px) and (max-width: 1199px) {
        .mainContainer {
            grid-template-columns: 28vw minmax(30vw, 10fr) min-content;
        }
    }

    @media only screen and (max-width: 949px) and (min-height: 275px) and (max-height: 449px) {
        .mainContainer {
            grid-template-rows: max-content minmax(100px, 1fr);
        }
    }

    @media only screen and (max-height: 274px) {
        .mainContainer {
            grid-template-rows: max-content minmax(100px, 1fr);
        }
    }

    @media only screen and (max-width: 500px) and (max-height: 249px) {
        .mainContainer {
            grid-template-columns: 0 minmax(40vw, 10fr) min-content;
        }
    }
</style>

<script>
    import logo from '~/components/logo';
    import foot from '~/components/foot';
    import loadingBar from '~/components/loading-bar';

    export default {
        components: {
            logo,
            foot,
            loadingBar
        },
        created() {
            this.$store.dispatch('loading/endLoadAfterDelay');
        }
    };
</script>
