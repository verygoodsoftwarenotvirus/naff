package metrics

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"internal/v1/metrics/runtime.go":      runtimeDotGo(pkgRoot, types),
		"internal/v1/metrics/types.go":        typesDotGo(pkgRoot, types),
		"internal/v1/metrics/counter_test.go": counterTestDotGo(pkgRoot, types),
		"internal/v1/metrics/doc.go":          docDotGo(),
		"internal/v1/metrics/meta_test.go":    metaTestDotGo(pkgRoot, types),
		"internal/v1/metrics/wire.go":         wireDotGo(pkgRoot, types),
		"internal/v1/metrics/counter.go":      counterDotGo(pkgRoot, types),
		"internal/v1/metrics/meta.go":         metaDotGo(pkgRoot, types),
		"internal/v1/metrics/runtime_test.go": runtimeTestDotGo(pkgRoot, types),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
