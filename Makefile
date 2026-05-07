# Windows-friendly Makefile
PROTO_FILES := $(wildcard proto/*/*.proto)

.PHONY: proto proto-tools proto-generate run-grpc run-api run-tcp run-udp run-ws run-all benchmark-tcp benchmark-udp

proto:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILES)

proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


build:
	go build ./...

run-grpc:
	cd cmd/grpc-server && go run main.go

run-api:
	cd cmd/api-server && go run main.go

run-tcp:
	cd cmd/tcp-server && go run main.go

run-udp:
	cd cmd/udp-server && go run main.go

run-ws:
	cd cmd/websocket-server && go run main.go

# Benchmarks
benchmark-tcp:
	go run internal/benchmarks/tcp_stress/main.go -conns 2000

benchmark-udp:
	go run internal/benchmarks/udp_reliability/main.go -n 2000

# Note: On Windows, it's highly recommended to run servers in separate terminals.
# This command will launch 5 separate PowerShell windows.
run-all:
	powershell -ExecutionPolicy Bypass -File run-all.ps1
