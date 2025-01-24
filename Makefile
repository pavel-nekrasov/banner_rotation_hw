
test-env-build:
	docker compose -f deploy/docker-compose.test.yml build
test-env-up:
	docker compose -f deploy/docker-compose.test.yml up -d rotator_db rotator_mq
test-env-down:
	docker compose -f deploy/docker-compose.test.yml down

test-migrate:
	docker compose -f deploy/docker-compose.test.yml up migrate

integration-test:
	make test-env-build
	make test-env-up
	make test-migrate
	docker compose -f deploy/docker-compose.test.yml up rotator_test
	make test-env-down

test:
	go test -race -v -count=1 -race -timeout=1m ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: test test-env-up test-env-down integration-test