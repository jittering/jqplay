
.DEFAULT_GOAL := run
SHELL := bash

clean:
	rm -f bin/jqplay
	rm -f */rice-box.go
	rm -rf */statik/
	rm -rf web/public/build/

build: clean
	set -eo pipefail
	cd web && npm install && npm run build && rm -f public/build/*.map && cd ..
	go generate ./...
	go build -ldflags="-X 'main.GinMode=release'" -o ./bin/jqplay ./cmd/jqplay
	echo "built bin/jqplay"

statik:
	go generate ./jq/...

run: statik
	rm -f ./server/rice-box.go
	go run ./cmd/jqplay -verbose -no-open

deps:
	cat cmd/tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

release:
	goreleaser release --rm-dist --parallelism 1 --skip-validate
