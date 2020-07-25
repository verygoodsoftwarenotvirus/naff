GOPATH             := $(GOPATH)
COVERAGE_OUT       := coverage.out
INSTALL_PATH       := ~/.bin
EMBEDDED_PACKAGE   := embedded
GO_FORMAT          := gofmt -s -w
THIS_PKG           := gitlab.com/verygoodsoftwarenotvirus/naff
EXAMPLE_OUTPUT_DIR := example_output
CURRENT_PROJECT    := gamut
NOW                := $(shell date +%s%N)
VERSION            := $(shell git rev-parse HEAD)
VERSION_FLAG       := -ldflags "-X main.Version=$(VERSION)_$(NOW)"
EXAMPLE_APP        := cmd/example_proj/main.go

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
	go build -o $(INSTALL_PATH)/naff $(VERSION_FLAG) $(THIS_PKG)/cmd/cli

.PHONY: example_output_subdirs
example_output_subdirs:
	for dir in `go list gitlab.com/verygoodsoftwarenotvirus/todo/...`; do mkdir -p `echo $$dir | sed -r 's/gitlab\.com\/verygoodsoftwarenotvirus\/todo/example_output/g')`; done

.PHONY: clean_example_output
clean_example_output:
	rm -rf $(EXAMPLE_OUTPUT_DIR)

$(EXAMPLE_OUTPUT_DIR):
	mkdir -p $(EXAMPLE_OUTPUT_DIR)

.PHONY: clean_todo
clean_todo: clean_example_output $(EXAMPLE_OUTPUT_DIR)
	PROJECT=todo OUTPUT_DIR=$(EXAMPLE_OUTPUT_DIR) go run $(EXAMPLE_APP)

.PHONY: compare_todo
compare_todo: clean_todo
	meld $(EXAMPLE_OUTPUT_DIR) ~/src/gitlab.com/verygoodsoftwarenotvirus/todo &

.PHONY: clean_gamut
clean_gamut: clean_example_output $(EXAMPLE_OUTPUT_DIR)
	PROJECT=gamut OUTPUT_DIR=$(EXAMPLE_OUTPUT_DIR) go run $(EXAMPLE_APP)

.PHONY: compare_gamut
compare_gamut: clean_gamut
	meld $(EXAMPLE_OUTPUT_DIR) ~/src/gitlab.com/verygoodsoftwarenotvirus/gamut &

.PHONY: install-tojen
install-tojen:
	go build -o $(INSTALL_PATH)/tojen $(VERSION_FLAG) $(THIS_PKG)/forks/tojen

.PHONY: format
format:
	for file in `find $(PWD) -name '*.go'`; do $(GO_FORMAT) $$file; done

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

generate_tests:
	 gotests -template_dir development/gotests_templates -all -w # fill this out