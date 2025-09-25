package main

import (
	"database/sql"
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
	_ "modernc.org/sqlite"
)

func fuzzySuggestions(db *sql.DB, input string) (fuzzy.Ranks, error) {
	var words []string
	rows, err := db.Query("SELECT word From words;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	matches := fuzzy.RankFind(input, words)
	sort.Sort(matches)

	return matches[:3], nil
}
