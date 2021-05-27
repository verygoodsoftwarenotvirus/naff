package builders

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "postgres"

	basePackagePath = "internal/database/querybuilding"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	for _, db := range proj.EnabledDatabases() {
		files := map[string]*jen.File{
			"account_user_memberships.go":      accountUserMembershipsDotGo(proj),
			"accounts_test.go":                 accountsTestDotGo(proj),
			"api_clients_test.go":              apiClientsTestDotGo(proj),
			"audit_log_entries_test.go":        auditLogEntriesTestDotGo(proj),
			"doc.go":                           docDotGo(proj),
			"generic_test.go":                  genericTestDotGo(proj),
			fmt.Sprintf("%s_test.go", db):      postgresTestDotGo(proj),
			"items.go":                         itemsDotGo(proj),
			"items_test.go":                    itemsTestDotGo(proj),
			"webhooks.go":                      webhooksDotGo(proj),
			"webhooks_test.go":                 webhooksTestDotGo(proj),
			"wire.go":                          wireDotGo(proj),
			"account_user_memberships_test.go": accountUserMembershipsTestDotGo(proj),
			"generic.go":                       genericDotGo(proj),
			"migrations.go":                    migrationsDotGo(proj),
			fmt.Sprintf("%s.go", db):           postgresDotGo(proj),
			"users_test.go":                    usersTestDotGo(proj),
			"accounts.go":                      accountsDotGo(proj),
			"api_clients.go":                   apiClientsDotGo(proj),
			"audit_log_entries.go":             auditLogEntriesDotGo(proj),
			"users.go":                         usersDotGo(proj),
		}

		for path, file := range files {
			if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, db, path), file); err != nil {
				return err
			}
		}
	}

	return nil
}
