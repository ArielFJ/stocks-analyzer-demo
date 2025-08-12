.PHONY: build run test clean migrate dev

# Build the application
build:
	go build -o bin/stock-api main.go

# Run the application
run:
	go run main.go

# Run in development mode
dev:
	go run main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# =============================================================================
# DATABASE COMMANDS - Local Installation
# =============================================================================

# Start CockroachDB locally (if installed)
db-start:
	cockroach start-single-node --insecure --listen-addr=localhost:26257 --store=type=mem,size=0.25 --background

# Create database locally
db-create:
	cockroach sql --insecure --execute="CREATE DATABASE IF NOT EXISTS stockdb;"

# Stop CockroachDB locally
db-stop:
	cockroach quit --insecure

# Setup development environment (local installation)
setup: db-start db-create
	@echo "Development environment setup complete"
	@echo "Run 'make run' to start the server"

# =============================================================================
# DOCKER COMMANDS - For remaining docker-compose.yml
# =============================================================================

# Start database with Docker
docker-db-up:
	docker compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 10
	@echo "Database is ready"

# Stop database
docker-db-down:
	docker-compose down

# Start database and show logs
docker-db-logs:
	docker-compose up

# Setup development environment with Docker
docker-setup: docker-db-up
	@echo "Docker development environment setup complete"
	@echo "Run 'make run' to start the API server"

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Run security scan (requires gosec)
security:
	gosec ./...

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/stock-api main.go