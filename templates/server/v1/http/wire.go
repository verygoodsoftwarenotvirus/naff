package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildWireVarDeclarations(proj)...)
	code.Add(buildProvideNamespace(proj)...)

	// if proj.EnableNewsman {
	code.Add(buildProvideNewsmanTypeNameManipulationFunc()...)
	// }

	return code
}

func buildWireVarDeclarations(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("Providers is our wire superset of providers this package offers."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(

				jen.ID("paramFetcherProviders"),
				jen.ID("ProvideServer"),
				jen.ID("ProvideNamespace"),
				func() jen.Code {
					if proj.EnableNewsman {
						return jen.ID("ProvideNewsmanTypeNameManipulationFunc")
					}
					return jen.Null()
				}(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideNamespace(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideNamespace provides a namespace."),
		jen.Line(),
		jen.Func().ID("ProvideNamespace").Params().Params(jen.Qual(proj.InternalMetricsV1Package(), "Namespace")).Block(
			jen.Return().ID("serverNamespace"),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideNewsmanTypeNameManipulationFunc() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideNewsmanTypeNameManipulationFunc provides an WebhookIDFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideNewsmanTypeNameManipulationFunc").Params().Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "TypeNameManipulationFunc")).Block(
			jen.Return().Func().Params(jen.ID("s").String()).Params(jen.String()).Block(
				jen.Return().ID("s"),
			),
		),
		jen.Line(),
	}

	return lines
}
