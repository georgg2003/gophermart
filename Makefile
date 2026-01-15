# Имя бинаря
BINARY := gophermart

# Путь к миграциям
MIGRATIONS_DIR := ./migrations

# Параметры базы
DB_URL := "postgres://gophermart:password@localhost:5435/gophermart?sslmode=disable"

# Go параметры
GO := go
GOFLAGS :=
GOMOD := $(shell go env GOMOD)
GOPATH := $(shell go env GOPATH)
PKG := ./...

export PATH := $(GOPATH)/bin:$(PATH)

.PHONY: all tidy deps build run test mock migrate-up migrate-down clean accrual

all: build

## ---------------------------
## Dependencies
## ---------------------------

tidy:
	$(GO) mod tidy

## ---------------------------
## Build & Run
## ---------------------------

build:
	$(GO) build -o bin/$(BINARY) $(PKG)

accrual:
	./cmd/accrual/accrual_darwin_arm64 -a :8000

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
	grep -vE "(mock\.go|/mock/|/*.gen.go)" coverage.out > coverage.filtered.out
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