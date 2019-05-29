TEST_LOG_LEVEL := $(or $(LOG_LEVEL), fatal)

build:
	go build main.go

fmt:
	go fmt ./...

run:
	go run main.go

test:
	LOG_LEVEL=$(TEST_LOG_LEVEL) go test ./...

help:
	@go run main.go -help