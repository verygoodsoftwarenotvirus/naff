package misc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(project *models.Project) error {
	files := map[string]func(projRoot, binaryName string) []byte{
		"dockerfiles/development.Dockerfile":                 developmentDotDockerfile,
		"dockerfiles/formatting.Dockerfile":                  formattingDotDockerfile,
		"dockerfiles/frontend-tests.Dockerfile":              frontendTestDotDockerfile,
		"dockerfiles/integration-coverage-server.Dockerfile": integrationCoverageServerDotDockerfile,
		"dockerfiles/integration-server.Dockerfile":          integrationServerDotDockerfile,
		"dockerfiles/integration-tests.Dockerfile":           integrationTestsDotDockerfile,
		"dockerfiles/load-tests.Dockerfile":                  loadTestsDotDockerfile,
		"dockerfiles/server.Dockerfile":                      serverDotDockerfile,
	}

	for filename, file := range files {
		fname := utils.BuildTemplatePath(project.OutputPath, filename)

		if mkdirErr := os.MkdirAll(filepath.Dir(fname), os.ModePerm); mkdirErr != nil {
			log.Printf("error making directory: %v\n", mkdirErr)
		}

		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("error opening file: %v", err)
			return err
		}

		bytes := file(project.OutputPath, project.Name.SingularPackageName())
		if _, err := f.Write(bytes); err != nil {
			log.Printf("error writing to file: %v", err)
			return err
		}
	}

	return nil
}

func formattingDotDockerfile(projRoot, binaryName string) []byte {
	return []byte(fmt.Sprintf(`FROM golang:stretch

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

CMD if [ $(gofmt -l . | grep -Ev '^vendor\/' | head -c1 | wc -c) -ne 0 ]; then exit 1; fi
`, projRoot))
}

func developmentDotDockerfile(projRoot, binaryName string) []byte {
	return []byte(fmt.Sprintf(`# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/%s

COPY . .
COPY --from=frontend-build-stage /app/public /frontend

RUN go build -trimpath -o /%s %s/cmd/server/v1

# final stage
FROM debian:stretch

COPY --from=build-stage /%s /%s
COPY config_files config_files

RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
USER appuser

ENV DOCKER=true

ENTRYPOINT ["/%s"]`, projRoot, binaryName, projRoot, binaryName, binaryName, binaryName))
}

func frontendTestDotDockerfile(projRoot, binaryName string) []byte {
	return []byte(fmt.Sprintf(`FROM golang:stretch

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

ENTRYPOINT [ "go", "test", "-v", "-failfast", "-parallel=1", "%s/tests/v1/frontend" ]
`, projRoot, projRoot))
}

func integrationCoverageServerDotDockerfile(projRoot, binaryName string) []byte {
	return []byte(fmt.Sprintf(`# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go test -o /integration-server -c -coverpkg \
	%s/internal/..., \
	%s/database/v1/..., \
	%s/services/v1/..., \
	%s/cmd/server/v1/ \
    %s/cmd/server/v1

# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# final stage
FROM debian:stable

COPY config_files config_files
COPY --from=build-stage /integration-server /integration-server

EXPOSE 80

ENTRYPOINT ["/integration-server", "-test.coverprofile=/home/integration-coverage.out"]

`, projRoot, projRoot, projRoot, projRoot, projRoot, projRoot))
}

func integrationServerDotDockerfile(projRoot, binaryName string) []byte {
	return []byte(fmt.Sprintf(`# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -trimpath -o /%s -v %s/cmd/server/v1

# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# final stage
FROM debian:stretch

RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
USER appuser

COPY config_files config_files
COPY --from=build-stage /%s /%s
COPY --from=frontend-build-stage /app/public /frontend

ENTRYPOINT ["/%s"]
`, projRoot, binaryName, projRoot, binaryName, binaryName, binaryName))
}

func integrationTestsDotDockerfile(projRoot, binaryName string) []byte {
	return []byte(fmt.Sprintf(`FROM golang:stretch

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

WORKDIR /go/src/%s

ADD . .

ENTRYPOINT [ "go", "test", "-v", "-failfast", "%s/tests/v1/integration" ]

# for a more specific test:
# ENTRYPOINT [ "go", "test", "-v", "%s/tests/v1/integration", "-run", "TestExport/Exporting/should_be_exportable" ]
`, projRoot, projRoot, projRoot))
}

func loadTestsDotDockerfile(projRoot, binaryName string) []byte {
	return []byte(fmt.Sprintf(`# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -o /loadtester %s/tests/v1/load

# final stage
FROM debian:stable

COPY --from=build-stage /loadtester /loadtester

ENV DOCKER=true

ENTRYPOINT ["/loadtester"]
`, projRoot, projRoot))
}

func serverDotDockerfile(projRoot, binaryName string) []byte {
	return []byte(fmt.Sprintf(`# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -trimpath -o /%s %s/cmd/server/v1

# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# final stage
FROM debian:stable

RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
USER appuser

COPY config_files config_files
COPY --from=build-stage /%s /%s
COPY --from=frontend-build-stage /app/public /frontend

ENV DOCKER=true

ENTRYPOINT ["/%s"]
`, projRoot, binaryName, projRoot, binaryName, binaryName, binaryName))
}
