.PHONY: build test lint

build:
	go build -o redifire .

test:
	go test ./...

lint:
	go vet ./...
