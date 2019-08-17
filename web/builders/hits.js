// # Only the hits builder
// This playlist builder is similar to the random walk standard builder
// but will only select "popular" tracks (as per Spotify rankings).

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

let avgPopularity;

self.hooks.buildTree = function (kdTree, {points}) {
    let totalPopularity = 0;
    points.forEach(point => totalPopularity += point.popularity);
    avgPopularity = totalPopularity / points.length;
    addLoggingDimension(POPULARITY);
    return new kdTree(points, [AETHEREALNESS, PRIMORDIALNESS, KEY, MODE, TEMPO], (a, b) => {
        return calculateEuclidianDistance(
            b.aetherealness - a.aetherealness,
            b.primordialness - a.primordialness,
            nonJarringHarmonicDifference(a, b),
            nonJarringTempoDifference(a, b)
        );
    });
}

self.hooks.getFirstTrack = function() {
    return pickRandom(getNodesWhere(({popularity}) => popularity > avgPopularity));
};

self.hooks.getNextTrack = function({tags}, previousTrack) {
    let neighbors = getNearestNeighborsByTrack(treeSize() * 0.005, previousTrack);
    return pickRandom(cullUnpopularTracks(cullTracksWithAlreadySeenTags(neighbors, tags)));
};

function cullUnpopularTracks(neighbors) {
    let unculledNeighbors = [];
    for (let nn of neighbors) {
        if (nn.point.popularity <= avgPopularity) {
            continue;
        }
        unculledNeighbors.push(nn);
    }
    return unculledNeighbors;
}
