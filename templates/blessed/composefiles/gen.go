package composefiles

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	pgImage = "postgres:latest"
)

var (
	nullLogger = &models.DockerComposeLogging{
		Driver: "none",
	}
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
func RenderPackage(pkgRoot string, projectName wordsmith.SuperPalabra, types []models.DataType) error {
	files := map[string]models.DockerComposeFile{
		"compose-files/development.json":          developmentDotJSON(projectName),
		"compose-files/frontend-tests.json":       frontendTestsDotJSON(projectName),
		"compose-files/integration-coverage.json": integrationCoverageDotJSON(projectName),
		"compose-files/production.json":           productionDotJSON(projectName),
	}

	for _, db := range []string{"postgres", "sqlite", "mariadb"} {
		files[fmt.Sprintf("compose-files/integration-tests-%s.json", db)] = integrationTestsDotJSON(projectName, GetDatabasePalabra(db))
		files[fmt.Sprintf("compose-files/load-tests-%s.json", db)] = loadTestsDotJSON(projectName, GetDatabasePalabra(db))
	}

	for filename, file := range files {
		fname := utils.BuildTemplatePath(pkgRoot, filename)

		if mkdirErr := os.MkdirAll(filepath.Dir(fname), os.ModePerm); mkdirErr != nil {
			log.Printf("error making directory: %v\n", mkdirErr)
		}

		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("error opening file: %v", err)
			return err
		}

		bytes, _ := json.MarshalIndent(file, "", "  ")
		if _, err := f.Write(bytes); err != nil {
			log.Printf("error writing to file: %v", err)
			return err
		}
	}

	return nil
}

func developmentDotJSON(projectName wordsmith.SuperPalabra) models.DockerComposeFile {
	serviceName := fmt.Sprintf("%s-server", projectName.KebabName())

	return models.DockerComposeFile{
		Version: "3.3",
		Services: map[string]models.DockerComposeService{
			"database": {
				Image: pgImage,
				Environment: map[string]string{
					"POSTGRES_DB":       "todo",
					"POSTGRES_PASSWORD": "hunter2",
					"POSTGRES_USER":     "dbuser",
				},
				Logging: nullLogger,
				Ports:   []string{"2345:5432"},
			},
			"grafana": {
				Image: "grafana/grafana",
				Links: []string{
					"prometheus",
				},
				Logging: nullLogger,
				Ports: []string{
					"3000:3000",
				},
				Volumes: []models.DockerVolume{
					{
						Source: "../deploy/grafana/local/provisioning",
						Target: "/etc/grafana/provisioning",
						Type:   "bind",
					},
					{
						Source: "../deploy/grafana/dashboards",
						Target: "/etc/grafana/dashboards",
						Type:   "bind",
					},
				},
			},
			"prometheus": {
				Command: "--config.file=/etc/prometheus/config.yaml --storage.tsdb.path=/prometheus",
				Image:   "quay.io/prometheus/prometheus:v2.0.0",
				Logging: nullLogger,
				Ports: []string{
					"9090:9090",
				},
				Volumes: []models.DockerVolume{
					{
						Source: "../deploy/prometheus/local/config.yaml",
						Target: "/etc/prometheus/config.yaml",
						Type:   "bind",
					},
				},
			},
			serviceName: {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/development.Dockerfile",
				},
				DependsOn: []string{
					"prometheus",
					"grafana",
				},
				Environment: map[string]string{
					"CONFIGURATION_FILEPATH":           "config_files/development.toml",
					"DOCKER":                           "true",
					"JAEGER_AGENT_HOST":                "tracing-server",
					"JAEGER_AGENT_PORT":                "6831",
					"JAEGER_SAMPLER_MANAGER_HOST_PORT": "tracing-server:5778",
					"JAEGER_SERVICE_NAME":              fmt.Sprintf("%s-server", projectName.KebabName()),
				},
				Links: []string{
					"tracing-server",
					"database",
				},
				Ports: []string{
					"80:8888",
				},
				Volumes: []models.DockerVolume{
					{
						Source: "../frontend/v1/public",
						Target: "/frontend",
						Type:   "bind",
					},
				},
			},
			"tracing-server": {
				Image:   "jaegertracing/all-in-one:latest",
				Logging: nullLogger,
				Ports: []string{
					"6831:6831/udp",
					"5778:5778",
					"16686:16686",
				},
			},
		},
	}
}

