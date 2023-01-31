generate-openapi:
	# openapi-generator-cli version-manager set 6.2.0
	#  generate server:
	rm -rf ./socialappapi
	openapi-generator-cli generate --input-spec openapi.yaml  \
	--generator-name go-server \
	--output socialappapi \
	--ignore-file-override \
	--additional-properties=packageName=socialappapi \
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
	openapi-generator-cli generate --input-spec openapi.yaml  \
        --generator-name go \
        --output client \
		--additional-properties=isGoSubmodule=true \
		--additional-properties=packageName=client \
		--additional-properties=structPrefix=true
	rm client/go.mod
	rm client/go.sum

	# generate documentation
	openapi-generator-cli generate -i openapi.yaml  \
	--generator-name markdown
	
	# generate schema
	rm -rf mysql-schema
	openapi-generator-cli generate -i openapi.yaml  \
	--generator-name mysql-schema \
	--output generated-schema \
	--additional-properties=namedParametersEnabled=true \
	--additional-properties=identifierNamingConvention=snake_case

	# format go code
	go fmt ./...
	goimports -w .
	
	# generate cli 
	rm -rf cli
	openapi-generator-cli generate --input-spec openapi.yaml \
		--generator-name bash \
        --output cli \
		--additional-properties=generateZshCompletion=true \
		--additional-properties=scriptName=socialapp-cli

	# generate jmeter
	rm -rf jmeter
	openapi-generator-cli generate --input-spec openapi.yaml \
		--generator-name jmeter \
		--output jmeter


generate-protobuf:
	# generate protobuf
	rm -rf protobuf-schema
	openapi-generator-cli generate --input-spec openapi.yaml \
		--generator-name protobuf-schema \
		--output protobuf-schema

	# add at the end of every file option go_package = \"github.com/igomez10/socialapp-go-client/openapitools\";
	
	#  for every proto file create a package with its name and add the option go_package
	#  to the file
	for file in protobuf-schema/models/*.proto; do \
		package_name=$$(basename $$file .proto); \
		echo "package $$package_name;" | cat - $$file > temp && mv temp $$file; \
	done

	echo "option go_package = \"./models\";" | tee -a protobuf-schema/models/*.proto
	echo "option go_package = \"./services\";" | tee -a protobuf-schema/services/*.proto

	cd protobuf-schema && \
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative ./services/*.proto
	
	cd protobuf-schema && \
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative ./models/*.proto
	echo "FINISHED GENERATING PROTOBUF"

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
