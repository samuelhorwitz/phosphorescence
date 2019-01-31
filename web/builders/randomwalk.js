// # Random Walk builder
// This playlist builder performs a random walk by finding "close" neighbors and
// and randomly selecting from them to choose the next track.

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

// Prior to executing user scripts, the code is wrapped in a closure, so you may
// put as many variables and functions in your script as you like without
// worrying about namespace pollution. Regardless, every user script is run in a
// single-use web worker anyway.

// First we define a simple _k_-d tree by setting our dimensions and our
// distance function. We will choose a random measure of evocativeness as well
// as use the supplied harmonic and tempo functions to ensure we lean towards
// non-jarring harmonic shifts and BPM changes.
self.hooks.buildTree = function (kdTree, {points}) {
    return new kdTree(points, [AETHEREALNESS, PRIMORDIALNESS, KEY, MODE, TEMPO], (a, b) => {
        return calculateEuclidianDistance(
            b.aetherealness - a.aetherealness,
            b.primordialness - a.primordialness,
            nonJarringHarmonicDifference(a, b),
            nonJarringTempoDifference(a, b)
        );
    });
}

// Choose a completely random first track.
self.hooks.getFirstTrack = function() {
    return getRandomTrack();
};

// All subsequent tracks are picked by feeding the last chosen track into the
// nearest neighbor function as the ideal point and searching for neighbors that
// are close-by in evocativeness as well as non-jarring in tempo and harmonic
// changes. Prior to picking a random track, we also make sure to clear out all
// tracks that may not be exact duplicates but are likely remixes or different
// releases of an already chosen track.
self.hooks.getNextTrack = function({tags}, previousTrack) {
    let neighbors = getNearestNeighborsByTrack(treeSize() * 0.01, previousTrack);
    return pickRandom(cullTracksWithAlreadySeenTags(neighbors, tags));
};
