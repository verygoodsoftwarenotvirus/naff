package model

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"tests/v1/testutil/rand/model/doc.go":           docDotGo(),
		"tests/v1/testutil/rand/model/oauth2clients.go": oauth2ClientsDotGo(pkgRoot, types),
		"tests/v1/testutil/rand/model/users.go":         usersDotGo(pkgRoot, types),
		"tests/v1/testutil/rand/model/webhooks.go":      webhooksDotGo(pkgRoot, types),
	}

	for _, typ := range types {
		files[fmt.Sprintf("tests/v1/testutil/rand/model/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(pkgRoot, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
