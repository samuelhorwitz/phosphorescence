// # Only the hits pruner
// This playlist builder prunes unpopular tracks as per the Spotify popularity
// ranking.

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

self.hooks.prune = function ({tracks}) {
    addLoggingDimension(ENERGY, LOUDNESS, TEMPO);
    let sleepTracks = {};
    Object.values(tracks).forEach(trackWrapper => {
        let {track, features} = trackWrapper;
        let {energy, loudness, tempo} = features;
        if (tempo <= 130 && energy <= 0.7 && loudness <= -10) {
            sleepTracks[track.id] = trackWrapper;
        }
    });
    return buildResponse(sleepTracks);
}
