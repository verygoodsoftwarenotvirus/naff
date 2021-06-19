package frontend

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
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
	files := map[string]string{
		"webhooks_test.go":      webhooksTestDotGo(proj),
		"accounts_test.go":      accountsTestDotGo(proj),
		"base_template_test.go": baseTemplateTestDotGo(proj),
		"config_test.go":        configTestDotGo(proj),
		"http_routes.go":        httpRoutesDotGo(proj),
		"webhooks.go":           webhooksDotGo(proj),
		"time_test.go":          timeTestDotGo(proj),
		"form_parsers.go":       formParsersDotGo(proj),
		"form_parsers_test.go":  formParsersTestDotGo(proj),
		"helper_test.go":        helperTestDotGo(proj),
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
		"static_assets_test.go": staticAssetsTestDotGo(proj),
		"http_routes_test.go":   httpRoutesTestDotGo(proj),
		"i18n_test.go":          i18NTestDotGo(proj),
		"languages.go":          languagesDotGo(proj),
		"languages_test.go":     languagesTestDotGo(proj),
		"settings_test.go":      settingsTestDotGo(proj),
		"api_clients.go":        apiClientsDotGo(proj),
		"config.go":             configDotGo(proj),
		"helpers.go":            helpersDotGo(proj),
		"settings.go":           settingsDotGo(proj),
	}

	staticFiles := map[string]string{
		"templates/partials/auth/login.gotpl":                loginAuthPartial(),
		"templates/partials/auth/register.gotpl":             registrationAuthPartial(),
		"templates/partials/auth/registration_success.gotpl": registrationSuccessAuthPartial(),
		"templates/partials/settings/account_settings.gotpl": accountSettingsPartial(),
		"templates/partials/settings/admin_settings.gotpl":   adminSettingsPartial(),
		"templates/partials/settings/user_settings.gotpl":    userSettingsPartial(),
		"templates/base_template.gotpl":                      baseTemplate(),
		"assets/favicon.svg":                                 favicon(),
		"translations/en.toml":                               englishTranslationsToml(),
	}

	jenFiles := map[string]*jen.File{
		"service.go":      serviceDotGo(proj),
		"service_test.go": serviceTestDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		jenFiles[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
		jenFiles[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	for path, file := range staticFiles {
		fp := utils.BuildTemplatePath(proj.OutputPath, filepath.Join(basePackagePath, path))
		dirToMake := filepath.Dir(fp)

		e := os.MkdirAll(dirToMake, 0777)
		_ = e
		if err := os.WriteFile(fp, []byte(file), 0644); err != nil {
			return err
		}
	}

	return nil
}

//go:embed helper_test.gotpl
var helperTestTemplate string

func helperTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helperTestTemplate, nil)
}

//go:embed form_parsers_test.gotpl
var formParsersTestTemplate string

func formParsersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, formParsersTestTemplate, nil)
}

//go:embed form_parsers.gotpl
var formParsersTemplate string

func formParsersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, formParsersTemplate, nil)
}

//go:embed webhooks_test.gotpl
var webhooksTestTemplate string

func webhooksTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTestTemplate, nil)
}

//go:embed accounts_test.gotpl
var accountsTestTemplate string

func accountsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTestTemplate, nil)
}

//go:embed base_template_test.gotpl
var baseTemplateTestTemplate string

func baseTemplateTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, baseTemplateTestTemplate, nil)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed http_routes.gotpl
var httpRoutesTemplate string

