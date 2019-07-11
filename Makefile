GOPATH            := $(GOPATH)
GO_PACKAGE        := gitlab.com/verygoodsoftwarenotvirus/naff
COVERAGE_OUT      := coverage.out
INSTALL_PATH      := ~/.bin

VERSION := $(shell git rev-parse --short HEAD)

## generic make stuff
.PHONY: clean
clean:
	rm -f naff_debug

## Project prerequisites
vendor:
	docker run --env GO111MODULE=on --env GOPATH=$(GOPATH) --volume `pwd`:`pwd` --workdir=`pwd` --workdir=`pwd` golang:latest /bin/sh -c "go mod vendor"

.PHONY: revendor
revendor:
	rm -rf vendor go.{mod,sum}
	GO111MODULE=on go mod init
	$(MAKE) vendor

.PHONY: test
test:
	docker build --tag coverage-todo:latest --file dockerfiles/coverage.Dockerfile .
	docker run --rm --volume `pwd`:`pwd` --workdir=`pwd` coverage-todo:latest

.PHONY: run
run:
	go run $(GO_PACKAGE)/cmd/cli generate

naff_debug:
	go build -o naff_debug $(GO_PACKAGE)/cmd/cli

templates:
	@rm -rf template/
	@mkdir template
	go run cmd/tools/template_builder/main.go

.PHONY: install
install:
	go build -o $(INSTALL_PATH)/naff -ldflags "-X main.Version=$(VERSION)" $(GO_PACKAGE)/cmd/cli

