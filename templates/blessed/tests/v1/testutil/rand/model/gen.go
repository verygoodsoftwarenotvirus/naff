package model

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"tests/v1/testutil/rand/model/doc.go":           docDotGo(),
		"tests/v1/testutil/rand/model/init.go":          initDotGo(proj),
		"tests/v1/testutil/rand/model/oauth2clients.go": oauth2ClientsDotGo(proj),
		"tests/v1/testutil/rand/model/users.go":         usersDotGo(proj),
		"tests/v1/testutil/rand/model/webhooks.go":      webhooksDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("tests/v1/testutil/rand/model/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
