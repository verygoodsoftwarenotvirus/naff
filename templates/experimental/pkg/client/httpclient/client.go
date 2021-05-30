package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("defaultTimeout").Op("=").Lit(30).Op("*").Qual("time", "Second"),
			jen.ID("clientName").Op("=").Lit("todo_client_v1"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("authMethod").Struct(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("cookieAuthMethod").Op("=").ID("new").Call(jen.ID("authMethod")),
			jen.ID("pasetoAuthMethod").Op("=").ID("new").Call(jen.ID("authMethod")),
			jen.ID("defaultContentType").Op("=").ID("encoding").Dot("ContentTypeJSON"),
			jen.ID("errInvalidResponseCode").Op("=").Qual("errors", "New").Call(jen.Lit("invalid response code")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Client").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("panicker").ID("panicking").Dot("Panicker"),
				jen.ID("url").Op("*").Qual("net/url", "URL"),
				jen.ID("requestBuilder").Op("*").ID("requests").Dot("Builder"),
				jen.ID("encoder").ID("encoding").Dot("ClientEncoder"),
				jen.ID("unauthenticatedClient").Op("*").Qual("net/http", "Client"),
				jen.ID("authedClient").Op("*").Qual("net/http", "Client"),
				jen.ID("authMethod").Op("*").ID("authMethod"),
				jen.ID("accountID").ID("uint64"),
				jen.ID("debug").ID("bool"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AuthenticatedClient returns the authenticated *httpclient.Client that we use to make most requests."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("AuthenticatedClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Body(
			jen.Return().ID("c").Dot("authedClient")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PlainClient returns the unauthenticated *httpclient.Client that we use to make certain requests."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("PlainClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Body(
			jen.Return().ID("c").Dot("unauthenticatedClient")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("URL provides the client's URL."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("URL").Params().Params(jen.Op("*").Qual("net/url", "URL")).Body(
			jen.Return().ID("c").Dot("url")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RequestBuilder provides the client's *requests.Builder."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("RequestBuilder").Params().Params(jen.Op("*").ID("requests").Dot("Builder")).Body(
			jen.Return().ID("c").Dot("requestBuilder")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewClient builds a new API client for us."),
		jen.Line(),
		jen.Func().ID("NewClient").Params(jen.ID("u").Op("*").Qual("net/url", "URL"), jen.ID("options").Op("...").ID("option")).Params(jen.Op("*").ID("Client"), jen.ID("error")).Body(
			jen.ID("l").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.If(jen.ID("u").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNoURLProvided"))),
			jen.ID("c").Op(":=").Op("&").ID("Client").Valuesln(jen.ID("url").Op(":").ID("u"), jen.ID("logger").Op(":").ID("l"), jen.ID("debug").Op(":").ID("false"), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("clientName")), jen.ID("panicker").Op(":").ID("panicking").Dot("NewProductionPanicker").Call(), jen.ID("encoder").Op(":").ID("encoding").Dot("ProvideClientEncoder").Call(
				jen.ID("l"),
				jen.ID("encoding").Dot("ContentTypeJSON"),
			), jen.ID("authedClient").Op(":").Op("&").Qual("net/http", "Client").Valuesln(jen.ID("Transport").Op(":").ID("buildWrappedTransport").Call(jen.ID("defaultTimeout")), jen.ID("Timeout").Op(":").ID("defaultTimeout")), jen.ID("unauthenticatedClient").Op(":").Op("&").Qual("net/http", "Client").Valuesln(jen.ID("Transport").Op(":").ID("buildWrappedTransport").Call(jen.ID("defaultTimeout")), jen.ID("Timeout").Op(":").ID("defaultTimeout"))),
			jen.List(jen.ID("requestBuilder"), jen.ID("err")).Op(":=").ID("requests").Dot("NewBuilder").Call(
				jen.ID("c").Dot("url"),
				jen.ID("c").Dot("logger"),
				jen.ID("encoding").Dot("ProvideClientEncoder").Call(
					jen.ID("l"),
					jen.ID("defaultContentType"),
				),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("c").Dot("requestBuilder").Op("=").ID("requestBuilder"),
			jen.For(jen.List(jen.ID("_"), jen.ID("opt")).Op(":=").Range().ID("options")).Body(
				jen.If(jen.ID("optionSetErr").Op(":=").ID("opt").Call(jen.ID("c")), jen.ID("optionSetErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("optionSetErr")))),
			jen.Return().List(jen.ID("c"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("closeResponseBody takes a given HTTP response and closes its body, logging if an error occurs."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("closeResponseBody").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Op("*").Qual("net/http", "Response")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithResponse").Call(jen.ID("res")),
			jen.If(jen.ID("res").Op("!=").ID("nil")).Body(
				jen.If(jen.ID("err").Op(":=").ID("res").Dot("Body").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("closing response body"),
					))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("loggerWithFilter prepares a logger from the Client logger that has relevant filter information."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("loggerWithFilter").Params(jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("logging").Dot("Logger")).Body(
			jen.If(jen.ID("filter").Op("==").ID("nil")).Body(
				jen.Return().ID("c").Dot("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("FilterIsNilKey"),
					jen.ID("true"),
				)),
			jen.Return().ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("FilterLimitKey"),
				jen.ID("filter").Dot("Limit"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("FilterPageKey"),
				jen.ID("filter").Dot("Page"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildURL builds standard service URLs."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("BuildURL").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("qp").Qual("net/url", "Values"), jen.ID("parts").Op("...").ID("string")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("u").Op(":=").ID("c").Dot("buildRawURL").Call(
				jen.ID("ctx"),
				jen.ID("qp"),
				jen.ID("parts").Op("..."),
			), jen.ID("u").Op("!=").ID("nil")).Body(
				jen.Return().ID("u").Dot("String").Call()),
			jen.Return().Lit(""),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildRawURL takes a given set of query parameters and url parts, and returns a parsed url object from them."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("buildRawURL").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("queryParams").Qual("net/url", "Values"), jen.ID("parts").Op("...").ID("string")).Params(jen.Op("*").Qual("net/url", "URL")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tu").Op(":=").Op("*").ID("c").Dot("url"),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("URLQueryKey"),
				jen.ID("queryParams").Dot("Encode").Call(),
			),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.Qual("path", "Join").Call(jen.ID("append").Call(
				jen.Index().ID("string").Valuesln(jen.Lit("api"), jen.Lit("v1")),
				jen.ID("parts").Op("..."),
			).Op("..."))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building URL"),
				),
				jen.Return().ID("nil"),
			),
			jen.If(jen.ID("queryParams").Op("!=").ID("nil")).Body(
				jen.ID("u").Dot("RawQuery").Op("=").ID("queryParams").Dot("Encode").Call()),
			jen.Return().ID("tu").Dot("ResolveReference").Call(jen.ID("u")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildVersionlessURL builds a url without the `/api/v1/` prefix. It should otherwise be identical to buildRawURL."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("buildVersionlessURL").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("qp").Qual("net/url", "Values"), jen.ID("parts").Op("...").ID("string")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tu").Op(":=").Op("*").ID("c").Dot("url"),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.Qual("path", "Join").Call(jen.ID("parts").Op("..."))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachErrorToSpan").Call(
					jen.ID("span"),
					jen.Lit("building url"),
					jen.ID("err"),
				),
				jen.Return().Lit(""),
			),
			jen.If(jen.ID("qp").Op("!=").ID("nil")).Body(
				jen.ID("u").Dot("RawQuery").Op("=").ID("qp").Dot("Encode").Call()),
			jen.Return().ID("tu").Dot("ResolveReference").Call(jen.ID("u")).Dot("String").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildWebsocketURL builds a standard url and then converts its scheme to the websocket protocol."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("BuildWebsocketURL").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("qp").Qual("net/url", "Values"), jen.ID("parts").Op("...").ID("string")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("u").Op(":=").ID("c").Dot("buildRawURL").Call(
				jen.ID("ctx"),
				jen.ID("qp"),
				jen.ID("parts").Op("..."),
			),
			jen.ID("u").Dot("Scheme").Op("=").Lit("ws"),
			jen.Return().ID("u").Dot("String").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("IsUp returns whether the service's health endpoint is returning 200s."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("IsUp").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("bool")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildHealthCheckRequest").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building health check request"),
				),
				jen.Return().ID("false"),
			),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("unauthenticatedClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("performing health check"),
				),
				jen.Return().ID("false"),
			),
			jen.ID("c").Dot("closeResponseBody").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return().ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusOK"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("fetchResponseToRequest takes a given *http.Request and executes it with the provided."),
		jen.Line(),
		jen.Func().Comment("client, alongside some debugging logging.").Params(jen.ID("c").Op("*").ID("Client")).ID("fetchResponseToRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("client").Op("*").Qual("net/http", "Client"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual("net/http", "Response"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.If(jen.List(jen.ID("command"), jen.ID("err")).Op(":=").ID("http2curl").Dot("GetCurlCommand").Call(jen.ID("req")), jen.ID("err").Op("==").ID("nil").Op("&&").ID("c").Dot("debug")).Body(
				jen.ID("logger").Op("=").ID("c").Dot("logger").Dot("WithValue").Call(
					jen.Lit("curl"),
					jen.ID("command").Dot("String").Call(),
				)),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("client").Dot("Do").Call(jen.ID("req").Dot("WithContext").Call(jen.ID("ctx"))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing request"),
				))),
			jen.Var().Defs(
				jen.ID("bdump").Index().ID("byte"),
			),
			jen.If(jen.List(jen.ID("bdump"), jen.ID("err")).Op("=").ID("httputil").Dot("DumpResponse").Call(
				jen.ID("res"),
				jen.ID("true"),
			), jen.ID("err").Op("==").ID("nil")).Body(
				jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
					jen.Lit("response_body"),
					jen.ID("string").Call(jen.ID("bdump")),
				)),
			jen.ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ResponseStatusKey"),
				jen.ID("res").Dot("StatusCode"),
			).Dot("Debug").Call(jen.Lit("request executed")),
			jen.Return().List(jen.ID("res"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("executeAndUnmarshal executes a request and unmarshalls it to the provided interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("executeAndUnmarshal").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("httpClient").Op("*").Qual("net/http", "Client"), jen.ID("out").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("executing request")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
				jen.ID("ctx"),
				jen.ID("httpClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing request"),
				)),
			jen.ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ResponseStatusKey"),
				jen.ID("res").Dot("StatusCode"),
			).Dot("Debug").Call(jen.Lit("request executed")),
			jen.If(jen.ID("err").Op("=").ID("errorFromResponse").Call(jen.ID("res")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing request"),
				)),
			jen.If(jen.ID("out").Op("!=").ID("nil")).Body(
				jen.If(jen.ID("err").Op("=").ID("c").Dot("unmarshalBody").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("out"),
				), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("loading %s %d response from server"),
						jen.ID("res").Dot("Request").Dot("Method"),
						jen.ID("res").Dot("StatusCode"),
					))),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("fetchAndUnmarshal takes a given request and executes it with the auth client."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("fetchAndUnmarshal").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("out").Interface()).Params(jen.ID("error")).Body(
			jen.Return().ID("c").Dot("executeAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("c").Dot("authedClient"),
				jen.ID("out"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("fetchAndUnmarshalWithoutAuthentication takes a given request and executes it with the plain client."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("fetchAndUnmarshalWithoutAuthentication").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("out").Interface()).Params(jen.ID("error")).Body(
			jen.Return().ID("c").Dot("executeAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("c").Dot("unauthenticatedClient"),
				jen.ID("out"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("responseIsOK executes an HTTP request and loads the response content into a bool."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("responseIsOK").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("authedClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("false"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing existence request"),
				))),
			jen.ID("c").Dot("closeResponseBody").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return().List(jen.ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusOK"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
