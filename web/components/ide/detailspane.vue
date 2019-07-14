<template>
    <aside>
        <section>
            <div class="tableWrapper" ref="tableWrapper">
                <table>
                    <thead>
                        <tr>
                            <th class="number">
                                #
                            </th>
                            <th>
                                Title
                            </th>
                            <th>
                                Artist
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(track, index) in $store.state.idetracks.playlist" @click="inspect(track)" :class="{selected: !!inspectedTrack && (track.track.id == inspectedTrack.track.id)}">
                            <td class="number" @mouseenter="hoverIndex = index" @mouseleave="hoverIndex = null">
                                <span v-if="hoverIndex == index"><svg @click="play(index, track)" xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title>Play Track</title><path d="M3,0.25V31.71L30.25,16ZM5,3.71L26.25,16,5,28.24V3.71Z"></path></svg></span>
                                <span>{{index + 1}}</span>
                            </td>
                            <td><a target="_blank" :href="track.track.external_urls.spotify">{{track.track.name}}</a></td>
                            <td>
                                <ol>
                                    <li v-for="artist in track.track.artists">
                                        <a target="_blank" :href="artist.external_urls.spotify">{{artist.name}}</a>
                                    </li>
                                </ol>
                            </td>
                        </tr>
                        <tr v-if="!$store.state.idetracks.playlist || $store.state.idetracks.playlist.length == 0">
                            <td colspan="3" class="noPlaylist">No playlist generated</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </section>
        <section class="monacoWrapper" :class="{isOpen: isOpen}">
            <details :open="isOpen">
                <summary>
                    <div class="summary" v-on:click.stop.prevent>
                        <button class="disclosure" @click="isOpen = true" :disabled="!inspectedTrack">Click To See Details</button>
                        <div v-if="inspectedTrack">
                            <span>
                                <a target="_blank" :href="inspectedTrack.track.external_urls.spotify">{{inspectedTrack.track.name}}</a> -
                                <ol>
                                    <li v-for="artist in inspectedTrack.track.artists">
                                        <a target="_blank" :href="artist.external_urls.spotify">{{artist.name}}</a>
                                    </li>
                                </ol> -
                                <a target="_blank" :href="inspectedTrack.track.album.external_urls.spotify">{{inspectedTrack.track.album.name}}</a>
                            </span>
                            <dl>
                                <template v-for="(val, dim) in dimensionsList">
                                    <dt>{{dim}}</dt>
                                    <dd>{{val}}</dd>
                                </template>
                            </dl>
                        </div>
                        <div v-if="!inspectedTrack">
                            Please select a track.
                        </div>
                    </div>
                    <button class="closer disclosure" @click="isOpen = false">Collapse</button>
                </summary>
                <div ref="monaco" class="monacoContainer"></div>
            </details>
        </section>
        <section class="runner">
            <div class="runnerButtons">
                <button @click="runBuilder" :disabled="isBuilding">Test Builder</button>
                <button @click="terminateBuilder" v-show="isBuilding">Terminate Builder</button>
            </div>
            <progress v-show="isBuilding"></progress>
            <div v-if="lastError" class="buildError">
                {{lastError}}
            </div>
            <div v-if="isBuilding && isLongRunning" class="buildError">
                The script is taking a while to run...
            </div>
            <div v-if="!isBuilding && isLongRunning" class="buildError">
                The script took a while to run, this might indicate an issue.
            </div>
            <div>
                <label>
                    Count
                    <input type="number" v-model="numberOfTracks">
                </label>
            </div>
        </section>
        <section class="player">
            <player></player>
        </section>
    </aside>
</template>

