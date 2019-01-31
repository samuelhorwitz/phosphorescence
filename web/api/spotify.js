const express = require('express');
const axios = require('axios');
const querystring = require('querystring');
const S3 = require('aws-sdk/clients/s3');
const s3 = new S3({
    endpoint: 'sfo2.digitaloceanspaces.com',
    regions: 'sfo2',
    accessKeyId: process.env.SPACES_ID,
    secretAccessKey: process.env.SPACES_SECRET
});

module.exports = function(app) {
    app.get('/api/spotify/authorize', (req, res) => {
        res.redirect('https://accounts.spotify.com/authorize?' + querystring.stringify({
            client_id: process.env.SPOTIFY_CLIENT_ID,
            response_type: 'code',
            scope: ['streaming', 'user-read-birthdate', 'user-read-email', 'user-read-private', 'user-read-playback-state', 'user-read-recently-played'].join(' '),
            redirect_uri: `${process.env.PHOSPHOR_ORIGIN}/auth/login`
        }));
    });

    app.get('/api/spotify/tokens', async (req, res) => {
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
                    redirect_uri: `${process.env.PHOSPHOR_ORIGIN}/auth/login`
                };
            }
            body = Object.assign(body, {
                client_id: process.env.SPOTIFY_CLIENT_ID,
                client_secret: process.env.SPOTIFY_SECRET
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

    app.get('/api/spotify/tracks', async (req, res) => {
        let spotifyToken = req.query.token;
        if (spotifyToken) {
            try {
                let {status} = await axios.get('https://api.spotify.com/v1/me', {headers: {Authorization: `Bearer ${spotifyToken}`}});
                if (status == 200) {
                    let redirectTo;
                    if (process.env.NODE_ENV === 'production') {
                        redirectTo = s3.getSignedUrl('getObject', {
                            Bucket: 'phosphorescence',
                            Key: 'tracks.json',
                            Expires: 60 * 2
                        });
                    }
                    else {
                        redirectTo = '/tracks.json';
                    }
                    res.redirect(302, redirectTo);
                }
                else {
                    res.sendStatus(403);
                }
            }
            catch ({response}) {
                res.sendStatus(response.status);
            }
        }
        else {
            res.sendStatus(403);
        }
    });

    return app;
};
