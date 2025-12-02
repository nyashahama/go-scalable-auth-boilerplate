.PHONY: help build run test clean migrate-up migrate-down sqlc docker-build docker-up docker-down

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/api cmd/api/main.go

run: ## Run the application
	go run cmd/api/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-integration: ## Run integration tests
	go test -v -tags=integration ./tests/integration/...

clean: ## Clean build artifacts
	rm -rf bin/ coverage.out coverage.html

sqlc: ## Generate sqlc code
	sqlc generate

migrate-up: ## Run database migrations up
	migrate -path migrations -database "$(DB_URL)" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "$(DB_URL)" down

migrate-create: ## Create a new migration (usage: make migrate-create name=add_users_table)
	migrate create -ext sql -dir migrations -seq $(name)

docker-build: ## Build docker image
	docker build -t user-auth-app:latest .

docker-up: ## Start docker compose
	docker-compose up -d

docker-down: ## Stop docker compose
	docker-compose down

docker-logs: ## View docker logs
	docker-compose logs -f

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	gofumpt -l -w .

tidy: ## Tidy dependencies
	go mod tidy

generate-jwt: ## Generate a secure JWT secret
	@openssl rand -base64 32

.DEFAULT_GOAL := help