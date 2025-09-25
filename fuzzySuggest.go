package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/lithammer/fuzzysearch/fuzzy"
	_ "modernc.org/sqlite"
)

func getFuzzyMatches(db *sql.DB, input string) (fuzzy.Ranks, error) {
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
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Distance < matches[j].Distance
	})

	return matches[:3], nil
}

func fuzzySuggestions(db *sql.DB, input string) (string, error) {
	matches, err := getFuzzyMatches(db, input)
	if err != nil {
		return "", err
	}
	fmt.Println("\033[1mDid you mean:\033[0m")
	for i, m := range matches {
		fmt.Printf("%d. %s\n", i+1, m.Target)
	}
	fmt.Println("0. Exit")

	reader := bufio.NewReader(os.Stdin)

	input, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = input[:len(input)-1]
	number, err := strconv.Atoi(input)
	if number == 0 {
		return "", nil
	}
	if err != nil || number < 0 || number > len(matches) {
		fmt.Println("Invalid input. Please enter a valid number.")
		return "", err
	}

	result := matches[number-1].Target

	return result, nil
}
