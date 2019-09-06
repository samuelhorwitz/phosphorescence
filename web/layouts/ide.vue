<template>
    <div class="ideWrapper">
        <sidebar class="sidebar"></sidebar>
        <main>
            <nuxt/>
        </main>
        <div class="playerWrapper">
            <details-pane></details-pane>
        </div>
    </div>
</template>

<style scoped>
    .ideWrapper {
        border: 7px inset magenta;
        position: absolute;
        top: 0px;
        left: 0px;
        right: 0px;
        bottom: 0px;
        display: flex;
    }
    .sidebar {
        flex: 2.5;
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        align-items: center;
    }
    main {
        flex: 8;
        display: flex;
    }

    .playerWrapper {
        flex: 4;
        display: flex;
        max-width: 25em;
        overflow: hidden;
    }
</style>

<script>
    import sidebar from '~/components/ide/sidebar';
    import detailsPane from '~/components/ide/detailspane';

    export default {
        components: {
            sidebar,
            detailsPane
        },
        middleware: ['authenticated', 'disabled'],
        created() {
            this.$store.commit('ide/restore');
        }
    };
</script>
