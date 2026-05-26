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

generate-openapi:
	set -e
	# openapi-generator-cli version-manager set 7.11.0
	#  generate server:
	rm -rf ./socialappapi
	openapi-generator generate --input-spec openapi.yaml  \
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

	# generate go client
	rm -rf client
	openapi-generator generate --input-spec openapi.yaml  \
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
	rm client/test/* # remove tests in client sdk

	# generate documentation
	rm -rf Models
	openapi-generator generate -i openapi.yaml  \
	--output Models \
	--generator-name markdown

	# generate schema
	rm -rf generated-schema
	openapi-generator generate -i openapi.yaml  \
	--generator-name mysql-schema \
	--output generated-schema \
	--additional-properties=namedParametersEnabled=true \
	--additional-properties=identifierNamingConvention=snake_case

	# format go code
	go fmt ./...
	goimports -w .

	# generate cli
	rm -rf cli
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name bash \
        --output cli \
		--additional-properties=generateZshCompletion=true \
		--additional-properties=scriptName=socialapp-cli

	# generate jmeter
	rm -rf jmeter
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name jmeter \
		--output jmeter

	# generate postman collection
	rm -rf postman
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name postman-collection \
		--output postman
	# generate frontend types
	$(MAKE) generate-frontend-types

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
