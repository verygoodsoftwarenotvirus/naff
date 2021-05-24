package tracing

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("Jaeger").Op("=").Lit("jaeger"),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("Config").Struct(
			jen.ID("Jaeger").Op("*").ID("JaegerConfig"),
			jen.ID("Provider").ID("string"),
			jen.ID("SpanCollectionProbability").ID("float64"),
		).Type().ID("JaegerConfig").Struct(
			jen.ID("CollectorEndpoint").ID("string"),
			jen.ID("ServiceName").ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Initialize provides an instrumentation handler."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Config")).ID("Initialize").Params(jen.ID("l").ID("logging").Dot("Logger")).Params(jen.ID("flushFunc").Params(), jen.ID("err").ID("error")).Body(
			jen.ID("logger").Op(":=").ID("l").Dot("WithValue").Call(
				jen.Lit("tracing_provider"),
				jen.ID("c").Dot("Provider"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("setting tracing provider")),
			jen.Switch(jen.Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("c").Dot("Provider")))).Body(
				jen.Case(jen.ID("Jaeger")).Body(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("setting up jaeger")), jen.Return().ID("c").Dot("SetupJaeger").Call()),
				jen.Default().Body(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("invalid tracing config")), jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("Config")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates the config struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("JaegerConfig")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("c"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("CollectorEndpoint"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("ServiceName"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("JaegerConfig")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates the config struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("JaegerConfig")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("c"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("CollectorEndpoint"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("ServiceName"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	return code
}
