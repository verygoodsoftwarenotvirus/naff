package mockencoding

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/encoding/mock"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"client_encoder.go": clientEncoderDotGo(proj),
		"doc.go":            docDotGo(proj),
		"encoding.go":       encodingDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed client_encoder.gotpl
var clientEncoderTemplate string

func clientEncoderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, clientEncoderTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed encoding.gotpl
var encodingTemplate string

func encodingDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, encodingTemplate, nil)
}
