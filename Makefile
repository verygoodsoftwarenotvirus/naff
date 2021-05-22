GOPATH             := $(GOPATH)
COVERAGE_OUT       := coverage.out
INSTALL_PATH       := ~/.bin
EMBEDDED_PACKAGE   := embedded
GO_FORMAT          := gofmt -s -w
THIS_PKG           := gitlab.com/verygoodsoftwarenotvirus/naff
ARTIFACTS_DIR      := artifacts
COVERAGE_OUT       := $(ARTIFACTS_DIR)/coverage.out
PACKAGE_LIST       := `go list gitlab.com/verygoodsoftwarenotvirus/naff/... | grep -Ev '(cmd|testprojects|example_models|example_output|forks)'`
EXAMPLE_OUTPUT_DIR := example_output
CURRENT_PROJECT    := gamut
NOW                := $(shell date +%s%N)
VERSION            := $(shell git rev-parse HEAD)
VERSION_FLAG       := -ldflags "-X main.Version=$(VERSION)_$(NOW)"
EXAMPLE_APP        := cmd/example_proj/main.go

$(ARTIFACTS_DIR):
	@mkdir -p $(ARTIFACTS_DIR)

## Project prerequisites

ensure-goimports:
ifndef $(shell command -v goimports 2> /dev/null)
	$(shell GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports)
endif

.PHONY: vendor-clean
vendor-clean:
	rm -rf vendor go.sum

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: revendor
revendor: vendor-clean vendor

.PHONY: install
install:
	go build -o $(INSTALL_PATH)/naff $(VERSION_FLAG) $(THIS_PKG)/cmd/cli

.PHONY: install-tojen-fork
install-tojen-fork:
	go build -o $(INSTALL_PATH)/tojen $(VERSION_FLAG) $(THIS_PKG)/forks/tojen

## tests

.PHONY: clean-coverage
clean-coverage:
	@rm -f $(COVERAGE_OUT) profile.out;

test: coverage

.PHONY: coverage
coverage: clean-coverage $(ARTIFACTS_DIR)
	go test -coverprofile=$(COVERAGE_OUT) -covermode=atomic -race $(PACKAGE_LIST)
	go tool cover -func=$(ARTIFACTS_DIR)/coverage.out | grep 'total:' | xargs | awk '{ print "COVERAGE: " $$3 }'

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

.PHONY: clean_every_type
clean_every_type: clean_example_output $(EXAMPLE_OUTPUT_DIR)
	PROJECT=every_type OUTPUT_DIR=$(EXAMPLE_OUTPUT_DIR) go run $(EXAMPLE_APP)

.PHONY: clean_forums
clean_forums: clean_example_output $(EXAMPLE_OUTPUT_DIR)
	PROJECT=forums OUTPUT_DIR=$(EXAMPLE_OUTPUT_DIR) go run $(EXAMPLE_APP)

.PHONY: compare_gamut
compare_gamut: clean_gamut
	meld $(EXAMPLE_OUTPUT_DIR) ~/src/gitlab.com/verygoodsoftwarenotvirus/gamut &

## CI output tests

.PHONY: ci-todo-project-unit-tests
ci-todo-project-unit-tests: clean_todo
	(cd $(EXAMPLE_OUTPUT_DIR) && $(MAKE) revendor rewire config_files quicktest)

.PHONY: ci-every-type-project-unit-tests
ci-every-type-project-unit-tests: clean_every_type
	(cd $(EXAMPLE_OUTPUT_DIR) && $(MAKE) revendor rewire config_files quicktest)

.PHONY: ci-forums-project-unit-tests
ci-forums-project-unit-tests: clean_forums
	(cd $(EXAMPLE_OUTPUT_DIR) && $(MAKE) revendor rewire config_files quicktest)

.PHONY: ci-todo-project-integration-tests
ci-todo-project-integration-tests: clean_todo
	(cd $(EXAMPLE_OUTPUT_DIR) && $(MAKE) revendor rewire config_files lintegration-tests)

.PHONY: ci-every-type-project-integration-tests
ci-every-type-project-integration-tests: clean_every_type
	(cd $(EXAMPLE_OUTPUT_DIR) && $(MAKE) revendor rewire config_files lintegration-tests)

## housekeeping

.PHONY: format
format:
	for file in `find $(PWD) -name '*.go'`; do $(GO_FORMAT) $$file; done

.PHONY: docker_image
docker_image:
	docker build --tag naff:latest --file Dockerfile .

# example commands

example_generate_tests:
	gotests -template_dir development/gotests_templates -all -w # relative path directory