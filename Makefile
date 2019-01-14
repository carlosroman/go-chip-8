.DEFAULT_GOAL := test

.PHONY: build test clean info fmt

test-clean: clean
	@mkdir -p target

test: test-clean
	@go test \
	    -v \
	    -race \
	    ./...

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
