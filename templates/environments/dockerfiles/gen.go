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
		"environments/local/Dockerfile":                                     developmentDotDockerfile,
		"environments/local/workers.Dockerfile":                             developmentWorkersDotDockerfile,
		"environments/testing/dockerfiles/workers.Dockerfile":               testingWorkersDotDockerfile,
		"environments/testing/dockerfiles/formatting.Dockerfile":            formattingDotDockerfile,
		"environments/testing/dockerfiles/frontend-tests-server.Dockerfile": frontendTestsServerDotDockerfile,
		"environments/testing/dockerfiles/integration-tests.Dockerfile":     integrationTestsDotDockerfile,
		"environments/testing/dockerfiles/integration-server.Dockerfile":    buildIntegrationServerDotDockerfile,
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
	return fmt.Sprintf(`FROM golang:1.17-stretch

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY . .

CMD if [ $(gofmt -l . | grep -Ev '^vendor\/' | head -c1 | wc -c) -ne 0 ]; then exit 1; fi
`, projRoot)
}

func developmentWorkersDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`# build stage
FROM golang:1.17-stretch

WORKDIR /go/src/%s

RUN	apt-get update && apt-get install -y \
	--no-install-recommends \
	entr \
	&& rm -rf /var/lib/apt/lists/*
ENV ENTR_INOTIFY_WORKAROUND=true

ENTRYPOINT echo "please wait for workers to start" && sleep 15 && find . -type f \( -iname "*.go*" ! -iname "*_test.go" \) | entr -r go run %s/cmd/workers
`, projRoot, projRoot)
}

func testingWorkersDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`# build stage
FROM golang:1.17-stretch AS build-stage

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY . .

RUN go build -trimpath -o /workers -v %s/cmd/workers

# final stage
FROM debian:bullseye

COPY --from=build-stage /workers /workers

ENTRYPOINT ["/workers"]
`, projRoot, projRoot)
}

func developmentDotDockerfile(projRoot, binaryName string) string {
	return fmt.Sprintf(`# build stage
FROM golang:1.17-stretch

WORKDIR /go/src/%s

RUN	apt-get update && apt-get install -y \
	--no-install-recommends \
	entr \
	&& rm -rf /var/lib/apt/lists/*
ENV ENTR_INOTIFY_WORKAROUND=true

ENTRYPOINT echo "please wait for server to start" && find . -type f \( -iname "*.go*" ! -iname "*_test.go" \) | entr -r go run %s/cmd/server
`, projRoot, projRoot)
}

func frontendTestDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`FROM golang:1.17-stretch

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

ENTRYPOINT [ "go", "test", "-v", "-failfast", "-parallel=1", "%s/tests/v1/frontend" ]
`, projRoot, projRoot)
}

func integrationCoverageServerDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`# build stage
FROM golang:1.17-stretch AS build-stage

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

func buildIntegrationServerDotDockerfile(projRoot, binaryName string) string {
	return fmt.Sprintf(`# build stage
FROM golang:1.17-stretch AS build-stage

WORKDIR /go/src/%s

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

COPY . .

RUN go build -trimpath -o /%s -v %s/cmd/server

# final stage
FROM debian:stretch

COPY --from=build-stage /%s /%s

ENTRYPOINT ["/%s"]
`, projRoot, binaryName, projRoot, binaryName, binaryName, binaryName)
}

func frontendTestsServerDotDockerfile(projRoot, binaryName string) string {
	return fmt.Sprintf(`# build stage
FROM golang:1.17-stretch AS build-stage

WORKDIR /go/src/%s

COPY . .

RUN go build -trimpath -o /%s -v %s/cmd/server

# final stage
FROM debian:stretch

COPY --from=build-stage /%s /%s

RUN mkdir /home/appuser
RUN groupadd --gid 999 appuser && \
    useradd --system --uid 999 --gid appuser appuser
RUN chown appuser /home/appuser
WORKDIR /home/appuser
USER appuser

COPY environments/testing/config_files/frontend-tests.toml /etc/config.toml

ENTRYPOINT ["/%s"]
`, projRoot, binaryName, projRoot, binaryName, binaryName, binaryName)
}

func integrationTestsDotDockerfile(projRoot, _ string) string {
	return fmt.Sprintf(`FROM golang:1.17-stretch

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

WORKDIR /go/src/%s

COPY . .

# to debug a specific test:
# ENTRYPOINT [ "go", "test", "-parallel", "1", "-v", "-failfast", "%s/tests/integration", "-run", "TestIntegration/TestSomething" ]

ENTRYPOINT [ "go", "test", "-v", "-failfast", "%s/tests/integration" ]
`, projRoot, projRoot, projRoot)
}
