package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(jen.Const().Defs(
		jen.ID("userAgentHeader").Op("=").Lit("User-Agent"),
		jen.ID("userAgent").Op("=").Lit("TODO Service Client"),
	),
		jen.Line(),
	)
	ret.Add(jen.Line())

	ret.Add(jen.Type().ID("defaultRoundTripper").Struct(
		jen.ID("baseTransport").Op("*").Qual("net/http", "Transport"),
	),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("newDefaultRoundTripper constructs a new http.RoundTripper"),
		jen.Line(),
		jen.Func().ID("newDefaultRoundTripper").Params().Params(jen.Op("*").ID("defaultRoundTripper")).Block(
			jen.Return(
				jen.Op("&").ID("defaultRoundTripper").Valuesln(
					jen.ID("baseTransport").Op(":").ID("buildDefaultTransport").Call(),
				),
			),
		),
		jen.Line(),
	)
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment("RoundTrip implements the http.RoundTripper interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("defaultRoundTripper")).ID("RoundTrip").Params(
			jen.ID("req").Op("*").Qual("net/http", "Request"),
		).Params(
			jen.Op("*").Qual("net/http", "Response"),
			jen.ID("error"),
		).Block(
			jen.ID("req").Dot("Header").Dot("Set").Call(
				jen.ID("userAgentHeader"),
				jen.ID("userAgent"),
			),
			jen.Return().ID("t").Dot("baseTransport").Dot("RoundTrip").Call(
				jen.ID("req"),
			),
		),
	)
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment("buildDefaultTransport constructs a new http.Transport"),
		jen.Line(),
		jen.Func().ID("buildDefaultTransport").Params().Params(jen.Op("*").Qual("net/http", "Transport")).Block(
			jen.Return().Op("&").Qual("net/http", "Transport").Valuesln(
				jen.ID("Proxy").Op(":").Qual("net/http", "ProxyFromEnvironment"),
				jen.ID("DialContext").Op(":").Parens(jen.Op("&").Qual("net", "Dialer").Valuesln(
					jen.ID("Timeout").Op(":").Lit(30).Op("*").Qual("time", "Second"),
					jen.ID("KeepAlive").Op(":").Lit(30).Op("*").Qual("time", "Second"),
					jen.ID("DualStack").Op(":").ID("true"),
				),
				).Dot("DialContext"),
				jen.ID("MaxIdleConns").Op(":").Lit(100),
				jen.ID("MaxIdleConnsPerHost").Op(":").Lit(100),
				jen.ID("IdleConnTimeout").Op(":").Lit(90).Op("*").Qual("time", "Second"),
				jen.ID("TLSHandshakeTimeout").Op(":").Lit(10).Op("*").Qual("time", "Second"),
				jen.ID("ExpectContinueTimeout").Op(":").Lit(1).Op("*").Qual("time", "Second"),
			),
		),
	)
	ret.Add(jen.Line())

	return ret
}
