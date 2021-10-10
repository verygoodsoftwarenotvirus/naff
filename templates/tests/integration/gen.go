package integration

import (
	_ "embed"
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
	files := map[string]string{
		"init.go":             initDotGo(proj),
		"meta_test.go":        metaTestDotGo(proj),
		"users_test.go":       usersTestDotGo(proj),
		"webhooks_test.go":    webhooksTestDotGo(proj),
		"auth_test.go":        authTestDotGo(proj),
		"doc.go":              docDotGo(proj),
		"helpers.go":          helpersDotGo(proj),
		"suite_test.go":       suiteTestDotGo(proj),
		"accounts_test.go":    accountsTestDotGo(proj),
		"api_clients_test.go": apiClientsTestDotGo(proj),
		"admin_test.go":       adminTestDotGo(proj),
		"frontend_test.go":    frontendTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{}
	for _, typ := range proj.DataTypes {
		jenFiles[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed init.gotpl
var initTemplate string

func initDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, initTemplate, nil)
}

//go:embed meta_test.gotpl
var metaTestTemplate string

func metaTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, metaTestTemplate, nil)
}

//go:embed users_test.gotpl
var usersTestTemplate string

func usersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTestTemplate, nil)
}

//go:embed webhooks_test.gotpl
var webhooksTestTemplate string

func webhooksTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTestTemplate, nil)
}

//go:embed auth_test.gotpl
var authTestTemplate string

func authTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed helpers.gotpl
var helpersTestTemplate string

func helpersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTestTemplate, nil)
}

//go:embed suite_test.gotpl
var suiteTestTemplate string

func suiteTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, suiteTestTemplate, nil)
}

//go:embed accounts_test.gotpl
var accountsTestTemplate string

func accountsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTestTemplate, nil)
}

//go:embed api_clients_test.gotpl
var apiClientsTestTemplate string

func apiClientsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientsTestTemplate, nil)
}

//go:embed admin_test.gotpl
var adminTestTemplate string

func adminTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminTestTemplate, nil)
}

//go:embed frontend_test.gotpl
var frontendTestTemplate string

func frontendTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, frontendTestTemplate, nil)
}
