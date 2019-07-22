<template>
    <footer>
        <ol>
            <li><span @click="flash" @mouseover="tosHover = true" @mouseout="tosHover = false" class="wrapper"><nuxt-link to="/legal/tos" target="_blank"><glitch :text="'[Terms of Service]'" :enable="tosHover"/></nuxt-link></span></li>
            <li ><span @click="flash" @mouseover="privacyHover = true" @mouseout="privacyHover = false" class="wrapper"><nuxt-link to="/legal/privacy" target="_blank"><glitch :text="'[Privacy]'" :enable="privacyHover"/></nuxt-link></span></li>
            <li><span @mouseover="githubHover = true" @mouseout="githubHover = false" class="wrapper"><a target="_blank" href="https://github.com/samuelhorwitz/phosphorescence"><glitch :text="'[GitHub]'" :enable="githubHover"/></a></span></li>
            <li class="spotify">
                <div class="spotifyLogo">
                    <a target="_blank" href="https://spotify.com"><span class="spotifyPreCopy">Content made available by</span><img alt="Spotify" class="spotifyLogoImg" src="/images/spotify_small.png"></a>
                </div>
            </li>
        </ol>
    </footer>
</template>

<style scoped>
    .spotifyLogo {
        position: absolute;
        bottom: -30px;
        right: 10px;
    }

    .spotifyLogo a {
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .spotifyPreCopy {
        font-size: .1em;
        margin-right: .5em;
        width: 70px;
        display: none;
    }

    .spotifyLogoImg {
        width: 70px;
    }

    footer {
        width: 100%;
        height: 3.2em;
        grid-column: 1 / 3;
        display: flex;
    }

    ol {
        display: flex;
        justify-content: center;
        align-items: center;
        margin: 0px;
        padding: 0px;
        font-family: courier;
        flex: 1;
    }

    li {
        display: inline;
        flex: 100;
        text-align: center;
        cursor: pointer;
        font-size: 1.5em;
    }

    li.spotify {
        position: relative;
        flex: 0;
        bottom: 0.5em;
    }

    .wrapper {
        display: inline-block;
        position: relative;
    }

    a {
        text-decoration: none;
        color: white;
    }

    @media only screen and (max-height: 449px) {
        footer {
            grid-column: 1 / 3;
            grid-row: 4 / 5;
            margin-top: auto;
        }

        li.spotify {
            position: inherit;
        }

        .spotifyLogo {
            position: absolute;
            bottom: 5px;
            right: 5px;
        }
    }

    @media only screen and (max-height: 249px) {
        footer {
            grid-row: 2 / 5;
        }
    }

    @media only screen and (max-width: 749px) and (max-height: 449px) {
        footer {
            height: 3em;
            font-size: 10px;
        }
    }

    @media only screen and (min-height: 450px) and (max-width: 1099px) {
        footer {
            height: 5em;
        }
    }

    @media only screen and (max-height: 449px) and (max-width: 1099px) {
        footer {
            height: 3em;
            font-size: 10px;
        }
    }

    @media only screen and (max-width: 1099px) {
        footer {
            grid-column: 1 / 2;
        }

        ol {
            justify-content: flex-start;
            margin-left: 2em;
        }

        li {
            flex: 0;
        }

        li.spotify {
            position: inherit;
        }

        .spotifyLogo {
            bottom: 5px;
            right: 5px;
        }
    }
</style>

<script>
    import glitch from '~/components/glitch';

    export default {
        components: {
            glitch
        },
        data() {
            return {
                tosHover: false,
                privacyHover: false,
                githubHover: false
            };
        },
        methods: {
            flash() {
                this.$store.dispatch('loading/loadFlash');
            }
        }
    };
</script>