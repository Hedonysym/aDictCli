# aDictCli

An offline cli dictionary with autocompletion/spell correction based on phoenetic matching, written in golang for a boot.dev solo project.

## Included external packages:

- modernc.org/sqlite and dependancies
- github.com/lithammer/fuzzysearch/fuzzy

A special thanks to https://github.com/ssvivian/WebstersDictionary for use of their json formatted webster's dictionary.

## Blueprint:

An sqlite database will first be generated from a json file to allow fast lookups.