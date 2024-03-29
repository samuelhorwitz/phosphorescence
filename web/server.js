const express = require('express');
const app = express();

app.use(function (req, res, next) {
    res.set({
        'Content-Security-Policy': `default-src 'none';media-src https://p.scdn.co;manifest-src 'self';child-src 'self';script-src 'self' 'unsafe-eval' 'unsafe-inline' https://sdk.scdn.co https://www.google.com https://www.gstatic.com https://www.google-analytics.com https://connect.facebook.net https://static.ads-twitter.com https://analytics.twitter.com;img-src 'self' https://www.google-analytics.com https://www.facebook.com https://t.co https://i.scdn.co data:;style-src 'self' 'unsafe-inline' https://fonts.gstatic.com https://fonts.googleapis.com;font-src 'self' https://fonts.gstatic.com;frame-src 'self' https://accounts.spotify.com https://open.spotify.com https://sdk.scdn.co https://www.google.com https://developers.google.com ${process.env.EOS_ORIGIN} ${process.env.API_ORIGIN};connect-src 'self' https://phosphorescence-tracks.nyc3.digitaloceanspaces.com https://api.spotify.com https://www.google-analytics.com ${process.env.API_ORIGIN} ${process.env.SCRIPTS_ORIGIN};worker-src 'self';base-uri 'none';form-action 'none';frame-ancestors 'self';block-all-mixed-content;navigate-to 'self' ${process.env.API_ORIGIN} https://accounts.spotify.com;`
    })
    next();
});

if (process.env.NODE_ENV === 'production') {
    const http = require('http');
    app.use(express.static('./public'));
    app.use('/editor/:id?', express.static('./public'));
    const server = http.createServer(app).listen('80', '0.0.0.0', () => {
        console.log('Server listening on `' + server.address().address + ':' + server.address().port + '`.');
    });
}
else {
    const { Nuxt, Builder } = require('nuxt');
    const host = process.env.HOST || '127.0.0.1';
    const port = process.env.PORT || 3000;
    const config = require('./nuxt.config.js');
    config.dev = true;
    const nuxt = new Nuxt(config);
    app.use(nuxt.render);
    new Builder(nuxt).build().then(() => {
        const https = require('https');
        const fs = require('fs');
        const privateKey  = fs.readFileSync(__dirname + '/phosphor.localhost.key', 'utf8');
        const certificate = fs.readFileSync(__dirname + '/phosphor.localhost.crt', 'utf8');
        const credentials = {key: privateKey, cert: certificate};
        const server = https.createServer(credentials, app).listen(port, host, () => {
            console.log('Server listening on `' + server.address().address + ':' + server.address().port + '`.');
        });
    });
}
