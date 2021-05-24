package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metaDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Null().Type().ID("MetaSettings").Struct(
			jen.ID("RunMode").ID("runMode"),
			jen.ID("Debug").ID("bool"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("MetaSettings")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates an MetaSettings struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").ID("MetaSettings")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.Op("&").ID("s"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("s").Dot("RunMode"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						jen.ID("TestingRunMode"),
						jen.ID("DevelopmentRunMode"),
						jen.ID("ProductionRunMode"),
					),
				),
			)),
		jen.Line(),
	)

	return code
}
