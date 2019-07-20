package spider

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type playlists struct {
	Whitelists       []string `json:"whitelists"`
	Seeders          []string `json:"seeders"`
	Blacklists       []string `json:"blacklists"`
	ArtistBlacklists []string `json:"artistBlacklists"`
}

func loadPlaylists() (playlists, error) {
	ex, err := os.Executable()
	if err != nil {
		return playlists{}, fmt.Errorf("Could not get executable path: %s", err)
	}
	exPath := filepath.Dir(ex)
	playlistsPath := filepath.Join(exPath, "playlists.json")
	playlistFile, err := os.Open(playlistsPath)
	if err != nil {
		return playlists{}, fmt.Errorf("Could not open playlists file: %s", err)
	}
	defer playlistFile.Close()
	var jsonLines []string
	scanner := bufio.NewScanner(playlistFile)
	for scanner.Scan() {
		pieces := strings.Split(scanner.Text(), "//")
		if len(pieces) > 0 {
			jsonLines = append(jsonLines, pieces[0])
		}
	}
	jsonData := strings.Join(jsonLines, "\n")
	var pl playlists
	err = json.Unmarshal([]byte(jsonData), &pl)
	if err != nil {
		return playlists{}, fmt.Errorf("Could not unmarshal playlists JSON: %s", err)
	}
	return pl, nil
}
