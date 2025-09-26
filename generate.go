//go:generate mkdir -p $HOME/.local/share/adict
//go:generate go -C . run ./cmd/gen -in ./data/dictionary.json -out $HOME/.local/share/adict/dictionary.db -schema ./cmd/gen/schema.sql
package tools
