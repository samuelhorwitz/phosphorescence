self.hooks.buildTree = function (kdTree, {points}) {
    return new kdTree(points, [AETHEREALNESS, PRIMORDIALNESS, KEY, MODE, TEMPO], (a, b) => 0);
}

self.hooks.getFirstTrack = function() {
    return pickRandom(getNodesWhere(({aetherealness, primordialness, energy, tempo}) => {
        return tempo <= 140 && aetherealness <= 0.6 && primordialness <= 0.5;
    }));
};

self.hooks.getNextTrack = null;