func httpRoutesDotGo(proj *models.Project) string {
	routes := []jen.Code{}
	for _, typ := range proj.DataTypes {
		sn := typ.Name.Singular()
		pn := typ.Name.Plural()
		prn := typ.Name.PluralRouteName()
		uvn := typ.Name.UnexportedVarName()

		routes = append(routes,
			jen.Newline(),
			jen.Newline(),
			jen.IDf("single%sPattern", sn).Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("numericIDPattern"),
				jen.IDf("%sIDURLParamKey", uvn),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Read%sPermission", pn))).
				Dotln("Get").Call(
				jen.Litf("/%s", prn),
				jen.ID("s").Dotf("build%sTableView", pn).Call(jen.ID("true")),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Read%sPermission", pn))).
				Dotln("Get").Call(
				jen.Litf("/dashboard_pages/%s", prn),
				jen.ID("s").Dotf("build%sTableView", pn).Call(jen.ID("false")),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Create%sPermission", pn))).
				Dotln("Get").Call(
				jen.Lit(fmt.Sprintf("/%s/", prn)+"new"),
				jen.ID("s").Dotf("build%sCreatorView", sn).Call(jen.ID("true")),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Create%sPermission", pn))).
				Dotln("Post").Call(
				jen.Litf("/%s/new/submit", prn),
				jen.ID("s").Dotf("handle%sCreationRequest", sn),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Archive%sPermission", pn))).
				Dotln("Delete").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit(fmt.Sprintf("/dashboard_pages/%s/", prn)+"%s"),
					jen.IDf("single%sPattern", sn),
				),
				jen.ID("s").Dotf("handle%sArchiveRequest", sn),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Archive%sPermission", pn))).
				Dotln("Get").Call(
				jen.Litf("/dashboard_pages/%s/new", prn),
				jen.ID("s").Dotf("build%sCreatorView", sn).Call(jen.ID("false")),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Update%sPermission", pn))).
				Dotln("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit(fmt.Sprintf("/%s/", prn)+"%s"),
					jen.IDf("single%sPattern", sn),
				),
				jen.ID("s").Dotf("build%sEditorView", sn).Call(jen.ID("true")),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Update%sPermission", pn))).
				Dotln("Put").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit(fmt.Sprintf("/dashboard_pages/%s/", prn)+"%s"),
					jen.IDf("single%sPattern", sn),
				),
				jen.ID("s").Dotf("handle%sUpdateRequest", sn),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Update%sPermission", pn))).
				Dotln("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit(fmt.Sprintf("/dashboard_pages/%s/", prn)+"%s"),
					jen.IDf("single%sPattern", sn),
				),
				jen.ID("s").Dotf("build%sEditorView", sn).Call(jen.ID("false")),
			),
		)
	}

	var b bytes.Buffer
	if err := jen.Null().Add(routes...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	generated := map[string]string{
		"typeRoutes": b.String(),
	}

	return models.RenderCodeFile(proj, httpRoutesTemplate, generated)
}

//go:embed webhooks.gotpl
var webhooksTemplate string

func webhooksDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTemplate, nil)
}

//go:embed time_test.gotpl
var timeTestTemplate string

func timeTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, timeTestTemplate, nil)
}

//go:embed wire_test.gotpl
var wireTestTemplate string

func wireTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTestTemplate, nil)
}

//go:embed accounts.gotpl
var accountsTemplate string

func accountsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTemplate, nil)
}

//go:embed api_clients_test.gotpl
var apiClientsTestTemplate string

func apiClientsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientsTestTemplate, nil)
}

//go:embed billing.gotpl
var billingTemplate string

func billingDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, billingTemplate, nil)
}

//go:embed helpers_test.gotpl
var helpersTestTemplate string

func helpersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTestTemplate, nil)
}

//go:embed static_assets.gotpl
var staticAssetsTemplate string

func staticAssetsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, staticAssetsTemplate, nil)
}

//go:embed users_test.gotpl
var usersTestTemplate string

func usersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTestTemplate, nil)
}

//go:embed auth.gotpl
var authTemplate string

func authDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTemplate, nil)
}

//go:embed base_template.gotpl
var baseTemplateTemplate string

func baseTemplateDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, baseTemplateTemplate, nil)
}

//go:embed time.gotpl
var timeTemplate string

func timeDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, timeTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

//go:embed users.gotpl
var usersTemplate string

func usersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTemplate, nil)
}

//go:embed auth_test.gotpl
var authTestTemplate string

func authTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTestTemplate, nil)
}

//go:embed billing_test.gotpl
var billingTestTemplate string

func billingTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, billingTestTemplate, nil)
}

//go:embed i18n.gotpl
var i18NTemplate string

func i18NDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, i18NTemplate, nil)
}

//go:embed static_assets_test.gotpl
var staticAssetsTestTemplate string

func staticAssetsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, staticAssetsTestTemplate, nil)
}

//go:embed http_routes_test.gotpl
var httpRoutesTestTemplate string

func httpRoutesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpRoutesTestTemplate, nil)
}

//go:embed i18n_test.gotpl
var i18NTestTemplate string

func i18NTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, i18NTestTemplate, nil)
}

//go:embed languages.gotpl
var languagesTemplate string

func languagesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, languagesTemplate, nil)
}

//go:embed languages_test.gotpl
var languagesTestTemplate string

func languagesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, languagesTestTemplate, nil)
}

//go:embed settings_test.gotpl
var settingsTestTemplate string

func settingsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, settingsTestTemplate, nil)
}

//go:embed api_clients.gotpl
var apiClientsTemplate string

func apiClientsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientsTemplate, nil)
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed helpers.gotpl
var helpersTemplate string

func helpersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTemplate, nil)
}

//go:embed settings.gotpl
var settingsTemplate string

func settingsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, settingsTemplate, nil)
}

//
////go:embed service.gotpl
//var serviceTemplate string
//
//func serviceDotGo(proj *models.Project) string {
//	return models.RenderCodeFile(proj, serviceTemplate, nil)
//}
//
////go:embed service_test.gotpl
//var serviceTestTemplate string
//
//func serviceTestDotGo(proj *models.Project) string {
//	return models.RenderCodeFile(proj, serviceTestTemplate, nil)
//}
