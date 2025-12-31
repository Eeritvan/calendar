include .env

GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=calendar
MIGRATIONS_ORIGIN=migrations

.PHONY: start build test clean generate

start: build
	./$(BINARY_NAME)

dev:
	air

build: clean
	$(GOBUILD) -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME) && rm -rf tmp

sqlc:
	sqlc generate

sqlc-vet:
	sqlc vet

sqlc-diff:
	sqlc diff

generate: sqlc
	$(GORUN) github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=cfg.yaml api.yaml

migrate:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_ORIGIN) up

migrate-undo:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_ORIGIN) down 1

migrate-version:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_ORIGIN) version

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_ORIGIN) $(name)

migrate-reset:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_ORIGIN) drop
