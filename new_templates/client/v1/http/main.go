package client

import jen "github.com/dave/jennifer/jen"

func mainDotGo() *jen.File {
	ret := jen.NewFile("client")
	ret.Add(jen.Null())

	addImports(ret)
	ret.Add(jen.Line())

	// consts
	ret.Add(
		jen.Const().Defs(
			jen.Id("defaultTimeout").Op("=").Lit(5).Op("*").Qual("time", "Second"),
			jen.Id("clientName").Op("=").Lit("v1_client"),
		),
	)

	// vars
	ret.Add(
		jen.Var().Defs(
			jen.Id("ErrNotFound").Op("=").Qual("errors", "New").Call(jen.Lit("404: not found")),
			jen.Id("ErrUnauthorized").Op("=").Qual("errors", "New").Call(jen.Lit("401: not authorized")),
		),
		jen.Line(),
	)

	// types
	ret.Add(
		jen.Type().Id(v1).Struct(
			jen.Id("plainClient").Op("*").Qual("net/http", "Client"),
			jen.Id("authedClient").Op("*").Qual("net/http", "Client"),
			jen.Id("logger").Qual("logging", "Logger"),
			jen.Id("Debug").Id("bool"),
			jen.Id("URL").Op("*").Qual("net/url", "URL"),
			jen.Id("Scopes").Index().Id("string"),
			jen.Id("tokenSource").Qual(coreOAuth2Pkg, "TokenSource"),
		),
		jen.Line(),
	)

	// NewClient
	ret.Add(
		jen.Comment("NewClient builds a new API client for us"),
		jen.Line(),
		jen.Func().Id("NewClient").Params(
			ctxParam(),
			jen.List(
				jen.Id("clientID"),
				jen.Id("clientSecret"),
			).Id("string"),
			jen.Id("address").Op("*").Qual("net/url", "URL"),
			jen.Id("logger").Qual(loggingPkg, "Logger"),
			jen.Id("hclient").Op("*").Qual("net/http", "Client"),
			jen.Id("scopes").Index().Id("string"),
			jen.Id("debug").Id("bool"),
		).Params(jen.Op("*").Id(v1),
			jen.Id("error")).Block(
			jen.Var().Id("client").Op("=").Id("hclient"),
			jen.If(jen.Id("client").Op("==").Id("nil")).Block(
				jen.Id("client").Op("=").Op("&").Qual("net/http", "Client").Values(jen.Dict{
					jen.Id("Timeout"): jen.Id("defaultTimeout"),
				}),
			),
			jen.If(jen.Id("client").Dot("Timeout").Op("==").Lit(0)).Block(
				jen.Id("client").Dot("Timeout").Op("=").Id("defaultTimeout"),
			),
			jen.Line(),
			jen.If(jen.Id("debug")).Block(
				jen.Id("logger").Dot("SetLevel").Call(
					jen.Qual(loggingPkg, "DebugLevel"),
				),
				jen.Id("logger").Dot("Debug").Call(
					jen.Lit("log level set to debug!"),
				),
			),
			jen.Line(),
			jen.List(
				jen.Id("ac"),
				jen.Id("ts"),
			).Op(":=").Id("buildOAuthClient").Call(
				jen.Id("ctx"),
				jen.Id("address"),
				jen.Id("clientID"),
				jen.Id("clientSecret"),
				jen.Id("scopes"),
			),
			jen.Line(),
			jen.Id("c").Op(":=").Op("&").Id(v1).Values(jen.Dict{
				jen.Id("URL"):          jen.Id("address"),
				jen.Id("plainClient"):  jen.Id("client"),
				jen.Id("logger"):       jen.Qual("logger", "WithName").Call(jen.Id("clientName")),
				jen.Id("Debug"):        jen.Id("debug"),
				jen.Id("authedClient"): jen.Id("ac"),
				jen.Id("tokenSource"):  jen.Id("ts"),
			}),
			jen.Line(),
			jen.Id("logger").Dot("WithValue").Call(
				jen.Lit("url"),
				jen.Id("address").Dot("String").Call(),
			).Dot("Debug").Call(
				jen.Lit("returning client"),
			),
			jen.Return().List(jen.Id("c"),
				jen.Id("nil")),
		),
		jen.Line(),
	)

	// buildOAuthClient
	ret.Add(
		jen.Comment("buildOAuthClient does too much"),
		jen.Line(),
		jen.Func().Id("buildOAuthClient").Params(
			ctxParam(),
			jen.Id("uri").Op("*").Qual("net/url", "URL"),
			jen.List(
				jen.Id("clientID"),
				jen.Id("clientSecret"),
			).Id("string"),
			jen.Id("scopes").Index().Id("string"),
		).Params(
			jen.Op("*").Qual("net/http", "Client"),
			jen.Id("oauth2").Dot("TokenSource"),
		).Block(
			jen.Id("conf").Op(":=").Qual("golang.org/x/oauth2/clientcredentials", "Config").Values(jen.Dict{
				jen.Id("ClientID"):     jen.Id("clientID"),
				jen.Id("ClientSecret"): jen.Id("clientSecret"),
				jen.Id("Scopes"):       jen.Id("scopes"),
				jen.Id("EndpointParams"): jen.Qual("net/url", "Values").Values(jen.Dict{
					jen.Lit("client_id"):     jen.Index().Id("string").Values(jen.Id("clientID")),
					jen.Lit("client_secret"): jen.Index().Id("string").Values(jen.Id("clientSecret")),
				}),
				jen.Id("TokenURL"): jen.Id("tokenEndpoint").Call(
					jen.Id("uri"),
				).Dot("TokenURL"),
			}),
			jen.Line(),
			jen.Id("ts").Op(":=").Id("oauth2").Dot("ReuseTokenSource").Call(
				jen.Id("nil"),
				jen.Id("conf").Dot("TokenSource").Call(
					jen.Id("ctx"),
				),
			),
			jen.Id("client").Op(":=").Op("&").Qual("net/http", "Client").Values(
				jen.Id("Transport").Op(":").Op("&").Id("oauth2").Dot("Transport").Values(
					jen.Id("Base").Op(":").Op("&").Qual("go.opencensus.io/plugin/ochttp", "Transport").Values(
						jen.Id("Base").Op(":").Id("newDefaultRoundTripper").Call(),
					),
					jen.Id("Source").Op(":").Id("ts"),
				),
				jen.Id("Timeout").Op(":").Lit(5).Op("*").Qual("time", "Second"),
			),
			jen.Line(),
			jen.Return().List(jen.Id("client"),
				jen.Id("ts")),
		),
		jen.Line(),
	)

	// tokenEndpoint
	ret.Add(
		jen.Comment("tokenEndpoint provides the oauth2 Endpoint for a given host"),
		jen.Line(),
		jen.Func().Id("tokenEndpoint").Params(
			jen.Id("baseURL").Op("*").Qual("net/url", "URL"),
		).Params(
			jen.Id("oauth2").Dot("Endpoint"),
		).Block(
			jen.List(
				jen.Id("tu"),
				jen.Id("au"),
			).Op(":=").List(jen.Op("*").Id("baseURL"), jen.Op("*").Id("baseURL")),
			jen.List(
				jen.Id("tu").Dot("Path"),
				jen.Id("au").Dot("Path"),
			).Op("=").List(
				jen.Lit("oauth2/token"),
				jen.Lit("oauth2/authorize"),
			),
			jen.Line(),
			jen.Return().Id("oauth2").Dot("Endpoint").Values(jen.Dict{
				jen.Id("TokenURL"): jen.Id("tu").Dot("String").Call(),
				jen.Id("AuthURL"):  jen.Id("au").Dot("String").Call(),
			}),
		),
		jen.Line(),
	)

	// NewSimpleClient
	ret.Add(comments(
		"NewSimpleClient is a client that is capable of much less than the normal client",
		"and has noops or empty values for most of its authentication and debug parts.",
		"Its purpose at the time of this writing is merely so I can make users (which",
		"is a route that doesn't require authentication)",
	)...)
	ret.Add(
		jen.Func().Id("NewSimpleClient").Params(
			ctxParam(),
			jen.Id("address").Op("*").Qual("net/url", "URL"),
			jen.Id("debug").Id("bool"),
		).Params(
			jen.Op("*").Id(v1),
			jen.Id("error"),
		).Block(
			jen.Id("l").Op(":=").Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.Id("h").Op(":=").Op("&").Qual("net/http", "Client").Values(
				jen.Id("Timeout").Op(":").Lit(5).Op("*").Qual("time", "Second"),
			),
			jen.List(
				jen.Id("c"),
				jen.Id("err"),
			).Op(":=").Id("NewClient").Call(
				jen.Id("ctx"),
				jen.Lit(""),
				jen.Lit(""),
				jen.Id("address"),
				jen.Id("l"),
				jen.Id("h"),
				jen.Index().Id("string").Values(jen.Lit("*")),
				jen.Id("debug"),
			),
			jen.Return().List(
				jen.Id("c"),
				jen.Id("err"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("AuthenticatedClient provides the client's authenticated HTTP client"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("AuthenticatedClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Block(
			jen.Return().Id("c").Dot("authedClient"),
		),
		jen.Line(),
	)

	// c.PlainClient
	ret.Add(
		jen.Comment("PlainClient provides the client's unauthenticated HTTP client"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("PlainClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Block(
			jen.Return().Id("c").Dot("plainClient"),
		),
		jen.Line(),
	)

	// c.TokenSource
	ret.Add(
		jen.Comment("TokenSource provides the client's token source"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("TokenSource").Params().Params(jen.Id("oauth2").Dot("TokenSource")).Block(
			jen.Return().Id("c").Dot("tokenSource"),
		),
		jen.Line(),
	)

	// c.executeRawRequest
	ret.Add(comments(
		"executeRawRequest takes a given *http.Request and executes it with the provided",
		"client, alongside some debugging logging.",
	)...)
	ret.Add(
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("executeRawRequest").Params(
			ctxParam(),
			jen.Id("client").Op("*").Qual("net/http", "Client"),
			jen.Id("req").Op("*").Qual("net/http", "Request"),
		).Params(jen.Op("*").Qual("net/http", "Response"),
			jen.Id("error")).Block(
			jen.Var().Id("logger").Op("=").Id("c").Dot("logger"),
			jen.If(jen.List(
				jen.Id("command"),
				jen.Id("err"),
			).Op(":=").Qual("github.com/moul/http2curl", "GetCurlCommand").Call(
				jen.Id("req"),
			),
				jen.Id("err").Op("==").Id("nil").Op("&&").Id("c").Dot("Debug"),
			).Block(
				jen.Id("logger").Op("=").Id("c").Dot("logger").Dot("WithValue").Call(
					jen.Lit("curl"),
					jen.Id("command").Dot("String").Call(),
				),
			),
			jen.Line(),
			jen.List(
				jen.Id("res"),
				jen.Id("err"),
			).Op(":=").Id("client").Dot("Do").Call(
				jen.Id("req").Dot("WithContext").Call(
					jen.Id("ctx"),
				),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("executing request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.Id("c").Dot("Debug")).Block(
				jen.List(
					jen.Id("bdump"),
					jen.Id("err"),
				).Op(":=").Id("httputil").Dot("DumpResponse").Call(
					jen.Id("res"),
					jen.Id("true"),
				),
				jen.If(jen.Id("err").Op("==").Id("nil").Op("&&").Id("req").Dot("Method").Op("!=").Qual("net/http", "MethodGet")).Block(
					jen.Id("logger").Op("=").Id("logger").Dot("WithValue").Call(
						jen.Lit("response_body"),
						jen.Id("string").Call(
							jen.Id("bdump"),
						),
					),
				),
				jen.Id("logger").Dot("Debug").Call(
					jen.Lit("request executed"),
				),
			),
			jen.Line(),
			jen.Return().List(jen.Id("res"),
				jen.Id("nil")),
		),
		jen.Line(),
	)

	// c.BuildURL
	ret.Add(
		jen.Comment("BuildURL builds standard service URLs"),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildURL").Params(
			jen.Id("qp").Qual("net/url", "Values"),
			jen.Id("parts").Op("...").Id("string"),
		).Params(jen.Id("string")).Block(
			jen.If(jen.Id("qp").Op("!=").Id("nil")).Block(
				jen.Return().Id("c").Dot("buildURL").Call(
					jen.Id("qp"),
					jen.Id("parts").Op("..."),
				).Dot("String").Call(),
			),
			jen.Return().Id("c").Dot("buildURL").Call(
				jen.Id("nil"),
				jen.Id("parts").Op("..."),
			).Dot("String").Call(),
		),
		jen.Line(),
	)

	// c.buildURL
	ret.Add(comments(
		"buildURL takes a given set of query parameters and URL parts, and returns",
		"a parsed URL object from them.",
	)...)
	ret.Add(
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("buildURL").Params(
			jen.Id("queryParams").Qual("net/url", "Values"),
			jen.Id("parts").Op("...").Id("string"),
		).Params(jen.Op("*").Qual("net/url", "URL")).Block(
			jen.Id("tu").Op(":=").Op("*").Id("c").Dot("URL"),
			jen.Line(),
			jen.Id("parts").Op("=").Id("append").Call(
				jen.Index().Id("string").Values(jen.Lit("api"), jen.Lit("v1")),
				jen.Id("parts").Op("..."),
			),
			jen.List(
				jen.Id("u"),
				jen.Id("err"),
			).Op(":=").Qual("net/url", "Parse").Call(
				jen.Qual("strings", "Join").Call(
					jen.Id("parts"),
					jen.Lit("/"),
				),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("panic").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("was asked to build an invalid URL: %v"),
						jen.Id("err"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.Id("queryParams").Op("!=").Id("nil")).Block(
				jen.Id("u").Dot("RawQuery").Op("=").Id("queryParams").Dot("Encode").Call(),
			),
			jen.Line(),
			jen.Return().Id("tu").Dot("ResolveReference").Call(
				jen.Id("u"),
			),
		),
		jen.Line(),
	)

	// c.buildVersionlessURL
	ret.Add(comments(
		"buildVersionlessURL builds a URL without the `/api/v1/` prefix. It should",
		"otherwise be identical to buildURL",
	)...)
	ret.Add(
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("buildVersionlessURL").Params(
			jen.Id("qp").Qual("net/url", "Values"),
			jen.Id("parts").Op("...").Id("string"),
		).Params(jen.Id("string")).Block(
			jen.Id("tu").Op(":=").Op("*").Id("c").Dot("URL"),
			jen.Line(),
			jen.List(
				jen.Id("u"),
				jen.Id("err"),
			).Op(":=").Qual("net/url", "Parse").Call(
				jen.Qual("path", "Join").Call(
					jen.Id("parts").Op("..."),
				),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("panic").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("user tried to build an invalid URL: %v"),
						jen.Id("err"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.Id("qp").Op("!=").Id("nil")).Block(
				jen.Id("u").Dot("RawQuery").Op("=").Id("qp").Dot("Encode").Call(),
			),
			jen.Line(),
			jen.Return().Id("tu").Dot("ResolveReference").Call(
				jen.Id("u"),
			).Dot("String").Call(),
		),
		jen.Line(),
	)

	// c.BuildWebsocketURL
	ret.Add(comments("BuildWebsocketURL builds a standard URL and then converts its scheme to the websocket protocol")...)
	ret.Add(
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("BuildWebsocketURL").Params(
			jen.Id("parts").Op("...").Id("string"),
		).Params(jen.Id("string")).Block(
			jen.Id("u").Op(":=").Id("c").Dot("buildURL").Call(
				jen.Id("nil"),
				jen.Id("parts").Op("..."),
			),
			jen.Id("u").Dot("Scheme").Op("=").Lit("ws"),
			jen.Line(),
			jen.Return().Id("u").Dot("String").Call(),
		),
		jen.Line(),
	)

	// c.BuildHealthCheckRequest
	ret.Add(comments("BuildHealthCheckRequest builds a health check HTTP Request")...)
	ret.Add(
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("BuildHealthCheckRequest").Params().Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("u").Op(":=").Op("*").Id("c").Dot("URL"),
			jen.Id("uri").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s://%s/_meta_/ready"),
				jen.Id("u").Dot("Scheme"),
				jen.Id("u").Dot("Host"),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	// c.IsUp
	ret.Add(comments("IsUp returns whether or not the service's health endpoint is returning 200s")...)
	ret.Add(
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("IsUp").Params().Params(jen.Id("bool")).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildHealthCheckRequest").Call(),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("c").Dot("logger").Dot("Error").Call(
					jen.Id("err"),
					jen.Lit("building request"),
				),
				jen.Return().Id("false"),
			),
			jen.Line(),
			jen.List(
				jen.Id("res"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("plainClient").Dot("Do").Call(
				jen.Id("req"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("c").Dot("logger").Dot("Error").Call(
					jen.Id("err"),
					jen.Lit("health check"),
				),
				jen.Return().Id("false"),
			),
			jen.Line(),
			jen.Defer().Func().Params().Block(
				jen.If(jen.Id("err").Op(":=").Id("res").Dot("Body").Dot("Close").Call(),
					jen.Id("err").Op("!=").Id("nil")).Block(
					jen.Id("c").Dot("logger").Dot("Error").Call(
						jen.Id("err"),
						jen.Lit("closing response body"),
					),
				),
			).Call(),
			jen.Line(),
			jen.Return().Id("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusOK"),
		),
		jen.Line(),
	)

	// c.buildDataRequest
	ret.Add(comments("buildDataRequest builds an HTTP request for a given method, URL, and body data.")...)
	ret.Add(
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("buildDataRequest").Params(
			jen.List(
				jen.Id("method"),
				jen.Id("uri"),
			).Id("string"),
			jen.Id("in").Interface(),
		).Params(jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error")).Block(
			jen.List(
				jen.Id("body"),
				jen.Id("err"),
			).Op(":=").Id("createBodyFromStruct").Call(
				jen.Id("in"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"),
					jen.Id("err")),
			),
			jen.Line(),
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Qual("net/http", "NewRequest").Call(
				jen.Id("method"),
				jen.Id("uri"),
				jen.Id("body"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Id("err"),
				),
			),
			jen.Line(),
			jen.Id("req").Dot("Header").Dot("Set").Call(
				jen.Lit("Content-type"),
				jen.Lit("application/json"),
			),
			jen.Return().List(jen.Id("req"),
				jen.Id("nil")),
		),
		jen.Line(),
	)

	// c.retrieve
	ret.Add(comments("retrieve executes an HTTP request and loads the response content into a struct")...)
	ret.Add(
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("retrieve").Params(
			ctxParam(),
			jen.Id("req").Op("*").Qual("net/http", "Request"),
			jen.Id("obj").Interface(),
		).Params(jen.Id("error")).Block(
			jen.If(jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(
				jen.Id("obj"),
			),
				jen.Id("err").Op("!=").Id("nil"),
			).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("struct to load must be a pointer: %w"),
					jen.Id("err"),
				),
			),
			jen.Line(),
			jen.List(
				jen.Id("res"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("executeRawRequest").Call(
				jen.Id("ctx"),
				jen.Id("c").Dot("authedClient"),
				jen.Id("req"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("executing request: %w"),
					jen.Id("err"),
				),
			),
			jen.Line(),
			jen.If(jen.Id("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().Id("ErrNotFound"),
			),
			jen.Line(),
			jen.Return().Id("unmarshalBody").Call(
				jen.Id("res"), jen.Op("&").Id("obj"),
			),
		),
		jen.Line(),
	)

	// c.executeRequest
	ret.Add(comments(
		"executeRequest takes a given request and executes it with the auth client. It returns some errors",
		"upon receiving certain status codes, but otherwise will return nil upon success.",
	)...)
	ret.Add(
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("executeRequest").Params(
			ctxParam(),
			jen.Id("req").Op("*").Qual("net/http", "Request"),
			jen.Id("out").Interface(),
		).Params(jen.Id("error")).Block(
			jen.List(
				jen.Id("res"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("executeRawRequest").Call(
				jen.Id("ctx"),
				jen.Id("c").Dot("authedClient"),
				jen.Id("req"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("executing request: %w"),
					jen.Id("err"),
				),
			),
			jen.Switch(jen.Id("res").Dot("StatusCode")).Block(
				jen.Case(jen.Qual("net/http", "StatusNotFound")).Block(
					jen.Return().Id("ErrNotFound"),
				),
				jen.Case(jen.Qual("net/http", "StatusUnauthorized")).Block(
					jen.Return().Id("ErrUnauthorized"),
				),
			),
			jen.Line(),
			jen.If(jen.Id("out").Op("!=").Id("nil")).Block(
				jen.Id("resErr").Op(":=").Id("unmarshalBody").Call(
					jen.Id("res"), jen.Op("&").Id("out"),
				),
				jen.If(jen.Id("resErr").Op("!=").Id("nil")).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("loading response from server: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Line(),
			jen.Return().Id("nil"),
		),
		jen.Line(),
	)

	// c.executeUnathenticatedDataRequest
	ret.Add(comments("c.executeUnathenticatedDataRequest")...)
	ret.Add(
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("executeUnathenticatedDataRequest").Params(
			ctxParam(),
			jen.Id("req").Op("*").Qual("net/http", "Request"),
			jen.Id("out").Interface(),
		).Params(jen.Id("error")).Block(
			jen.If(jen.Id("out").Op("!=").Id("nil")).Block(
				jen.If(jen.List(
					jen.Id("np"),
					jen.Id("err"),
				).Op(":=").Id("argIsNotPointer").Call(
					jen.Id("out"),
				),
					jen.Id("np").Op("||").Id("err").Op("!=").Id("nil"),
				).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("struct to load must be a pointer: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.List(
				jen.Id("res"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("executeRawRequest").Call(
				jen.Id("ctx"),
				jen.Id("c").Dot("plainClient"),
				jen.Id("req"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("executing request: %w"),
					jen.Id("err"),
				),
			),
			jen.If(jen.Id("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().Id("ErrNotFound"),
			),
			jen.If(jen.Id("out").Op("!=").Id("nil")).Block(
				jen.Id("resErr").Op(":=").Id("unmarshalBody").Call(
					jen.Id("res"), jen.Op("&").Id("out"),
				),
				jen.If(jen.Id("resErr").Op("!=").Id("nil")).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("loading response from server: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Return().Id("nil"),
		),
		jen.Line(),
	)

	return ret
}
