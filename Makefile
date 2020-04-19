all: go-proto

go-proto:
	protoc --proto_path=./service --go_out=./service service.proto

build:
	go build -o build/draft-server cmd/server/*.go

test:
	go test -v ./... -cover