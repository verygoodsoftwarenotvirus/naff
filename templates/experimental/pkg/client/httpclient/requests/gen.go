package requests

import (
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
		"accounts_test_.go":          accountsTestDotGo(proj),
		"auth.go":                    authDotGo(proj),
		"test_helpers_test_.go":      testHelpersTestDotGo(proj),
		"webhooks_test_.go":          webhooksTestDotGo(proj),
		"auth_test_.go":              authTestDotGo(proj),
		"builder.go":                 builderDotGo(proj),
		"errors.go":                  errorsDotGo(proj),
		"items_test_.go":             itemsTestDotGo(proj),
		"paseto_test_.go":            pasetoTestDotGo(proj),
		"users_test_.go":             usersTestDotGo(proj),
		"admin.go":                   adminDotGo(proj),
		"audit_log_entries_test_.go": auditLogEntriesTestDotGo(proj),
		"users.go":                   usersDotGo(proj),
		"items.go":                   itemsDotGo(proj),
		"paseto.go":                  pasetoDotGo(proj),
		"accounts.go":                accountsDotGo(proj),
		"admin_test_.go":             adminTestDotGo(proj),
		"api_clients.go":             apiClientsDotGo(proj),
		"api_clients_test_.go":       apiClientsTestDotGo(proj),
		"audit_log_entries.go":       auditLogEntriesDotGo(proj),
		"builder_test_.go":           builderTestDotGo(proj),
		"webhooks.go":                webhooksDotGo(proj),
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
