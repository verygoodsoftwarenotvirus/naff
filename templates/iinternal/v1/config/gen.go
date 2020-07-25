package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"internal/v1/config/wire.go":         wireDotGo(proj),
		"internal/v1/config/config.go":       configDotGo(proj),
		"internal/v1/config/config_test.go":  configTestDotGo(proj),
		"internal/v1/config/database.go":     databaseDotGo(proj),
		"internal/v1/config/doc.go":          docDotGo(),
		"internal/v1/config/metrics.go":      metricsDotGo(proj),
		"internal/v1/config/metrics_test.go": metricsTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
