.PHONY: setup
setup:
	go mod download

.PHONY: build
build:
	go build -ldflags="-w -s" -o ./bin/main ./cmd/user-rest-service/main.go

.PHONY: run
run:
	env GO_ENV=local LOCAL_SERVER_PORT=8080 LOCAL_CORS_ALLOWED_ORIGINS="http://localhost:3000" LOCAL_COOKIE_DOMAIN="localhost" LOCAL_ACCOUNT_HOST="localhost" LOCAL_ACCOUNT_PORT=8081 LOCAL_MYSQL_DSN="root:root@tcp(127.0.0.1:3306)/test_db" LOCAL_REDIS_DSN="127.0.0.1:6379" go run cmd/user-rest-service/main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	env GO_ENV=local LOCAL_SERVER_PORT=8080 LOCAL_CORS_ALLOWED_ORIGINS="http://localhost:3000" LOCAL_COOKIE_DOMAIN="localhost" LOCAL_ACCOUNT_HOST="localhost" LOCAL_ACCOUNT_PORT=8081 LOCAL_MYSQL_DSN="root:root@tcp(127.0.0.1:3306)/test_db" LOCAL_REDIS_DSN="127.0.0.1:6379" go test -race -v ./...
