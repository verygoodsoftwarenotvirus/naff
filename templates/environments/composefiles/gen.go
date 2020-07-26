package composefiles

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func GetDatabasePalabra(vendor string) wordsmith.SuperPalabra {
	switch vendor {
	case "postgres":
		return wordsmith.FromSingularPascalCase("Postgres")
	case "sqlite":
		return wordsmith.FromSingularPascalCase("Sqlite")
	case "mariadb", "maria_db":
		return &wordsmith.ManualWord{
			SingularStr:                           "MariaDB",
			PluralStr:                             "MariaDBs",
			RouteNameStr:                          "mariadb",
			KebabNameStr:                          "mariadb",
			PluralRouteNameStr:                    "mariadbs",
			UnexportedVarNameStr:                  "mariaDB",
			PluralUnexportedVarNameStr:            "mariaDBs",
			PackageNameStr:                        "mariadbs",
			SingularPackageNameStr:                "mariadb",
			SingularCommonNameStr:                 "maria DB",
			ProperSingularCommonNameWithPrefixStr: "a Maria DB",
			PluralCommonNameStr:                   "maria DBs",
			SingularCommonNameWithPrefixStr:       "maria DB",
			PluralCommonNameWithPrefixStr:         "maria DBs",
		}
	default:
		panic(fmt.Sprintf("unknown vendor: %q", vendor))
	}
}

// RenderPackage renders the package
func RenderPackage(project *models.Project) error {
	files := map[string]string{
		"environments/local/docker-compose.yaml":                                         developmentDotYaml(project.Name),
		"environments/testing/compose_files/frontend-tests.yaml":                         frontendTestsDotYAML(project.Name),
		"environments/testing/compose_files/integration_tests/integration-coverage.yaml": integrationCoverageDotYAML(project.Name),
	}

	for _, db := range project.EnabledDatabases() {
		_ = db
		files[fmt.Sprintf("environments/testing/compose_files/integration_tests/integration-tests-%s.yaml", db)] = integrationTestsDotYAML(project.Name, GetDatabasePalabra(db))
		files[fmt.Sprintf("environments/testing/compose_files/load_tests/load-tests-%s.yaml", db)] = loadTestsDotYAML(project.Name, GetDatabasePalabra(db))
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

		if _, err := f.WriteString(file); err != nil {
			log.Printf("error writing to file: %v", err)
			return err
		}
	}

	return nil
}

func developmentDotYaml(projectName wordsmith.SuperPalabra) string {
	return fmt.Sprintf(`version: "3.3"
services:
    database:
        image: postgres:latest
        environment:
            POSTGRES_DB: '%s'
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
    %s-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
            JAEGER_AGENT_HOST: 'tracing-server'
            JAEGER_AGENT_PORT: '6831'
            JAEGER_SAMPLER_MANAGER_HOST_PORT: 'tracing-server:5778'
            JAEGER_SERVICE_NAME: '%s-server'
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
`, projectName.RouteName(), projectName.KebabName(), projectName.KebabName())
}

func integrationTestsDotYAML(projectName, dbName wordsmith.SuperPalabra) string {
	switch dbName.RouteName() {
	case "postgres":
		return fmt.Sprintf(`version: "3.3"
services:
    database:
        image: postgres:latest
        environment:
            POSTGRES_DB: '%s'
            POSTGRES_PASSWORD: 'hunter2'
            POSTGRES_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 2345:5432
    %s-server:
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
            TARGET_ADDRESS: 'http://%s-server:8888'
        links:
            - %s-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-tests.Dockerfile'
        container_name: 'postgres_integration_tests'
`, projectName.RouteName(), projectName.KebabName(), projectName.KebabName(), projectName.KebabName())
	case "mariadb", "maria_db":
		return fmt.Sprintf(`version: "3.3"
services:
    database:
        image: "mariadb:latest"
        environment:
            MYSQL_ALLOW_EMPTY_PASSWORD: 'no'
            MYSQL_DATABASE: '%s'
            MYSQL_PASSWORD: 'hunter2'
            MYSQL_RANDOM_ROOT_PASSWORD: 'yes'
            MYSQL_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 3306:3306
    %s-server:
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
            TARGET_ADDRESS: 'http://%s-server:8888'
        links:
            - %s-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-tests.Dockerfile'
        container_name: 'mariadb_integration_tests'
`, projectName.RouteName(), projectName.KebabName(), projectName.KebabName(), projectName.KebabName())
	case "sqlite":
		return fmt.Sprintf(`version: '3.3'
services:
    %s-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-sqlite.Dockerfile'
    test:
        environment:
            TARGET_ADDRESS: 'http://%s-server:8888'
        links:
            - %s-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-tests.Dockerfile'
        container_name: 'sqlite_integration_tests'
`, projectName.KebabName(), projectName.KebabName(), projectName.KebabName())
	}

	panic("invalid db")
}

