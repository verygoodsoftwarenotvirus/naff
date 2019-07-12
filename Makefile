GOPATH            := $(GOPATH)
GO_PACKAGE        := gitlab.com/verygoodsoftwarenotvirus/naff
COVERAGE_OUT      := coverage.out
INSTALL_PATH      := ~/.bin

VERSION := $(shell git rev-parse --short HEAD)

## Project prerequisites
.PHONY: vendor-clean
vendor-clean:
	rm -rf vendor go.sum

.PHONY: vendor
vendor: template-clean
	GO111MODULE=on go mod vendor

.PHONY: revendor
revendor: vendor-clean vendor

.PHONY: test
test:
	docker build --tag coverage-todo:latest --file dockerfiles/coverage.Dockerfile .
	docker run --rm --volume `pwd`:`pwd` --workdir=`pwd` coverage-todo:latest

.PHONY: run
run:
	go run $(GO_PACKAGE)/cmd/cli generate

naff_debug:
	go build -o naff_debug $(GO_PACKAGE)/cmd/cli

.PHONY: template-clean
template-clean:
	rm -rf template
	mkdir -p template

templates: template-clean
	go run cmd/tools/template_builder/main.go

.PHONY: install
install:
	go build -o $(INSTALL_PATH)/naff -ldflags "-X main.Version=$(VERSION)" $(GO_PACKAGE)/cmd/cli

