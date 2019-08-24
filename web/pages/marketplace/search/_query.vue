<template>
    <section>
        <h2>Search Results</h2>
        <ol v-if="searchResults">
            <li v-for="(searchResult, index) of searchResults">
                <h3 v-if="names[index] && authorNames[index]">
                    <span class="name" v-html="names[index]"></span>&nbsp;&horbar;&nbsp;<span class="authorName" v-html="authorNames[index]"></span>
                </h3>
                <h3 v-if="!names[index] || !authorNames[index]">
                    <span class="name">{{searchResult.name}}</span>&nbsp;&horbar;&nbsp;<span class="authorName">{{searchResult.authorName}}</span>
                </h3>
                <p v-if="descriptions[index]" v-html="descriptions[index]"></p>
                <p v-if="!descriptions[index]">&hellip;{{searchResult.description}}&hellip;</p>
            </li>
        </ol>
        <footer v-if="searchResults">
            End of results
        </footer>
        <aside v-if="!searchResults">
            No search results found
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
    import {debounce} from 'lodash';

    function buildMarkedText(text, marks) {
        if (!marks || marks.length === 0 || marks.length % 2 !== 0) {
            return text;
        }
        let el = document.createElement('div');
        let lead = document.createTextNode(text.substring(0, marks[0]));
        el.appendChild(lead);
        let lastMark;
        for (let i = 0; i < marks.length; i += 2) {
            let startMark = marks[i];
            let endMark = marks[i + 1];
            if (lastMark) {
                let t = document.createTextNode(text.substring(lastMark, startMark));
                el.appendChild(t);
            }
            let mark = document.createElement('mark');
            mark.classList.add('searchResult');
            mark.innerText = text.substring(startMark, endMark);
            el.appendChild(mark);
            lastMark = endMark;
        }
        let tail = document.createTextNode(text.substring(marks[marks.length - 1]));
        el.appendChild(tail);
        return el.innerHTML;
    }

    export default {
        layout: 'marketplace',
        async fetch({store, params, error}) {
            await getAccessToken();
            let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
            if (!userResponse.ok) {
                return error({statusCode: userResponse.status, message: 'Could not get user information'});
            }
            let {user} = await userResponse.json();
            store.commit('user/user', user);
            let {query} = params;
            store.commit('marketplace/setQuery', query);
            let searchResponse = await fetch(`${process.env.API_ORIGIN}/scripts/search?query=${encodeURIComponent(query)}`, {credentials: 'include'});
            if (!searchResponse.ok) {
                return error({statusCode: userResponse.status, message: 'Could not perform search'});
            }
            let {results} = await searchResponse.json();
            store.commit('marketplace/setSearchResults', results);
        },
        data() {
            return {
                names: [],
                authorNames: [],
                descriptions: []
            };
        },
        computed: {
            searchResults: {
                get() { return this.$store.state.marketplace.searchResults; },
                set(searchResults) { this.$store.commit('marketplace/setSearchResults', searchResults); }
            }
        },
        watch: {
            searchResults: {
                immediate: true,
                handler: 'buildMarkedText'
            }
        },
        methods: {
            buildMarkedText: debounce(function (searchResults) {
                let names = [], authorNames = [], descriptions = [];
                for (let i in searchResults) {
                    let searchResult = searchResults[i];
                    names[i] = buildMarkedText(searchResult.name, searchResult.nameMarks);
                    authorNames[i] = buildMarkedText(searchResult.authorName, searchResult.authorNameMarks);
                    descriptions[i] = '&hellip;' + buildMarkedText(searchResult.description, searchResult.descriptionMarks) + '&hellip;';
                }
                this.names = names;
                this.authorNames = authorNames;
                this.descriptions = descriptions;
            }, 200)
        },
        beforeDestroy() {
            this.$store.commit('marketplace/clearQuery');
        }
    };
</script>