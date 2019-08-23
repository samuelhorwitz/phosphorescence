<template>
    <div>
        <input type="text" v-model="query" v-on:keydown.stop>
        <ul>
            <li v-for="recommended of recommendedQueries">
                {{recommended}}
            </li>
        </ul>
    </div>
</template>

<style>
    li {
        appearance: none;
        color: white;
    }
</style>

<script>
    import {debounce} from 'lodash';
    import {getAccessToken} from '~/assets/session';

    export default {
        async fetch({store, error}) {
            await getAccessToken();
            let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
            if (!userResponse.ok) {
                return error({statusCode: userResponse.status, message: "Could not get user information"});
            }
            let {user} = await userResponse.json();
            store.commit('user/user', user);
        },
        data() {
            return {
                query: "",
                recommendedQueries: null
            };
        },
        watch: {
            query: 'getRecommendedQueries'
        },
        methods: {
            getRecommendedQueries: debounce(async function (query) {
                let queryRecommendationResponse = await fetch(`${process.env.API_ORIGIN}/scripts/query-recommendation?query=${encodeURIComponent(query)}`, {credentials: 'include'});
                this.recommendedQueries = null;
                if (!queryRecommendationResponse.ok) {
                    return;
                }
                let {recommended} = await queryRecommendationResponse.json();
                this.recommendedQueries = recommended;
            }, 200)
        }
    };
</script>