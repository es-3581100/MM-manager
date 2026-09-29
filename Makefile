.PHONY: test vet smoke github-live-replay build

test:
	go test ./...

vet:
	go vet ./...

build:
	mkdir -p bin
	go build -trimpath -o bin/appdir-matrix ./cmd/appdir-matrix

smoke:
	bash scripts/smoke.sh

github-live-replay:
	bash scripts/github-live-replay.sh
