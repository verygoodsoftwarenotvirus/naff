package client

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"database/v1/client/client.go":              clientDotGo(),
		"database/v1/client/doc.go":                 docDotGo(),
		"database/v1/client/oauth2_clients_test.go": oauth2ClientsTestDotGo(),
		"database/v1/client/users.go":               usersDotGo(),
		"database/v1/client/users_test.go":          usersTestDotGo(),
		"database/v1/client/webhooks_test.go":       webhooksTestDotGo(),
		"database/v1/client/client_test.go":         clientTestDotGo(),
		"database/v1/client/oauth2_clients.go":      oauth2ClientsDotGo(),
		"database/v1/client/webhooks.go":            webhooksDotGo(),
		"database/v1/client/wire.go":                wireDotGo(),
	}

	for _, typ := range types {
		files[fmt.Sprintf("database/v1/client/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(typ)
		files[fmt.Sprintf("database/v1/client/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(typ)
	}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
