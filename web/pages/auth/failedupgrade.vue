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
        layout: 'empty',
        head() {
            return {
                title: 'Session Upgrade Failed | Phosphorescence',
                meta: [
                    {
                        hid: 'og:type',
                        name: 'og:type',
                        content: 'website'
                    },
                    {
                        hid: 'og:site_name',
                        name: 'og:site_name',
                        content: 'Phosphorescence'
                    },
                    {
                        hid: 'og:image',
                        name: 'og:image',
                        content: 'https://phosphor.me/og.png'
                    },
                    {
                        hid: 'og:description',
                        name: 'og:description',
                        content: 'Build trance and chill-out playlists for Spotify'
                    },
                    {
                        hid: 'og:url',
                        name: 'og:url',
                        content: 'https://phosphor.me/auth/failedupgrade'
                    },
                    {
                        hid: 'og:title',
                        name: 'og:title',
                        content: 'Session Upgrade Failed'
                    }
                ]
            };
        }
    };
</script>
