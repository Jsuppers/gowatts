PROJECT_NAME := "gowatts"

.PHONY: build lint test

all: dependencies

dependencies:
	export GO111MODULE="on"
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

coverage: generate
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out
	goveralls -coverprofile=cover.out -service=travis-ci -repotoken ${COVERALLS_TOKEN}
	rm cover.out

coverhtml: generate
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o coverage.html
	rm cover.out

