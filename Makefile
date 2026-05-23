DB_DSN ?= host=localhost port=5432 user=api_user password=api_pass dbname=api_structure sslmode=disable TimeZone=UTC
GOOSE ?= goose
COMPOSE ?= docker compose

.PHONY: up down ps run test migrate-up migrate-down migrate-status

up:
	$(COMPOSE) up -d postgres

down:
	$(COMPOSE) down

ps:
	$(COMPOSE) ps

run:
	DB_DSN="$(DB_DSN)" go run ./cmd/server

test:
	go test ./...

migrate-up:
	$(GOOSE) -dir ./migrations postgres "$(DB_DSN)" up

migrate-down:
	$(GOOSE) -dir ./migrations postgres "$(DB_DSN)" down

migrate-status:
	$(GOOSE) -dir ./migrations postgres "$(DB_DSN)" status
