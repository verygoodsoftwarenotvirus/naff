package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(jen.Const().Defs(
		jen.ID("userAgentHeader").Equals().Lit("User-Agent"),
		jen.ID("userAgent").Equals().Lit("TODO Service Client"),
	),
		jen.Line(),
	)
	ret.Add(jen.Line())

	ret.Add(jen.Type().ID("defaultRoundTripper").Struct(
		jen.ID("baseTransport").ParamPointer().Qual("net/http", "Transport"),
	),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("newDefaultRoundTripper constructs a new http.RoundTripper"),
		jen.Line(),
		jen.Func().ID("newDefaultRoundTripper").Params().Params(jen.Op("*").ID("defaultRoundTripper")).Block(
			jen.Return(
				jen.VarPointer().ID("defaultRoundTripper").Valuesln(
					jen.ID("baseTransport").MapAssign().ID("buildDefaultTransport").Call(),
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
			jen.ID("req").ParamPointer().Qual("net/http", "Request"),
		).Params(
			jen.ParamPointer().Qual("net/http", "Response"),
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
		jen.Func().ID("buildDefaultTransport").Params().Params(jen.ParamPointer().Qual("net/http", "Transport")).Block(
			jen.Return().VarPointer().Qual("net/http", "Transport").Valuesln(
				jen.ID("Proxy").MapAssign().Qual("net/http", "ProxyFromEnvironment"),
				jen.ID("DialContext").MapAssign().Parens(jen.VarPointer().Qual("net", "Dialer").Valuesln(
					jen.ID("Timeout").MapAssign().ID("defaultTimeout"),
					jen.ID("KeepAlive").MapAssign().Lit(30).Times().Qual("time", "Second"),
				),
				).Dot("DialContext"),
				jen.ID("MaxIdleConns").MapAssign().Lit(100),
				jen.ID("MaxIdleConnsPerHost").MapAssign().Lit(100),
				jen.ID("TLSHandshakeTimeout").MapAssign().Lit(10).Times().Qual("time", "Second"),
				jen.ID("ExpectContinueTimeout").MapAssign().Lit(2).Times().ID("defaultTimeout"),
				jen.ID("IdleConnTimeout").MapAssign().Lit(3).Times().ID("defaultTimeout"),
			),
		),
	)
	ret.Add(jen.Line())

	return ret
}
