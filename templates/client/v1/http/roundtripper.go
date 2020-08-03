package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildRoundtripperConstDecls()...)
	code.Add(buildDefaultRoundTripper()...)
	code.Add(buildNewDefaultRoundTripper()...)
	code.Add(buildRoundTrip()...)
	code.Add(buildBuildDefaultTransport()...)

	return code
}

func buildRoundtripperConstDecls() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("userAgentHeader").Equals().Lit("User-Agent"),
			jen.ID("userAgent").Equals().Lit("TODO Service Client"),
		),
		jen.Line(),
	}

	return lines
}

func buildDefaultRoundTripper() []jen.Code {
	lines := []jen.Code{
		jen.Type().ID("defaultRoundTripper").Struct(
			jen.ID("baseTransport").PointerTo().Qual("net/http", "Transport"),
		),
		jen.Line(),
	}

	return lines
}

func buildNewDefaultRoundTripper() []jen.Code {
	lines := []jen.Code{
		jen.Comment("newDefaultRoundTripper constructs a new http.RoundTripper."),
		jen.Line(),
		jen.Func().ID("newDefaultRoundTripper").Params().Params(jen.PointerTo().ID("defaultRoundTripper")).Body(
			jen.Return(
				jen.AddressOf().ID("defaultRoundTripper").Valuesln(
					jen.ID("baseTransport").MapAssign().ID("buildDefaultTransport").Call(),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildRoundTrip() []jen.Code {
	lines := []jen.Code{
		jen.Comment("RoundTrip implements the http.RoundTripper interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").PointerTo().ID("defaultRoundTripper")).ID("RoundTrip").Params(
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params(
			jen.PointerTo().Qual("net/http", "Response"),
			jen.Error(),
		).Body(
			jen.ID(constants.RequestVarName).Dot("Header").Dot("Set").Call(
				jen.ID("userAgentHeader"),
				jen.ID("userAgent"),
			),
			jen.Return().ID("t").Dot("baseTransport").Dot("RoundTrip").Call(
				jen.ID(constants.RequestVarName),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildDefaultTransport() []jen.Code {
	lines := []jen.Code{
		jen.Comment("buildDefaultTransport constructs a new http.Transport."),
		jen.Line(),
		jen.Func().ID("buildDefaultTransport").Params().Params(jen.PointerTo().Qual("net/http", "Transport")).Body(
			jen.Return().AddressOf().Qual("net/http", "Transport").Valuesln(
				jen.ID("Proxy").MapAssign().Qual("net/http", "ProxyFromEnvironment"),
				jen.ID("DialContext").MapAssign().Parens(jen.AddressOf().Qual("net", "Dialer").Valuesln(
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
		jen.Line(),
	}

	return lines
}