func integrationTestsDotJSON(projectName, dbName wordsmith.SuperPalabra) models.DockerComposeFile {
	serviceName := fmt.Sprintf("%s-server", projectName.KebabName())

	dcf := models.DockerComposeFile{
		Version: "3.3",
		Services: map[string]models.DockerComposeService{
			"test": {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/integration-tests.Dockerfile",
				},
				Environment: map[string]string{
					"DOCKER":         "true",
					"TARGET_ADDRESS": fmt.Sprintf("http://%s:8888", serviceName),
				},
				Links: []string{serviceName},
			},
			serviceName: {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/integration-server.Dockerfile",
				},
				Environment: map[string]string{
					"DOCKER":                 "true",
					"CONFIGURATION_FILEPATH": fmt.Sprintf("config_files/integration-tests-%s.toml", dbName.KebabName()),
				},
				Ports: []string{
					"80:8888",
				},
			},
		},
	}

	switch dbName.RouteName() {
	case "postgres":
		dcf.Services["database"] = models.DockerComposeService{
			Image: pgImage,
			Environment: map[string]string{
				"POSTGRES_DB":       "todo",
				"POSTGRES_PASSWORD": "hunter2",
				"POSTGRES_USER":     "dbuser",
			},
			Ports:   []string{"2345:5432"},
			Logging: nullLogger,
		}
		dcf.Services[serviceName] = models.DockerComposeService{
			Build: &models.DockerComposeBuild{
				Context:    "../",
				Dockerfile: "dockerfiles/integration-server.Dockerfile",
			},
			Environment: map[string]string{
				"DOCKER":                 "true",
				"CONFIGURATION_FILEPATH": fmt.Sprintf("config_files/integration-tests-%s.toml", dbName.KebabName()),
			},
			Links: []string{"database"},
			Ports: []string{
				"80:8888",
			},
		}
	case "mariadb", "maria_db":
		dcf.Services["database"] = models.DockerComposeService{
			Image: "mariadb:latest",
			Environment: map[string]string{
				"MYSQL_ALLOW_EMPTY_PASSWORD": "no",
				"MYSQL_DATABASE":             "todo",
				"MYSQL_PASSWORD":             "hunter2",
				"MYSQL_RANDOM_ROOT_PASSWORD": "yes",
				"MYSQL_USER":                 "dbuser",
			},
			Ports:   []string{"3306:3306"},
			Logging: nullLogger,
		}
		dcf.Services[serviceName] = models.DockerComposeService{
			Build: &models.DockerComposeBuild{
				Context:    "../",
				Dockerfile: "dockerfiles/integration-server.Dockerfile",
			},
			Environment: map[string]string{
				"DOCKER":                 "true",
				"CONFIGURATION_FILEPATH": fmt.Sprintf("config_files/integration-tests-%s.toml", dbName.KebabName()),
			},
			Links: []string{"database"},
			Ports: []string{
				"80:8888",
			},
		}
	}

	return dcf
}

