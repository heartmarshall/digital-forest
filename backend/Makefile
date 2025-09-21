# Digital Forest Backend Makefile

.PHONY: help build run test test-unit test-integration test-e2e test-coverage clean docker-build docker-run

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-e2e       - Run end-to-end tests only"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run application in Docker"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/app ./cmd/app

# Run the application
run:
	@echo "Running application..."
	go run ./cmd/app

# Run all tests
test:
	@echo "Running all tests..."
	go test -v ./...

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v -short ./internal/usecase/... ./internal/transport/http/handlers/... ./internal/transport/http/dto/...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	go test -v -run "Integration" ./internal/repository/... ./internal/transport/http/...

# Run end-to-end tests only
test-e2e:
	@echo "Running end-to-end tests..."
	go test -v -run "E2E" ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./...

# Run tests with benchmarks
test-bench:
	@echo "Running benchmarks..."
	go test -v -bench=. -benchmem ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t digital-forest-backend .

# Run application in Docker
docker-run:
	@echo "Running application in Docker..."
	docker-compose up --build

# Run tests in Docker
docker-test:
	@echo "Running tests in Docker..."
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	# Add migration commands here when you have a migration tool

migrate-down:
	@echo "Rolling back database migrations..."
	# Add rollback commands here when you have a migration tool

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	$(MAKE) deps
	$(MAKE) fmt
	@echo "Development environment ready!"

# CI/CD pipeline
ci:
	@echo "Running CI pipeline..."
	$(MAKE) deps
	$(MAKE) fmt
	$(MAKE) lint
	$(MAKE) test-coverage
	$(MAKE) build
	@echo "CI pipeline completed successfully!"