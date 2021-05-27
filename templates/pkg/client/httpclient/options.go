package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func optionsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("option").Func().Params(jen.Op("*").ID("Client")).Params(jen.ID("error")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SetOptions sets a new option on the client."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("SetOptions").Params(jen.ID("opts").Op("...").ID("option")).Params(jen.ID("error")).Body(
			jen.For(jen.List(jen.ID("_"), jen.ID("opt")).Op(":=").Range().ID("opts")).Body(
				jen.If(jen.ID("err").Op(":=").ID("opt").Call(jen.ID("c")), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().ID("err"))),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UsingJSON sets the url on the client."),
		jen.Line(),
		jen.Func().ID("UsingJSON").Params().Params(jen.Func().Params(jen.Op("*").ID("Client")).Params(jen.ID("error"))).Body(
			jen.Return().Func().Params(jen.ID("c").Op("*").ID("Client")).Params(jen.ID("error")).Body(
				jen.List(jen.ID("requestBuilder"), jen.ID("err")).Op(":=").ID("requests").Dot("NewBuilder").Call(
					jen.ID("c").Dot("url"),
					jen.ID("c").Dot("logger"),
					jen.ID("encoding").Dot("ProvideClientEncoder").Call(
						jen.ID("c").Dot("logger"),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().ID("err")),
				jen.ID("c").Dot("requestBuilder").Op("=").ID("requestBuilder"),
				jen.ID("c").Dot("encoder").Op("=").ID("encoding").Dot("ProvideClientEncoder").Call(
					jen.ID("c").Dot("logger"),
					jen.ID("encoding").Dot("ContentTypeJSON"),
				),
				jen.Return().ID("nil"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UsingXML sets the url on the client."),
		jen.Line(),
		jen.Func().ID("UsingXML").Params().Params(jen.Func().Params(jen.Op("*").ID("Client")).Params(jen.ID("error"))).Body(
			jen.Return().Func().Params(jen.ID("c").Op("*").ID("Client")).Params(jen.ID("error")).Body(
				jen.List(jen.ID("requestBuilder"), jen.ID("err")).Op(":=").ID("requests").Dot("NewBuilder").Call(
					jen.ID("c").Dot("url"),
					jen.ID("c").Dot("logger"),
					jen.ID("encoding").Dot("ProvideClientEncoder").Call(
						jen.ID("c").Dot("logger"),
						jen.ID("encoding").Dot("ContentTypeXML"),
					),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().ID("err")),
				jen.ID("c").Dot("requestBuilder").Op("=").ID("requestBuilder"),
				jen.ID("c").Dot("encoder").Op("=").ID("encoding").Dot("ProvideClientEncoder").Call(
					jen.ID("c").Dot("logger"),
					jen.ID("encoding").Dot("ContentTypeXML"),
				),
				jen.Return().ID("nil"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UsingLogger sets the logger on the client."),
		jen.Line(),
		jen.Func().ID("UsingLogger").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.Func().Params(jen.Op("*").ID("Client")).Params(jen.ID("error"))).Body(
			jen.Return().Func().Params(jen.ID("c").Op("*").ID("Client")).Params(jen.ID("error")).Body(
				jen.ID("c").Dot("logger").Op("=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")),
				jen.Return().ID("nil"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UsingDebug sets the debug value on the client."),
		jen.Line(),
		jen.Func().ID("UsingDebug").Params(jen.ID("debug").ID("bool")).Params(jen.Func().Params(jen.Op("*").ID("Client")).Params(jen.ID("error"))).Body(
			jen.Return().Func().Params(jen.ID("c").Op("*").ID("Client")).Params(jen.ID("error")).Body(
				jen.ID("c").Dot("debug").Op("=").ID("debug"),
				jen.If(jen.ID("debug")).Body(
					jen.ID("c").Dot("logger").Dot("SetLevel").Call(jen.ID("logging").Dot("DebugLevel"))),
				jen.Return().ID("nil"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UsingTimeout sets the debug value on the client."),
		jen.Line(),
		jen.Func().ID("UsingTimeout").Params(jen.ID("timeout").Qual("time", "Duration")).Params(jen.Func().Params(jen.Op("*").ID("Client")).Params(jen.ID("error"))).Body(
			jen.Return().Func().Params(jen.ID("c").Op("*").ID("Client")).Params(jen.ID("error")).Body(
				jen.If(jen.ID("timeout").Op("==").Lit(0)).Body(
					jen.ID("timeout").Op("=").ID("defaultTimeout")),
				jen.ID("c").Dot("authedClient").Dot("Timeout").Op("=").ID("timeout"),
				jen.ID("c").Dot("unauthenticatedClient").Dot("Timeout").Op("=").ID("timeout"),
				jen.Return().ID("nil"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UsingCookie sets the authCookie value on the client."),
		jen.Line(),
		jen.Func().ID("UsingCookie").Params(jen.ID("cookie").Op("*").Qual("net/http", "Cookie")).Params(jen.Func().Params(jen.Op("*").ID("Client")).Params(jen.ID("error"))).Body(
			jen.Return().Func().Params(jen.ID("c").Op("*").ID("Client")).Params(jen.ID("error")).Body(
				jen.If(jen.ID("cookie").Op("==").ID("nil")).Body(
					jen.Return().ID("ErrCookieRequired")),
				jen.ID("c").Dot("authMethod").Op("=").ID("cookieAuthMethod"),
				jen.ID("c").Dot("authedClient").Dot("Transport").Op("=").ID("newCookieRoundTripper").Call(
					jen.ID("c"),
					jen.ID("cookie"),
				),
				jen.ID("c").Dot("authedClient").Op("=").ID("buildRetryingClient").Call(
					jen.ID("c").Dot("authedClient"),
					jen.ID("c").Dot("logger"),
					jen.ID("c").Dot("tracer"),
				),
				jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("set client auth cookie")),
				jen.Return().ID("nil"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UsingPASETO sets the authCookie value on the client."),
		jen.Line(),
		jen.Func().ID("UsingPASETO").Params(jen.ID("clientID").ID("string"), jen.ID("secretKey").Index().ID("byte")).Params(jen.Func().Params(jen.Op("*").ID("Client")).Params(jen.ID("error"))).Body(
			jen.Return().Func().Params(jen.ID("c").Op("*").ID("Client")).Params(jen.ID("error")).Body(
				jen.ID("c").Dot("authMethod").Op("=").ID("pasetoAuthMethod"),
				jen.ID("c").Dot("authedClient").Dot("Transport").Op("=").ID("newPASETORoundTripper").Call(
					jen.ID("c"),
					jen.ID("clientID"),
					jen.ID("secretKey"),
				),
				jen.ID("c").Dot("authedClient").Op("=").ID("buildRetryingClient").Call(
					jen.ID("c").Dot("authedClient"),
					jen.ID("c").Dot("logger"),
					jen.ID("c").Dot("tracer"),
				),
				jen.Return().ID("nil"),
			)),
		jen.Line(),
	)

	return code
}
