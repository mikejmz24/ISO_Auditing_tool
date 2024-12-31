# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@templ generate
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests/... -v

# Test short results
test-short:
	@go test -v ./tests/... 2>&1 | grep -v "warning" | grep -E "=== RUN |--- (PASS|FAIL|SKIP)" | grep -v "=== RUN" | \
		sed -e 's/--- PASS/\x1b[32m--- PASS\x1b[0m/' \
		    -e 's/--- FAIL/\x1b[31m--- FAIL\x1b[0m/' \
		    -e 's/--- SKIP/\x1b[33m--- SKIP\x1b[0m/' | \
				awk '{if(index($$0, "/") > 0) {split($$0, parts, "/"); split(parts[1], colin_sub, ": "); print colin_sub[1] ": " parts[length(parts)]} else {print $$0}}'

# Run the database migrations
migrate:
	@go run cmd/api/main.go migrate

# Seed the database
seed:
	@go run cmd/api/main.go seed

# Truncate the database tables
truncate:
	@go run cmd/api/main.go truncate

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
