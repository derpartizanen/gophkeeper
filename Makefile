DEFAULT_GOAL := help

BUILD_FOLDER = dist
CRT_FOLDER = ssl/ca

# Build info
CLIENT_VERSION ?= 0.1.0

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: proto
proto: ## Generate gRPC protobuf bindings
	go ./proto/generate.go

.PHONY: ssl
ssl: ## Generate SSL certificates for secure communications
	./scripts/gen-ca
	./scripts/issue-crt

.PHONY: keeperd ## Build keeperd
keeperd:
	go build -o $(BUILD_FOLDER)/$@ cmd/$@/*.go

.PHONY: keeperctl ## Build keeperctl
keeperctl:
	./scripts/build-client $(CLIENT_VERSION)

.PHONY: all ## Build whole product.
all: keeperd keeperctl

.PHONY: download
download: ## Download go.mod dependencies
	echo Downloading go.mod dependencies
	go mod download

.PHONY: up
up: ## Run the project in docker compose
	docker compose -f deploy/docker-compose.yaml up -d --build

.PHONY: down
down: ## Stop the running project and destroy containers
	docker compose -f deploy/docker-compose.yaml down

.PHONY: clean
clean: down
	rm -rf $(BUILD_FOLDER) $(CRT_FOLDER)

.PHONY: test
test: ## Run unit tests
	@go test -v -cover ./...
