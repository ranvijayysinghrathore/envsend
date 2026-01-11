.PHONY: help build build-cli build-api build-worker test test-cli test-backend test-integration test-crypto clean docker-build docker-push deploy dev-setup migrate-up migrate-down lint fmt

# Variables
CLI_BINARY=bin/envsend
API_BINARY=bin/api-server
WORKER_BINARY=bin/worker
DOCKER_REGISTRY?=localhost:5000
VERSION?=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: build-cli build-api build-worker ## Build all binaries

build-cli: ## Build CLI binary
	@echo "Building CLI..."
	@mkdir -p bin
	go build $(LDFLAGS) -o $(CLI_BINARY) ./cli

build-api: ## Build API server binary
	@echo "Building API server..."
	@mkdir -p bin
	go build $(LDFLAGS) -o $(API_BINARY) ./backend/cmd/server

build-worker: ## Build worker binary
	@echo "Building worker..."
	@mkdir -p bin
	go build $(LDFLAGS) -o $(WORKER_BINARY) ./backend/cmd/worker

test: test-cli test-backend ## Run all tests

test-cli: ## Run CLI tests
	@echo "Running CLI tests..."
	cd cli && go test -v -race -cover ./...

test-backend: ## Run backend tests
	@echo "Running backend tests..."
	cd backend && go test -v -race -cover ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/integration/...

test-crypto: ## Run cryptography verification tests
	@echo "Running cryptography tests..."
	cd cli/crypto && go test -v -race -cover -bench=. ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf dist/
	go clean -cache -testcache

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	docker build -f Dockerfile.api -t $(DOCKER_REGISTRY)/envsend-api:$(VERSION) .
	docker build -f Dockerfile.worker -t $(DOCKER_REGISTRY)/envsend-worker:$(VERSION) .
	docker build -f Dockerfile.cli -t $(DOCKER_REGISTRY)/envsend-cli:$(VERSION) .

docker-push: ## Push Docker images to registry
	@echo "Pushing Docker images..."
	docker push $(DOCKER_REGISTRY)/envsend-api:$(VERSION)
	docker push $(DOCKER_REGISTRY)/envsend-worker:$(VERSION)
	docker push $(DOCKER_REGISTRY)/envsend-cli:$(VERSION)

dev-setup: ## Set up local development environment
	@echo "Setting up development environment..."
	@echo "Starting PostgreSQL, Redis, and MinIO..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	sleep 5
	@echo "Running migrations..."
	$(MAKE) migrate-up
	@echo "Development environment ready!"

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	migrate -path migrations -database "postgres://envsend:envsend@localhost:5432/envsend?sslmode=disable" up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	migrate -path migrations -database "postgres://envsend:envsend@localhost:5432/envsend?sslmode=disable" down

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@echo "Creating migration: $(NAME)"
	migrate create -ext sql -dir migrations -seq $(NAME)

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

deploy: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f deployments/kubernetes/

helm-install: ## Install via Helm
	@echo "Installing via Helm..."
	helm install envsend deployments/helm/ --create-namespace --namespace envsend

helm-upgrade: ## Upgrade Helm deployment
	@echo "Upgrading Helm deployment..."
	helm upgrade envsend deployments/helm/ --namespace envsend

install-cli: build-cli ## Install CLI to system
	@echo "Installing CLI..."
	cp $(CLI_BINARY) /usr/local/bin/envsend
	ln -sf /usr/local/bin/envsend /usr/local/bin/envreceive
	@echo "CLI installed successfully!"

run-api: build-api ## Run API server locally
	@echo "Starting API server..."
	./$(API_BINARY)

run-worker: build-worker ## Run worker locally
	@echo "Starting worker..."
	./$(WORKER_BINARY)
