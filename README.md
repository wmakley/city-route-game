# City Route Game

Implementation of a board game similar to Hansa Teutonica (© 2020 Pegasus Spiele GmbH) for pedagogical purposes.

## Requirements

* Sqlite 3

## Setup

```sh
make
make run
```

Command line flags are optional. Listens on ":8080" by default.

## Development

```sh
# Run this OUTSIDE the project directory:
go get github.com/cespare/reflex

# Inside project directory again:
make watch
```

Simple script that watches go source files for changes and re-compiles and restarts server.
