package spider

import (
	"encoding/json"
	"fmt"
	"log"
)

func GetTracks(cfg *Config) (map[string]*TrackEnvelope, error) {
	playlists, err := loadPlaylists()
	if err != nil {
		return nil, fmt.Errorf("Could not load playlists: %s", err)
	}
	err = initializeToken(cfg.SpotifyClientID, cfg.SpotifySecret)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize token: %s", err)
	}
	blacklist, err := buildBlacklist(playlists)
	if err != nil {
		return nil, fmt.Errorf("Could not build blacklist: %s", err)
	}
	artistBlacklist, err := buildArtistBlacklist(playlists)
	if err != nil {
		return nil, fmt.Errorf("Could not build artist blacklist: %s", err)
	}
	allTracks, err := buildTracks(playlists, blacklist, artistBlacklist)
	if err != nil {
		return nil, fmt.Errorf("Could not build tracks: %s", err)
	}
	finalTracks, err := getTrackFeaturesInBatches(allTracks)
	if err != nil {
		return nil, fmt.Errorf("Could not get track features: %s", err)
	}
	return finalTracks, nil
}

func buildBlacklist(playlists playlists) (map[string]bool, error) {
	blacklist := make(map[string]bool)
	for _, blacklistPlaylist := range playlists.Blacklists {
		blacklistedTracks, err := getTracksFromPlaylist(blacklistPlaylist)
		if err != nil {
			return nil, fmt.Errorf("Could not get blacklisted tracks: %s", err)
		}
		for id := range blacklistedTracks {
			blacklist[id] = true
		}
	}
	return blacklist, nil
}

func buildArtistBlacklist(playlists playlists) (map[string]bool, error) {
	artistBlacklist := make(map[string]bool)
	for _, artistBlacklistPlaylist := range playlists.ArtistBlacklists {
		blacklistedArtistTracks, err := getTracksFromPlaylist(artistBlacklistPlaylist)
		if err != nil {
			return nil, fmt.Errorf("Could not get blacklisted artist trackss: %s", err)
		}
		for id, blacklistedArtistTrack := range blacklistedArtistTracks {
			var track struct {
				Artists []struct {
					ID string `json:"id"`
				} `json:"artists"`
			}
			err := json.Unmarshal(blacklistedArtistTrack.Track, &track)
			if err != nil {
				return nil, fmt.Errorf("Could not parse Spotify track: %s", err)
			}
			if len(track.Artists) == 0 {
				return nil, fmt.Errorf("Could not get primary artist, no artists on track %s", id)
			}
			artistBlacklist[track.Artists[0].ID] = true
		}
	}
	return artistBlacklist, nil
}

func buildTracks(playlists playlists, blacklist map[string]bool, artistBlacklist map[string]bool) (map[string]*TrackEnvelope, error) {
	allTracks := make(map[string]*TrackEnvelope)
	// Whitelists are seeders with no exclusions applied
	for _, whitelistPlaylist := range playlists.Whitelists {
		whitelistedTracks, err := getTracksFromPlaylist(whitelistPlaylist)
		if err != nil {
			return nil, fmt.Errorf("Could not get whitelisted tracks: %s", err)
		}
		for id, whitelistedTrack := range whitelistedTracks {
			if id == "" {
				continue
			}
			allTracks[id] = whitelistedTrack
		}
	}
	// Seeders exclude any blacklisted tracks as well as anything that is deemed a
	// "mix cut" which is basically any track already part of a mixed set.
	for _, seedPlaylist := range playlists.Seeders {
		seedTracks, err := getTracksFromPlaylist(seedPlaylist)
		if err != nil {
			return nil, fmt.Errorf("Could not get seed tracks: %s", err)
		}
		for id, seedTrack := range seedTracks {
			if id == "" {
				continue
			}
			if blacklist[id] {
				continue
			}
			isArtistBlacklisted, artistName, err := checkIfBlacklistedArtist(seedTrack.Track, artistBlacklist)
			if err != nil {
				return nil, fmt.Errorf("Could not check if artist is blacklisted: %s", err)
			}
			if isArtistBlacklisted {
				log.Printf(`Skipping blacklisted artist "%s"`, artistName)
				continue
			}
			isMixCut, trackName, err := checkIfMixCut(seedTrack.Track)
			if err != nil {
				return nil, fmt.Errorf("Could not check if mix cut: %s", err)
			}
			if isMixCut {
				log.Printf(`Skipping mixed track "%s"`, trackName)
				continue
			}
			allTracks[id] = seedTrack
		}
	}
	return allTracks, nil
}
