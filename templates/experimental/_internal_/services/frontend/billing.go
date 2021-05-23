package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func billingDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("pricePanicker").Op("=").ID("panicking").Dot("NewProductionPanicker").Call(),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("priceDivisor").Op("=").Lit(0.01),
		jen.ID("arbitraryPriceMax").Op("=").Lit(100000),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NOTE: this function panics when it receives a number > 100,000, as it is not meant to handle those kinds of prices."),
		jen.Line(),
		jen.Func().ID("renderPrice").Params(jen.ID("p").ID("uint32")).Params(jen.ID("string")).Body(
			jen.If(jen.ID("p").Op(">=").ID("arbitraryPriceMax")).Body(
				jen.ID("pricePanicker").Dot("Panic").Call(jen.Lit("price to be rendered is too large!"))),
			jen.Return().Qual("fmt", "Sprintf").Call(
				jen.Lit("$%.2f"),
				jen.ID("float64").Call(jen.ID("p")).Op("*").ID("priceDivisor"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleCheckoutSessionStart").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("checkout session route called")),
			jen.ID("selectedPlan").Op(":=").ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.Lit("plan")),
			jen.If(jen.ID("selectedPlan").Op("==").Lit("")).Body(
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.List(jen.ID("sessionID"), jen.ID("err")).Op(":=").ID("s").Dot("paymentManager").Dot("CreateCheckoutSession").Call(
				jen.ID("ctx"),
				jen.ID("selectedPlan"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("error creating checkout session token"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("logger").Dot("WithValue").Call(
				jen.Lit("sessionID"),
				jen.ID("sessionID"),
			).Dot("Debug").Call(jen.Lit("session id fetched")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleCheckoutSuccess").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("checkout session success route called")),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusTooEarly")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleCheckoutCancel").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("checkout session cancellation route called")),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusTooEarly")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleCheckoutFailure").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("checkout session failure route called")),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusTooEarly")),
		),
		jen.Line(),
	)

	return code
}
