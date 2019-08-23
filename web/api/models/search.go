package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"regexp"
)

var (
	markTag          = regexp.MustCompile(`(?i)</?mark>`)
	hashtagFirstRune = regexp.MustCompile(`[\pL\pN]`)
)

// rank real, id uuid, type searchable_type, name text, description text, author_name text, likes bigint
type searchResult struct {
	Rank             float64    `json:"rank"`
	ID               uuid.UUID  `json:"id"`
	ResultType       resultType `json:"resultType"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	AuthorName       string     `json:"authorName"`
	LikeCount        uint64     `json:"likeCount"`
	NameMarks        []int      `json:"nameMarks"`
	DescriptionMarks []int      `json:"descriptionMarks"`
	AuthorNameMarks  []int      `json:"authorNameMarks"`
}

type resultType string

const (
	scriptResultType      resultType = "script"
	scriptChainResultType resultType = "script_chain"
)

func Query(q string) (searchResults []searchResult, _ error) {
	rows, err := postgresDB.Query("select rank, id, type, name, description, author_name, likes from search($1) limit 50", q)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Search failed: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var searchResult searchResult
		var name string
		var authorName string
		var description string
		err := rows.Scan(&searchResult.Rank, &searchResult.ID, &searchResult.ResultType, &name, &description, &authorName, &searchResult.LikeCount)
		if err != nil {
			return nil, fmt.Errorf("Could not scan row: %s", err)
		}
		searchResult.Name, searchResult.NameMarks = getPlaintextAndMarkIndices(name)
		searchResult.AuthorName, searchResult.AuthorNameMarks = getPlaintextAndMarkIndices(authorName)
		searchResult.Description, searchResult.DescriptionMarks = getPlaintextAndMarkIndices(description)
		searchResults = append(searchResults, searchResult)
	}
	return searchResults, nil
}

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
