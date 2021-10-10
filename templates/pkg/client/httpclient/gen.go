package httpclient

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "httpclient"

	basePackagePath = "pkg/client/httpclient"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"mock_read_closer_test.go":    mockReadCloserTestDotGo(proj),
		"roundtripper_paseto.go":      roundtripperPasetoDotGo(proj),
		"webhooks_test.go":            webhooksTestDotGo(proj),
		"admin_test.go":               adminTestDotGo(proj),
		"helpers_test.go":             helpersTestDotGo(proj),
		"auth.go":                     authDotGo(proj),
		"roundtripper_cookie.go":      roundtripperCookieDotGo(proj),
		"users.go":                    usersDotGo(proj),
		"users_test.go":               usersTestDotGo(proj),
		"options.go":                  optionsDotGo(proj),
		"options_test.go":             optionsTestDotGo(proj),
		"paseto_test.go":              pasetoTestDotGo(proj),
		"api_clients.go":              apiClientsDotGo(proj),
		"client.go":                   clientDotGo(proj),
		"roundtripper_paseto_test.go": roundtripperPasetoTestDotGo(proj),
		"test_helpers_test.go":        testHelpersTestDotGo(proj),
		"accounts_test.go":            accountsTestDotGo(proj),
		"api_clients_test.go":         apiClientsTestDotGo(proj),
		"roundtripper_base.go":        roundtripperBaseDotGo(proj),
		"doc.go":                      docDotGo(proj),
		"roundtripper_cookie_test.go": roundtripperCookieTestDotGo(proj),
		"accounts.go":                 accountsDotGo(proj),
		"admin.go":                    adminDotGo(proj),
		"auth_test.go":                authTestDotGo(proj),
		"paseto.go":                   pasetoDotGo(proj),
		"roundtripper_base_test.go":   roundtripperBaseTestDotGo(proj),
		"errors.go":                   errorsDotGo(proj),
		"helpers.go":                  helpersDotGo(proj),
		"client_test.go":              clientTestDotGo(proj),
		"webhooks.go":                 webhooksDotGo(proj),
		"websockets.go":               websocketsDotGo(proj),
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

//go:embed mock_read_closer_test.gotpl
var mockReadCloserTestTemplate string

func mockReadCloserTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockReadCloserTestTemplate, nil)
}

//go:embed roundtripper_paseto.gotpl
var roundtripperPasetoTemplate string

func roundtripperPasetoDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, roundtripperPasetoTemplate, nil)
}

//go:embed webhooks_test.gotpl
var webhooksTestTemplate string

func webhooksTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTestTemplate, nil)
}

//go:embed admin_test.gotpl
var adminTestTemplate string

func adminTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminTestTemplate, nil)
}

//go:embed helpers_test.gotpl
var helpersTestTemplate string

func helpersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTestTemplate, nil)
}

//go:embed auth.gotpl
var authTemplate string

func authDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTemplate, nil)
}

//go:embed roundtripper_cookie.gotpl
var roundtripperCookieTemplate string

func roundtripperCookieDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, roundtripperCookieTemplate, nil)
}

//go:embed users.gotpl
var usersTemplate string

func usersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTemplate, nil)
}

//go:embed users_test.gotpl
var usersTestTemplate string

func usersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTestTemplate, nil)
}

//go:embed options.gotpl
var optionsTemplate string

func optionsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, optionsTemplate, nil)
}

//go:embed options_test.gotpl
var optionsTestTemplate string

func optionsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, optionsTestTemplate, nil)
}

//go:embed paseto_test.gotpl
var pasetoTestTemplate string

func pasetoTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, pasetoTestTemplate, nil)
}

//go:embed api_clients.gotpl
var apiClientsTemplate string

func apiClientsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientsTemplate, nil)
}

//go:embed client.gotpl
var clientTemplate string

func clientDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, clientTemplate, nil)
}

//go:embed roundtripper_paseto_test.gotpl
var roundtripperPasetoTestTemplate string

func roundtripperPasetoTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, roundtripperPasetoTestTemplate, nil)
}

//go:embed test_helpers_test.gotpl
var testHelpersTestTemplate string

func testHelpersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, testHelpersTestTemplate, nil)
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

//go:embed roundtripper_base.gotpl
var roundtripperBaseTemplate string

func roundtripperBaseDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, roundtripperBaseTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed roundtripper_cookie_test.gotpl
var roundtripperCookieTestTemplate string

func roundtripperCookieTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, roundtripperCookieTestTemplate, nil)
}

//go:embed accounts.gotpl
var accountsTemplate string

func accountsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTemplate, nil)
}

//go:embed admin.gotpl
var adminTemplate string

func adminDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminTemplate, nil)
}

//go:embed auth_test.gotpl
var authTestTemplate string

func authTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTestTemplate, nil)
}

//go:embed paseto.gotpl
var pasetoTemplate string

func pasetoDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, pasetoTemplate, nil)
}

//go:embed roundtripper_base_test.gotpl
var roundtripperBaseTestTemplate string

func roundtripperBaseTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, roundtripperBaseTestTemplate, nil)
}

//go:embed errors.gotpl
var errorsTemplate string

func errorsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, errorsTemplate, nil)
}

//go:embed helpers.gotpl
var helpersTemplate string

func helpersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTemplate, nil)
}

//go:embed client_test.gotpl
var clientTestTemplate string

func clientTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, clientTestTemplate, nil)
}

//go:embed webhooks.gotpl
var webhooksTemplate string

func webhooksDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTemplate, nil)
}

//go:embed websockets.gotpl
var websocketsTemplate string

func websocketsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, websocketsTemplate, nil)
}
