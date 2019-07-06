const express = require('express');
const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();
const axios = require('axios');
const uuidv4 = require('uuid/v4');
const uuidv5 = require('uuid/v5');
const zlib = require('zlib');
const S3 = require('aws-sdk/clients/s3');
const s3 = new S3({
    endpoint: 'nyc3.digitaloceanspaces.com',
    regions: 'nyc3',
    accessKeyId: process.env.SPACES_ID,
    secretAccessKey: process.env.SPACES_SECRET
});
const db = require('knex')({
  client: 'pg',
  connection: process.env.PG_CONNECTION_STRING
});
const phosphorNamespace = uuidv5('phosphor.me', uuidv5.DNS);
const scriptsNamespace = uuidv5('scripts', phosphorNamespace);

module.exports = function(app) {
    app.get('/api/scripts/my', async (req, res) => {
        let spotifyToken = req.cookies['spotify_access'];
        if (!spotifyToken) {
            res.sendStatus(403);
            return;
        }
        let spotifyId;
        try {
            let {status, data} = await axios.get('https://api.spotify.com/v1/me', {headers: {Authorization: `Bearer ${spotifyToken}`}});
            if (status != 200) {
                res.sendStatus(403);
                return;
            }
            spotifyId = data.id;
        }
        catch ({response}) {
            res.sendStatus(response.status);
            return;
        }
        let scripts;
        try {
            scripts = await db.select()
                .column('scripts.id', 'scripts.is_private')
                .from('scripts')
                .innerJoin('users', 'users.id', 'scripts.author_id')
                .where('users.spotify_id', spotifyId)
                .whereNull('scripts.deleted_at');
        }
        catch (e) {
            res.sendStatus(500);
            return;
        }
        res.send(JSON.stringify(scripts));
    });

    app.post('/api/scripts', jsonParser, async (req, res) => {
        let spotifyToken = req.cookies['spotify_access'];
        if (!spotifyToken) {
            res.sendStatus(403);
            return;
        }
        if (!req.body.script || !req.body.type) {
            res.sendStatus(400);
            return;
        }
        let spotifyId;
        try {
            let {status, data} = await axios.get('https://api.spotify.com/v1/me', {headers: {Authorization: `Bearer ${spotifyToken}`}});
            if (status != 200) {
                res.sendStatus(403);
                return;
            }
            spotifyId = data.id;
        }
        catch ({response}) {
            res.sendStatus(response.status);
            return;
        }
        let scriptId = uuidv5(req.body.script, scriptsNamespace);
        let existsOnS3;
        try {
            existsOnS3 = await new Promise((resolve, reject) => {
                s3.headObject({
                    Bucket: 'phosphorescence-scripts',
                    Key: scriptId
                }, (err, data) => {
                    if (err) {
                        if (err.code === 'NotFound') {
                            resolve(false);
                            return;
                        }
                        reject(err);
                        return;
                    }
                    resolve(true);
                });
            });
        }
        catch (e) {
            res.sendStatus(500);
            return;
        }
        if (!existsOnS3) {
            try {
                await new Promise((resolve, reject) => {
                    s3.putObject({
                        Bucket: 'phosphorescence-scripts',
                        ACL: 'public-read',
                        ContentType: 'application/javascript',
                        ContentEncoding: 'gzip',
                        Key: scriptId,
                        Body: zlib.gzipSync(req.body.script)
                    }, (err, data) => {
                        if (err) {
                            reject(err);
                            return;
                        }
                        resolve(data);
                    });
                });
            }
            catch (e) {
                res.sendStatus(500);
                return;
            }
        }
        let scriptInternalId;
        try {
            await db.transaction(async trx => {
                let users = await trx.select()
                    .from('users')
                    .where('spotify_id', spotifyId);
                let userId;
                if (users.length === 0) {
                    userId = uuidv4();
                    await trx('users').insert({
                        id: userId,
                        spotify_id: spotifyId
                    });
                }
                else {
                    userId = users[0].id;
                }
                scriptInternalId = uuidv4();
                await trx('scripts').insert({
                    id: scriptInternalId,
                    author_id: userId
                });
                await trx('script_versions').insert({
                    script_id: scriptInternalId,
                    type: req.body.type,
                    file_id: scriptId
                });
            });
        }
        catch (e) {
            res.sendStatus(500);
            return;
        }
        res.send({id: scriptInternalId, mostRecentVersion: scriptId});
    });

    return app;
};
