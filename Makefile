GOPATH            := $(GOPATH)
GO_PACKAGE        := gitlab.com/verygoodsoftwarenotvirus/naff
COVERAGE_OUT      := coverage.out
INSTALL_PATH      := ~/.bin
EMBEDDED_PACKAGE  := embedded
GO_FORMAT         := gofmt -s -w

VERSION := $(shell git rev-parse --short HEAD)

## Project prerequisites
.PHONY: deps
deps:
	GO111MODULE=off go get -u github.com/UnnoTed/fileb0x

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

template-dirs:
	for dir in `find ../todo -type d -not -path "*\.git*" -not -path "*node_modules*" -not -path "*vendor*" | sed -r 's/\.\.\/todo\/?/templates\/experimental\//g'`; do mkdir -p $$dir; done
	for dir in `find ../todo -type d -not -path "*\.git*" -not -path "*node_modules*" -not -path "*vendor*" | sed -r 's/\.\.\/todo\/?/templates\/blessed\//g'`; do mkdir -p $$dir; done

fix-template-test-files:
	find templates/ -name "*_test.go" -exec bash -c 'mv "$1" `echo "$1" | sed -r "s/_test\.go/test_\.go/g"` ' - '{}' \;

.PHONY: $(EMBEDDED_PACKAGE)
$(EMBEDDED_PACKAGE): templates
	fileb0x b0x.yaml

.PHONY: install
install: $(EMBEDDED_PACKAGE)
	go build -o $(INSTALL_PATH)/naff -ldflags "-X main.Version=$(VERSION)" $(GO_PACKAGE)/cmd/cli

.PHONY: example_output_subdirs
example_output_subdirs:
	for dir in `go list gitlab.com/verygoodsoftwarenotvirus/todo/...`; do mkdir -p `echo $$dir | sed -r 's/gitlab\.com\/verygoodsoftwarenotvirus\/todo/example_output/g')`; done

.PHONY: example_output
example_output:
	@go run cmd/todoproj/main.go

.PHONY: install-tojen
install-tojen:
	go build -o $(INSTALL_PATH)/tojen -ldflags "-X main.Version=$(VERSION)" $(GO_PACKAGE)/forks/tojen

.PHONY: format
format:
	for file in `find $(PWD) -name '*.go'`; do $(GO_FORMAT) $$file; done
