package frontend

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "frontend"

	basePackagePath = "internal/services/frontend"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"webhooks_test.go":      webhooksTestDotGo(proj),
		"accounts_test.go":      accountsTestDotGo(proj),
		"base_template_test.go": baseTemplateTestDotGo(proj),
		"config_test.go":        configTestDotGo(proj),
		"http_routes.go":        httpRoutesDotGo(proj),
		"webhooks.go":           webhooksDotGo(proj),
		"time_test.go":          timeTestDotGo(proj),
		"wire_test.go":          wireTestDotGo(proj),
		"accounts.go":           accountsDotGo(proj),
		"api_clients_test.go":   apiClientsTestDotGo(proj),
		"billing.go":            billingDotGo(proj),
		"helpers_test.go":       helpersTestDotGo(proj),
		"static_assets.go":      staticAssetsDotGo(proj),
		"users_test.go":         usersTestDotGo(proj),
		"auth.go":               authDotGo(proj),
		"base_template.go":      baseTemplateDotGo(proj),
		"time.go":               timeDotGo(proj),
		"wire.go":               wireDotGo(proj),
		"users.go":              usersDotGo(proj),
		"auth_test.go":          authTestDotGo(proj),
		"billing_test.go":       billingTestDotGo(proj),
		"i18n.go":               i18NDotGo(proj),
		"items.go":              itemsDotGo(proj),
		"static_assets_test.go": staticAssetsTestDotGo(proj),
		"service.go":            serviceDotGo(proj),
		"http_routes_test.go":   httpRoutesTestDotGo(proj),
		"i18n_test.go":          i18NTestDotGo(proj),
		"items_test.go":         itemsTestDotGo(proj),
		"languages.go":          languagesDotGo(proj),
		"languages_test.go":     languagesTestDotGo(proj),
		"settings_test.go":      settingsTestDotGo(proj),
		"api_clients.go":        apiClientsDotGo(proj),
		"config.go":             configDotGo(proj),
		"helpers.go":            helpersDotGo(proj),
		"service_test.go":       serviceTestDotGo(proj),
		"settings.go":           settingsDotGo(proj),
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
