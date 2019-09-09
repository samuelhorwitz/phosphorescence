// # "Hard" seeder

/* Creative Commons 0 Dedication
 *
 * This work is dedicated under the Creative Commons 0 dedication.
 * To the extent possible under law, the person who associated CC0 with this
 * work has waived all copyright and related or neighboring rights to this work.
 * https://creativecommons.org/publicdomain/zero/1.0/
 *
 * This is contrary to the majority of code in this repository which is licensed
 * under the Apache 2.0 license with a retained copyright. Only files such as this one
 * which are explicitly licensed differently should be considered licensed under
 * the file-specific license described within. All other files are implicitly
 * licensed under the repository's Apache 2.0 license.
 */

self.hooks.getFirstTrack = function({points}) {
    [AETHEREALNESS, PRIMORDIALNESS, ENERGY, TEMPO].forEach(addLoggingDimension);
    return buildResponse(buildPoint(pickRandom(points.filter(({aetherealness, primordialness, energy, tempo}) => {
        return energy >= 0.8 && tempo >= 135 && aetherealness <= 0.5 && primordialness <= 0.6;
    }))));
};
