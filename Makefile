# Имя бинаря
BINARY := gophermart

# Путь к миграциям
MIGRATIONS_DIR := ./migrations

# Параметры базы
DB_URL := "postgres://gophermart:password@localhost:5432/gophermart?sslmode=disable"

# Go параметры
GO := go
GOFLAGS :=
GOMOD := $(shell go env GOMOD)
GOPATH := $(shell go env GOPATH)
PKG := ./...

export PATH := $(GOPATH)/bin:$(PATH)

.PHONY: all tidy deps build run test mock migrate-up migrate-down clean

all: build

## ---------------------------
## Dependencies
## ---------------------------

tidy:
	$(GO) mod tidy

deps:
	$(GO) install go.uber.org/mock/mockgen@latest
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## ---------------------------
## Build & Run
## ---------------------------

build:
	$(GO) build -o bin/$(BINARY) $(PKG)

run:
	$(GO) run $(PKG) -d $(DB_URL)

clean:
	rm -rf bin

## ---------------------------
## Testing
## ---------------------------

test:
	$(GO) test $(PKG) -cover -v

coverage: 
	$(GO) test $(PKG) -covermode=count -coverpkg=$(PKG) -coverprofile=coverage.out
	grep -vE "(mock\.go|/mock/)" coverage.out > coverage.filtered.out
	mv coverage.filtered.out coverage.out
	$(GO) tool cover -func=coverage.out

## ---------------------------
## Generate mocks and stuff
## ---------------------------

generate:
	$(GO) generate $(PKG)

## ---------------------------
## Migrations
## ---------------------------

migrate-up:
	$(GOPATH)/bin/migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) up

migrate-down:
	$(GOPATH)/bin/migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) down

migrate-new:
	$(GOPATH)/bin/migrate -database $(DB_URL) create -ext sql -seq -digits 6 -dir $(MIGRATIONS_DIR) $(NAME)

migrate-force:
	$(GOPATH)/bin/migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) force $(VER)