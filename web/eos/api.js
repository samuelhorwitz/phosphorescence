// # Phosphorescence API
// The Phosphorescence API allows users to create their own playlist building
// functionality right in their browser. Phosphorescence does all of the heavy
// lifting including building out a space-partitioned tree of thousands of
// tracks which can be used to find "nearest neighbor" tracks using your own
// custom distance functions. Here, we will detail what global constructs are
// available to you to use in your scripts.

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

// ## Legal
// Please note that if you develop custom scripts on our website, you are
// relinquishing these scripts into the public domain as per
// [CC0](https://creativecommons.org/share-your-work/public-domain/cc0/). You
// may clone the project from Github and run your own instance instead,
// if you do not wish to release your work into the public domain.
//
// ---

// ### Public constants

// Spotify's [mode](https://developer.spotify.com/documentation/web-api/reference/tracks/get-audio-features/)
// is a numeric value representing "major" or "minor".
Object.defineProperties(self, {
    MINOR: {
        value: 0,
        writable: false,
        configurable: false
    },
    MAJOR: {
        value: 1,
        writable: false,
        configurable: false
    }
});

// Spotify's [key](https://developer.spotify.com/documentation/web-api/reference/tracks/get-audio-features/)
// is a numeric value representing
// [pitch class](https://en.wikipedia.org/wiki/Pitch_class#Other_ways_to_label_pitch_classes).
[['C', 'B_SHARP', 'D_DOUBLE_FLAT'],
 ['C_SHARP', 'D_FLAT', 'B_DOUBLE_SHARP'],
 ['D', 'C_DOUBLE_SHARP', 'E_DOUBLE_FLAT'],
 ['D_SHARP', 'E_FLAT', 'F_DOUBLE_FLAT'],
 ['E', 'D_DOUBLE_SHARP', 'F_FLAT'],
 ['F', 'E_SHARP', 'G_DOUBLE_FLAT'],
 ['F_SHARP', 'G_FLAT', 'E_DOUBLE_SHARP'],
 ['G', 'F_DOUBLE_SHARP', 'A_DOUBLE_FLAT'],
 ['G_SHARP', 'A_FLAT'],
 ['A', 'G_DOUBLE_SHARP', 'B_DOUBLE_FLAT'],
 ['A_SHARP', 'B_FLAT', 'C_DOUBLE_FLAT'],
 ['B', 'A_DOUBLE_SHARP', 'C_FLAT']].forEach((pitch, i) => {
    pitch.forEach(representation => {
        Object.defineProperty(self, representation, {
            value: i,
            writable: false,
            configurable: false
        });
    });
 });

// Dimensions that can be used to partition the tree. Most are supplied by
// Spotify, "aetherealness" and "primordialness" are our own. You can read more
// about all the different track features in
// [Spotify's API documentation](https://developer.spotify.com/documentation/web-api/reference/tracks/get-audio-features/).
Object.defineProperties(self, {
    AETHEREALNESS: {
        value: 'aetherealness',
        writable: false,
        configurable: false
    },
    PRIMORDIALNESS: {
        value: 'primordialness',
        writable: false,
        configurable: false
    },
    KEY: {
        value: 'key',
        writable: false,
        configurable: false
    },
    MODE: {
        value: 'mode',
        writable: false,
        configurable: false
    },
    TEMPO: {
        value: 'tempo',
        writable: false,
        configurable: false
    },
    VALENCE: {
        value: 'valence',
        writable: false,
        configurable: false
    },
    ENERGY: {
        value: 'energy',
        writable: false,
        configurable: false
    },
    DANCEABILITY: {
        value: 'danceability',
        writable: false,
        configurable: false
    },
    LOUDNESS: {
        value: 'loudness',
        writable: false,
        configurable: false
    },
    SPEECHINESS: {
        value: 'speechiness',
        writable: false,
        configurable: false
    },
    ACOUSTICNESS: {
        value: 'acousticness',
        writable: false,
        configurable: false
    },
    INSTRUMENTALNESS: {
        value: 'instrumentalness',
        writable: false,
        configurable: false
    },
    LIVENESS: {
        value: 'liveness',
        writable: false,
        configurable: false
    },
    TIME_SIGNATURE: {
        value: 'timeSignature',
        writable: false,
        configurable: false
    },
    DURATION_MS: {
        value: 'duration',
        writable: false,
        configurable: false
    },
    POPULARITY: {
        value: 'popularity',
        writable: false,
        configurable: false
    }
});

