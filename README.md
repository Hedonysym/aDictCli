# aDictCli

An offline cli dictionary with autocompletion/spell correction, written in golang for a boot.dev solo project.

## Included external packages:

- modernc.org/sqlite and dependancies
- github.com/mattn/go-sqlite3

A special thanks to https://github.com/ssvivian/WebstersDictionary for use of their json formatted webster's dictionary.

## Install

first download the repo files and place them in the desired location on your filesystem

then run ```go generate ./... && go install ./cmd/adict``` from the root of the repo (the folder with go.mod)
this will install the dictionary sql at ~/.local/share/adict by default and can be changed with a .env file at the root of the repo

```ADICT_DB="$HOME/.local/share/adict/dictionary.db"```

after that you can run the command from anywhere with ```adict <args>```

you can also use the following to customize your installation

go -C abs/path/to/repo run ./cmd/gen -in ./data/dictionary.json -out "abs/path/to/database/destination" -schema ./cmd/gen/schema.sql

## Enjoy
