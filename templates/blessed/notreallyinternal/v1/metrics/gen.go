package metrics

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"internal/v1/metrics/runtime.go":      runtimeDotGo(),
		"internal/v1/metrics/types.go":        typesDotGo(),
		"internal/v1/metrics/counter_test.go": counterTestDotGo(),
		"internal/v1/metrics/doc.go":          docDotGo(),
		"internal/v1/metrics/meta_test.go":    metaTestDotGo(),
		"internal/v1/metrics/wire.go":         wireDotGo(),
		"internal/v1/metrics/counter.go":      counterDotGo(),
		"internal/v1/metrics/meta.go":         metaDotGo(),
		"internal/v1/metrics/runtime_test.go": runtimeTestDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
