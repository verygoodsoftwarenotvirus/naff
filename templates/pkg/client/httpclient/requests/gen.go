package requests

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "pkg/client/httpclient/requests"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"webhooks_test.go":          webhooksTestDotGo(proj),
		"admin.go":                  adminDotGo(proj),
		"api_clients.go":            apiClientsDotGo(proj),
		"builder.go":                builderDotGo(proj),
		"items.go":                  itemsDotGo(proj),     // DELETE ME
		"items_test.go":             itemsTestDotGo(proj), // DELETE ME
		"paseto_test.go":            pasetoTestDotGo(proj),
		"test_helpers_test.go":      testHelpersTestDotGo(proj),
		"users.go":                  usersDotGo(proj),
		"accounts.go":               accountsDotGo(proj),
		"accounts_test.go":          accountsTestDotGo(proj),
		"admin_test.go":             adminTestDotGo(proj),
		"audit_log_entries_test.go": auditLogEntriesTestDotGo(proj),
		"auth_test.go":              authTestDotGo(proj),
		"errors.go":                 errorsDotGo(proj),
		"paseto.go":                 pasetoDotGo(proj),
		"users_test.go":             usersTestDotGo(proj),
		"webhooks.go":               webhooksDotGo(proj),
		"api_clients_test.go":       apiClientsTestDotGo(proj),
		"audit_log_entries.go":      auditLogEntriesDotGo(proj),
		"auth.go":                   authDotGo(proj),
		"builder_test.go":           builderTestDotGo(proj),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName)] = iterablesDotGo(typ)
	//	files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName)] = iterablesTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed webhooks_test.gotpl
var webhooksTestTemplate string

func webhooksTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTestTemplate, nil)
}

//go:embed admin.gotpl
var adminTemplate string

func adminDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminTemplate, nil)
}

//go:embed api_clients.gotpl
var apiClientsTemplate string

func apiClientsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientsTemplate, nil)
}

//go:embed builder.gotpl
var builderTemplate string

func builderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, builderTemplate, nil)
}

//go:embed items.gotpl
var itemsTemplate string

func itemsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, itemsTemplate, nil)
}

//go:embed items_test.gotpl
var itemsTestTemplate string

func itemsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, itemsTestTemplate, nil)
}

//go:embed paseto_test.gotpl
var pasetoTestTemplate string

func pasetoTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, pasetoTestTemplate, nil)
}

//go:embed test_helpers_test.gotpl
var testHelpersTestTemplate string

func testHelpersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, testHelpersTestTemplate, nil)
}

//go:embed users.gotpl
var usersTemplate string

func usersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTemplate, nil)
}

//go:embed accounts.gotpl
var accountsTemplate string

func accountsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTemplate, nil)
}

//go:embed accounts_test.gotpl
var accountsTestTemplate string

func accountsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTestTemplate, nil)
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

//go:embed auth_test.gotpl
var authTestTemplate string

func authTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTestTemplate, nil)
}

//go:embed errors.gotpl
var errorsTemplate string

func errorsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, errorsTemplate, nil)
}

//go:embed paseto.gotpl
var pasetoTemplate string

func pasetoDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, pasetoTemplate, nil)
}

//go:embed users_test.gotpl
var usersTestTemplate string

func usersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTestTemplate, nil)
}

//go:embed webhooks.gotpl
var webhooksTemplate string

func webhooksDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTemplate, nil)
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

//go:embed auth.gotpl
var authTemplate string

func authDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTemplate, nil)
}

//go:embed builder_test.gotpl
var builderTestTemplate string

func builderTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, builderTestTemplate, nil)
}
