package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperBaseDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("userAgentHeader").Op("=").Lit("User-Agent"),
			jen.ID("userAgent").Op("=").Lit("TODO Service Client"),
			jen.ID("maxRetryCount").Op("=").Lit(5),
			jen.ID("minRetryWait").Op("=").Lit(100).Op("*").Qual("time", "Millisecond"),
			jen.ID("maxRetryWait").Op("=").Qual("time", "Second"),
			jen.ID("keepAlive").Op("=").Lit(30).Op("*").Qual("time", "Second"),
			jen.ID("tlsHandshakeTimeout").Op("=").Lit(10).Op("*").Qual("time", "Second"),
			jen.ID("expectContinueTimeout").Op("=").Lit(2).Op("*").ID("defaultTimeout"),
			jen.ID("idleConnTimeout").Op("=").Lit(3).Op("*").ID("defaultTimeout"),
			jen.ID("maxIdleConns").Op("=").Lit(100),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("defaultRoundTripper").Struct(jen.ID("baseRoundTripper").Qual("net/http", "RoundTripper")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("newDefaultRoundTripper constructs a new http.RoundTripper."),
		jen.Line(),
		jen.Func().ID("newDefaultRoundTripper").Params(jen.ID("timeout").Qual("time", "Duration")).Params(jen.Qual("net/http", "RoundTripper")).Body(
			jen.Return().Op("&").ID("defaultRoundTripper").Valuesln(jen.ID("baseRoundTripper").Op(":").ID("buildWrappedTransport").Call(jen.ID("timeout")))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RoundTrip implements the http.RoundTripper interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("defaultRoundTripper")).ID("RoundTrip").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual("net/http", "Response"), jen.ID("error")).Body(
			jen.ID("req").Dot("Header").Dot("Set").Call(
				jen.ID("userAgentHeader"),
				jen.ID("userAgent"),
			),
			jen.Return().ID("t").Dot("baseRoundTripper").Dot("RoundTrip").Call(jen.ID("req")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildWrappedTransport constructs a new http.Transport."),
		jen.Line(),
		jen.Func().ID("buildWrappedTransport").Params(jen.ID("timeout").Qual("time", "Duration")).Params(jen.Qual("net/http", "RoundTripper")).Body(
			jen.If(jen.ID("timeout").Op("==").Lit(0)).Body(
				jen.ID("timeout").Op("=").ID("defaultTimeout")),
			jen.ID("t").Op(":=").Op("&").Qual("net/http", "Transport").Valuesln(jen.ID("Proxy").Op(":").Qual("net/http", "ProxyFromEnvironment"), jen.ID("DialContext").Op(":").Parens(jen.Op("&").Qual("net", "Dialer").Valuesln(jen.ID("Timeout").Op(":").ID("timeout"), jen.ID("KeepAlive").Op(":").ID("keepAlive"))).Dot("DialContext"), jen.ID("MaxIdleConns").Op(":").ID("maxIdleConns"), jen.ID("MaxIdleConnsPerHost").Op(":").ID("maxIdleConns"), jen.ID("TLSHandshakeTimeout").Op(":").ID("tlsHandshakeTimeout"), jen.ID("ExpectContinueTimeout").Op(":").ID("expectContinueTimeout"), jen.ID("IdleConnTimeout").Op(":").ID("idleConnTimeout")),
			jen.Return().ID("otelhttp").Dot("NewTransport").Call(
				jen.ID("t"),
				jen.ID("otelhttp").Dot("WithSpanNameFormatter").Call(jen.ID("tracing").Dot("FormatSpan")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildRequestLogHook").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.Func().Params(jen.ID("retryablehttp").Dot("Logger"), jen.Op("*").Qual("net/http", "Request"), jen.ID("int"))).Body(
			jen.ID("logger").Op("=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")),
			jen.Return().Func().Params(jen.ID("_").ID("retryablehttp").Dot("Logger"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("numRetries").ID("int")).Body(
				jen.If(jen.ID("req").Op("!=").ID("nil")).Body(
					jen.ID("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("WithValue").Call(
						jen.Lit("retry_count"),
						jen.ID("numRetries"),
					).Dot("Debug").Call(jen.Lit("making request")))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildResponseLogHook").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.Func().Params(jen.ID("retryablehttp").Dot("Logger"), jen.Op("*").Qual("net/http", "Response"))).Body(
			jen.ID("logger").Op("=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")),
			jen.Return().Func().Params(jen.ID("_").ID("retryablehttp").Dot("Logger"), jen.ID("res").Op("*").Qual("net/http", "Response")).Body(
				jen.If(jen.ID("res").Op("!=").ID("nil")).Body(
					jen.ID("logger").Dot("WithResponse").Call(jen.ID("res")).Dot("Debug").Call(jen.Lit("received response")))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildCheckRetryFunc").Params(jen.ID("tracer").ID("tracing").Dot("Tracer")).Params(jen.Func().Params(jen.Qual("context", "Context"), jen.Op("*").Qual("net/http", "Response"), jen.ID("error")).Params(jen.ID("bool"), jen.ID("error"))).Body(
			jen.Return().Func().Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Op("*").Qual("net/http", "Response"), jen.ID("err").ID("error")).Params(jen.ID("bool"), jen.ID("error")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracer").Dot("StartCustomSpan").Call(
					jen.ID("ctx"),
					jen.Lit("CheckRetry"),
				),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.If(jen.ID("res").Op("!=").ID("nil")).Body(
					jen.ID("tracing").Dot("AttachResponseToSpan").Call(
						jen.ID("span"),
						jen.ID("res"),
					)),
				jen.Return().ID("retryablehttp").Dot("DefaultRetryPolicy").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildErrorHandler").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.Func().Params(jen.ID("res").Op("*").Qual("net/http", "Response"), jen.ID("err").ID("error"), jen.ID("numTries").ID("int")).Params(jen.Op("*").Qual("net/http", "Response"), jen.ID("error"))).Body(
			jen.ID("logger").Op("=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")),
			jen.Return().Func().Params(jen.ID("res").Op("*").Qual("net/http", "Response"), jen.ID("err").ID("error"), jen.ID("numTries").ID("int")).Params(jen.Op("*").Qual("net/http", "Response"), jen.ID("error")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("try_number"),
					jen.ID("numTries"),
				).Dot("WithResponse").Call(jen.ID("res")).Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("executing request"),
				),
				jen.Return().List(jen.ID("res"), jen.ID("err")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildRetryingClient").Params(jen.ID("client").Op("*").Qual("net/http", "Client"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("tracer").ID("tracing").Dot("Tracer")).Params(jen.Op("*").Qual("net/http", "Client")).Body(
			jen.ID("rc").Op(":=").Op("&").ID("retryablehttp").Dot("Client").Valuesln(jen.ID("HTTPClient").Op(":").ID("client"), jen.ID("Logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")), jen.ID("RetryWaitMin").Op(":").ID("minRetryWait"), jen.ID("RetryWaitMax").Op(":").ID("maxRetryWait"), jen.ID("RetryMax").Op(":").ID("maxRetryCount"), jen.ID("RequestLogHook").Op(":").ID("buildRequestLogHook").Call(jen.ID("logger")), jen.ID("ResponseLogHook").Op(":").ID("buildResponseLogHook").Call(jen.ID("logger")), jen.ID("CheckRetry").Op(":").ID("buildCheckRetryFunc").Call(jen.ID("tracer")), jen.ID("Backoff").Op(":").ID("retryablehttp").Dot("DefaultBackoff"), jen.ID("ErrorHandler").Op(":").ID("buildErrorHandler").Call(jen.ID("logger"))),
			jen.ID("c").Op(":=").ID("rc").Dot("StandardClient").Call(),
			jen.ID("c").Dot("Timeout").Op("=").ID("defaultTimeout"),
			jen.Return().ID("c"),
		),
		jen.Line(),
	)

	return code
}
