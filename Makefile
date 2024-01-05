.PHONY: generate mini-grpc-gateway fmt run

GRG_SRCS := $(patsubst ./%,%,$(shell find . -path "*/grpc-rest-gateway/*.go"))

generate: mini-grpc-gateway
	@echo "Generating new files..."
	@protoc --go_out=./pb/post --go_opt=paths=source_relative --go-grpc_out=./pb/post --go-grpc_opt=paths=source_relative service.proto 
	@protoc --mini-grpc-gateway_out=./pb/post service.proto && echo "Generated successfully!"

mini-grpc-gateway:
	@echo "Building protoc-gen-mini-grpc-gateway..."
	@go build -o $(GOPATH)/bin/protoc-gen-mini-grpc-gateway main.go && echo "Built successful!"

fmt: $(GRG_SRCS)
	@echo "Formatting files..."
	@gofmt -s -w $(GRG_SRCS) && goimports -w $(GRG_SRCS) && echo "Formatted successfully!"

backend: generate fmt
	go run ./server

run: generate
	go run ./grpc-rest-gateway
	
