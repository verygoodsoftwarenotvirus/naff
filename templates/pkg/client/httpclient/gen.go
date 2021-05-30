package httpclient

import (
	_ "embed"
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
	files := map[string]*jen.File{
		"mock_read_closer_test.go":    mockReadCloserTestDotGo(proj),
		"roundtripper_paseto.go":      roundtripperPasetoDotGo(proj),
		"webhooks_test.go":            webhooksTestDotGo(proj),
		"admin_test.go":               adminTestDotGo(proj),
		"audit_log_test.go":           auditLogTestDotGo(proj),
		"helpers_test.go":             helpersTestDotGo(proj),
		"auth.go":                     authDotGo(proj),
		"items_test.go":               itemsTestDotGo(proj),
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
		"items.go":                    itemsDotGo(proj),
		"audit_log.go":                auditLogDotGo(proj),
		"client_test.go":              clientTestDotGo(proj),
		"webhooks.go":                 webhooksDotGo(proj),
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
