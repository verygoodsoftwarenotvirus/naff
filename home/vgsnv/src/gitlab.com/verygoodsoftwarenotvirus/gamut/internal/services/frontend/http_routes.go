package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("SetupRoutes sets up the routes."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("SetupRoutes").Params(jen.ID("router").ID("routing").Dot("Router")).Body(
			jen.ID("router").Op("=").ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("UserAttributionMiddleware")),
			jen.List(jen.ID("staticFileServer"), jen.ID("err")).Op(":=").ID("s").Dot("StaticDir").Call(jen.Lit("/frontend")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("s").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("establishing static file server"),
				)),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/*"),
				jen.ID("staticFileServer"),
			),
		),
		jen.Newline(),
	)

	return code
}
