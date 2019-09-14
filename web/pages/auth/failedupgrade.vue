<template>
    <div class="pageWrapper">
        <article class="authPage">
            <h2 class="pageHeader">Authentication Failed</h2>
            <p>It appears your magic link didn't work. Perhaps it's expired?</p>
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
        layout: 'empty'
    };
</script>
