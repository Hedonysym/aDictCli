package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

var notFoundErr = fmt.Errorf("not found")

func main() {
	var depth int = 3

	repoFlag := flag.String("repo", "", "path to repository")
	dbFlag := flag.String("db", "", "path to dictionary.db")
	helpFlag := flag.Bool("help", false, "show help")
	verboseFlag := flag.Bool("verbose", false, "show full results")
	generateFlag := flag.Bool("generate", false, "generate the database")
	flag.Parse()

	dbPath, err := resolveDBPath(*dbFlag)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return err
	}

	if *helpFlag {
		fmt.Println("Welcome to aDictCli\n\nusage: adict --<optional flag> <word>")
		fmt.Println("optional flags:\n--verbose: show full results(otherwise defaults to 3)\n--generate: generate the database\n")
		os.Exit(0)
	}
	if *verboseFlag {
		depth = 100
	}
	if *generateFlag {
		repo := *repoFlag
		if repo == "" {
			if r, err := findRepoRoot(); err == nil {
				repo = r
			} else if env := os.Getenv("ADICT_REPO"); env != "" {
				repo = env
			} else {
				log.Fatal("repo root not found: pass --repo or set ADICT_REPO")
			}
		}
		in := filepath.Join(repo, "data", "dictionary.json")
		out := dbPath
		schema := filepath.Join(repo, "cmd", "gen", "schema.sql")

		log.Printf("Generating DB\n  repo: %s\n  in: %s\n  out: %s\n  schema: %s", repo, in, dbPath, schema)
		cmd := exec.Command("go", "-C", repo, "run", "./cmd/gen", "-in", in, "-out", out, "-schema", schema)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		log.Println("Generation complete")
		os.Exit(0)
	}
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: adict --<optional flag> <word> \nuse adict --help for more info")
		os.Exit(0)
	}
	word := strings.ToLower(strings.TrimSpace(flag.Args()[0]))

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=ro&immutable=1", dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = dictionaryLookup(db, word, depth)
	if err == notFoundErr {
		suggestion, err := fuzzySuggestions(db, word)
		if err != nil {
			log.Fatal(err)
		}
		if suggestion == "" {
			os.Exit(0)
		}
		err = dictionaryLookup(db, suggestion, depth)

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

func resolveDBPath(cli string) (string, error) {
	if cli != "" {
		return filepath.Abs(cli)
	}
	if env := os.Getenv("ADICT_DB"); env != "" {
		return filepath.Abs(env)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".local", "share", "adict")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Join(dir, "dictionary.db"))
}
func findRepoRoot() (string, error) {
	// start from CWD
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		p := filepath.Dir(dir)
		if p == dir {
			return "", fmt.Errorf("go.mod not found; set --repo")
		}
		dir = p
	}
}
