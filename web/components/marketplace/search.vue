<template>
    <div>
        <input type="text" ref="query" v-model="query" @focus="showRecommended" @blur="hideRecommended" @keydown.arrow-up="moveCursorUp" @keydown.arrow-down="moveCursorDown" @keydown.enter="handleEnter" @keydown.esc="handleEscape">
        <ul v-if="isFocused && !recommendedHidden && recommendedQueries && recommendedQueries.length > 1" ref="dropdown">
            <li v-for="(recommended, index) of recommendedQueries" :class="{selected: index + 1 === cursor}" tabindex="0" @click="selectRecommendation(index)" @focus="setCursor(index)" @keydown.arrow-up="moveCursorUp" @keydown.arrow-down="moveCursorDown" @keydown.enter="handleEnter" @keydown.esc="handleEscape">
                {{recommended}}
            </li>
        </ul>
    </div>
</template>

<style scoped>
    div {
        flex: 1;
        display: flex;
        position: relative;
    }

    input {
        flex: 1;
        margin: 0px;
        padding: 0px;
        padding-left: 0.25em;
        font-size: 2em;
    }

    input:focus {
        outline: none;
    }

    ul {
        position: absolute;
        top: 3em;
        list-style: none;
        margin: 0px;
        padding: 0px;
        width: 100%;
        background-color: white;
        border: 2px solid rgb(238, 238, 238);
        border-top: 0px;
        box-sizing: border-box;
    }

    li {
        padding: 0px 0.25em;
        cursor: pointer;
    }

    li:hover {
        background-color: cyan;
        color: white;
    }

    li:focus {
        outline: none;
    }

    li.selected {
        background-color: cyan;
        color: white;
    }
</style>

<script>
    import {mapState} from 'vuex';
    import {debounce} from 'lodash';

    export default {
        data() {
            return {
                recommendedHidden: false,
                recommendedQueries: null,
                cursor: 0,
                isFocused: false
            };
        },
        computed: {
            query: {
                get() { return this.$store.state.marketplace.query; },
                set(query) { this.$store.commit('marketplace/setQuery', query); }
            }
        },
        watch: {
            query: 'getRecommendedQueries'
        },
        methods: {
            showRecommended() {
                this.isFocused = true;
                this.recommendedHidden = false;
            },
            hideRecommended(e) {
                if (e.relatedTarget && e.relatedTarget.parentNode === this.$refs.dropdown) {
                    return;
                }
                this.isFocused = false;
                this.recommendedHidden = true;
                this.cursor = 0;
            },
            moveCursorDown(e) {
                if (!this.recommendedQueries || this.cursor === this.recommendedQueries.length) {
                    return;
                }
                e.preventDefault();
                this.cursor++;
            },
            moveCursorUp(e) {
                if (this.cursor === 0) {
                    return;
                }
                e.preventDefault();
                this.cursor--;
            },
            setCursor(index) {
                this.cursor = index + 1;
            },
            handleEnter() {
                if (this.cursor === 0) {
                    this.search();
                } else {
                    this.searchFromCursor();
                }
            },
            handleEscape() {
                this.recommendedHidden = true;
                this.cursor = 0;
            },
            search() {
                if (!this.query) {
                    this.$router.push({name: 'marketplace', params: {query: this.query}});
                } else {
                    this.$router.push({name: 'marketplace-search-query', params: {query: this.query}});
                }
                this.$refs.query.blur();
            },
            searchFromCursor() {
                if (this.cursor === 0 || this.cursor >= this.recommendedQueries.length) {
                    return;
                }
                this.query = this.recommendedQueries[this.cursor - 1] + ' ';
                this.recommendedHidden = true;
                this.search();
            },
            selectRecommendation(index) {
                this.query = this.recommendedQueries[index] + ' ';
                this.search();
            },
            getRecommendedQueries: debounce(async function (query) {
                if (!query) {
                    return;
                }
                let queryRecommendationResponse = await fetch(`${process.env.API_ORIGIN}/scripts/query-recommendation?query=${encodeURIComponent(query)}`, {credentials: 'include'});
                if (!queryRecommendationResponse.ok) {
                    this.recommendedQueries = null;
                    return;
                }
                let {recommended} = await queryRecommendationResponse.json();
                this.recommendedQueries = recommended;
                this.recommendedHidden = false;
                this.cursor = 0;
            }, 200)
        }
    };
</script>
