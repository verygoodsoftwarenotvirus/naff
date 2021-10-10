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
		"environments/local/docker-compose-base.yaml":                                      localComposeBase(project),
		"environments/local/docker-compose-services.yaml":                                  localComposeServices(project),
		"environments/testing/compose_files/integration_tests/integration-tests-base.yaml": integrationTestsBaseDotYAML(project.Name),
	}

	for _, db := range project.EnabledDatabases() {
		x := getDatabasePalabra(db)
		files[fmt.Sprintf("environments/testing/compose_files/integration_tests/integration-tests-%s.yaml", x.RouteName())] = integrationTestsDotYAML(project.Name, x)
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
	case string(models.Postgres):
		return wordsmith.FromSingularPascalCase("Postgres")
	case string(models.MySQL):
		return &wordsmith.ManualWord{
			SingularStr:                           "MySQL",
			PluralStr:                             "MySQLs",
			RouteNameStr:                          "mysql",
			KebabNameStr:                          "mysql",
			PluralRouteNameStr:                    "mysqls",
			UnexportedVarNameStr:                  "mariaDB",
			PluralUnexportedVarNameStr:            "mariaDBs",
			PackageNameStr:                        "mysqls",
			SingularPackageNameStr:                "mysql",
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

func localComposeBase(project *models.Project) string {
	return fmt.Sprintf(`version: "3.8"
services:
    worker_queue:
        image: redis:6-buster
        container_name: worker_queue
    postgres:
        hostname: pgdatabase
        container_name: database
        image: postgres:13
        environment:
            POSTGRES_DB: '%s'
            POSTGRES_PASSWORD: 'hunter2'
            POSTGRES_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - '2345:5432'
    elasticsearch:
      image: elasticsearch:7.14.1
      ports:
        - '9200:9200'
        - '9300:9300'
      environment:
        discovery.type: 'single-node'
    tracing-server:
        image: jaegertracing/all-in-one:1.22.0
        logging:
            driver: none
        ports:
            - "5775:5775/udp"
            - "6831:6831/udp"
            - "6832:6832/udp"
            - "5778:5778"
            - "16686:16686"
            - "14268:14268"
            - "9411:9411"
    prometheus:
        image: quay.io/prometheus/prometheus:v2.25.0
        command: '--config.file=/etc/prometheus/config.yaml --storage.tsdb.path=/prometheus --log.level=debug'
        logging:
          driver: none
        ports:
            - '9090:9090'
        volumes:
            - source: "../../environments/local/prometheus/config.yaml"
              target: "/etc/prometheus/config.yaml"
              type: 'bind'
    grafana:
      image: grafana/grafana:7.4.3
      logging:
        driver: none
      ports:
        - '3000:3000'
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
`, project.Name.RouteName())
}

func localComposeServices(project *models.Project) string {
	return fmt.Sprintf(`version: "3.8"
services:
    workers:
        container_name: workers
        volumes:
            - source: '../../environments/testing/config_files/integration-tests-postgres.config'
              target: '/etc/service.config'
              type: 'bind'
            - source: '../../'
              target: '/go/src/%s/cmd/workers'
              type: 'bind'
        build:
            context: '../../'
            dockerfile: 'environments/local/workers.Dockerfile'
    %s-server:
        container_name: api_server
        environment:
            %s_SERVER_LOCAL_CONFIG_STORE_KEY: 'SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU='
            CONFIGURATION_FILEPATH: '/etc/service.config'
            JAEGER_DISABLED: 'false'
        ports:
            - '8888:8888'
        volumes:
            - source: '../../environments/local/service.config'
              target: '/etc/service.config'
              type: 'bind'
            - source: '../../'
              target: '/go/src/%s'
              type: 'bind'
        build:
            context: '../../'
            dockerfile: 'environments/local/Dockerfile'
`, project.OutputPath, project.Name.KebabName(), strings.ToUpper(project.Name.Singular()), project.OutputPath)
}

func integrationTestsDotYAML(projectName, dbName wordsmith.SuperPalabra) string {
	switch dbName.Singular() {
	case string(models.Postgres):
		return fmt.Sprintf(`version: "3.8"
services:
    workers:
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/workers.Dockerfile'
        environment:
            CONFIGURATION_FILEPATH: '/etc/service.config'
            %s_WORKERS_LOCAL_CONFIG_STORE_KEY: 'SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU='
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-postgres.config'
              target: '/etc/service.config'
              type: 'bind'
    api_server:
        depends_on:
            - workers
        environment:
            USE_NOOP_LOGGER: 'nope'
            %s_SERVER_LOCAL_CONFIG_STORE_KEY: 'SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU='
            CONFIGURATION_FILEPATH: '/etc/service.config'
        ports:
            - '8888:8888'
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server.Dockerfile'
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-postgres.config'
              target: '/etc/service.config'
              type: 'bind'
`, strings.ToUpper(projectName.Singular()), strings.ToUpper(projectName.Singular()))
	case string(models.MySQL):
		return fmt.Sprintf(`version: "3.8"
services:
    workers:
        container_name: workers
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/workers.Dockerfile'
        environment:
            CONFIGURATION_FILEPATH: '/etc/service.config'
            %s_WORKERS_LOCAL_CONFIG_STORE_KEY: 'SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU='
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-mysql.config'
              target: '/etc/service.config'
              type: 'bind'
    api_server:
        depends_on:
            - workers
        environment:
            %s_SERVER_LOCAL_CONFIG_STORE_KEY: 'SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU='
            CONFIGURATION_FILEPATH: '/etc/service.config'
            JAEGER_DISABLED: 'false'
        ports:
            - '8888:8888'
        build:
            context: '../../../../'
            dockerfile: 'environments/testing/dockerfiles/integration-server.Dockerfile'
        volumes:
            - source: '../../../../environments/testing/config_files/integration-tests-mysql.config'
              target: '/etc/service.config'
              type: 'bind'
`, strings.ToUpper(projectName.Singular()), strings.ToUpper(projectName.Singular()))
	}

	panic(fmt.Sprintf("invalid db: ", dbName.RouteName()))
}

func integrationTestsBaseDotYAML(projectName wordsmith.SuperPalabra) string {
	return fmt.Sprintf(`version: "3.8"
services:
    redis:
        hostname: worker_queue
        image: redis:6-buster
        container_name: redis
        ports:
            - '6379:6379'
    postgres:
        container_name: postgres
        hostname: pgdatabase
        image: postgres:13
        environment:
            POSTGRES_DB: '%s'
            POSTGRES_PASSWORD: 'hunter2'
            POSTGRES_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - '5432:5432'
    mysql:
        container_name: mysql
        hostname: mysqldatabase
        image: "mysql:8"
        environment:
            MYSQL_ALLOW_EMPTY_PASSWORD: 'no'
            MYSQL_DATABASE: '%s'
            MYSQL_PASSWORD: 'hunter2'
            MYSQL_RANDOM_ROOT_PASSWORD: 'yes'
            MYSQL_USER: 'dbuser'
        logging:
            driver: none
        ports:
            - '3306:3306'
    elasticsearch:
        image: elasticsearch:7.14.1
        ports:
            - '9200:9200'
            - '9300:9300'
        environment:
            discovery.type: 'single-node'
#    prometheus:
#        image: quay.io/prometheus/prometheus:v2.25.0
#        logging:
#            driver: none
#        ports:
#            - '9090:9090'
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
#            - '3000:3000'
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
`, projectName.RouteName(), projectName.RouteName())
}
