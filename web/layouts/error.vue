<template>
    <div class="container">
        <h2 v-if="notFound">Page not found</h2>
        <nuxt-link v-if="notFound" to="/" class="home">Return home</nuxt-link>
        <h1 v-if="!notFound">An error occurred</h1>
        <p v-if="!notFound">{{error.message}}</p>
        <p v-if="!$store.state.user.user">It is possible we think you are a robot. If you are not a bot, you can try to refresh, however you may have to <nuxt-link to="/auth">log in</nuxt-link> with your Spotify account rather than browse anonymously.</p>
    </div>
</template>

<style scoped>
    .container {
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;
        width: 100%;
        height: 100%;
        color: white;
        flex: 1;
    }

    h2 {
        margin: 0px;
        font-size: 5em;
    }

    a {
        color: aqua;
        text-decoration: none;
    }

    a.home {
        margin-top: 2em;
        font-size: 1.75em;
    }

    a:hover {
        text-decoration: underline;
    }

    p {
        margin: 0px;
        margin-top: 2em;
        max-width: 30em;
        font-size: 1.75em;
        text-align: center;
    }
</style>

<script>
    export default {
        props: ['error'],
        computed: {
            notFound() {
                return this.error && this.error.statusCode === 404;
            }
        },
        beforeCreate() {
            this.$nuxt.setLayout('empty'); // bug that we cannot just use layout property
        },
        created() {
            this.$store.dispatch('loading/endLoadAfterDelay');
        }
    };
</script>
