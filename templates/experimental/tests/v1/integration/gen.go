package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"tests/v1/integration/init.go":          initDotGo(),
		"tests/v1/integration/items_test.go":    itemsTestDotGo(),
		"tests/v1/integration/meta_test.go":     metaTestDotGo(),
		"tests/v1/integration/oauth2_test.go":   oauth2TestDotGo(),
		"tests/v1/integration/users_test.go":    usersTestDotGo(),
		"tests/v1/integration/webhooks_test.go": webhooksTestDotGo(),
		"tests/v1/integration/auth_test.go":     authTestDotGo(),
		"tests/v1/integration/doc.go":           docDotGo(),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
