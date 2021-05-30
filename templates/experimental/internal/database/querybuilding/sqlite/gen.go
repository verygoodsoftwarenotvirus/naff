package sqlite

import (
	_ "embed"
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
		"accounts.go":                      accountsDotGo(proj),
		"api_clients_test.go":              apiClientsTestDotGo(proj),
		"audit_log_entries.go":             auditLogEntriesDotGo(proj),
		"generic.go":                       genericDotGo(proj),
		"users_test.go":                    usersTestDotGo(proj),
		"doc.go":                           docDotGo(proj),
		"generic_test.go":                  genericTestDotGo(proj),
		"sqlite_test.go":                   sqliteTestDotGo(proj),
		"webhooks.go":                      webhooksDotGo(proj),
		"wire.go":                          wireDotGo(proj),
		"account_user_memberships.go":      accountUserMembershipsDotGo(proj),
		"account_user_memberships_test.go": accountUserMembershipsTestDotGo(proj),
		"items_test.go":                    itemsTestDotGo(proj),
		"migrations.go":                    migrationsDotGo(proj),
		"sqlite.go":                        sqliteDotGo(proj),
		"users.go":                         usersDotGo(proj),
		"webhooks_test.go":                 webhooksTestDotGo(proj),
		"accounts_test.go":                 accountsTestDotGo(proj),
		"api_clients.go":                   apiClientsDotGo(proj),
		"audit_log_entries_test.go":        auditLogEntriesTestDotGo(proj),
		"items.go":                         itemsDotGo(proj),
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
