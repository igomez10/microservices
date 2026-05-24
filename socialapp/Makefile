# Socialapp Makefile
#
# Common targets for building, testing, generating code, and operating the
# social-networking API service. Run `make help` to list every target with
# its description.

.DEFAULT_GOAL := help

# ----------------------------------------------------------------------------
# Help
# ----------------------------------------------------------------------------

.PHONY: help
help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } \
		/^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2 }' \
		$(MAKEFILE_LIST)

##@ Setup

.PHONY: setup-init
setup-init: ## Create the Docker networks shared with docker-elk
	docker network create docker-elk_elk
	docker network create socialapp_default

##@ Code generation

.PHONY: generate-openapi
generate-openapi: ## Generate server stubs, clients, docs, CLI, JMeter and Postman from openapi.yaml
	set -e
	# openapi-generator-cli version-manager set 7.11.0

	# Server stubs
	rm -rf ./socialappapi
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name go-server \
		--output socialappapi \
		--ignore-file-override \
		--git-repo-id=microservices/socialapp \
		--git-user-id=igomez10 \
		--additional-properties=packageName=openapi \
		--additional-properties=outputAsLibrary=true \
		--additional-properties=packageVersion=1.18.0 \
		--additional-properties=serverPort=8080 \
		--additional-properties=featureCORS=true \
		--additional-properties=hideGenerationTimestamp=true \
		--additional-properties=addResponseHeaders=true \
		--additional-properties=onlyInterfaces=true \
		--additional-properties=router=chi
	rm -rf socialappapi/openapi
	mv socialappapi/go socialappapi/openapi

	# Go client SDK
	rm -rf client
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name go \
		--output client \
		--git-repo-id=microservices/socialapp \
		--git-user-id=igomez10 \
		--additional-properties=packageName=client \
		--additional-properties=isGoSubmodule=true \
		--additional-properties=packageName=client \
		--additional-properties=structPrefix=true
	rm client/go.mod
	rm client/go.sum
	rm client/test/*

	# Markdown docs
	rm -rf Models
	openapi-generator generate -i openapi.yaml \
		--output Models \
		--generator-name markdown

	# MySQL schema
	rm -rf generated-schema
	openapi-generator generate -i openapi.yaml \
		--generator-name mysql-schema \
		--output generated-schema \
		--additional-properties=namedParametersEnabled=true \
		--additional-properties=identifierNamingConvention=snake_case

	# Format generated Go code
	go fmt ./...
	goimports -w .

	# Bash CLI
	rm -rf cli
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name bash \
		--output cli \
		--additional-properties=generateZshCompletion=true \
		--additional-properties=scriptName=socialapp-cli

	# JMeter test plan
	rm -rf jmeter
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name jmeter \
		--output jmeter

	# Postman collection
	rm -rf postman
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name postman-collection \
		--output postman

	# Frontend TypeScript types
	$(MAKE) generate-frontend-types

.PHONY: generate-frontend-types
generate-frontend-types: ## Generate frontend TypeScript types from openapi.yaml
	cd frontend && npm run generate:api

.PHONY: sqlc-generate
sqlc-generate: ## Generate type-safe Go DB code from db/query.sql
	sqlc generate --file db/sqlc.yaml

.PHONY: generate-scopes
generate-scopes: ## Generate OAuth2 scope constants from openapi.yaml
	cd tools/scopegen && go run main.go -input ../../openapi.yaml -output ../../pkg/scopes/scopes.go
	go fmt ./pkg/scopes/...
	goimports -w ./pkg/scopes/

##@ Build & dev server

.PHONY: build
build: ## Build the socialapp binaries
	go build ./...

# Tip: `export $$(cat .env | xargs) && make start-dev-server` to load env vars first.
.PHONY: start-dev-server
start-dev-server: ## Hot-reload dev server via reflex
	reflex -r .go -s --decoration="fancy" -- go run ./...

# Loads .env_vscode_debug (not .env) so host-side URL shortener and IDE debug
# settings stay in one place.
.PHONY: start-dev-server-debug
start-dev-server-debug: ## Hot-reload dev server under Delve (listens on :2345)
	mkdir -p logs
	reflex -r .go -s --decoration="fancy" -- sh -c '\
		if [ -f ./.env_vscode_debug ]; then set -a; . ./.env_vscode_debug; set +a; fi; \
		PORT=$${PORT:-8085} \
		LOG_LEVEL=$${LOG_LEVEL:-debug} \
		DATABASE_URL=$${DATABASE_URL:-postgres://postgres:password@localhost:5432/socialapp?sslmode=disable} \
		REDIS_URL=$${REDIS_URL:-redis://localhost:6379} \
		OTEL_EXPORTER_OTLP_ENDPOINT=$${OTEL_EXPORTER_OTLP_ENDPOINT:-http://localhost:4317} \
		LOGSTASH_HOST=$${LOGSTASH_HOST:-tcp://localhost:53000} \
		PROPERTIES_SUBDOMAIN=$${PROPERTIES_SUBDOMAIN:-http://properties.localhost} \
		PUTTYKNIFE_URL=$${PUTTYKNIFE_URL:-http://localhost:8096} \
		PUTTYKNIFE_DOMAIN=$${PUTTYKNIFE_DOMAIN:-puttyknife-server.localhost} \
		URL_SHORTENER=$${URL_SHORTENER:-http://localhost:8089} \
		URLSHORTENER_URL=$${URLSHORTENER_URL:-http://localhost:8089} \
		URLSHORTENER_SUBDOMAIN=$${URLSHORTENER_SUBDOMAIN:-http://localhost:8089} \
		SOCIALAPP_SUBDOMAIN=$${SOCIALAPP_SUBDOMAIN:-http://localhost:$${PORT:-8085}} \
		JWT_SECRET=$${JWT_SECRET:-secret} \
		KIBANA_URL=$${KIBANA_URL:-http://localhost:5601} \
		KIBANA_SUBDOMAIN=$${KIBANA_SUBDOMAIN:-kibana.localhost} \
		LOCAL_SUBDOMAIN=$${LOCAL_SUBDOMAIN:-localhost:$${PORT:-8085}} \
		dlv debug ./cmd --headless --listen=:2345 --api-version=2 --accept-multiclient --continue 2>&1 | tee logs/socialapp-debug.log'

.PHONY: frontend-debug
frontend-debug: ## Start the Vite frontend with the dev proxy pointed at the local API
	mkdir -p logs
	cd frontend && VITE_DEV_PROXY_TARGET=http://localhost:8085 npm run dev -- --host 127.0.0.1 --debug 2>&1 | tee ../logs/frontend-debug.log

##@ Docker Compose

.PHONY: start
start: ## docker compose up -d (this service only)
	docker compose up -d

.PHONY: compose-up
compose-up: ## Recreate this service with docker compose (down + up, foreground)
	docker compose down && docker compose up

.PHONY: compose-all-up
compose-all-up: ## Bring up the full local stack (docker-elk, socialapp, puttyknife, urlshortener)
	cd ../../docker-elk && docker compose up setup -d
	sleep 10
	cd ../../docker-elk && docker compose up -d
	docker compose up -d
	cd ../../puttyknife && docker compose up -d
	cd ../urlshortener && docker compose up -d

.PHONY: compose-all-down
compose-all-down: ## Tear down the full local stack
	cd ../urlshortener && docker compose down
	cd ../../puttyknife && docker compose down
	docker compose down
	cd ../../docker-elk && docker compose down

.PHONY: compose-prune
compose-prune: ## Tear down the full local stack and remove volumes/orphans
	docker compose down --volumes --remove-orphans
	cd ../urlshortener && docker compose down --volumes --remove-orphans
	cd ../../docker-elk && docker compose down --volumes --remove-orphans

##@ Testing

.PHONY: test
test: test-unit test-integration test-e2e test-contract ## Run unit, integration, e2e and contract tests

.PHONY: test-ci
test-ci: test-unit ## Test target run in CI (unit only)

.PHONY: test-unit
test-unit: ## Run unit tests (excludes e2e, integration, contracts)
	go test $$(go list ./... | grep -v -E '/(e2e|integration_tests|contracts|pkg/controller/url)$$')

# Defaults APP_HOST to start-dev-server-debug's PORT (8085).
.PHONY: test-integration
test-integration: ## Run integration tests against $$APP_HOST (default localhost:8085)
	APP_HOST=$${APP_HOST:-http://localhost:8085} go test -v ./integration_tests/... -timeout 300s

.PHONY: test-integration-local
test-integration-local: ## Run integration tests against http://localhost:8085
	APP_HOST=http://localhost:8085 go test -v ./integration_tests/... -timeout 300s

.PHONY: test-e2e
test-e2e: ## Run Go end-to-end tests
	go test -v ./e2e/... -timeout 300s

.PHONY: test-contract
test-contract: ## Run Pact consumer contract tests
	../scripts/pact-check.sh
	go test ./contracts -run Test.*Contract -count=1
	go test ./pkg/controller/url -run TestURLApiService_GetUrlDataContract -count=1

.PHONY: publish-contracts
publish-contracts: ## Publish Pact contracts to the broker
	../scripts/publish-contracts.sh \
		contracts/pacts/socialapp-client-socialapp.json \
		contracts/pacts/socialapp-urlshortener.json

# Playwright browser E2E (see frontend/playwright.config.ts). Requires
# E2E_USERNAME and E2E_PASSWORD (export them or put them in ./.env). Optional:
# E2E_API_BASE_URL, E2E_BASE_URL. Vite is started automatically. First-time
# setup: make playwright-install
.PHONY: frontend-test-e2e playwright playwright-install
frontend-test-e2e: ## Run Playwright browser E2E suite (needs E2E_* env)
	if [ -f ./.env ]; then set -a; . ./.env; set +a; fi; cd frontend && npm run test:e2e

playwright: frontend-test-e2e ## Alias for frontend-test-e2e

playwright-install: ## Install Playwright browsers (one-time setup)
	cd frontend && npx playwright install --with-deps chromium

##@ Auxiliary tools

.PHONY: start-cady
start-cady: ## Reverse-proxy localhost:4545 -> localhost:8080 via Caddy
	caddy reverse-proxy --from localhost:4545 --to localhost:8080

.PHONY: start-jmeter
start-jmeter: ## Run the JMeter UserApi plan and tail its log
	jmeter -H localhost -P 9091 -t jmeter/UserApi.jmx & tail -f jmeter.log

.PHONY: build-rust-client
build-rust-client: ## Build the generated Rust client
	cargo build --manifest-path generated-rust-client/Cargo.toml

.PHONY: run-rust-client
run-rust-client: ## Run the generated Rust client
	cargo run --manifest-path generated-rust-client/Cargo.toml

##@ Infrastructure (Terraform)

.PHONY: tf-init
tf-init: ## terraform init in ./infrastructure
	cd infrastructure && terraform init

.PHONY: tf-plan
tf-plan: ## terraform plan in ./infrastructure
	cd infrastructure && terraform plan

.PHONY: tf-apply
tf-apply: ## terraform apply -auto-approve in ./infrastructure
	cd infrastructure && terraform apply -auto-approve

##@ Deployment

.PHONY: build-push-integration-tests
build-push-integration-tests: ## Build and push the integration-test Docker image
	docker build . -f Integration.Dockerfile --platform=linux/arm64 -t igomeza/socialapptests
	docker push igomeza/socialapptests

# Requires SSH_PORT, SSH_USERNAME and SSH_HOST to be set in the environment.
.PHONY: deploy
deploy: ## SSH into the host and run ./deploy.sh remotely
	ssh -o StrictHostKeyChecking=no \
		-p $(SSH_PORT) \
		$(SSH_USERNAME)@$(SSH_HOST) \
		'bash -s' < ./deploy.sh
