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

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/user/register` | Register new user |
| POST | `/api/user/login` | Authenticate |
| POST | `/api/user/orders` | Submit order number |
| GET | `/api/user/orders` | List orders with statuses |
| GET | `/api/user/balance` | Current balance |
| POST | `/api/user/balance/withdraw` | Withdraw points |
| GET | `/api/user/withdrawals` | Withdrawal history |

## Configuration

| Flag | Env | Description |
|------|-----|-------------|
| `-a` | `RUN_ADDRESS` | HTTP server address |
| `-d` | `DATABASE_URI` | PostgreSQL DSN |
| `-r` | `ACCRUAL_SYSTEM_ADDRESS` | Accrual service URL |
