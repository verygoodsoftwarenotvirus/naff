package sqlite

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "sqlite"

	basePackagePath = "internal/database/querybuilding/sqlite"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"migrations.go":                     migrationsDotGo(proj),
		"webhooks_test_.go":                 webhooksTestDotGo(proj),
		"items.go":                          itemsDotGo(proj),
		"doc.go":                            docDotGo(proj),
		"items_test_.go":                    itemsTestDotGo(proj),
		"sqlite.go":                         sqliteDotGo(proj),
		"users_test_.go":                    usersTestDotGo(proj),
		"api_clients_test_.go":              apiClientsTestDotGo(proj),
		"webhooks.go":                       webhooksDotGo(proj),
		"generic_test_.go":                  genericTestDotGo(proj),
		"account_user_memberships_test_.go": accountUserMembershipsTestDotGo(proj),
		"accounts.go":                       accountsDotGo(proj),
		"accounts_test_.go":                 accountsTestDotGo(proj),
		"api_clients.go":                    apiClientsDotGo(proj),
		"audit_log_entries.go":              auditLogEntriesDotGo(proj),
		"audit_log_entries_test_.go":        auditLogEntriesTestDotGo(proj),
		"generic.go":                        genericDotGo(proj),
		"account_user_memberships.go":       accountUserMembershipsDotGo(proj),
		"users.go":                          usersDotGo(proj),
		"wire.go":                           wireDotGo(proj),
		"sqlite_test_.go":                   sqliteTestDotGo(proj),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
