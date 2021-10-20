package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func genDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Var().Defs(
			jen.ID("packageName").Op("=").Lit("frontend"),
			jen.ID("basePackagePath").Op("=").Lit("home/vgsnv/src/gitlab.com/verygoodsoftwarenotvirus/gamut/internal/services/frontend"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("RenderPackage renders the package"),
		jen.Newline(),
		jen.Func().ID("RenderPackage").Params(jen.ID("proj").Op("*").ID("models").Dot("Project")).Params(jen.ID("error")).Body(
			jen.ID("files").Op(":=").Map(jen.ID("string")).Op("*").ID("jen").Dot("File").Valuesln(jen.Lit("service_test.go").Op(":").ID("serviceTestDotGo").Call(jen.ID("proj")), jen.Lit("static_assets.go").Op(":").ID("staticAssetsDotGo").Call(jen.ID("proj")), jen.Lit("config_test.go").Op(":").ID("configTestDotGo").Call(jen.ID("proj")), jen.Lit("gen.go").Op(":").ID("genDotGo").Call(jen.ID("proj")), jen.Lit("helper_test.go").Op(":").ID("helperTestDotGo").Call(jen.ID("proj")), jen.Lit("helpers.go").Op(":").ID("helpersDotGo").Call(jen.ID("proj")), jen.Lit("http_routes.go").Op(":").ID("httpRoutesDotGo").Call(jen.ID("proj")), jen.Lit("wire_test.go").Op(":").ID("wireTestDotGo").Call(jen.ID("proj")), jen.Lit("config.go").Op(":").ID("configDotGo").Call(jen.ID("proj")), jen.Lit("helpers_test.go").Op(":").ID("helpersTestDotGo").Call(jen.ID("proj")), jen.Lit("http_routes_test.go").Op(":").ID("httpRoutesTestDotGo").Call(jen.ID("proj")), jen.Lit("service.go").Op(":").ID("serviceDotGo").Call(jen.ID("proj")), jen.Lit("wire.go").Op(":").ID("wireDotGo").Call(jen.ID("proj"))),
			jen.For(jen.List(jen.ID("path"), jen.ID("file")).Op(":=").Range().ID("files")).Body(
				jen.If(jen.ID("err").Op(":=").ID("utils").Dot("RenderGoFile").Call(
					jen.ID("proj"),
					jen.Qual("path/filepath", "Join").Call(
						jen.ID("basePackagePath"),
						jen.ID("path"),
					),
					jen.ID("file"),
				), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().ID("err"))),
			jen.Return().ID("nil"),
		),
		jen.Newline(),
	)

	return code
}
