GOPATH            := $(GOPATH)
COVERAGE_OUT      := coverage.out
INSTALL_PATH      := ~/.bin
EMBEDDED_PACKAGE  := embedded
GO_FORMAT         := gofmt -s -w
THIS_PKG          := gitlab.com/verygoodsoftwarenotvirus/naff
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

.PHONY: run
run:
	go run $(THIS_PKG)/cmd/cli generate

naff_debug:
	go build -o naff_debug $(THIS_PKG)/cmd/cli

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

.PHONY: test
test: test-wordsmith

.PHONY: test-wordsmith
test-wordsmith:
	go test -v ./lib/wordsmith/

.PHONY: $(EMBEDDED_PACKAGE)
$(EMBEDDED_PACKAGE): templates
	fileb0x b0x.yaml

.PHONY: install
install: $(EMBEDDED_PACKAGE)
	go build -o $(INSTALL_PATH)/naff -ldflags "-X main.Version=$(VERSION)" $(THIS_PKG)/cmd/cli

.PHONY: example_output_subdirs
example_output_subdirs:
	for dir in `go list gitlab.com/verygoodsoftwarenotvirus/todo/...`; do mkdir -p `echo $$dir | sed -r 's/gitlab\.com\/verygoodsoftwarenotvirus\/todo/example_output/g')`; done

.PHONY: example_output
example_output:
	@go run cmd/todoproj/main.go

.PHONY: install-tojen
install-tojen:
	go build -o $(INSTALL_PATH)/tojen -ldflags "-X main.Version=$(VERSION)" $(THIS_PKG)/forks/tojen

.PHONY: format
format:
	for file in `find $(PWD) -name '*.go'`; do $(GO_FORMAT) $$file; done

.PHONY: compare
compare:
	meld example_output ~/src/gitlab.com/verygoodsoftwarenotvirus/todo &
