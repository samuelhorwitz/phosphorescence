<template>
    <main>
        <h1>Get some test data</h1>
        <textarea width="500" height="1000" v-model="testDataJson"></textarea>
    </main>
</template>

<script>
    import {getPlaylist, shuffle, getLocalTokens, tryToRefresh} from '~/util/spotify';
    const testPlaylists = ['5WdZ7SDAgBE8kK0EYkn3xj', '5lSgExb6yTKfLusGag7bm7'];
    const ignoreDueToTrained = ["006UwN5OTJdw3mLJn7KTlf","029GZRvHYKbxwkyKKYorph","08TD9e3yZEfCtdBUc2qEZ8","0Dm2PvF9qx7RnCeM1x1r1I","0DRjqCEf5bflI8luIt9PNY","0ejEe9ipfqsUV7ZuJlULXH","0IEQf3fLNfl7GSxIMZBTIL","0nTs7JsCo6lJJ3Yhzz5atT","0p2BvTX1CoUjU3wNgOwWlx","0Rppw4AfieqmQ38I2TcFpO","0soeGxCAy8LwxjGDjOShf6","0sSrY0wjDZzS1QJAunfaeG","0suBYnerhXJMu3hLigdPDm","0wrutF2Svuprj6yTAJkJtw","0XW39C18hnak93aUDUgID5","0yVaUTYKxdjQotqw1OGQVx","10MpdSRHZmSPR1m0SUh1R5","19QUNfdeb5aWHObd54vpvu","1bnyOUx6faiyu5h6byy73C","1CqLOu8hIIn1iuh40bbzXJ","1dby8pyzTUhKDeJZEhmGxR","1DSKB58inLQfTbmApg3920","1EsDihsTfSKXLVqid1hAr0","1HVmUtidpyBEvuJ6pZEx3b","1L0gJuJqW9PaTG7VqSIiZ6","1mAUds7ZBlGyRrA5MuwVBx","1ndMtTEDyYlUNtGtCFfPkX","1onhfO4MIeXQHEHoGIPUIW","1PhRbqbfsa4mUwkCmRRAk1","1SJNqDkfZRx0keSUMU8dTw","1ssrzlbRgjhsiBiM9t7Mad","21DBxXQkvkgpJ0ks7532NI","21mtSLaMvm8N5FvnzO8TT5","25JPB0o9ZQqj7tf1vmBBQx","296hRskwVjQhca2ZQZPAXQ","29Fyyu7x4emFq8CfSI5yXr","2chHJwFPYmskM393ZpZl8W","2foFln1kOOWgWF41bDeMGe","2GdKmL9WkVag0UP23yQCDJ","2GLRwpqSfzcJmjQlSGKsC7","2JcqYs3baaHYqMmlwUn0CQ","2l6x22Bl6xhhMuDBPhGgO7","2M0L7xPvQjAUHK950HBQcz","2oDxVGK5bBFIWR13VQIxBr","2q3D7qCMGgQ6mX0sXQukPW","2rK7hzFsL3KqlzHW13JSNE","2uTKlGCBd9GQ1CEOgZ5qRL","2xAOSQbRN4nV31OQ2TW1Gz","2XCA5r8qacW1g8sNBcNn6L","2XmsZtjm0iBEXdQDns4ys2","2z96opOidetDscj34oTT1H","30gEmd0lWxWKfwNWUDI1ry","33ivKU3ch2Uex4ZrpB4tkS","34XFMmrIsOWPiWryoRmNXs","34yMP8IIBXz9WxprFrovZO","38W1IohVo9EB2gaGIozGig","3B9aLGtTYRaFusDf0gMI0C","3CBuO83LHWYycIKoclJrDP","3CpoqtUPzfoBpvfyQ5xQYK","3Dg26dk5afwaDPkaoAA4dU","3DkUeJKLMYjP9qJMOZM8UN","3Dp7YaABByyUikCf8q3Z8E","3e8KIcAyuJauXu4Y0u8emb","3GOKEjQFQdSauRn2AJRZI1","3ldGEFvu9oJc432DSVhIZt","3NZwvPydbzUchK271JYFr9","3ovw6hZBoSPufBGLUVis7Q","3OxfzKG0uRePXc7Xj9YSF0","3Q5ddPoj7he3HAvjKpAaC1","3rZkRFzfbV4ZRxg5Xpg68c","3SdjsAwSP5FEB2QWL6qKvm","3SvodsZQvZtwy5szNUeHsK","3x8nnNDfa0aLvbq52SHwri","3ZVVhtqSnZZNJhOwINkJA9","443jJ9JZQmnMBGn2ZhNYxK","45phaDzVD89YsjDeXJrlbJ","48S9Ueh0n0nUH1Sgh45eZF","493sOxeUzIv1yLod4Z6l4n","4acv9cybr3QdBcQ5WpyiI1","4blexrVa2uCT6o2EKlJdWi","4bMSGNx8pbJZYHlAXAH0wu","4CkQiW6QKUFbcuWKuIc3mx","4co2NU8qyWqWxShE2Vfpx6","4dHL3zs6u7CkE8HDuwM7o5","4dOefkhSbaRgGfryQ67t8D","4fu3sxrZ9RFMiTSeBiVLkQ","4fV1k0OyGMZgYIdK2gHJgc","4GvCnjNmn5PRXOegpWQMTP","4hUxfQlRQc15gXuoT7ajs3","4jFRrJgkKg3bUL4IDZQ9Cs","4jy7tsWx2FZjXspLV5uZkQ","4Qme2qmSqORtRRM9hTC5hH","4RGWfjbxxW9Z8E5DZMaGcV","4rsHDST3FgBWEQk55M9i6F","4v1UjUdFQ7o2C2SQuOVpGn","4VpO7lgq7K6YtctgkjifZU","4VTIWTdwAzMQspQtpcVti8","4xB6DpaRgtx6mvSeJbCx5D","4ZAkzg3RYVeVyCqLiiaEx6","532ztzOKzl1Bl8NRF7OuRY","56RyGLm97SwPbMsofPrPtW","5eUUGuDmaBEychzDsDjfwK","5GZY6SKdheLR2m1pGLMkdb","5h6CNFzBLC9Ca3MgtTMnMm","5hhZfa9bYiHKvNJLyfNvHS","5kQ9eRG4btu43OSvsD1zrx","5lbXgwzMBiUxzd4gpy5bX6","5nZQXnDmVWKchJiXIK5VxZ","5svKBdzQoAhy2ml2wjuoRp","5syh7hAzAnpLAHlWjyFWvz","5t8RM0ORGLtufvsGMeG0zp","5TqgdZqGl8l25UMq9xX67E","5VIVZAQVSfJRD8f4XFbgBL","5ZhOkgKfw8eyOc7Cu17Pvb","6276TOaayx29M6Sgf1bd4J","6bH9QbhgR0QGFQyEtKt3GR","6d9K7tOQ7VP1TAWKRaQJwI","6JqMCvsqcS4ZrOsylFmYg0","6K2gQ5DIl97b8UOxAJxljn","6k2pBpf10f3DNWKykGGJjn","6KfbtqiuVB8GoH0CpnJpR2","6N0EGpoggDg04bthrj2jGs","6TBwdLNcuCScR9kzesDbve","6wv5lDxfIh5e27tgvi5diu","6xnUficHSWNOyY02bNea3G","6Z5hYMAqFiqiq3VgFp8hGz","6Z9D86FCHE6PDlpUJKwuKP","6ZfzbvseI7sxis7uFUEAHx","6ZTi2Kk6ax3bvBl9S9dU2r","753buqZYXLZ6OjfePwOhHY","7BhB11GDezuGWN0XpzVako","7CFSUxz0EGR7YivejoVErE","7GEkkoIG1x3qCUvqgsUq6w","7lBCnpuqVD51GbRtsupg8u","7LrmcBh9qpuaQ8lEXMjuj2","7sbV1r4ONyyIrDvde7GCA6","7sDP6oLTlyJOtDymeNMn03","7uksvEjKFDo8ZYoUSWM0p3"];

    export default {
        data() {
            return {
                limit: 1500,
                testData: null
            };
        },
        computed: {
            testDataJson() {
                if (!this.testData) {
                    return '';
                }
                return JSON.stringify(this.testData, null, 4);
            }
        },
        async created() {
            if (!process.browser) {
                return;
            }
            let {access, refresh, expires, okay} = getLocalTokens();
            if (!okay) {
                let response = await tryToRefresh(refresh);
                access = response.access;
                refresh = response.refresh;
                expires = response.expires;
            }
            this.$store.dispatch('session/tokens', {access, refresh, expires});
            let allTracks = [];
            let trackAnalysis = [];
            for (let playlist of testPlaylists) {
                let {tracks, allTracksAnalysis} = await getPlaylist(playlist);
                allTracks = allTracks.concat(tracks);
                trackAnalysis = trackAnalysis.concat(allTracksAnalysis);
            }
            allTracks = shuffle(allTracks);
            let analysis = trackAnalysis.reduce((acc, cur) => {
                if (!cur) return acc;
                acc[cur.id] = cur;
                return acc;
            }, {});
            let idx = 0;
            this.testData = {};
            while (Object.keys(this.testData).length < this.limit) {
                let track = allTracks[idx++];
                if (ignoreDueToTrained.includes(track.track.id) || this.testData[track.track.id]) {
                    continue;
                }
                this.testData[track.track.id] = {track: track, analysis: analysis[track.track.id]};
            }
        }
    };
</script>