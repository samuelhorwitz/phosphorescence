package spider

import (
	"fmt"
	"regexp"
)

var mixCutMatcher = regexp.MustCompile(`(?i)([\[(](mix cut|mixed)[\])]|\bmix cut)`)

// This should have a very low rate of false positives, and a hard to determine
// false negative rate. Basically a lot of tracks may be pre-mixed and unlabeled,
// but it's exceedingly unlikely any tracks are NOT premixed but somehow labeled
// "mix cut" erroneously.
func isMixCut(track track) bool {
	return mixCutMatcher.MatchString(track.Name)
}

func isBlacklistedArtist(track track, artistBlacklist map[string]bool) (bool, error) {
	if len(track.Artists) == 0 {
		return false, fmt.Errorf("No artists on track %s", track.ID)
	}
	return artistBlacklist[track.Artists[0].ID], nil
}

func isRemovedFromSpotify(track track) bool {
	if len(track.AvailableMarkets) == 0 {
		return true
	}
	return false
}
