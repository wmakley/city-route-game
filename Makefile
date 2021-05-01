.PHONY: clean clean_db test run watch run_clean_db

GO_SRCS := $(shell find . -type f -name '*.go' -and -not -name '*_test.go')

REFLEX := $(shell command -v reflex 2> /dev/null)

all: bin/admin

bin/admin: $(GO_SRCS)
	go build -o $@ cmd/admin/main.go

clean:
	rm -fv bin/*

run: bin/admin
	bin/admin $(SERVER_FLAGS)

clean-db:
	rm -f ./admin.sqlite

watch:
ifndef REFLEX
	$(error "reflex not found in PATH; you may need to run 'go get github.com/cespare/reflex' (be sure to do this outside of the api directory so it isn't added to go.mod)")
endif
	reflex -s -r '\.go$$' make run

test:
	go test ./admin
