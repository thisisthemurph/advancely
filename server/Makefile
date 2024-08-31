.PHONY: migrate migrate-down

api:
	go run ./cmd/api

build:
	go build -tags netgo -ldflags '-s -w' -o ./bin/app ./cmd/api/.

migrate:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down