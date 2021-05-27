package logging

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ProviderZerolog").Op("=").Lit("zerolog"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Config").Struct(
				jen.ID("Name").ID("string"),
				jen.ID("Level").ID("Level"),
				jen.ID("Provider").ID("string"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideLogger builds a Logger according to the provided config."),
		jen.Line(),
		jen.Func().ID("ProvideLogger").Params(jen.ID("cfg").ID("Config")).Params(jen.ID("Logger")).Body(
			jen.Var().Defs(
				jen.ID("l").ID("Logger"),
			),
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				jen.Case(jen.ID("ProviderZerolog")).Body(
					jen.ID("l").Op("=").ID("NewZerologLogger").Call()),
				jen.Default().Body(
					jen.ID("l").Op("=").ID("NewNoopLogger").Call()),
			),
			jen.ID("l").Dot("SetLevel").Call(jen.ID("cfg").Dot("Level")),
			jen.ID("l").Op("=").ID("l").Dot("WithName").Call(jen.ID("cfg").Dot("Name")),
			jen.Return().ID("l"),
		),
		jen.Line(),
	)

	return code
}
