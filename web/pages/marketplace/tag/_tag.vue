<template>
    <section>
        <h2>#{{$route.params.tag}}</h2>
        <ol v-if="tags">
            <li v-for="(tag, index) of tags" tabindex="0">
                <h3>
                    <router-link :to="'/marketplace/' + builderTypePathFragment[tag.resultType] + '/' + encodeURIComponent(tag.id)"><span class="name">{{tag.name}}</span></router-link>&nbsp;&horbar;&nbsp;<span class="authorName">{{tag.authorName}}</span>
                </h3>
                <p v-html="descriptions[index]" @click="handleClicks"></p>
            </li>
        </ol>
        <footer v-if="tags">
            End of results
        </footer>
        <aside v-if="!tags">
            No results for tag #{{$route.params.tag}}
        </aside>
    </section>
</template>

<style scoped>
    ol {
        list-style: none;
        margin: 0px;
    }

    li {
        margin-bottom: 2em;
    }

    h2 {
        margin: 0px;
        margin-bottom: 1em;
        font-weight: bolder;
        font-variant: all-small-caps;
    }

    h3 {
        margin: 0px;
        margin-bottom: 1em;
    }

    p {
        margin: 0px;
        margin-left: 1em;
    }

    aside {
        font-size: 3em;
    }

    footer {
        font-style: italic;
        font-size: 1.25em;
        margin-bottom: 2em;
    }
</style>

<script>
    import {getSafeHtml, buildTagMarker, handleClicks} from '~/assets/safehtml';

    export default {
        layout: 'marketplace',
        async asyncData({params, error}) {
            let {tag} = params;
            if (!tag) {
                return error({statusCode: 400, message: 'No tag'});
            }
            let tagResponse = await fetch(`${process.env.API_ORIGIN}/search/tag/${tag}`, {credentials: 'include'});
            if (!tagResponse.ok) {
                return error({statusCode: tagResponse.status, message: 'Could not get tag'});
            }
            let {results: tags} = await tagResponse.json();
            let descriptions = [];
            for (let i in tags) {
                let tag = tags[i];
                descriptions[i] = getSafeHtml(tag.description, buildTagMarker(tag.description));
            }
            return {
                tags,
                descriptions
            };
        },
        watchQuery: ['tag'],
        methods: {
            handleClicks
        },
        created() {
            this.builderTypePathFragment = {'script': 'script', 'script_chain': 'scriptchain'};
        }
    };
</script>