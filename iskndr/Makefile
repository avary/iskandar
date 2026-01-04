.PHONY: build

build:
	go build -o bin/iskndr ./cmd/iskndr
	chmod +x bin/iskndr

test:
	go test ./... -v -race
