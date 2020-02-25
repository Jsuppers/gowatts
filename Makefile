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

