.DEFAULT_GOAL := test

.PHONY: build test test-ci clean info fmt

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

build/ebiten : clean
build/ebiten : export CGO_ENABLED=1
build/ebiten : export CXX=x86_64-w64-mingw32-g++
build/ebiten : export CC=x86_64-w64-mingw32-gcc
build/ebiten : export GOOS=windows
build/ebiten :
	@go build \
	    -o target/ebiten.exe \
	    cmd/ebiten/main.go