func loadTestsDotYAML(projectName, dbName wordsmith.SuperPalabra) string {
	switch dbName.RouteName() {
	case "postgres":
		return fmt.Sprintf(`---
version: '3.3'
services:
    database:
        image: postgres:latest
        environment:
            POSTGRES_DB: '%s'
            POSTGRES_PASSWORD: 'hunter2'
            POSTGRES_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 2345:5432
    load-tests:
        environment:
            TARGET_ADDRESS: 'http://%s-server:8888'
        links:
            - %s-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/load-tests.Dockerfile'
    %s-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        links:
            - database
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-postgres.Dockerfile'
`, projectName.RouteName(), projectName.KebabName(), projectName.KebabName(), projectName.KebabName())
	case "mariadb", "maria_db":
		return fmt.Sprintf(`version: '3.3'
services:
    database:
        image: mariadb:latest
        environment:
            MYSQL_ALLOW_EMPTY_PASSWORD: 'no'
            MYSQL_DATABASE: '%s'
            MYSQL_PASSWORD: 'hunter2'
            MYSQL_RANDOM_ROOT_PASSWORD: 'yes'
            MYSQL_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 3306:3306
    load-tests:
        environment:
            TARGET_ADDRESS: 'http://%s-server:8888'
        links:
            - %s-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/load-tests.Dockerfile'
    %s-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        links:
            - database
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-mariadb.Dockerfile'
`, projectName.RouteName(), projectName.KebabName(), projectName.KebabName(), projectName.KebabName())
	case "sqlite":
		return fmt.Sprintf(`version: '3.3'
services:
    load-tests:
        environment:
            TARGET_ADDRESS: 'http://%s-server:8888'
        links:
            - %s-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/load-tests.Dockerfile'
    %s-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server-sqlite.Dockerfile'
`, projectName.KebabName(), projectName.KebabName(), projectName.KebabName())
	}

	panic("invalid db")
}

func frontendTestsDotYAML(projectName wordsmith.SuperPalabra) string {
	return fmt.Sprintf(`version: "3.3"
services:
    chrome:
        image: selenium/node-chrome:3.141.59-oxygen
        environment:
            HUB_HOST: 'selenium-hub'
            HUB_PORT: '4444'
        logging:
            driver: none
        links:
            - selenium-hub
        volumes:
            - source: '/dev/shm'
              target: '/dev/shm'
              type: 'bind'
    database:
        image: postgres:latest
        environment:
            POSTGRES_DB: '%s'
            POSTGRES_PASSWORD: 'hunter2'
            POSTGRES_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - 2345:5432
    firefox:
        image: selenium/node-firefox:3.141.59-oxygen
        environment:
            HUB_HOST: 'selenium-hub'
            HUB_PORT: '4444'
        logging:
            driver: none
        links:
            - selenium-hub
        volumes:
            - source: '/dev/shm'
              target: '/dev/shm'
              type: 'bind'
    selenium-hub:
        image: selenium/hub:3.141.59-oxygen
        logging:
            driver: none
        ports:
            - 4444:4444
    test:
        environment:
            DOCKER: 'true'
            TARGET_ADDRESS: 'http://%s-server:8888'
        links:
            - selenium-hub
            - %s-server
        build:
            context: '../../../'
            dockerfile: 'environments/testing/dockerfiles/frontend-tests.Dockerfile'
        depends_on:
            - firefox
            - chrome
    %s-server:
        build:
            context: '../../../'
            dockerfile: 'environments/testing/dockerfiles/frontend-tests-server.Dockerfile'
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
        links:
            - database
`, projectName.RouteName(), projectName.KebabName(), projectName.KebabName(), projectName.KebabName())
}

func integrationCoverageDotYAML(projectName wordsmith.SuperPalabra) string {
	return fmt.Sprintf(`version: '3.3'
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
            POSTGRES_DB: '%s'
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
`, projectName.RouteName())
}
