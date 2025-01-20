
test_env_up:
	docker compose -f docker-compose.test.yml up -d
test_env_down:
	docker compose -f docker-compose.test.yml down
test_run: test_env_up
	go test -race ./internal/...

test: test_env_up test_run test_env_down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: test test_env_up test_run test_env_down