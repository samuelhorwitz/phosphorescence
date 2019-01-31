<template>
    <aside>
        <header class="logoContainer">
            <h1 @click="flash" :class="{flashing: isFlashing}"><glitch :text="'Phorland'" :enable="isGlitching" /></h1>
            <h2 @click="flash" :class="{flashing: isMiamiVicing}">JavaScript IDE</h2>
        </header>
        <section>
            <p>Welcome to the Phosphorescence IDE, Phorland. Here you may edit your own playlist builders.</p>
            <p>Please check out <a target="_blank" href="/docs/eos/api.html">our documentation</a>. Debugging is done via your browsers native developer tools.</p>
            <p><em>Please Note:</em> If you save code to our servers, you are relinquishing all rights to that code under the <a target="_blank" href="https://creativecommons.org/share-your-work/public-domain/cc0/">Creative Commons 0</a> public domain dedication.</p>
            <p>This is so that we do not have to worry about code attribution and contributor management when allowing users to fork your work. If you do not wish to release your work into the public domain, you may <a target="_blank" href="https://github.com/samuelhorwitz/phosphorescence">find us on GitHub</a> and develop from your own instance of the app.</p>
        </section>
    </aside>
</template>

<style scoped>
    .logoContainer {
        position: relative;
    }

    aside {
        background: rgb(0,8,123);
        background: linear-gradient(180deg, rgba(0,8,123,1) 21%, rgba(193,193,193,1) 100%);
        border-right: 1px solid black;
        padding: 1em;
    }

    h1 {
        font-family: Impact;
        font-size: 4em;
        color: teal;
        text-shadow: -1px -1px 0 black, 1px -1px 0 black, -1px 1px 0 black, 1px 1px 0 black;
        margin: 0px;
        text-align: center;
        cursor: pointer;
    }

    h1.flashing {
        color: aqua;
        text-shadow: none;
    }

    h2 {
        font-family: Impact;
        font-size: 2em;
        color: rgb(0, 8, 123);
        text-shadow: -1px -1px 0 rgb(193,193,193), 1px -1px 0 rgb(193,193,193), -1px 1px 0 rgb(193,193,193), 1px 1px 0 rgb(193,193,193);
        margin: 0px;
        text-align: center;
        cursor: pointer;
        height: 1.5em;
    }

    h2.flashing {
        color: magenta;
        text-shadow: -1px -1px 0 aquamarine, 1px -1px 0 aquamarine, -1px 1px 0 aquamarine, 1px 1px 0 aquamarine;
        font-family: 'Caveat';
        text-decoration: underline;
    }

    p {
        color: white;
    }

    a {
        color: lightgray;
    }

    section {
        overflow: hidden;
        overflow-y: scroll;
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
                isGlitching: false,
                isFlashing: false,
                isMiamiVicing: false
            };
        },
        methods: {
            flash() {
                if (this.isGlitching) {
                    return;
                }
                this.isGlitching = true;
                this.isFlashing = true;
                this.isMiamiVicing = true;
                setTimeout(() => this.isGlitching = false, 2000);
                setTimeout(() => this.isFlashing = false, 1000);
                setTimeout(() => this.isMiamiVicing = false, 1700);
            }
        }
    };
</script>
