package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"internal/v1/config/wire.go":         wireDotGo(pkgRoot, types),
		"internal/v1/config/config.go":       configDotGo(pkgRoot, types),
		"internal/v1/config/config_test.go":  configTestDotGo(pkgRoot, types),
		"internal/v1/config/database.go":     databaseDotGo(pkgRoot, types),
		"internal/v1/config/doc.go":          docDotGo(),
		"internal/v1/config/metrics.go":      metricsDotGo(pkgRoot, types),
		"internal/v1/config/metrics_test.go": metricsTestDotGo(pkgRoot, types),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
