.PHONY: client

client:
	openapi-generator generate -i ../docs/openapi/swagger.yaml -g go  -o ./client --additional-properties packageName=client
	rm ./client/go.mod
	rm ./client/go.sum