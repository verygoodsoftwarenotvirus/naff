package mariadb

import (
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
		"accounts_test_.go":                 accountsTestDotGo(proj),
		"audit_log_entries.go":              auditLogEntriesDotGo(proj),
		"audit_log_entries_test_.go":        auditLogEntriesTestDotGo(proj),
		"users.go":                          usersDotGo(proj),
		"account_user_memberships.go":       accountUserMembershipsDotGo(proj),
		"account_user_memberships_test_.go": accountUserMembershipsTestDotGo(proj),
		"migrations.go":                     migrationsDotGo(proj),
		"webhooks_test_.go":                 webhooksTestDotGo(proj),
		"api_clients.go":                    apiClientsDotGo(proj),
		"mariadb_test_.go":                  mariadbTestDotGo(proj),
		"users_test_.go":                    usersTestDotGo(proj),
		"accounts.go":                       accountsDotGo(proj),
		"api_clients_test_.go":              apiClientsTestDotGo(proj),
		"doc.go":                            docDotGo(proj),
		"generic.go":                        genericDotGo(proj),
		"generic_test_.go":                  genericTestDotGo(proj),
		"items.go":                          itemsDotGo(proj),
		"items_test_.go":                    itemsTestDotGo(proj),
		"mariadb.go":                        mariadbDotGo(proj),
		"webhooks.go":                       webhooksDotGo(proj),
		"wire.go":                           wireDotGo(proj),
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
