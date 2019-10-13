package model

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"tests/v1/testutil/rand/model/doc.go":           docDotGo(),
		"tests/v1/testutil/rand/model/items.go":         itemsDotGo(),
		"tests/v1/testutil/rand/model/oauth2clients.go": oauth2ClientsDotGo(),
		"tests/v1/testutil/rand/model/users.go":         usersDotGo(),
		"tests/v1/testutil/rand/model/webhooks.go":      webhooksDotGo(),
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
