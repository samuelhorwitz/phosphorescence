self.hooks.getFirstTrack = function({points}) {
    [AETHEREALNESS, PRIMORDIALNESS, KEY, MODE, ENERGY, TEMPO].forEach(addLoggingDimension);
    return buildResponse(buildPoint(pickRandom(points.filter(({aetherealness, primordialness, energy, tempo}) => {
        return energy <= 0.7 && tempo <= 130 && aetherealness >= 0.65 && primordialness <= 0.65;
    }))));
};
