GOPATH            := $(GOPATH)
COVERAGE_OUT      := coverage.out

SERVER_PRIV_KEY := dev_files/certs/server/key.pem
SERVER_CERT_KEY := dev_files/certs/server/cert.pem
CLIENT_PRIV_KEY := dev_files/certs/client/key.pem
CLIENT_CERT_KEY := dev_files/certs/client/cert.pem

## generic make stuff

.PHONY: clean
clean:
	rm -f naff_debug

.PHONY: dev_deps
dev_deps:
	go get -u github.com/gobuffalo/packr/v2/packr2

## Project prerequisites
vendor:
	docker run --env GO111MODULE=on --env GOPATH=$(GOPATH) --volume `pwd`:`pwd` --workdir=`pwd` --workdir=`pwd` golang:latest /bin/sh -c "go mod vendor"

.PHONY: revendor
revendor:
	rm -rf vendor go.{mod,sum}
	GO111MODULE=on go mod init
	$(MAKE) vendor

.PHONY: test
test:
	docker build --tag coverage-todo:latest --file dockerfiles/coverage.Dockerfile .
	docker run --rm --volume `pwd`:`pwd` --workdir=`pwd` coverage-todo:latest

.PHONY: run
run:
	go run gitlab.com/verygoodsoftwarenotvirus/naff/cmd/cli generate

naff_debug:
	go build -o naff_debug gitlab.com/verygoodsoftwarenotvirus/naff/cmd/cli

templates:
	@rm -rf template/
	@mkdir template
	go run cmd/tools/template_builder/main.go
