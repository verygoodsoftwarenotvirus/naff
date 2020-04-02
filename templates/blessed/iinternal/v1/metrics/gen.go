package metrics

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"internal/v1/metrics/runtime.go":      runtimeDotGo(proj),
		"internal/v1/metrics/types.go":        typesDotGo(proj),
		"internal/v1/metrics/counter_test.go": counterTestDotGo(proj),
		"internal/v1/metrics/doc.go":          docDotGo(),
		"internal/v1/metrics/meta_test.go":    metaTestDotGo(proj),
		"internal/v1/metrics/wire.go":         wireDotGo(proj),
		"internal/v1/metrics/counter.go":      counterDotGo(proj),
		"internal/v1/metrics/meta.go":         metaDotGo(proj),
		"internal/v1/metrics/runtime_test.go": runtimeTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
