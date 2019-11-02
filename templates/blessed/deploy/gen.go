package deploy

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, projectName wordsmith.SuperPalabra, types []models.DataType) error {
	files := map[string]func(wordsmith.SuperPalabra, []models.DataType) []byte{
		"deploy/prometheus/local/config.yaml":                    prometheusLocalConfigDotYAML,
		"deploy/grafana/dashboards/dashboard.json":               dashboardDotJSON,
		"deploy/grafana/local/provisioning/dashboards/all.yaml":  grafanaLocalProvisioningDashboardsAllDotYAML,
		"deploy/grafana/local/provisioning/datasources/all.yaml": grafanLocalProvisioningDataSourcesAllDotYAML,
	}

	for filename, file := range files {
		fname := utils.BuildTemplatePath(filename)

		if mkdirErr := os.MkdirAll(filepath.Dir(fname), os.ModePerm); mkdirErr != nil {
			log.Printf("error making directory: %v\n", mkdirErr)
		}

		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("error opening file: %v", err)
			return err
		}

		bytes := file(projectName, types)
		if _, err := f.Write(bytes); err != nil {
			log.Printf("error writing to file: %v", err)
			return err
		}
	}

	return nil
}

func prometheusLocalConfigDotYAML(service wordsmith.SuperPalabra, _ []models.DataType) []byte {
	serviceName := service.KebabName()

	return []byte(fmt.Sprintf(`# my global config
global:
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `+"`"+`job=<job_name>`+"`"+` to any timeseries scraped from this config.
  - job_name: '%s-server'

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.
    static_configs:
      - targets: ['%s-server']

    tls_config:
      insecure_skip_verify: true
`, serviceName, serviceName))
}

func grafanaLocalProvisioningDashboardsAllDotYAML(_ wordsmith.SuperPalabra, _ []models.DataType) []byte {
	return []byte(`- name: 'default' # name of this dashboard configuration (not dashboard itself)
  org_id: 1 # id of the org to hold the dashboard
  folder: '' # name of the folder to put the dashboard (http://docs.grafana.org/v5.0/reference/dashboard_folders/)
  type: 'file' # type of dashboard description (json files)
  options:
    folder: '/etc/grafana/dashboards' # where dashboards are
`)
}

func grafanLocalProvisioningDataSourcesAllDotYAML(_ wordsmith.SuperPalabra, _ []models.DataType) []byte {
	return []byte(`apiVersion: 1

# Gracias a https://ops.tips/blog/initialize-grafana-with-preconfigured-dashboards/#configuring-grafana
datasources:
  - access: 'proxy' # make grafana perform the requests
    version: 1 # well, versioning
    is_default: true # whether this should be the default DS
    name: 'prometheus' # name of the datasource
    type: 'prometheus' # type of the data source
    org_id: 1 # id of the organization to tie this datasource to
    url: 'http://prometheus:9090' # url of the prom instance
`)
}
