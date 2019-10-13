package sqlite

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"database/v1/queriers/sqlite/migrations.go":          migrationsDotGo(),
		"database/v1/queriers/sqlite/oauth2_clients.go":      oauth2ClientsDotGo(),
		"database/v1/queriers/sqlite/oauth2_clients_test.go": oauth2ClientsTestDotGo(),
		"database/v1/queriers/sqlite/sqlite_test.go":         sqliteTestDotGo(),
		"database/v1/queriers/sqlite/users_test.go":          usersTestDotGo(),
		"database/v1/queriers/sqlite/wire.go":                wireDotGo(),
		"database/v1/queriers/sqlite/items_test.go":          itemsTestDotGo(),
		"database/v1/queriers/sqlite/items.go":               itemsDotGo(),
		"database/v1/queriers/sqlite/sqlite.go":              sqliteDotGo(),
		"database/v1/queriers/sqlite/users.go":               usersDotGo(),
		"database/v1/queriers/sqlite/webhooks.go":            webhooksDotGo(),
		"database/v1/queriers/sqlite/webhooks_test.go":       webhooksTestDotGo(),
		"database/v1/queriers/sqlite/doc.go":                 docDotGo(),
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
