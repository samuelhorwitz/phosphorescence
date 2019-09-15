<template>
    <div>
        <div class="mainContainer">
            <loadingBar></loadingBar>
            <logo v-if="showLogo"></logo>
            <main :class="{noLogo: !showLogo}">
                <nuxt/>
            </main>
            <foot></foot>
        </div>
    </div>
</template>

<style scoped>
    .mainContainer {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }

    body.ios .mainContainer {
        min-height: 84vh;
    }

    @media only screen and (display-mode: fullscreen) {
        body.ios .mainContainer {
            min-height: 100vh;
        }
    }

    main {
        display: flex;
    }

    main.noLogo {
        flex: 1;
        align-items: center;
    }

    footer {
        margin-top: auto;
    }

    @media only screen and (max-height: 449px) {
        main {
            align-items: center;
            justify-content: center;
            margin: 1em 0;
        }
    }

    @media only screen and (max-width: 500px) {
        main.noLogo {
            align-items: flex-start;
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
        computed: {
            showLogo() {
                let showLogo = this.$route.matched.map((r) => {
                    return (r.components.default.options ? r.components.default.options.showLogo : r.components.default.showLogo)
                })[0];
                if (showLogo === false) {
                    return false;
                }
                return true;
            }
        },
        created() {
            this.$store.dispatch('loading/endLoadAfterDelay');
        },
        head() {
            let existingClasses = document.body.getAttribute('class') || '';
            let additionalClasses = ['scrollable', 'vividSunrise'];
            if (additionalClasses.length === 0) {
                return {bodyAttrs: {class: existingClasses}};
            }
            let missingClasses = [];
            for (let additionalClass of additionalClasses) {
                if (!document.body.classList.contains(additionalClass)) {
                    missingClasses.push(additionalClass);
                }
            }
            if (missingClasses.length === 0) {
                return {bodyAttrs: {class: existingClasses}};
            }
            return {
                bodyAttrs: {
                    class: `${existingClasses} ${missingClasses.join(' ')}`
                }
            }
        }
    };
</script>
