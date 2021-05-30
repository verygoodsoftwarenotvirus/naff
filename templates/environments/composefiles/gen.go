package composefiles

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(project *models.Project) error {
	files := map[string]string{
		"environments/local/docker-compose.yaml":                                           developmentDotYaml(project),
		"environments/testing/compose_files/integration_tests/integration-tests-base.yaml": integrationTestsBaseDotYAML(project.Name),
		"environments/testing/compose_files/load_tests/load-tests-base.yaml":               loadTestsBaseDotYAML(project.Name),
	}

	for _, db := range project.EnabledDatabases() {
		_ = db
		files[fmt.Sprintf("environments/testing/compose_files/integration_tests/integration-tests-%s.yaml", db)] = integrationTestsDotYAML(project.Name, getDatabasePalabra(db))
		files[fmt.Sprintf("environments/testing/compose_files/load_tests/load-tests-%s.yaml", db)] = loadTestsDotYAML(project.Name, getDatabasePalabra(db))
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

func getDatabasePalabra(vendor string) wordsmith.SuperPalabra {
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

func developmentDotYaml(project *models.Project) string {
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
    todo-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
            JAEGER_DISABLED: 'false'
        ports:
            - 8888:8888
        links:
#            - tracing-server
            - database
#            - prometheus
        volumes:
            - source: '../../environments/local/service.config'
              target: '/etc/config.toml'
              type: 'bind'
            - source: '../../'
              target: '/go/src/%s'
              type: 'bind'
        build:
            context: '../../'
            dockerfile: 'environments/local/Dockerfile'
#    tracing-server:
#        image: jaegertracing/all-in-one:1.22.0
#        logging:
#            driver: none
#        ports:
#            - "5775:5775/udp"
#            - "6831:6831/udp"
#            - "6832:6832/udp"
#            - "5778:5778"
#            - "16686:16686"
#            - "14268:14268"
#            - "9411:9411"
#    prometheus:
#        image: quay.io/prometheus/prometheus:v2.25.0
#        command: '--config.file=/etc/prometheus/config.yaml --storage.tsdb.path=/prometheus --log.level=debug'
#        logging:
#          driver: none
#        ports:
#            - 9090:9090
#        volumes:
#            - source: "../../environments/local/prometheus/config.yaml"
#              target: "/etc/prometheus/config.yaml"
#              type: 'bind'
#    grafana:
#      image: grafana/grafana:7.4.3
#      logging:
#        driver: none
#      ports:
#        - 3000:3000
#      links:
#        - prometheus
#      volumes:
#        - source: '../../environments/local/grafana/grafana.ini'
#          target: '/etc/grafana/grafana.ini'
#          type: 'bind'
#        - source: '../../environments/local/grafana/datasources.yaml'
#          target: '/etc/grafana/provisioning/datasources/datasources.yml'
#          type: 'bind'
#        - source: '../../environments/local/grafana/dashboards.yaml'
#          target: '/etc/grafana/provisioning/dashboards/dashboards.yml'
#          type: 'bind'
#        - source: '../../environments/local/grafana/dashboards'
#          target: '/etc/grafana/provisioning/dashboards/dashboards'
#          type: 'bind'
`, project.Name.RouteName(), project.OutputPath)
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
            %s_SERVICE_LOCAL_SECRET_STORE_KEY: 'SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU='
            USE_NOOP_LOGGER: 'nope'
        links:
            - database
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-postgres.config'
              target: '/etc/config.toml'
              type: 'bind'
    test:
        container_name: 'postgres_integration_tests'
`, projectName.RouteName(), projectName.KebabName(), strings.ToUpper(projectName.RouteName()))
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
        links:
            - database
        environment:
            %s_SERVICE_LOCAL_SECRET_STORE_KEY: 'SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU='
            USE_NOOP_LOGGER: 'nope'
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-mariadb.config'
              target: '/etc/config.toml'
              type: 'bind'
    test:
        container_name: 'mariadb_integration_tests'
`, projectName.RouteName(), projectName.KebabName(), strings.ToUpper(projectName.RouteName()))
	case "sqlite":
		return fmt.Sprintf(`version: '3.3'
services:
    %s-server:
        environment:
            %s_SERVICE_LOCAL_SECRET_STORE_KEY: 'SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU='
            USE_NOOP_LOGGER: 'nope'
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-sqlite.config'
              target: '/etc/config.toml'
              type: 'bind'
    test:
        container_name: 'sqlite_integration_tests'
`, projectName.KebabName(), strings.ToUpper(projectName.RouteName()))
	}

	panic("invalid db")
}

func loadTestsBaseDotYAML(projectName wordsmith.SuperPalabra) string {
	return fmt.Sprintf(`---
version: '3.3'
services:
    load-tests:
        links:
            - %s-server
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/load-tests.Dockerfile'
    %s-server:
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server.Dockerfile'
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
        ports:
            - 80:8888
#    tracing-server:
#        image: jaegertracing/all-in-one:1.22.0
#        logging:
#            driver: none
#        ports:
#            - 6831:6831/udp
#            - 5778:5778
#            - 16686:16686
#    prometheus:
#        image: quay.io/prometheus/prometheus:v2.0.0
#        logging:
#            driver: none
#        ports:
#            - 9090:9090
#        volumes:
#            - source: "../../../../environments/testing/prometheus/config.yaml"
#              target: "/etc/prometheus/config.yaml"
#              type: 'bind'
#        command: '--config.file=/etc/prometheus/config.yaml --storage.tsdb.path=/prometheus'
#    grafana:
#        image: grafana/grafana
#        logging:
#            driver: none
#        ports:
#            - 3000:3000
#        links:
#            - prometheus
#        volumes:
#            - source: '../../../../environments/testing/grafana/grafana.ini'
#              target: '/etc/grafana/grafana.ini'
#              type: 'bind'
#            - source: '../../../../environments/testing/grafana/datasources.yaml'
#              target: '/etc/grafana/provisioning/datasources/datasources.yml'
#              type: 'bind'
#            - source: '../../../../environments/testing/grafana/dashboards.yaml'
#              target: '/etc/grafana/provisioning/dashboards/dashboards.yml'
#              type: 'bind'
#            - source: '../../../../environments/testing/grafana/dashboards'
#              target: '/etc/grafana/provisioning/dashboards/dashboards'
#              type: 'bind'
`, projectName.KebabName(), projectName.KebabName())
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
        links:
            - database
            - tracing-server
            - prometheus
            - loki
            - promtail
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-postgres.toml'
              target: '/etc/config.toml'
              type: 'bind'
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
    %s-server:
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server.Dockerfile'
        links:
            - database
            - tracing-server
            - prometheus
            - loki
            - promtail
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-mariadb.toml'
              target: '/etc/config.toml'
              type: 'bind'
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
            dockerfile: 'environments/testing/dockerfiles/integration-server.Dockerfile'
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-sqlite.toml'
              target: '/etc/config.toml'
              type: 'bind'
`, projectName.KebabName(), projectName.KebabName(), projectName.KebabName())
	}

	panic("invalid db")
}

