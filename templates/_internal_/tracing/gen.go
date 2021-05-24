package tracing

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName     = "tracing"
	basePackagePath = "internal/observability/tracing"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":            docDotGo(),
		"spans.go":          spansDotGo(proj),
		"span_attachers.go": spanAttachersDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
