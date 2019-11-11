package integration

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"tests/v1/integration/init.go":          initDotGo(pkg),
		"tests/v1/integration/meta_test.go":     metaTestDotGo(pkg),
		"tests/v1/integration/oauth2_test.go":   oauth2TestDotGo(pkg),
		"tests/v1/integration/users_test.go":    usersTestDotGo(pkg),
		"tests/v1/integration/webhooks_test.go": webhooksTestDotGo(pkg),
		"tests/v1/integration/auth_test.go":     authTestDotGo(pkg),
		"tests/v1/integration/doc.go":           docDotGo(),
	}

	for _, typ := range pkg.DataTypes {
		files[fmt.Sprintf("tests/v1/integration/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(pkg, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg.OutputPath, path, file); err != nil {
			return err
		}
	}

	return nil
}
