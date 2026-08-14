APP_CONFIG ?= configs/config.local.yaml
MIGRATIONS_DIR ?= migrations
TEST_DATABASE_DSN ?=

.PHONY: help tidy fmt test test-integration run migrate-up migrate-down migrate-status docker-up docker-down docker-logs openapi

help:
	@echo "Available targets:"
	@echo "  make tidy              - tidy Go modules"
	@echo "  make fmt               - format Go code"
	@echo "  make test              - run unit tests"
	@echo "  make test-integration  - run integration tests with TEST_DATABASE_DSN"
	@echo "  make run               - run API server"
	@echo "  make migrate-up        - apply migrations"
	@echo "  make migrate-down      - roll back one migration"
	@echo "  make migrate-status    - show migration status"
	@echo "  make docker-up         - start local stack"
	@echo "  make docker-down       - stop local stack"
	@echo "  make docker-logs       - follow app logs"
	@echo "  make openapi           - check OpenAPI document exists"

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

test-integration:
	GOWEB_TEST_DATABASE_DSN="$(TEST_DATABASE_DSN)" go test -tags=integration ./...

run:
	go run ./cmd/api -config $(APP_CONFIG)

migrate-up:
	go run ./cmd/migrate -config $(APP_CONFIG) -dir $(MIGRATIONS_DIR) up

migrate-down:
	go run ./cmd/migrate -config $(APP_CONFIG) -dir $(MIGRATIONS_DIR) down

migrate-status:
	go run ./cmd/migrate -config $(APP_CONFIG) -dir $(MIGRATIONS_DIR) status

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

openapi:
	test -f api/openapi/openapi.yaml
