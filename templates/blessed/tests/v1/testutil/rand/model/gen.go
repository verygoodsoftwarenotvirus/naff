package model

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"tests/v1/testutil/rand/model/doc.go":           docDotGo(),
		"tests/v1/testutil/rand/model/oauth2clients.go": oauth2ClientsDotGo(pkg),
		"tests/v1/testutil/rand/model/users.go":         usersDotGo(pkg),
		"tests/v1/testutil/rand/model/webhooks.go":      webhooksDotGo(pkg),
	}

	for _, typ := range pkg.DataTypes {
		files[fmt.Sprintf("tests/v1/testutil/rand/model/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(pkg, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg, path, file); err != nil {
			return err
		}
	}

	return nil
}
