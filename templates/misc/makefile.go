package misc

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func makefile(proj *models.Project) string {
	return fmt.Sprintf(`PWD                           := $(shell pwd)
GOPATH                        := $(GOPATH)
ARTIFACTS_DIR                 := artifacts
COVERAGE_OUT                  := $(ARTIFACTS_DIR)/coverage.out
SEARCH_INDICES_DIR            := $(ARTIFACTS_DIR)/search_indices
GO                            := docker run --interactive --tty --volume $(PWD):$(PWD) --workdir $(PWD) --user $(shell id -u):$(shell id -g) golang:1.17-stretch go
GO_FORMAT                     := gofmt -s -w
THIS                          := %s
TOTAL_PACKAGE_LIST            := `+"`"+`go list $(THIS)/...`+"`"+`
TESTABLE_PACKAGE_LIST         := `, proj.OutputPath) + "`" + `go list $(THIS)/... | grep -Ev '(cmd|tests|testutil|mock|fake)'` + "`" + fmt.Sprintf(`
ENVIRONMENTS_DIR              := environments
TEST_ENVIRONMENT_DIR          := $(ENVIRONMENTS_DIR)/testing
TEST_DOCKER_COMPOSE_FILES_DIR := $(TEST_ENVIRONMENT_DIR)/compose_files
FRONTEND_DIR                  := frontend
FRONTEND_TOOL                 := pnpm

## non-PHONY folders/files

clear:
	@printf "\033[2J\033[3J\033[1;1H"

clean:
	rm -rf $(ARTIFACTS_DIR)

$(ARTIFACTS_DIR):
	@mkdir --parents $(ARTIFACTS_DIR)

clean-$(ARTIFACTS_DIR):
	@rm -rf $(ARTIFACTS_DIR)

$(SEARCH_INDICES_DIR):
	@mkdir --parents $(SEARCH_INDICES_DIR)

clean-search-indices:
	@rm -rf $(SEARCH_INDICES_DIR)

.PHONY: setup
setup: $(ARTIFACTS_DIR) $(SEARCH_INDICES_DIR) revendor frontend_vendor rewire configs

.PHONY: configs
configs:
	go run cmd/tools/config_gen/main.go

## prerequisites

# not a bad idea to do this either:
## GO111MODULE=off go install golang.org/x/tools/...

ensure_wire_installed:
ifndef $(shell command -v wire 2> /dev/null)
	$(shell GO111MODULE=off go install github.com/google/wire/cmd/wire@latest)
endif

ensure_fieldalign_installed:
ifndef $(shell command -v wire 2> /dev/null)
	$(shell GO111MODULE=off go get -u golang.org/x/tools/...)
endif

ensure_scc_installed:
ifndef $(shell command -v scc 2> /dev/null)
	$(shell GO111MODULE=off go install github.com/boyter/scc@latest)
endif

ensure_pnpm_installed:
ifndef $(shell command -v pnpm 2> /dev/null)
	$(shell npm install -g pnpm)
endif

.PHONY: clean_vendor
clean_vendor:
	rm -rf vendor go.sum

vendor:
	if [ ! -f go.mod ]; then go mod init; fi
	go mod vendor

.PHONY: revendor
revendor: clean_vendor vendor # frontend_vendor

## dependency injection

.PHONY: clean_wire
clean_wire:
	rm -f $(THIS)/internal/build/server/wire_gen.go

.PHONY: wire
wire: ensure_wire_installed vendor
	wire gen $(THIS)/internal/build/server

.PHONY: rewire
rewire: ensure_wire_installed clean_wire wire

## Frontend stuff

.PHONY: clean_frontend
clean_frontend:
	@(cd $(FRONTEND_DIR) && rm -rf dist/build/)

.PHONY: frontend_vendor
frontend_vendor:
	@(cd $(FRONTEND_DIR) && $(FRONTEND_TOOL) install)

.PHONY: dev_frontend
dev_frontend: ensure_pnpm_installed clean-frontend
	@(cd $(FRONTEND_DIR) && $(FRONTEND_TOOL) run autobuild)

# frontend-only runs a simple static server that powers the frontend of the application. In this mode, all API calls are
# skipped, and data on the page is faked. This is useful for making changes that don't require running the entire service.
.PHONY: frontend_only
frontend_only: ensure_pnpm_installed clean-frontend
	@(cd $(FRONTEND_DIR) && $(FRONTEND_TOOL) run start:frontend-only)

## formatting

.PHONY: format_frontend
format_frontend:
	(cd $(FRONTEND_DIR) && $(FRONTEND_TOOL) run format)

.PHONY: format_backend
format_backend:
	for file in `) + "`" + `find $(PWD) -name '*.go'` + "`" + `; do $(GO_FORMAT) $$file; done

.PHONY: fmt
fmt: format

.PHONY: format
format: format_backend format_frontend

.PHONY: check_backend_formatting
check_backend_formatting: vendor
	docker build --tag check_formatting --file environments/testing/dockerfiles/formatting.Dockerfile .
	docker run --interactive --tty --rm check_formatting

.PHONY: check-frontend-formatting
check-frontend-formatting:
	(cd $(FRONTEND_DIR) && $(FRONTEND_TOOL) run format:check)

.PHONY: check_formatting
check_formatting: check_backend_formatting check-frontend-formatting

## Testing things

.PHONY: pre_lint
pre_lint:
	@until fieldalignment -fix ./...; do true; done > /dev/null

.PHONY: docker_lint
docker_lint:
	@docker pull openpolicyagent/conftest:v0.21.0
	docker run --interactive --tty --rm --volume $(PWD):$(PWD) --workdir=$(PWD) openpolicyagent/conftest:v0.21.0 test --policy docker_security.rego ` + "`" + `find . -type f -name "*.Dockerfile"` + "`" + fmt.Sprintf(`

.PHONY: lint
lint: pre_lint docker_lint
	@docker pull golangci/golangci-lint:v1.42
	docker run \
		--rm \
		--volume $(PWD):$(PWD) \
		--workdir=$(PWD) \
		golangci/golangci-lint:v1.42 golangci-lint run --config=.golangci.yml ./...

.PHONY: clean_coverage
clean_coverage:
	@rm -f $(COVERAGE_OUT) profile.out;

.PHONY: coverage
coverage: clean_coverage $(ARTIFACTS_DIR)
	@go test -coverprofile=$(COVERAGE_OUT) -shuffle=on -covermode=atomic -race $(TESTABLE_PACKAGE_LIST) > /dev/null
	@go tool cover -func=$(ARTIFACTS_DIR)/coverage.out | grep 'total:' | xargs | awk '{ print "COVERAGE: " $$3 }'

.PHONY: quicktest # basically only running once instead of with -count 5 or whatever
quicktest: $(ARTIFACTS_DIR) vendor clear
	go build $(TOTAL_PACKAGE_LIST)
	go test -cover -shuffle=on -race -failfast $(TESTABLE_PACKAGE_LIST)

.PHONY: frontend_tests
frontend_tests:
	docker-compose --file $(TEST_DOCKER_COMPOSE_FILES_DIR)/frontend_tests.yaml up \
	--build \
	--quiet-pull \
	--force-recreate \
	--remove-orphans \
	--renew-anon-volumes \
	--always-recreate-deps \
	--abort-on-container-exit

## Integration tests

.PHONY: wipe_docker
wipe_docker:
	@docker stop $(shell docker ps -aq) && docker rm $(shell docker ps -aq)

.PHONY: docker_wipe
docker_wipe:
	@docker stop $(shell docker ps -aq) && docker rm $(shell docker ps -aq)

.PHONY: ensure_mysql_is_up
ensure_mysql_is_up:
	@echo "waiting for mysql"
	@until docker exec --interactive --tty mysql mysql --user='dbuser' --password='hunter2' -e 'SELECT 1' %s; do true; done > /dev/null

.PHONY: ensure_postgres_is_up
ensure_postgres_is_up:
	@echo "waiting for postgres"
	@until docker exec --interactive --tty postgres psql 'postgres://dbuser:hunter2@localhost:5432/%s?sslmode=disable' -c 'SELECT 1'; do true; done > /dev/null

.PHONY: ensure_elasticsearch_is_up
ensure_elasticsearch_is_up:
	@echo "waiting for elasticsearch"
	@until docker run --interactive --tty --network=host curlimages/curl:7.79.1 curl --head "http://localhost:9200/*" ; do true; done > /dev/null

.PHONY: wipe_local_postgres
wipe_local_postgres: ensure_postgres_is_up
	@echo "wiping postgres"
	@until docker exec --interactive --tty postgres psql 'postgres://dbuser:hunter2@localhost:5432/%s?sslmode=disable' -c 'DROP SCHEMA public CASCADE; CREATE SCHEMA public;'; do true; done > /dev/null

.PHONY: wipe_local_mysql
wipe_local_mysql: ensure_mysql_is_up
	@echo "wiping mysql"
	docker exec --interactive --tty mysql mysql --user='dbuser' --password='hunter2' -e "DROP DATABASE %s; CREATE DATABASE %s"

.PHONY: wipe_local_elasticsearch
wipe_local_elasticsearch:
	@echo "wiping elasticsearch"
	@docker run --interactive --tty --network=host curlimages/curl:7.79.1 curl -X DELETE "http://localhost:9200/*"

.PHONY: deploy_base_infra
deploy_base_infra:
	docker-compose \
	--file $(ENVIRONMENTS_DIR)/local/docker-compose-base.yaml up \
	--quiet-pull \
	--no-recreate \
	--always-recreate-deps \
	--detach

.PHONY: lintegration_tests # this is just a handy lil' helper I use sometimes
lintegration_tests: lint clear integration_tests

.PHONY: integration-tests
integration-tests: integration-tests-postgres integration-tests-mysql

.PHONY: integration-tests-
integration-tests-%%:
	docker-compose \
	--file $(TEST_DOCKER_COMPOSE_FILES_DIR)/integration_tests/integration-tests-base.yaml \
	--file $(TEST_DOCKER_COMPOSE_FILES_DIR)/integration_tests/integration-tests-$*.yaml up \
	--build \
	--force-recreate \
	--remove-orphans \
	--renew-anon-volumes \
	--always-recreate-deps $(if $(filter y Y yes YES true TRUE plz sure yup YUP,$(LET_HANG)),, --abort-on-container-exit)

## Running

.PHONY: dev
dev: $(ARTIFACTS_DIR) $(SEARCH_INDICES_DIR) deploy_base_infra
	docker-compose --file $(ENVIRONMENTS_DIR)/local/docker-compose-services.yaml up \
	--quiet-pull \
	--build \
	--force-recreate \
	--renew-anon-volumes \
	--detach

.PHONY: dev_user
dev_user:
	go run $(THIS)/cmd/tools/data_scaffolder --url=http://localhost --count=1 --single-user-mode --debug

.PHONY: load_data_for_admin
load_data_for_admin:
	go run $(THIS)/cmd/tools/data_scaffolder --url=http://localhost --count=5 --debug

## misc

.PHONY: tree
tree:
	tree -d -I vendor

.PHONY: cloc
cloc: ensure_scc_installed
	@scc --include-ext go --exclude-dir vendor
`,
		proj.Name.SingularPackageName(),
		proj.Name.SingularPackageName(),
		proj.Name.SingularPackageName(),
		proj.Name.SingularPackageName(),
		proj.Name.SingularPackageName(),
	)
}
