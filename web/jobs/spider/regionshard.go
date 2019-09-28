package spider

func regionShard(tracks []*TrackEnvelope) map[string][]*TrackEnvelope {
	shards := make(map[string][]*TrackEnvelope)
	for _, trackEnvelope := range tracks {
		for _, market := range trackEnvelope.Track.AvailableMarkets {
			shards[market] = append(shards[market], trackEnvelope)
		}
	}
	return shards
}
