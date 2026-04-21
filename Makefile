DDB_REPO := $(shell realpath ../../dinnerdonebetter/dinnerdonebetter)
ARTIFACTS := artifacts
DDB_OUTPUT := $(ARTIFACTS)/ddb
KITCHEN_SINK_OUTPUT := $(ARTIFACTS)/kitchen_sink

.PHONY: build test generate generate-ddb generate-kitchen-sink diff-ddb clean

build:
	go build ./...

test:
	go test ./... -count=1

# Post-pivot, `generate` emits the verbatim DDB mirror regardless of YAML
# spec contents. Pointed at testdata/ddb.yaml for clarity.
generate: build
	@rm -rf $(ARTIFACTS)/output
	@mkdir -p $(ARTIFACTS)/output
	go run ./cmd/naff/ generate -c testdata/ddb.yaml -o $(ARTIFACTS)/output

generate-ddb: build
	@rm -rf $(DDB_OUTPUT)
	@mkdir -p $(DDB_OUTPUT)
	go run ./cmd/naff/ generate -c testdata/ddb.yaml -o $(DDB_OUTPUT)

generate-kitchen-sink: build
	@rm -rf $(KITCHEN_SINK_OUTPUT)
	@mkdir -p $(KITCHEN_SINK_OUTPUT)
	go run ./cmd/naff/ generate -c testdata/kitchen_sink.yaml -o $(KITCHEN_SINK_OUTPUT)

ADDENDUM := backend/internal/domain/mealplanning

compare-ddb: # generate-ddb
	meld $(DDB_OUTPUT)/$(ADDENDUM) $(DDB_REPO)/$(ADDENDUM) &

clean:
	rm -rf $(ARTIFACTS)/
