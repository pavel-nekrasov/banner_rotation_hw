
BIN := "./bin/rotator_server"
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

docker-build-test:
	docker compose -f deploy/docker-compose.test.yml build
docker-env-up-test:
	docker compose -f deploy/docker-compose.test.yml up -d rotator_db rotator_mq
docker-stop-test:
	docker compose -f deploy/docker-compose.test.yml down

docker-migrate-test:
	docker compose -f deploy/docker-compose.test.yml up rotator_test_migrate

integration-test:
	make docker-build-test
	make docker-env-up-test
	make docker-migrate-test
	docker compose -f deploy/docker-compose.test.yml up rotator_test
	make docker-stop-test



docker-build:
	docker compose -f deploy/docker-compose.yml build
docker-env-up:
	docker compose -f deploy/docker-compose.yml up -d rotator_db rotator_mq
docker-migrate:
	docker compose -f deploy/docker-compose.yml up migrate
run:
	make docker-build
	make docker-env-up
	make docker-migrate
	docker compose -f deploy/docker-compose.yml up -d rotator_server
stop:
	docker compose -f deploy/docker-compose.yml down



build-server:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/server

build: build-server



install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test -count=10 -race -timeout=5m ./internal/...



proto-generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/server/grpc/pb \
		--go-grpc_out=internal/server/grpc/pb \
		api/*.proto

.PHONY: test test-env-up test-env-down test-migrate integration-test proto-generate env-build env-up migrate run stop