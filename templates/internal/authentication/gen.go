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
		"bcrypt.go":             bcryptDotGo(proj),
		"bcrypt_test.go":        bcryptTestDotGo(proj),
		"doc.go":                docDotGo(),
		"mock_authenticator.go": mockAuthenticatorDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
