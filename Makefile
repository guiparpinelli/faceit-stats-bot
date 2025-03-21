# Default shell
SHELL := /bin/bash

# Default target
.PHONY: default
default: run

# Load environment variables and run the bot
.PHONY: run
run:
	@echo "Starting FACEIT Stats Bot..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found!"; \
		echo "Please create a .env file with DISCORD_TOKEN and FACEIT_API_KEY"; \
		exit 1; \
	fi
	@source .env && go run main.go -t $$DISCORD_TOKEN -k $$FACEIT_API_KEY

# Build the binary
.PHONY: build
build:
	@echo "Building FACEIT Stats Bot..."
	@go build -o faceit-stats-bot main.go
	@echo "Build complete: ./faceit-stats-bot"

# Test the code
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Show help information
.PHONY: help
help:
	@echo "FACEIT Stats Discord Bot Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run        Run the bot (default)"
	@echo "  test       Run tests"
	@echo "  help       Show this help information"
	@echo ""
	@echo "Note: You need a .env file with DISCORD_TOKEN and FACEIT_API_KEY defined"
