GOPATH             := $(GOPATH)
COVERAGE_OUT       := coverage.out
INSTALL_PATH       := ~/.bin
EMBEDDED_PACKAGE   := embedded
GO_FORMAT          := gofmt -s -w
THIS_PKG           := gitlab.com/verygoodsoftwarenotvirus/naff
VERSION            := $(shell git rev-parse HEAD)
EXAMPLE_OUTPUT_DIR := example_output

## Project prerequisites
.PHONY: deps
deps:
	GO111MODULE=off go get -u github.com/UnnoTed/fileb0x

.PHONY: vendor-clean
vendor-clean:
	rm -rf vendor go.sum

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: revendor
revendor: vendor-clean vendor

.PHONY: run
run:
	go run $(THIS_PKG)/cmd/cli generate

naff_debug:
	go build -o naff_debug $(THIS_PKG)/cmd/cli

.PHONY: test
test: test-models test-wordsmith test-http-client

.PHONY: test-wordsmith
test-wordsmith:
	go test -race -cover -v ./lib/wordsmith/

.PHONY: test-models
test-models:
	go test -race -cover -v ./models/

.PHONY: test-http-client
test-http-client:
	go test -race -cover -v ./templates/blessed/client/v1/http/

.PHONY: install
install:
	go build -o $(INSTALL_PATH)/naff -ldflags "-X main.Version=$(VERSION)" $(THIS_PKG)/cmd/cli

.PHONY: example_output_subdirs
example_output_subdirs:
	for dir in `go list gitlab.com/verygoodsoftwarenotvirus/todo/...`; do mkdir -p `echo $$dir | sed -r 's/gitlab\.com\/verygoodsoftwarenotvirus\/todo/example_output/g')`; done

$(EXAMPLE_OUTPUT_DIR):
	go run cmd/todoproj/main.go

.PHONY: clean_example_output
clean_example_output:
	rm -rf $(EXAMPLE_OUTPUT_DIR)

.PHONY: install-tojen
install-tojen:
	go build -o $(INSTALL_PATH)/tojen -ldflags "-X main.Version=$(VERSION)" $(THIS_PKG)/forks/tojen

.PHONY: format
format:
	for file in `find $(PWD) -name '*.go'`; do $(GO_FORMAT) $$file; done

.PHONY: compare
compare: clean_example_output $(EXAMPLE_OUTPUT_DIR)
	meld $(EXAMPLE_OUTPUT_DIR) ~/src/gitlab.com/verygoodsoftwarenotvirus/todo &

.PHONY: docker_image
docker_image:
	docker build --tag naff:latest --file Dockerfile .

.PHONY: example_run
example_run: clean_example_output $(EXAMPLE_OUTPUT_DIR)
	(cd $(EXAMPLE_OUTPUT_DIR) && $(MAKE) rewire quicktest integration-tests-postgres)
	@# (cd $(EXAMPLE_OUTPUT_DIR) && $(MAKE) rewire && go test -cover -v gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/independents)

ensure-goimports:
ifndef $(shell command -v goimports 2> /dev/null)
	$(shell GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports)
endif
