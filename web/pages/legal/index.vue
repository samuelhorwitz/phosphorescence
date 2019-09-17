<template>
    <legal>
        <h2>Legal</h2>
        <ul>
            <li><nuxt-link to="/legal/tos">Terms of Service</nuxt-link></li>
            <li><nuxt-link to="/legal/privacy">Privacy Policy</nuxt-link></li>
            <li><nuxt-link to="/legal/cookies">Cookie Policy</nuxt-link></li>
            <li><nuxt-link to="/legal/trademarks">Trademark Policy</nuxt-link></li>
        </ul>
    </legal>
</template>

<script>
    import legal from '~/components/legal/legal';

    export default {
        layout: 'empty',
        components: {legal}
    };
</script>