<template>
    <div class="pageWrapper">
        <article class="authPage">
            <h2 class="pageHeader">Login Failed</h2>
            <p><nuxt-link to="/auth">Try Again</nuxt-link></p>
        </article>
    </div>
</template>

<style scoped>
    a {
        color: indigo;
    }

    a:hover {
        color: magenta;
    }
</style>

<script>
    export default {
        layout: 'unauthorized'
    };
</script>
