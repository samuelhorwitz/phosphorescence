self.hooks.getFirstTrack = function({points}) {
    [AETHEREALNESS, PRIMORDIALNESS, KEY, MODE, ENERGY, TEMPO].forEach(addLoggingDimension);
    return buildResponse(buildPoint(pickRandom(points.filter(({aetherealness, primordialness, energy, tempo}) => {
        return energy >= 0.75 && tempo >= 125 && aetherealness >= 0.5 && primordialness >= 0.6;
    }))));
};
