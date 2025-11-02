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
	# openapi-generator version-manager set 7.11.0
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
	reflex -r .go -s --decoration="fancy" --  go run ./...

build:
	go build ./...

test:
	go test ./...

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
