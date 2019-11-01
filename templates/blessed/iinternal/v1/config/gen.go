package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"internal/v1/config/wire.go":         wireDotGo(),
		"internal/v1/config/config.go":       configDotGo(pkgRoot),
		"internal/v1/config/config_test.go":  configTestDotGo(pkgRoot),
		"internal/v1/config/database.go":     databaseDotGo(pkgRoot),
		"internal/v1/config/doc.go":          docDotGo(),
		"internal/v1/config/metrics.go":      metricsDotGo(pkgRoot),
		"internal/v1/config/metrics_test.go": metricsTestDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
