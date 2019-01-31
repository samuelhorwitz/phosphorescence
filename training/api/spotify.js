import express from 'express';
import axios from 'axios';
import querystring from 'querystring';
const {clientId, secret} = require('../spotify_credentials.json');
const app = express();

app.get('/authorize', (req, res) => {
    res.redirect('https://accounts.spotify.com/authorize?' + querystring.stringify({
        client_id: clientId,
        response_type: 'code',
        scope: ['streaming', 'user-read-birthdate', 'user-read-email', 'user-read-private'].join(' '),
        redirect_uri: 'http://localhost:3000/tokens'
    }));
});

app.get('/tokens', async (req, res) => {
    try { 
        let body;
        if (req.query.type == 'refresh') {
            body = {
                grant_type: 'refresh_token',
                refresh_token: req.query.code
            };
        } else {
            body = {
                grant_type: 'authorization_code',
                code: req.query.code,
                redirect_uri: 'http://localhost:3000/tokens'
            };
        }
        body = Object.assign(body, {
            client_id: clientId,
            client_secret: secret
        });
        let {status, statusText, data} = await axios.post('https://accounts.spotify.com/api/token', querystring.stringify(body));
        res.setHeader('Content-Type', 'application/json');
        res.send(JSON.stringify({access: data.access_token, refresh: data.refresh_token, expires: data.expires_in}));
    } catch ({response, request, message}) {
        if (response) {
            console.error(response.data);
            console.error(response.status);
            console.error(response.headers);
        } else if (request) {
            console.error(request);
        } else {
            console.error(message);
        }
        res.sendStatus(500);
    }
});

export default {
    path: '/api/spotify',
    handler: app
};
