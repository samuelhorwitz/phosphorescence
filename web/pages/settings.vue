<template>
    <article>
        <div>
            <h2>Settings</h2>
            <p>You are logged in as the Spotify user <a target="_blank" rel="external noopener" :href="'https://open.spotify.com/user/' + $store.state.user.user.spotifyId">{{$store.state.user.user.name}}</a>.</p>
            <p>Currently, all settings are stored locally on the device and automatically updated as they are changed.</p>
            <p>You may disconnect your Spotify account from this application at any time by visiting the <a target="_blank" rel="external noopener" href="https://www.spotify.com/us/account/apps/">Spotify user account page</a>.</p>
            <p>You may destroy every active session and log out of every device including this one: <button :disabled="logOutEverywhereClicked" @click="logoutEverywhere">Destroy All Sessions</button></p>
            <p>You may empty your cache on this device: <button :disabled="cacheClearState" @click="destroyCaches">{{clearCacheButtonText}}</button></p>
            <p>You may force a reload of the app in a context where you are unable to reload using your browser (PWA, etc). <button @click="reload">Reload</button></p>
        </div>
        <p>In order to use certain advanced features of Phosphorescence, you must authenticate your session. To do this, we will send an email to the address listed on your Spotify account. By clicking the button below, you agree to allow us to send you that single email. We will not save your email or use it for anything else. Once you receive the email, follow the enclosed directions.
        <button class="large" :disabled="authenticateClicked || alreadyAuthenticated" @click="sendAuthEmail">{{authButtonText}}</button>
        </p>
    </article>
</template>

<style scoped>
    article {
        margin: 0px 0.5em;
        font-size: 16px;
        background-color: teal;
        flex: 1;
        overflow-y: scroll;
        -webkit-overflow-scrolling: touch;
        border: 5px outset magenta;
        max-height: 100%;
        box-sizing: border-box;
        padding: 1em;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }

    h2 {
        color: white;
        font-size: 2em;
        margin: 0px;
        margin-bottom: 0.3em;
    }

    a {
        color: indigo;
    }

    a:hover {
        color: magenta;
    }

    button {
        border: 2px outset darkgray;
        background-color: gray;
        -webkit-appearance: none;
        color: black;
        cursor: pointer;
        padding: 0.3em;
    }

    button.large {
        border-width: 7px;
        width: 100%;
        margin: 0px;
        margin-top: 1em;
        font-size: 2em;
        padding: 0.5em 0;
    }

    button:hover, button:disabled {
        border-style: inset;
    }

    button:disabled {
        color: darkgray;
        cursor: not-allowed;
    }
</style>

<script>
    export default {
        data() {
            return {
                alreadyAuthenticated: this.$store.state.user.user.authenticated,
                authenticateClicked: false,
                logOutEverywhereClicked: false,
                logOutEverywhereFailed: false,
                emailSentAddress: "",
                cacheClearState: null
            };
        },
        computed: {
            authButtonText() {
                if (this.alreadyAuthenticated) {
                    return 'Already authenticated';
                }
                if (this.authenticateClicked) {
                    if (this.emailSentAddress) {
                        return `Email Sent to ${this.emailSentAddress}`;
                    }
                    return 'Sending email...';
                }
                return 'Authenticate Session';
            },
            logoutButtonText() {
                if (this.logOutEverywhereClicked) {
                    if (this.logOutEverywhereFailed) {
                        return 'Failed to destroy all sessions';
                    }
                    return 'Destroying sessions...';
                }
                return 'Destroy All Sessions';
            },
            clearCacheButtonText() {
                if (this.cacheClearState === 'clearing') {
                    return 'Clearing Cache...'
                } else if (this.cacheClearState === 'cleared') {
                    return 'Cache cleared'
                }
                return 'Clear Caches';
            }
        },
        methods: {
            async sendAuthEmail() {
                this.authenticateClicked = true;
                let emailResponse = await fetch(`${process.env.API_ORIGIN}/authenticate?utcOffsetMinutes=-${new Date().getTimezoneOffset()}`, {
                    method: 'POST',
                    credentials: 'include'
                });
                let {email} = await emailResponse.json();
                this.emailSentAddress = email;
            },
            async logoutEverywhere() {
                this.logOutEverywhereClicked = true;
                let {status} = await fetch(`${process.env.API_ORIGIN}/authenticate/logoutall`, {
                    method: 'POST',
                    credentials: 'include'
                });
                if (status === 200) {
                    window.location.href = '/';
                } else {
                    this.logOutEverywhereFailed = true;
                }
            },
            async destroyCaches() {
                this.cacheClearState = 'clearing';
                for (let cacheKey of await caches.keys()) {
                    await caches.delete(cacheKey);
                }
                this.cacheClearState = 'cleared';
            },
            reload() {
                location.reload(true);
            }
        }
    };
</script>