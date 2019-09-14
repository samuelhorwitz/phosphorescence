<template>
    <article class="legal">
        <slot></slot>
    </article>
</template>

<style scoped>
    .legal {
        margin: 1em;
        font-size: 16px;
        background-color: teal;
        padding: 1em;
        flex: 1;
        overflow-y: scroll;
        -webkit-overflow-scrolling: touch;
        border: 5px outset magenta;
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