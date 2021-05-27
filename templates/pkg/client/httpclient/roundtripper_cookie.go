package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperCookieDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("cookieRoundtripper").Struct(
				jen.ID("cookie").Op("*").Qual("net/http", "Cookie"),
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("base").Qual("net/http", "RoundTripper"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("newCookieRoundTripper").Params(jen.ID("client").Op("*").ID("Client"), jen.ID("cookie").Op("*").Qual("net/http", "Cookie")).Params(jen.Op("*").ID("cookieRoundtripper")).Body(
			jen.Return().Op("&").ID("cookieRoundtripper").Valuesln(jen.ID("cookie").Op(":").ID("cookie"), jen.ID("logger").Op(":").ID("client").Dot("logger"), jen.ID("tracer").Op(":").ID("client").Dot("tracer"), jen.ID("base").Op(":").ID("newDefaultRoundTripper").Call(jen.ID("client").Dot("unauthenticatedClient").Dot("Timeout")))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RoundTrip authorizes and authenticates the request with a cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("cookieRoundtripper")).ID("RoundTrip").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual("net/http", "Response"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("t").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("reqBodyClosed").Op(":=").ID("false"),
			jen.ID("logger").Op(":=").ID("t").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.If(jen.ID("req").Dot("Body").Op("!=").ID("nil")).Body(
				jen.Defer().Func().Params().Body(
					jen.If(jen.Op("!").ID("reqBodyClosed")).Body(
						jen.If(jen.ID("err").Op(":=").ID("req").Dot("Body").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
							jen.ID("observability").Dot("AcknowledgeError").Call(
								jen.ID("err"),
								jen.ID("t").Dot("logger"),
								jen.ID("span"),
								jen.Lit("closing response body"),
							)))).Call()),
			jen.If(jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("req").Dot("Cookie").Call(jen.ID("t").Dot("cookie").Dot("Name")), jen.ID("c").Op("==").ID("nil").Op("||").ID("err").Op("!=").ID("nil")).Body(
				jen.ID("req").Dot("AddCookie").Call(jen.ID("t").Dot("cookie"))),
			jen.ID("reqBodyClosed").Op("=").ID("true"),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("t").Dot("base").Dot("RoundTrip").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing RoundTrip"),
				))),
			jen.If(jen.ID("responseCookies").Op(":=").ID("res").Dot("Cookies").Call(), jen.ID("len").Call(jen.ID("responseCookies")).Op(">=").Lit(1)).Body(
				jen.ID("t").Dot("cookie").Op("=").ID("responseCookies").Index(jen.Lit(0))),
			jen.Return().List(jen.ID("res"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
