run-all:
	go run ./cmd/authgateway & \
	go run ./cmd/user

run-user:
	go run ./cmd/user
run-authgateway:
	go run ./cmd/authgateway

protoc:
	protoc --go_out=. --go-grpc_out=. proto/user.proto

build:
	go build -o bin/authgateway ./cmd/authgateway & \
	go build -o bin/user ./cmd/user