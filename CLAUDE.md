# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

A collection of Go microservices demonstrating cloud-native development practices. Contains three main services:
- **socialapp**: Social networking API (users, comments, followers, RBAC, OAuth2)
- **urlshortener**: URL shortening service with analytics
- **bitcoinprice**: Bitcoin price tracker using Buffalo framework

## Common Commands

### Socialapp (primary service)
```bash
cd socialapp
make build                    # Build the service
make test                     # Run unit tests
make test-e2e                 # Run end-to-end tests
make generate-openapi         # Regenerate API code from openapi.yaml
make sqlc-generate            # Regenerate SQL code from db/query.sql
make generate-scopes          # Generate OAuth2 scopes from openapi.yaml
make start-dev-server         # Hot reload dev server (requires reflex)
docker compose up -d          # Start full stack with dependencies
```

### URL Shortener
```bash
cd urlshortener
make build                    # Build the service
make test                     # Run unit tests
make generate-openapi         # Regenerate API code
make sqlc-generate            # Regenerate SQL code
make start-dev-server         # Hot reload dev server
docker compose up -d          # Start with dependencies
```

### Bitcoin Price (Buffalo)
```bash
cd bitcoinprice
buffalo dev                   # Start development server
buffalo test                  # Run tests
```

### Running a Single Test
```bash
go test -v -run TestName ./path/to/package
```

## Architecture

### Code Generation Pipeline
Both socialapp and urlshortener follow an API-first approach:
1. **OpenAPI spec** (`openapi.yaml`) defines the API contract
2. **openapi-generator** creates server stubs, Go clients, CLI tools, docs, and Postman collections
3. **SQLC** generates type-safe database code from SQL queries

Generated code locations:
- socialapp: `socialappapi/openapi/` (server), `client/` (Go client), `Models/` (docs)
- urlshortener: `generated/server/`, `generated/clients/go/client/`

### Service Structure (socialapp)
```
cmd/main.go              # Entry point, config, telemetry setup
internal/
  server/                # HTTP server and router assembly
  routers/               # Chi router with middleware chain
  middlewares/           # Authorization, caching, request ID, etc.
pkg/
  controller/            # Business logic (authentication, comment, role, scope, url, user)
  dbpgx/                 # SQLC-generated database queries
  scopes/                # OAuth2 scope definitions (auto-generated)
```

### Database
- PostgreSQL with pgx driver (socialapp uses connection pooling via `ForcedConnectionPool`)
- Schema in `db/setup/schema.sql`, queries in `db/query.sql`
- SQLC config in `db/sqlc.yaml`

### Authentication & Authorization
- OAuth2 with JWT tokens
- Scopes defined in openapi.yaml and auto-generated to `pkg/scopes/`
- Authorization middleware validates required scopes per endpoint

### Observability Stack
- **Tracing**: OpenTelemetry -> OTLP exporter -> Tempo/Jaeger
- **Metrics**: Prometheus (port 9090), Grafana Mimir
- **Logging**: zerolog -> Grafana Loki (Docker driver)
- **Profiling**: Pyroscope, pprof at `/debug/pprof/`

### Key Dependencies
- Router: chi/v5
- Database: pgx/v5, SQLC
- Logging: zerolog
- Config: jessevdk/go-flags (CLI args and env vars)
- HTTP client: hashicorp/go-retryablehttp (with retry/backoff)
- Telemetry: OpenTelemetry SDK

## CI/CD

- **GitHub Actions**: Deploys socialapp and urlshortener on push to `mainline` when respective directories change
- **CircleCI**: Builds bitcoinprice and loadgenerator

## Docker Logging Standard

Services use Grafana Loki driver with these options:
```yaml
logging:
  driver: loki
  options:
    loki-url: "http://localhost:3100/loki/api/v1/push"
    loki-batch-size: "5000"
    loki-retries: "3"
    max-size: "10m"
    max-file: "3"
```
