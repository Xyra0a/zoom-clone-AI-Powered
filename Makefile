# Makefile for Zoom Clone Backend

.PHONY: install run build test clean help

help:
	@echo "Available commands:"
	@echo "  make install   - Download dependencies"
	@echo "  make run       - Run the server"
	@echo "  make build     - Build the binary"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make lint      - Run linter"
	@echo "  make fmt       - Format code"

install:
	@echo "Downloading dependencies..."
	go mod download

run: install
	@echo "Starting server..."
	go run main.go

build: install
	@echo "Building binary..."
	go build -o zoom-clone-backend .

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -f zoom-clone-backend
	go clean

lint:
	@echo "Running linter..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
