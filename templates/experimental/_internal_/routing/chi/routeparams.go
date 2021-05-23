package chi

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func routeparamsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("chiRouteParamManager").Struct(),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewRouteParamManager provides a new RouteParamManager."),
		jen.Line(),
		jen.Func().ID("NewRouteParamManager").Params().Params(jen.ID("routing").Dot("RouteParamManager")).Body(
			jen.Return().Op("&").ID("chiRouteParamManager").Values()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildRouteParamIDFetcher builds a function that fetches a given key from a path with variables added by a router."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("chiRouteParamManager")).ID("BuildRouteParamIDFetcher").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.List(jen.ID("key"), jen.ID("logDescription")).ID("string")).Params(jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("uint64")).Body(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
					jen.ID("chi").Dot("URLParam").Call(
						jen.ID("req"),
						jen.ID("key"),
					),
					jen.Lit(10),
					jen.Lit(64),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil").Op("&&").ID("len").Call(jen.ID("logDescription")).Op(">").Lit(0)).Body(
					jen.ID("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("fetching %s ID from request"),
							jen.ID("logDescription"),
						),
					)),
				jen.Return().ID("u"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildRouteParamStringIDFetcher builds a function that fetches a given key from a path with variables added by a router."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("chiRouteParamManager")).ID("BuildRouteParamStringIDFetcher").Params(jen.ID("key").ID("string")).Params(jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("string")).Body(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
				jen.Return().ID("chi").Dot("URLParam").Call(
					jen.ID("req"),
					jen.ID("key"),
				))),
		jen.Line(),
	)

	return code
}
