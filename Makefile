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


# ifdef update
#         u=-u
# endif
# 
# deps:
#                 go get ${u} github.com/golang/dep/cmd/dep
#                 dep ensure
# 
# devel-deps: deps
#                 go get ${u} github.com/golang/lint/golint
#                 go get ${u} github.com/mattn/goveralls
#                 go get ${u} github.com/motemen/gobump
#                 go get ${u} github.com/laher/goxc
#                 go get ${u} github.com/Songmu/ghch
# 
# test: deps
#                 go test
# 
# lint: devel-deps
#                 go vet
#                 golint -set_exit_status
# 
# major-vaersionup:
#                 gobump major -w
#                 gobump show -r
# 
# minor-vaersionup:
#                 gobump minor -w
#                 gobump show -r
# 
# crossbuild: devel-deps
#                 mkdir -p ./dist
#                 goxc 
# 
# release:
#         ghr -u takihito -r mackerel-plugin-photon-stats v$(shell gobump show -r) dist/snapshot/
# 
# install:
#         mkr plugin install takihito/mackerel-plugin-photon-stats@v$(shell gobump show -r)
# 
# .PHONY: test deps devel-deps lint cover crossbuild release
