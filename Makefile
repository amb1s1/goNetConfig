BINARY_NAME := gonetconfig
PROTO_DIR := proto

.PHONY: all build run test test-v vet clean proto entity-init help

all: vet test build ## Run vet, test, and build

build: ## Build the application binary
	go build -o $(BINARY_NAME) .

run: ## Run the application
	go run .

test: ## Run all tests
	go test ./...

test-v: ## Run all tests with verbose output
	go test -v ./...

vet: ## Run go vet on all packages
	go vet ./...

clean: ## Remove build artifacts
	rm -f $(BINARY_NAME)

proto: ## Regenerate gRPC code from proto files
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/service.proto

entity-init: ## Copy template entity config to ~/.gonetconfig/
	mkdir -p $(HOME)/.gonetconfig
	cp configs/entity.yml $(HOME)/.gonetconfig/entity.yml
	@echo "Entity config created at $(HOME)/.gonetconfig/entity.yml"

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
