self.hooks.buildTree = function (kdTree, {points}) {
    return new kdTree(points, [AETHEREALNESS, PRIMORDIALNESS, KEY, MODE, ENERGY, TEMPO], (a, b) => 0);
}

self.hooks.getFirstTrack = function() {
    return pickRandom(getNodesWhere(({aetherealness, primordialness, energy, tempo}) => {
        return energy >= 0.75 && tempo >= 125 && aetherealness >= 0.5 && primordialness >= 0.6;
    }));
};

self.hooks.getNextTrack = null;
