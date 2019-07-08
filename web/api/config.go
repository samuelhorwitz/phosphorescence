package main

type config struct {
	isProduction             bool
	phosphorOrigin           string
	spotifyClientID          string
	spotifySecret            string
	apiOrigin                string
	spacesID                 string
	spacesSecret             string
	spacesTracksEndpoint     string
	spacesTracksRegion       string
	spacesScriptsEndpoint    string
	spacesScriptsRegion      string
	postgresConnectionString string
}
