dependencies:
	go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	swag init -g docs/global_openapi.go -ot yaml -o docs/openapi/ --parseDependency
	swag fmt
fmt:
	find . -name "*.go" | xargs -L1 gofumpt -e -w
	find . -name "*.go" | xargs -L1 gofumports -w



dev_devspace:
	devspace --config ./tools/devspace.yaml dev
purge_devspace:
	devspace --config ./tools/devspace.yaml purge