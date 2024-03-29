package composefiles

import (
	"os"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

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

	T.Run("mysql", func(t *testing.T) {
		dbName := "mysql"

		expected := `MySQL`
		actual := getDatabasePalabra(dbName).Singular()

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

	T.Run("mysql", func(t *testing.T) {
		exampleProjectName := wordsmith.FromSingularPascalCase("Whatever")
		dbName := getDatabasePalabra("mysql")

		expected := `version: "3.3"
services:
    database:
        image: "mysql:latest"
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
            dockerfile: 'environments/testing/dockerfiles/integration-server-mysql.Dockerfile'
    test:
        environment:
            TARGET_ADDRESS: 'http://whatever-server:8888'
        links:
            - whatever-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-tests.Dockerfile'
        container_name: 'mysql_integration_tests'
`
		actual := integrationTestsDotYAML(exampleProjectName, dbName)

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
