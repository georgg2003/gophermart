# Gophermart

A loyalty points accumulation system — REST API service built with Go.

Users register, submit order numbers, and earn points via an external accrual service. Points can be withdrawn to pay for new orders.

## Features

- User registration and authentication (JWT)
- Order submission and status tracking
- Async polling of external accrual service
- Balance management: accruals and withdrawals
- Idempotent order processing

## Stack

Go · PostgreSQL · Chi · JWT · golang-migrate

## Running locally

```bash
# Start PostgreSQL
docker compose up -d postgres

# Run the service
go run ./cmd/gophermart \
  -d "postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable" \
  -a "localhost:8080" \
  -r "http://localhost:8088"
```

## API

API endpoints can be found in [api/swagger.yaml](api/swagger.yaml).

Code is generated from this file with oapi-codegen

## Configuration

| Flag | Env | Description |
|------|-----|-------------|
| `-a` | `RUN_ADDRESS` | HTTP server address |
| `-d` | `DATABASE_URI` | PostgreSQL DSN |
| `-r` | `ACCRUAL_SYSTEM_ADDRESS` | Accrual service URL |
