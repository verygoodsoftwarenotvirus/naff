package composefiles

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"os"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"

	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		project := testprojects.BuildTodoApp()
		project.OutputPath = os.TempDir()

		assert.NoError(t, RenderPackage(project))
	})

	T.Run("with invalid output directory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.OutputPath = `/dev/null`

		assert.Error(t, RenderPackage(proj))
	})
}

func Test_getDatabasePalabra(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbName := "postgres"

		expected := `Postgres`
		actual := getDatabasePalabra(dbName).Singular()

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbName := "sqlite"

		expected := `Sqlite`
		actual := getDatabasePalabra(dbName).Singular()

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbName := "mariadb"

		expected := `MariaDB`
		actual := getDatabasePalabra(dbName).Singular()

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_developmentDotYaml(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleProjectName := wordsmith.FromSingularPascalCase("Whatever")

		expected := `version: "3.3"
services:
    database:
        image: postgres:latest
        environment:
            POSTGRES_DB: 'whatever'
            POSTGRES_PASSWORD: 'hunter2'
            POSTGRES_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 2345:5432
    grafana:
        image: grafana/grafana
        logging:
            driver: none
        ports:
            - 3000:3000
        links:
            - prometheus
        volumes:
            - source: '../../environments/local/grafana/grafana.ini'
              target: '/etc/grafana/grafana.ini'
              type: 'bind'
            - source: '../../environments/local/grafana/datasources.yaml'
              target: '/etc/grafana/provisioning/datasources/datasources.yml'
              type: 'bind'
            - source: '../../environments/local/grafana/dashboards.yaml'
              target: '/etc/grafana/provisioning/dashboards/dashboards.yml'
              type: 'bind'
            - source: '../../environments/local/grafana/dashboards'
              target: '/etc/grafana/provisioning/dashboards/dashboards'
              type: 'bind'
    prometheus:
        image: quay.io/prometheus/prometheus:v2.0.0
        logging:
            driver: none
        ports:
            - 9090:9090
        volumes:
            - source: "../../environments/local/prometheus/config.yaml"
              target: "/etc/prometheus/config.yaml"
              type: 'bind'
        command: '--config.file=/etc/prometheus/config.yaml --storage.tsdb.path=/prometheus'
    whatever-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
            JAEGER_AGENT_HOST: 'tracing-server'
            JAEGER_AGENT_PORT: '6831'
            JAEGER_SAMPLER_MANAGER_HOST_PORT: 'tracing-server:5778'
            JAEGER_SERVICE_NAME: 'whatever-server'
        ports:
            - 80:8888
        links:
            - tracing-server
            - database
        volumes:
            - source: '../../frontend/v1/public'
              target: '/frontend'
              type: 'bind'
            - source: '../../environments/local/config.toml'
              target: '/etc/config.toml'
              type: 'bind'
        build:
            context: '../../'
            dockerfile: 'environments/local/Dockerfile'
        depends_on:
            - prometheus
            - grafana
    tracing-server:
        image: jaegertracing/all-in-one:latest
        logging:
            driver: none
        ports:
            - 6831:6831/udp
            - 5778:5778
            - 16686:16686
`
		actual := developmentDotYaml(exampleProjectName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_integrationTestsDotYAML(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		exampleProjectName := wordsmith.FromSingularPascalCase("Whatever")
		dbName := getDatabasePalabra("postgres")

		expected := `version: "3.3"
services:
    database:
        image: postgres:latest
        environment:
            POSTGRES_DB: 'whatever'
            POSTGRES_PASSWORD: 'hunter2'
            POSTGRES_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 2345:5432
    whatever-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        links:
            - database
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-postgres.Dockerfile'
    test:
        environment:
            TARGET_ADDRESS: 'http://whatever-server:8888'
        links:
            - whatever-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-tests.Dockerfile'
        container_name: 'postgres_integration_tests'
`
		actual := integrationTestsDotYAML(exampleProjectName, dbName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		exampleProjectName := wordsmith.FromSingularPascalCase("Whatever")
		dbName := getDatabasePalabra("sqlite")

		expected := `version: '3.3'
services:
    whatever-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-sqlite.Dockerfile'
    test:
        environment:
            TARGET_ADDRESS: 'http://whatever-server:8888'
        links:
            - whatever-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-tests.Dockerfile'
        container_name: 'sqlite_integration_tests'
`
		actual := integrationTestsDotYAML(exampleProjectName, dbName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		exampleProjectName := wordsmith.FromSingularPascalCase("Whatever")
		dbName := getDatabasePalabra("mariadb")

		expected := `version: "3.3"
services:
    database:
        image: "mariadb:latest"
        environment:
            MYSQL_ALLOW_EMPTY_PASSWORD: 'no'
            MYSQL_DATABASE: 'whatever'
            MYSQL_PASSWORD: 'hunter2'
            MYSQL_RANDOM_ROOT_PASSWORD: 'yes'
            MYSQL_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 3306:3306
    whatever-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        links:
            - database
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-mariadb.Dockerfile'
    test:
        environment:
            TARGET_ADDRESS: 'http://whatever-server:8888'
        links:
            - whatever-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-tests.Dockerfile'
        container_name: 'mariadb_integration_tests'
`
		actual := integrationTestsDotYAML(exampleProjectName, dbName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_loadTestsDotYAML(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		exampleProjectName := wordsmith.FromSingularPascalCase("Whatever")
		dbName := getDatabasePalabra("postgres")

		expected := `---
version: '3.3'
services:
    database:
        image: postgres:latest
        environment:
            POSTGRES_DB: 'whatever'
            POSTGRES_PASSWORD: 'hunter2'
            POSTGRES_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 2345:5432
    load-tests:
        environment:
            TARGET_ADDRESS: 'http://whatever-server:8888'
        links:
            - whatever-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/load-tests.Dockerfile'
    whatever-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        links:
            - database
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-postgres.Dockerfile'
`
		actual := loadTestsDotYAML(exampleProjectName, dbName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		exampleProjectName := wordsmith.FromSingularPascalCase("Whatever")
		dbName := getDatabasePalabra("sqlite")

		expected := `version: '3.3'
services:
    load-tests:
        environment:
            TARGET_ADDRESS: 'http://whatever-server:8888'
        links:
            - whatever-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/load-tests.Dockerfile'
    whatever-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-sqlite.Dockerfile'
`
		actual := loadTestsDotYAML(exampleProjectName, dbName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		exampleProjectName := wordsmith.FromSingularPascalCase("Whatever")
		dbName := getDatabasePalabra("mariadb")

		expected := `version: '3.3'
services:
    database:
        image: mariadb:latest
        environment:
            MYSQL_ALLOW_EMPTY_PASSWORD: 'no'
            MYSQL_DATABASE: 'whatever'
            MYSQL_PASSWORD: 'hunter2'
            MYSQL_RANDOM_ROOT_PASSWORD: 'yes'
            MYSQL_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 3306:3306
    load-tests:
        environment:
            TARGET_ADDRESS: 'http://whatever-server:8888'
        links:
            - whatever-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/load-tests.Dockerfile'
    whatever-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        links:
            - database
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-mariadb.Dockerfile'
`
		actual := loadTestsDotYAML(exampleProjectName, dbName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_integrationCoverageDotYAML(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleProjectName := wordsmith.FromSingularPascalCase("Whatever")

		expected := `version: '3.3'
services:
    coverage-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
            RUNTIME_DURATION: '30s'
        ports:
            - 80:8888
        links:
            - database
        volumes:
            - source: '../../../../artifacts/'
              target: '/home/'
              type: 'bind'
            - source: '../../../../environments/testing/config_files/coverage.toml'
              target: '/etc/config.toml'
              type: 'bind'
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-coverage-server.Dockerfile'
    database:
        image: postgres:latest
        environment:
            POSTGRES_DB: 'whatever'
            POSTGRES_PASSWORD: 'hunter2'
            POSTGRES_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 2345:5432
    test:
        environment:
            JAEGER_AGENT_HOST: 'tracing-server'
            JAEGER_AGENT_PORT: '6831'
            JAEGER_SAMPLER_MANAGER_HOST_PORT: 'tracing-server:5778'
            JAEGER_SERVICE_NAME: 'coverage-server'
            TARGET_ADDRESS: 'http://coverage-server'
            WAIT_FOR_COVERAGE: 'yes'
        links:
            - coverage-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-tests.Dockerfile'
`
		actual := integrationTestsBaseDotYAML(exampleProjectName)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
