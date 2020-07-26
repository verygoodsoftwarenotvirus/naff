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

ensure-goimports:
ifndef $(shell command -v goimports 2> /dev/null)
	$(shell GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports)
endif

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

## tests

.PHONY: quicktest
quicktest:
	go test -race -failfast `go list gitlab.com/verygoodsoftwarenotvirus/naff/... | grep -Ev '(cmd|testprojects|example_models|example_output|forks)'`

.PHONY: install
install:
	go build -o $(INSTALL_PATH)/naff $(VERSION_FLAG) $(THIS_PKG)/cmd/cli

## local run testing

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

## housekeeping

.PHONY: format
format:
	for file in `find $(PWD) -name '*.go'`; do $(GO_FORMAT) $$file; done

.PHONY: docker_image
docker_image:
	docker build --tag naff:latest --file Dockerfile .

generate_tests:
	 gotests -template_dir development/gotests_templates -all -w # fill this out