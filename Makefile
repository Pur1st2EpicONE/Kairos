.PHONY: all up down reset migrate-up migrate-down local lint test postgres rabbit app_logs postgres_logs rabbit_logs_logs queues .env .env.example help
.POSIX:
.SILENT:

-include .env.example .env

GOOSE_CMD = goose -dir ./migrations postgres "user=${DB_USER} password=${DB_PASSWORD} dbname=kairos-db host=localhost port=5433 sslmode=disable"

all: up

up:	
	if [ ! -f .env ] && [ ! -f .env.example ]; then \
		echo "Missing environment file: .env or .env.example is required."; \
		exit 1; \
	fi
	if [ ! -f .env ]; then cat .env.example > .env; fi
	if [ ! -f config.yaml ]; then cp ./configs/config.full.yaml ./config.yaml; fi
	if [ ! -f docker-compose.yaml ]; then cp ./deployments/docker-compose.full.yaml ./docker-compose.yaml; fi
	if [ ! -f Dockerfile ]; then cp ./deployments/Dockerfile ./Dockerfile; fi
	COMPOSE_BAKE=true docker compose up -d postgres rabbitmq app
	rm -f Dockerfile

down:
	docker compose down 2>/dev/null || true 
	rm -f Dockerfile docker-compose.yaml config.yaml

reset:
	docker volume rm kairos_postgres_data

migrate-up:
	@if command -v goose > /dev/null 2>&1; then $(GOOSE_CMD) up; else echo "You need Goose migration tool to use this command!"; fi

migrate-down:
	@if command -v goose > /dev/null 2>&1; then $(GOOSE_CMD) down; else echo "You need Goose migration tool to use this command!"; fi

local:
	if [ ! -f .env ]; then cat .env.example > .env; fi 
	if [ ! -f config.yaml ]; then cp ./configs/config.dev.yaml ./config.yaml; fi 
	if [ ! -f docker-compose.yaml ]; then cp ./deployments/docker-compose.dev.yaml ./docker-compose.yaml; fi
	COMPOSE_BAKE=true docker compose up -d postgres rabbitmq
	until docker exec rabbitmq rabbitmqctl status > /dev/null 2>&1; do sleep 0.5; done
	bash -c 'trap "exit 0" INT; go run ./cmd/kairos/main.go'

lint:
	golangci-lint run ./...

test:
	if [ ! -f .env ]; then cat .env.example > .env	; fi 
	if [ ! -f config.yaml ]; then cp ./configs/config.test.yaml ./config.yaml; fi 
	if [ ! -f docker-compose.yaml ]; then cp ./deployments/docker-compose.test.yaml ./docker-compose.yaml; fi
	COMPOSE_BAKE=true docker compose -f docker-compose.yaml up -d postgres-test
	until docker exec postgres-test pg_isready -U ${DB_USER} -d kairos_test > /dev/null 2>&1; do sleep 0.5; done
	echo "Running tests, please be patient (≈2 min)"
	COMPOSE_BAKE=true docker compose -f docker-compose.yaml run --rm app-test > .temp 2>/dev/null
	cat .temp; rm -f .temp
	docker compose -f docker-compose.yaml down -v > /dev/null 2>&1
	rm -f docker-compose.yaml config.yaml .env

postgres:
	docker compose exec postgres psql -U ${DB_USER} -d kairos-db

rabbit:
	docker compose exec rabbitmq bash

app_logs:
	docker compose logs --tail 10 app

postgres_logs:
	docker compose logs --tail 10 postgres

rabbit_logs:
	docker compose logs --tail 10 rabbitmq

queues:
	docker compose exec rabbitmq rabbitmqctl list_queues

.env:
	@:

help:
	@echo " ———————————————————————————————————————————————————————————————————————————————————— "
	@echo "| up             | Start all services (postgres, rabbitmq, app) in background        |"
	@echo "| down           | Stop and remove all containers, networks, and temporary files     |"
	@echo "| reset          | Remove postgres Docker volume                                     |"
	@echo "| migrate-up     | Apply all database migrations (Goose migration tool required)     |"
	@echo "| migrate-down   | Rollback all database migrations (Goose migration tool required)  |"
	@echo "| local          | Start local dev environment (go 1.25.1 required)                  |"
	@echo "| lint           | Run golangci-lint                                                 |"
	@echo "| test           | Run unit and integration tests                                    |"
	@echo "| postgres       | Open psql shell inside postgres container                         |"
	@echo "| rabbit         | Open shell inside rabbitmq container                              |"
	@echo "| app_logs       | Show last 10 lines of app logs                                    |"
	@echo "| postgres_logs  | Show last 10 lines of postgres logs                               |"
	@echo "| rabbit_logs    | Show last 10 lines of rabbitmq logs                               |"
	@echo "| queues         | List queues in rabbitmq                                           |"
	@echo " ———————————————————————————————————————————————————————————————————————————————————— "

.DEFAULT:
	@echo " No rule to make target '$@'. Available make targets:"
	@$(MAKE) --no-print-directory help