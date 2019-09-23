const express = require('express');
const app = express();

app.use(function (req, res, next) {
    res.set({
        'Content-Security-Policy': `default-src 'none';child-src 'self';script-src 'self' 'unsafe-eval' blob:;style-src 'self' 'unsafe-inline';worker-src 'self';connect-src blob:;base-uri 'none';form-action 'none';frame-ancestors ${process.env.PHOSPHOR_ORIGIN};block-all-mixed-content;navigate-to 'none';sandbox allow-scripts allow-same-origin;`
    })
    next();
});

if (process.env.NODE_ENV === 'production') {
    const http = require('http');
    app.use(express.static('./public', {setHeaders}));
    const server = http.createServer(app).listen('80', '0.0.0.0', () => {
        console.log('Server listening on `' + server.address().address + ':' + server.address().port + '`.');
    });
}
else {
    const host = process.env.HOST || '127.0.0.1';
    const port = process.env.PORT || 3001;
    const https = require('https');
    const fs = require('fs');
    const privateKey  = fs.readFileSync(__dirname + '/eos.localhost.key', 'utf8');
    const certificate = fs.readFileSync(__dirname + '/eos.localhost.crt', 'utf8');
    const credentials = {key: privateKey, cert: certificate};
    app.use(express.static('./dist', {setHeaders}));
    const server = https.createServer(credentials, app).listen(port, host, () => {
        console.log('Server listening on `' + server.address().address + ':' + server.address().port + '`.');
    });
}

function setHeaders(res, path) {
    if (path.match('worker')) {
        res.setHeader('Cache-Control', 'public, max-age=31536000');
    }
}
