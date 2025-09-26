package main

import (
	"database/sql"
	"log"
	"testing"
)

func TestFuzzy(t *testing.T) {
	db, err := sql.Open("sqlite", "file:assets/dictionary.db?mode=ro&immutable=1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	matches, err := fuzzySuggestions(db, "nigga")
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range matches {
		t.Log(m.Source, m.Target)
	}
	if len(matches) != 3 {
		t.Errorf("expected 3 matches, got %d", len(matches))
	}
}
