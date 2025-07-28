# Makefile for IU-K8s Backend API

.PHONY: help build run test clean generate deps docker-build docker-run

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Variables
BINARY_NAME=server
DOCKER_IMAGE=iu-k8s-api
PORT?=8080

# Build targets
build: ## Build the application
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) cmd/server/main.go
	@echo "Build complete: bin/$(BINARY_NAME)"

build-linux: ## Build for Linux (useful for Docker)
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux cmd/server/main.go
	@echo "Linux build complete: bin/$(BINARY_NAME)-linux"

# Run targets
run: ## Run the application
	@echo "Starting server on port $(PORT)..."
	@PORT=$(PORT) go run cmd/server/main.go

dev: ## Run in development mode with live reload (requires air)
	@echo "Starting development server..."
	@air

# Code generation
generate: ## Generate API code from OpenAPI spec
	@echo "Generating API code..."
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config openapi_config.yaml openapi.yaml
	@echo "Code generation complete"

# Dependencies
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# Testing
test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@go test -race -v ./...

# Linting and formatting
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p $(PORT):8080 --rm $(DOCKER_IMAGE)

docker-shell: ## Get shell access to Docker container
	@docker run -it --rm $(DOCKER_IMAGE) /bin/sh

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# OpenAPI
openapi-validate: ## Validate OpenAPI specification
	@echo "Validating OpenAPI spec..."
	@npx @apidevtools/swagger-parser validate openapi.yaml || echo "Note: Requires npm and @apidevtools/swagger-parser"

openapi-serve: ## Serve OpenAPI documentation
	@echo "Serving OpenAPI docs on http://localhost:8080..."
	@npx @redocly/cli preview-docs openapi.yaml || echo "Note: Requires npm and @redocly/cli"

# Installation helpers
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed"

# Database (for future use)
migrate-up: ## Run database migrations up (placeholder)
	@echo "Running migrations up..."
	@echo "Not implemented yet"

migrate-down: ## Run database migrations down (placeholder)
	@echo "Running migrations down..."
	@echo "Not implemented yet"

# Production helpers
deploy: build-linux ## Deploy to production (placeholder)
	@echo "Deploying to production..."
	@echo "Not implemented yet"

health-check: ## Check if server is healthy
	@echo "Checking server health..."
	@curl -f http://localhost:$(PORT)/health || echo "Server is not responding"
