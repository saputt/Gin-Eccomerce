# Gin E-Commerce API

A learning project: an e-commerce backend written in Go with Gin, GORM, and PostgreSQL. The goal is to go beyond "it works" and practice real layered architecture — SOLID boundaries, dependency injection via interfaces, request/response DTOs, transactions, and versioned SQL migrations.

> **Status: work in progress.** Not all modules are wired up yet. The project currently does not compile end-to-end (see [Known issues](#known-issues)).

## Tech stack

| Concern | Choice |
| --- | --- |
| Language | Go |
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) (PostgreSQL driver) |
| Database | PostgreSQL (Neon) |
| Auth | JWT (`golang-jwt/jwt/v5`) + bcrypt |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) (raw SQL) |
| Payments | Xendit (provider adapter) |

## Architecture

Domain-oriented layout, one package per business domain:

```
cmd/api/                      # entrypoint: config, DI wiring, router
internal/
  common/
    helper/                   # transaction context helpers
    middleware/               # auth, role guard, rate limiter, semaphore
    response/                 # uniform success/error envelopes
  domain/
    auth/                     # login, register, JWT issuing/validation
    user/                     # user model + repository
    product/                  # product CRUD
    inventory/                # stock per product
    order/                    # order + suborder per seller
    payment/                  # payment records + Xendit provider
  infrasctructure/
    database/                 # PostgreSQL connection (GORM)
migration/                    # paired up/down SQL migrations
```

Each domain follows the same shape:

- `model.go` — persistence entities
- `dto.go` — request/response payloads, separated from models
- `handler.go` — HTTP binding, status codes, route registration
- `service.go` — business logic
- `interface.go` / `repository.go` — contracts the service depends on
- `repository/gorm.go` — GORM implementation of the contract

Services depend on interfaces rather than concrete repositories, and cross-domain collaborators are declared as narrow contracts (e.g. `ProductServiceContract`) so each package only sees what it needs. Transactions are passed through `context.Context` via `helper.InjectTX` / `helper.ExtractTX`, letting repositories participate in a transaction started by the service without leaking `*gorm.DB` into signatures.

## Endpoints

Base path: `/api`

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| POST | `/auth/register` | — | Register a user |
| POST | `/auth/login` | — | Login, returns JWT |
| GET | `/products` | — | List products |
| GET | `/products/:productId` | — | Product detail |
| POST | `/products` | Bearer (SELLER) | Create product + inventory |
| PATCH | `/products/:productId` | Bearer (SELLER) | Update product |
| DELETE | `/products/:productId` | Bearer (SELLER) | Delete product |
| PATCH | `/inventory/:productId/stock` | — | Update stock |
| GET | `/inventory/:productId` | — | Stock by product |
| POST | `/orders` | Bearer | Create order (split into suborders per seller) |

Auth uses `Authorization: Bearer <token>`. JWT claims carry `id`, `email`, and `role`.

## Getting started

Prerequisites: Go 1.26+, PostgreSQL, and the `migrate` CLI.

```bash
cp .env.example .env
# JWT_SECRET=...
# DATABASE_URL=postgresql://user:pass@host/db?sslmode=require
```

Run migrations and start the server:

```bash
make migrate-up
go run ./cmd/api
```

Server listens on `:3000`.

### Make targets

| Target | Purpose |
| --- | --- |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back one migration |
| `make migrate-new name=add_x` | Create a new up/down pair |
| `make migrate-status` | Show current migration version |
| `make migrate-force version=N` | Force version (recover a dirty state) |
| `make migrate-fix` | Force to the last applied version |

## Known issues

- `cmd/api/main.go` still imports the flat package layout; those imports were repathed to `internal/domain/...` but the wiring itself is incomplete.
- `OrderHandler.CancelOrder` calls `CancelOrder()` without arguments, and `ProductHandler.RegisterRoutes` is called without the auth middleware — the build fails on these.
- The order and payment modules are not registered in `main.go` yet.

## Roadmap

- Finish order flow: create → pay → cancel, with stock rollback
- Register order/payment handlers and providers
- Stock reservation and race-safe updates
- Tests around services and repositories
- Request validation and consistent error mapping

## Note on secrets

`.env` is git-ignored. Copy `.env.example` and fill in your own values. The file was committed in the first commit of this repository, so the old database URL and JWT secret still exist in history — rotate them if it was ever a live credential.
