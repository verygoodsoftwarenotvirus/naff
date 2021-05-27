package encoding

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "encoding"

	basePackagePath = "internal/encoding"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"client_encoder.go":              clientEncoderDotGo(proj),
		"client_encoder_test.go":         clientEncoderTestDotGo(proj),
		"config.go":                      configDotGo(proj),
		"config_test.go":                 configTestDotGo(proj),
		"content_type.go":                contentTypeDotGo(proj),
		"content_type_test.go":           contentTypeTestDotGo(proj),
		"doc.go":                         docDotGo(proj),
		"server_encoder_decoder.go":      serverEncoderDecoderDotGo(proj),
		"server_encoder_decoder_test.go": serverEncoderDecoderTestDotGo(proj),
		"wire.go":                        wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
