<template>
    <div class="mainContainer">
        <logo></logo>
        <toolbar></toolbar>
        <main>
            <nuxt/>
        </main>
        <album></album>
        <player></player>
        <foot></foot>
    </div>
</template>

<style scoped>
    .mainContainer {
        display: grid;
        grid-template-columns: 10fr 1fr;
        grid-template-rows: min-content min-content minmax(100px, 1fr) min-content min-content;
        height: 100vh;
    }

    main {
        grid-column: 1 / 2;
    }

    @media only screen and (min-height: 450px) and (max-width: 1099px) {
        .mainContainer {
            grid-template-columns: 100%;
            grid-template-rows: min-content min-content minmax(100px, 1fr) 100px min-content min-content;
        }

        main {
            grid-column: 1 / 2;
        }
    }

    @media only screen and (max-height: 449px) {
        .mainContainer {
            grid-template-columns: 28vw minmax(40vw, 10fr) min-content;
            grid-template-rows: max-content minmax(100px, 1fr) max-content max-content;
        }

        main {
            grid-column: 1 / 3;
            grid-row: 2 / 3;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 1em 0;
        }
    }

    @media only screen and (max-width: 949px) and (min-height: 275px) and (max-height: 449px) {
        main {
            display: none;
        }
    }

    @media only screen and (max-height: 274px) {
        main {
            display: none;
        }
    }

    @media only screen and (max-width: 500px) and (max-height: 249px) {
        .mainContainer {
            grid-template-columns: 0 minmax(40vw, 10fr) min-content;
        }
    }
</style>

<script>
    import {initialize, builders, loadNewPlaylist} from '~/assets/recordcrate';
    import {accessTokenExists, refreshUser, getUsersCountry} from '~/assets/session';
    import {spider} from '~/assets/spotify';
    import logo from '~/components/logo';
    import toolbar from '~/components/toolbar';
    import album from '~/components/album';
    import player from '~/components/player';
    import foot from '~/components/foot';

    export default {
        components: {
            logo,
            toolbar,
            album,
            player,
            foot
        },
        middleware: 'authenticated',
        beforeCreate() {
            if (!accessTokenExists()) {
                refreshUser();
            }
        },
        async created() {
            this.$store.commit('loading/startLoad');
            this.$store.commit('preferences/restore');
            this.$store.commit('tracks/restore');
            await initialize(await getUsersCountry());
            if (!this.$store.getters['tracks/playlistLoaded']) {
                let {playlist} = await loadNewPlaylist(this.$store.state.preferences.tracksPerPlaylist, builders.randomwalk, builders[this.$store.state.preferences.seedStyle]);
                this.$store.dispatch('tracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
            }
            this.$store.dispatch('loading/endLoadAfterDelay');
        }
    };
</script>
