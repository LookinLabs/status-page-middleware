include .env

DEFAULT_GOAL := build

build:
	go build -o status_page_middleware main.go

run:
	go run main.go

mod-vendor:
	go mod vendor

linter:
	@golangci-lint run

gosec:
	@gosec -quiet ./...

tests:
	go test -v ./tests/...

validate: linter gosec tests