package postgres

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"database/v1/queriers/postgres/oauth2_clients.go":      oauth2ClientsDotGo(),
		"database/v1/queriers/postgres/postgres.go":            postgresDotGo(),
		"database/v1/queriers/postgres/webhooks.go":            webhooksDotGo(),
		"database/v1/queriers/postgres/wire.go":                wireDotGo(),
		"database/v1/queriers/postgres/doc.go":                 docDotGo(),
		"database/v1/queriers/postgres/postgres_test.go":       postgresTestDotGo(),
		"database/v1/queriers/postgres/users.go":               usersDotGo(),
		"database/v1/queriers/postgres/users_test.go":          usersTestDotGo(),
		"database/v1/queriers/postgres/webhooks_test.go":       webhooksTestDotGo(),
		"database/v1/queriers/postgres/migrations.go":          migrationsDotGo(types),
		"database/v1/queriers/postgres/oauth2_clients_test.go": oauth2ClientsTestDotGo(),
	}

	for _, typ := range types {
		files[fmt.Sprintf("database/v1/queriers/postgres/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(typ)
		files[fmt.Sprintf("database/v1/queriers/postgres/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(typ)
	}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
