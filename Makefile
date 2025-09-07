run-all:
	go run ./cmd/authgateway & \
	go run ./cmd/user & \
	go run ./cmd/room

run-user:
	go run ./cmd/user
run-authgateway:
	go run ./cmd/authgateway
run-room:
	go run ./cmd/room

protoc:
	protoc --go_out=. --go-grpc_out=. proto/user.proto &\
	protoc --go_out=. --go-grpc_out=. proto/room.proto

build:
	go build -o bin/authgateway ./cmd/authgateway & \
	go build -o bin/user ./cmd/user