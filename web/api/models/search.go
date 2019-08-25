package models

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
	"regexp"
)

var (
	markTag          = regexp.MustCompile(`(?i)</?mark>`)
	hashtagFirstRune = regexp.MustCompile(`[\pL\pN]`)
	punctuation      = regexp.MustCompile(`[\p{Pi}\x{ff02}\x{ff07}"'#]`)
	quoted           = regexp.MustCompile(`(?i)[\p{Pi}\x{ff02}\x{ff07}"'](.*?)[\p{Pf}\x{ff02}\x{ff07}"']`)
	hashtag          = regexp.MustCompile(`#((?:[\pL\pN][\pM\x{200C}\x{200D}]*)+(?:[\p{Pc}\p{Pd}](?:[\pL\pN][\pM\x{200C}\x{200D}]*)+)*)`)
	finalWord        = regexp.MustCompile(`^((?:\PZ*[\pZ\p{Pi}\p{Pf}\x{ff02}\x{ff07}"'#]+)*)(\PZ*)$`)
)

// rank real, id uuid, type searchable_type, name text, description text, author_name text, likes bigint
type searchResult struct {
	Rank             float64    `json:"rank,omitempty"`
	ID               uuid.UUID  `json:"id"`
	ResultType       resultType `json:"resultType"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	AuthorName       string     `json:"authorName"`
	LikeCount        uint64     `json:"likeCount"`
	NameMarks        []int      `json:"nameMarks,omitempty"`
	DescriptionMarks []int      `json:"descriptionMarks,omitempty"`
	AuthorNameMarks  []int      `json:"authorNameMarks,omitempty"`
}

type resultType string

const (
	scriptResultType      resultType = "script"
	scriptChainResultType resultType = "script_chain"
)

func Query(q string) (searchResults []searchResult, _ error) {
	q, strictMatches, tags := parseQuery(q)
	rows, err := postgresDB.Query("select rank, id, type, name, description, author_name, likes from search($1, $2, $3) limit 50", q, pq.Array(strictMatches), pq.Array(tags))
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

func QueryTag(tag string) (searchResults []searchResult, _ error) {
	rows, err := postgresDB.Query("select id, type, name, unmodified_description, author_name, likes from search_hashtag($1) limit 100", tag)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Tag search failed: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var searchResult searchResult
		err := rows.Scan(&searchResult.ID, &searchResult.ResultType, &searchResult.Name, &searchResult.Description, &searchResult.AuthorName, &searchResult.LikeCount)
		if err != nil {
			return nil, fmt.Errorf("Could not scan row: %s", err)
		}
		searchResults = append(searchResults, searchResult)
	}
	return searchResults, nil
}

func RecommendedQuery(q string) (queries []string, err error) {
	firstPart, lastWord := splitQueryLastWord(q)
	if lastWord == "" {
		return []string{firstPart}, nil
	}
	fn := "get_recommended_word"
	if len(firstPart) > 0 && string(firstPart[len(firstPart)-1]) == "#" {
		fn = "get_recommended_tag"
	}
	var recommendedWords []string
	err = postgresDB.QueryRow(fmt.Sprintf("select %s($1, 5)", fn), lastWord).Scan(pq.Array(&recommendedWords))
	if err != nil {
		return nil, fmt.Errorf("Couldn't execute query: %s", err)
	}
	for _, recommendedWord := range recommendedWords {
		queries = append(queries, fmt.Sprintf("%s%s", firstPart, recommendedWord))
	}
	return queries, nil
}

func parseQuery(q string) (_ string, strictMatches, tags []string) {
	for _, match := range quoted.FindAllStringSubmatch(q, -1) {
		strictMatches = append(strictMatches, match[1])
	}
	for _, match := range hashtag.FindAllStringSubmatch(q, -1) {
		tags = append(tags, match[1])
	}
	return punctuation.ReplaceAllString(q, " "), strictMatches, tags
}

func getPlaintextAndMarkIndices(markedUp string) (plain string, marks []int) {
	offset := 0
	lastEnd := 0
	markSubstrings := markTag.FindAllStringIndex(markedUp, -1)
	if len(markSubstrings) == 0 {
		return markedUp, nil
	}
	plain += markedUp[0:markSubstrings[0][0]]
	for i, index := range markSubstrings {
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

func splitQueryLastWord(q string) (string, string) {
	matches := finalWord.FindStringSubmatch(q)
	return matches[1], matches[2]
}
