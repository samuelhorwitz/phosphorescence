const nano = require('nano')('http://localhost:5984');
const rawDataFiles = require('require-all')(__dirname + '/raw');
// const rawTestFiles = require('require-all')(__dirname + '/test');
// const allFiles = Object.assign({}, rawDataFiles, rawTestFiles);
const allFiles = rawDataFiles;
(async () => {
    try {
        await nano.db.destroy('phosphorescence-training');
    } catch (e) {}
    await nano.db.create('phosphorescence-training');
    const trainingDatabase = nano.db.use('phosphorescence-training');
    const ratingCount = {};
    for (let filename in allFiles) {
        let rawDataFile = allFiles[filename];
        for (let [id, track] of Object.entries(rawDataFile)) {
            if (!id || !track || id == "null") {
                console.warn(`Skipping ${id}, seems to be invalid...`);
                continue;
            }
            if (ratingCount[id]) {
                let oldTrack = await trainingDatabase.get(id);
                track._rev = oldTrack._rev;
                track.track.aetherealness = runningAverage(ratingCount[id], oldTrack.track.aetherealness, track.track.aetherealness);
                track.track.primordialness = runningAverage(ratingCount[id], oldTrack.track.primordialness, track.track.primordialness);
                console.log(`Updating ${id}...`);
            } else {
                console.log(`Inserting ${id}...`);
            }
            try {
                await trainingDatabase.insert(track, id);
                if (!ratingCount[id]) {
                    ratingCount[id] = 1;
                } else {
                    ratingCount[id]++;
                }
            } catch (e) {
                console.warn(`Failed to insert ${id}: ${e.reason}`);
            }
        }
    }
    await trainingDatabase.insert(require(__dirname + '/track_couchdb_view.json'), '_design/tracks');
})();

function runningAverage(count, o, n) {
    return o + ((n - o) / (count + 1));
}
