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
		"api_clients_test.go":       apiClientsTestDotGo(proj),
		"paseto_test.go":            pasetoTestDotGo(proj),
		"accounts.go":               accountsDotGo(proj),
		"accounts_test.go":          accountsTestDotGo(proj),
		"builder.go":                builderDotGo(proj),
		"builder_test.go":           builderTestDotGo(proj),
		"errors.go":                 errorsDotGo(proj),
		"items_test.go":             itemsTestDotGo(proj),
		"test_helpers_test.go":      testHelpersTestDotGo(proj),
		"webhooks_test.go":          webhooksTestDotGo(proj),
		"admin_test.go":             adminTestDotGo(proj),
		"auth_test.go":              authTestDotGo(proj),
		"auth.go":                   authDotGo(proj),
		"items.go":                  itemsDotGo(proj),
		"paseto.go":                 pasetoDotGo(proj),
		"webhooks.go":               webhooksDotGo(proj),
		"admin.go":                  adminDotGo(proj),
		"audit_log_entries.go":      auditLogEntriesDotGo(proj),
		"users.go":                  usersDotGo(proj),
		"users_test.go":             usersTestDotGo(proj),
		"api_clients.go":            apiClientsDotGo(proj),
		"audit_log_entries_test.go": auditLogEntriesTestDotGo(proj),
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
