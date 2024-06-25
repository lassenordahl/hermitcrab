# Makefile for HermitCrab project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=hermitcrab

# Mockgen
MOCKGEN=mockgen

# Colors
CYAN=\033[0;36m
RESET=\033[0m

.PHONY: all build test clean run deps generate generate-mocks

all: test build

build:
	@echo "🏗️  Building..."
	@$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

test:
	@echo "🧪 Testing..."
	@$(GOTEST) -v ./...

clean:
	@echo "🧹 Cleaning..."
	@rm -f $(BINARY_NAME)
	@$(GOCMD) clean

run: build
	@echo "🚀 Running..."
	@./$(BINARY_NAME)

deps:
	@echo "📦 Installing dependencies..."
	@$(GOGET) -v -t -d ./...
	@$(GOGET) github.com/golang/mock/mockgen@v1.6.0

generate-mocks:
	@echo "🎭 Generating mocks..."
	@$(MOCKGEN) -source=pkg/hermitcrab/bucket/interface.go -destination=pkg/hermitcrab/bucket/mock_bucket_manager.go -package=bucket

generate: generate-mocks
	@echo "✨ Generated all artifacts"

help:
	@echo "$(CYAN)Available commands:$(RESET)"
	@echo "  make build         - Build the binary"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make run           - Build and run the binary"
	@echo "  make deps          - Install dependencies"
	@echo "  make generate-mocks - Generate mock for BucketManager"
	@echo "  make generate      - Generate all artifacts"
	@echo "  make help          - Show this help message"
