package requests

import (
	_ "embed"
	"fmt"
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
	files := map[string]string{
		"accounts.go":          accountsDotGo(proj),
		"accounts_test.go":     accountsTestDotGo(proj),
		"api_clients.go":       apiClientsDotGo(proj),
		"api_clients_test.go":  apiClientsTestDotGo(proj),
		"auth.go":              authDotGo(proj),
		"admin.go":             adminDotGo(proj),
		"builder.go":           builderDotGo(proj),
		"paseto_test.go":       pasetoTestDotGo(proj),
		"test_helpers_test.go": testHelpersTestDotGo(proj),
		"users.go":             usersDotGo(proj),
		"admin_test.go":        adminTestDotGo(proj),
		"auth_test.go":         authTestDotGo(proj),
		"errors.go":            errorsDotGo(proj),
		"paseto.go":            pasetoDotGo(proj),
		"users_test.go":        usersTestDotGo(proj),
		"builder_test.go":      builderTestDotGo(proj),
		"webhooks.go":          webhooksDotGo(proj),
		"webhooks_test.go":     webhooksTestDotGo(proj),
		"websockets.go":        websocketsDotGo(proj),
		"websockets_test.go":   websocketsTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{}

	for _, typ := range proj.DataTypes {
		jenFiles[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
		jenFiles[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed websockets.gotpl
var websocketsTemplate string

func websocketsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, websocketsTemplate, nil)
}

//go:embed websockets_test.gotpl
var websocketsTestTemplate string

func websocketsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, websocketsTestTemplate, nil)
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
