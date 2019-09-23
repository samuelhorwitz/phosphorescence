package spider

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var mixCutMatcher = regexp.MustCompile(`(?i)([\[(](mix cut|mixed)[\])]|\bmix cut)`)

// This should have a very low rate of false positives, and a hard to determine
// false negative rate. Basically a lot of tracks may be pre-mixed and unlabeled,
// but it's exceedingly unlikely any tracks are NOT premixed but somehow labeled
// "mix cut" erroneously.
func checkIfMixCut(trackJSON json.RawMessage) (bool, string, error) {
	var track struct {
		Name string `json:"name"`
	}
	err := json.Unmarshal(trackJSON, &track)
	if err != nil {
		return false, "", fmt.Errorf("Could not parse Spotify track: %s", err)
	}
	return mixCutMatcher.MatchString(track.Name), track.Name, nil
}

func checkIfBlacklistedArtist(trackJSON json.RawMessage, artistBlacklist map[string]bool) (bool, string, error) {
	var track struct {
		ID      string `json:"id"`
		Artists []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"artists"`
	}
	err := json.Unmarshal(trackJSON, &track)
	if err != nil {
		return false, "", fmt.Errorf("Could not parse Spotify track: %s", err)
	}
	if len(track.Artists) == 0 {
		return false, "", fmt.Errorf("No artists on track %s", track.ID)
	}
	return artistBlacklist[track.Artists[0].ID], track.Artists[0].Name, nil
}

func checkIfRemovedFromSpotify(trackJSON json.RawMessage) (bool, string, error) {
	var track struct {
		Name             string   `json:"name"`
		AvailableMarkets []string `json:"available_markets"`
	}
	err := json.Unmarshal(trackJSON, &track)
	if err != nil {
		return false, "", fmt.Errorf("Could not parse Spotify track: %s", err)
	}
	if len(track.AvailableMarkets) == 0 {
		return false, track.Name, nil
	}
	return true, track.Name, nil
}
