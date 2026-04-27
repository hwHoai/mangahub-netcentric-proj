PROTO_FILES := $(shell find proto -name "*.proto")
export PATH := $(PATH):$(shell go env GOPATH)/bin

.PHONY: proto proto-tools proto-generate

proto:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILES)

proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto-generate:
	go generate ./proto
