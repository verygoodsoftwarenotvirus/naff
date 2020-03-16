package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"internal/v1/config/wire.go":         wireDotGo(pkg),
		"internal/v1/config/config.go":       configDotGo(pkg),
		"internal/v1/config/config_test.go":  configTestDotGo(pkg),
		"internal/v1/config/database.go":     databaseDotGo(pkg),
		"internal/v1/config/doc.go":          docDotGo(),
		"internal/v1/config/metrics.go":      metricsDotGo(pkg),
		"internal/v1/config/metrics_test.go": metricsTestDotGo(pkg),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg, path, file); err != nil {
			return err
		}
	}

	return nil
}
