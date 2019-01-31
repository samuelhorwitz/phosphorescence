self.hooks.buildTree = function (kdTree, {points}) {
    return new kdTree(points, [AETHEREALNESS, PRIMORDIALNESS, KEY, MODE, ENERGY, TEMPO], (a, b) => 0);
}

self.hooks.getFirstTrack = function() {
    return pickRandom(getNodesWhere(({aetherealness, primordialness, energy, tempo}) => {
        return energy <= 0.9 && tempo <= 140 && aetherealness >= 0.7 && primordialness <= 0.3;
    }));
};

self.hooks.getNextTrack = null;
