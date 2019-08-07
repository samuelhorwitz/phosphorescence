<template>
    <div class="pageWrapper">
        <article class="authPage">
            <h2 class="pageHeader">Settings</h2>
            <p>You are logged in as the Spotify user <a target="_blank" :href="'https://open.spotify.com/user/' + $store.state.user.user.spotifyId">{{$store.state.user.user.name}}</a>.</p>
            <p>Currently, all settings are stored locally on the device and automatically updated as they are changed.</p>
            <p>You may disconnect your Spotify account from this application at any time by visiting the <a target="_blank" href="https://www.spotify.com/us/account/apps/">Spotify user account page</a>.</p>
            <p>In order to use certain advanced features of Phosphorescence, you must authenticate your session. To do this, we will send an email to the address listed on your Spotify account. By clicking the button below, you agree to allow us to send you that single email. We will not save your email or use it for anything else. Once you receive the email, follow the enclosed directions.</p>
            <button :disabled="authenticateClicked || alreadyAuthenticated" @click="sendAuthEmail">{{buttonText}}</button>
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
        border: 7px outset darkgray;
        background-color: gray;
        -webkit-appearance: none;
        width: 100%;
        font-size: 2em;
        color: black;
        padding: 0.5em 0;
        cursor: pointer;
        margin-top: 1em;
        margin-bottom: 1em;
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
                emailSentAddress: ""
            };
        },
        computed: {
            buttonText() {
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
            }
        }
    };
</script>