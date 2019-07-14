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
            <hr>
            <h3>Save Your Script</h3>
            <h4 v-if="$store.state.ide.scriptData != null && !$store.getters['ide/isScriptOwnedByUser']">This script is not owned by you, it is read-only. Fork it if you wish to make changes.</h4>
            <form @submit.prevent>
                <label>
                    Name
                    <input type="text" v-model="scriptName" placeholder="Untitled Script" :disabled="$store.state.ide.scriptData != null && !$store.getters['ide/isScriptOwnedByUser']">
                </label>
                <hr>
                <div v-if="$store.state.ide.scriptData == null" class="buttonContainer">
                    <button @click="saveNew">Save New Script</button>
                </div>
                <div v-if="$store.state.ide.scriptData != null && $store.state.user.user && $store.getters['ide/isScriptOwnedByUser']" class="buttonContainer">
                    <button @click="saveDraft">Save Draft</button>
                    <button @click="publish">Publish</button>
                    <button @click="duplicate">Duplicate</button>
                </div>
                <div v-if="$store.state.ide.scriptData != null && $store.state.user.user && !$store.getters['ide/isScriptOwnedByUser']" class="buttonContainer">
                    <button @click="fork">Fork</button>
                </div>
            </form>
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

    h3 {
        color: white;
        text-align: center;
    }

    h4 {
        color: aquamarine;
        font-size: 1em;
    }

    label {
        color: white;
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

    input {
        width: 100%;
    }

    .buttonContainer button {
        font-weight: bold;
        appearance: none;
        margin: 0px;
        height: 2.5em;
        border: 1px outset gray;
        background-color: gray;
        margin-bottom: 1em;
        cursor: pointer;
        outline: 0px;
        flex: 1;
        margin-left: 1px;
        margin-right: 1px;
    }

    .buttonContainer {
        display: flex;
        width: 100%;
    }
</style>

<script>
    import glitch from '~/components/glitch';

    export default {
        components: {
            glitch
        },
        data() {
            let initialScriptName;
            if (this.$store.state.ide.scriptData) {
                initialScriptName = this.$store.state.ide.scriptData.name;
            }
            return {
                isGlitching: false,
                isFlashing: false,
                isMiamiVicing: false,
                scriptName: initialScriptName
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
            },
            async saveNew() {
                let createResponse = await fetch(`${process.env.API_ORIGIN}/script`, {
                    method: 'POST',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name: this.scriptName,
                        script: this.$store.state.ide.script
                    })
                });
                let {create} = await createResponse.json();
                this.$router.push({path: `/editor/${create.id}`});
            },
            async saveDraft() {
                let updateResponse = await fetch(`${process.env.API_ORIGIN}/script/${this.$store.state.ide.scriptData.id}`, {
                    method: 'PUT',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name: this.scriptName,
                        script: this.$store.state.ide.script
                    })
                });
                let {update} = await updateResponse.json();
            },
            async publish() {
                let updateResponse = await fetch(`${process.env.API_ORIGIN}/script/${this.$store.state.ide.scriptData.id}/publish`, {
                    method: 'PUT',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name: this.scriptName,
                        script: this.$store.state.ide.script
                    })
                });
                let {update} = await updateResponse.json();
            },
            async fork() {
                let forkResponse = await fetch(`${process.env.API_ORIGIN}/script/${this.$store.state.ide.scriptData.id}/version/${this.$store.state.ide.scriptData.mostRecent.createdAt}/fork`, {
                    method: 'POST',
                    credentials: 'include'
                });
                let {fork} = await forkResponse.json();
                this.$router.push({path: `/editor/${fork.id}`});
            },
            async duplicate() {
                let forkResponse = await fetch(`${process.env.API_ORIGIN}/script/${this.$store.state.ide.scriptData.id}/duplicate`, {
                    method: 'POST',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({name: this.scriptName})
                });
                let {fork} = await forkResponse.json();
                await fetch(`${process.env.API_ORIGIN}/script/${fork.id}`, {
                    method: 'PUT',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({script: this.$store.state.ide.script})
                })
                this.$router.push({path: `/editor/${fork.id}`});
            }
        }
    };
</script>
