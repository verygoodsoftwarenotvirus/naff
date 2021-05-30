package encoded_qr_code_generator

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "encoded_qr_code_generator"

	basePackagePath = "cmd/tools/encoded_qr_code_generator"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"main.go": mainDotGoString(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed main.gotpl
var mainTemplate string

func mainDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, mainTemplate, nil)
}
