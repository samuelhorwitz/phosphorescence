const poolSize = 20;

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

self.hooks.getNextTrack = function({tags}, previousTrack) {
    let neighbors = getNearestNeighborsByTrack(poolSize, previousTrack);
    return pickRandom(cullTracksWithAlreadySeenTags(neighbors, tags));
};
