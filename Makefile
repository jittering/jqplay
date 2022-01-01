
.DEFAULT_GOAL := go
SHELL := bash

clean-web:
	rm -rf web/public/build/

clean: clean-web
	rm -rf dist/
	rm -f bin/jqplay
	rm -f */rice-box.go
	rm -rf */statik/

build: clean
	set -eo pipefail
	cd web && npm install && npm run build && rm -f public/build/*.map && cd ..
	go generate ./...
	mkdir -p dist
	go build -ldflags="-X 'main.GinMode=release'" -o ./dist/jqplay ./cmd/jqplay
	echo "built dist/jqplay"

statik:
	go generate ./jq/...

run: statik
	rm -f ./server/rice-box.go
	mkdir -p web/public/build
	go run ./cmd/jqplay -verbose -no-open

deps:
	cat cmd/tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

release:
	goreleaser release --rm-dist --parallelism 1 --skip-validate

web: clean-web
	cd web && npm run dev

go:
	forego start

.PHONY: web build clean statik run release
