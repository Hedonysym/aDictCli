package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runDict(path, word string) error {
	cmd := exec.Command("go", "run", path, word)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	absPath, err := filepath.Abs("cmd/dict/main.go")
	if err != nil {
		log.Fatal(err)
	}
	if err := runDict(absPath, os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
