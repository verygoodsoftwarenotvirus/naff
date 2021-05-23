package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func frontendServiceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildFrontendConstantDefs()...)
	code.Add(buildFrontendTypeDefs(proj)...)
	code.Add(buildProvideFrontendService(proj)...)

	return code
}

func buildFrontendConstantDefs() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("serviceName").Equals().Lit("frontend_service"),
		),
		jen.Line(),
	}

	return lines
}
func buildFrontendTypeDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("Service is responsible for serving HTML (and other static resources)"),
			jen.ID("Service").Struct(
				proj.LoggerParam(),
				jen.ID("config").Qual(proj.InternalConfigPackage(), "FrontendSettings"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideFrontendService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideFrontendService provides the frontend service to dependency injection."),
		jen.Line(),
		jen.Func().ID("ProvideFrontendService").Params(proj.LoggerParam(), jen.ID("cfg").Qual(proj.InternalConfigPackage(), "FrontendSettings")).Params(jen.PointerTo().ID("Service")).Body(
			jen.ID("svc").Assign().AddressOf().ID("Service").Valuesln(
				jen.ID("config").MapAssign().ID("cfg"),
				jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.ID("serviceName")),
			),
			jen.Return().ID("svc"),
		),
		jen.Line(),
	}

	return lines
}
