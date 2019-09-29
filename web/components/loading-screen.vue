<template>
    <aside class="loading" v-show="!$store.getters['tracks/playlistLoaded']">
        <h2><b class="appname">Phosphorescence</b>
        builds coherent playlists for <b class="genre">trance</b> and <b class="genre">chill-out</b> listeners.
        </h2>
        <div class="wrapper">
            <ul>
                <li class="loadingMessage" v-for="loadMessage in $store.state.loading.descriptions" :class="{done: loadMessage.done}">
                    {{loadMessage.description}}...
                </li>
            </ul>
        </div>
        <div class="boilerplate">
            <p>We are an <a target="_blank" rel="external noopener" href="https://github.com/samuelhorwitz/phosphorescence">open source project</a> and 100% free to use.</p>
            <p>Spotify provides the music, however we are <em>not</em> affiliated with Spotify in any way.</p>
        </div>
    </aside>
</template>

<style scoped>
    aside.loading {
        overflow: auto;
        color: white;
        padding: 1em;
        flex: 1;
        z-index: 10000000;
        display: flex;
        flex-direction: column;
        align-items: center;
        font-size: 16px;
        width: 100vw;
        height: 100%;
        box-sizing: border-box;
    }

    .wrapper {
        display: flex;
        width: 100%;
        flex: 1;
        justify-content: center;
        align-items: center;
        margin-bottom: auto;
    }

    ul {
        padding: 0px;
        font-size: 0.4em;
        margin: 0px;
        width: 375px;
        height: 100%;
        max-height: 200px;
    }

    li.loadingMessage {
        display: block;
        font-size: 3em;
    }

    li.loadingMessage:not(.done)::before {
        content: '‣';
        color: magenta;
        display: inline-block;
        width: 1ex;
    }

    li.loadingMessage.done::before {
        content: '✓';
        color: cyan;
        display: inline-block;
        width: 1em;
    }

    .boilerplate {
        font-size: 0.75em;
        color: #eee;
        width: 100%;
        max-width: 65em;
    }

    h2 {
        margin-top: 0px;
        margin-bottom: 1ex;
        font-size: 2.2em;
        max-width: 22em;
    }

    a {
        color: cyan;
    }

    a:hover {
        color: magenta;
    }

    b.appname {
        font-variant: small-caps;
        font-weight: bolder;
        font-family: 'Varela';
    }

    b.genre {
        font-weight: inherit;
        color: aqua;
    }

    @media only screen and (max-height: 999px) {
        aside.loading {
            flex: 1;
        }
    }
</style>