package requests

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "requests"

	basePackagePath = "pkg/client/httpclient/requests"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"webhooks_test.go":          webhooksTestDotGo(proj),
		"admin.go":                  adminDotGo(proj),
		"api_clients.go":            apiClientsDotGo(proj),
		"builder.go":                builderDotGo(proj),
		"items_test.go":             itemsTestDotGo(proj),
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
		"items.go":                  itemsDotGo(proj),
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
