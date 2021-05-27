package integration

import (
	"fmt"
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
		"init.go":             initDotGo(proj),
		"meta_test.go":        metaTestDotGo(proj),
		"users_test.go":       usersTestDotGo(proj),
		"webhooks_test.go":    webhooksTestDotGo(proj),
		"auth_test.go":        authTestDotGo(proj),
		"doc.go":              docDotGo(proj),
		"helpers_test.go":     helpersTestDotGo(proj),
		"suite_test.go":       suiteTestDotGo(proj),
		"accounts_test.go":    accountsTestDotGo(proj),
		"api_clients_test.go": apiClientsTestDotGo(proj),
		"admin_test.go":       adminTestDotGo(proj),
		"audit_log_test.go":   auditLogTestDotGo(proj),
		"frontend_test.go":    frontendTestDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
