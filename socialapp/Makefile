setup-init:
	docker network create docker-elk_elk
	docker network create socialapp_default

compose-all-up:
	cd ../../docker-elk && docker compose up setup -d
	sleep 10
	cd ../../docker-elk && docker compose up -d
	docker compose up -d
	cd ../../puttyknife && docker compose up -d
	cd ../urlshortener && docker compose up -d

compose-all-down:
	cd ../urlshortener && docker compose down
	cd ../../puttyknife && docker compose down
	docker compose down
	cd ../../docker-elk && docker compose down

compose-prune:
	docker compose down --volumes  --remove-orphans
	cd ../urlshortener && docker compose down --volumes --remove-orphans
	cd ../../docker-elk && docker compose down --volumes  --remove-orphans

# ---------------------------------------------------------------------------
# OpenAPI code generation
#
# Every artifact below is generated from openapi.yaml and committed to the repo.
# Run `make generate-openapi` after editing the spec, or a single sub-target
# (e.g. `make generate-openapi-server`) while iterating.
# `make check-openapi-drift` verifies the committed output still matches the spec.
#
# Output differs between generator versions. The version is pinned in
# openapitools.json and fetched by the npm wrapper, so CI needs only Node.
# Override OPENAPI_GENERATOR to use a locally installed binary instead.
# ---------------------------------------------------------------------------
OPENAPI_GENERATOR         ?= npx @openapitools/openapi-generator-cli
OPENAPI_GENERATOR_VERSION ?= 7.25.0
OPENAPI_SPEC              := openapi.yaml
OPENAPI_GIT_USER_ID       := igomez10
OPENAPI_GIT_REPO_ID       := microservices/socialapp
OPENAPI_OUTPUTS           := socialappapi client Models postman jmeter \
                             frontend/src/api/types.ts

.PHONY: generate-openapi generate-openapi-server generate-openapi-client \
	generate-openapi-docs generate-openapi-jmeter generate-openapi-postman \
	generate-frontend-types validate-openapi check-openapi-generator \
	check-openapi-drift

check-openapi-generator:
	@v="$$($(OPENAPI_GENERATOR) version 2>/dev/null | tail -1)"; \
	if [ "$$v" != "$(OPENAPI_GENERATOR_VERSION)" ]; then \
		echo "openapi-generator $(OPENAPI_GENERATOR_VERSION) required, found $$v"; \
		echo "Regenerating with another version produces unrelated diffs."; \
		echo "Install the pinned version, or bump OPENAPI_GENERATOR_VERSION and regenerate everything."; \
		exit 1; \
	fi

validate-openapi: check-openapi-generator
	$(OPENAPI_GENERATOR) validate -i $(OPENAPI_SPEC)

generate-openapi: generate-openapi-server generate-openapi-client \
	generate-openapi-docs generate-openapi-jmeter generate-openapi-postman \
	generate-frontend-types

# Server: interfaces and chi controllers only (onlyInterfaces=true).
# The implementations are hand-written in pkg/controller.
generate-openapi-server: validate-openapi
	rm -rf socialappapi
	$(OPENAPI_GENERATOR) generate --input-spec $(OPENAPI_SPEC) \
		--generator-name go-server \
		--output socialappapi \
		--git-repo-id=$(OPENAPI_GIT_REPO_ID) \
		--git-user-id=$(OPENAPI_GIT_USER_ID) \
		--additional-properties=packageName=openapi \
		--additional-properties=outputAsLibrary=true \
		--additional-properties=serverPort=8080 \
		--additional-properties=featureCORS=true \
		--additional-properties=hideGenerationTimestamp=true \
		--additional-properties=addResponseHeaders=true \
		--additional-properties=onlyInterfaces=true \
		--additional-properties=router=chi
	# outputAsLibrary emits the package into socialappapi/go; the import path
	# used across the module is socialappapi/openapi.
	mv socialappapi/go socialappapi/openapi
	gofmt -w socialappapi
	goimports -w socialappapi

generate-openapi-client: validate-openapi
	rm -rf client
	$(OPENAPI_GENERATOR) generate --input-spec $(OPENAPI_SPEC) \
		--generator-name go \
		--output client \
		--git-repo-id=$(OPENAPI_GIT_REPO_ID) \
		--git-user-id=$(OPENAPI_GIT_USER_ID) \
		--additional-properties=packageName=client \
		--additional-properties=isGoSubmodule=true \
		--additional-properties=structPrefix=true \
		--global-property=apiTests=false,modelTests=false
	# the client ships as part of this module rather than its own; the generator
	# emits go.mod/go.sum regardless of isGoSubmodule
	rm -f client/go.mod client/go.sum
	gofmt -w client
	goimports -w client

