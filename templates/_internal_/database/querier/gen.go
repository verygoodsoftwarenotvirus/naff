package querier

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "querier"

	basePackagePath = "internal/database/querier"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"client.go":              clientDotGo(proj),
		"doc.go":                 docDotGo(),
		"oauth2_clients_test.go": oauth2ClientsTestDotGo(proj),
		"users.go":               usersDotGo(proj),
		"users_test.go":          usersTestDotGo(proj),
		"webhooks_test.go":       webhooksTestDotGo(proj),
		"client_test.go":         clientTestDotGo(proj),
		"oauth2_clients.go":      oauth2ClientsDotGo(proj),
		"webhooks.go":            webhooksDotGo(proj),
		"wire.go":                wireDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
		files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
