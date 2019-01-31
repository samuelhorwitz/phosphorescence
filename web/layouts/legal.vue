<template>
    <div>
        <logo></logo>
        <main>
            <nuxt/>
        </main>
        <foot></foot>
    </div>
</template>

<script>
    import logo from '~/components/logo';
    import foot from '~/components/foot';

    export default {
        components: {
            logo,
            foot
        },
        created() {
            this.$store.dispatch('loading/endLoadAfterDelay');
        }
    };
</script>
