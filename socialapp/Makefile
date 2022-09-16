generate-openapi:
	#  generate server:
	rm -rf ./socialappapi
	openapi-generator generate --input-spec openapi.yaml  \
	--generator-name go-server \
	--output socialappapi \
	--ignore-file-override \
	--additional-properties=packageName=socialappapi \
	--additional-properties=outputAsLibrary=true \
	--additional-properties=packageVersion=1.18.0 \
	--additional-properties=serverPort=8080 \
	--additional-properties=featureCORS=true \
	--additional-properties=hideGenerationTimestamp=true \
	--additional-properties=onlyInterfaces=true \
	--additional-properties=router=chi
	rm -rf socialappapi/openapi
	mv socialappapi/go socialappapi/openapi
	
	# generate client
	openapi-generator generate --input-spec openapi.yaml  \
        --generator-name go \
        --output client \
		--additional-properties=isGoSubmodule=true \
		--additional-properties=packageName=client \
		--additional-properties=structPrefix=true
	rm client/go.mod
	rm client/go.sum

	# generate documentation
	openapi-generator generate -i openapi.yaml  \
	--generator-name markdown
	
	# generate schema
	rm -rf mysql-schema
	openapi-generator generate -i openapi.yaml  \
	--generator-name mysql-schema \
	--output generated-schema \
	--additional-properties=namedParametersEnabled=true
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
	
	## generate k6
	rm -rf k6
	openapi-generator generate --input-spec openapi.yaml \
		--generator-name k6 \
		--output k6


sqlc-generate:
	sqlc generate --file db/sqlc.yaml

compose-up:
	docker compose down && docker-compose up

start-dev-server:
	reflex -r .go -s --decoration="fancy" --  go run ./...

build:
	go build ./...

test:
	go test ./...

start-cady:
	 caddy reverse-proxy --from localhost:4545  --to localhost:8080
