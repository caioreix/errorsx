.PHONY: mocks

ROOT_PATH := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

install:
	@go install github.com/vektra/mockery/v2@v2.53.3

test:
	@go test $(go list ./... | grep -v _mock) -race -cover -coverprofile=coverage.out -count=1
	@go tool cover -func coverage.out

mocks:
	@mockery --config $(ROOT_PATH)/.mockery.yml
