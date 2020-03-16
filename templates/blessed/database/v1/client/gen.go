package client

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"database/v1/client/client.go":              clientDotGo(pkg),
		"database/v1/client/doc.go":                 docDotGo(),
		"database/v1/client/oauth2_clients_test.go": oauth2ClientsTestDotGo(pkg),
		"database/v1/client/users.go":               usersDotGo(pkg),
		"database/v1/client/users_test.go":          usersTestDotGo(pkg),
		"database/v1/client/webhooks_test.go":       webhooksTestDotGo(pkg),
		"database/v1/client/client_test.go":         clientTestDotGo(pkg),
		"database/v1/client/oauth2_clients.go":      oauth2ClientsDotGo(pkg),
		"database/v1/client/webhooks.go":            webhooksDotGo(pkg),
		"database/v1/client/wire.go":                wireDotGo(pkg),
	}

	for _, typ := range pkg.DataTypes {
		files[fmt.Sprintf("database/v1/client/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(pkg, typ)
		files[fmt.Sprintf("database/v1/client/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(pkg, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg, path, file); err != nil {
			return err
		}
	}

	return nil
}
