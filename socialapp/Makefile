generate-openapi:
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

sqlc-generate:
	sqlc generate --file db/sqlc.yaml

restart-db:
	docker compose down && docker-compose up

start-dev-server:
	reflex -r .go -s --  go run ./...
