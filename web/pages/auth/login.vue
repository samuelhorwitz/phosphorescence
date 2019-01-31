<template>
</template>

<script>
    import {login} from '~/assets/session';

    export default {
        layout: 'unauthorized',
        async created() {
            this.$store.dispatch('loading/endLoadAfterDelay');
            try {
                await login(this.$route.query.code);
                if (parent != window) {
                    window.close();
                }
                else {
                    this.$router.push('/');
                }
            }
            catch (e) {
                console.error(e);
                this.$router.push('/auth/failed');
            }
        }
    };
</script>