func loadTestsDotJSON(projectName, dbName wordsmith.SuperPalabra) models.DockerComposeFile {
	serviceName := fmt.Sprintf("%s-server", projectName.KebabName())

	dbrn := dbName.RouteName()
	var links []string
	if dbrn != "sqlite" {
		links = []string{
			"database",
		}
	}

	dcf := models.DockerComposeFile{
		Version: "3.3",
		Services: map[string]models.DockerComposeService{
			"load-tests": {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/load-tests.Dockerfile",
				},
				Environment: map[string]string{
					"TARGET_ADDRESS": "http://todo-server:8888",
				},
				Links: []string{
					serviceName,
				},
			},
			serviceName: {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/integration-server.Dockerfile",
				},
				Environment: map[string]string{
					"CONFIGURATION_FILEPATH": fmt.Sprintf("config_files/integration-tests-%s.toml", dbName.KebabName()),
					"DOCKER":                 "true",
				},
				Ports: []string{
					"80:8888",
				},
				Links: links,
			},
		},
	}

	switch dbrn {
	case "postgres":
		dcf.Services["database"] = models.DockerComposeService{
			Image: pgImage,
			Environment: map[string]string{
				"POSTGRES_DB":       "todo",
				"POSTGRES_PASSWORD": "hunter2",
				"POSTGRES_USER":     "dbuser",
			},
			Ports:   []string{"2345:5432"},
			Logging: nullLogger,
		}
	case "mariadb", "maria_db":
		dcf.Services["database"] = models.DockerComposeService{
			Image: "mariadb:latest",
			Environment: map[string]string{
				"MYSQL_ALLOW_EMPTY_PASSWORD": "no",
				"MYSQL_DATABASE":             "todo",
				"MYSQL_PASSWORD":             "hunter2",
				"MYSQL_RANDOM_ROOT_PASSWORD": "yes",
				"MYSQL_USER":                 "dbuser",
			},
			Ports:   []string{"3306:3306"},
			Logging: nullLogger,
		}
	}

	return dcf
}

func frontendTestsDotJSON(projectName wordsmith.SuperPalabra) models.DockerComposeFile {
	hubServiceName := "selenium-hub"
	serviceName := fmt.Sprintf("%s-server", projectName.KebabName())

	return models.DockerComposeFile{
		Version: "3.3",
		Services: map[string]models.DockerComposeService{
			"chrome": {
				Image: "selenium/node-chrome:3.141.59-oxygen",
				Environment: map[string]string{
					"HUB_HOST": "selenium-hub",
					"HUB_PORT": "4444",
				},
				Logging: nullLogger,
				Links:   []string{hubServiceName},
				Volumes: []models.DockerVolume{
					{
						Source: "/dev/shm",
						Target: "/dev/shm",
						Type:   "bind",
					},
				},
			},
			"firefox": {
				Image: "selenium/node-firefox:3.141.59-oxygen",
				Environment: map[string]string{
					"HUB_HOST": "selenium-hub",
					"HUB_PORT": "4444",
				},
				Logging: nullLogger,
				Links:   []string{hubServiceName},
				Volumes: []models.DockerVolume{
					{
						Source: "/dev/shm",
						Target: "/dev/shm",
						Type:   "bind",
					},
				},
			},
			"database": {
				Image: pgImage,
				Environment: map[string]string{
					"POSTGRES_DB":       "todo",
					"POSTGRES_PASSWORD": "hunter2",
					"POSTGRES_USER":     "dbuser",
				},
				Logging: nullLogger,
				Ports:   []string{"2345:5432"},
			},
			hubServiceName: {
				Image:         "selenium/hub:3.141.59-oxygen",
				ContainerName: hubServiceName,
				Logging:       nullLogger,
				Ports: []string{
					"4444:4444",
				},
			},
			"test": {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/frontend-tests.Dockerfile",
				},
				DependsOn: []string{"firefox", "chrome"},
				Environment: map[string]string{
					"DOCKER":         "true",
					"TARGET_ADDRESS": fmt.Sprintf("http://%s:8888", serviceName),
				},
				Links: []string{hubServiceName, serviceName},
			},
			serviceName: {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/server.Dockerfile",
				},
				Environment: map[string]string{
					"CONFIGURATION_FILEPATH": "config_files/production.toml",
					"DOCKER":                 "true",
				},
				Links: []string{
					"database",
				},
				Ports: []string{
					"80:8888",
				},
			},
		},
	}
}

