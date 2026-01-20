ENV?=local
include .env.$(ENV)

GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=calendar
MIGRATIONS_ORIGIN=migrations
DB_DRIVER=postgres
DOCKER_IMAGE=calendar-backend

.PHONY: start build test clean generate

# dev
dev air:
	air

start: build
	./$(BINARY_NAME)

build: clean
	GOEXPERIMENT=jsonv2 GOEXPERIMENT=greenteagc $(GOBUILD) -o $(BINARY_NAME)

clean:
	rm -f $(binary_name) && rm -rf tmp

build-docker:
	docker buildx build -f Dockerfile . -t $(DOCKER_IMAGE)

# generate
sqlc: sqlc-vet
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
