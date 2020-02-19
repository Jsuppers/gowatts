PROJECT_NAME := "gowatts"

build:
	go get -v -d ./...
	CGO_ENABLED=0 go build -o ./bin/${PROJECT_NAME} .