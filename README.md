# Salary Advance Loan Service

Backend service for salary advance loan challenge implemented in Go.

## Features

- JWT-based authentication and authorization
- Role-based access control: `admin`, `uploader`
- Login rate limiting (failed-attempt window)
- Customer data validation with detailed batch logs
- Verified customer persistence in PostgreSQL
- Transaction mapping, synthetic transaction generation, and customer rating (1-10)
- Idempotent admin seeding on startup
- Unit and integration-style tests for auth, handlers, repositories, and services

## Project Structure

- `cmd/server.go`: app entrypoint and route wiring
- `cmd/utils.go`: environment/config helpers, database bootstrap, admin seed logic
- `internal/interfaces/http`: handlers and auth middleware
- `internal/interfaces/dto`: file readers for customer/transaction/sample data
- `internal/services`: auth, validation, and rating business logic
- `internal/repository`: PostgreSQL repository implementation
- `internal/testutil`: shared PostgreSQL test helpers
- `data`: input files

## Security Measures

- Passwords are hashed with bcrypt before storage
- JWT access tokens are signed with HMAC SHA-256
- Token claims include role and expiration
- Token denylist is used for logout invalidation
- Login attempts are rate-limited per source key (IP + username)
- Protected routes require valid bearer token
- Admin-only routes enforce role checks

## Run Locally

1. Set environment values (optional):

```bash
export PORT=8080
export JWT_SECRET='secret'
export ADMIN_USERNAME='admin'
export ADMIN_PASSWORD='Admin@1234'
export POSTGRES_DB='loan_service'
export POSTGRES_USER='loan_user'
export POSTGRES_PASSWORD='12345678'
export DB_HOST='localhost'
export DB_PORT='5432'
export DB_SSLMODE='disable'
export CUSTOMERS_FILE='data/customers.json'
export TRANSACTIONS_FILE='data/transactions.json'
export SAMPLE_FILE='data/sample_customers.csv'
```

Note: when running inside Docker Compose, `DB_HOST` should be `postgres`.

2. Start service:

```bash
go run ./cmd
```

## Docker

```bash
docker compose up --build
```

## API Endpoints

- `POST /auth/register` (public): register uploader
- `POST /auth/login` (public): returns JWT token
- `POST /auth/register-admin` (admin only): create admin
- `POST /auth/logout` (authenticated): invalidate token
- `GET /api/validate` (authenticated): run customer sample validation and persist verified customers
- `GET /api/verified-customers` (authenticated): list verified customers
- `GET /api/customer-ratings` (admin only): calculate and return customer ratings

Use `Authorization: Bearer <token>` for protected endpoints.

## Validation Rules

For each sample customer record:

- account number must match format `^\\d{10,13}$`
- account number must exist in canonical customer list
- customer name must match canonical name (trimmed, case-insensitive)

Per record output includes:

- `record_index`
- `verified`
- `errors`
- `normalized_record` (for verified records)

Batch output includes:

- `batch_id`
- `contains_errors`
- `failure_reason` for unverified batches

## Rating Calculation

Customer rating is calculated with weighted components:

- transaction count score (25%)
- total transaction volume score (30%)
- duration score between first and last transaction (20%)
- balance stability score using coefficient of variation (25%)

Formula:

```text
rating = 0.25*count_score + 0.30*volume_score + 0.20*duration_score + 0.25*stability_score
```

Final rating is clamped to range `[1, 10]`.

If a customer has no transactions in `transactions.json`, synthetic transactions are generated in a deterministic way based on account number while preventing negative balances (unless overdraft is enabled).

## Tests

Run all tests:

```bash
go test ./...
```

Test coverage areas:

- authentication and middleware behavior
- DTO parsing paths (CSV sample ingestion)
- repository behavior (PostgreSQL)
- validation and rating service behavior
