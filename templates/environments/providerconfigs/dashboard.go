package providerconfigs

import (
	_ "embed"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

//go:embed grafana_dashboard.jsontpl
var grafanaDashboard string

func dashboardDotJSON(proj *models.Project) string {
	return models.RenderCodeFile(proj, grafanaDashboard)
}
