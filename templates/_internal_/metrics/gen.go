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
		"runtime.go":      runtimeDotGo(proj),
		"types.go":        typesDotGo(proj),
		"counter_test.go": counterTestDotGo(proj),
		"doc.go":          docDotGo(),
		"wire.go":         wireDotGo(proj),
		"counter.go":      counterDotGo(proj),
		"runtime_test.go": runtimeTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
