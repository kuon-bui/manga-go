# manga-go

A RESTful API backend for a manga reading platform, built with Go. The service handles user authentication, manga catalog management, chapter delivery, reading history tracking, and role-based access control.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Clone the Repository](#1-clone-the-repository)
  - [Configure Environment](#2-configure-environment)
  - [Start Infrastructure Services](#3-start-infrastructure-services)
  - [Run Migrations](#4-run-migrations)
  - [Start the API Server](#5-start-the-api-server)
- [Project Structure](#project-structure)
- [API Resources](#api-resources)
- [API Documentation](#api-documentation)
- [Development](#development)
  - [Live Reload](#live-reload)
  - [Async Worker](#async-worker)
  - [Generate Swagger Docs](#generate-swagger-docs)
  - [Makefile Commands](#makefile-commands)
- [Infrastructure Services](#infrastructure-services)
- [Architecture Overview](#architecture-overview)

---

## Features

- 🔐 **Authentication** – JWT-based auth with HTTP-only cookies, refresh token rotation, and Redis blacklist
- 🛡️ **Authorization** – Role-Based Access Control (RBAC) via Casbin
- 📚 **Manga Catalog** – Full CRUD for Comics, Chapters, Pages, Authors, Genres, and Tags
- 🌍 **Translation Groups** – Manage translation team groups linked to chapters
- 📖 **Reading History** – Track per-user reading progress
- 📧 **Email** – Transactional emails (password reset, etc.) via go-mail + Mailpit for local testing
- 🗄️ **File Storage** – S3-compatible object storage (MinIO)
- 📈 **Observability** – Distributed tracing (OpenTelemetry + Jaeger) and Prometheus metrics
- 📝 **API Docs** – Auto-generated Swagger/OpenAPI documentation

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22+ |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| Dependency Injection | [Uber Fx](https://github.com/uber-go/fx) |
| ORM | [GORM](https://gorm.io) + PostgreSQL driver |
| Database | PostgreSQL 17 |
| Cache / Queue Backend | Redis 8 |
| Task Queue | [Asynq](https://github.com/hibiken/asynq) |
| Migrations | [sql-migrate](https://github.com/rubenv/sql-migrate) |
| Authentication | [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) |
| Authorization | [Casbin v3](https://github.com/casbin/casbin) |
| Logger | [Uber Zap](https://github.com/uber-go/zap) |
| Config | [Viper](https://github.com/spf13/viper) |
| Tracing | [OpenTelemetry](https://opentelemetry.io) → Jaeger |
| Metrics | Prometheus |
| Storage | S3-compatible (MinIO) |
| Email | [go-mail](https://github.com/wneessen/go-mail) + Mailpit |
| API Docs | [Swaggo](https://github.com/swaggo/swag) |

---

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/)
- [`sql-migrate`](https://github.com/rubenv/sql-migrate) CLI – for running migrations manually

```bash
go install github.com/rubenv/sql-migrate/...@latest
```

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/kuon-bui/manga-go.git
cd manga-go
```

### 2. Configure Environment

Copy the example environment file and update values as needed:

```bash
cp .env.example .env
cp config.yml.example config.yml
```

Key environment variables in `.env`:

```env
SERVICE_PORT=8080
DB_USER=admin
DB_PASSWORD=postgres
DB_NAME=mydb
DB_PORT=5432
REDIS_PORT=6379
REDIS_PASSWORD=
JAEGER_PORT=4317
```

### 3. Start Infrastructure Services

Spin up PostgreSQL, Redis, MinIO, Mailpit, and Jaeger using Docker Compose:

```bash
docker compose up -d postgresql redis minio jaeger mailpit
```

Or start the full stack (including the API server and queue worker):

```bash
docker compose up -d
```

### 4. Run Migrations

```bash
sql-migrate up -config=configs/dbconfig.yml
```

Check migration status:

```bash
sql-migrate status -config=configs/dbconfig.yml
```

### 5. Start the API Server

```bash
go run cmd/dev/main.go
```

The API will be available at `http://localhost:8080`.

---

## Project Structure

```
manga-go/
├── cmd/
│   ├── dev/main.go          # API server entrypoint
│   └── queue/main.go        # Async worker entrypoint
├── configs/
│   └── dbconfig.yml         # sql-migrate configuration
├── docs/                    # Auto-generated Swagger docs
├── internal/
│   ├── app/
│   │   ├── fx.go            # Root Fx module
│   │   ├── api/
│   │   │   ├── common/      # Shared route/response utilities
│   │   │   ├── route/       # Route handlers (one package per resource)
│   │   │   └── server/      # Gin setup and HTTP server lifecycle
│   │   ├── asynq/           # Async task runner setup
│   │   └── middleware/      # Auth and slug middleware
│   └── pkg/
│       ├── common/          # Shared types (SqlModel, Paging)
│       ├── config/          # Config loader (Viper)
│       ├── model/           # GORM entity models
│       ├── repo/            # Data repositories
│       ├── request/         # Request DTOs per resource
│       ├── services/        # Business logic per resource
│       ├── casbin/          # RBAC engine setup
│       ├── gorm/            # Database connection
│       ├── redis/           # Redis client
│       ├── jwt_provider/    # JWT generation & validation
│       ├── logger/          # Zap logger wrapper
│       ├── mail/            # Email service
│       ├── object_storage/  # S3/MinIO client
│       ├── tracer/          # OpenTelemetry setup
│       └── queue/           # Asynq task definitions
├── migrations/              # SQL migration files
├── resources/               # Static assets (e.g., Casbin policy)
├── docker-compose.yml
├── Dockerfile
├── makefile
├── api.air.toml             # Air live-reload config for API
├── queue.air.toml           # Air live-reload config for worker
├── config.yml.example
└── .env.example
```

---

## API Resources

| Resource | Path | Description |
|---|---|---|
| **User** | `/users` | Sign up, sign in, sign out, password reset, token refresh |
| **Author** | `/authors` | Manage manga authors |
| **Genre** | `/genres` | Manage manga genres |
| **Tag** | `/tags` | Manage content tags |
| **Comic** | `/comics` | Manga titles with metadata (authors, genres, tags, age rating) |
| **Chapter** | `/chapters` | Chapters linked to comics and translation groups |
| **Translation Group** | `/translation-groups` | Translation team management |
| **File** | `/files` | File upload endpoint |
| **Role** | `/roles` | User role management (RBAC) |
| **Permission** | `/permissions` | Permission definition (RBAC) |
| **Reading History** | `/reading-histories` | Track user reading progress per comic/chapter |

---

## API Documentation

Swagger UI is served at:

```
http://localhost:8080/swagger/index.html
```

To regenerate the docs after changing handler annotations:

```bash
make swagger
```

---

## Development

### Live Reload

Install [Air](https://github.com/cosmtrek/air) for live reload:

```bash
go install github.com/air-verse/air@latest
```

Start with live reload:

```bash
# API server
air -c api.air.toml

# Queue worker
air -c queue.air.toml
```

### Async Worker

```bash
go run cmd/queue/main.go
```

Monitor tasks via Asynqmon UI at `http://localhost:8091`.

### Generate Swagger Docs

```bash
make swagger
```

This runs `swag init` and writes generated files to `docs/`.

### Makefile Commands

| Command | Description |
|---|---|
| `make start` | Start the API server |
| `make migrate` | Apply pending migrations |
| `make migrate-rollback` | Roll back the last migration |
| `make migrate-rollback-all` | Roll back all migrations |
| `make migration` | Create a new migration file |
| `make swagger` | Regenerate Swagger documentation |

---

## Infrastructure Services

All services are defined in `docker-compose.yml`:

| Service | URL / Port | Description |
|---|---|---|
| **API Server** | `http://localhost:8080` | Main REST API |
| **PostgreSQL** | `localhost:5432` | Primary database |
| **Redis** | `localhost:6379` | Cache & task queue backend |
| **MinIO** | `http://localhost:9001` | S3-compatible object storage (UI) |
| **Mailpit** | `http://localhost:8025` | Local email testing UI |
| **Jaeger** | `http://localhost:16686` | Distributed tracing UI |
| **PgAdmin** | `http://localhost:8081` | PostgreSQL admin UI |
| **Asynqmon** | `http://localhost:8091` | Async task monitoring UI |

---

## Architecture Overview

The project follows a layered architecture enforced through **Uber Fx** dependency injection:

```
HTTP Request
     │
     ▼
 Middleware (Auth, Slug)
     │
     ▼
 Route Handler (internal/app/api/route/<resource>/)
     │
     ▼
 Service Layer (internal/pkg/services/<resource>/)
     │
     ▼
 Repository Layer (internal/pkg/repo/<resource>/)
     │
     ▼
 GORM + PostgreSQL
```

**Key patterns:**
- Every package exposes a `var Module = fx.Module(...)` consumed by its parent module
- All models embed `common.SqlModel` (UUID primary key, `created_at`, `updated_at`, soft-delete `deleted_at`)
- All repositories embed `base.BaseRepository[T]` providing generic CRUD operations
- All service methods return `response.Result` for consistent HTTP responses
- Every handler includes Swagger annotations for auto-generated API documentation
