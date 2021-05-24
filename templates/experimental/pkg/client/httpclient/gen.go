package httpclient

import (
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
		"client.go":                    clientDotGo(proj),
		"doc.go":                       docDotGo(proj),
		"helpers.go":                   helpersDotGo(proj),
		"roundtripper_base_test_.go":   roundtripperBaseTestDotGo(proj),
		"roundtripper_paseto.go":       roundtripperPasetoDotGo(proj),
		"users.go":                     usersDotGo(proj),
		"webhooks_test_.go":            webhooksTestDotGo(proj),
		"helpers_test_.go":             helpersTestDotGo(proj),
		"mock_read_closer_test_.go":    mockReadCloserTestDotGo(proj),
		"options_test_.go":             optionsTestDotGo(proj),
		"roundtripper_cookie.go":       roundtripperCookieDotGo(proj),
		"roundtripper_paseto_test_.go": roundtripperPasetoTestDotGo(proj),
		"webhooks.go":                  webhooksDotGo(proj),
		"accounts_test_.go":            accountsTestDotGo(proj),
		"api_clients_test_.go":         apiClientsTestDotGo(proj),
		"options.go":                   optionsDotGo(proj),
		"paseto.go":                    pasetoDotGo(proj),
		"roundtripper_cookie_test_.go": roundtripperCookieTestDotGo(proj),
		"auth.go":                      authDotGo(proj),
		"roundtripper_base.go":         roundtripperBaseDotGo(proj),
		"users_test_.go":               usersTestDotGo(proj),
		"audit_log_test_.go":           auditLogTestDotGo(proj),
		"client_test_.go":              clientTestDotGo(proj),
		"errors.go":                    errorsDotGo(proj),
		"items_test_.go":               itemsTestDotGo(proj),
		"test_helpers_test_.go":        testHelpersTestDotGo(proj),
		"admin_test_.go":               adminTestDotGo(proj),
		"api_clients.go":               apiClientsDotGo(proj),
		"paseto_test_.go":              pasetoTestDotGo(proj),
		"accounts.go":                  accountsDotGo(proj),
		"admin.go":                     adminDotGo(proj),
		"auth_test_.go":                authTestDotGo(proj),
		"items.go":                     itemsDotGo(proj),
		"audit_log.go":                 auditLogDotGo(proj),
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
