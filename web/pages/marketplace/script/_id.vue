<template>
    <section>
        <h3>
            {{script.name}}&nbsp;&horbar;&nbsp;{{script.authorName}}<spotifyUserLink :id="script.authorSpotifyId" :name="script.authorName" :isAuthor="true"/>
        </h3>
        <p v-html="description" @click="handleClicks"></p>
    </section>
</template>

<style scoped>
    h3 {
        margin: 0px;
    }

    p {
        margin: 0px;
        margin-top: 1em;
    }

    img.spotifyLink {
        width: 1em;
        height: 1em;
        min-width: 21px;
        min-height: 21px;
        transform: translateX(50%);
        vertical-align: middle;
    }
</style>

<script>
    import {getSafeHtml, buildTagMarker, handleClicks} from '~/assets/safehtml';
    import spotifyUserLink from '~/components/marketplace/spotifyuserlink';

    export default {
        layout: 'marketplace',
        components: {
            spotifyUserLink
        },
        async asyncData({params, error}) {
            let {id} = params;
            if (!id) {
                return error({statusCode: 400, message: 'No script id'});
            }
            let scriptResponse = await fetch(`${process.env.API_ORIGIN}/scripts/${id}`, {credentials: 'include'});
            if (!scriptResponse.ok) {
                return error({statusCode: scriptResponse.status, message: 'Could not get script'});
            }
            let {script} = await scriptResponse.json();
            return {
                script,
                description: getSafeHtml(script.description, buildTagMarker(script.description))
            };
        },
        watchQuery: ['id'],
        methods: {
            handleClicks
        }
    };
</script>