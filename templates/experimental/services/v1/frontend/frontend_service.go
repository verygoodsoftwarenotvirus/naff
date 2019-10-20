package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func frontendServiceDotGo() *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("serviceName").Op("=").Lit("frontend_service"),
	jen.Line(),
	)

	ret.Add(
		jen.Type().ID("Service").Struct(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
		"Logger",
	),
	jen.ID("config").ID("config").Dot(
		"FrontendSettings",
	)),
	jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideFrontendService provides the frontend service to dependency injection"),
		jen.Line(),
		jen.Func().ID("ProvideFrontendService").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
		"Logger",
	),
	jen.ID("cfg").ID("config").Dot(
		"FrontendSettings",
	)).Params(jen.Op("*").ID("Service")).Block(
		jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(
	jen.ID("config").Op(":").ID("cfg"), jen.ID("logger").Op(":").ID("logger").Dot(
			"WithName",
		).Call(jen.ID("serviceName"))),
		jen.Return().ID("svc"),
	),
	jen.Line(),
	)
	return ret
}
