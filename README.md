# Hansa Teutonica

Implementation of board game Hansa Teutonica for pedagogical purposes.

## Requirements

* Sqlite 3

## Setup

```sh
make
make run

# Or without make:
go build -o hansa .
./hansa -listenaddr 0.0.0.0 -port 8080
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
