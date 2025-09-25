package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

var notFoundErr = fmt.Errorf("not found")

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: adict --<optional flag> <word> \nuse adict --help for more info")
		os.Exit(0)
	}
	word := strings.ToLower(strings.TrimSpace(os.Args[1]))

	db, err := sql.Open("sqlite", "file:assets/dictionary.db?mode=ro&immutable=1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = dictionaryLookup(db, word, 3)
	if err == notFoundErr {
		suggestion, err := fuzzySuggestions(db, word)
		if suggestion == "" {
			os.Exit(0)
		}
		if err != nil {
			log.Fatal(err)
		}
		err = dictionaryLookup(db, suggestion, 3)

		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func dictionaryLookup(db *sql.DB, word string, depth int) error {
	rows, err := db.Query(`SELECT pos, defs FROM entries WHERE word = ? ORDER BY pos`, word)
	if err != nil {
		return err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		var pos, defsJSON string
		if err := rows.Scan(&pos, &defsJSON); err != nil {
			return err
		}
		var defs []string
		_ = json.Unmarshal([]byte(defsJSON), &defs)
		fmt.Printf("\033[1m%s (%s)\033[0m\n", word, pos)
		for i, d := range defs {
			if i >= depth {
				break
			}
			fmt.Printf("  %d. %s\n", i+1, d)
		}
	}
	if !found {
		return notFoundErr
	}
	return nil
}
