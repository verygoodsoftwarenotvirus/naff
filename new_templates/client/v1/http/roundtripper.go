package client

import jen "github.com/dave/jennifer/jen"

func roundtripperDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(jen.Const().Defs(
		jen.Id("userAgentHeader").Op("=").Lit("User-Agent"),
		jen.Id("userAgent").Op("=").Lit("TODO Service Client"),
	),
		jen.Line(),
	)
	ret.Add(jen.Line())

	ret.Add(jen.Type().Id("defaultRoundTripper").Struct(
		jen.Id("baseTransport").Op("*").Qual("net/http", "Transport"),
	),
		jen.Line(),
	)
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Id("newDefaultRoundTripper").Params().Params(jen.Op("*").Id("defaultRoundTripper")).Block(
			jen.Return(
				jen.Op("&").Id("defaultRoundTripper").Values(jen.Dict{
					jen.Id("baseTransport"): jen.Id("buildDefaultTransport").Call(),
				}),
			),
		),
		jen.Line(),
	)
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(jen.Id("t").Op("*").Id("defaultRoundTripper")).Id("RoundTrip").Params(
			jen.Id("req").Op("*").Qual("net/http", "Request"),
		).Params(
			jen.Op("*").Qual("net/http", "Response"),
			jen.Id("error"),
		).Block(
			jen.Id("req").Dot("Header").Dot("Set").Call(
				jen.Id("userAgentHeader"),
				jen.Id("userAgent"),
			),
			jen.Line(),
			jen.Return().Id("t").Dot("baseTransport").Dot("RoundTrip").Call(
				jen.Id("req"),
			),
		),
	)
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Id("buildDefaultTransport").Params().Params(jen.Op("*").Qual("net/http", "Transport")).Block(
			jen.Return().Op("&").Qual("net/http", "Transport").Values(
				jen.Dict{
					jen.Id("Proxy"): jen.Qual("net/http", "ProxyFromEnvironment"),
					jen.Id("DialContext"): jen.Parens(jen.Op("&").Qual("net", "Dialer").Values(
						jen.Dict{
							jen.Id("Timeout"):   jen.Lit(30).Op("*").Qual("time", "Second"),
							jen.Id("KeepAlive"): jen.Lit(30).Op("*").Qual("time", "Second"),
							jen.Id("DualStack"): jen.Id("true"),
						},
					),
					).Dot("DialContext"),
					jen.Id("MaxIdleConns"):          jen.Lit(100),
					jen.Id("MaxIdleConnsPerHost"):   jen.Lit(100),
					jen.Id("IdleConnTimeout"):       jen.Lit(90).Op("*").Qual("time", "Second"),
					jen.Id("TLSHandshakeTimeout"):   jen.Lit(10).Op("*").Qual("time", "Second"),
					jen.Id("ExpectContinueTimeout"): jen.Lit(1).Op("*").Qual("time", "Second"),
				},
			),
		),
	)
	ret.Add(jen.Line())

	return ret
}
