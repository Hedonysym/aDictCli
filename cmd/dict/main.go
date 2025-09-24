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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dict <word>")
		os.Exit(1)
	}
	q := strings.ToLower(strings.TrimSpace(os.Args[1]))

	db, err := sql.Open("sqlite", "file:assets/dictionary.db?mode=ro&immutable=1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT pos, defs FROM entries WHERE word = ? ORDER BY pos`, q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		var pos, defsJSON string
		if err := rows.Scan(&pos, &defsJSON); err != nil {
			log.Fatal(err)
		}
		var defs []string
		_ = json.Unmarshal([]byte(defsJSON), &defs)
		fmt.Printf("%s (%s)\n", q, pos)
		for i, d := range defs {
			fmt.Printf("  %d. %s\n", i+1, d)
		}
	}
	if !found {
		log.Fatal("word not found")
	}
}
