package capitalism

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("StripeProvider").Op("=").Lit("stripe"),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Config").Struct(
				jen.ID("Stripe").Op("*").ID("StripeConfig"),
				jen.ID("Provider").ID("string"),
				jen.ID("Enabled").ID("bool"),
			),
			jen.ID("StripeConfig").Struct(
				jen.ID("APIKey").ID("string"),
				jen.ID("SuccessURL").ID("string"),
				jen.ID("CancelURL").ID("string"),
				jen.ID("WebhookSecret").ID("string"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("Config")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a StripeConfig struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("StripeConfig")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("APIKey"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("StripeConfig")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a StripeConfig struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("StripeConfig")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("APIKey"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	return code
}
