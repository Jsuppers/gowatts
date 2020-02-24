PROJECT_NAME := "gowatts"

.PHONY: build lint

generate: 
	go generate ./...

build:
	go get -v -d ./...
	CGO_ENABLED=0 go build -o ./bin/${PROJECT_NAME} .

lint: generate 
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.22.2
	bin/golangci-lint run
