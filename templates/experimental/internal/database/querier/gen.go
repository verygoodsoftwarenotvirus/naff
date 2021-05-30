package querier

import (
	_ "embed"
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
		"admin.go":                         adminDotGo(proj),
		"admin_test.go":                    adminTestDotGo(proj),
		"items_test.go":                    itemsTestDotGo(proj),
		"migrate.go":                       migrateDotGo(proj),
		"querier.go":                       querierDotGo(proj),
		"querier_test.go":                  querierTestDotGo(proj),
		"account_user_memberships.go":      accountUserMembershipsDotGo(proj),
		"account_user_memberships_test.go": accountUserMembershipsTestDotGo(proj),
		"users.go":                         usersDotGo(proj),
		"audit_log_entries_test.go":        auditLogEntriesTestDotGo(proj),
		"users_test.go":                    usersTestDotGo(proj),
		"webhooks_test.go":                 webhooksTestDotGo(proj),
		"accounts.go":                      accountsDotGo(proj),
		"audit_log_entries.go":             auditLogEntriesDotGo(proj),
		"accounts_test.go":                 accountsTestDotGo(proj),
		"items.go":                         itemsDotGo(proj),
		"doc.go":                           docDotGo(proj),
		"errors.go":                        errorsDotGo(proj),
		"migrate_test.go":                  migrateTestDotGo(proj),
		"webhooks.go":                      webhooksDotGo(proj),
		"api_clients.go":                   apiClientsDotGo(proj),
		"api_clients_test.go":              apiClientsTestDotGo(proj),
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
