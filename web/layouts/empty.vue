<template>
    <div>
        <div class="mainContainer">
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
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }

    body.ios .mainContainer {
        min-height: 84vh;
    }

    main {
        display: flex;
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
        },
    };
</script>
