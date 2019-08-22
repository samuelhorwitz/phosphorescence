package search

import "regexp"

var (
	markTag          = regexp.MustCompile(`(?i)</?mark>`)
	hashtagFirstRune = regexp.MustCompile(`[\pL\pN]`)
)

func getPlaintextAndMarkIndices(markedUp string) (plain string, marks []int) {
	offset := 0
	lastEnd := 0
	for i, index := range markTag.FindAllStringIndex(markedUp, -1) {
		start := index[0]
		end := index[1]
		realStart := start
		// If we see the beginning of hashtag, then we want to include the hash
		// mark in the marked up tag (Postgres doesn't include punctuation in
		// our search). So we push the start back by one to include it.
		if start > 0 && i%2 == 0 && string(markedUp[start-1]) == "#" && hashtagFirstRune.MatchString(string(markedUp[end])) {
			start--
		}
		marks = append(marks, start-offset)
		if lastEnd > 0 {
			plain += markedUp[lastEnd:realStart]
		}
		lastEnd = end
		offset += end - realStart
	}
	plain += markedUp[lastEnd:len(markedUp)]
	return
}
