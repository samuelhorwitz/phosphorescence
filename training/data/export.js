const fs = require('fs');
const dataStream = fs.createWriteStream(__dirname + '/csv/data.csv');
// const testStream = fs.createWriteStream(__dirname + '/csv/test.csv');
const nano = require('nano')('http://localhost:5984');
const window = require('svgdom');
const SVG = require('svg.js')(window);
const el1 = window.document.createElement('div');
// const el2 = window.document.createElement('div');

function shuffle(array) {
    let currentIndex = array.length, temporaryValue, randomIndex;
    // While there remain elements to shuffle...
    while (0 !== currentIndex) {
        // Pick a remaining element...
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex -= 1;

        // And swap it with the current element.
        temporaryValue = array[currentIndex];
        array[currentIndex] = array[randomIndex];
        array[randomIndex] = temporaryValue;
    }
    return array;
}

(async () => {
    const trainingDatabase = nano.db.use('phosphorescence-training');
    const drawSize = 500;
    const distribution = SVG(el1).size(drawSize, drawSize);
    // const scaledDistribution = SVG(el2).size(drawSize, drawSize);
    distribution.rect(drawSize, drawSize);
    // scaledDistribution.rect(drawSize, drawSize);
    const trainedWith = [];
    let body = await trainingDatabase.view('tracks', 'data', {
      'include_docs': false
    });
    dataStream.write('mode,time_signature,key_c,key_cs,key_d,key_ds,key_e,key_f,key_fs,key_g,key_gs,key_a,key_as,key_b,danceability,energy,loudness,speechiness,acousticness,instrumentalness,liveness,valence,tempo,aetherealness,primordialness\n');
    shuffle(body.rows).forEach((doc) => {
        console.log(`Export ${doc.value.name} - ${doc.value.artists.join(', ')}`);
        let {danceability, energy, key, loudness, mode, speechiness, acousticness, instrumentalness, liveness, valence, tempo, time_signature, aetherealness, primordialness} = doc.value;
        let keyC = key == 0 ? 1 : 0;
        let keyCs = key == 1 ? 1 : 0;
        let keyD = key == 2 ? 1 : 0;
        let keyDs = key == 3 ? 1 : 0;
        let keyE = key == 4 ? 1 : 0;
        let keyF = key == 5 ? 1 : 0;
        let keyFs = key == 6 ? 1 : 0;
        let keyG = key == 7 ? 1 : 0;
        let keyGs = key == 8 ? 1 : 0;
        let keyA = key == 9 ? 1 : 0;
        let keyAs = key == 10 ? 1 : 0;
        let keyB = key == 11 ? 1 : 0;
        dataStream.write(`${mode},${time_signature},${keyC},${keyCs},${keyD},${keyDs},${keyE},${keyF},${keyFs},${keyG},${keyGs},${keyA},${keyAs},${keyB},${danceability},${energy},${loudness},${speechiness},${acousticness},${instrumentalness},${liveness},${valence},${tempo},${aetherealness},${primordialness}\n`);
        distribution.circle(10).attr({ fill: '#f06' }).x(aetherealness * drawSize).y(primordialness * drawSize);
        // scaledDistribution.circle(10).attr({ fill: '#f06' }).x(aetherealness*drawSize).y(primordialness*drawSize);
        trainedWith.push(doc.id);
    });
    dataStream.end();
    // let testBody = await trainingDatabase.view('tracks', 'test', {
    //   'include_docs': false
    // });
    // testStream.write('danceability,energy,key,loudness,mode,speechiness,acousticness,instrumentalness,liveness,valence,tempo,time_signature\n');
    // testBody.rows.forEach((doc) => {
    //     console.log(`Export ${doc.value.name} - ${doc.value.artists.join(', ')}`);
    //     let {danceability, energy, key, loudness, mode, speechiness, acousticness, instrumentalness, liveness, valence, tempo, time_signature} = doc.value;
    //     testStream.write(`${danceability},${energy},${key},${loudness},${mode},${speechiness},${acousticness},${instrumentalness},${liveness},${valence},${tempo},${time_signature}\n`);
    // });
    // testStream.end();
    await new Promise(resolve => {
        fs.writeFile(__dirname + '/distribution.svg', distribution.svg(), resolve); 
    });
    // await new Promise(resolve => {
    //     fs.writeFile(__dirname + '/distribution_scaled.svg', scaledDistribution.svg(), resolve); 
    // });
    // await new Promise(resolve => {
    //     fs.writeFile(__dirname + '/trainedwith.json', JSON.stringify(trainedWith), resolve); 
    // });
})();
