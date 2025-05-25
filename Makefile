custom-gcl-build:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint custom

lint: custom-gcl-build
	bin/custom-gcl run ./...

test:
	mkdir -p data
	go test -v -covermode=atomic -coverprofile=data/coverage.out ./...
	grep -v "mock" data/coverage.out > data/coverage.out.tmp
	go tool cover -html data/coverage.out.tmp -o data/coverage.html
	open data/coverage.html

build-api: lint test
	mkdir -p ./bin && go build -o bin/monitoring-api cmd/monitoring-api/main.go

.PHONY: run build-api test lint custom-gcl-build

.DEFAULT_GOAL:=run
