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
		"admin_test.go":               adminTestDotGo(proj),
		"client_test.go":              clientTestDotGo(proj),
		"webhooks.go":                 webhooksDotGo(proj),
		"options_test.go":             optionsTestDotGo(proj),
		"roundtripper_cookie.go":      roundtripperCookieDotGo(proj),
		"accounts_test.go":            accountsTestDotGo(proj),
		"api_clients_test.go":         apiClientsTestDotGo(proj),
		"audit_log_test.go":           auditLogTestDotGo(proj),
		"items.go":                    itemsDotGo(proj),
		"paseto_test.go":              pasetoTestDotGo(proj),
		"roundtripper_cookie_test.go": roundtripperCookieTestDotGo(proj),
		"roundtripper_paseto.go":      roundtripperPasetoDotGo(proj),
		"admin.go":                    adminDotGo(proj),
		"errors.go":                   errorsDotGo(proj),
		"paseto.go":                   pasetoDotGo(proj),
		"roundtripper_base_test.go":   roundtripperBaseTestDotGo(proj),
		"doc.go":                      docDotGo(proj),
		"helpers_test.go":             helpersTestDotGo(proj),
		"items_test.go":               itemsTestDotGo(proj),
		"options.go":                  optionsDotGo(proj),
		"accounts.go":                 accountsDotGo(proj),
		"auth.go":                     authDotGo(proj),
		"auth_test.go":                authTestDotGo(proj),
		"client.go":                   clientDotGo(proj),
		"roundtripper_base.go":        roundtripperBaseDotGo(proj),
		"roundtripper_paseto_test.go": roundtripperPasetoTestDotGo(proj),
		"api_clients.go":              apiClientsDotGo(proj),
		"helpers.go":                  helpersDotGo(proj),
		"users.go":                    usersDotGo(proj),
		"webhooks_test.go":            webhooksTestDotGo(proj),
		"audit_log.go":                auditLogDotGo(proj),
		"mock_read_closer_test.go":    mockReadCloserTestDotGo(proj),
		"test_helpers_test.go":        testHelpersTestDotGo(proj),
		"users_test.go":               usersTestDotGo(proj),
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
