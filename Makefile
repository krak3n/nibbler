PATH := $(CURDIR)/bin:$(PATH)

DB_HOST            ?= localhost
DB_NAME            ?= nibbler
DB_USER            ?= postgres
DB_PORT            ?= 5432
DB_PASSWORD        ?= postgres
DB_SSL_MODE        ?= disable
DB_MIGRATIONS_PATH ?= $(CURDIR)/migrations

protoc: tools
	@$(MAKE) -C api/proto all

TOOLS += github.com/pressly/goose/cmd/goose
tools: $(TOOLS)

$(TOOLS): %:
	GOBIN=$(CURDIR)/bin go install $*

nibbler:
	go build -o bin/nibbler ./cmd/nibbler

image:
	docker build -f build/Dockerfile -t krak3n/nibbler:latest .

test:
	go test -v ./...

test-integration:
	go test -v -tags integration ./...

PSQL_URI = "host=$(DB_HOST) user=$(DB_USER) dbname=$(DB_NAME) sslmode=$(DB_SSL_MODE) password=$(DB_PASSWORD)"

migrate_up:
	goose -dir migrations postgres $(PSQL_URI) up

migrate_down:
	goose -dir migrations postgres $(PSQL_URI) down

up:
	docker-compose -f deployments/docker-compose.yml up -d

down:
	docker-compose -f deployments/docker-compose.yml down
