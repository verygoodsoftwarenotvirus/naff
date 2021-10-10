package integration

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "tests/utils"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"doc.go":                       docDotGo(proj),
		"mock_handler.go":              mockHandlerDotGo(proj),
		"mock_http_response_writer.go": mockHTTPResponseWriterDotGo(proj),
		"mock_matchers.go":             mockMatchersDotGo(proj),
		"mock_reader.go":               mockReaderDotGo(proj),
		"testutil.go":                  testutilDotGo(proj),
		"utils.go":                     utilsDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed mock_handler.gotpl
var mockHandlerTemplate string

func mockHandlerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockHandlerTemplate, nil)
}

//go:embed mock_http_response_writer.gotpl
var mockHTTPResponseWriterTemplate string

func mockHTTPResponseWriterDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockHTTPResponseWriterTemplate, nil)
}

//go:embed mock_matchers.gotpl
var mockMatchersTemplate string

func mockMatchersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockMatchersTemplate, nil)
}

//go:embed mock_reader.gotpl
var mockReaderTemplate string

func mockReaderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockReaderTemplate, nil)
}

//go:embed testutil.gotpl
var testutilTemplate string

func testutilDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, testutilTemplate, nil)
}

//go:embed utils.gotpl
var utilsTemplate string

func utilsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, utilsTemplate, nil)
}
