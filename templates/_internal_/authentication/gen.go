package authentication

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName     = "authentication"
	testPackageName = "authentication_test"
	basePackagePath = "internal/authentication"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"authenticator.go":      authenticatorDotGo(proj),
		"authenticator_test.go": authenticatorTestDotGo(proj),
		"doc.go":                docDotGo(),
		"wire.go":               wireDotGo(),
		"mock_authenticator.go": mockAuthenticatorDotGo(proj),
		"config.go":             configDotGo(proj),
		"config_test.go":        configTestDotGo(proj),
		"argon2.go":             argon2DotGo(proj),
		"argon2_test.go":        argon2TestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
