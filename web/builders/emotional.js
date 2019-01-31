self.hooks.buildTree = function (kdTree, {points}) {
    return new kdTree(points, [AETHEREALNESS, PRIMORDIALNESS, KEY, MODE, ENERGY, TEMPO], (a, b) => 0);
}

self.hooks.getFirstTrack = function() {
    return pickRandom(getNodesWhere(({aetherealness, primordialness, energy, tempo}) => {
        return energy >= 0.85 && tempo >= 130 && aetherealness >= 0.6 && primordialness <= 0.4;
    }));
};

self.hooks.getNextTrack = null;