func integrationCoverageDotJSON(projectName wordsmith.SuperPalabra) models.DockerComposeFile {
	return models.DockerComposeFile{
		Version: "3.3",
		Services: map[string]models.DockerComposeService{
			"database": {
				Image: pgImage,
				Environment: map[string]string{
					"POSTGRES_DB":       "todo",
					"POSTGRES_PASSWORD": "hunter2",
					"POSTGRES_USER":     "dbuser",
				},
				Logging: nullLogger,
				Ports:   []string{"2345:5432"},
			},
			"test": {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/integration-tests.Dockerfile",
				},
				Links: []string{
					"coverage-server",
				},
				Environment: map[string]string{
					"DOCKER":                           "true",
					"JAEGER_AGENT_HOST":                "tracing-server",
					"JAEGER_AGENT_PORT":                "6831",
					"JAEGER_SAMPLER_MANAGER_HOST_PORT": "tracing-server:5778",
					"JAEGER_SERVICE_NAME":              "coverage-server",
					"TARGET_ADDRESS":                   "http://coverage-server",
					"WAIT_FOR_COVERAGE":                "yes",
				},
			},
			"coverage-server": {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/integration-coverage-server.Dockerfile",
				},
				Environment: map[string]string{
					"CONFIGURATION_FILEPATH": "config_files/coverage.toml",
					"DOCKER":                 "true",
					"RUNTIME_DURATION":       "30s",
				},
				Links: []string{
					"database",
				},
				Ports: []string{
					"80:8888",
				},
				Volumes: []models.DockerVolume{
					{
						Source: "../artifacts/",
						Target: "/home/",
						Type:   "bind",
					},
				},
			},
		},
	}
}

func productionDotJSON(projectName wordsmith.SuperPalabra) models.DockerComposeFile {
	serviceName := fmt.Sprintf("%s-server", projectName.KebabName())

	return models.DockerComposeFile{
		Version: "3.3",
		Services: map[string]models.DockerComposeService{
			"database": {
				Image: pgImage,
				Environment: map[string]string{
					"POSTGRES_DB":       "todo",
					"POSTGRES_PASSWORD": "hunter2",
					"POSTGRES_USER":     "dbuser",
				},
				Logging: nullLogger,
				Ports:   []string{"2345:5432"},
			},
			serviceName: {
				Build: &models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/server.Dockerfile",
				},
				DependsOn: []string{
					"prometheus",
					"grafana",
				},
				Environment: map[string]string{
					"CONFIGURATION_FILEPATH":           "config_files/production.toml",
					"DOCKER":                           "true",
					"JAEGER_AGENT_HOST":                "tracing-server",
					"JAEGER_AGENT_PORT":                "6831",
					"JAEGER_SAMPLER_MANAGER_HOST_PORT": "tracing-server:5778",
					"JAEGER_SERVICE_NAME":              serviceName,
				},
				Links: []string{
					"tracing-server",
					"database",
				},
				Ports: []string{
					"80:8888",
				},
			},
			"tracing-server": {
				Image:   "jaegertracing/all-in-one:latest",
				Logging: nullLogger,
				Ports: []string{
					"5775:5775/udp",
					"6831:6831/udp",
					"6832:6832/udp",
					"5778:5778",
					"16686:16686",
					"14268:14268",
					"9411:9411",
				},
			},
			"grafana": {
				Image:   "grafana/grafana",
				Logging: nullLogger,
				Links: []string{
					"prometheus",
				},
				Ports: []string{
					"3000:3000",
				},
				Volumes: []models.DockerVolume{
					{
						Source: "../deploy/grafana/local/provisioning",
						Target: "/etc/grafana/provisioning",
						Type:   "bind",
					},
					{
						Source: "../deploy/grafana/dashboards",
						Target: "/etc/grafana/dashboards",
						Type:   "bind",
					},
				},
			},
			"prometheus": {
				Image:   "quay.io/prometheus/prometheus:v2.0.0",
				Command: "--config.file=/etc/prometheus/config.yaml --storage.tsdb.path=/prometheus",
				Logging: nullLogger,
				Ports: []string{
					"9090:9090",
				},
				Volumes: []models.DockerVolume{
					{
						Source: "../deploy/prometheus/local/config.yaml",
						Target: "/etc/prometheus/config.yaml",
						Type:   "bind",
					},
				},
			},
		},
	}
}
