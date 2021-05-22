package dockerfiles

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
	files := map[string]func(projRoot, binaryName string) string{
		"environments/local/Dockerfile":                                           developmentDotDockerfile,
		"environments/testing/dockerfiles/formatting.Dockerfile":                  formattingDotDockerfile,
		"environments/testing/dockerfiles/frontend-tests.Dockerfile":              frontendTestDotDockerfile,
		"environments/testing/dockerfiles/integration-coverage-server.Dockerfile": integrationCoverageServerDotDockerfile,
		"environments/testing/dockerfiles/frontend-tests-server.Dockerfile":       frontendTestsServerDotDockerfile,
		"environments/testing/dockerfiles/integration-tests.Dockerfile":           integrationTestsDotDockerfile,
		"environments/testing/dockerfiles/load-tests.Dockerfile":                  loadTestsDotDockerfile,
	}

	for _, db := range project.EnabledDatabases() {
		files[fmt.Sprintf("environments/testing/dockerfiles/integration-server-%s.Dockerfile", db)] = buildIntegrationServerDotDockerfile(db)
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
		if _, err := f.WriteString(bytes); err != nil {
			log.Printf("error writing to file: %v", err)
			return err
		}
	}

	return nil
}

func formattingDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`FROM golang:stretch

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD ../../../dockerfiles .

CMD if [ $(gofmt -l . | grep -Ev '^vendor\/' | head -c1 | wc -c) -ne 0 ]; then exit 1; fi
`, projRoot)
}

func developmentDotDockerfile(projRoot, binaryName string) string {
	return fmt.Sprintf(`# frontend-build-stage
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

RUN mkdir /home/appuser
RUN groupadd --gid 999 appuser && \
    useradd --system --uid 999 --gid appuser appuser
RUN chown appuser /home/appuser
WORKDIR /home/appuser
USER appuser

ENV DOCKER=true

ENTRYPOINT ["/%s"]`, projRoot, binaryName, projRoot, binaryName, binaryName, binaryName)
}

func frontendTestDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`FROM golang:stretch

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

ENTRYPOINT [ "go", "test", "-v", "-failfast", "-parallel=1", "%s/tests/v1/frontend" ]
`, projRoot, projRoot)
}

func integrationCoverageServerDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`# build stage
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

COPY --from=build-stage /integration-server /integration-server

EXPOSE 80

ENTRYPOINT ["/integration-server", "-test.coverprofile=/home/integration-coverage.out"]

`, projRoot, projRoot, projRoot, projRoot, projRoot, projRoot)
}

func buildIntegrationServerDotDockerfile(dbName string) func(projRoot, binaryName string) string {
	return func(projRoot, binaryName string) string {
		return fmt.Sprintf(`# build stage
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

RUN mkdir /home/appuser
RUN groupadd --gid 999 appuser && \
    useradd --system --uid 999 --gid appuser appuser
RUN chown appuser /home/appuser
WORKDIR /home/appuser
USER appuser

COPY environments/testing/config_files/integration-tests-%s.toml /etc/config.toml
COPY --from=build-stage /%s /%s
COPY --from=frontend-build-stage /app/public /frontend

ENTRYPOINT ["/%s"]
`, projRoot, binaryName, projRoot, dbName, binaryName, binaryName, binaryName)
	}
}

func frontendTestsServerDotDockerfile(projRoot, binaryName string) string {
	return fmt.Sprintf(`# build stage
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

RUN mkdir /home/appuser
RUN groupadd --gid 999 appuser && \
    useradd --system --uid 999 --gid appuser appuser
RUN chown appuser /home/appuser
WORKDIR /home/appuser
USER appuser

COPY environments/testing/config_files/frontend-tests.toml /etc/config.toml
COPY --from=build-stage /%s /%s
COPY --from=frontend-build-stage /app/public /frontend

ENTRYPOINT ["/%s"]
`, projRoot, binaryName, projRoot, binaryName, binaryName, binaryName)
}

func integrationTestsDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`FROM golang:stretch

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

WORKDIR /go/src/%s

ADD . .

ENTRYPOINT [ "go", "test", "-v", "-failfast", "%s/tests/v1/integration" ]

# for a more specific test:
# ENTRYPOINT [ "go", "test", "-v", "%s/tests/v1/integration", "-run", "InsertTestNameHere" ]
`, projRoot, projRoot, projRoot)
}

func loadTestsDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`# build stage
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
`, projRoot, projRoot)
}
