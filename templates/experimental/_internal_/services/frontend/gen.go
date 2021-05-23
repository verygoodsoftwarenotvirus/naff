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
		"webhooks_test_.go":      webhooksTestDotGo(proj),
		"accounts_test_.go":      accountsTestDotGo(proj),
		"base_template_test_.go": baseTemplateTestDotGo(proj),
		"config_test_.go":        configTestDotGo(proj),
		"http_routes.go":         httpRoutesDotGo(proj),
		"webhooks.go":            webhooksDotGo(proj),
		"time_test_.go":          timeTestDotGo(proj),
		"wire_test_.go":          wireTestDotGo(proj),
		"accounts.go":            accountsDotGo(proj),
		"api_clients_test_.go":   apiClientsTestDotGo(proj),
		"billing.go":             billingDotGo(proj),
		"helpers_test_.go":       helpersTestDotGo(proj),
		"static_assets.go":       staticAssetsDotGo(proj),
		"users_test_.go":         usersTestDotGo(proj),
		"auth.go":                authDotGo(proj),
		"base_template.go":       baseTemplateDotGo(proj),
		"time.go":                timeDotGo(proj),
		"wire.go":                wireDotGo(proj),
		"users.go":               usersDotGo(proj),
		"auth_test_.go":          authTestDotGo(proj),
		"billing_test_.go":       billingTestDotGo(proj),
		"i18n.go":                i18NDotGo(proj),
		"items.go":               itemsDotGo(proj),
		"static_assets_test_.go": staticAssetsTestDotGo(proj),
		"service.go":             serviceDotGo(proj),
		"http_routes_test_.go":   httpRoutesTestDotGo(proj),
		"i18n_test_.go":          i18NTestDotGo(proj),
		"items_test_.go":         itemsTestDotGo(proj),
		"languages.go":           languagesDotGo(proj),
		"languages_test_.go":     languagesTestDotGo(proj),
		"settings_test_.go":      settingsTestDotGo(proj),
		"api_clients.go":         apiClientsDotGo(proj),
		"config.go":              configDotGo(proj),
		"helpers.go":             helpersDotGo(proj),
		"service_test_.go":       serviceTestDotGo(proj),
		"settings.go":            settingsDotGo(proj),
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
