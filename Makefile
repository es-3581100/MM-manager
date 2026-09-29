.PHONY: test vet smoke build

test:
	go test ./...

vet:
	go vet ./...

build:
	mkdir -p bin
	go build -trimpath -o bin/appdir-matrix ./cmd/appdir-matrix

smoke:
	bash scripts/smoke.sh
