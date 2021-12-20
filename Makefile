
.DEFAULT_GOAL := run
SHELL := bash

clean:
	rm -f bin/jqplay
	rm -f */rice-box.go
	rm -rf */statik/
	npm exec grunt clean

build: clean
	set -eo pipefail
	npm install
	# cd server && rice embed-go && cd ..
	go generate ./...
	go build -ldflags="-X 'main.GinMode=release'" -o ./bin/jqplay ./cmd/jqplay
	echo "built bin/jqplay"

run:
	go run ./cmd/jqplay -verbose

deps:
	cat cmd/tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

