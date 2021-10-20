package frontend

import (
	"bytes"
	_ "embed"
	"fmt"
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
		"http_routes.go":  httpRoutesDotGo(proj),
		"helper_test.go":  helperTestDotGo(proj),
		"wire_test.go":    wireTestDotGo(proj),
		"helpers_test.go": helpersTestDotGo(proj),
		//"static_assets.go":    staticAssetsDotGo(proj),
		"wire.go":             wireDotGo(proj),
		"http_routes_test.go": httpRoutesTestDotGo(proj),
		"config.go":           configDotGo(proj),
		"config_test.go":      configTestDotGo(proj),
		"helpers.go":          helpersDotGo(proj),
	}

	jenFiles := map[string]*jen.File{
		"service.go":       serviceDotGo(proj),
		"service_test.go":  serviceTestDotGo(proj),
		"static_assets.go": staticAssetsDotGo(proj),
	}

	//for _, typ := range proj.DataTypes {
	//	jenFiles[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
	//	jenFiles[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	//}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
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
			jen.IDf("single%sPattern", sn).Assign().Qual("fmt", "Sprintf").Call(
				jen.ID("numericIDPattern"),
				jen.IDf("%sIDURLParamKey", uvn),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Read%sPermission", pn))).
				Dotln("Get").Call(
				jen.Litf("/%s", prn),
				jen.ID("s").Dotf("build%sTableView", pn).Call(jen.True()),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Read%sPermission", pn))).
				Dotln("Get").Call(
				jen.Litf("/dashboard_pages/%s", prn),
				jen.ID("s").Dotf("build%sTableView", pn).Call(jen.False()),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Create%sPermission", pn))).
				Dotln("Get").Call(
				jen.Lit(fmt.Sprintf("/%s/", prn)+"new"),
				jen.ID("s").Dotf("build%sCreatorView", sn).Call(jen.True()),
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
				jen.ID("s").Dotf("build%sCreatorView", sn).Call(jen.False()),
			),
			jen.Newline(),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").
				Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dotf("Update%sPermission", pn))).
				Dotln("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit(fmt.Sprintf("/%s/", prn)+"%s"),
					jen.IDf("single%sPattern", sn),
				),
				jen.ID("s").Dotf("build%sEditorView", sn).Call(jen.True()),
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
				jen.ID("s").Dotf("build%sEditorView", sn).Call(jen.False()),
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

//go:embed wire_test.gotpl
var wireTestTemplate string

func wireTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTestTemplate, nil)
}

//go:embed helpers_test.gotpl
var helpersTestTemplate string

func helpersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTestTemplate, nil)
}

////go:embed static_assets.gotpl
//var staticAssetsTemplate string
//
//func staticAssetsDotGo(proj *models.Project) string {
//	return models.RenderCodeFile(proj, staticAssetsTemplate, nil)
//}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

//go:embed http_routes_test.gotpl
var httpRoutesTestTemplate string

func httpRoutesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpRoutesTestTemplate, nil)
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed helpers.gotpl
var helpersTemplate string

func helpersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTemplate, nil)
}
