.DEFAULT_GOAL := test

.PHONY: build test test-ci clean info fmt build-pixel

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

clean-pixel: test-clean

build-pixel: clean-pixel
build-pixel: export CGO_ENABLED=1
build-pixel: export CXX=x86_64-w64-mingw32-g++
build-pixel: export CC=x86_64-w64-mingw32-gcc
build-pixel: export GOOS=windows
build-pixel:
	@go build \
	    -o target/pixel.exe \
	    cmd/pixel/main.go

info:
	@env | sort -i
	@go version

fmt:
	@go fmt \
	    ./...
