<template>
    <div class="mainContainer">
        <aside>
            <logo min="true"></logo>
        </aside>
        <div class="mainWrapper">
            <navigation></navigation>
            <main>
                <nuxt/>
            </main>
        </div>
    </div>
</template>

<style scoped>
    .mainContainer {
        background-color: white;
        display: flex;
        width: 100vw;
        height: 100vh;
    }

    aside {
        display: flex;
        align-items: flex-start;
        width: 20em;
        height: 100%;
        background: linear-gradient(to top, rgb(244, 187, 251), rgb(0, 27, 244) 50%, rgb(40, 27, 61) 80%);
    }

    main {
        flex: 1;
        margin: 0px 2em;
    }

    .mainWrapper {
        display: flex;
        flex-direction: column;
        flex: 1;
        overflow-y: auto;
    }
</style>

<script>
    import logo from '~/components/logo';
    import navigation from '~/components/marketplace/navigation';

    export default {
        components: {
            logo,
            navigation
        },
        middleware: ['authenticated', 'disabled'],
        created() {
            this.$store.dispatch('loading/endLoadAfterDelay');
        }
    };
</script>
