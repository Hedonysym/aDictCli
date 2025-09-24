package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

type entry struct {
	Word        string   `json:"word"`
	POS         string   `json:"pos"`
	Definitions []string `json:"definitions"`
}

func main() {
	in := flag.String("in", "./data/webster.json", "input JSON")
	out := flag.String("out", "./assets/dictionary.db", "output SQLite DB")
	schema := flag.String("schema", "./cmd/gen/schema.sql", "schema SQL file")
	flag.Parse()

	if err := run(*in, *out, *schema); err != nil {
		log.Fatal(err)
	}
}

func run(inPath, outPath, schemaPath string) error {
	_ = os.MkdirAll("assets", 0o755)

	db, err := sql.Open("sqlite", outPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// pragmas for fast build and apply schema

	if _, err = db.Exec(`PRAGMA journal_mode=OFF; PRAGMA synchonous=OFF; PRAGMA temp_store=MEMORY;`); err != nil {
		return err
	}
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}
	if _, err = db.Exec(string(schemaSQL)); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	insEntry, err := tx.Prepare(`INSERT INTO entries(word,pos,defs) VALUES(?,?,?) on conflict (word,pos) do update set defs=excluded.defs`)
	if err != nil {
		return err
	}
	defer insEntry.Close()
	selDefs, _ := tx.Prepare(`SELECT defs FROM entries WHERE word=? AND pos=?`)
	defer selDefs.Close()

	updDefs, _ := tx.Prepare(`UPDATE entries SET defs=? WHERE word=? AND pos=?`)
	defer updDefs.Close()

	f, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if d, ok := tok.(json.Delim); !ok || d != '[' {
		return io.ErrUnexpectedEOF
	}

	count := 0
	for dec.More() {
		var e entry
		if err := dec.Decode(&e); err != nil {
			return err
		}

		word := strings.ToLower(strings.TrimSpace(e.Word))
		pos := strings.TrimSpace(e.POS)

		defsJSON, err := json.Marshal(e.Definitions)
		if err != nil {
			return err
		}

		var res sql.Result

		if res, err = insEntry.Exec(word, wordCase(pos), string(defsJSON)); err != nil {
			return err
		}
		aff, _ := res.RowsAffected()
		if aff == 0 {
			// conflict: merge
			var curr string
			if err := selDefs.QueryRow(word, pos).Scan(&curr); err != nil {
				return err
			}
			var arr []string
			_ = json.Unmarshal([]byte(curr), &arr)
			arr = append(arr, e.Definitions...)
			merged, _ := json.Marshal(arr)
			if _, err := updDefs.Exec(string(merged), word, pos); err != nil {
				return err
			}
		}
		count++
		if count%5000 == 0 {
			log.Printf("inserted %d entries", count)
		}
	}
	if _, err := dec.Token(); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	_, _ = db.Exec(`VACUUM; PRAGMA optimize;`)
	log.Printf("done, total entries: %d", count)
	return nil
}

func wordCase(s string) string { return s }
