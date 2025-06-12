run-gateway:
	go run ./cmd/gateway
run-user:
	go run ./cmd/user

build:
	go build -o bin/gateway ./cmd/gateway
	go build -o bin/user ./cmd/user