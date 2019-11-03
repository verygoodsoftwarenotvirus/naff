package client

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"database/v1/client/client.go":              clientDotGo(pkgRoot, types),
		"database/v1/client/doc.go":                 docDotGo(),
		"database/v1/client/oauth2_clients_test.go": oauth2ClientsTestDotGo(pkgRoot, types),
		"database/v1/client/users.go":               usersDotGo(pkgRoot, types),
		"database/v1/client/users_test.go":          usersTestDotGo(pkgRoot, types),
		"database/v1/client/webhooks_test.go":       webhooksTestDotGo(pkgRoot, types),
		"database/v1/client/client_test.go":         clientTestDotGo(pkgRoot, types),
		"database/v1/client/oauth2_clients.go":      oauth2ClientsDotGo(pkgRoot, types),
		"database/v1/client/webhooks.go":            webhooksDotGo(pkgRoot, types),
		"database/v1/client/wire.go":                wireDotGo(pkgRoot, types),
	}

	for _, typ := range types {
		files[fmt.Sprintf("database/v1/client/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(pkgRoot, typ)
		files[fmt.Sprintf("database/v1/client/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(pkgRoot, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
