# URL Shortener Makefile
#
# Common targets for building, testing, generating code, and operating the
# URL shortener service. Run `make help` to list every target with its
# description.

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

##@ Code generation

.PHONY: generate-openapi
generate-openapi: ## Generate server stubs, clients, docs, CLI, JMeter and Postman from openapi.yaml
	set -e

	# Pin the openapi-generator version used below.
	openapi-generator-cli version-manager set 7.11.0

	# Server stubs
	rm -rf ./generated/server
	openapi-generator-cli generate -i openapi.yaml \
		-g go-server \
		-o ./generated/server \
		-p packageName=server \
		--git-user-id=igomez10 \
		--git-repo-id=microservices \
		--additional-properties=serverPort=8080 \
		--additional-properties=featureCORS=true \
		--additional-properties=hideGenerationTimestamp=true \
		--additional-properties=addResponseHeaders=true \
		--additional-properties=onlyInterfaces=true \
		--additional-properties=router=chi \
		--additional-properties=packageName=server \
		--additional-properties=withGoMod=false \
		--additional-properties=withGoCodegenComment=false \
		--additional-properties=sourceFolder=. \
		--additional-properties=outputAsLibrary=true

	# Go client SDK
	rm -rf generated/clients/go/client
	openapi-generator-cli generate --input-spec openapi.yaml \
		--git-user-id=igomez10 \
		--git-repo-id=microservices \
		--generator-name go \
		--output generated/clients/go/client \
		-p packageName=client \
		--additional-properties=packageName=client \
		--additional-properties=withGoMod=false \
		--additional-properties=withGoCodegenComment=false \
		--additional-properties=withGoTests=false
	rm generated/clients/go/client/test/*

	# Markdown docs
	openapi-generator-cli generate -i openapi.yaml \
		--output generated/docs \
		--generator-name markdown

	# MySQL schema
	rm -rf mysql-schema
	openapi-generator-cli generate -i openapi.yaml \
		--generator-name mysql-schema \
		--output generated/schema \
		--additional-properties=namedParametersEnabled=true \
		--additional-properties=identifierNamingConvention=snake_case

	# Format generated Go code
	go fmt ./...
	goimports -w .

	# Bash CLI
	rm -rf cli
	openapi-generator-cli generate --input-spec openapi.yaml \
		--generator-name bash \
		--output generated/cli \
		--additional-properties=generateZshCompletion=true \
		--additional-properties=scriptName=urlshortener-cli

	# JMeter test plan
	rm -rf jmeter
	openapi-generator-cli generate --input-spec openapi.yaml \
		--generator-name jmeter \
		--output generated/jmeter

	# Postman collection
	rm -rf postman
	openapi-generator-cli generate --input-spec openapi.yaml \
		--generator-name postman-collection \
		--output generated/postman

.PHONY: sqlc-generate
sqlc-generate: ## Generate type-safe Go DB code from db/query.sql
	sqlc generate --file db/sqlc.yaml

##@ Build & dev server

.PHONY: build
build: ## Build the urlshortener binaries
	go build ./...

# Tip: `export $$(cat .env | xargs) && make start-dev-server` to load env vars first.
.PHONY: start-dev-server
start-dev-server: ## Hot-reload dev server via reflex
	reflex -r .go -s --decoration="fancy" -- go run ./...

##@ Docker Compose

.PHONY: start
start: ## docker compose up -d
	docker compose up -d

.PHONY: compose-up
compose-up: ## Recreate the service with docker compose (down + up, foreground)
	docker compose down && docker compose up

##@ Testing

.PHONY: test
test: ## Run all Go tests
	go test ./...

.PHONY: contract-test
contract-test: ## Run Pact provider verification
	../scripts/pact-check.sh
	go test ./contracts -run TestProvider_VerifyUrlShortener -count=1

##@ Auxiliary tools

.PHONY: start-cady
start-cady: ## Reverse-proxy localhost:4545 -> localhost:8080 via Caddy
	caddy reverse-proxy --from localhost:4545 --to localhost:8080

.PHONY: start-jmeter
start-jmeter: ## Run the JMeter UserApi plan and tail its log
	jmeter -H localhost -P 9091 -t jmeter/UserApi.jmx & tail -f jmeter.log

.PHONY: build-rust-client
build-rust-client: ## Build the generated Rust client
	cargo build --manifest-path generated/rust-client/Cargo.toml

.PHONY: run-rust-client
run-rust-client: ## Run the generated Rust client
	cargo run --manifest-path generated/rust-client/Cargo.toml

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
	docker build . -f Integration.Dockerfile --platform=linux/arm64 -t igomeza/urlshortenertests
	docker push igomeza/urlshortenertests

# Requires SSH_PORT, SSH_USERNAME and SSH_HOST to be set in the environment.
.PHONY: deploy
deploy: ## SSH into the host and run ./deploy.sh remotely
	ssh -o StrictHostKeyChecking=no \
		-p $(SSH_PORT) \
		$(SSH_USERNAME)@$(SSH_HOST) \
		'bash -s' < ./deploy.sh
