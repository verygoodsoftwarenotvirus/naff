package integration

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"tests/v1/integration/init.go":          initDotGo(pkgRoot),
		"tests/v1/integration/meta_test.go":     metaTestDotGo(),
		"tests/v1/integration/oauth2_test.go":   oauth2TestDotGo(pkgRoot),
		"tests/v1/integration/users_test.go":    usersTestDotGo(pkgRoot),
		"tests/v1/integration/webhooks_test.go": webhooksTestDotGo(pkgRoot),
		"tests/v1/integration/auth_test.go":     authTestDotGo(pkgRoot),
		"tests/v1/integration/doc.go":           docDotGo(),
	}

	for _, typ := range types {
		files[fmt.Sprintf("tests/v1/integration/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(pkgRoot, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
