<template>
    <div class="mainContainer">
        <aside>
            <logo min="true"></logo>
        </aside>
        <main>
            <navigation></navigation>
            <nuxt/>
        </main>
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
        flex: 1;
        height: 100%;
        background: linear-gradient(to top, rgb(244, 187, 251), rgb(0, 27, 244) 50%, rgb(40, 27, 61) 80%);
    }

    main {
        flex: 4;
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
        middleware: 'authenticated',
        created() {
            this.$store.dispatch('loading/endLoadAfterDelay');
        }
    };
</script>
