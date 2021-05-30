package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperPasetoDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("pasetoRoundTripper").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("base").Qual("net/http", "RoundTripper"),
				jen.ID("client").Op("*").ID("Client"),
				jen.ID("clientID").ID("string"),
				jen.ID("secretKey").Index().ID("byte"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("newPASETORoundTripper").Params(jen.ID("client").Op("*").ID("Client"), jen.ID("clientID").ID("string"), jen.ID("secretKey").Index().ID("byte")).Params(jen.Op("*").ID("pasetoRoundTripper")).Body(
			jen.Return().Op("&").ID("pasetoRoundTripper").Valuesln(jen.ID("clientID").Op(":").ID("clientID"), jen.ID("secretKey").Op(":").ID("secretKey"), jen.ID("logger").Op(":").ID("client").Dot("logger"), jen.ID("tracer").Op(":").ID("client").Dot("tracer"), jen.ID("base").Op(":").ID("newDefaultRoundTripper").Call(jen.ID("client").Dot("unauthenticatedClient").Dot("Timeout")), jen.ID("client").Op(":").ID("client"))),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("pasetoRoundTripperClient").Op("=").ID("buildRetryingClient").Call(
				jen.Op("&").Qual("net/http", "Client").Valuesln(jen.ID("Timeout").Op(":").ID("defaultTimeout")),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
				jen.ID("tracing").Dot("NewTracer").Call(jen.Lit("PASETO_roundtripper")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RoundTrip authorizes and authenticates the request with a PASETO."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("pasetoRoundTripper")).ID("RoundTrip").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual("net/http", "Response"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("t").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("reqBodyClosed").Op(":=").ID("false"),
			jen.ID("logger").Op(":=").ID("t").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.If(jen.ID("req").Dot("Body").Op("!=").ID("nil")).Body(
				jen.Defer().Func().Params().Body(
					jen.If(jen.Op("!").ID("reqBodyClosed")).Body(
						jen.If(jen.ID("err").Op(":=").ID("req").Dot("Body").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
							jen.ID("observability").Dot("AcknowledgeError").Call(
								jen.ID("err"),
								jen.ID("logger"),
								jen.ID("span"),
								jen.Lit("closing response body"),
							)))).Call()),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("t").Dot("client").Dot("fetchAuthTokenForAPIClient").Call(
				jen.ID("ctx"),
				jen.ID("pasetoRoundTripperClient"),
				jen.ID("t").Dot("clientID"),
				jen.ID("t").Dot("secretKey"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching prerequisite PASETO"),
				))),
			jen.ID("reqBodyClosed").Op("=").ID("true"),
			jen.ID("req").Dot("Header").Dot("Add").Call(
				jen.Lit("Authorization"),
				jen.ID("token"),
			),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("t").Dot("base").Dot("RoundTrip").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing PASETO-authorized request"),
				))),
			jen.Return().List(jen.ID("res"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