// ### API Declarations
// Defining an API function puts it on the global (`self`) object but prevents
// it from being overwritten or altered.
const defineAPIFunction = (key, fn) => Object.defineProperty(self, key, {
    value: fn,
    writable: false,
    configurable: false
});

// ### Web Security
// We use dependency injection to keep most things private and expose only
// certain desirable globals. All the user code runs in a web worker on an
// origin separate from the main Phosphorescence website, for security. Closures
// are used inside of the worker to prevent leaking anything to the global scope
// except what is explicitly desired. This way, the user script is sandboxed
// from the main origin, the UI thread, and bootstrapping code of the sandbox
// itself.
export function injector({getTree, getIdToTagMap, registerDimension}) {
    // ### The _k_-d tree
    // [_k_-d trees](https://en.wikipedia.org/wiki/K-d_tree) are trees that are
    // partitioned across _k_ dimensions. The use of these trees is for finding
    // "nearest neighbors" which are close-by points where closeness is defined
    // by the distance function applied to the two points during a comparison.
    // _k_-d trees are a much more optimal way of finding neighbors than the
    // naive way of comparing every point with every other point.

    // ### Hooks
    // Writing a playlist builder script requires you to write hook functions so
    // that the script runner is able to hook your logic functionality into the
    // system. For an example of these hooks, see
    // [the random walk builder](../builders/randomwalk.html) and
    // [the phosphorescence builder](../builders/phosphorescence.html).
    // You may return promises and use `async` functions if desired, for all
    // hooks. It will be properly handled by the runner.
    self.hooks = {
        // This hook is called prior to tree-building and allows you to prune
        // the list of tracks. It is not required.
        prune({tracks, idToTagMap, unprunedTracks}) {return buildResponse(tracks)},
        // This hook allows you to define your tree building function. The tree
        // is the main datastructure used in the playlist building search. All
        // available tracks are already transformed into `points`, however if
        // you want to define you own points we provide `tracks` as well as a
        // map of IDs to tags as the raw starting point.
        buildTree(kdTreeCtor, {points, tracks, idToTagMap}) {return buildResponse(null)},
        // This hook lets you specify logic for getting the first track of the
        // playlist.
        getFirstTrack({playlist, tags, goalTracks, points, tracks}, tree) {throw new Error('You must define a getFirstTrack hook!')},
        // And this one let's you specify logic for getting all subsequent
        // tracks, one at a time.
        getNextTrack({playlist, tags, goalTracks, points, previousTrack, tracks}, tree) {throw new Error('You must define a getNextTrack hook!')}
    };

    // Build response envelope for hook.
    defineAPIFunction('buildResponse', function(data) {
        return {data};
    });

    // Wrap point data correctly.
    defineAPIFunction('buildPoint', function(point) {
        return {point};
    });

    // Get _k_ nearest neighbors relative to a track. The track will be turned
    // into a standard point for you. If you are using your own custom point
    // style you cannot use this function.
    defineAPIFunction('getNearestNeighborsByTrack', function(k, track) {
        return getTree().nearest(Math.floor(k), getPointFromTrack(track));
    });

    // Get _k_ nearest neighbors relative to a raw point. This function is used
    // mainly when getting the first track where there is no prior track you are
    // using as comparison, or when you have defined your own point shape.
    defineAPIFunction('getNearestNeighbors', function(k, point) {
        return getTree().nearest(Math.floor(k), point);
    });

    // This simply returns a completely random value from the tree.
    defineAPIFunction('getRandomTrack', function() {
        return getTree().getRandomNode();
    });

    // This traverse every node using a custom check function.
    defineAPIFunction('getNodesWhere', function(fn) {
        return getTree().getNodesWhere(fn);
    });

    // This traverse every node in the tree for analysis purposes.
    defineAPIFunction('forEachNode', function(fn) {
        return getTree().forEach(fn);
    });

    // This returns how many nodes are left in the tree.
    defineAPIFunction('treeSize', function(fn) {
        return getTree().length();
    });

    // Allows a non-tree dimension to be registered for logging.
    defineAPIFunction('addLoggingDimension', function(dim) {
        registerDimension(dim);
    });

    // When searching for nearest neighbors,
    // [Euclidian distance](https://en.wikipedia.org/wiki/Euclidean_distance)
    // is a good distance function to use.
    defineAPIFunction('calculateEuclidianDistance', function(...distances) {
        if (distances.length == 0) {
            return 0;
        }
        else if (distances.length == 1) {
            return Math.abs(distances[0]);
        }
        else if (distances.includes(Infinity)) {
            return Infinity;
        }
        return Math.sqrt(distances.reduce((acc, cur) => acc + Math.pow(cur, 2), 0));
    });

    // ### Utility Functions

    // Constants used for
    // [harmonic mixing](https://mixedinkey.com/harmonic-mixing-guide/). These
    // are not exported, but rather used by the below functions.
    const minorsCircle = [A_FLAT, E_FLAT, B_FLAT, F, C, G, D, A, E, B, F_SHARP, D_FLAT];
    const majorsCircle = [B, F_SHARP, D_FLAT, A_FLAT, E_FLAT, B_FLAT, F, C, G, D, A, E];
    const minorsPitchToPositionMap = minorsCircle.reduce((acc, cur, i) => {acc[cur] = i; return acc;}, []);
    const majorsPitchToPositionMap = majorsCircle.reduce((acc, cur, i) => {acc[cur] = i; return acc;}, []);
    const majorToMinor = [A, A_SHARP, B, C, C_SHARP, D, D_SHARP, E, F, F_SHARP, G, G_SHARP];
    const minorToMajor = [D_SHARP, E, F, F_SHARP, G, G_SHARP, A, A_SHARP, B, C, C_SHARP, D];
    const keyUp = [G, G_SHARP, A, A_SHARP, B, C, C_SHARP, D, D_SHARP, E, F, F_SHARP];
    const keyDown = [F, F_SHARP, G, G_SHARP, A, A_SHARP, B, C, C_SHARP, D, D_SHARP, E];

    // An easy way to check if two tracks will be harmonically compatible. It is
    // debatable whether or not this is actually worthwhile. Many DJs apparently
    // do not do this or rarely do. However, sometimes you can get some pretty
    // cool results where a track seems to carry right through to the next. It
    // also may be interesting to make your own harmonic mappings for edgier
    // harmonic shifts.
    defineAPIFunction('nonJarringHarmonicDifference', function(a, b) {
        if (!a.mode || !a.key || !b.mode || !b.key) {
            return 1;
        }
        let diff;
        if (a.mode == b.mode) {
            // We can use major here even if they are both minor because
            // distance is the same regardless of offset on the wheel.
            diff = Math.abs(majorsPitchToPositionMap[a.key] - majorsPitchToPositionMap[b.key]);
            if (diff > 12 / 2) {
                diff = 12 % diff;
            }
        }
        else {
            let diff;
            if (a.mode == MINOR && b.mode == MAJOR) {
                diff = Math.abs(minorsPitchToPositionMap[a.key] - majorsPitchToPositionMap[b.key]);
            }
            else {
                diff = Math.abs(majorsPitchToPositionMap[a.key] - minorsPitchToPositionMap[b.key]);
            }
            if (diff > 12 / 2) {
                diff = 12 % diff;
            }
            diff += 1;
        }
        if (diff == 0) {
            return 0;
        }
        return 1 - (1 / (1 + (0.5 * Math.pow(diff, 3))));
    });

    // If you want to do harmonics comparisons outside of the tree search and
    // not as distance functions but boolean checks, these are some useful
    // starting point helpers.
    defineAPIFunction('sameHarmonics', function(a, b) {
        return a.mode == b.mode && a.key == b.key;
    });

    defineAPIFunction('sameModeAndNeighborKeyChange', function(a, b) {
        return a.mode == b.mode && (keyUp[a.key] == b.key || keyDown[a.key] == b.key);
    });

    defineAPIFunction('differentModeAndNeighborKeyChange', function(a, b) {
        return a.mode != b.mode && ((a.mode == MAJOR && majorToMinor[a.key] == b.key) || (a.mode == MINOR && minorToMajor[a.key] == b.key));
    });

    // This function seems to be reasonable for translating BPM changes into
    // a distance from 0 to 1.
    defineAPIFunction('nonJarringTempoDifference', function(a, b) {
        if (!a.tempo || !b.tempo) {
            return 1;
        }
        return Math.min(1, Math.max(0, 0.57 * Math.log(Math.abs(b.tempo - a.tempo))));
    });

    // [Fisher-Yates shuffle](https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle).
    // Mutates the original array and returns it back.
    const {shuffle} = require('../common/random.js');
    defineAPIFunction('shuffle', function(arr) {
        return shuffle(arr);
    });

    // Get a random integer within a range.
    const {getRandomInt} = require('../common/random.js');
    defineAPIFunction('getRandomInt', function(min, max) {
        return getRandomInt(min, max);
    });

    // Get random value from an array.
    defineAPIFunction('pickRandom', function(arr) {
        return arr[getRandomInt(0, arr.length - 1)];
    });

    // A helper function that simulates an n-sided dice roll with a target
    // minimum value.
    defineAPIFunction('rollDice', function(minTarget, sides) {
        if (minTarget > sides) {
            return true;
        }
        let rand = getRandomInt(0, sides - 1);
        return rand < minTarget;
    });

    // [Root mean square](https://en.wikipedia.org/wiki/Root_mean_square)
    // calculation.
    defineAPIFunction('calculateRMS', function(...vals) {
        return Math.sqrt(vals.reduce((acc, cur) => acc + Math.pow(cur, 2), 0) / vals.length);
    });

    // Tags are track IDs that can be used to identify tracks that are different
    // in Spotify's database but likely the same, or very similar, otherwise.
    // Remixes, mix cuts, alternate releases, and more all should end up getting
    // the same tag. We automatically keep track of the tags of all chosen
    // tracks during a playlist build, but we don't do anything with that by
    // default. If you want to prune these potential duplicates, you can use
    // this function. Other cool things to use tags for might be purposefully
    // calling back to a track that was previously played by finding a remix to
    // play later, if one exists in the track listings.
    defineAPIFunction('cullTracksWithAlreadySeenTags', function(neighbors, previousTags) {
        let unculledNeighbors = [];
        for (let nn of neighbors) {
            if (previousTags[nn.point.tag]) {
                continue;
            }
            unculledNeighbors.push(nn);
        }
        return unculledNeighbors;
    });

    // These functions are helpers for using the baseline track to point
    // conversion method. By default, this is done for you and handed off to the
    // tree building function for you to use as needed.
    defineAPIFunction('tracksToPoints', function(tracks) {
        let evocativenessPoints = [];
        for (let trackWrapper of Object.values(tracks)) {
            evocativenessPoints.push(getPointFromTrack(trackWrapper));
        }
        return evocativenessPoints;
    });

    defineAPIFunction('getPointFromTrack', function(trackWrapper) {
        let idsToTags = getIdToTagMap();
        let {id, track, features, evocativeness} = trackWrapper;
        let {popularity} = track;
        let {
            key,
            mode,
            tempo,
            energy,
            valence,
            liveness,
            loudness,
            speechiness,
            acousticness,
            danceability,
            instrumentalness,
            duration_ms: duration,
            time_signature: timeSignature
        } = features;
        let {aetherealness, primordialness} = evocativeness;
        let tag = idsToTags[id];
        return {
            id,
            tag,
            key,
            mode,
            tempo,
            energy,
            valence,
            liveness,
            loudness,
            duration,
            popularity,
            speechiness,
            acousticness,
            danceability,
            timeSignature,
            aetherealness,
            primordialness,
            instrumentalness
        };
    });
};
