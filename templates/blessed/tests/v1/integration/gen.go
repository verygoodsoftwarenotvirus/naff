package integration

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"tests/v1/integration/init.go":          initDotGo(proj),
		"tests/v1/integration/meta_test.go":     metaTestDotGo(proj),
		"tests/v1/integration/oauth2_test.go":   oauth2TestDotGo(proj),
		"tests/v1/integration/users_test.go":    usersTestDotGo(proj),
		"tests/v1/integration/webhooks_test.go": webhooksTestDotGo(proj),
		"tests/v1/integration/auth_test.go":     authTestDotGo(proj),
		"tests/v1/integration/doc.go":           docDotGo(),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("tests/v1/integration/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj.OutputPath, path, file); err != nil {
			return err
		}
	}

	return nil
}
