<template>
    <aside>
        <div class="container">
            <div class="wrapper" :class="{deviceMenuVisible: devicesMenu}">
                <menu>
                    <li v-if="playerReadyAndConnected"><button @click="play" v-if="$store.getters['tracks/stopped']" :class="{disabled: !$store.getters['tracks/playlistLoaded']}">
                        <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Play Track</title><path d="M3,0.25V31.71L30.25,16ZM5,3.71L26.25,16,5,28.24V3.71Z"></path></svg>
                    </button></li>
                    <li v-if="playerReadyAndConnected"><button @click="resume" v-if="$store.getters['tracks/paused']" :class="{disabled: !$store.getters['tracks/playlistLoaded']}">
                        <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Play Track</title><path d="M3,0.25V31.71L30.25,16ZM5,3.71L26.25,16,5,28.24V3.71Z"></path></svg>
                    </button></li>
                    <li v-if="playerReadyAndConnected"><button @click="pause" v-if="$store.getters['tracks/playing']" :class="{disabled: !$store.getters['tracks/playlistLoaded']}">
                        <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Pause Track</title><path d="M13.76,1H6.63V31h7.14V1Zm-2,28H8.63V3h3.14V29Z"></path><path d="M25.37,1H18.24V31h7.14V1Zm-2,28H20.24V3h3.14V29Z"></path></svg>
                    </button></li>
                    <li v-if="playerReadyAndConnected"><button @click="previous" :class="{disabled: !canSeek || !$store.getters['tracks/canSkipBackward']}">
                        <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Previous Track</title><path d="M3.46,2h-2V30h2V16.08L30.54,31.71V0.25L3.46,15.89V2ZM28.54,3.72V28.25L7.3,16Z"></path></svg>
                    </button></li>
                    <li v-if="playerReadyAndConnected"><button @click="next" :class="{disabled: !canSeek || !$store.getters['tracks/canSkipForward']}">
                        <svg xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 32 32" x="0px" y="0px" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Next Track</title><path d="M28.54,15.89L1.46,0.25V31.71L28.54,16.08V30h2V2h-2v13.9ZM3.46,28.25V3.72L24.7,16Z"></path></svg>
                    </button></li>
                    <li v-if="premiumUser" class="devicePicker">
                        <button @click="toggleDevicePicker">
                            <span v-if="!activeDevice">Choose a device</span>
                            <span v-if="activeDevice && !hideActiveDevice">{{playingOnText}} <em>{{activeDevice.name}}</em></span>
                            <span v-if="hideActiveDevice">Loading...</span>
                        </button>
                    </li>
                    <li v-if="!premiumUser" class="devicePicker savePlaylist">
                        <button @click="createPlaylist" :disabled="savePlaylistSuccessful" :class="{success: savePlaylistSuccessful, failure: savePlaylistFailed}">
                            <div class="deviceIcon phosphorLogo">+</div>
                            <span>{{savePlaylistButtonText}}</span>
                        </button>
                    </li>
                </menu>
                <div class="trackData" ref="trackData" :class="{stopped: $store.getters['tracks/stopped']}">
                    <span class="trackDetails" :class="{scrollingBanner: isTrackDataScrolling}" ref="trackDetails" v-if="playerReadyAndConnected && !$store.getters['tracks/stopped']">
                        <span class="trackName">
                            <a target="_blank" rel="external noopener" :href="currentTrackUrl">{{currentTrackName}}</a>
                        </span>
                        <ol class="artistsNames">
                            <li v-for="artist in currentTrackArtists">
                                <a target="_blank" rel="external noopener" :href="artist.external_urls.spotify">{{artist.name}}</a>
                            </li>
                        </ol>
                        <span class="albumName">
                            <a target="_blank" rel="external noopener" :href="currentAlbumUrl">{{currentAlbumName}}</a>
                        </span>
                    </span>
                    <span class="nothingPlaying scrollingBanner" v-if="playerReadyAndConnected && $store.getters['tracks/stopped'] && !$store.state.tracks.spotifyAppearsDown">
                        ... Welcome to Phosphorescence ... Please Click "Play" To Listen ... 💿 💻 ... Hint: You can drag and drop a track from Spotify to seed the playlist builder ...
                    </span>
                    <span class="nothingPlaying scrollingBanner" v-if="playerReadyAndConnected && $store.getters['tracks/stopped'] && $store.state.tracks.spotifyAppearsDown">
                        ... 😢😢😢 Spotify's Playback API Appears To Be Down Right Now 😢 ...
                    </span>
                </div>
            </div>
        </div>
        <div v-if="devicesMenu" class="devicesMenu">
            <ol>
                <li v-for="device in devices" :class="{active: device.active}" @click="transferDevice(device.id)" v-if="device.type != 'Phosphorescence' || device.isPrimary">
                    <button class="deviceButton">
                        <div v-if="device.type == 'Phosphorescence'" class="deviceIcon phosphorLogo">P</div>
                        <div v-if="device.type == 'Computer'" class="deviceIcon"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0px" y="0px" viewBox="0 0 64 64" xml:space="preserve" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Select {{device.name}}</title><g><g><g><path d="M39.703,39.322H7.936c-5.02,0-6.035-0.821-6.035-4.885V13.483c0-4.226,1.251-4.4,6.035-4.4h31.768     c4.784,0,6.035,0.175,6.035,4.4v21.44C45.738,39.147,44.487,39.322,39.703,39.322z M7.936,11.083c-1.532,0-3.269,0-3.745,0.322     c-0.243,0.163-0.29,0.886-0.29,2.078v20.954c0,2.776,0,2.885,4.035,2.885h31.768c1.532,0,3.269,0,3.746-0.322     c0.242-0.163,0.29-0.885,0.29-2.076v-21.44c0-1.192-0.047-1.915-0.29-2.078c-0.477-0.322-2.213-0.322-3.746-0.322H7.936z"></path></g><g><path d="M31.378,43.978H16.261l1.605-6.647h11.905L31.378,43.978z M18.801,41.978h10.036l-0.64-2.647H19.44L18.801,41.978z"></path></g></g><g><path d="M46.212,54.917H2.236l5.502-10.389h32.973L46.212,54.917z M5.559,52.917H42.89l-3.382-6.389H8.943L5.559,52.917z"></path></g><g><path d="M58.412,53.536c-0.279,0-0.583-0.006-0.914-0.012c-0.382-0.007-0.802-0.015-1.262-0.015h-3.89    c-3.972,0-5.862-0.748-5.862-5.001V14.675c0-4.333,2.01-4.535,5.862-4.535h3.89c3.853,0,5.863,0.202,5.863,4.535V49.6    C62.1,52.927,60.858,53.536,58.412,53.536z M52.347,12.14c-3.754,0-3.862,0.071-3.862,2.535v33.834    c0,2.548,0.286,3.001,3.862,3.001h3.89c0.474,0,0.905,0.008,1.299,0.015c0.318,0.006,0.609,0.012,0.877,0.012    c1.624,0,1.688,0,1.688-1.937V14.675c0-2.464-0.108-2.535-3.863-2.535H52.347z"></path></g><g><path d="M55.423,39.413c0,0.683-0.495,1.233-1.107,1.233l0,0c-0.61,0-1.107-0.551-1.107-1.233l0,0    c0-0.682,0.497-1.233,1.107-1.233l0,0C54.928,38.18,55.423,38.731,55.423,39.413L55.423,39.413z"></path></g><g><path d="M55.423,44.525c0,0.683-0.495,1.235-1.107,1.235l0,0c-0.61,0-1.107-0.553-1.107-1.235l0,0c0-0.68,0.497-1.232,1.107-1.232    l0,0C54.928,43.293,55.423,43.846,55.423,44.525L55.423,44.525z"></path></g><g><rect x="50.352" y="16.598" width="7.884" height="2"></rect></g><g><rect x="50.352" y="20.688" width="7.884" height="2"></rect></g><g><rect x="50.352" y="24.778" width="7.884" height="2"></rect></g></g></svg></div>
                        <div v-if="device.type == 'Smartphone'" class="deviceIcon"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0px" y="0px" viewBox="0 0 64 64" xml:space="preserve" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Select {{device.name}}</title><g><path d="M42.152,61.572H21.849c-6.655,0-6.655-2.627-6.655-10.426L15.192,7.534c-0.009-0.711-0.032-2.377,1.209-3.635   c0.963-0.977,2.457-1.472,4.438-1.472h22.323c1.98,0,3.474,0.495,4.437,1.472c1.241,1.258,1.219,2.924,1.209,3.635l-0.001,0.142   v43.471C48.807,58.945,48.807,61.572,42.152,61.572z M20.839,4.428c-1.426,0-2.44,0.295-3.014,0.876   c-0.654,0.663-0.641,1.627-0.633,2.203l0.001,0.169v43.471c0,8.11,0.174,8.426,4.655,8.426h20.304c4.48,0,4.654-0.315,4.654-8.426   l0.002-43.64c0.008-0.576,0.021-1.54-0.634-2.203c-0.573-0.581-1.587-0.876-3.013-0.876H20.839z"></path><path d="M48.807,49.602H15.193V11.976c0-0.959,0.034-1.773,0.104-2.488l0.087-0.903h33.228l0.09,0.9   c0.07,0.704,0.104,1.519,0.104,2.491V49.602z M17.193,47.602h29.613V11.976c0-0.511-0.01-0.973-0.03-1.392H17.223   c-0.02,0.423-0.03,0.884-0.03,1.392V47.602z"></path><path d="M32,58.174c-1.953,0-3.541-1.589-3.541-3.541c0-1.951,1.588-3.539,3.541-3.539c1.952,0,3.541,1.588,3.541,3.539   C35.541,56.585,33.952,58.174,32,58.174z M32,53.094c-0.85,0-1.541,0.69-1.541,1.539c0,0.85,0.691,1.541,1.541,1.541   s1.541-0.691,1.541-1.541C33.541,53.784,32.85,53.094,32,53.094z"></path></g></svg></div>
                        <div v-if="device.type != 'Smartphone' && device.type != 'Computer' && device.type != 'Phosphorescence'" class="deviceIcon"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.0" x="0px" y="0px" viewBox="0 0 24 24" xml:space="preserve" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Select {{device.name}}</title><path d="M7,9.7l15-2V4c0-1.2-1.1-2.1-2.3-2l-11,1.5c-1,0.1-1.7,1-1.7,2V9.7z"></path><path d="M9,18c0,1.7-1.6,3.4-3.5,3.9S2,21.4,2,19.7s1.6-3.4,3.5-3.9C7.4,15.4,9,16.3,9,18z"></path><line stroke-width="2" stroke-miterlimit="10" x1="8" y1="18" x2="8" y2="6"></line><path d="M22,17c0,1.7-1.6,3.4-3.5,3.9S15,20.4,15,18.7s1.6-3.4,3.5-3.9S22,15.3,22,17z"></path><line stroke-width="2" stroke-miterlimit="10" x1="21" y1="17" x2="21" y2="4"></line></svg></div>
                        <div class="deviceName">{{device.name}} <span v-if="device.type == 'Phosphorescence' && device.isPrimary">(This Window)</span></div>
                    </button>
                </li>
                <li @click="createPlaylist">
                    <button :disabled="savePlaylistSuccessful" class="deviceButton" :class="{success: savePlaylistSuccessful, failure: savePlaylistFailed}">
                        <div class="deviceIcon phosphorLogo">+</div>
                        <div class="deviceName">{{savePlaylistButtonText}}</div>
                    </button>
                </li>
                <li @click="refreshDevices">
                    <button class="deviceButton">
                        <div class="deviceIcon"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0px" y="0px" viewBox="0 0 100 100" xml:space="preserve" aria-labelledby="uniqueTitleID" role="img"><title id="uniqueTitleID">Refresh Devices</title><g><path d="M29.455,58.993c-3.851-8.647-2.062-18.514,4.554-25.136c6.619-6.616,16.492-8.405,25.146-4.556   c0.082,0.036,0.166,0.051,0.248,0.082l-4.238,3.327c-1.77,1.388-2.08,3.944-0.691,5.715c0.803,1.022,1.998,1.559,3.205,1.559   c0.883,0,1.766-0.284,2.512-0.868l10.744-8.428c1.219-0.955,1.789-2.518,1.473-4.033l-2.797-13.415   c-0.457-2.199-2.609-3.612-4.814-3.154c-2.201,0.458-3.615,2.614-3.156,4.816l1.291,6.197c-0.035-0.018-0.064-0.042-0.102-0.058   c-12.1-5.383-25.92-2.864-35.217,6.423c-9.285,9.292-11.805,23.112-6.42,35.209c0.749,1.683,2.401,2.683,4.134,2.683   c0.614,0,1.239-0.127,1.837-0.392C29.446,63.947,30.469,61.274,29.455,58.993z"></path><path d="M78.814,37.026c-1.012-2.283-3.686-3.31-5.969-2.296c-2.283,1.012-3.311,3.685-2.295,5.967   c3.844,8.656,2.057,18.523-4.561,25.138c-6.482,6.482-16.081,8.321-24.601,4.774l4.231-3.317c1.767-1.388,2.079-3.947,0.688-5.718   c-1.387-1.767-3.946-2.079-5.714-0.688l-10.746,8.428c-1.218,0.955-1.79,2.518-1.473,4.031l2.796,13.413   C31.57,88.68,33.262,90,35.15,90c0.274,0,0.555-0.028,0.833-0.084c2.2-0.461,3.615-2.615,3.157-4.817l-1.285-6.167   c4.023,1.685,8.218,2.517,12.367,2.517c8.159,0,16.122-3.178,22.167-9.226C81.67,62.948,84.193,49.128,78.814,37.026z"></path></g></svg></div>
                        <div class="deviceName">Refresh Devices</div>
                    </button>
                </li>
            </ol>
        </div>
    </aside>
