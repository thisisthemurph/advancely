.PHONY: migrate migrate-down

test:
	go test -v ./...

api:
	go run ./cmd/api

build:
	go build -tags netgo -ldflags '-s -w' -o ./bin/app ./cmd/api/.

lint:
	golangci-lint run

migrate:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down