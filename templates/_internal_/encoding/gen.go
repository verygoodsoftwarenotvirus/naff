package encoding

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/encoding"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
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
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed client_encoder.gotpl
var clientEncoderTemplate string

func clientEncoderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, clientEncoderTemplate)
}

//go:embed client_encoder_test.gotpl
var clientEncoderTestTemplate string

func clientEncoderTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, clientEncoderTestTemplate)
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate)
}

//go:embed content_type.gotpl
var contentTypeTemplate string

func contentTypeDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, contentTypeTemplate)
}

//go:embed content_type_test.gotpl
var contentTypeTestTemplate string

func contentTypeTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, contentTypeTestTemplate)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate)
}

//go:embed server_encoder_decoder.gotpl
var serverEncoderDecoderTemplate string

func serverEncoderDecoderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, serverEncoderDecoderTemplate)
}

//go:embed server_encoder_decoder_test.gotpl
var serverEncoderDecoderTestTemplate string

func serverEncoderDecoderTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, serverEncoderDecoderTestTemplate)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate)
}
