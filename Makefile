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
		echo "Please create a .env file with DISCORD_TOKEN"; \
		exit 1; \
	fi
	@source .env && go run cmd/bot/*.go -t $$DISCORD_TOKEN

# Build the binary
.PHONY: build
build:
	@echo "Building FACEIT Stats Bot..."
	@go build -o faceit-stats-bot cmd/bot/*.go
	@echo "Build complete: ./faceit-stats-bot"

# Test the code
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Generate SQL queries
.PHONY: queries
queries:
	@echo "Generating SQL queries..."
	@cd internal/infrastructure/db && sqlc generate

.PHONY: migrate
migrate:
	@echo "Migrating database..."
	@migrate -path internal/infrastructure/db/migrations -database "sqlite3://app.db" up

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
	@echo "  build      Build the bot binary"
	@echo "  test       Run tests"
	@echo "  queries    Generate SQL queries"
	@echo "  migrate    Migrate the database"
	@echo "  help       Show this help information"
	@echo ""
	@echo "Note: You need a .env file with DISCORD_TOKEN defined"
