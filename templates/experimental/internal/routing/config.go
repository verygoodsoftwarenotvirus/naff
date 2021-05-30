package routing

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ChiProvider").Op("=").Lit("chi"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Config").Struct(
				jen.ID("Provider").ID("string"),
				jen.ID("AllowedOrigins").Index().ID("string"),
				jen.ID("AllowedHeaders").Index().ID("string"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("Config")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a router config struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("Provider"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(jen.ID("ChiProvider")),
				),
			)),
		jen.Line(),
	)

	return code
}
