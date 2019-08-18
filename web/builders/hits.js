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

// Prune functions are called before tree building and track generation. They
// may also be called alone for builders that only act as pruners, in a chain of
// builder calls. They are responsible for reducing the tracks available. They
// are only allowed to return a subset of the passed in `tracks`; attempting to
// add tracks will result in an error. The `unprunedTracks` value is passed in
// solely for analysis, such as here, where we want the average popularity of
// all tracks in order to filter.
self.hooks.prune = function ({tracks, unprunedTracks}) {
    addLoggingDimension(POPULARITY);
    let totalPopularity = 0;
    let unprunedTracksArr = Object.values(unprunedTracks);
    unprunedTracksArr.forEach(({track}) => totalPopularity += track.popularity);
    let avgPopularity = totalPopularity / unprunedTracksArr.length;
    let popularTracks = {};
    Object.values(tracks).forEach(trackWrapper => {
        let {track} = trackWrapper;
        if (track.popularity > avgPopularity) {
            popularTracks[track.id] = trackWrapper;
        }
    });
    return buildResponse(popularTracks);
}