generate-openapi-docs: validate-openapi
	rm -rf Models
	$(OPENAPI_GENERATOR) generate --input-spec $(OPENAPI_SPEC) \
		--generator-name markdown \
		--output Models

generate-openapi-jmeter: validate-openapi
	rm -rf jmeter
	$(OPENAPI_GENERATOR) generate --input-spec $(OPENAPI_SPEC) \
		--generator-name jmeter \
		--output jmeter

generate-openapi-postman: validate-openapi
	rm -rf postman
	$(OPENAPI_GENERATOR) generate --input-spec $(OPENAPI_SPEC) \
		--generator-name postman-collection \
		--output postman

# Fails if the committed generated code no longer matches openapi.yaml.
# Intended for CI.
check-openapi-drift: generate-openapi
	git diff --exit-code -- $(OPENAPI_OUTPUTS)

generate-frontend-types:
	cd frontend && npm run generate:api

# Playwright browser E2E (see frontend/playwright.config.ts). Requires E2E_USERNAME and E2E_PASSWORD
# (export them or put them in ./.env). Optional: E2E_API_BASE_URL, E2E_BASE_URL. Vite is started automatically.
# First-time setup: make playwright-install
.PHONY: playwright frontend-test-e2e playwright-install
frontend-test-e2e playwright:
	if [ -f ./.env ]; then set -a; . ./.env; set +a; fi; cd frontend && npm run test:e2e

playwright-install:
	cd frontend && npx playwright install --with-deps chromium

frontend-debug:
	mkdir -p logs
	cd frontend && VITE_DEV_PROXY_TARGET=http://localhost:8085 npm run dev -- --host 127.0.0.1 --debug 2>&1 | tee ../logs/frontend-debug.log


sqlc-generate:
	sqlc generate --file db/sqlc.yaml

generate-scopes:
	cd tools/scopegen && go run main.go -input ../../openapi.yaml -output ../../pkg/scopes/scopes.go
	go fmt ./pkg/scopes/...
	goimports -w ./pkg/scopes/

compose-up:
	docker compose down && docker-compose up

# export $(cat .env | xargs) &&  make start-dev-server
start-dev-server:
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
	go run ./cmd

# Loads .env_vscode_debug (not .env) so host-side URL shortener and IDE debug settings stay in one place.
start-dev-server-debug:
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

build:
	go build ./...

# Run only unit tests (excludes e2e, integration, and contract tests)
test-unit:
	go test $$(go list ./... | grep -v -E '/(e2e|integration_tests|contracts|pkg/controller/url)$$')

# Run integration tests (requires a running API; default APP_HOST matches start-dev-server-debug PORT)
test-integration:
	APP_HOST=$${APP_HOST:-http://localhost:8085} go test -v ./integration_tests/... -timeout 300s

# Run integration tests against local server
test-integration-local:
	APP_HOST=http://localhost:8085 go test -v ./integration_tests/... -timeout 300s

# Run end-to-end tests
test-e2e:
	go test -v ./e2e/... -timeout 300s

# Run contract tests
test-contract:
	../scripts/pact-check.sh
	go test ./contracts -run Test.*Contract -count=1
	go test ./pkg/controller/url -run TestURLApiService_GetUrlDataContract -count=1

# Run all tests
test: test-unit test-integration test-e2e test-contract

# Add CI specific test target
test-ci: test-unit

publish-contracts:
	../scripts/publish-contracts.sh \
		contracts/pacts/socialapp-client-socialapp.json \
		contracts/pacts/socialapp-urlshortener.json

start:
	docker compose up -d

start-cady:
	caddy reverse-proxy --from localhost:4545  --to localhost:8080

start-jmeter:
	jmeter -H localhost -P 9091 -t jmeter/UserApi.jmx & tail -f jmeter.log
build-rust-client:
	cargo build --manifest-path generated-rust-client/Cargo.toml
run-rust-client:
	cargo run --manifest-path generated-rust-client/Cargo.toml

tf-init:
	cd infrastructure && terraform init

tf-plan:
	cd infrastructure && terraform plan

tf-apply:
	cd infrastructure && terraform apply -auto-approve

build-push-integration-tests:
	docker build . -f Integration.Dockerfile --platform=linux/arm64 -t igomeza/socialapptests
	docker push igomeza/socialapptests

deploy:
	# ssh into the server and run the deploy script
	 ssh -o StrictHostKeyChecking=no \
            -p $(SSH_PORT) \
            $(SSH_USERNAME)@$(SSH_HOST) \
			'bash -s' < ./deploy.sh
