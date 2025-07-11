# Makefile for the Go project

# Binary name
BINARY_NAME=redis-like

# Default command
.DEFAULT_GOAL := help

## =============================================================================
## Help
## =============================================================================

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run         Run the server"
	@echo "  dev         Run the server with live-reloading (requires air)"
	@echo "  build       Build the binary for production"
	@echo "  install-air Install air (live-reloading tool)"

## =============================================================================
## Development
## =============================================================================

run:
	@echo "Running the server..."
	@go run ./cmd/server/main.go

dev:
	@echo "Running the server with live-reloading..."
	@air

test:
	@echo "Running tests..."
	@go test ./...

## =============================================================================
## Production
## =============================================================================

build:
	@echo "Building the binary for production..."
	@go build -o ./bin/$(BINARY_NAME) ./cmd/server/main.go

## =============================================================================
## Tools
## =============================================================================

install-air:
	@echo "Installing air..."
	@go install github.com/air-verse/air@latest

.PHONY: help run dev build install-air
