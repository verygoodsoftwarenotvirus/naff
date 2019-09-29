package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func roundtripperDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

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
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().ID("newDefaultRoundTripper").Params().Params(jen.Op("*").ID("defaultRoundTripper")).Block(
			jen.Return(
				jen.Op("&").ID("defaultRoundTripper").Values(jen.Dict{
					jen.ID("baseTransport"): jen.ID("buildDefaultTransport").Call(),
				}),
			),
		),
		jen.Line(),
	)
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(jen.ID(T).Op("*").ID("defaultRoundTripper")).ID("RoundTrip").Params(
			jen.ID("req").Op("*").Qual("net/http", "Request"),
		).Params(
			jen.Op("*").Qual("net/http", "Response"),
			jen.ID("error"),
		).Block(
			jen.ID("req").Dot("Header").Dot("Set").Call(
				jen.ID("userAgentHeader"),
				jen.ID("userAgent"),
			),
			jen.Line(),
			jen.Return().ID(T).Dot("baseTransport").Dot("RoundTrip").Call(
				jen.ID("req"),
			),
		),
	)
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().ID("buildDefaultTransport").Params().Params(jen.Op("*").Qual("net/http", "Transport")).Block(
			jen.Return().Op("&").Qual("net/http", "Transport").Values(
				jen.Dict{
					jen.ID("Proxy"): jen.Qual("net/http", "ProxyFromEnvironment"),
					jen.ID("DialContext"): jen.Parens(jen.Op("&").Qual("net", "Dialer").Values(
						jen.Dict{
							jen.ID("Timeout"):   jen.Lit(30).Op("*").Qual("time", "Second"),
							jen.ID("KeepAlive"): jen.Lit(30).Op("*").Qual("time", "Second"),
							jen.ID("DualStack"): jen.ID("true"),
						},
					),
					).Dot("DialContext"),
					jen.ID("MaxIdleConns"):          jen.Lit(100),
					jen.ID("MaxIdleConnsPerHost"):   jen.Lit(100),
					jen.ID("IdleConnTimeout"):       jen.Lit(90).Op("*").Qual("time", "Second"),
					jen.ID("TLSHandshakeTimeout"):   jen.Lit(10).Op("*").Qual("time", "Second"),
					jen.ID("ExpectContinueTimeout"): jen.Lit(1).Op("*").Qual("time", "Second"),
				},
			),
		),
	)
	ret.Add(jen.Line())

	return ret
}
