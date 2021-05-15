.PHONY: clean clean_db test run watch run_clean_db

GO_SRCS := $(shell find . -type f -name '*.go' -and -not -name '*_test.go')
JS_SRCS := $(shell find javascript -type f -name '*.js' -or -name '*.css' -or -name '*.vue')
REFLEX := $(shell command -v reflex 2> /dev/null)

all: bin/admin static/admin/admin.bundle.js

bin/admin: $(GO_SRCS)
	go build -o $@ cmd/admin/main.go

static/admin/admin.bundle.js: node_modules/.timestamp $(JS_SRCS)
	pnpm run build
	touch $@

node_modules/.timestamp:
	pnpm install
	touch node_modules/.timestamp

clean:
	rm node_modules/.timestamp
	rm -fv bin/*
	rm static/*.js
	rm static/*.css
	rm static/*.LICENSE.txt
	rm static/admin/*.js
	rm static/admin/*.css
	rm static/admin/*.LICENSE.txt

run: SERVER_FLAGS=-assethost http://localhost:9000
run: bin/admin
	bin/admin $(SERVER_FLAGS)

clean-db:
	rm -fv ./data/*.sqlite

migrate: bin/admin
	bin/admin -migrate

clean-migrate: clean_db migrate

watch:
ifndef REFLEX
	$(error "reflex not found in PATH; you may need to run 'go get github.com/cespare/reflex' (be sure to do this outside of the api directory so it isn't added to go.mod)")
endif
	reflex -s -r '\.go$$' -R 'node_modules|javascript|static' make run

test:
	go test ./admin
