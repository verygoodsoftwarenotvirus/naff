package auth

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"internal/v1/auth/authenticator.go":      authenticatorDotGo(proj),
		"internal/v1/auth/authenticator_test.go": authenticatorTestDotGo(proj),
		"internal/v1/auth/bcrypt.go":             bcryptDotGo(proj),
		"internal/v1/auth/bcrypt_test.go":        bcryptTestDotGo(proj),
		"internal/v1/auth/doc.go":                docDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
