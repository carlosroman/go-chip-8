.DEFAULT_GOAL := test

.PHONY: build test test-ci clean info fmt dep

dep:
	@go mod vendor

test-clean: clean
	@mkdir -p target

test: test-clean
	@go test \
	    -v \
	    -race \
	    ./...

test-ci: test-clean
	@scripts/coverage.sh

clean:
	@rm -rf target

build: export CGO_ENABLED=0
build:
	@echo "build app"

info:
	@env | sort -i
	@go version

fmt:
	@go fmt \
	    ./...
