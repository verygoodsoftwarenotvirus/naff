package utils

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "utils"

	basePackagePath = "tests/utils"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"mock_matchers.go":             mockMatchersDotGo(proj),
		"mock_reader.go":               mockReaderDotGo(proj),
		"testutil.go":                  testutilDotGo(proj),
		"utils.go":                     utilsDotGo(proj),
		"doc.go":                       docDotGo(proj),
		"mock_handler.go":              mockHandlerDotGo(proj),
		"mock_http_response_writer.go": mockHTTPResponseWriterDotGo(proj),
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
