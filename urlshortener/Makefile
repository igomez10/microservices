generate-openapi:
	#  generate server:
	rm -rf ./generated/server
	openapi-generator generate -i openapi.yaml  \
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

	# generate go client
	rm -rf generated/clients/go/client
	openapi-generator generate --input-spec openapi.yaml  \
		--git-user-id=igomez10 \
		--git-repo-id=microservices \
        --generator-name go \
        --output generated/clients/go/client \
		-p packageName=client \
		--additional-properties=packageName=client \
		--additional-properties=withGoMod=false \
		--additional-properties=withGoCodegenComment=false \
		--additional-properties=withGoTests=false

	
	rm generated/clients/go/client/test/* # remove tests in client sdk

	# generate documentation
	openapi-generator generate -i openapi.yaml  \
	--output generated/docs \
	--generator-name markdown
	
	# generate schema
	rm -rf mysql-schema
	openapi-generator generate -i openapi.yaml  \
	--generator-name mysql-schema \
	--output generated/schema \
	--additional-properties=namedParametersEnabled=true \
	--additional-properties=identifierNamingConvention=snake_case

	# format go code
	go fmt ./...
	goimports -w .
	
	# generate cli 
	rm -rf cli
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name bash \
        --output generated/cli \
		--additional-properties=generateZshCompletion=true \
		--additional-properties=scriptName=urlshortener-cli

	# generate jmeter
	rm -rf jmeter
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name jmeter \
		--output generated/jmeter

	# generate postman collection
	rm -rf postman
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name postman-collection \
		--output generated/postman


sqlc-generate:
	sqlc generate --file db/sqlc.yaml

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
	cargo build --manifest-path generated/rust-client/Cargo.toml
run-rust-client:
	cargo run --manifest-path generated/rust-client/Cargo.toml

tf-init:
	cd infrastructure && terraform init

tf-plan:
	cd infrastructure && terraform plan

tf-apply:
	cd infrastructure && terraform apply -auto-approve

build-push-integration-tests:
	docker build . -f Integration.Dockerfile --platform=linux/arm64 -t igomeza/urlshortenertests 	
	docker push igomeza/urlshortenertests

deploy:
	# ssh into the server and run the deploy script
	 ssh -o StrictHostKeyChecking=no \
            -p $(SSH_PORT) \
            $(SSH_USERNAME)@$(SSH_HOST) \
			'bash -s' < ./deploy.sh
