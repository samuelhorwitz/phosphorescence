const express = require('express');
const app = express();

app.use(function (req, res, next) {
    let unsafeEval = '';
    // Webpack dev mode uses eval
    if (process.env.NODE_ENV !== 'production') {
        unsafeEval = `'unsafe-eval'`;
    }
    res.set({
        'Content-Security-Policy': `default-src 'none';script-src ${process.env.EOS_ORIGIN} ${unsafeEval} blob:;worker-src blob:;connect-src blob:;base-uri 'none';form-action 'none';frame-ancestors ${process.env.PHOSPHOR_ORIGIN};block-all-mixed-content;navigate-to 'none';sandbox allow-scripts;`
    })
    next();
});

if (process.env.NODE_ENV === 'production') {
    const http = require('http');
    app.use(express.static('./public'));
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
    app.use(express.static('./dist'));
    const server = https.createServer(credentials, app).listen(port, host, () => {
        console.log('Server listening on `' + server.address().address + ':' + server.address().port + '`.');
    });
}
