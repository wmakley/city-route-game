.PHONY: clean clean_db test run watch run_clean_db

GO_SRCS := $(shell find . -type f -name '*.go' -and -not -name '*_test.go')

REFLEX := $(shell command -v reflex 2> /dev/null)

hansa: $(GO_SRCS)
	go build -o $@

clean:
	rm -f hansa

run: hansa
	./hansa $(SERVER_FLAGS)

clean_db:
	rm -f ./hansa.sqlite

run_clean_db: clean_db run

watch:
ifndef REFLEX
	$(error "reflex not found in PATH; you may need to run 'go get github.com/cespare/reflex' (be sure to do this outside of the api directory so it isn't added to go.mod)")
endif
	reflex -s -r '\.go$$' make run

test:
	go test .
