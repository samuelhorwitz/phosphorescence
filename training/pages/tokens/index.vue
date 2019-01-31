<template>
    <main>
        <div v-if="!failed">Loading...</div>
        <div v-if="failed">Failed!</div>
    </main>
</template>

<script>
    import {getTokens} from '~/util/spotify';

    export default {
        data() {
            return {
                failed: false
            };
        },
        async created() {
            if (!process.browser) {
                return;
            }
            try {
                let {access, refresh, expires} = await getTokens(this.$route.query.code);
                this.$store.dispatch('session/tokens', {access, refresh, expires});
                this.$router.push('/');
            }
            catch (e) {
                this.failed = true;
                console.error(e);
            }
        }
    };
</script>