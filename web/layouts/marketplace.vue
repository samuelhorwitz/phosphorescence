<template>
    <div class="mainContainer">
        <aside>
            <logo min="true"></logo>
        </aside>
        <main>
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
    }
</style>

<script>
    import logo from '~/components/logo';

    export default {
        components: {
            logo
        },
        middleware: 'authenticated',
        created() {
            this.$store.dispatch('loading/endLoadAfterDelay');
        }
    };
</script>
