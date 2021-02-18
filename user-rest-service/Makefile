RED=\033[31m
GREEN=\033[32m
RESET=\033[0m

COLORIZE_PASS=sed ''/PASS/s//$$(printf "$(GREEN)PASS$(RESET)")/''
COLORIZE_FAIL=sed ''/FAIL/s//$$(printf "$(RED)FAIL$(RESET)")/''

.PHONY: setup
setup:
	go mod download

.PHONY: build
build:
	go build -ldflags="-w -s" -o ./bin/main ./cmd/user-rest-service/main.go

.PHONY: run
run:
	GO_ENV=local LOCAL_SERVER_PORT=8080 LOCAL_CORS_ALLOWED_ORIGINS="http://localhost:3000" LOCAL_COOKIE_NAME="local_session_id" LOCAL_COOKIE_DOMAIN="localhost" LOCAL_COOKIE_SECURE=false LOCAL_COOKIE_EXPIRATION="2592000s" LOCAL_ACCOUNT_HOST="localhost" LOCAL_ACCOUNT_PORT=8081 LOCAL_MYSQL_DSN="root:root@tcp(127.0.0.1:3306)/test_db" LOCAL_REDIS_DSN="127.0.0.1:6379" go run cmd/user-rest-service/main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	GO_ENV=local LOCAL_SERVER_PORT=8080 LOCAL_CORS_ALLOWED_ORIGINS="http://localhost:3000" LOCAL_COOKIE_NAME="local_session_id" LOCAL_COOKIE_DOMAIN="localhost" LOCAL_COOKIE_SECURE=false LOCAL_COOKIE_EXPIRATION="2592000s" LOCAL_ACCOUNT_HOST="localhost" LOCAL_ACCOUNT_PORT=8081 LOCAL_MYSQL_DSN="root:root@tcp(127.0.0.1:3306)/test_db" LOCAL_REDIS_DSN="127.0.0.1:6379" go test -race -v ./... | $(COLORIZE_PASS) | $(COLORIZE_FAIL)
