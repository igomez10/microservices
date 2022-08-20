generate-openapi:
	openapi-generator generate --input-spec openapi.yaml  \
	--generator-name go-server \
	--output socialappapi \
	--ignore-file-override \
	--additional-properties=packageName=socialappapi \
	--additional-properties=outputAsLibrary=false \
	--additional-properties=packageVersion=1.18.0 \
	--additional-properties=serverPort=9000 \
	--additional-properties=featureCORS=yes \
	--additional-properties=hideGenerationTimestamp=false \
	--additional-properties=onlyInterfaces=false \
	--additional-properties=router=chi
