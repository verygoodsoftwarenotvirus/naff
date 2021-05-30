package testutil

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "testutils"

	basePackagePath = "tests/utils"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":                       docDotGo(proj),
		"testutil.go":                  testutilDotGo(proj),
		"mock_matchers.go":             mockMatchersDotGo(proj),
		"mock_reader.go":               mockReaderDotGo(proj),
		"utils.go":                     utilsDotGo(proj),
		"mock_handler.go":              mockHandlerDotGo(proj),
		"mock_http_response_writer.go": mockHTTPResponseWriterDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
