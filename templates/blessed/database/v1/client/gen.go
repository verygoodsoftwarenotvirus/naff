package client

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"database/v1/client/client.go":              clientDotGo(proj),
		"database/v1/client/doc.go":                 docDotGo(),
		"database/v1/client/oauth2_clients_test.go": oauth2ClientsTestDotGo(proj),
		"database/v1/client/users.go":               usersDotGo(proj),
		"database/v1/client/users_test.go":          usersTestDotGo(proj),
		"database/v1/client/webhooks_test.go":       webhooksTestDotGo(proj),
		"database/v1/client/client_test.go":         clientTestDotGo(proj),
		"database/v1/client/oauth2_clients.go":      oauth2ClientsDotGo(proj),
		"database/v1/client/webhooks.go":            webhooksDotGo(proj),
		"database/v1/client/wire.go":                wireDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("database/v1/client/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
		files[fmt.Sprintf("database/v1/client/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
