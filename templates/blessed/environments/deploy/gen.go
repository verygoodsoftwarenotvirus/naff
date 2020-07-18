package deploy

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(project *models.Project) error {
	files := map[string]func(project *models.Project) []byte{
		"deploy/prometheus/local/config.yaml":                    prometheusLocalConfigDotYAML,
		"deploy/grafana/dashboards/dashboard.json":               dashboardDotJSON,
		"deploy/grafana/local/provisioning/dashboards/all.yaml":  grafanaLocalProvisioningDashboardsAllDotYAML,
		"deploy/grafana/local/provisioning/datasources/all.yaml": grafanLocalProvisioningDataSourcesAllDotYAML,
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

		if _, err := f.Write(file(project)); err != nil {
			log.Printf("error writing to file: %v", err)
			return err
		}
	}

	return nil
}

func prometheusLocalConfigDotYAML(project *models.Project) []byte {
	serviceName := project.Name.KebabName()

	return []byte(fmt.Sprintf(`global:
  # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  scrape_interval:     15s
  # Evaluate rules every 15 seconds. The default is every 1 minute.
  evaluation_interval: 15s
  # scrape_timeout is set to the global default (10s).

scrape_configs:
  - job_name: '%s-server'

    static_configs:
      - targets: ['%s-server:8888']

    # How frequently to scrape targets from this job.
    scrape_interval: 15s

    # Per-scrape timeout when scraping this job.
    scrape_timeout: 15s

    # The HTTP resource path on which to fetch metrics from targets.
    metrics_path: '/metrics'

    # Configures the protocol scheme used for requests.
    scheme: 'http'

    # Optional HTTP URL parameters.
    # params:
    #   key: 'value'

    # Sets the `+"`"+`Authorization`+"`"+` header on every scrape request with the
    # configured username and password.
    # password and password_file are mutually exclusive.
    #  basic_auth:
    #    username: <string>
    #    password: <secret>
    #    password_file: <string>

    # Sets the `+"`"+`Authorization`+"`"+` header on every scrape request with
    # the configured bearer token. It is mutually exclusive with `+"`"+`bearer_token_file`+"`"+`.
    # bearer_token: <secret>

    # Sets the `+"`"+`Authorization`+"`"+` header on every scrape request with the bearer token
    # read from the configured file. It is mutually exclusive with `+"`"+`bearer_token`+"`"+`.
    #  bearer_token_file: /path/to/bearer/token/file

    # Configures the scrape request's TLS settings.
    tls_config:
      # CA certificate to validate API server certificate with.
      #   ca_file: '/path/to/file'

      # Certificate and key files for client cert authentication to the server.
      #   cert_file: '/path/to/file'
      #   key_file: '/path/to/file'

      # ServerName extension to indicate the name of the server.
      # https://tools.ietf.org/html/rfc4366#section-3.1
      #   server_name: ''

      # Disable validation of the server certificate.
      insecure_skip_verify: true

    # Optional proxy URL.
    #  proxy_url: ''

    # Per-scrape limit on number of scraped samples that will be accepted.
    # If more than this number of samples are present after metric relabelling
    # the entire scrape will be treated as failed. 0 means no limit.
    sample_limit:  0


#alerting:
#  alertmanagers:
#    - scheme: https
#      static_configs:
#        - targets:
#            - "1.2.3.4:9093"
#            - "1.2.3.5:9093"
#            - "1.2.3.6:9093"
`, serviceName, serviceName))
}

func grafanaLocalProvisioningDashboardsAllDotYAML(project *models.Project) []byte {
	return []byte(`apiVersion: 1

providers:
  # <string> an unique provider name
  - name: 'default'
    # <int> org id. will default to orgId 1 if not specified
    orgId: 1
    # <string, required> name of the dashboard folder. Required
    folder: '/etc/grafana/dashboards'
    # <string> folder UID. will be automatically generated if not specified
    folderUid: ''
    # <string, required> provider type. Required
    type: file
    # <bool> disable dashboard deletion
    disableDeletion: false
    # <bool> enable dashboard editing
    editable: true
    # <int> how often Grafana will scan for changed dashboards
    updateIntervalSeconds: 10
    # <bool> allow updating provisioned dashboards from the UI
    allowUiUpdates: false
    options:
      # <string, required> path to dashboard files on disk. Required
      path: '/etc/grafana/dashboards'
`)
}

func grafanLocalProvisioningDataSourcesAllDotYAML(project *models.Project) []byte {
	return []byte(`apiVersion: 1

# Thanks to https://ops.tips/blog/initialize-grafana-with-preconfigured-dashboards/#configuring-grafana
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
