<template>
    <div class="wrapper">
        <div class="monacoContainer" ref="monaco"></div>
    </div>
</template>

<style scoped>
    .wrapper {
        display: flex;
        flex: 1;
        align-items: center;
        justify-content: center;
        position: relative;
    }

    .monacoContainer {
        position: absolute;
        width: 100%;
        height: 100%;
    }
</style>

<script>
    import * as monaco from 'monaco-editor';
    import basewalker from '!raw-loader!~/assets/ide/basewalker.js';
    import api from '!raw-loader!~/eos/api.js';

    export default {
        data() {
            return {
                editor: null,
                relayoutFn: null
            };
        },
        mounted() {
            let initialValue = this.$store.state.ide.script;
            if (!initialValue) {
                initialValue = basewalker;
            }
            this.editor = monaco.editor.create(this.$refs.monaco, {
                value: initialValue,
                language: 'javascript'
            });
            this.saveState();
            this.editor.onDidChangeModelContent(() => {
                this.saveState();
            });
            this.relayoutFn = () => this.updateDimensions();
            addEventListener('resize', this.relayoutFn);
        },
        beforeDestroy() {
            this.editor.dispose();
            removeEventListener('resize', this.relayoutFn);
        },
        methods: {
            updateDimensions() {
                this.editor.layout();
            },
            saveState() {
                this.$store.commit('ide/save', this.editor.getValue());
            }
        }
    };
</script>