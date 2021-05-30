package frontend

import (
	_ "embed"
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
		"languages_test.go":     languagesTestDotGo(proj),
		"settings.go":           settingsDotGo(proj),
		"time.go":               timeDotGo(proj),
		"http_routes_test.go":   httpRoutesTestDotGo(proj),
		"items_test.go":         itemsTestDotGo(proj),
		"billing_test.go":       billingTestDotGo(proj),
		"items.go":              itemsDotGo(proj),
		"languages.go":          languagesDotGo(proj),
		"settings_test.go":      settingsTestDotGo(proj),
		"users_test.go":         usersTestDotGo(proj),
		"accounts_test.go":      accountsTestDotGo(proj),
		"api_clients_test.go":   apiClientsTestDotGo(proj),
		"wire_test.go":          wireTestDotGo(proj),
		"time_test.go":          timeTestDotGo(proj),
		"users.go":              usersDotGo(proj),
		"i18n_test.go":          i18NTestDotGo(proj),
		"webhooks_test.go":      webhooksTestDotGo(proj),
		"http_routes.go":        httpRoutesDotGo(proj),
		"i18n.go":               i18NDotGo(proj),
		"base_template_test.go": baseTemplateTestDotGo(proj),
		"static_assets.go":      staticAssetsDotGo(proj),
		"wire.go":               wireDotGo(proj),
		"billing.go":            billingDotGo(proj),
		"config.go":             configDotGo(proj),
		"base_template.go":      baseTemplateDotGo(proj),
		"config_test.go":        configTestDotGo(proj),
		"helpers.go":            helpersDotGo(proj),
		"helpers_test.go":       helpersTestDotGo(proj),
		"service_test.go":       serviceTestDotGo(proj),
		"accounts.go":           accountsDotGo(proj),
		"auth.go":               authDotGo(proj),
		"service.go":            serviceDotGo(proj),
		"static_assets_test.go": staticAssetsTestDotGo(proj),
		"webhooks.go":           webhooksDotGo(proj),
		"api_clients.go":        apiClientsDotGo(proj),
		"auth_test.go":          authTestDotGo(proj),
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
