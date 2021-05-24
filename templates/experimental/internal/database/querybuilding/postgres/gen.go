package postgres

import (
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
		"accounts.go":                       accountsDotGo(proj),
		"api_clients.go":                    apiClientsDotGo(proj),
		"items.go":                          itemsDotGo(proj),
		"postgres.go":                       postgresDotGo(proj),
		"accounts_test_.go":                 accountsTestDotGo(proj),
		"postgres_test_.go":                 postgresTestDotGo(proj),
		"users_test_.go":                    usersTestDotGo(proj),
		"webhooks.go":                       webhooksDotGo(proj),
		"wire.go":                           wireDotGo(proj),
		"audit_log_entries.go":              auditLogEntriesDotGo(proj),
		"doc.go":                            docDotGo(proj),
		"generic.go":                        genericDotGo(proj),
		"items_test_.go":                    itemsTestDotGo(proj),
		"migrations.go":                     migrationsDotGo(proj),
		"users.go":                          usersDotGo(proj),
		"webhooks_test_.go":                 webhooksTestDotGo(proj),
		"account_user_memberships.go":       accountUserMembershipsDotGo(proj),
		"account_user_memberships_test_.go": accountUserMembershipsTestDotGo(proj),
		"api_clients_test_.go":              apiClientsTestDotGo(proj),
		"audit_log_entries_test_.go":        auditLogEntriesTestDotGo(proj),
		"generic_test_.go":                  genericTestDotGo(proj),
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
