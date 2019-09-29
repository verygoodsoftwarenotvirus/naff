package client

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

const (
	apiV1 = "api/v1"
)

func newClientMethod(name string) *jen.Statement {
	return jen.Func().Params(jen.ID("c").Op("*").ID(v1)).ID(name)
}

func mainDotGo() *jen.File {
	ret := jen.NewFile("client")
	ret.Add(jen.Null())

	addImports(ret)
	ret.Add(jen.Line())

	// consts
	ret.Add(
		jen.Const().Defs(
			jen.ID("defaultTimeout").Op("=").Lit(5).Op("*").Qual("time", "Second"),
			jen.ID("clientName").Op("=").Lit("v1_client"),
		),
	)

	// vars
	ret.Add(
		jen.Var().Defs(
			jen.Comment("ErrNotFound is a handy error to return when we receive a 404 response"),
			jen.ID("ErrNotFound").Op("=").Qual("errors", "New").Call(jen.Lit("404: not found")),
			jen.Line(),
			jen.Comment("ErrUnauthorized is a handy error to return when we receive a 404 response"),
			jen.ID("ErrUnauthorized").Op("=").Qual("errors", "New").Call(jen.Lit("401: not authorized")),
		),
		jen.Line(),
	)

	// types
	ret.Add(comments(fmt.Sprintf("%s is a client for interacting with v1 of our REST API", v1))...)
	ret.Add(
		jen.Type().ID(v1).Struct(
			jen.ID("plainClient").Op("*").Qual("net/http", "Client"),
			jen.ID("authedClient").Op("*").Qual("net/http", "Client"),
			jen.ID("logger").Qual(loggingPkg, "Logger"),
			jen.ID("Debug").ID("bool"),
			jen.ID("URL").Op("*").Qual("net/url", "URL"),
			jen.ID("Scopes").Index().ID("string"),
			jen.ID("tokenSource").Qual(coreOAuth2Pkg, "TokenSource"),
		),
		jen.Line(),
	)

	// NewClient
	ret.Add(
		jen.Comment("NewClient builds a new API client for us"),
		jen.Line(),
		jen.Func().ID("NewClient").Params(
			ctxParam(),
			jen.List(
				jen.ID("clientID"),
				jen.ID("clientSecret"),
			).ID("string"),
			jen.ID("address").Op("*").Qual("net/url", "URL"),
			jen.ID("logger").Qual(loggingPkg, "Logger"),
			jen.ID("hclient").Op("*").Qual("net/http", "Client"),
			jen.ID("scopes").Index().ID("string"),
			jen.ID("debug").ID("bool"),
		).Params(
			jen.Op("*").ID(v1),
			jen.ID("error"),
		).Block(
			jen.Var().ID("client").Op("=").ID("hclient"),
			jen.If(jen.ID("client").Op("==").ID("nil")).Block(
				jen.ID("client").Op("=").Op("&").Qual("net/http", "Client").Values(jen.Dict{
					jen.ID("Timeout"): jen.ID("defaultTimeout"),
				}),
			),
			jen.If(jen.ID("client").Dot("Timeout").Op("==").Lit(0)).Block(
				jen.ID("client").Dot("Timeout").Op("=").ID("defaultTimeout"),
			),
			jen.Line(),
			jen.If(jen.ID("debug")).Block(
				jen.ID("logger").Dot("SetLevel").Call(
					jen.Qual(loggingPkg, "DebugLevel"),
				),
				jen.ID("logger").Dot("Debug").Call(
					jen.Lit("log level set to debug!"),
				),
			),
			jen.Line(),
			jen.List(
				jen.ID("ac"),
				jen.ID("ts"),
			).Op(":=").ID("buildOAuthClient").Call(
				jen.ID("ctx"),
				jen.ID("address"),
				jen.ID("clientID"),
				jen.ID("clientSecret"),
				jen.ID("scopes"),
			),
			jen.Line(),
			jen.ID("c").Op(":=").Op("&").ID(v1).Values(jen.Dict{
				jen.ID("URL"):          jen.ID("address"),
				jen.ID("plainClient"):  jen.ID("client"),
				jen.ID("logger"):       jen.Qual("logger", "WithName").Call(jen.ID("clientName")),
				jen.ID("Debug"):        jen.ID("debug"),
				jen.ID("authedClient"): jen.ID("ac"),
				jen.ID("tokenSource"):  jen.ID("ts"),
			}),
			jen.Line(),
			jen.ID("logger").Dot("WithValue").Call(
				jen.Lit("url"),
				jen.ID("address").Dot("String").Call(),
			).Dot("Debug").Call(
				jen.Lit("returning client"),
			),
			jen.Return().List(jen.ID("c"),
				jen.ID("nil")),
		),
		jen.Line(),
	)

	// buildOAuthClient
	ret.Add(
		jen.Comment("buildOAuthClient does too much"),
		jen.Line(),
		jen.Func().ID("buildOAuthClient").Params(
			ctxParam(),
			jen.ID("uri").Op("*").Qual("net/url", "URL"),
			jen.List(
				jen.ID("clientID"),
				jen.ID("clientSecret"),
			).ID("string"),
			jen.ID("scopes").Index().ID("string"),
		).Params(
			jen.Op("*").Qual("net/http", "Client"),
			jen.ID("oauth2").Dot("TokenSource"),
		).Block(
			jen.ID("conf").Op(":=").Qual("golang.org/x/oauth2/clientcredentials", "Config").Values(jen.Dict{
				jen.ID("ClientID"):     jen.ID("clientID"),
				jen.ID("ClientSecret"): jen.ID("clientSecret"),
				jen.ID("Scopes"):       jen.ID("scopes"),
				jen.ID("EndpointParams"): jen.Qual("net/url", "Values").Values(jen.Dict{
					jen.Lit("client_id"):     jen.Index().ID("string").Values(jen.ID("clientID")),
					jen.Lit("client_secret"): jen.Index().ID("string").Values(jen.ID("clientSecret")),
				}),
				jen.ID("TokenURL"): jen.ID("tokenEndpoint").Call(
					jen.ID("uri"),
				).Dot("TokenURL"),
			}),
			jen.Line(),
			jen.ID("ts").Op(":=").ID("oauth2").Dot("ReuseTokenSource").Call(
				jen.ID("nil"),
				jen.ID("conf").Dot("TokenSource").Call(
					jen.ID("ctx"),
				),
			),
			jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(jen.Dict{
				jen.ID("Transport"): jen.Op("&").ID("oauth2").Dot("Transport").Values(jen.Dict{
					jen.ID("Base"): jen.Op("&").Qual("go.opencensus.io/plugin/ochttp", "Transport").Values(jen.Dict{
						jen.ID("Base"): jen.ID("newDefaultRoundTripper").Call(),
					}),
					jen.ID("Source"): jen.ID("ts"),
				}),
				jen.ID("Timeout"): jen.Lit(5).Op("*").Qual("time", "Second"),
			}),
			jen.Line(),
			jen.Return().List(jen.ID("client"),
				jen.ID("ts")),
		),
		jen.Line(),
	)

	// tokenEndpoint
	ret.Add(
		jen.Comment("tokenEndpoint provides the oauth2 Endpoint for a given host"),
		jen.Line(),
		jen.Func().ID("tokenEndpoint").Params(
			jen.ID("baseURL").Op("*").Qual("net/url", "URL"),
		).Params(
			jen.ID("oauth2").Dot("Endpoint"),
		).Block(
			jen.List(
				jen.ID("tu"),
				jen.ID("au"),
			).Op(":=").List(jen.Op("*").ID("baseURL"), jen.Op("*").ID("baseURL")),
			jen.List(
				jen.ID("tu").Dot("Path"),
				jen.ID("au").Dot("Path"),
			).Op("=").List(
				jen.Lit("oauth2/token"),
				jen.Lit("oauth2/authorize"),
			),
			jen.Line(),
			jen.Return().ID("oauth2").Dot("Endpoint").Values(jen.Dict{
				jen.ID("TokenURL"): jen.ID("tu").Dot("String").Call(),
				jen.ID("AuthURL"):  jen.ID("au").Dot("String").Call(),
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
		jen.Func().ID("NewSimpleClient").Params(
			ctxParam(),
			jen.ID("address").Op("*").Qual("net/url", "URL"),
			jen.ID("debug").ID("bool"),
		).Params(
			jen.Op("*").ID(v1),
			jen.ID("error"),
		).Block(
			jen.ID("l").Op(":=").Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("h").Op(":=").Op("&").Qual("net/http", "Client").Values(jen.Dict{
				jen.ID("Timeout"): jen.Lit(5).Op("*").Qual("time", "Second"),
			}),
			jen.List(
				jen.ID("c"),
				jen.ID("err"),
			).Op(":=").ID("NewClient").Call(
				jen.ID("ctx"),
				jen.Lit(""),
				jen.Lit(""),
				jen.ID("address"),
				jen.ID("l"),
				jen.ID("h"),
				jen.Index().ID("string").Values(jen.Lit("*")),
				jen.ID("debug"),
			),
			jen.Return().List(
				jen.ID("c"),
				jen.ID("err"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("AuthenticatedClient provides the client's authenticated HTTP client"),
		jen.Line(),
		newClientMethod("AuthenticatedClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Block(
			jen.Return().ID("c").Dot("authedClient"),
		),
		jen.Line(),
	)

	// c.PlainClient
	ret.Add(
		jen.Comment("PlainClient provides the client's unauthenticated HTTP client"),
		jen.Line(),
		newClientMethod("PlainClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Block(
			jen.Return().ID("c").Dot("plainClient"),
		),
		jen.Line(),
	)

	// c.TokenSource
	ret.Add(
		jen.Comment("TokenSource provides the client's token source"),
		jen.Line(),
		newClientMethod("TokenSource").Params().Params(jen.ID("oauth2").Dot("TokenSource")).Block(
			jen.Return().ID("c").Dot("tokenSource"),
		),
		jen.Line(),
	)

	// c.executeRawRequest
	ret.Add(comments(
		"executeRawRequest takes a given *http.Request and executes it with the provided",
		"client, alongside some debugging logging.",
	)...)
	ret.Add(
		newClientMethod("executeRawRequest").Params(
			ctxParam(),
			jen.ID("client").Op("*").Qual("net/http", "Client"),
			jen.ID("req").Op("*").Qual("net/http", "Request"),
		).Params(jen.Op("*").Qual("net/http", "Response"),
			jen.ID("error")).Block(
			jen.Var().ID("logger").Op("=").ID("c").Dot("logger"),
			jen.If(jen.List(
				jen.ID("command"),
				jen.ID("err"),
			).Op(":=").Qual("github.com/moul/http2curl", "GetCurlCommand").Call(
				jen.ID("req"),
			),
				jen.ID("err").Op("==").ID("nil").Op("&&").ID("c").Dot("Debug"),
			).Block(
				jen.ID("logger").Op("=").ID("c").Dot("logger").Dot("WithValue").Call(
					jen.Lit("curl"),
					jen.ID("command").Dot("String").Call(),
				),
			),
			jen.Line(),
			jen.List(
				jen.ID("res"),
				jen.ID("err"),
			).Op(":=").ID("client").Dot("Do").Call(
				jen.ID("req").Dot("WithContext").Call(
					jen.ID("ctx"),
				),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("executing request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.ID("c").Dot("Debug")).Block(
				jen.List(
					jen.ID("bdump"),
					jen.ID("err"),
				).Op(":=").ID("httputil").Dot("DumpResponse").Call(
					jen.ID("res"),
					jen.ID("true"),
				),
				jen.If(jen.ID("err").Op("==").ID("nil").Op("&&").ID("req").Dot("Method").Op("!=").Qual("net/http", "MethodGet")).Block(
					jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
						jen.Lit("response_body"),
						jen.ID("string").Call(
							jen.ID("bdump"),
						),
					),
				),
				jen.ID("logger").Dot("Debug").Call(
					jen.Lit("request executed"),
				),
			),
			jen.Line(),
			jen.Return().List(jen.ID("res"),
				jen.ID("nil")),
		),
		jen.Line(),
	)

	// c.BuildURL
	ret.Add(
		jen.Comment("BuildURL builds standard service URLs"),
		jen.Line(),
		newClientMethod("BuildURL").Params(
			jen.ID("qp").Qual("net/url", "Values"),
			jen.ID("parts").Op("...").ID("string"),
		).Params(jen.ID("string")).Block(
			jen.If(jen.ID("qp").Op("!=").ID("nil")).Block(
				jen.Return().ID("c").Dot("buildURL").Call(
					jen.ID("qp"),
					jen.ID("parts").Op("..."),
				).Dot("String").Call(),
			),
			jen.Return().ID("c").Dot("buildURL").Call(
				jen.ID("nil"),
				jen.ID("parts").Op("..."),
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
		newClientMethod("buildURL").Params(
			jen.ID("queryParams").Qual("net/url", "Values"),
			jen.ID("parts").Op("...").ID("string"),
		).Params(jen.Op("*").Qual("net/url", "URL")).Block(
			jen.ID("tu").Op(":=").Op("*").ID("c").Dot("URL"),
			jen.Line(),
			jen.ID("parts").Op("=").ID("append").Call(
				jen.Index().ID("string").Values(jen.Lit("api"), jen.Lit("v1")),
				jen.ID("parts").Op("..."),
			),
			jen.List(
				jen.ID("u"),
				jen.ID("err"),
			).Op(":=").Qual("net/url", "Parse").Call(
				jen.Qual("strings", "Join").Call(
					jen.ID("parts"),
					jen.Lit("/"),
				),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("was asked to build an invalid URL: %v"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.ID("queryParams").Op("!=").ID("nil")).Block(
				jen.ID("u").Dot("RawQuery").Op("=").ID("queryParams").Dot("Encode").Call(),
			),
			jen.Line(),
			jen.Return().ID("tu").Dot("ResolveReference").Call(
				jen.ID("u"),
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
		newClientMethod("buildVersionlessURL").Params(
			jen.ID("qp").Qual("net/url", "Values"),
			jen.ID("parts").Op("...").ID("string"),
		).Params(jen.ID("string")).Block(
			jen.ID("tu").Op(":=").Op("*").ID("c").Dot("URL"),
			jen.Line(),
			jen.List(
				jen.ID("u"),
				jen.ID("err"),
			).Op(":=").Qual("net/url", "Parse").Call(
				jen.Qual("path", "Join").Call(
					jen.ID("parts").Op("..."),
				),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("user tried to build an invalid URL: %v"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.ID("qp").Op("!=").ID("nil")).Block(
				jen.ID("u").Dot("RawQuery").Op("=").ID("qp").Dot("Encode").Call(),
			),
			jen.Line(),
			jen.Return().ID("tu").Dot("ResolveReference").Call(
				jen.ID("u"),
			).Dot("String").Call(),
		),
		jen.Line(),
	)

	// c.BuildWebsocketURL
	ret.Add(comments("BuildWebsocketURL builds a standard URL and then converts its scheme to the websocket protocol")...)
	ret.Add(
		newClientMethod("BuildWebsocketURL").Params(
			jen.ID("parts").Op("...").ID("string"),
		).Params(jen.ID("string")).Block(
			jen.ID("u").Op(":=").ID("c").Dot("buildURL").Call(
				jen.ID("nil"),
				jen.ID("parts").Op("..."),
			),
			jen.ID("u").Dot("Scheme").Op("=").Lit("ws"),
			jen.Line(),
			jen.Return().ID("u").Dot("String").Call(),
		),
		jen.Line(),
	)

	// c.BuildHealthCheckRequest
	ret.Add(comments("BuildHealthCheckRequest builds a health check HTTP Request")...)
	ret.Add(
		newClientMethod("BuildHealthCheckRequest").Params().Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("u").Op(":=").Op("*").ID("c").Dot("URL"),
			jen.ID("uri").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s://%s/_meta_/ready"),
				jen.ID("u").Dot("Scheme"),
				jen.ID("u").Dot("Host"),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	// c.IsUp
	ret.Add(comments("IsUp returns whether or not the service's health endpoint is returning 200s")...)
	ret.Add(
		newClientMethod("IsUp").Params().Params(jen.ID("bool")).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildHealthCheckRequest").Call(),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("c").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("building request"),
				),
				jen.Return().ID("false"),
			),
			jen.Line(),
			jen.List(
				jen.ID("res"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("plainClient").Dot("Do").Call(
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("c").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("health check"),
				),
				jen.Return().ID("false"),
			),
			jen.Line(),
			jen.Defer().Func().Params().Block(
				jen.If(jen.ID("err").Op(":=").ID("res").Dot("Body").Dot("Close").Call(),
					jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("c").Dot("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Lit("closing response body"),
					),
				),
			).Call(),
			jen.Line(),
			jen.Return().ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusOK"),
		),
		jen.Line(),
	)

	// c.buildDataRequest
	ret.Add(comments("buildDataRequest builds an HTTP request for a given method, URL, and body data.")...)
	ret.Add(
		newClientMethod("buildDataRequest").Params(
			jen.List(
				jen.ID("method"),
				jen.ID("uri"),
			).ID("string"),
			jen.ID("in").Interface(),
		).Params(jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error")).Block(
			jen.List(
				jen.ID("body"),
				jen.ID("err"),
			).Op(":=").ID("createBodyFromStruct").Call(
				jen.ID("in"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.ID("err")),
			),
			jen.Line(),
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").Qual("net/http", "NewRequest").Call(
				jen.ID("method"),
				jen.ID("uri"),
				jen.ID("body"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.ID("req").Dot("Header").Dot("Set").Call(
				jen.Lit("Content-type"),
				jen.Lit("application/json"),
			),
			jen.Return().List(jen.ID("req"),
				jen.ID("nil")),
		),
		jen.Line(),
	)

	// c.retrieve
	ret.Add(comments("retrieve executes an HTTP request and loads the response content into a struct")...)
	ret.Add(
		newClientMethod("retrieve").Params(
			ctxParam(),
			jen.ID("req").Op("*").Qual("net/http", "Request"),
			jen.ID("obj").Interface(),
		).Params(jen.ID("error")).Block(
			jen.If(jen.ID("err").Op(":=").ID("argIsNotPointerOrNil").Call(
				jen.ID("obj"),
			),
				jen.ID("err").Op("!=").ID("nil"),
			).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("struct to load must be a pointer: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.List(
				jen.ID("res"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("executeRawRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("authedClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("executing request: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.If(jen.ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().ID("ErrNotFound"),
			),
			jen.Line(),
			jen.Return().ID("unmarshalBody").Call(
				jen.ID("res"), jen.Op("&").ID("obj"),
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
		newClientMethod("executeRequest").Params(
			ctxParam(),
			jen.ID("req").Op("*").Qual("net/http", "Request"),
			jen.ID("out").Interface(),
		).Params(jen.ID("error")).Block(
			jen.List(
				jen.ID("res"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("executeRawRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("authedClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("executing request: %w"),
					jen.ID("err"),
				),
			),
			jen.Switch(jen.ID("res").Dot("StatusCode")).Block(
				jen.Case(jen.Qual("net/http", "StatusNotFound")).Block(
					jen.Return().ID("ErrNotFound"),
				),
				jen.Case(jen.Qual("net/http", "StatusUnauthorized")).Block(
					jen.Return().ID("ErrUnauthorized"),
				),
			),
			jen.Line(),
			jen.If(jen.ID("out").Op("!=").ID("nil")).Block(
				jen.ID("resErr").Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"), jen.Op("&").ID("out"),
				),
				jen.If(jen.ID("resErr").Op("!=").ID("nil")).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("loading response from server: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	// c.executeUnathenticatedDataRequest
	ret.Add(comments("c.executeUnathenticatedDataRequest")...)
	ret.Add(
		newClientMethod("executeUnathenticatedDataRequest").Params(
			ctxParam(),
			jen.ID("req").Op("*").Qual("net/http", "Request"),
			jen.ID("out").Interface(),
		).Params(jen.ID("error")).Block(
			jen.If(jen.ID("out").Op("!=").ID("nil")).Block(
				jen.If(jen.List(
					jen.ID("np"),
					jen.ID("err"),
				).Op(":=").ID("argIsNotPointer").Call(
					jen.ID("out"),
				),
					jen.ID("np").Op("||").ID("err").Op("!=").ID("nil"),
				).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("struct to load must be a pointer: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.List(
				jen.ID("res"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("executeRawRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("plainClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("executing request: %w"),
					jen.ID("err"),
				),
			),
			jen.If(jen.ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().ID("ErrNotFound"),
			),
			jen.If(jen.ID("out").Op("!=").ID("nil")).Block(
				jen.ID("resErr").Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"), jen.Op("&").ID("out"),
				),
				jen.If(jen.ID("resErr").Op("!=").ID("nil")).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("loading response from server: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return ret
}
