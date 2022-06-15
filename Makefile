deploy:
	docker buildx build --platform linux/arm64 -t kpture/admission:latest .
	docker push kpture/admission:latest
	helm uninstall tls || true
	helm install tls  ./charts/tls
dependencies:
	go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	swag init -g docs/global_openapi.go -ot yaml  --parseDependency
	swag fmt
unit_test:
	go test -run=TestUnit -bench=. -benchmem -v -cover ./...
profile:
	go tool pprof -http=:9090 -no_browser -nodecount 30 http://localhost:8080/debug/pprof/heap