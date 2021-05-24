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
		"client_encoder_test_.go":         clientEncoderTestDotGo(proj),
		"config.go":                       configDotGo(proj),
		"content_type_test_.go":           contentTypeTestDotGo(proj),
		"server_encoder_decoder.go":       serverEncoderDecoderDotGo(proj),
		"server_encoder_decoder_test_.go": serverEncoderDecoderTestDotGo(proj),
		"wire.go":                         wireDotGo(proj),
		"client_encoder.go":               clientEncoderDotGo(proj),
		"content_type.go":                 contentTypeDotGo(proj),
		"doc.go":                          docDotGo(proj),
		"config_test_.go":                 configTestDotGo(proj),
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
