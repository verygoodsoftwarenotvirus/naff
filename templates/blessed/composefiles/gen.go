package composefiles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

var (
	nullLogger = &models.DockerComposeLogging{
		Driver: "none",
	}
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, projectName wordsmith.SuperPalabra, types []models.DataType) error {
	composeFiles := map[string]models.DockerComposeFile{
		"compose-files/development.json": buildDevelopmentComposeFile(projectName),
	}

	for filename, file := range composeFiles {
		fn := utils.BuildTemplatePath(filename)

		f, _ := json.MarshalIndent(file, "", "  ")
		if err := ioutil.WriteFile(fn, f, 0644); err != nil {
			log.Printf("error rendering %q: %v\n", filename, err)
		}
	}

	return nil
}

func buildDevelopmentComposeFile(projectName wordsmith.SuperPalabra) models.DockerComposeFile {
	return models.DockerComposeFile{
		Version: "3.3",
		Services: map[string]models.DockerComposeService{
			"database": {
				Image: "postgres:latest",
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
			fmt.Sprintf("%s-server", projectName.KebabName()): {
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
