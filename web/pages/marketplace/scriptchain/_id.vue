<template>
    <section>
        <h3>
            {{scriptChain.name}}&nbsp;&horbar;&nbsp;{{scriptChain.authorName}}<spotifyUserLink :id="scriptChain.authorSpotifyId" :name="scriptChain.authorName" :isAuthor="true"/>
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
                return error({statusCode: 400, message: 'No script chain id'});
            }
            let scriptChainResponse = await fetch(`${process.env.API_ORIGIN}/script-chains/${id}`, {credentials: 'include'});
            if (!scriptChainResponse.ok) {
                return error({statusCode: scriptChainResponse.status, message: 'Could not get script chain'});
            }
            let {scriptChain} = await scriptChainResponse.json();
            return {
                scriptChain,
                description: getSafeHtml(scriptChain.description, buildTagMarker(scriptChain.description))
            };
        },
        watchQuery: ['id'],
        methods: {
            handleClicks
        }
    };
</script>