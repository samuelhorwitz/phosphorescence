package spider

import (
	"fmt"
	"log"
)

func GetTracks(cfg *Config) ([]*TrackEnvelope, map[string][]*TrackEnvelope, error) {
	playlists, err := loadPlaylists()
	if err != nil {
		return nil, nil, fmt.Errorf("Could not load playlists: %s", err)
	}
	err = initializeToken(cfg.SpotifyClientID, cfg.SpotifySecret)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not initialize token: %s", err)
	}
	blacklist, err := buildBlacklist(playlists)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not build blacklist: %s", err)
	}
	artistBlacklist, err := buildArtistBlacklist(playlists)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not build artist blacklist: %s", err)
	}
	allTracks, err := buildTracks(playlists, blacklist, artistBlacklist)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not build tracks: %s", err)
	}
	finalTracks, err := getTrackFeaturesInBatches(allTracks)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not get track features: %s", err)
	}
	var tracksArr []*TrackEnvelope
	for id, trackEnvelope := range finalTracks {
		trackEnvelope.ID = id
		tracksArr = append(tracksArr, trackEnvelope)
	}
	shardedTracks := regionShard(tracksArr)
	return tracksArr, shardedTracks, nil
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
			if len(blacklistedArtistTrack.Track.Artists) == 0 {
				return nil, fmt.Errorf("Could not get primary artist, no artists on track %s", id)
			}
			artistBlacklist[blacklistedArtistTrack.Track.Artists[0].ID] = true
		}
	}
	return artistBlacklist, nil
}

func buildTracks(playlists playlists, blacklist map[string]bool, artistBlacklist map[string]bool) (map[string]*TrackEnvelope, error) {
	allTracks := make(map[string]*TrackEnvelope)
	// Whitelists are seeders with no exclusions applied except for no
	// longer being available in any region
	for _, whitelistPlaylist := range playlists.Whitelists {
		whitelistedTracks, err := getTracksFromPlaylist(whitelistPlaylist)
		if err != nil {
			return nil, fmt.Errorf("Could not get whitelisted tracks: %s", err)
		}
		for id, whitelistedTrack := range whitelistedTracks {
			if id == "" {
				continue
			}
			if isRemovedFromSpotify(whitelistedTrack.Track) {
				log.Printf(`Skipping track that no longer is available "%s"`, whitelistedTrack.Track.Name)
				continue
			}
			allTracks[id] = whitelistedTrack
		}
	}
	// Seeders exclude any blacklisted tracks as well as anything that is deemed a
	// "mix cut" which is basically any track already part of a mixed set. They
	// also exclude tracks without any playable region.
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
			if isRemovedFromSpotify(seedTrack.Track) {
				log.Printf(`Skipping track that no longer is available "%s"`, seedTrack.Track.Name)
				continue
			}
			isArtistBlacklisted, err := isBlacklistedArtist(seedTrack.Track, artistBlacklist)
			if err != nil {
				return nil, fmt.Errorf("Could not check if artist is blacklisted: %s", err)
			}
			if isArtistBlacklisted {
				log.Printf(`Skipping blacklisted artist "%s"`, seedTrack.Track.Artists[0].Name)
				continue
			}
			if isMixCut(seedTrack.Track) {
				log.Printf(`Skipping mixed track "%s"`, seedTrack.Track.Name)
				continue
			}
			allTracks[id] = seedTrack
		}
	}
	return allTracks, nil
}
