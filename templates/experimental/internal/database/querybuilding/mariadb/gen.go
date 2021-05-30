package mariadb

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "mariadb"

	basePackagePath = "internal/database/querybuilding/mariadb"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"account_user_memberships.go":      accountUserMembershipsDotGo(proj),
		"accounts_test.go":                 accountsTestDotGo(proj),
		"items.go":                         itemsDotGo(proj),
		"users.go":                         usersDotGo(proj),
		"audit_log_entries.go":             auditLogEntriesDotGo(proj),
		"generic_test.go":                  genericTestDotGo(proj),
		"mariadb_test.go":                  mariadbTestDotGo(proj),
		"migrations.go":                    migrationsDotGo(proj),
		"webhooks_test.go":                 webhooksTestDotGo(proj),
		"wire.go":                          wireDotGo(proj),
		"account_user_memberships_test.go": accountUserMembershipsTestDotGo(proj),
		"api_clients.go":                   apiClientsDotGo(proj),
		"api_clients_test.go":              apiClientsTestDotGo(proj),
		"generic.go":                       genericDotGo(proj),
		"mariadb.go":                       mariadbDotGo(proj),
		"webhooks.go":                      webhooksDotGo(proj),
		"accounts.go":                      accountsDotGo(proj),
		"audit_log_entries_test.go":        auditLogEntriesTestDotGo(proj),
		"doc.go":                           docDotGo(proj),
		"items_test.go":                    itemsTestDotGo(proj),
		"users_test.go":                    usersTestDotGo(proj),
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
