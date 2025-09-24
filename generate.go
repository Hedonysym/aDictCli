//go:build tools

//go:generate go run ./cmd/gen -in ./data/webster.json -out ./assets/dictionary.db

package tools
