package routing

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func routerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("Middleware").Func().Params(jen.Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
			jen.ID("Router").Interface(
				jen.ID("LogRoutes").Params(),
				jen.ID("Handler").Params().Params(jen.Qual("net/http", "Handler")),
				jen.ID("Handle").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "Handler")),
				jen.ID("HandleFunc").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("WithMiddleware").Params(jen.ID("middleware").Op("...").ID("Middleware")).Params(jen.ID("Router")),
				jen.ID("Route").Params(jen.ID("pattern").ID("string"), jen.ID("fn").Func().Params(jen.ID("r").ID("Router"))).Params(jen.ID("Router")),
				jen.ID("Connect").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("Delete").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("Get").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("Head").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("Options").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("Patch").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("Patch").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("Post").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("Put").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("Trace").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")),
				jen.ID("AddRoute").Params(jen.List(jen.ID("method"), jen.ID("path")).ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc"), jen.ID("middleware").Op("...").ID("Middleware")).Params(jen.ID("error")),
			),
			jen.ID("RouteParamManager").Interface(
				jen.ID("BuildRouteParamIDFetcher").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.List(jen.ID("key"), jen.ID("logDescription")).ID("string")).Func().Params(jen.Op("*").Qual("net/http", "Request")).ID("uint64"),
				jen.ID("BuildRouteParamStringIDFetcher").Params(jen.ID("key").ID("string")).Params(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")), jen.ID("string")),
			),
		),
		jen.Line(),
	)

	return code
}
