.PHONY: build test generate clean

build:
	go build ./...

test:
	go test ./... -count=1

generate: build
	@rm -rf output/
	@mkdir -p output/
	go run ./cmd/naff/ generate -c testdata/integration_project.yaml -o output/

clean:
	rm -rf output/
