<template>
    <nav>
        <ul>
            <li>
                <search></search>
            </li>
        </ul>
    </nav>
</template>

<style scoped>
    nav {
        height: 5em;
        padding: 1em;
        box-sizing: border-box;
        display: flex;
    }

    ul {
        list-style: none;
        display: flex;
        margin: 0px;
        padding: 0px;
        flex: 1;
    }

    li {
        display: flex;
        margin: 0px;
        padding: 0px;
        flex: 1;
    }
</style>

<script>
    import search from '~/components/marketplace/search';

    export default {
        components: {
            search
        }
    };
</script>
