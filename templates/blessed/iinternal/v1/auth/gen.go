package auth

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"internal/v1/auth/authenticator.go":      authenticatorDotGo(pkgRoot, types),
		"internal/v1/auth/authenticator_test.go": authenticatorTestDotGo(pkgRoot, types),
		"internal/v1/auth/bcrypt.go":             bcryptDotGo(pkgRoot, types),
		"internal/v1/auth/bcrypt_test.go":        bcryptTestDotGo(pkgRoot, types),
		"internal/v1/auth/doc.go":                docDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