<style scoped>
    .tableWrapper {
        max-height: 400px;
        overflow-y: scroll;
    }

    aside {
        display: flex;
        flex-direction: column;
        background-color: rgb(200, 200, 200);
        border-left: 1px solid black;
        position: relative;
        width: 100%;
    }

    ol {
        list-style-type: none;
        margin: 0px;
        padding: 0px;
        display: inline;
    }

    li {
        display: inline;
    }

    li:not(:last-child)::after {
        content: ', ';
        display: inline;
    }

    table {
        table-layout: fixed;
        width: 100%;
        height: 100%;
        border-spacing: 0px;
        border-collapse: unset;
        color: black;
    }

    thead tr:nth-child(1) th{
        background: gray;
        position: sticky;
        top: 0;
        z-index: 10;
    }

    th {
        background-color: rgba(193, 193, 193);
        border: 3px outset rgb(170, 170, 170);
    }

    th.number {
        width: 3em;
    }

    tr {
        color: black;
        background-color: white;
    }

    td {
        padding: 0.5em;
    }

    td:first-child:not(.noPlaylist) {
        display: flex;
    }

    td:first-child span {
        flex: 1;
    }

    tr:hover {
        cursor: pointer;
    }

    tbody tr:hover, tbody tr.selected {
        background-color: rgb(0, 14, 118);
        color: white;
    }

    tr:hover a, tr.selected a {
        color: white;
    }

    td:not(.number) {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    td.number {
        text-align: right;
    }

    a {
        text-decoration: none;
    }

    a:hover {
        text-decoration: underline;
    }

    .monacoWrapper {
        display: flex;
        flex: 1;
        align-items: flex-start;
        justify-content: center;
        position: relative;
        max-height: 375px;
        border-top: 1px solid black;
    }

    .monacoWrapper:not(.isOpen) {
        overflow-y: scroll;
    }

    .monacoContainer {
        position: absolute;
        width: 100%;
        top: 27.5px;
        bottom: 0px;
    }

    details {
        display: flex;
        flex: 1;
        list-style: none;
    }

    summary::-webkit-details-marker {
        display: none;
    }

    summary {
        outline: 0;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    details:not([open]) .closer {
        display: none;
    }

    details:not([open]) .summary {
        display: inline-block;
    }

    details[open] .closer {
        display: inline;
    }

    details[open] .summary {
        display: none;
    }

    .summary {
        flex: 1;
        padding: 1em;
        cursor: default;
    }

    button {
        font-weight: bold;
        appearance: none;
        margin: 0px;
        height: 2.5em;
        border: 1px outset gray;
        background-color: gray;
        margin-bottom: 1em;
        cursor: pointer;
        outline: 0px;
    }

    .closer {
        width: 100%;
        margin: 0px;
    }

    dl {
        cursor: default;
    }

    dt {
        cursor: text;
    }

    dd {
        font-family: courier;
        background-color: white;
        border: 3px inset gray;
        cursor: text;
    }

    .runner {
        padding: 1em;
        border-top: 1px solid black;
    }

    button[disabled] {
        color: white;
        background-color: lightgray;
        border-color: lightgray;
        border-style: inset;
        cursor: not-allowed;
    }

    .runnerButtons {
        display: flex;
    }

    .buildError {
        color: red;
        font-family: courier;
    }

    .player {
        display: flex;
        justify-content: center;
        margin-top: auto;
    }

    .number svg {
        max-height: 0.8em;
        fill: white;
        stroke: white;
    }

    svg:hover {
        fill: magenta;
    }

    .noPlaylist {
        display: table-cell;
    }
</style>

<script>
    import {getUsersCountry} from '~/assets/session';
    import {initialize, loadNewPlaylist} from '~/assets/recordcrate';
    import {terminatePlaylistBuilding} from '~/assets/eos';
    import player from '~/components/ide/player';
    import * as monaco from 'monaco-editor';

    const keysWithSharps = ['C', 'C♯', 'D', 'D♯', 'E', 'F', 'F♯', 'G', 'G♯', 'A', 'A♯', 'B'];
    const keysWithFlats = ['C', 'D♭', 'D', 'E♭', 'E', 'F', 'G♭', 'G', 'A♭', 'A', 'B♭', 'B'];

    export default {
        components: {
            player
        },
        data() {
            return {
                relayoutFn: null,
                editor: null,
                inspectedTrack: null,
                isOpen: false,
                isBuilding: false,
                lastError: null,
                isLongRunning: false,
                hoverIndex: null
            };
        },
        computed: {
            dimensionsList() {
                if (!this.inspectedTrack) {
                    return null;
                }
                let dimMap = {};
                for (let dim of this.$store.state.ide.dimensions) {
                    let name = dim;
                    let postfix = '';
                    let val = this.inspectedTrack.features[dim];
                    if (dim === 'tempo') {
                        dimMap.BPM = this.inspectedTrack.features.tempo;
                    }
                    else if (dim === 'key' || dim === 'mode') {
                        if (!dimMap['Harmonic']) {
                            let key = keysWithSharps[this.inspectedTrack.features.key];
                            let mode = this.inspectedTrack.features.mode == 1 ? 'major' : 'minor';
                            dimMap.Harmonic = `${key} ${mode}`;
                        }
                    }
                    else if (dim === 'timeSignature') {
                        dimMap['Time Signature'] = `${this.inspectedTrack.features['time_signature']} / 4`;
                    }
                    else if (dim === 'duration') {
                        dimMap.Duration = this.inspectedTrack.features.duration_ms;
                    }
                    else if (dim === 'primordialness' || dim === 'aetherealness') {
                        dimMap[dim.charAt(0).toUpperCase() + dim.slice(1)] = this.inspectedTrack.evocativeness[dim];
                    }
                    else {
                        dimMap[dim.charAt(0).toUpperCase() + dim.slice(1)] = this.inspectedTrack.features[dim];
                    }
                }
                return dimMap;
            },
            numberOfTracks: {
                get() {
                    return this.$store.state.ide.numberOfTracks;
                },
                set(newValue) {
                    this.$store.commit('ide/numberOfTracks', newValue);
                }
            }
        },
        watch: {
            inspectedTrack(track) {
                this.editor.setValue(JSON.stringify(track, null, 4));
            },
            isOpen(open) {
                if (open) {
                    this.$nextTick(() => {
                        this.updateDimensions();
                    });
                }
            }
        },
        async created() {
            await initialize(await getUsersCountry());
            this.$store.commit('idetracks/restore');
        },
        mounted() {
            this.editor = monaco.editor.create(this.$refs.monaco, {
                value: '',
                language: 'json',
                readOnly: true
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
            inspect(track) {
                this.inspectedTrack = track;
            },
            async runBuilder() {
                let playlist, dimensions;
                this.isBuilding = true;
                this.lastError = null;
                this.isLongRunning = false;
                let timeoutId = setTimeout(() => {
                    this.isLongRunning = true;
                }, 10000);
                try {
                    let response = await loadNewPlaylist(this.$store.state.ide.numberOfTracks, this.$store.state.ide.script);
                    if (!response.error) {
                        playlist = response.playlist;
                        dimensions = response.dimensions;
                    }
                    else {
                        this.lastError = response.error;
                        this.isBuilding = false;
                        clearTimeout(timeoutId);
                        return;
                    }
                }
                catch (e) {
                    this.lastError = e.message;
                    this.isBuilding = false;
                    clearTimeout(timeoutId);
                    return;
                }
                clearTimeout(timeoutId);
                this.$store.commit('ide/dimensions', dimensions);
                this.$store.dispatch('idetracks/loadPlaylist', JSON.parse(JSON.stringify(playlist)));
                this.inspectedTrack = playlist[0];
                this.isBuilding = false;
            },
            terminateBuilder() {
                terminatePlaylistBuilding();
            },
            play(i, track) {
                this.inspect(track);
                this.$store.dispatch('idetracks/seekTrack', i);
                this.$store.dispatch('idetracks/play');
            }
        }
    };
</script>