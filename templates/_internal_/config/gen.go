package config

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "config"

	basePackagePath = "internal/config"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"wire.go":         wireDotGo(proj),
		"config.go":       configDotGo(proj),
		"config_test.go":  configTestDotGo(proj),
		"database.go":     databaseDotGo(proj),
		"doc.go":          docDotGo(),
		"metrics.go":      metricsDotGo(proj),
		"metrics_test.go": metricsTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
