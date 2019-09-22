<template>
    <footer>
        <ol>
            <li><span @click="flash" @mouseover="blogHover = true" @mouseout="blogHover = false" class="wrapper"><a target="_blank" rel="external noopener" href="https://blog.phosphor.me/"><glitch text="[Blog]" :enable="blogHover"/></a></span></li>
            <li><span @click="flash" @mouseover="legalHover = true" @mouseout="legalHover = false" class="wrapper"><nuxt-link to="/legal" target="_blank"><glitch text="[Legal]" :enable="legalHover"/></nuxt-link></span></li>
            <li class="socialWrapper">
                <ul>
                    <li class="social" title="Facebook"><span @mouseover="facebookHover = true" @mouseout="facebookHover = false" class="wrapper"><a target="_blank" rel="external noopener" href="https://facebook.com/phosphorme/"><glitch class="fab" raw-text="&#xf09a;" :enable="facebookHover"/></a></span></li>
                    <li class="social" title="Twitter"><span @mouseover="twitterHover = true" @mouseout="twitterHover = false" class="wrapper"><a target="_blank" rel="external noopener" href="https://twitter.com/phosphor_me/"><glitch class="fab" raw-text="&#xf099;" :enable="twitterHover"/></a></span></li>
                    <li class="social" title="Instagram"><span @mouseover="instagramHover = true" @mouseout="instagramHover = false" class="wrapper"><a target="_blank" rel="external noopener" href="https://instagram.com/phosphor_me/"><glitch class="fab" raw-text="&#xf16d;" :enable="instagramHover"/></a></span></li>
                    <li class="social" title="Spotify"><span @mouseover="spotifyHover = true" @mouseout="spotifyHover = false" class="wrapper"><a target="_blank" rel="external noopener" href="https://open.spotify.com/user/spv1hpk4dj1qwneuxwg4yg0cn"><glitch class="fab" raw-text="&#xf1bc;" :enable="spotifyHover"/></a></span></li>
                    <li class="social" title="GitHub"><span @mouseover="githubHover = true" @mouseout="githubHover = false" class="wrapper"><a target="_blank" rel="external noopener" href="https://github.com/samuelhorwitz/phosphorescence"><glitch class="fab" raw-text="&#xf09b;" :enable="githubHover"/></a></span></li>
                </ul>
            </li>
            <li class="spotify">
                <div class="spotifyLogo">
                    <a target="_blank" rel="external noopener" href="https://spotify.com"><span class="spotifyPreCopy">Content made available by</span><img alt="Spotify" class="spotifyLogoImg" src="/images/spotify_small.png"></a>
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
        position: relative;
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

    ul {
        margin: 0px;
        padding: 0px;
        display: flex;
        flex-direction: row;
    }

    li {
        display: inline;
        flex: 100;
        text-align: center;
        cursor: pointer;
        font-size: 1.5em;
    }

    li.social {
        flex: unset;
        margin-left: 3ex;
    }

    li.spotify {
        position: relative;
        flex: 0;
        bottom: 0.5em;
    }

    li.socialWrapper {
        font-size: 1.25em;
        margin-left: 1ex;
        margin-top: 0.25ex;
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
            grid-column: 1 / 4;
            grid-row: 4 / 5;
            margin-top: auto;
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
            grid-column: 1 / 4;
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

    @media only screen and (max-width: 1299px) {
        li.social {
            margin-left: 1ex;
        }
    }

    @media only screen and (display-mode: fullscreen) and (orientation: portrait) {
        .spotifyLogo {
            bottom: 0px;
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
                legalHover: false,
                githubHover: false,
                facebookHover: false,
                twitterHover: false,
                instagramHover: false,
                spotifyHover: false,
                blogHover: false
            };
        },
        methods: {
            flash() {
                this.$store.dispatch('loading/loadFlash');
            }
        }
    };
</script>