all: go-proto

go-proto:
	protoc --proto_path=./service --go_out=plugins=grpc:./service service.proto

# js-proto:
# 	mkdir -p ./service/js
# 	protoc --proto_path=./service --js_out=import_style=commonjs:./service/js service.proto

build:
	go build -o build/draft-server cmd/server/*.go

fmt:
	gofmt -w ./..

test:
	go test -v ./... -cover