package postgres

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "postgres"

	basePackagePath = "internal/database/querybuilding/postgres"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"postgres.go":                      postgresDotGo(proj),
		"postgres_test.go":                 postgresTestDotGo(proj),
		"items.go":                         itemsDotGo(proj),
		"audit_log_entries_test.go":        auditLogEntriesTestDotGo(proj),
		"generic.go":                       genericDotGo(proj),
		"generic_test.go":                  genericTestDotGo(proj),
		"items_test.go":                    itemsTestDotGo(proj),
		"migrations.go":                    migrationsDotGo(proj),
		"users_test.go":                    usersTestDotGo(proj),
		"webhooks.go":                      webhooksDotGo(proj),
		"accounts.go":                      accountsDotGo(proj),
		"webhooks_test.go":                 webhooksTestDotGo(proj),
		"accounts_test.go":                 accountsTestDotGo(proj),
		"api_clients_test.go":              apiClientsTestDotGo(proj),
		"users.go":                         usersDotGo(proj),
		"account_user_memberships.go":      accountUserMembershipsDotGo(proj),
		"api_clients.go":                   apiClientsDotGo(proj),
		"audit_log_entries.go":             auditLogEntriesDotGo(proj),
		"doc.go":                           docDotGo(proj),
		"wire.go":                          wireDotGo(proj),
		"account_user_memberships_test.go": accountUserMembershipsTestDotGo(proj),
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
