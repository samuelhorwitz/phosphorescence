if (process.env.NODE_ENV !== 'production') {
    require('dotenv').config({path: require('path').resolve(process.cwd(), '../.env')});
}
const fs = require('fs');
const zlib = require('zlib');
const S3 = require('aws-sdk/clients/s3');
const s3 = new S3({
    endpoint: 'sfo2.digitaloceanspaces.com',
    regions: 'sfo2',
    accessKeyId: process.env.SPACES_ID,
    secretAccessKey: process.env.SPACES_SECRET
});

(async function() {
    await new Promise((resolve, reject) => {
        s3.copyObject({
            Bucket: 'phosphorescence',
            CopySource: '/phosphorescence/tracks.json',
            Key: `old/tracks-${new Date().getTime()}.json`
        }, function(err, data) {
            if (err) {
                reject(err);
            }
            else {
                resolve();
            }
        });
    });
    await new Promise((resolve, reject) => {
        s3.upload({
            Bucket: 'phosphorescence',
            Key: 'tracks.json',
            Body: fs.createReadStream(`${__dirname}${process.env.TRACKS_JSON}`).pipe(zlib.createGzip()),
            ContentType: 'application/json',
            ContentEncoding: 'gzip',
            ACL: 'private'
        }, function(err, data) {
            if (err) {
                reject(err);
            }
            else {
                resolve();
            }
        });
    })
})();