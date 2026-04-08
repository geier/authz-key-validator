.PHONY: build test clean

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -v -race

lint:
	golangci-lint run

docker-build:
	docker build -t authz-key-validator:latest .
