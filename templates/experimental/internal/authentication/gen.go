package authentication

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "authentication"

	basePackagePath = "internal/authentication"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"config_test_.go":        configTestDotGo(proj),
		"doc.go":                 docDotGo(proj),
		"mock_authenticator.go":  mockAuthenticatorDotGo(proj),
		"argon2.go":              argon2DotGo(proj),
		"argon2_test_.go":        argon2TestDotGo(proj),
		"authenticator.go":       authenticatorDotGo(proj),
		"authenticator_test_.go": authenticatorTestDotGo(proj),
		"config.go":              configDotGo(proj),
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
