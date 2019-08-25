<template>
    <section>
        <h2>#{{$route.params.tag}}</h2>
        <ol v-if="tags">
            <li v-for="(tag, index) of tags">
                <h3>
                    <span class="name">{{tag.name}}</span>&nbsp;&horbar;&nbsp;<span class="authorName">{{tag.authorName}}</span>
                </h3>
                <p v-if="descriptions[index]" v-html="descriptions[index]" @click="handleClicks"></p>
                <p v-if="!descriptions[index]">{{tag.description}}</p>
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
    section {
        margin-left: 2em;
    }

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

<style>
    mark.searchResult {
        background-color: inherit;
        font-weight: bold;
        text-decoration: underline;
        background-color: magenta;
        color: white;
    }
</style>

<script>
    import {getAccessToken} from '~/assets/session';
    import {getSafeHtml, buildTagMarker, handleClicks} from '~/assets/safehtml';
    import {debounce} from 'lodash';

    export default {
        layout: 'marketplace',
        async fetch({store, error}) {
            await getAccessToken();
            let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
            if (!userResponse.ok) {
                return error({statusCode: userResponse.status, message: "Could not get user information"});
            }
            let {user} = await userResponse.json();
            store.commit('user/user', user);
        },
        watchQuery: ['tag'],
        data() {
            return {
                tags: [],
                descriptions: []
            };
        },
        watch: {
            tags: {
                immediate: true,
                handler: 'buildMarkedText'
            }
        },
        async created() {
            let tag = this.$route.params.tag;
            if (!tag) {
                return;
            }
            let tagResponse = await fetch(`${process.env.API_ORIGIN}/scripts/search-tag?tag=${tag}`, {credentials: 'include'});
            if (!tagResponse.ok) {
                return;
            }
            let {results} = await tagResponse.json();
            this.tags = results;
        },
        methods: {
            buildMarkedText: debounce(function (tags) {
                let descriptions = [];
                for (let i in tags) {
                    let tag = tags[i];
                    descriptions[i] = getSafeHtml(tag.description, buildTagMarker(tag.description));
                }
                this.descriptions = descriptions;
            }, 200),
            handleClicks
        }
    };
</script>