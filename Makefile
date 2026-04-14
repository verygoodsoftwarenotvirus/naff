DDB_REPO := $(shell realpath ../../dinnerdonebetter/dinnerdonebetter/backend)
DDB_OUTPUT := .ddb_generated

.PHONY: build test generate generate-ddb diff-ddb clean

build:
	go build ./...

test:
	go test ./... -count=1

generate: build
	@rm -rf output/
	@mkdir -p output/
	go run ./cmd/naff/ generate -c testdata/integration_project.yaml -o output/

generate-ddb: build
	@rm -rf $(DDB_OUTPUT)
	@mkdir -p $(DDB_OUTPUT)
	go run ./cmd/naff/ generate -c testdata/ddb.yaml -o $(DDB_OUTPUT)

diff-ddb: generate-ddb
	@if ! command -v meld >/dev/null 2>&1; then \
		echo "meld not found, install it or use another diff tool"; \
		echo "  brew install --cask meld"; \
		echo ""; \
		echo "Generated output is at: $(DDB_OUTPUT)"; \
		echo "DDB backend is at:      $(DDB_REPO)"; \
		exit 1; \
	fi
	meld $(DDB_OUTPUT) $(DDB_REPO)

clean:	
	rm -rf output/ $(DDB_OUTPUT)
