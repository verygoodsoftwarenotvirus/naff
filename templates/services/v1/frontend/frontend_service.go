package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func frontendServiceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("frontend")

	utils.AddImports(proj, code)

	code.Add(
		jen.Const().Defs(
			jen.ID("serviceName").Equals().Lit("frontend_service"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.Comment("Service is responsible for serving HTML (and other static resources)"),
			jen.ID("Service").Struct(
				constants.LoggerParam(),
				jen.ID("config").Qual(proj.InternalConfigV1Package(), "FrontendSettings"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideFrontendService provides the frontend service to dependency injection."),
		jen.Line(),
		jen.Func().ID("ProvideFrontendService").Params(constants.LoggerParam(), jen.ID("cfg").Qual(proj.InternalConfigV1Package(), "FrontendSettings")).Params(jen.PointerTo().ID("Service")).Block(
			jen.ID("svc").Assign().AddressOf().ID("Service").Valuesln(
				jen.ID("config").MapAssign().ID("cfg"),
				jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.ID("serviceName")),
			),
			jen.Return().ID("svc"),
		),
		jen.Line(),
	)

	return code
}
