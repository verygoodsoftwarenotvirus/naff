package metrics

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"internal/v1/metrics/runtime.go":      runtimeDotGo(pkg),
		"internal/v1/metrics/types.go":        typesDotGo(pkg),
		"internal/v1/metrics/counter_test.go": counterTestDotGo(pkg),
		"internal/v1/metrics/doc.go":          docDotGo(),
		"internal/v1/metrics/meta_test.go":    metaTestDotGo(pkg),
		"internal/v1/metrics/wire.go":         wireDotGo(pkg),
		"internal/v1/metrics/counter.go":      counterDotGo(pkg),
		"internal/v1/metrics/meta.go":         metaDotGo(pkg),
		"internal/v1/metrics/runtime_test.go": runtimeTestDotGo(pkg),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg.OutputPath, path, file); err != nil {
			return err
		}
	}

	return nil
}
