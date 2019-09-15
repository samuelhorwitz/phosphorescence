<template>
    <article class="legal">
        <div class="container">
            <slot></slot>
        </div>
    </article>
</template>

<style scoped>
    .legal {
        margin: 1em;
        font-size: 16px;
        background-color: teal;
        padding: 1em;
        flex: 1;
        border: 5px outset magenta;
        display: flex;
        justify-content: center;
    }

    .container {
        max-width: 40em;
    }

    a {
        color: indigo;
    }

    a:hover {
        color: magenta;
    }
</style>

<script>
    export default {};
</script>