GO_EXECUTABLE ?= go
VERSION = `git describe --always --tags --abbrev=0 | tr -d "[\r\n]"`
TIME = `date +%FT%T%z`
MODULE = github.com/musenwill/exercise

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
BINARY := "exercise"

LDFLAGS= -ldflags "-X ${MODULE}/common.Version=${VERSION} -X ${MODULE}/common.BuildTime=${TIME} -X ${MODULE}/common.AppName=${BINARY}"

UNAME = $(shell uname)
ifeq (${UNAME}, Darwin)
	os=darwin
else
	os=linux
endif

build: check game-24
	${GO_EXECUTABLE} build ${LDFLAGS} -o ${BINARY}

check:
	golint ./... | grep -v "exported" | exit 0
	go vet ./...
	gofmt -d -s `find . -name "*.go" -type f`
	go test ./...

game-24:
	${GO_EXECUTABLE} build ${LDFLAGS} -o game24/main/game24 game24/main/main.go

clean:
	rm -f ${BINARY}
	rm -rf dist

build-all:
	gox -verbose \
	${LDFLAGS} \
	-os="linux darwin windows freebsd openbsd netbsd" \
	-arch="amd64 386 armv5 armv6 armv7 arm64" \
	-osarch="!darwin/arm64" \
	-output="dist/{{.OS}}-{{.Arch}}/${BINARY}" .

build-os:
	gox -verbose \
	${LDFLAGS} \
	-os="${os}" \
	-arch="amd64" \
	-output="dist/{{.OS}}-{{.Arch}}/${BINARY}" .

.PHONY: build build-all build-os clean check
