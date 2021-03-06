ifdef update
	u=-u
endif

deps:
		go get ${u} github.com/golang/dep/cmd/dep
		dep ensure

devel-deps: deps
		go get ${u} github.com/golang/lint/golint
		go get ${u} github.com/mattn/goveralls
		go get ${u} github.com/motemen/gobump
		go get ${u} github.com/laher/goxc
		go get ${u} github.com/Songmu/ghch

test: deps
		go test

lint: devel-deps
		go vet
		golint -set_exit_status

major-vaersionup:
		gobump major -w
		gobump show -r

minor-vaersionup:
		gobump minor -w
		gobump show -r

crossbuild: devel-deps
		mkdir -p ./dist
		goxc 

release:
	ghr -u takihito -r mackerel-plugin-photon-stats v$(shell gobump show -r) dist/snapshot/

install:
	mkr plugin install takihito/mackerel-plugin-photon-stats@v$(shell gobump show -r)

.PHONY: test deps devel-deps lint cover crossbuild release

