generate-openapi:
	#  generate server:
	openapi-generator generate --input-spec openapi.yaml  \
	--generator-name go-server \
	--output socialappapi \
	--ignore-file-override \
	--additional-properties=packageName=socialappapi \
	--additional-properties=outputAsLibrary=true \
	--additional-properties=packageVersion=1.18.0 \
	--additional-properties=serverPort=8080 \
	--additional-properties=featureCORS=true \
	--additional-properties=hideGenerationTimestamp=false \
	--additional-properties=onlyInterfaces=false \
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
	openapi-generator generate -i openapi.yaml  \
	--generator-name mysql-schema \
	--output generated-schema \
	--additional-properties=namedParametersEnabled=true

	go fmt ./...
	
sqlc-generate:
	sqlc generate --file db/sqlc.yaml

compose-up:
	docker compose down && docker-compose up

start-dev-server:
	reflex -r .go -s --  go run ./...

build:
	go build ./...

test:
	go test ./...
