<template>
    <div class="pageWrapper">
        <article class="authPage">
            <h2 class="pageHeader">Settings</h2>
            <p>You are logged in as the Spotify user <a target="_blank" :href="'https://open.spotify.com/user/' + $store.state.user.user.spotifyId">{{$store.state.user.user.name}}</a>.</p>
            <p>Currently, all settings are stored locally on the device and automatically updated as they are changed.</p>
            <p>You may disconnect your Spotify account from this application at any time by visiting the <a target="_blank" href="https://www.spotify.com/us/account/apps/">Spotify user account page</a>.</p>
            <p>You may destroy every active session and log out of every device including this one: <button :disabled="logOutEverywhereClicked" @click="logoutEverywhere">Destroy All Sessions</button></p>
            <p>In order to use certain advanced features of Phosphorescence, you must authenticate your session. To do this, we will send an email to the address listed on your Spotify account. By clicking the button below, you agree to allow us to send you that single email. We will not save your email or use it for anything else. Once you receive the email, follow the enclosed directions.</p>
            <button class="large" :disabled="authenticateClicked || alreadyAuthenticated" @click="sendAuthEmail">{{authButtonText}}</button>
        </article>
    </div>
</template>

<style scoped>
    h2 {
        margin-bottom: 0.3em;
    }

    .pageWrapper {
        margin-left: 2em;
        margin-right: 2em;
    }

    .authPage {
        margin: 0px;
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
        margin-top: 1em;
        margin-bottom: 1em;
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
        async fetch({store, error}) {
            let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
            if (!userResponse.ok) {
                return error({statusCode: userResponse.status, message: "Could not get user information"});
            }
            let {user} = await userResponse.json();
            store.commit('user/user', user);
        },
        data() {
            return {
                alreadyAuthenticated: this.$store.state.user.user.authenticated,
                authenticateClicked: false,
                logOutEverywhereClicked: false,
                logOutEverywhereFailed: false,
                emailSentAddress: ""
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
            }
        }
    };
</script>