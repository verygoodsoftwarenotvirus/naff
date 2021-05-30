package integration

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "integration"

	basePackagePath = "tests/integration"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"auth_test.go":        authTestDotGo(proj),
		"suite_test.go":       suiteTestDotGo(proj),
		"webhooks_test.go":    webhooksTestDotGo(proj),
		"api_clients_test.go": apiClientsTestDotGo(proj),
		"items_test.go":       itemsTestDotGo(proj),
		"meta_test.go":        metaTestDotGo(proj),
		"admin_test.go":       adminTestDotGo(proj),
		"audit_log_test.go":   auditLogTestDotGo(proj),
		"frontend_test.go":    frontendTestDotGo(proj),
		"init.go":             initDotGo(proj),
		"users_test.go":       usersTestDotGo(proj),
		"accounts_test.go":    accountsTestDotGo(proj),
		"helpers_test.go":     helpersTestDotGo(proj),
		"doc.go":              docDotGo(proj),
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
