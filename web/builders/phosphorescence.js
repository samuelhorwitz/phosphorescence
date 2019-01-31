// # Phosphorescence builder
// Phosphorescence's official playlist builder.

/* Creative Commons 0 Dedication
 *
 * This work is dedicated under the Creative Commons 0 dedication.
 * To the extent possible under law, the person who associated CC0 with this
 * work has waived all copyright and related or neighboring rights to this work.
 * https://creativecommons.org/publicdomain/zero/1.0/
 *
 * This is contrary to the majority of code in this repository which is licensed
 * under the MIT license with a retained copyright. Only files such as this one
 * which are explicitly licensed differently should be considered licensed under
 * the file-specific license described within. All other files are implicitly
 * licensed under the repository's MIT license.
 */

self.hooks.buildTree = function (kdTree, {points}) {
    return new kdTree(points, [AETHEREALNESS, PRIMORDIALNESS, KEY, MODE, TEMPO, ENERGY, VALENCE], (a, b) => {
        let valenceDiff = 1;
        if (a.valence && b.valence) {
            valenceDiff = b.valence - a.valence;
        }
        return calculateEuclidianDistance(
            b.aetherealness - a.aetherealness,
            b.primordialness - a.primordialness,
            b.energy - a.energy,
            valenceDiff,
            nonJarringHarmonicDifference(a, b),
            nonJarringTempoDifference(a, b)
        );
    });
}

self.hooks.getFirstTrack = function() {
    let neighbors = getNearestNeighbors(300, {
        aetherealness: Math.random(),
        primordialness: Math.random(),
        tempo: 125,
        energy: 0.45
    });
    let lowEnergyNeighbors = [];
    for (let neighbor of neighbors) {
        if (neighbor.point.tempo < 135) {
            lowEnergyNeighbors.push(neighbor);
        }
    }
    return pickRandom(lowEnergyNeighbors);
};

const LOOK_BACK = 4;
const EVOCATIVE_SWITCH_A_DROP = 0;
const EVOCATIVE_SWITCH_P_DROP = 1;
const EVOCATIVE_SWITCH_AP_DROP = 2;
const JARRING_KEY_CHANGE_DROP = 3;
// Not really a drop but double down on the energy and postpone the drop.
const SECOND_WIND_DROP = 4;

const runningValues = {
    valence: 0,
    energy: 0,
    tempo: 0,
    pitchFatigue: 0,
    vibeGoal: 4,
    lastDrop: 0
};

