include .env

GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=calendar
MIGRATIONS_ORIGIN=migrations
DB_DRIVER=postgres

.PHONY: start build test clean generate

# dev
dev:
	air

start: build
	./$(BINARY_NAME)

build: clean
	$(GOBUILD) -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME) && rm -rf tmp

# generate
sqlc:
	sqlc generate

sqlc-vet:
	sqlc vet

sqlc-diff:
	sqlc diff

generate: sqlc
	$(GORUN) github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=cfg.yaml api.yaml

# migrations
migrate-create:
	goose create $(name) sql -dir $(MIGRATIONS_ORIGIN) -s

migrate-up:
	goose -dir $(MIGRATIONS_ORIGIN) $(DB_DRIVER) $(DB_URL) up

migrate-down:
	goose -dir $(MIGRATIONS_ORIGIN) $(DB_DRIVER) $(DB_URL) down

migrate-reset:
	goose -dir $(MIGRATIONS_ORIGIN) $(DB_DRIVER) $(DB_URL) reset

migrate-version:
	goose -dir $(MIGRATIONS_ORIGIN) $(DB_DRIVER) $(DB_URL) version

migrate-status:
	goose -dir $(MIGRATIONS_ORIGIN) $(DB_DRIVER) $(DB_URL) status
