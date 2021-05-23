package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func routeParamManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("NewRouteParamManager returns a new RouteParamManager."),
		jen.Line(),
		jen.Func().ID("NewRouteParamManager").Params().Params(jen.Op("*").ID("RouteParamManager")).Body(
			jen.Return().Op("&").ID("RouteParamManager").Values()),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("RouteParamManager").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserIDFetcherFromSessionContextData satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("RouteParamManager")).ID("UserIDFetcherFromSessionContextData").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("req")).Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("FetchContextFromRequest satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("RouteParamManager")).ID("FetchContextFromRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("req")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("SessionContextData")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildRouteParamIDFetcher satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("RouteParamManager")).ID("BuildRouteParamIDFetcher").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.List(jen.ID("key"), jen.ID("logDescription")).ID("string")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("uint64")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("logger"),
				jen.ID("key"),
				jen.ID("logDescription"),
			).Dot("Get").Call(jen.Lit(0)).Assert(jen.Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildRouteParamStringIDFetcher satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("RouteParamManager")).ID("BuildRouteParamStringIDFetcher").Params(jen.ID("key").ID("string")).Params(jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("string")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("key")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")))),
		jen.Line(),
	)

	return code
}
