
BIN := "./bin/rotator_server"
BIN_MIGRATOR := "./bin/rotator_migrator"
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

test-env-build:
	docker compose -f deploy/docker-compose.test.yml build
test-env-up:
	docker compose -f deploy/docker-compose.test.yml up -d rotator_db rotator_mq
test-env-down:
	docker compose -f deploy/docker-compose.test.yml down

test-migrate:
	docker compose -f deploy/docker-compose.test.yml up rotator_test_migrate

integration-test:
	make test-env-build
	make test-env-up
	make test-migrate
	docker compose -f deploy/docker-compose.test.yml up rotator_test
	make test-env-down

test:
	go test -count=10 -race -timeout=5m ./internal/...


env-build:
	docker compose -f deploy/docker-compose.yml build
env-up:
	docker compose -f deploy/docker-compose.yml up -d rotator_db rotator_mq
migrate:
	docker compose -f deploy/docker-compose.yml up migrate

run:
	make env-build
	make env-up
	make migrate
	docker compose -f deploy/docker-compose.yml up -d rotator_server
stop:
	docker compose -f deploy/docker-compose.yml down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

build-migrator:
	go build -v -o $(BIN_MIGRATOR) -ldflags "$(LDFLAGS)" ./cmd/migrate

build-server:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/server

build: build-server build-migrator

proto-generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/server/grpc/pb \
		--go-grpc_out=internal/server/grpc/pb \
		api/*.proto

.PHONY: test test-env-up test-env-down test-migrate integration-test proto-generate env-build env-up migrate run stop