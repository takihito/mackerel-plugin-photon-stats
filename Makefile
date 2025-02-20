GIT_VER ?= $(shell git describe --tags | sed -e 's/-/+/')
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
export GO111MODULE := on

.PHONY: test build

init_mod:
	go mod init github.com/takihito/mackerel-plugin-photon-stats
	go get github.com/mackerelio/go-mackerel-plugin-helper

build: 
	mkdir -p bin
	go build -ldflags "-s -w -X main.Version=${GIT_VER} -X main.buildDate=${DATE}" -trimpath -o bin/mackerel-plugin-photon-stats

test:
	go test -race

# use goreleaser
# brew install goreleaser/tap/goreleaser

check/release:
	goreleaser check

build/release:
	goreleaser release --snapshot --clean