func integrationTestsBaseDotYAML(projectName wordsmith.SuperPalabra) string {
	return fmt.Sprintf(`version: "3.3"
services:
    %s-server:
        environment:
            CONFIGURATION_FILEPATH: '/etc/config.toml'
            JAEGER_DISABLED: 'true'
        ports:
            - 80:8888
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server.Dockerfile'
    test:
        environment:
            TARGET_ADDRESS: 'http://%s-server:8888'
        links:
            - %s-server
#            - tracing-server
#            - prometheus
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-tests.Dockerfile'
        container_name: 'integration_tests'
#    tracing-server:
#        image: jaegertracing/all-in-one:1.22.0
#        logging:
#            driver: none
#        ports:
#            - "5775:5775/udp"
#            - "6831:6831/udp"
#            - "6832:6832/udp"
#            - "5778:5778"
#            - "16686:16686"
#            - "14268:14268"
#            - "9411:9411"
#    prometheus:
#        image: quay.io/prometheus/prometheus:v2.25.0
#        logging:
#            driver: none
#        ports:
#            - 9090:9090
#        volumes:
#            - source: "../../../../environments/testing/prometheus/config.yaml"
#              target: "/etc/prometheus/config.yaml"
#              type: 'bind'
#        command: '--config.file=/etc/prometheus/config.yaml --storage.tsdb.path=/prometheus'
#    grafana:
#        image: grafana/grafana
#        logging:
#            driver: none
#        ports:
#            - 3000:3000
#        links:
#            - prometheus
#        volumes:
#            - source: '../../../../environments/testing/grafana/grafana.ini'
#              target: '/etc/grafana/grafana.ini'
#              type: 'bind'
#            - source: '../../../../environments/testing/grafana/datasources.yaml'
#              target: '/etc/grafana/provisioning/datasources/datasources.yml'
#              type: 'bind'
#            - source: '../../../../environments/testing/grafana/dashboards.yaml'
#              target: '/etc/grafana/provisioning/dashboards/dashboards.yml'
#              type: 'bind'
#            - source: '../../../../environments/testing/grafana/dashboards'
#              target: '/etc/grafana/provisioning/dashboards/dashboards'
#              type: 'bind'
`, projectName.KebabName(), projectName.KebabName(), projectName.KebabName())
}
