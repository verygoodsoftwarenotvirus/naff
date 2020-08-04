package misc

import (
	"os"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"

	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		project := testprojects.BuildTodoApp()
		project.OutputPath = os.TempDir()

		assert.NoError(t, RenderPackage(project))
	})

	//	T.Run("with invalid output directory", func(t *testing.T) {
	//		t.Parallel()
	//
	//		proj := testprojects.BuildTodoApp()
	//		proj.OutputPath = `/\0/\0/\0`
	//
	//		assert.Error(t, RenderPackage(proj))
	//	})
}

func Test_formattingDotDockerfile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"

		expected := `FROM golang:stretch

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD ../../../dockerfiles .

CMD if [ $(gofmt -l . | grep -Ev '^vendor\/' | head -c1 | wc -c) -ne 0 ]; then exit 1; fi
`
		actual := formattingDotDockerfile(projRoot, "")

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_developmentDotDockerfile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"
		binaryName := "binary"

		expected := `# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

COPY . .
COPY --from=frontend-build-stage /app/public /frontend

RUN go build -trimpath -o /binary gitlab.com/verygoodsoftwarenotvirus/example/cmd/server/v1

# final stage
FROM debian:stretch

COPY --from=build-stage /binary /binary

RUN mkdir /home/appuser
RUN groupadd --gid 999 appuser && \
    useradd --system --uid 999 --gid appuser appuser
RUN chown appuser /home/appuser
WORKDIR /home/appuser
USER appuser

ENV DOCKER=true

ENTRYPOINT ["/binary"]`
		actual := developmentDotDockerfile(projRoot, binaryName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_frontendTestDotDockerfile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"

		expected := `FROM golang:stretch

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

ENTRYPOINT [ "go", "test", "-v", "-failfast", "-parallel=1", "gitlab.com/verygoodsoftwarenotvirus/example/tests/v1/frontend" ]
`
		actual := frontendTestDotDockerfile(projRoot, "")

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_integrationCoverageServerDotDockerfile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"

		expected := `# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go test -o /integration-server -c -coverpkg \
	gitlab.com/verygoodsoftwarenotvirus/example/internal/..., \
	gitlab.com/verygoodsoftwarenotvirus/example/database/v1/..., \
	gitlab.com/verygoodsoftwarenotvirus/example/services/v1/..., \
	gitlab.com/verygoodsoftwarenotvirus/example/cmd/server/v1/ \
    gitlab.com/verygoodsoftwarenotvirus/example/cmd/server/v1

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

`
		actual := integrationCoverageServerDotDockerfile(projRoot, "")

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIntegrationServerDotDockerfile(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"
		binaryName := "binary"

		expected := `# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -trimpath -o /binary -v gitlab.com/verygoodsoftwarenotvirus/example/cmd/server/v1

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

COPY environments/testing/config_files/integration-tests-postgres.toml /etc/config.toml
COPY --from=build-stage /binary /binary
COPY --from=frontend-build-stage /app/public /frontend

ENTRYPOINT ["/binary"]
`
		actual := buildIntegrationServerDotDockerfile("postgres")(projRoot, binaryName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"
		binaryName := "binary"

		expected := `# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -trimpath -o /binary -v gitlab.com/verygoodsoftwarenotvirus/example/cmd/server/v1

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

COPY environments/testing/config_files/integration-tests-sqlite.toml /etc/config.toml
COPY --from=build-stage /binary /binary
COPY --from=frontend-build-stage /app/public /frontend

ENTRYPOINT ["/binary"]
`
		actual := buildIntegrationServerDotDockerfile("sqlite")(projRoot, binaryName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"
		binaryName := "binary"

		expected := `# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -trimpath -o /binary -v gitlab.com/verygoodsoftwarenotvirus/example/cmd/server/v1

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

COPY environments/testing/config_files/integration-tests-mariadb.toml /etc/config.toml
COPY --from=build-stage /binary /binary
COPY --from=frontend-build-stage /app/public /frontend

ENTRYPOINT ["/binary"]
`
		actual := buildIntegrationServerDotDockerfile("mariadb")(projRoot, binaryName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_frontendTestsServerDotDockerfile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"
		binaryName := "binary"

		expected := `# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -trimpath -o /binary -v gitlab.com/verygoodsoftwarenotvirus/example/cmd/server/v1

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
COPY --from=build-stage /binary /binary
COPY --from=frontend-build-stage /app/public /frontend

ENTRYPOINT ["/binary"]
`
		actual := frontendTestsServerDotDockerfile(projRoot, binaryName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_integrationTestsDotDockerfile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"

		expected := `FROM golang:stretch

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

ADD . .

ENTRYPOINT [ "go", "test", "-v", "-failfast", "gitlab.com/verygoodsoftwarenotvirus/example/tests/v1/integration" ]

# for a more specific test:
# ENTRYPOINT [ "go", "test", "-v", "gitlab.com/verygoodsoftwarenotvirus/example/tests/v1/integration", "-run", "InsertTestNameHere" ]
`
		actual := integrationTestsDotDockerfile(projRoot, "")

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_loadTestsDotDockerfile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		projRoot := "gitlab.com/verygoodsoftwarenotvirus/example"

		expected := `# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/verygoodsoftwarenotvirus/example

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go build -o /loadtester gitlab.com/verygoodsoftwarenotvirus/example/tests/v1/load

# final stage
FROM debian:stable

COPY --from=build-stage /loadtester /loadtester

ENV DOCKER=true

ENTRYPOINT ["/loadtester"]
`
		actual := loadTestsDotDockerfile(projRoot, "")

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
