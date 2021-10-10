package querier

import (
	_ "embed"
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
	files := map[string]*jen.File{}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = newIterablesDotGo(proj, typ)
		files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	stringFiles := map[string]string{
		"migrate.go":                       migrateDotGo(proj),
		"migrate_test.go":                  migrateTestDotGo(proj),
		"webhooks.go":                      webhooksDotGo(proj),
		"webhooks_test.go":                 webhooksTestDotGo(proj),
		"accounts.go":                      accountsDotGo(proj),
		"api_clients_test.go":              apiClientsTestDotGo(proj),
		"audit_log_entries.go":             auditLogEntriesDotGo(proj),
		"users_test.go":                    usersTestDotGo(proj),
		"admin_test.go":                    adminTestDotGo(proj),
		"audit_log_entries_test.go":        auditLogEntriesTestDotGo(proj),
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

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed migrate.gotpl
var migrateTemplate string

func migrateDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, migrateTemplate, nil)
}

//go:embed migrate_test.gotpl
var migrateTestTemplate string

func migrateTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, migrateTestTemplate, nil)
}

//go:embed webhooks.gotpl
var webhooksTemplate string

func webhooksDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTemplate, nil)
}

//go:embed webhooks_test.gotpl
var webhooksTestTemplate string

func webhooksTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTestTemplate, nil)
}

//go:embed accounts.gotpl
var accountsTemplate string

func accountsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTemplate, nil)
}

//go:embed api_clients_test.gotpl
var apiClientsTestTemplate string

func apiClientsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientsTestTemplate, nil)
}

//go:embed audit_log_entries.gotpl
var auditLogEntriesTemplate string

func auditLogEntriesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, auditLogEntriesTemplate, nil)
}

//go:embed users_test.gotpl
var usersTestTemplate string

func usersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTestTemplate, nil)
}

//go:embed admin_test.gotpl
var adminTestTemplate string

func adminTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminTestTemplate, nil)
}

//go:embed audit_log_entries_test.gotpl
var auditLogEntriesTestTemplate string

func auditLogEntriesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, auditLogEntriesTestTemplate, nil)
}

//go:embed account_user_memberships.gotpl
var accountUserMembershipsTemplate string

func accountUserMembershipsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipsTemplate, nil)
}

//go:embed account_user_memberships_test.gotpl
var accountUserMembershipsTestTemplate string

func accountUserMembershipsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipsTestTemplate, nil)
}

//go:embed admin.gotpl
var adminTemplate string

func adminDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed querier.gotpl
var querierTemplate string

func querierDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, querierTemplate, nil)
}

//go:embed users.gotpl
var usersTemplate string

func usersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTemplate, nil)
}

//go:embed accounts_test.gotpl
var accountsTestTemplate string

func accountsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTestTemplate, nil)
}

//go:embed api_clients.gotpl
var apiClientsTemplate string

func apiClientsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientsTemplate, nil)
}

//go:embed errors.gotpl
var errorsTemplate string

func errorsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, errorsTemplate, nil)
}

//go:embed querier_test.gotpl
var querierTestTemplate string

func querierTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, querierTestTemplate, nil)
}
