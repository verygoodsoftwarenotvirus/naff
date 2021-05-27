package querier

import (
	"fmt"
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
		"webhooks.go":                      webhooksDotGo(proj),
		"webhooks_test.go":                 webhooksTestDotGo(proj),
		"accounts.go":                      accountsDotGo(proj),
		"api_clients_test.go":              apiClientsTestDotGo(proj),
		"audit_log_entries.go":             auditLogEntriesDotGo(proj),
		"migrate_test.go":                  migrateTestDotGo(proj),
		"users_test.go":                    usersTestDotGo(proj),
		"admin_test.go":                    adminTestDotGo(proj),
		"audit_log_entries_test.go":        auditLogEntriesTestDotGo(proj),
		"migrate.go":                       migrateDotGo(proj),
		"account_user_memberships.go":      accountUserMembershipsDotGo(proj),
		"account_user_memberships_test.go": accountUserMembershipsTestDotGo(proj),
		"admin.go":                         adminDotGo(proj),
		"doc.go":                           docDotGo(proj),
		"querier.go":                       querierDotGo(proj),
		"users.go":                         usersDotGo(proj),
		"accounts_test.go":                 accountsTestDotGo(proj),
		"api_clients.go":                   apiClientsDotGo(proj),
		"errors.go":                        errorsDotGo(proj),
		"querier_test.go":                  querierTestDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
		files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
