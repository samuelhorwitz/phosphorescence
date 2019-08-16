self.hooks.buildTree = function (kdTree, {points}) {
    return new kdTree(points, [AETHEREALNESS, PRIMORDIALNESS], (a, b) => {
        return calculateEuclidianDistance(
            b.aetherealness - a.aetherealness,
            b.primordialness - a.primordialness
        );
    });
}

self.hooks.getFirstTrack = function() {
    return getRandomTrack();
};

self.hooks.getNextTrack = function({playlist, goalTracks, tags}) {
    let firstTrack = playlist[0];
    let neighbors = getNearestNeighborsByTrack(goalTracks * 10, firstTrack);
    let culledNeighbors = cullTracksWithAlreadySeenTags(neighbors, tags);
    return culledNeighbors[playlist.length - 1];
};
