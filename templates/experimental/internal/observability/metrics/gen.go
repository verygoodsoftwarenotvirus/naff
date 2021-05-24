package metrics

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "metrics"

	basePackagePath = "internal/observability/metrics"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":           docDotGo(proj),
		"types.go":         typesDotGo(proj),
		"types_test_.go":   typesTestDotGo(proj),
		"wire.go":          wireDotGo(proj),
		"config.go":        configDotGo(proj),
		"config_test_.go":  configTestDotGo(proj),
		"counter.go":       counterDotGo(proj),
		"counter_test_.go": counterTestDotGo(proj),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
