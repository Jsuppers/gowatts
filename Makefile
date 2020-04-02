PROJECT_NAME := "gowatts"

.PHONY: build lint test

all: dependencies

dependencies:
	@go get -v -d ./...

generate: 
	go generate ./...

build:
	CGO_ENABLED=0 go build -o ./bin/${PROJECT_NAME} .

lint: generate 
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.22.2
	bin/golangci-lint run

test: generate
	go test -race ./...

coverage: 
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out
	rm cover.out

coverhtml: generate
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o coverage.html
	rm cover.out

