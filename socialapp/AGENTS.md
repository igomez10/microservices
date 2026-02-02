# Repository Guidelines

## Project Structure & Module Organization
- `cmd/` contains service entrypoints (main packages) for binaries.
- `internal/` holds application logic intended for internal use only.
- `pkg/` provides shared, reusable packages (e.g., scopes and helpers).
- `db/` includes SQL schema/config (see `db/sqlc.yaml`) and generated queries.
- `e2e/` and `integration_tests/` contain end-to-end and integration tests.
- `socialappapi/`, `client/`, `Models/`, `generated-schema/`, `cli/`, `jmeter/`, and `postman/` are generated artifacts from `openapi.yaml`.
- Infrastructure and observability live in `docker-compose.yml`, `kubernetes/`, `grafana/`, `prometheus/`, `loki/`, `tempo/`, `otel-collector/`, and `traefik/`.

## Build, Test, and Development Commands
- `make start-dev-server`: run the Go service locally with live reload.
- `make build`: build all Go packages.
- `make test`: run unit tests (`go test ./...`).
- `make test-e2e`: run end-to-end tests (`go test -v ./e2e/...`).
- `make start`: bring up the docker-compose stack.
- `make sqlc-generate`: regenerate SQLC code from `db/sqlc.yaml`.
- `make generate-openapi`: regenerate OpenAPI server, client, docs, and schemas.

## Coding Style & Naming Conventions
- Go code is formatted with `gofmt` and imports are organized with `goimports`.
- Use standard Go naming: `camelCase` for vars, `PascalCase` for exported types.
- Test files follow `*_test.go` and table-driven tests where appropriate.
- Generated code should not be edited manually; update sources like `openapi.yaml` or SQL and regenerate.

## Testing Guidelines
- Unit tests run via `go test ./...`.
- End-to-end tests live in `e2e/` and run with `make test-e2e`.
- Name tests descriptively (e.g., `TestCreateUser`, `TestURLLifecycle`).

## Commit & Pull Request Guidelines
- Recent commits use conventional prefixes like `feat:`, `fix:`, `docs:`, and `chore:`; follow this style when possible.
- Keep commits focused and include context in the message (what/why).
- PRs should include a clear description, linked issues (if any), and test evidence; include screenshots for UI changes under `client/` or `streamlit/`.

## Configuration & Deployment Notes
- Local orchestration is via `docker-compose.yml` and supporting infra folders.
- Deployment scripts live in `deploy.sh`; ensure required env vars are set before running.