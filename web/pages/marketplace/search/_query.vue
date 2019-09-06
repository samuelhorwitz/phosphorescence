<template>
    <section>
        <h2 v-if="searchResults">Search Results</h2>
        <h2 v-if="!searchResults">No search results found</h2>
        <ol v-if="searchResults">
            <li v-for="(searchResult, index) of searchResults">
                <h3>
                    <router-link :to="'/marketplace/' + builderTypePathFragment[searchResult.resultType] + '/' + encodeURIComponent(searchResult.id)"><span class="name" v-html="names[index]"></span></router-link>&nbsp;&horbar;&nbsp;<span class="authorName" v-html="authorNames[index]"></span>
                </h3>
                <p v-html="descriptions[index]" @click="handleClicks"></p>
            </li>
        </ol>
        <footer v-if="searchResults">
            End of results
        </footer>
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
        position: relative;
    }

    p {
        margin: 0px;
        margin-left: 1em;
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
        background-color: teal;
        color: white;
    }
</style>

<script>
    import {getSafeHtml, buildMarker, buildTagMarker, combineMarkers, handleClicks} from '~/assets/safehtml';
    import verifiedBadge from '~/components/marketplace/verifiedbadge';

    export default {
        layout: 'marketplace',
        components: {verifiedBadge},
        async asyncData({store, params, error}) {
            let {query} = params;
            store.commit('marketplace/setQuery', query);
            let searchResponse = await fetch(`${process.env.API_ORIGIN}/search/${encodeURIComponent(query)}`, {credentials: 'include'});
            if (!searchResponse.ok) {
                return error({statusCode: userResponse.status, message: 'Could not perform search'});
            }
            let {results: searchResults} = await searchResponse.json();
            let names = [], authorNames = [], descriptions = [];
            let node = document.createElement('mark');
            node.classList.add('searchResult');
            for (let i in searchResults) {
                let searchResult = searchResults[i];
                names[i] = getSafeHtml(searchResult.name, buildMarker(JSON.parse(JSON.stringify(searchResult.nameMarks || [])), node));
                authorNames[i] = getSafeHtml(searchResult.authorName, buildMarker(JSON.parse(JSON.stringify(searchResult.authorNameMarks || [])), node));
                let descriptionMarker = buildMarker(JSON.parse(JSON.stringify(searchResult.descriptionMarks || [])), node);
                let descriptionTagMarker = buildTagMarker(searchResult.description);
                descriptions[i] = getSafeHtml(searchResult.description, combineMarkers(descriptionTagMarker, descriptionMarker));
                if (searchResult.partialDescription) {
                    descriptions[i] = `&hellip;${descriptions[i]}&hellip;`;
                }
            }
            return {
                searchResults,
                names,
                authorNames,
                descriptions
            };
        },
        watchQuery: ['query'],
        methods: {
            handleClicks
        },
        beforeRouteLeave(to, from, next) {
            this.$store.commit('marketplace/clearQuery');
            next();
        },
        created() {
            this.builderTypePathFragment = {'script': 'script', 'script_chain': 'scriptchain'};
        }
    };
</script>