</template>

<style scoped>
    aside {
        grid-column: 1 / 3;
        margin-top: 1em;
        margin-left: 2em;
        margin-right: 2em;
        position: relative;
        min-width: 0;
    }

    menu {
        margin: 0px;
        cursor: pointer;
        padding-left: .5em;
        display: flex;
        align-items: center;
    }

    menu li {
        display: flex;
    }

    .devicesMenu {
        z-index: 10000;
        position: absolute;
        bottom: 2.7em;
        background-color: rgba(26, 17, 16, 1);
        color: white;
        min-width: 15em;
        display: flex;
        border-top-left-radius: 5px;
        border-top-right-radius: 5px;
        border: 5px teal inset;
        border-bottom: 0px;
        border-bottom-right-radius: 5px;
        font-size: 16px
    }

    .devicesMenu ol {
        display: flex;
        flex: 1;
        flex-direction: column;
        list-style-type: none;
        margin: 0px;
        padding: 0px;
    }

    .devicesMenu ol li {
        display: flex;
        flex: 1;
        align-items: center;
        justify-content: flex-start;
        padding: 1em;
        cursor: pointer;
        user-select: none;
    }

    .devicesMenu ol li:hover {
        background-color: aqua;
    }

    .devicesMenu ol li:hover .deviceButton {
        color: blue;
    }

    .devicesMenu ol li:hover .deviceButton .deviceIcon svg {
        fill: blue;
        stroke: blue;
    }

    .devicesMenu ol li:hover .deviceButton .deviceIcon.phosphorLogo {
        color: blue;
        border-color: blue;
    }

    .container {
        text-align: center;
        display: flex;
        justify-content: center;
        height: 3em;
        font-size: 16px;
    }

    .wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: rgba(26, 17, 16, 1);
        border: 5px inset teal;
        border-radius: 10px;
        width: 100%;
        position: relative;
        box-sizing: border-box;
    }

    .wrapper.deviceMenuVisible {
        border-top-left-radius: 0px;
    }

    .trackName::after, .artistsNames::after {
        content: ' - ';
        color: white;
    }

    .nothingPlaying {
        color: white;
        cursor: pointer;
    }

    .trackDetails {
        display: inline-block;
    }

    .scrollingBanner {
        animation: marquee 15s linear infinite;
        padding-left: 100%;
        display: inline-block;
    }

    .scrollingBanner:hover {
        animation-play-state: paused;
    }

    .devicePicker button {
        color: white;
        cursor: pointer;
        padding-left: 0.5em;
    }

    .devicePicker:hover button span {
        color: aqua;
        text-decoration: underline;
    }

    .devicePicker em {
        font-style: normal;
    }

    .devicePicker.savePlaylist button {
        display: flex;
        flex-direction: row;
    }

    .devicePicker.savePlaylist button span {
        margin-left: 1em;
    }

    .devicePicker.savePlaylist button.success span,
    .devicePicker.savePlaylist:hover button.success span {
        color: limegreen;
    }

    .devicePicker.savePlaylist button.success .deviceIcon.phosphorLogo {
        color: limegreen;
        border: 1px solid limegreen;
    }

    .devicePicker.savePlaylist button.failure span,
    .devicePicker.savePlaylist:hover button.failure span {
        color: indianred;
    }

    .devicePicker.savePlaylist button.failure .deviceIcon.phosphorLogo {
        color: indianred;
        border: 1px solid indianred;
    }

    .deviceButton {
        color: white;
        cursor: pointer;
        display: flex;
    }

    .deviceButton .deviceIcon {
        width: 1.5em;
        height: 1.5em;
        margin-right: 1em;
    }

    .deviceButton .deviceIcon svg {
        fill: white;
        stroke: white;
    }

    .devicePicker.savePlaylist button .deviceIcon.phosphorLogo,
    .deviceButton .deviceIcon.phosphorLogo {
        color: white;
        font-family: 'Varela';
        text-transform: uppercase;
        font-weight: bold;
        border: 1px solid white;
        width: 1.25em;
        height: 1.25em;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .deviceButton.success {
        color: limegreen;
    }

    .deviceButton.success .deviceIcon svg {
        fill: limegreen;
        stroke: limegreen;
    }

    .deviceButton.success .deviceIcon.phosphorLogo {
        color: limegreen;
        border: 1px solid limegreen;
    }

    .deviceButton.failure {
        color: indianred;
    }

    .deviceButton.failure .deviceIcon svg {
        fill: indianred;
        stroke: indianred;
    }

    .deviceButton.failure .deviceIcon.phosphorLogo {
        color: indianred;
        border: 1px solid indianred;
    }

    @keyframes marquee {
        0%   { transform: translate(0, 0); }
        100% { transform: translate(-100%, 0); }
    }

    a {
        color: white;
        text-decoration: none;
    }

    a:hover {
        text-decoration: underline;
    }

    ol.artistsNames {
        list-style-type: none;
        margin: 0px;
        padding: 0px;
        display: inline;
    }

    ol.artistsNames li {
        display: inline;
    }

    ol.artistsNames li:not(:last-child)::after {
        content: ', ';
        display: inline;
        color: white;
    }

    .trackData {
        font-family: VT323;
        font-size: 1.3em;
        flex: 1;
        white-space: nowrap;
        overflow-x: scroll;
        overflow-y: hidden;
        padding-right: 0.5em;
        padding-bottom: 0.1em;
        position: relative;
    }

    .trackData.stopped {
        overflow: hidden;
    }

    .trackData::-webkit-scrollbar { 
        display: none; 
    }

    .trackData::before {
        content: '';
        background: linear-gradient(90deg, rgba(26,17,16,0.9) 46%, rgba(255,255,255,0) 100%);
        width: 5px;
        position: absolute;
        left: 0px;
        top: 5px;
        height: 70%;
    }

    .trackData::after {
        content: '';
        background: linear-gradient(270deg, rgba(26,17,16,0.9) 46%, rgba(255,255,255,0) 100%);
        width: 5px;
        position: absolute;
        top: 5px;
        right: 0px;
        height: 70%;
    }

    button {
        appearance: none;
        border: 0px;
        background-color: transparent;
        margin: 0px;
        padding: 0px;
        padding-right: 0.25em;
        outline: none;
    }

    menu button svg {
        width: 2em;
        fill: white;
        stroke: white;
        cursor: pointer;
        stroke-linejoin: round;
    }

    menu button:not(.disabled) svg:hover {
        fill: aquamarine;
        stroke: magenta;
    }

    button.disabled svg {
        fill: gray;
        stroke: gray;
    }

    button.disabled svg:hover {
        cursor: not-allowed;
    }

    @media only screen and (max-height: 449px) {
        aside {
            grid-column: 1 / 3;
            grid-row: 3 / 4;
            justify-content: center;
            margin: 0 1em;
        }

        .wrapper {
            flex: 1;
        }
    }

    @media only screen and (max-height: 249px) {
        aside {
            grid-row: 2 / 4;
        }
    }

    @media only screen and (min-height: 450px) and (max-width: 1099px) {
        aside {
            grid-column: 1 / 2;
        }
    }
</style>

<script>
    import {initializePlayer} from '~/assets/spotify';

    export default {
        data() {
            return {
                webPlayerReady: false,
                destroyer: null,
                devicesMenu: false,
                devicesLoaded: false,
                hideActiveDevice: false,
                isTrackDataScrolling: false,
                devices: [],
                savePlaylistState: null
            };
        },
        computed: {
            tracks() {
                return this.$store.state.tracks.playlist;
            },
            currentTrack() {
                return this.$store.getters['tracks/currentTrack'];
            },
            currentTrackName() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.name;
            },
            currentTrackUrl() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.external_urls.spotify;
            },
            currentAlbumName() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.album.name;
            },
            currentAlbumUrl() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.album.external_urls.spotify;
            },
            currentTrackArtists() {
                let track = this.currentTrack;
                if (!track) {
                    return null;
                }
                return track.track.artists;
            },
            canSeek() {
                return this.playerReadyAndConnected && !this.$store.state.tracks.neighborSeekLocked;
            },
            playingOnText() {
                if (this.$store.getters['tracks/stopped']) {
                    return 'Will play on';
                }
                return 'Playing on';
            },
            playerReadyAndConnected() {
                return this.webPlayerReady && this.$store.getters['tracks/isPlayerConnected'];
            },
            activeDevice() {
                for (let device of this.devices) {
                    if (device.active) {
                        return device;
                    }
                }
                return null;
            },
            savePlaylistSuccessful() {
                return this.savePlaylistState === 'SUCCESS';
            },
            savePlaylistFailed() {
                return this.savePlaylistState === 'FAILED';
            },
            savePlaylistButtonText() {
                if (this.savePlaylistSuccessful) {
                    return 'Playlist Saved';
                } else if (this.savePlaylistFailed) {
                    return 'Failed To Save';
                }
                return 'Save Playlist';
            },
            premiumUser() {
                return this.$store.state.user.user.product === 'premium';
            }
        },
        watch: {
            currentTrackName() {
                this.checkShouldScroll();
            }
        },
        methods: {
            async play() {
                let wasStopped = this.$store.getters['tracks/stopped'];
                await this.$store.dispatch('tracks/play');
                if (wasStopped) {
                    this.checkShouldScroll();
                }
            },
            resume() {
                this.$store.dispatch('tracks/resume');
            },
            pause() {
                this.$store.dispatch('tracks/pause');
            },
            next() {
                this.$store.dispatch('tracks/next');
            },
            previous() {
                this.$store.dispatch('tracks/previous');
            },
            async toggleDevicePicker() {
                this.savePlaylistState = null;
                this.devicesMenu = !this.devicesMenu;
                if (this.devicesLoaded) {
                    return;
                }
                if (this.devicesMenu) {
                    await this.refreshDevices();
                }
                this.devicesLoaded = true;
            },
            async refreshDevices() {
                let devicesResponse = await fetch(`${process.env.API_ORIGIN}/devices`, {credentials: 'include'});
                let {devices} = await devicesResponse.json();
                if (!devices) {
                    return;
                }
                this.setDevices(devices);
            },
            setDevices(devices) {
                let primaryDevice;
                let allDevices = [];
                for (let device of devices) {
                    let d = {
                        id: device.id,
                        name: device.name,
                        type: device.name == 'Phosphorescence' ? 'Phosphorescence' : device.type,
                        active: device.is_active
                    }
                    if (device.id == this.$store.state.tracks.deviceId) {
                        d.type = 'Phosphorescence';
                        d.isPrimary = true;
                        primaryDevice = d;
                    } else {
                        allDevices.push(d);
                    }
                }
                if (primaryDevice) {
                    allDevices.unshift(primaryDevice);
                }
                this.devices = allDevices;
            },
            async transferDevice(deviceId) {
                this.hideActiveDevice = true;
                let shouldPlayAfterTransfer = true;
                let playState = 'paused';
                if (this.$store.getters['tracks/playing'] || this.$store.getters['tracks/paused']) {
                    this.$store.commit('tracks/play');
                    playState = 'play';
                    shouldPlayAfterTransfer = false;
                }
                let devicesResponse = await fetch(`${process.env.API_ORIGIN}/device/${deviceId}?playState=play`, {
                    method: 'PUT',
                    credentials: 'include'
                });
                if (shouldPlayAfterTransfer) {
                    this.$store.dispatch('tracks/play');
                }
                let {devices} = await devicesResponse.json();
                this.setDevices(devices);
                this.hideActiveDevice = false;
                this.devicesMenu = false;
                this.checkShouldScroll();
            },
            checkShouldScroll() {
                setTimeout(() => {
                    if (!(this.$refs.trackData && this.$refs.trackDetails)) {
                        return;
                    }
                    if (!this.isTrackDataScrolling) {
                        if (this.$refs.trackData.clientWidth > this.$refs.trackDetails.clientWidth) {
                            return;
                        }
                        this.isTrackDataScrolling = true;
                    } else {
                        if (this.$refs.trackDetails.clientWidth - this.$refs.trackData.clientWidth > this.$refs.trackData.clientWidth) {
                            return;
                        }
                        this.isTrackDataScrolling = false;
                    }
                });
            },
            async createPlaylist() {
                let tracks = this.tracks.map(track => {
                    let uri = track.track.uri;
                    if (track.track.linked_from) {
                        uri = track.track.linked_from.uri;
                    }
                    return {
                        name: track.track.name,
                        uri
                    };
                });
                let savePlaylistResponse = await fetch(`${process.env.API_ORIGIN}/users/me/playlist`, {
                    method: 'POST',
                    credentials: 'include',
                    body: JSON.stringify({
                        tracks,
                        utcOffsetMinutes: -(new Date().getTimezoneOffset())
                    })
                });
                if (savePlaylistResponse.ok) {
                    this.savePlaylistState = 'SUCCESS';
                } else {
                    this.savePlaylistState = 'FAILED';
                }
            },
            handleKeyPress(e) {
                if (!this.webPlayerReady) {
                    return;
                }
                if (e.code === 'Space') {
                    e.stopPropagation();
                    e.preventDefault();
                    if (this.$store.getters['tracks/stopped']) {
                        this.play();
                    } else if (this.$store.getters['tracks/paused']) {
                        this.resume();
                    } else if (this.$store.getters['tracks/playing']) {
                        this.pause();
                    }
                }
                // left arrow
                else if (e.keyCode === 37 && (e.metaKey || e.ctrlKey)) {
                    e.stopPropagation();
                    e.preventDefault();
                    this.previous();
                }
                // right arrow
                else if (e.keyCode == 39 && (e.metaKey || e.ctrlKey)) {
                    e.stopPropagation();
                    e.preventDefault();
                    this.next();
                }
            }
        },
        mounted() {
            addEventListener('resize', this.checkShouldScroll);
            document.addEventListener('keydown', this.handleKeyPress);
        },
        async created() {
            if (!this.premiumUser) {
                console.debug('Free Spotify user; disabling web player');
                return;
            }
            this.$store.commit('loading/startLoad');
            let messageId = await this.$store.dispatch('loading/pushMessage', 'Initializing Spotify web player');
            this.$store.commit('loading/initializeProgress', {id: 'player', weight: 5});
            try {
                let playerWrapper = await initializePlayer(this.$store, 'tracks');
                this.destroyer = playerWrapper.destroyer;
                this.$store.dispatch('tracks/registerPlayer', playerWrapper.player);
                this.devices = [{
                    name: this.$store.state.tracks.deviceName,
                    id: this.$store.state.tracks.deviceId,
                    type: 'Phosphorescence',
                    active: true
                }];
                this.webPlayerReady = true;
            }
            catch (e) {
                this.webPlayerReady = false;
            }
            this.checkShouldScroll();
            this.$store.commit('loading/completeProgress', {id: 'player'});
            this.$store.commit('loading/clearMessage', messageId);
            this.$store.dispatch('loading/endLoadAfterDelay');
        },
        beforeDestroy() {
            if (this.destroyer) {
                this.destroyer();
            }
            removeEventListener('resize', this.checkShouldScroll);
            document.removeEventListener('keydown', this.handleKeyPress);
        }
    };
</script>