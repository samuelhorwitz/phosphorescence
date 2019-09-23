<template>
    <legal>
        <h2>Legal</h2>
        <ul>
            <li><nuxt-link to="/legal/tos">Terms of Service</nuxt-link></li>
            <li><nuxt-link to="/legal/privacy">Privacy Policy</nuxt-link></li>
            <li><nuxt-link to="/legal/cookies">Cookie Policy</nuxt-link></li>
            <li><nuxt-link to="/legal/trademarks">Trademark Policy</nuxt-link></li>
        </ul>
        <p>
            This site is protected by reCAPTCHA and the Google
                <a target="_blank" rel="external noopener" href="https://policies.google.com/privacy">Privacy Policy</a> and
                <a target="_blank" rel="external noopener" href="https://policies.google.com/terms">Terms of Service</a> apply.
        </p>
    </legal>
</template>

<script>
    import legal from '~/components/legal/legal';

    export default {
        layout: 'empty',
        components: {legal}
    };
</script>