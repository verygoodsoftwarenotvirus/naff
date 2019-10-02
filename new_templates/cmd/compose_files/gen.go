package compose_files

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	outputPath = "compose_files"
	bind       = "bind"

	db     = "database"
	tracer = "tracing-server"
)

var (
	noLogging = models.DockerComposeLogging{
		Driver: "none",
	}

	developmentComposeFile = &models.DockerComposeFile{
		Version: "3.3",
		Services: map[string]models.DockerComposeService{
			db: {
				Image: "postgres:latest",
				Environment: map[string]string{
					"POSTGRES_DB":       "todo",
					"POSTGRES_PASSWORD": "hunter2",
					"POSTGRES_USER":     "dbuser",
				},
				Ports:   []string{"2345:5432"},
				Logging: noLogging,
			},
			"grafana": {
				Image: "grafana/grafana",
				Links: []string{"prometheus"},
				Volumes: []models.DockerVolume{
					{
						Source: "../deploy/grafana/local/provisioning",
						Target: "/etc/grafana/local/provisioning",
						Type:   bind,
					},
					{
						Source: "../deploy/grafana/dashboards",
						Target: "/etc/grafana/dashboards",
						Type:   bind,
					},
				},
				Logging: noLogging,
			},
			"prometheus": {
				Image:   "quay.io/prometheus/prometheus:v2.0.0",
				Command: "--config.file=/etc/prometheus/config.yaml --storage.tsdb.path=/prometheus",
				Ports:   []string{"9090:9090"},
				Volumes: []models.DockerVolume{
					{
						Source: "../deploy/prometheus/local/config.yaml",
						Target: "/etc/prometheus/config.yaml",
						Type:   bind,
					},
				},
				Logging: noLogging,
			},
			"todo-server": {
				Build: models.DockerComposeBuild{
					Context:    "../",
					Dockerfile: "dockerfiles/development.Dockerfile",
				},
				Ports:     []string{"80:8888"},
				Links:     []string{tracer, db},
				DependsOn: []string{"prometheus", "grafana"},
				Environment: map[string]string{
					"DOCKER":                           "true",
					"CONFIGURATION_FILEPATH":           "config_files/development.toml",
					"JAEGER_AGENT_HOST":                "tracing-server",
					"JAEGER_AGENT_PORT":                "6831",
					"JAEGER_SAMPLER_MANAGER_HOST_PORT": "tracing-server:5778",
					"JAEGER_SERVICE_NAME":              "todo-server",
				},
				Volumes: []models.DockerVolume{
					{
						Source: "../frontend/v1/public",
						Target: "/frontend",
						Type:   bind,
					},
				},
			},
			tracer: {
				Image: "jaegertracing/all-in-one:latest",
				Ports: []string{
					"6831:6831/udp",
					"5778:5778",
					"16686:16686",
				},
				Logging: noLogging,
			},
		},
	}
)

// RenderFolder renders the folder
func RenderFolder(basePath string) error {
	fileMap := map[string]func(path string) error{
		"development.json": developmentDotJSON,
	}

	if err := os.MkdirAll(filepath.Join(basePath, outputPath), os.ModeDir); err != nil {
		return err
	}

	for path, fun := range fileMap {
		if err := fun(path); err != nil {
			return err
		}
	}

	return nil
}

func developmentDotJSON(path string) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&developmentComposeFile); err != nil {
		return err
	}

	return ioutil.WriteFile(path, b.Bytes(), 0644)
}
