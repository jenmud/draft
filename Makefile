.PHONY: proto
go-proto:
	protoc --proto_path=./proto --micro_out=./proto --go_out=./proto service.proto

python-proto:
	python3 -m grpc_tools.protoc --proto_path=./proto --python_out=./python --grpc_python_out=./python service.proto

.PHONY: build
build:
	go build -o build/draft-server cmd/server/*.go

.PHONY: test
test:
	go test -v ./... -cover