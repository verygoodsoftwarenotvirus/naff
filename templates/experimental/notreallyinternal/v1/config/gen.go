package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"internal/v1/config/wire.go":         wireDotGo(),
		"internal/v1/config/config.go":       configDotGo(),
		"internal/v1/config/config_test.go":  configTestDotGo(),
		"internal/v1/config/database.go":     databaseDotGo(),
		"internal/v1/config/doc.go":          docDotGo(),
		"internal/v1/config/metrics.go":      metricsDotGo(),
		"internal/v1/config/metrics_test.go": metricsTestDotGo(),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
