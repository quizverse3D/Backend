run-all:
	go run ./cmd/authgateway

run-authgateway:
	go run ./cmd/authgateway

build:
	go build -o bin/authgateway ./cmd/authgateway