# Phosphorescence

https://phosphor.me

Phosphorescence is a web application that integrates with Spotify which produces smart Trance music sets/playlists. It works off of a machine learning model that classifies Spotify audio features into a 2-axis spectrum of "evocativeness" trained by me, based off a classification concept I came up with for Trance music in general that is supposed to ignore artificial things like release date/genere boundaries/artist and focus on "evocativeness". One axis is primordialness to transcendentalness and the other axis is chthonicness to aetherealness. For more info, you can find the file EVOCATIVE_SUBGENRES.md in the training folder.

The web application builds out a _k_-d tree across some dimensions and the playlist is built by randomly selecting from a pool of nearest neighbors. The reason for the vague description is because actually, this logic is programmable. A script is written that let's the user specify the tree building, the first track selection, and all subsequent track selection, using various Spotify features as well as my evocativeness features. The simplest implementation (which we include bundled) is the "random walk" which selects a random track to start and then uses that track to randomly walk to a nearest neighbor out of X nearest neighbors where the nearness is defined by evocativeness distance across 2 dimensions as well as harmonic distance across mode and key dimensions and tempo distance as well. This means the selected track should flow from the previous track relatively well harmonically speaking as well as BPM-wise and also be evocatively similar while walking around the k-space from track 1 to track Y.

Built in is the deletion of tracks that have already been picked from the tree as well as some helper functionality that allows ignoring tracks which are similarly "tagged". A tag is basically a reduction of a track name/artist into something which will be equivalent across releases/remixes/etc to avoid the obvious issue of 50 remixes of the same track just coming one after another (as this sort of releasing is common in the genre).

It is possible for a user to program their own playlist builders using all of the tooling from right in the browser itself.

Also, when a user visits the site they are handed a cache of Spotify tracks to build the tree from, which is filtered client-side by region (so a user will not get tracks unplayable in their region). This cache is re-populated nightly and in the US has about 15k reduced to 11k tracks, after region filtering, at the time of this writing.
