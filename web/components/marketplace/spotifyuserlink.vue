<template>
    <a :href="'https://open.spotify.com/user/' + id" :rel="'external noopener' + (isAuthor ? ' author' : '')" target="_blank"><img :alt="'Spotify Profile for ' + name" :title="'Spotify Profile for ' + name" class="spotifyLink" src="/images/spotify_green.png"></a>
</template>

<style scoped>
    img.spotifyLink {
        display: inline;
        width: 1em;
        height: 1em;
        min-width: 21px;
        min-height: 21px;
        transform: translateX(50%);
        vertical-align: middle;
    }
</style>

<script>
    export default {
        props: {id: String, name: String, isAuthor: Boolean}
    };
</script>