self.hooks.getNextTrack = function({tags, playlist, goalTracks}, previousTrack) {
    let nearingTheEnd = playlist.length > goalTracks - 1 - 4;
    let sickOfPitch = false;
    let drop = false;
    let dropStyle;
    if (!runningValues.vibeGoal && runningValues.lastDrop > 3 && runningValues.energy > 0.8) {
        console.log(`Track ${playlist.length} peaking!`);
        drop = rollDice(1, Math.ceil((1 - runningValues.energy) / 0.1));
    }
    if (!drop) {
        sickOfPitch = rollDice(runningValues.pitchFatigue, 7);
        if (sickOfPitch) {
            console.log(`Track ${playlist.length} sick of pitch!`);
        }
    }
    else {
        dropStyle = getRandomInt(0, 4);
        console.log(`Track ${playlist.length} should be a drop of style ${dropStyle}!`);
    }
    if (drop && dropStyle != SECOND_WIND_DROP) {
        runningValues.lastDrop = 0;
    }
    else {
        runningValues.lastDrop++;
    }
    let neighbors = getNearestNeighborsByTrack(300, previousTrack);
    /*
    (a, b) => {
        let aetherealnessDiff = b.aetherealness - a.aetherealness;
        let primordialnessDiff = b.primordialness - a.primordialness;
        let pitchDiff = calculateHarmonicDifference(a, b);
        let energyAndTempoDiff = 0;
        let pitchFatigueDiff = 0;
        if (runningValues.vibeGoal > 0) {
            if (b.energy < a.energy && b.tempo < a.tempo) {
                energyAndTempoDiff = Infinity;
            }
        }
        else if (runningValues.vibeGoal < 0) {
            if (a.energy < b.energy && a.tempo < b.tempo) {
                energyAndTempoDiff = Infinity;
            }
        }
        else if (drop) {
            switch (dropStyle) {
                case EVOCATIVE_SWITCH_A_DROP:
                aetherealnessDiff = 1 - aetherealnessDiff;
                if (b.energy > 0.9) {
                    energyAndTempoDiff = 10;
                }
                break;
                case EVOCATIVE_SWITCH_P_DROP:
                primordialnessDiff = 1 - primordialnessDiff;
                if (b.energy > 0.9) {
                    energyAndTempoDiff = 10;
                }
                break;
                case EVOCATIVE_SWITCH_AP_DROP:
                aetherealnessDiff = 1 - aetherealnessDiff;
                primordialnessDiff = 1 - primordialnessDiff;
                if (b.energy > 0.9) {
                    energyAndTempoDiff = 10;
                }
                break;
                case JARRING_KEY_CHANGE_DROP:
                pitchDiff = jarringHarmonicDifference(a, b)
                if (b.energy > 0.9) {
                    energyAndTempoDiff = 10;
                }
                break;
                case SECOND_WIND_DROP:
                if (b.energy < 0.7) {
                    energyAndTempoDiff = 10;
                }
                break;
            }
        }
        if (sickOfPitch && sameHarmonics(a, b)) {
            pitchFatigueDiff = 10;
        }
        return calculateEuclidianDistance(
            aetherealnessDiff,
            primordialnessDiff,
            pitchDiff,
            nonJarringTempoDifference(a, b),
            energyAndTempoDiff,
            pitchFatigueDiff
        );
    }
    let unseenNeighbors = cullTracksWithAlreadySeenTags(neighbors, tags);*/
    let neighbor = pickRandom(cullTracksWithAlreadySeenTags(neighbors, tags))
    runningValues.valence = calculateRMS(...getLastValues(playlist, 'valence', LOOK_BACK));
    runningValues.energy = calculateRMS(...getLastValues(playlist, 'energy', LOOK_BACK));
    runningValues.tempo = calculateRMS(...getLastValues(playlist, 'tempo', LOOK_BACK));
    if (runningValues.vibeGoal > 0) {
        runningValues.vibeGoal--;
    }
    else if (runningValues.vibeGoal < 0) {
        runningValues.vibeGoal++;
    }
    if (sickOfPitch) {
        runningValues.pitchFatigue = 0;
    }
    else if (sameHarmonics(neighbor, previousTrack)) {
        runningValues.pitchFatigue += 2;
    }
    else if (sameModeAndNonJarringKeyChange(neighbor, previousTrack) || differentModeAndNonJarringKeyChange(neighbor, previousTrack)) {
        runningValues.pitchFatigue += 1;
    }
    else {
        runningValues.pitchFatigue = 0;
    }
    return neighbor;
}

const edgyKeyUp = [G, D, A, E, B, F_SHARP, C_SHARP, G_SHARP, D_SHARP, A_SHARP, F, C];
const edgyKeyDown = [F, A_SHARP, D_SHARP, G_SHARP, C_SHARP, F_SHARP, B, E, A, D, G, C];

function calculateHarmonicDifference(a, b) {
    if (nonJarringHarmonicShift(a, b)) {
        return 0;
    }
    else if (a.mode == b.mode && (edgyKeyUp[a.key] == b.key || edgyKeyDown[a.key] == b.key)) {
        return 1 / 7;
    }
    return Infinity;
}

function jarringHarmonicDifference(a, b) {
    if (nonJarringHarmonicShift(a, b)) {
        return Infinity;
    }
    return 0;
}

function getLastValues(tracks, key, count) {
    return tracks.slice(-count).map(track => track.features[key]);
}
