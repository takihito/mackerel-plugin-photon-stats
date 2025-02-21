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

add_release_version_tag:
	export RELEASE_VERSION=`go run . version | cut -f 2 -d " "` && git tag -a v$$RELEASE_VERSION -m "Release v$$RELEASE_VERSION" && git push origin v$$RELEASE_VERSION

check/release:
	goreleaser check

build/release:
	goreleaser release --snapshot --clean

