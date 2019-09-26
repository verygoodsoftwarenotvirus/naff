package client

import (
	jen "github.com/dave/jennifer/jen"
)

var (
	// Files are all the available files to generate
	Files = map[string]*jen.File{
		"client/v1/http/main.go":                mainDotGo(),
		"client/v1/http/main_test.go":           mainTestDotGo(),
		"client/v1/http/helpers.go":             helpersDotGo(),
		"client/v1/http/helpers_test.go":        helpersTestDotGo(),
		"client/v1/http/users.go":               usersDotGo(),
		"client/v1/http/users_test.go":          usersTestDotGo(),
		"client/v1/http/roundtripper.go":        roundtripperDotGo(),
		"client/v1/http/webhooks.go":            webhooksDotGo(),
		"client/v1/http/webhooks_test.go":       webhooksTestDotGo(),
		"client/v1/http/oauth2_clients.go":      oauth2ClientsDotGo(),
		"client/v1/http/oauth2_clients_test.go": oauth2ClientsTestDotGo(),
	}
)

func mainDotGo() *jen.File {
	ret := jen.NewFile("client")
	ret.Add(jen.Null())

	ret.ImportNames(map[string]string{
		"context":           "context",
		"fmt":               "fmt",
		"net/http":          "http",
		"net/http/httputil": "httputil",
		"net/url":           "url",
		"path":              "path",
		"strings":           "strings",
		"time":              "time",
		//
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/logging/v1":      "logging",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/logging/v1/noop": "noop",
		//
		"github.com/moul/http2curl":             "http2curl",
		"github.com/pkg/errors":                 "errors",
		"go.opencensus.io/plugin/ochttp":        "ochttp",
		"golang.org/x/oauth2":                   "oauth2",
		"golang.org/x/oauth2/clientcredentials": "clientcredentials",
	})
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
			jen.Id("ErrNotFound").Op("=").Id("errors").Dot("New").Call(jen.Lit("404: not found")),
			jen.Id("ErrUnauthorized").Op("=").Id("errors").Dot("New").Call(jen.Lit("401: not authorized")),
		),
		jen.Line(),
	)

	// types
	ret.Add(
		jen.Type().Id("V1Client").Struct(
			jen.Id("plainClient").Op("*").Qual("net/http", "Client"),
			jen.Id("authedClient").Op("*").Qual("net/http", "Client"),
			jen.Id("logger").Qual("logging", "Logger"),
			jen.Id("Debug").Id("bool"),
			jen.Id("URL").Op("*").Qual("net/url", "URL"),
			jen.Id("Scopes").Index().Id("string"),
			jen.Id("tokenSource").Qual("oauth2", "TokenSource"),
		),
		jen.Line(),
	)

	// NewClient
	ret.Add(
		jen.Comment("NewClient builds a new API client for us"),
		jen.Line(),
		jen.Func().Id("NewClient").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.List(jen.Id("clientID"), jen.Id("clientSecret")).Id("string"),
			jen.Id("address").Op("*").Qual("net/url", "URL"),
			jen.Id("logger").Id("logging").Dot("Logger"),
			jen.Id("hclient").Op("*").Qual("net/http", "Client"),
			jen.Id("scopes").Index().Id("string"),
			jen.Id("debug").Id("bool"),
		).Params(jen.Op("*").Id("V1Client"), jen.Id("error")).Block(
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
				jen.Id("logger").Dot("SetLevel").Call(jen.Id("logging").Dot("DebugLevel")),
				jen.Id("logger").Dot("Debug").Call(jen.Lit("log level set to debug!")),
			),
			jen.Line(),
			jen.List(jen.Id("ac"), jen.Id("ts")).Op(":=").Id("buildOAuthClient").Call(
				jen.Id("ctx"),
				jen.Id("address"),
				jen.Id("clientID"),
				jen.Id("clientSecret"),
				jen.Id("scopes"),
			),
			jen.Line(),
			jen.Id("c").Op(":=").Op("&").Id("V1Client").Values(jen.Dict{
				jen.Id("URL"):          jen.Id("address"),
				jen.Id("plainClient"):  jen.Id("client"),
				jen.Id("logger"):       jen.Qual("logger", "WithName").Call(jen.Id("clientName")),
				jen.Id("Debug"):        jen.Id("debug"),
				jen.Id("authedClient"): jen.Id("ac"),
				jen.Id("tokenSource"):  jen.Id("ts"),
			}),
			jen.Line(),
			jen.Id("logger").Dot("WithValue").Call(jen.Lit("url"), jen.Id("address").Dot("String").Call()).Dot("Debug").Call(jen.Lit("returning client")),
			jen.Return().List(jen.Id("c"), jen.Id("nil")),
		),
		jen.Line(),
	)

	// buildOAuthClient
	ret.Add(
		jen.Comment("buildOAuthClient does too much"),
		jen.Line(),
		jen.Func().Id("buildOAuthClient").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("uri").Op("*").Qual("net/url", "URL"),
			jen.List(
				jen.Id("clientID"),
				jen.Id("clientSecret"),
			).Id("string"),
			jen.Id("scopes").Index().Id("string"),
		).Params(
			jen.Op("*").Qual("net/http", "Client"),
			jen.Id("oauth2").Dot("TokenSource")).Block(
			jen.Id("conf").Op(":=").Id("clientcredentials").Dot("Config").Values(
				jen.Id("ClientID").Op(":").Id("clientID"),
				jen.Id("ClientSecret").Op(":").Id("clientSecret"),
				jen.Id("Scopes").Op(":").Id("scopes"),
				jen.Id("EndpointParams").Op(":").Qual("net/url", "Values").Values(
					jen.Lit("client_id").Op(":").Index().Id("string").Values(jen.Id("clientID")),
					jen.Lit("client_secret").Op(":").Index().Id("string").Values(jen.Id("clientSecret")),
				),
				jen.Id("TokenURL").Op(":").Id("tokenEndpoint").Call(jen.Id("uri")).Dot("TokenURL"),
			),
			jen.Line(),
			jen.Id("ts").Op(":=").Id("oauth2").Dot("ReuseTokenSource").Call(
				jen.Id("nil"),
				jen.Id("conf").Dot("TokenSource").Call(jen.Id("ctx")),
			),
			jen.Id("client").Op(":=").Op("&").Qual("net/http", "Client").Values(
				jen.Id("Transport").Op(":").Op("&").Id("oauth2").Dot("Transport").Values(
					jen.Id("Base").Op(":").Op("&").Id("ochttp").Dot("Transport").Values(
						jen.Id("Base").Op(":").Id("newDefaultRoundTripper").Call(),
					),
					jen.Id("Source").Op(":").Id("ts"),
				),
				jen.Id("Timeout").Op(":").Lit(5).Op("*").Qual("time", "Second"),
			),
			jen.Line(),
			jen.Return().List(jen.Id("client"), jen.Id("ts")),
		),
		jen.Line(),
	)

	// tokenEndpoint
	ret.Add(
		jen.Comment("tokenEndpoint provides the oauth2 Endpoint for a given host"),
		jen.Line(),
		jen.Func().Id("tokenEndpoint").Params(
			jen.Id("baseURL").Op("*").Qual("net/url", "URL")).Params(
			jen.Id("oauth2").Dot("Endpoint")).Block(
			jen.List(jen.Id("tu"), jen.Id("au")).Op(":=").List(jen.Op("*").Id("baseURL"), jen.Op("*").Id("baseURL")),
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
	ret.Add(
		jen.Comment("NewSimpleClient is a client that is capable of much less than the normal client"),
		jen.Line(),
		jen.Comment("and has noops or empty values for most of its authentication and debug parts."),
		jen.Line(),
		jen.Comment("Its purpose at the time of this writing is merely so I can make users (which"),
		jen.Line(),
		jen.Comment("is a route that doesn't require authentication)"),
		jen.Line(),
		jen.Func().Id("NewSimpleClient").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("address").Op("*").Qual("net/url", "URL"), jen.Id("debug").Id("bool")).Params(
			jen.Op("*").Id("V1Client"), jen.Id("error")).Block(
			jen.Id("l").Op(":=").Id("noop").Dot("ProvideNoopLogger").Call(),
			jen.Id("h").Op(":=").Op("&").Qual("net/http", "Client").Values(jen.Id("Timeout").Op(":").Lit(5).Op("*").Qual("time", "Second")),
			jen.List(jen.Id("c"), jen.Id("err")).Op(":=").Id("NewClient").Call(
				jen.Id("ctx"),
				jen.Lit(""),
				jen.Lit(""),
				jen.Id("address"),
				jen.Id("l"),
				jen.Id("h"),
				jen.Index().Id("string").Values(jen.Lit("*")),
				jen.Id("debug"),
			),
			jen.Return().List(jen.Id("c"), jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("AuthenticatedClient provides the client's authenticated HTTP client"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("AuthenticatedClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Block(
			jen.Return().Id("c").Dot("authedClient"),
		),
		jen.Line(),
	)

	// c.PlainClient
	ret.Add(
		jen.Comment("PlainClient provides the client's unauthenticated HTTP client"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("PlainClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Block(
			jen.Return().Id("c").Dot("plainClient"),
		),
		jen.Line(),
	)

	// c.TokenSource
	ret.Add(
		jen.Comment("TokenSource provides the client's token source"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("TokenSource").Params().Params(jen.Id("oauth2").Dot("TokenSource")).Block(
			jen.Return().Id("c").Dot("tokenSource"),
		),
		jen.Line(),
	)

	// c.executeRawRequest
	ret.Add(
		jen.Comment("executeRawRequest takes a given *http.Request and executes it with the provided"),
		jen.Line(),
		jen.Comment("client, alongside some debugging logging."),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("executeRawRequest").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("client").Op("*").Qual("net/http", "Client"),
			jen.Id("req").Op("*").Qual("net/http", "Request"),
		).Params(jen.Op("*").Qual("net/http", "Response"), jen.Id("error")).Block(
			jen.Var().Id("logger").Op("=").Id("c").Dot("logger"),
			jen.If(jen.List(jen.Id("command"), jen.Id("err")).Op(":=").Qual("http2curl", "GetCurlCommand").Call(jen.Id("req")),
				jen.Id("err").Op("==").Id("nil").Op("&&").Id("c").Dot("Debug"),
			).Block(
				jen.Id("logger").Op("=").Id("c").Dot("logger").Dot("WithValue").Call(
					jen.Lit("curl"),
					jen.Id("command").Dot("String").Call(),
				),
			),
			jen.Line(),
			jen.List(jen.Id("res"), jen.Id("err")).Op(":=").Id("client").Dot("Do").Call(
				jen.Id("req").Dot("WithContext").Call(jen.Id("ctx")),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("executing request: %w"), jen.Id("err"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.Id("c").Dot("Debug")).Block(
				jen.List(jen.Id("bdump"), jen.Id("err")).Op(":=").Id("httputil").Dot("DumpResponse").Call(jen.Id("res"), jen.Id("true")),
				jen.If(jen.Id("err").Op("==").Id("nil").Op("&&").Id("req").Dot("Method").Op("!=").Qual("net/http", "MethodGet")).Block(
					jen.Id("logger").Op("=").Id("logger").Dot("WithValue").Call(jen.Lit("response_body"), jen.Id("string").Call(jen.Id("bdump"))),
				),
				jen.Id("logger").Dot("Debug").Call(jen.Lit("request executed")),
			),
			jen.Line(),
			jen.Return().List(jen.Id("res"), jen.Id("nil")),
		),
		jen.Line(),
	)

	// c.BuildURL
	ret.Add(
		jen.Comment("BuildURL builds standard service URLs"),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildURL").Params(
			jen.Id("qp").Qual("net/url", "Values"),
			jen.Id("parts").Op("...").Id("string"),
		).Params(jen.Id("string")).Block(
			jen.If(jen.Id("qp").Op("!=").Id("nil")).Block(
				jen.Return().Id("c").Dot("buildURL").Call(
					jen.Id("qp"), jen.Id("parts").Op("...")).Dot("String").Call(),
			),
			jen.Return().Id("c").Dot("buildURL").Call(
				jen.Id("nil"),
				jen.Id("parts").Op("...")).Dot("String").Call(),
		),
		jen.Line(),
	)

	// c.buildURL
	ret.Add(
		jen.Comment("buildURL takes a given set of query parameters and URL parts, and returns"),
		jen.Line(),
		jen.Comment("a parsed URL object from them."),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("buildURL").Params(
			jen.Id("queryParams").Qual("net/url", "Values"),
			jen.Id("parts").Op("...").Id("string"),
		).Params(jen.Op("*").Qual("net/url", "URL")).Block(
			jen.Id("tu").Op(":=").Op("*").Id("c").Dot("URL"),
			jen.Line(),
			jen.Id("parts").Op("=").Id("append").Call(jen.Index().Id("string").Values(jen.Lit("api"), jen.Lit("v1")),
				jen.Id("parts").Op("...")), jen.List(jen.Id("u"), jen.Id("err")).Op(":=").Qual("net/url", "Parse").Call(
				jen.Qual("strings", "Join").Call(
					jen.Id("parts"), jen.Lit("/"),
				),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("panic").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("was asked to build an invalid URL: %v"), jen.Id("err"),
				)),
			),
			jen.Line(),
			jen.If(jen.Id("queryParams").Op("!=").Id("nil")).Block(
				jen.Id("u").Dot("RawQuery").Op("=").Id("queryParams").Dot("Encode").Call(),
			),
			jen.Line(),
			jen.Return().Id("tu").Dot("ResolveReference").Call(jen.Id("u")),
		),
		jen.Line(),
	)

	// c.buildVersionlessURL
	ret.Add(
		jen.Comment("buildVersionlessURL builds a URL without the `/api/v1/` prefix. It should"),
		jen.Line(),
		jen.Comment("otherwise be identical to buildURL"),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("buildVersionlessURL").Params(
			jen.Id("qp").Qual("net/url", "Values"),
			jen.Id("parts").Op("...").Id("string"),
		).Params(jen.Id("string")).Block(
			jen.Id("tu").Op(":=").Op("*").Id("c").Dot("URL"),
			jen.Line(),
			jen.List(jen.Id("u"), jen.Id("err")).Op(":=").Qual("net/url", "Parse").Call(
				jen.Qual("path", "Join").Call(
					jen.Id("parts").Op("..."),
				),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("panic").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("user tried to build an invalid URL: %v"), jen.Id("err"))),
			),
			jen.Line(),
			jen.If(jen.Id("qp").Op("!=").Id("nil")).Block(
				jen.Id("u").Dot("RawQuery").Op("=").Id("qp").Dot("Encode").Call(),
			),
			jen.Line(),
			jen.Return().Id("tu").Dot("ResolveReference").Call(jen.Id("u")).Dot("String").Call(),
		),
		jen.Line(),
	)

	// c.BuildWebsocketURL
	ret.Add(
		jen.Comment("BuildWebsocketURL builds a standard URL and then converts its scheme to the websocket protocol"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("BuildWebsocketURL").Params(
			jen.Id("parts").Op("...").Id("string"),
		).Params(jen.Id("string")).Block(
			jen.Id("u").Op(":=").Id("c").Dot("buildURL").Call(jen.Id("nil"), jen.Id("parts").Op("...")),
			jen.Id("u").Dot("Scheme").Op("=").Lit("ws"),
			jen.Line(),
			jen.Return().Id("u").Dot("String").Call(),
		),
		jen.Line(),
	)

	// c.BuildHealthCheckRequest
	ret.Add(
		jen.Comment("BuildHealthCheckRequest builds a health check HTTP Request"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("BuildHealthCheckRequest").Params().Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("u").Op(":=").Op("*").Id("c").Dot("URL"),
			jen.Id("uri").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s://%s/_meta_/ready"),
				jen.Id("u").Dot("Scheme"), jen.Id("u").Dot("Host")),
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
	ret.Add(
		jen.Comment("IsUp returns whether or not the service's health endpoint is returning 200s"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("IsUp").Params().Params(jen.Id("bool")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildHealthCheckRequest").Call(),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("c").Dot("logger").Dot("Error").Call(jen.Id("err"), jen.Lit("building request")),
				jen.Return().Id("false"),
			),
			jen.Line(),
			jen.List(jen.Id("res"), jen.Id("err")).Op(":=").Id("c").Dot("plainClient").Dot("Do").Call(jen.Id("req")),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("c").Dot("logger").Dot("Error").Call(jen.Id("err"), jen.Lit("health check")),
				jen.Return().Id("false"),
			),
			jen.Line(),
			jen.Defer().Func().Params().Block(
				jen.If(jen.Id("err").Op(":=").Id("res").Dot("Body").Dot("Close").Call(), jen.Id("err").Op("!=").Id("nil")).Block(
					jen.Id("c").Dot("logger").Dot("Error").Call(jen.Id("err"),
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
	ret.Add(
		jen.Comment("buildDataRequest builds an HTTP request for a given method, URL, and body data."),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("buildDataRequest").Params(
			jen.List(jen.Id("method"), jen.Id("uri")).Id("string"),
			jen.Id("in").Interface(),
		).Params(jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.List(jen.Id("body"), jen.Id("err")).Op(":=").Id("createBodyFromStruct").Call(jen.Id("in")),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("err")),
			),
			jen.Line(),
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(
				jen.Id("method"),
				jen.Id("uri"),
				jen.Id("body"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Id("err"),
				)),
			jen.Line(),
			jen.Id("req").Dot("Header").Dot("Set").Call(
				jen.Lit("Content-type"),
				jen.Lit("application/json"),
			),
			jen.Return().List(jen.Id("req"), jen.Id("nil")),
		),
		jen.Line(),
	)

	// c.retrieve
	ret.Add(
		jen.Comment("retrieve executes an HTTP request and loads the response content into a struct"),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("retrieve").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("req").Op("*").Qual("net/http", "Request"),
			jen.Id("obj").Interface(),
		).Params(jen.Id("error")).Block(
			jen.If(jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(jen.Id("obj")), jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("struct to load must be a pointer: %w"), jen.Id("err")),
			),
			jen.Line(),
			jen.List(jen.Id("res"), jen.Id("err")).Op(":=").Id("c").Dot("executeRawRequest").Call(
				jen.Id("ctx"),
				jen.Id("c").Dot("authedClient"),
				jen.Id("req"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("executing request: %w"), jen.Id("err")),
			),
			jen.Line(),
			jen.If(jen.Id("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().Id("ErrNotFound"),
			),
			jen.Line(),
			jen.Return().Id("unmarshalBody").Call(jen.Id("res"), jen.Op("&").Id("obj")),
		),
		jen.Line(),
	)

	// c.executeRequest
	ret.Add(
		jen.Comment("executeRequest takes a given request and executes it with the auth client. It returns some errors"),
		jen.Line(),
		jen.Comment("upon receiving certain status codes, but otherwise will return nil upon success."),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("V1Client")).Id("executeRequest").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("req").Op("*").Qual("net/http", "Request"),
			jen.Id("out").Interface(),
		).Params(jen.Id("error")).Block(
			jen.List(jen.Id("res"), jen.Id("err")).Op(":=").Id("c").Dot("executeRawRequest").Call(
				jen.Id("ctx"),
				jen.Id("c").Dot("authedClient"),
				jen.Id("req"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("executing request: %w"), jen.Id("err")),
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
					jen.Id("res"), jen.Op("&").Id("out")),
				jen.If(jen.Id("resErr").Op("!=").Id("nil")).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("loading response from server: %w"), jen.Id("err")),
				),
			),
			jen.Line(),
			jen.Return().Id("nil"),
		),
		jen.Line(),
	)

	// c.executeUnathenticatedDataRequest
	ret.Add(
		jen.Comment("c.executeUnathenticatedDataRequest"),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("executeUnathenticatedDataRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("req").Op("*").Qual("net/http", "Request"), jen.Id("out").Interface()).Params(
			jen.Id("error")).Block(
			jen.If(jen.Id("out").Op("!=").Id("nil")).Block(
				jen.If(jen.List(jen.Id("np"), jen.Id("err")).Op(":=").Id("argIsNotPointer").Call(jen.Id("out")), jen.Id("np").Op("||").Id("err").Op("!=").Id("nil")).Block(
					jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("struct to load must be a pointer: %w"), jen.Id("err")))), jen.List(jen.Id("res"), jen.Id("err")).Op(":=").Id("c").Dot("executeRawRequest").Call(jen.Id("ctx"), jen.Id("c").Dot("plainClient"), jen.Id("req")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("executing request: %w"), jen.Id("err"))), jen.If(jen.Id("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().Id("ErrNotFound")), jen.If(jen.Id("out").Op("!=").Id("nil")).Block(
				jen.Id("resErr").Op(":=").Id("unmarshalBody").Call(jen.Id("res"), jen.Op("&").Id("out")), jen.If(jen.Id("resErr").Op("!=").Id("nil")).Block(
					jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("loading response from server: %w"), jen.Id("err")))), jen.Return().Id("nil"),
		),
		jen.Line(),
	)

	return ret
}

func mainTestDotGo() *jen.File {
	ret := jen.NewFile("client")
	ret.Add(jen.Null())

	ret.ImportNames(map[string]string{
		"github.com/stretchr/testify/assert":  "assert",
		"github.com/stretchr/testify/mock":    "mock",
		"github.com/stretchr/testify/require": "require",
	})

	// vars
	ret.Add(jen.Var().Id("exampleURI").Op("=").Lit("https://todo.verygoodsoftwarenotvirus.ru"))

	// types
	ret.Add(
		jen.Type().Id("argleBargle").Struct(
			jen.Id("Name").Id("string"),
		),
		jen.Line(),
	)
	ret.Add(jen.Type().Id("valuer").Map(jen.Id("string")).Index().Id("string"))

	ret.Add(
		jen.Func().Params(
			jen.Id("v").Id("valuer")).Id("ToValues").Params().Params(
			jen.Qual("net/url", "Values")).Block(
			jen.Return().Qual("net/url", "Values").Call(jen.Id("v")),
		),
		jen.Line(),
	)

	// funcs
	ret.Add(
		jen.Func().Id("mustParseURL").Params(
			jen.Id("uri").Id("string")).Params(
			jen.Op("*").Qual("net/url", "URL")).Block(
			jen.List(jen.Id("u"), jen.Id("err")).Op(":=").Qual("net/url", "Parse").Call(jen.Id("uri")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("panic").Call(jen.Id("err"))), jen.Return().Id("u"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("buildTestClient").Params(
			jen.Id("t").Op("*").Qual("testing", "T"), jen.Id("ts").Op("*").Id("httptest").Dot("Server")).Params(
			jen.Op("*").Id("V1Client")).Block(
			jen.Id("t").Dot("Helper").Call(),
			jen.Id("l").Op(":=").Id("noop").Dot("ProvideNoopLogger").Call(),
			jen.Id("u").Op(":=").Id("mustParseURL").Call(jen.Id("ts").Dot("URL")),
			jen.Return().Op("&").Id("V1Client").Values(jen.Dict{
				jen.Id("URL"):          jen.Id("u"),
				jen.Id("plainClient"):  jen.Id("ts").Dot("Client").Call(),
				jen.Id("logger"):       jen.Id("l"),
				jen.Id("Debug"):        jen.Id("true"),
				jen.Id("authedClient"): jen.Id("ts").Dot("Client").Call(),
			}),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_AuthenticatedClient").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Id("actual").Op(":=").Id("c").Dot("AuthenticatedClient").Call(),
				jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("ts").Dot("Client").Call(), jen.Id("actual"), jen.Lit("AuthenticatedClient should return the assigned authedClient")),
			),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_PlainClient").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Id("actual").Op(":=").Id("c").Dot("PlainClient").Call(),
				jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("ts").Dot("Client").Call(), jen.Id("actual"), jen.Lit("PlainClient should return the assigned plainClient")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_TokenSource").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(
				jen.Lit("obligatory"), jen.Func().Params(
					jen.Id("t").Op("*").Qual("testing", "T")).Block(
					jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")),
					jen.List(jen.Id("c"), jen.Id("err")).Op(":=").Id("NewClient").Call(
						jen.Qual("context", "Background").Call(),
						jen.Lit(""),
						jen.Lit(""),
						jen.Id("mustParseURL").Call(jen.Id("exampleURI")),
						jen.Id("noop").Dot("ProvideNoopLogger").Call(),
						jen.Id("ts").Dot("Client").Call(),
						jen.Index().Id("string").Values(jen.Lit("*")), jen.Id("false"),
					),
					jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
					jen.Id("actual").Op(":=").Id("c").Dot("TokenSource").Call(),
					jen.Id("assert").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestNewClient").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")),
				jen.List(jen.Id("c"), jen.Id("err")).Op(":=").Id("NewClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.Id("mustParseURL").Call(jen.Id("exampleURI")),
					jen.Id("noop").Dot("ProvideNoopLogger").Call(),
					jen.Id("ts").Dot("Client").Call(),
					jen.Index().Id("string").Values(jen.Lit("*")), jen.Id("false")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("c")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")))),

			jen.Id("T").Dot("Run").Call(jen.Lit("with client but invalid timeout"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.Id("c"), jen.Id("err")).Op(":=").Id("NewClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.Id("mustParseURL").Call(jen.Id("exampleURI")),
					jen.Id("noop").Dot("ProvideNoopLogger").Call(),
					jen.Op("&").Qual("net/http", "Client").Values(
						jen.Id("Timeout").Op(":").Lit(0),
					),
					jen.Index().Id("string").Values(jen.Lit("*")), jen.Id("true"),
				),

				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("c")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("assert").Dot("Equal").Call(
					jen.Id("t"),
					jen.Id("c").Dot("plainClient").Dot("Timeout"),
					jen.Id("defaultTimeout"),
					jen.Lit("NewClient should set the default timeout"),
				),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestNewSimpleClient").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.Id("c"), jen.Id("err")).Op(":=").Id("NewSimpleClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Id("mustParseURL").Call(jen.Id("exampleURI")), jen.Id("true"),
				),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("c")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_executeRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("with error"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
						jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Id("expectedMethod")),
						jen.Qual("time", "Sleep").Call(jen.Lit(10).Op("*").Qual("time", "Hour"))),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Id("expectedMethod"), jen.Id("ts").Dot("URL"), jen.Id("nil")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")), jen.List(jen.Id("res"), jen.Id("err")).Op(":=").Id("c").Dot("executeRawRequest").Call(jen.Id("ctx"), jen.Op("&").Qual("net/http", "Client").Values(jen.Id("Timeout").Op(":").Qual("time", "Second")), jen.Id("req")), jen.Id("assert").Dot("Nil").Call(jen.Id("t"), jen.Id("res")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestBuildURL").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("various urls"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("t").Dot("Parallel").Call(),
				jen.List(jen.Id("u"), jen.Id("_")).Op(":=").Qual("net/url", "Parse").Call(jen.Id("exampleURI")),
				jen.List(jen.Id("c"), jen.Id("err")).Op(":=").Id("NewClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""), jen.Lit(""),
					jen.Id("u"),
					jen.Id("noop").Dot("ProvideNoopLogger").Call(),
					jen.Id("nil"),
					jen.Index().Id("string").Values(jen.Lit("*")), jen.Id("false"),
				),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("testCases").Op(":=").Index().Struct(
					jen.Id("expectation").Id("string"),
					jen.Id("inputParts").Index().Id("string"),
					jen.Id("inputQuery").Id("valuer"),
				).Values(
					jen.Values(jen.Dict{
						jen.Id("expectation"): jen.Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things"),
						jen.Id("inputParts"):  jen.Index().Id("string").Values(jen.Lit("things")),
					}),
					jen.Values(jen.Dict{
						jen.Id("expectation"): jen.Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/stuff?key=value"),
						jen.Id("inputQuery"): jen.Map(jen.Id("string")).Index().Id("string").Values(jen.Dict{
							jen.Lit("key"): jen.Values(jen.Lit("value")),
						}),
						jen.Id("inputParts"): jen.Index().Id("string").Values(jen.Lit("stuff")),
					}),
					jen.Values(jen.Id("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff?key=value1&key=value2&yek=eulav"),
						jen.Id("inputQuery").Op(":").Map(jen.Id("string")).Index().Id("string").Values(
							jen.Lit("key").Op(":").Values(
								jen.Lit("value1"),
								jen.Lit("value2")),
							jen.Lit("yek").Op(":").Values(jen.Lit("eulav")),
						),
						jen.Id("inputParts").Op(":").Index().Id("string").Values(jen.Lit("things"), jen.Lit("and"), jen.Lit("stuff")),
					),
				),
				jen.For(jen.List(jen.Id("_"), jen.Id("tc")).Op(":=").Range().Id("testCases")).Block(
					jen.Id("actual").Op(":=").Id("c").Dot("BuildURL").Call(
						jen.Id("tc").Dot("inputQuery").Dot("ToValues").Call(),
						jen.Id("tc").Dot("inputParts").Op("..."),
					),
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("tc").Dot("expectation"), jen.Id("actual")),
				),
			),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildWebsocketURL").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.Id("u"), jen.Id("_")).Op(":=").Qual("net/url", "Parse").Call(jen.Id("exampleURI")),
				jen.List(jen.Id("c"), jen.Id("err")).Op(":=").Id("NewClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.Id("u"),
					jen.Id("noop").Dot("ProvideNoopLogger").Call(),
					jen.Id("nil"),
					jen.Index().Id("string").Values(jen.Lit("*")),
					jen.Id("false"),
				),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("expected").Op(":=").Lit("ws://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff"),
				jen.Id("actual").Op(":=").Id("c").Dot("BuildWebsocketURL").Call(
					jen.Lit("things"),
					jen.Lit("and"),
					jen.Lit("stuff"),
				),
				jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual"))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildHealthCheckRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodGet"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildHealthCheckRequest").Call(),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")),
				jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_IsUp").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
						jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")),
						jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Id("actual").Op(":=").Id("c").Dot("IsUp").Call(),
				jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Id("actual"))),
			),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("with bad status code"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")),
					jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError"))),
				)),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Id("actual").Op(":=").Id("c").Dot("IsUp").Call(),
				jen.Id("assert").Dot("False").Call(jen.Id("t"), jen.Id("actual")),
			)),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(
				jen.Lit("with timeout"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
					jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
						jen.Qual("net/http", "HandlerFunc").Call(
							jen.Func().Params(jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
								jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")),
								jen.Qual("time", "Sleep").Call(jen.Lit(10).Op("*").Qual("time", "Hour"))),
						),
					),
					jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
					jen.Id("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
					jen.Id("actual").Op(":=").Id("c").Dot("IsUp").Call(),
					jen.Id("assert").Dot("False").Call(jen.Id("t"), jen.Id("actual"))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_buildDataRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(
				jen.Lit("happy path"), jen.Func().Params(
					jen.Id("t").Op("*").Qual("testing", "T")).Block(
					jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")),
					jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
					jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
					jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("buildDataRequest").Call(
						jen.Id("expectedMethod"),
						jen.Id("ts").Dot("URL"),
						jen.Op("&").Id("testingType").Values(jen.Id("Name").Op(":").Lit("name")),
					),
					jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
					jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expectedMethod"), jen.Id("req").Dot("Method"))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_makeRequest").Params(jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
							jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Id("expectedMethod")),
							jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(
								jen.Op("&").Id("argleBargle").Values(jen.Id("Name").Op(":").Lit("name"))),
							),
						),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Id("expectedMethod"), jen.Id("ts").Dot("URL"), jen.Id("nil")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("argleBargle").Values()),
				jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
			)),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("with 404"), jen.Func().Params(
				jen.Line(),
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
							jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Id("expectedMethod")),
							jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
						),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Id("expectedMethod"), jen.Id("ts").Dot("URL"), jen.Id("nil")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Line(),
				jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("ErrNotFound"), jen.Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("argleBargle").Values())),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_makeUnauthedDataRequest").Params(jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(jen.Id("res").Qual("net/http", "ResponseWriter"),
							jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
							jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Id("expectedMethod")),
							jen.Id("require").Dot("NoError").Call(
								jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(
									jen.Id("res")).Dot("Encode").Call(
									jen.Op("&").Id("argleBargle").Values(
										jen.Id("Name").Op(":").Lit("name"),
									),
								),
							),
						),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Line(),
				jen.List(jen.Id("in"), jen.Id("out")).Op(":=").List(jen.Op("&").Id("argleBargle").Values(), jen.Op("&").Id("argleBargle").Values()),
				jen.Line(),
				jen.List(jen.Id("body"), jen.Id("err")).Op(":=").Id("createBodyFromStruct").Call(jen.Id("in")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("body")),
				jen.Line(),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Id("expectedMethod"), jen.Id("ts").Dot("URL"), jen.Id("body")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("executeUnathenticatedDataRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Id("out")),
				jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
			)),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("with 404"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.Id("res").Qual("net/http", "ResponseWriter"),
							jen.Id("req").Op("*").Qual("net/http", "Request"),
						).Block(
							jen.Id("assert").Dot("Equal").Call(
								jen.Id("t"), jen.Id("req").Dot("Method"), jen.Id("expectedMethod")),
							jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
						),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Line(),
				jen.List(jen.Id("in"), jen.Id("out")).Op(":=").List(jen.Op("&").Id("argleBargle").Values(), jen.Op("&").Id("argleBargle").Values()),
				jen.Line(),
				jen.List(jen.Id("body"), jen.Id("err")).Op(":=").Id("createBodyFromStruct").Call(jen.Id("in")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("body")),
				jen.Line(),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("body"),
				),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Id("out"),
				),
				jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("ErrNotFound"), jen.Id("err")),
			)),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("with timeout"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
							jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Id("expectedMethod")),
							jen.Qual("time", "Sleep").Call(jen.Lit(10).Op("*").Qual("time", "Hour")),
						),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Line(),
				jen.List(jen.Id("in"), jen.Id("out")).Op(":=").List(jen.Op("&").Id("argleBargle").Values(), jen.Op("&").Id("argleBargle").Values()),
				jen.Line(),
				jen.List(jen.Id("body"), jen.Id("err")).Op(":=").Id("createBodyFromStruct").Call(jen.Id("in")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("body")),
				jen.Line(),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("body"),
				),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
				jen.Line(),
				jen.Id("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.Id("assert").Dot("Error").Call(
					jen.Id("t"),
					jen.Id("c").Dot("executeUnathenticatedDataRequest").Call(
						jen.Id("ctx"),
						jen.Id("req"),
						jen.Id("out"),
					),
				),
			)),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("with nil as output"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Line(),
				jen.Id("in").Op(":=").Op("&").Id("argleBargle").Values(),
				jen.Line(),
				jen.List(jen.Id("body"), jen.Id("err")).Op(":=").Id("createBodyFromStruct").Call(jen.Id("in")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("body")),
				jen.Line(),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("body"),
				),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Id("testingType").Values(),
				),
				jen.Id("assert").Dot("Error").Call(
					jen.Id("t"),
					jen.Id("err"),
				),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_retrieve").Params(jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").
					Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
						jen.Id("assert").Dot("Equal").Call(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Id("expectedMethod"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"),
							jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res")).Dot("Encode").Call(
								jen.Op("&").Id("argleBargle").Values(jen.Dict{
									jen.Id("Name"): jen.Lit("name"),
								}),
							),
						),
					))),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Line(),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Op("&").Id("argleBargle").Values(),
				),
				jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
			)),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("with nil passed in"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Line(),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				jen.Id("require").Dot("NotNil").Call(
					jen.Id("t"),
					jen.Id("req"),
				),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Id("nil")),
				jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")),
			)),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("with timeout"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
						jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Id("expectedMethod")),
						jen.Qual("time", "Sleep").Call(jen.Lit(10).Op("*").Qual("time", "Hour")),
					)),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")),
				jen.Line(),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Line(),
				jen.Id("c").Dot("authedClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("argleBargle").Values()),
				jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")),
			)),
			jen.Line(),
			jen.Id("T").Dot("Run").Call(jen.Lit("with 404"), jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.Id("res").Qual("net/http", "ResponseWriter"),
							jen.Id("req").Op("*").Qual("net/http", "Request"),
						).Block(
							jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Id("expectedMethod")),
							jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
						),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Line(),
				jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"), jen.Id("nil")),
				jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")),
				jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")),
				jen.Line(),
				jen.Id("assert").Dot("Equal").Call(
					jen.Id("t"),
					jen.Id("ErrNotFound"),
					jen.Id("c").Dot("retrieve").Call(
						jen.Id("ctx"),
						jen.Id("req"),
						jen.Op("&").Id("argleBargle").Values(),
					),
				),
			)),
		),
	)

	return ret
}

func helpersDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.ImportNames(map[string]string{
		"bytes":         "bytes",
		"fmt":           "fmt",
		"encoding/json": "json",
		"io":            "io",
		"io/ioutil":     "ioutil",
		"net/http":      "http",
		"reflect":       "reflect",
	})
	ret.Add(jen.Line())
	ret.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "models")
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment("argIsNotPointer checks an argument and returns whether or not it is a pointer"),
		jen.Line(),
		jen.Func().Id("argIsNotPointer").Params(jen.Id("i").Interface()).Params(jen.Id("notAPointer").Id("bool"), jen.Id("err").Id("error")).Block(
			jen.If(jen.Id("i").Op("==").Id("nil").Op("||").Qual("reflect", "TypeOf").Call(jen.Id("i")).Dot("Kind").Call().Op("!=").Qual("reflect", "Ptr")).Block(
				jen.Return().List(jen.Id("true"), jen.Qual("errors", "New").Call(jen.Lit("value is not a pointer"))),
			),
			jen.Return().List(jen.Id("false"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("argIsNotNil checks an argument and returns whether or not it is nil"),
		jen.Line(),
		jen.Func().Id("argIsNotNil").Params(jen.Id("i").Interface()).Params(jen.Id("isNil").Id("bool"), jen.Id("err").Id("error")).Block(
			jen.If(jen.Id("i").Op("==").Id("nil")).Block(
				jen.Return().List(jen.Id("true"), jen.Id("errors").Dot("New").Call(jen.Lit("value is nil")))), jen.Return().List(jen.Id("false"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("argIsNotPointerOrNil does what it says on the tin. This function"),
		jen.Line(),
		jen.Comment("is primarily useful for detecting if a destination value is valid"),
		jen.Line(),
		jen.Comment("before decoding an HTTP response, for instance"),
		jen.Line(),
		jen.Func().Id("argIsNotPointerOrNil").Params(jen.Id("i").Interface()).Params(jen.Id("error")).Block(
			jen.If(
				jen.List(jen.Id("nn"), jen.Id("err")).Op(":=").Id("argIsNotNil").Call(jen.Id("i")),
				jen.Id("nn").Op("||").Id("err").Op("!=").Id("nil"),
			).Block(jen.Return().Id("err")),
			jen.If(jen.List(jen.Id("np"), jen.Id("err")).Op(":=").Id("argIsNotPointer").Call(jen.Id("i")), jen.Id("np").Op("||").Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Id("err")), jen.Return().Id("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("unmarshalBody takes an HTTP response and JSON decodes its"),
		jen.Line(),
		jen.Comment("body into a destination value. `dest` must be a non-nil"),
		jen.Line(),
		jen.Comment("pointer to an object. Ideally, response is also not nil."),
		jen.Line(),
		jen.Comment("The error returned here should only ever be received in"),
		jen.Line(),
		jen.Comment("testing, and should never be encountered by an end-user."),
		jen.Line(),
		jen.Func().Id("unmarshalBody").Params(
			jen.Id("res").Op("*").Qual("net/http", "Response"), jen.Id("dest").Interface()).Params(
			jen.Id("error")).Block(
			jen.If(jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(jen.Id("dest")), jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Id("err")), jen.List(jen.Id("bodyBytes"), jen.Id("err")).Op(":=").Qual("io/ioutil", "ReadAll").Call(jen.Id("res").Dot("Body")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Id("err")), jen.If(jen.Id("res").Dot("StatusCode").Op(">=").Qual("net/http", "StatusBadRequest")).Block(
				jen.Id("apiErr").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "ErrorResponse").Values(), jen.If(jen.Id("err").Op("=").Qual("encoding/json", "Unmarshal").Call(jen.Id("bodyBytes"), jen.Op("&").Id("apiErr")), jen.Id("err").Op("!=").Id("nil")).Block(
					jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("unmarshaling error: %w"), jen.Id("err")),
				),
				jen.Return().Id("apiErr")), jen.If(jen.Id("err").Op("=").Qual("encoding/json", "Unmarshal").Call(jen.Id("bodyBytes"), jen.Op("&").Id("dest")), jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("unmarshaling body: %w"), jen.Id("err"))), jen.Return().Id("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("createBodyFromStruct takes any value in and returns an io.Reader"),
		jen.Line(),
		jen.Comment("for placement within http.NewRequest's last argument."),
		jen.Line(),
		jen.Func().Id("createBodyFromStruct").Params(
			jen.Id("in").Interface()).Params(
			jen.Qual("io", "Reader"), jen.Id("error")).Block(
			jen.List(jen.Id("out"), jen.Id("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id("in")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("err"))), jen.Return().List(jen.Qual("bytes", "NewReader").Call(jen.Id("out")), jen.Id("nil")),
		),
		jen.Line(),
	)

	return ret
}

func helpersTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())
	ret.ImportNames(map[string]string{
		"github.com/stretchr/testify/assert":  "assert",
		"github.com/stretchr/testify/mock":    "mock",
		"github.com/stretchr/testify/require": "require",
	})

	ret.Add(
		jen.Type().Id("testingType").Struct(
			jen.Id("Name").Id("string"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestArgIsNotPointerOrNil").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("expected use"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(jen.Op("&").Id("testingType").Values()), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should not be returned when a pointer is provided")))), jen.Id("T").Dot("Run").Call(jen.Lit("with non-pointer"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(jen.Id("testingType").Values()), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should be returned when a non-pointer is provided")))), jen.Id("T").Dot("Run").Call(jen.Lit("with nil"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(jen.Id("nil")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should be returned when nil is provided")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestArgIsNotPointer").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("expected use"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.Id("notAPointer"), jen.Id("err")).Op(":=").Id("argIsNotPointer").Call(jen.Op("&").Id("testingType").Values()), jen.Id("assert").Dot("False").Call(jen.Id("t"), jen.Id("notAPointer"), jen.Lit("expected `false` when a pointer is provided")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should not be returned when a pointer is provided")))), jen.Id("T").Dot("Run").Call(jen.Lit("with non-pointer"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.Id("notAPointer"), jen.Id("err")).Op(":=").Id("argIsNotPointer").Call(jen.Id("testingType").Values()), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Id("notAPointer"), jen.Lit("expected `true` when a non-pointer is provided")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should be returned when a non-pointer is provided")))), jen.Id("T").Dot("Run").Call(jen.Lit("with nil"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.Id("notAPointer"), jen.Id("err")).Op(":=").Id("argIsNotPointer").Call(jen.Id("nil")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Id("notAPointer"), jen.Lit("expected `true` when nil is provided")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should be returned when nil is provided")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestArgIsNotNil").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("without nil"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.Id("isNil"), jen.Id("err")).Op(":=").Id("argIsNotNil").Call(jen.Op("&").Id("testingType").Values()), jen.Id("assert").Dot("False").Call(jen.Id("t"), jen.Id("isNil"), jen.Lit("expected `false` when a pointer is provided")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should not be returned when a pointer is provided")))), jen.Id("T").Dot("Run").Call(jen.Lit("with non-pointer"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.Id("isNil"), jen.Id("err")).Op(":=").Id("argIsNotNil").Call(jen.Id("testingType").Values()), jen.Id("assert").Dot("False").Call(jen.Id("t"), jen.Id("isNil"), jen.Lit("expected `true` when a non-pointer is provided")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should not be returned when a non-pointer is provided")))), jen.Id("T").Dot("Run").Call(jen.Lit("with nil"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.Id("isNil"), jen.Id("err")).Op(":=").Id("argIsNotNil").Call(jen.Id("nil")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Id("isNil"), jen.Lit("expected `true` when nil is provided")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should be returned when nil is provided")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestUnmarshalBody").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("expected use"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Lit("example"), jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(jen.Id("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(jen.Qual("strings", "NewReader").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(`{"name": %q}`), jen.Id("expected")))), jen.Id("StatusCode").Op(":").Qual("net/http", "StatusOK")), jen.Var().Id("out").Id("testingType"), jen.Id("err").Op(":=").Id("unmarshalBody").Call(jen.Id("res"), jen.Op("&").Id("out")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("out").Dot("Name"), jen.Id("expected"), jen.Lit("expected marshaling to work")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be encountered unmarshaling into a valid struct")))), jen.Id("T").Dot("Run").Call(jen.Lit("with good status but unmarshallable response"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(jen.Id("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(jen.Qual("strings", "NewReader").Call(jen.Lit(`

					BLAH

				`))), jen.Id("StatusCode").Op(":").Qual("net/http", "StatusOK")), jen.Var().Id("out").Id("testingType"), jen.Id("err").Op(":=").Id("unmarshalBody").Call(jen.Id("res"), jen.Op("&").Id("out")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be encountered unmarshaling into a valid struct")))), jen.Id("T").Dot("Run").Call(jen.Lit("with an erroneous error code"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(jen.Id("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(jen.Qual("strings", "NewReader").Call(jen.Func().Params().Params(
					jen.Id("string")).Block(
					jen.Id("er").Op(":=").Op("&").Id("models").Dot("ErrorResponse").Values(), jen.List(jen.Id("bs"), jen.Id("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id("er")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Id("err")), jen.Return().Id("string").Call(jen.Id("bs"))).Call())), jen.Id("StatusCode").Op(":").Qual("net/http", "StatusBadRequest")), jen.Var().Id("out").Op("*").Id("testingType"), jen.Id("err").Op(":=").Id("unmarshalBody").Call(jen.Id("res"), jen.Op("&").Id("out")), jen.Id("assert").Dot("Nil").Call(jen.Id("t"), jen.Id("out"), jen.Lit("expected nil to be returned")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should be returned from the API")))), jen.Id("T").Dot("Run").Call(jen.Lit("with an erroneous error code and unmarshallable body"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(jen.Id("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(jen.Qual("strings", "NewReader").Call(jen.Lit(`

				BLAH

				`))), jen.Id("StatusCode").Op(":").Qual("net/http", "StatusBadRequest")), jen.Var().Id("out").Op("*").Id("testingType"), jen.Id("err").Op(":=").Id("unmarshalBody").Call(jen.Id("res"), jen.Op("&").Id("out")), jen.Id("assert").Dot("Nil").Call(jen.Id("t"), jen.Id("out"), jen.Lit("expected nil to be returned")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should be returned from the unmarshaller")))), jen.Id("T").Dot("Run").Call(jen.Lit("with nil target variable"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("err").Op(":=").Id("unmarshalBody").Call(jen.Id("nil"), jen.Id("nil")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("error should be encountered when passed nil")))), jen.Id("T").Dot("Run").Call(jen.Lit("with erroneous reader"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")), jen.Id("rc").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/mock", "NewMockReadCloser").Call(), jen.Id("rc").Dot("On").Call(jen.Lit("Read"), jen.Id("mock").Dot("Anything")).Dot("Return").Call(jen.Lit(0), jen.Id("expected")), jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(jen.Id("Body").Op(":").Id("rc"), jen.Id("StatusCode").Op(":").Qual("net/http", "StatusOK")), jen.Var().Id("out").Id("testingType"), jen.Id("err").Op(":=").Id("unmarshalBody").Call(jen.Id("res"), jen.Op("&").Id("out")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("err")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be encountered unmarshaling into a valid struct")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Id("testBreakableStruct").Struct(
			jen.Id("Thing").Qual("encoding/json", "Number"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestCreateBodyFromStruct").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("expected use"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Lit(`{"name":"expected"}`), jen.Id("x").Op(":=").Op("&").Id("testingType").Values(jen.Id("Name").Op(":").Lit("expected")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("createBodyFromStruct").Call(jen.Id("x")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("expected no error creating JSON from valid struct")), jen.List(jen.Id("bs"), jen.Id("err")).Op(":=").Qual("io/ioutil", "ReadAll").Call(jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("expected no error reading JSON from valid struct")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("string").Call(jen.Id("bs")), jen.Lit("expected and actual JSON bodies don't match")))), jen.Id("T").Dot("Run").Call(jen.Lit("with unmarshallable struct"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("x").Op(":=").Op("&").Id("testBreakableStruct").Values(jen.Id("Thing").Op(":").Lit("stuff")), jen.List(jen.Id("_"), jen.Id("err")).Op(":=").Id("createBodyFromStruct").Call(jen.Id("x")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err"), jen.Lit("expected no error creating JSON from valid struct")))),
		),
		jen.Line(),
	)
	return ret
}

func oauth2ClientsDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())
	ret.Add(jen.Var().Id("oauth2ClientsBasePath").Op("=").Lit("oauth2/clients"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildGetOAuth2ClientRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("oauth2ClientsBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("id"), jen.Lit(10))), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Id("uri"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("GetOAuth2Client").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Id("oauth2Client").Op("*").Id("models").Dot("OAuth2Client"), jen.Id("err").Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetOAuth2ClientRequest").Call(jen.Id("ctx"), jen.Id("id")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("oauth2Client")), jen.Return().List(jen.Id("oauth2Client"), jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildGetOAuth2ClientsRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("filter").Op("*").Id("models").Dot("QueryFilter")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("filter").Dot("ToValues").Call(), jen.Id("oauth2ClientsBasePath")), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Id("uri"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("GetOAuth2Clients").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("filter").Op("*").Id("models").Dot("QueryFilter")).Params(
			jen.Op("*").Id("models").Dot("OAuth2ClientList"), jen.Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetOAuth2ClientsRequest").Call(jen.Id("ctx"), jen.Id("filter")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.Var().Id("oauth2Clients").Op("*").Id("models").Dot("OAuth2ClientList"), jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("oauth2Clients")), jen.Return().List(jen.Id("oauth2Clients"), jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildCreateOAuth2ClientRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("cookie").Op("*").Qual("net/http", "Cookie"), jen.Id("body").Op("*").Id("models").Dot("OAuth2ClientCreationInput")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(jen.Id("nil"), jen.Lit("oauth2"), jen.Lit("client")), jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("buildDataRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Id("uri"), jen.Id("body")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("err"))), jen.Id("req").Dot("AddCookie").Call(jen.Id("cookie")), jen.Return().List(jen.Id("req"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("CreateOAuth2Client").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("cookie").Op("*").Qual("net/http", "Cookie"), jen.Id("input").Op("*").Id("models").Dot("OAuth2ClientCreationInput")).Params(
			jen.Op("*").Id("models").Dot("OAuth2Client"), jen.Id("error")).Block(
			jen.Var().Id("oauth2Client").Op("*").Id("models").Dot("OAuth2Client"), jen.If(jen.Id("cookie").Op("==").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("errors").Dot("New").Call(jen.Lit("cookie required for request")))), jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildCreateOAuth2ClientRequest").Call(jen.Id("ctx"), jen.Id("cookie"), jen.Id("input")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("err"))), jen.List(jen.Id("res"), jen.Id("err")).Op(":=").Id("c").Dot("executeRawRequest").Call(jen.Id("ctx"), jen.Id("c").Dot("plainClient"), jen.Id("req")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("executing request: %w"), jen.Id("err")))), jen.If(jen.Id("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("ErrNotFound"))), jen.If(jen.Id("resErr").Op(":=").Id("unmarshalBody").Call(jen.Id("res"), jen.Op("&").Id("oauth2Client")), jen.Id("resErr").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("errors").Dot("Wrap").Call(jen.Id("resErr"), jen.Lit("loading response from server")))), jen.Return().List(jen.Id("oauth2Client"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildArchiveOAuth2ClientRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("oauth2ClientsBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("id"), jen.Lit(10))), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodDelete"), jen.Id("uri"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("ArchiveOAuth2Client").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildArchiveOAuth2ClientRequest").Call(jen.Id("ctx"), jen.Id("id")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")),
			),
			jen.Return().Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Id("nil")),
		),
		jen.Line(),
	)

	return ret
}

func oauth2ClientsTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())
	ret.ImportNames(map[string]string{
		"github.com/stretchr/testify/assert":  "assert",
		"github.com/stretchr/testify/mock":    "mock",
		"github.com/stretchr/testify/require": "require",
	})

	ret.Add(
		jen.Func().Id("TestV1Client_BuildGetOAuth2ClientRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodGet"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.Id("expectedID").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetOAuth2ClientRequest").Call(jen.Id("ctx"), jen.Id("expectedID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("actual").Dot("URL").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.Id("expectedID")))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_GetOAuth2Client").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("OAuth2Client").Values(jen.Id("ID").Op(":").Lit(1), jen.Id("ClientID").Op(":").Lit("example"), jen.Id("ClientSecret").Op(":").Lit("blah")), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("req").Dot("URL").Dot("String").Call(), jen.Qual("strconv", "Itoa").Call(jen.Id("int").Call(jen.Id("expected").Dot("ID"))))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/api/v1/oauth2/clients/%d"), jen.Id("expected").Dot("ID")), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected")))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("GetOAuth2Client").Call(jen.Id("ctx"), jen.Id("expected").Dot("ID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildGetOAuth2ClientsRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodGet"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetOAuth2ClientsRequest").Call(jen.Id("ctx"), jen.Id("nil")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_GetOAuth2Clients").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("OAuth2ClientList").Values(jen.Id("Clients").Op(":").Index().Id("models").Dot("OAuth2Client").Values(jen.Values(jen.Id("ID").Op(":").Lit(1), jen.Id("ClientID").Op(":").Lit("example"), jen.Id("ClientSecret").Op(":").Lit("blah")))), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/api/v1/oauth2/clients"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected")))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("GetOAuth2Clients").Call(jen.Id("ctx"), jen.Id("nil")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildCreateOAuth2ClientRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(jen.Id("Username").Op(":").Lit("username"), jen.Id("Password").Op(":").Lit("password"), jen.Id("TOTPToken").Op(":").Lit("123456"))), jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildCreateOAuth2ClientRequest").Call(jen.Id("ctx"), jen.Op("&").Qual("net/http", "Cookie").Values(), jen.Id("exampleInput")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Qual("net/http", "MethodPost"), jen.Id("req").Dot("Method")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_CreateOAuth2Client").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(jen.Id("Username").Op(":").Lit("username"), jen.Id("Password").Op(":").Lit("password"), jen.Id("TOTPToken").Op(":").Lit("123456"))), jen.Id("exampleOutput").Op(":=").Op("&").Id("models").Dot("OAuth2Client").Values(jen.Id("ClientID").Op(":").Lit("EXAMPLECLIENTID"), jen.Id("ClientSecret").Op(":").Lit("EXAMPLECLIENTSECRET")), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Lit("/oauth2/client"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("exampleOutput")))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("oac"), jen.Id("err")).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(jen.Id("ctx"), jen.Op("&").Qual("net/http", "Cookie").Values(), jen.Id("exampleInput")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err")), jen.Id("assert").Dot("NotNil").Call(jen.Id("t"), jen.Id("oac")))), jen.Id("T").Dot("Run").Call(jen.Lit("with invalid body"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(jen.Id("Username").Op(":").Lit("username"), jen.Id("Password").Op(":").Lit("password"), jen.Id("TOTPToken").Op(":").Lit("123456"))), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/oauth2/client"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.List(jen.Id("_"), jen.Id("err")).Op(":=").Id("res").Dot("Write").Call(jen.Index().Id("byte").Call(jen.Lit(`


					BLAH

				`))), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("oac"), jen.Id("err")).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(jen.Id("ctx"), jen.Op("&").Qual("net/http", "Cookie").Values(), jen.Id("exampleInput")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")), jen.Id("assert").Dot("Nil").Call(jen.Id("t"), jen.Id("oac")))), jen.Id("T").Dot("Run").Call(jen.Lit("with timeout"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(jen.Id("Username").Op(":").Lit("username"), jen.Id("Password").Op(":").Lit("password"), jen.Id("TOTPToken").Op(":").Lit("123456"))), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/oauth2/client"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Qual("time", "Sleep").Call(jen.Lit(10).Op("*").Qual("time", "Hour"))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.Id("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"), jen.List(jen.Id("oac"), jen.Id("err")).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(jen.Id("ctx"), jen.Op("&").Qual("net/http", "Cookie").Values(), jen.Id("exampleInput")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")), jen.Id("assert").Dot("Nil").Call(jen.Id("t"), jen.Id("oac")))), jen.Id("T").Dot("Run").Call(jen.Lit("with 404"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(jen.Id("Username").Op(":").Lit("username"), jen.Id("Password").Op(":").Lit("password"), jen.Id("TOTPToken").Op(":").Lit("123456"))), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/oauth2/client"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound"))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("oac"), jen.Id("err")).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(jen.Id("ctx"), jen.Op("&").Qual("net/http", "Cookie").Values(), jen.Id("exampleInput")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("err"), jen.Id("ErrNotFound")), jen.Id("assert").Dot("Nil").Call(jen.Id("t"), jen.Id("oac")))), jen.Id("T").Dot("Run").Call(jen.Lit("with no cookie"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("_"), jen.Id("err")).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(jen.Id("ctx"), jen.Id("nil"), jen.Id("nil")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildArchiveOAuth2ClientRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodDelete"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("expectedID").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildArchiveOAuth2ClientRequest").Call(jen.Id("ctx"), jen.Id("expectedID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual").Dot("URL")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("actual").Dot("URL").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.Id("expectedID")))), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_ArchiveOAuth2Client").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/api/v1/oauth2/clients/%d"), jen.Id("expected")), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodDelete")), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("err").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")).Dot("ArchiveOAuth2Client").Call(jen.Id("ctx"), jen.Id("expected")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")))),
		),
		jen.Line(),
	)

	return ret
}

func roundtripperDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())
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
		jen.Func().Params(jen.Id("t").Op("*").Id("defaultRoundTripper")).Id("RoundTrip").Params(jen.Id("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual("net/http", "Response"), jen.Id("error")).Block(
			jen.Id("req").Dot("Header").Dot("Set").Call(
				jen.Id("userAgentHeader"),
				jen.Id("userAgent"),
			),
			jen.Line(),
			jen.Return().Id("t").Dot("baseTransport").Dot("RoundTrip").Call(jen.Id("req")),
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
					)).Dot("DialContext"),
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

func usersDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())
	ret.Add(jen.Var().Id("usersBasePath").Op("=").Lit("users"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildGetUserRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("userID").Id("uint64")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(jen.Id("nil"), jen.Id("usersBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("userID"), jen.Lit(10))), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Id("uri"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("GetUser").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("userID").Id("uint64")).Params(
			jen.Id("user").Op("*").Id("models").Dot("User"), jen.Id("err").Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetUserRequest").Call(jen.Id("ctx"), jen.Id("userID")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("user")), jen.Return().List(jen.Id("user"), jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildGetUsersRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("filter").Op("*").Id("models").Dot("QueryFilter")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(jen.Id("filter").Dot("ToValues").Call(), jen.Id("usersBasePath")), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Id("uri"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("GetUsers").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("filter").Op("*").Id("models").Dot("QueryFilter")).Params(
			jen.Op("*").Id("models").Dot("UserList"), jen.Id("error")).Block(
			jen.Id("users").Op(":=").Op("&").Id("models").Dot("UserList").Values(), jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetUsersRequest").Call(jen.Id("ctx"), jen.Id("filter")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("users")), jen.Return().List(jen.Id("users"), jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildCreateUserRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("body").Op("*").Id("models").Dot("UserInput")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(jen.Id("nil"), jen.Id("usersBasePath")), jen.Return().Id("c").Dot("buildDataRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Id("uri"), jen.Id("body")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("CreateUser").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("input").Op("*").Id("models").Dot("UserInput")).Params(
			jen.Op("*").Id("models").Dot("UserCreationResponse"), jen.Id("error")).Block(
			jen.Id("user").Op(":=").Op("&").Id("models").Dot("UserCreationResponse").Values(), jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildCreateUserRequest").Call(jen.Id("ctx"), jen.Id("input")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.Id("err").Op("=").Id("c").Dot("executeUnathenticatedDataRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("user")), jen.Return().List(jen.Id("user"), jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildArchiveUserRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("userID").Id("uint64")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(jen.Id("nil"), jen.Id("usersBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("userID"), jen.Lit(10))), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodDelete"), jen.Id("uri"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("ArchiveUser").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("userID").Id("uint64")).Params(jen.Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildArchiveUserRequest").Call(jen.Id("ctx"), jen.Id("userID")),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("building request"), jen.Id("err"))), jen.Return().Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildLoginRequest").Params(
			jen.List(jen.Id("username"), jen.Id("password"), jen.Id("totpToken")).Id("string")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.List(jen.Id("body"), jen.Id("err")).Op(":=").Id("createBodyFromStruct").Call(jen.Op("&").Id("models").Dot("UserLoginInput").Values(jen.Id("Username").Op(":").Id("username"), jen.Id("Password").Op(":").Id("password"), jen.Id("TOTPToken").Op(":").Id("totpToken"))), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("creating body from struct"), jen.Id("err")))), jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(jen.Id("nil"), jen.Id("usersBasePath"), jen.Lit("login")), jen.Return().Id("c").Dot("buildDataRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Id("uri"), jen.Id("body")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("Login").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.List(jen.Id("username"), jen.Id("password"), jen.Id("totpToken")).Id("string")).Params(
			jen.Op("*").Qual("net/http", "Cookie"), jen.Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildLoginRequest").Call(jen.Id("username"), jen.Id("password"), jen.Id("totpToken")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("err"))), jen.List(jen.Id("res"), jen.Id("err")).Op(":=").Id("c").Dot("plainClient").Dot("Do").Call(jen.Id("req")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("encountered error executing login request: %w"), jen.Id("err")))), jen.If(jen.Id("c").Dot("Debug")).Block(
				jen.List(jen.Id("b"), jen.Id("err")).Op(":=").Id("httputil").Dot("DumpResponse").Call(jen.Id("res"), jen.Id("true")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
					jen.Id("c").Dot("logger").Dot("Error").Call(jen.Id("err"), jen.Lit("dumping response"))), jen.Id("c").Dot("logger").Dot("WithValue").Call(jen.Lit("response"), jen.Id("string").Call(jen.Id("b"))).Dot("Debug").Call(jen.Lit("login response received"))), jen.Defer().Func().Params().Block(
				jen.If(jen.Id("err").Op(":=").Id("res").Dot("Body").Dot("Close").Call(), jen.Id("err").Op("!=").Id("nil")).Block(
					jen.Id("c").Dot("logger").Dot("Error").Call(jen.Id("err"), jen.Lit("closing response body")))).Call(), jen.Id("cookies").Op(":=").Id("res").Dot("Cookies").Call(), jen.If(jen.Id("len").Call(jen.Id("cookies")).Op(">").Lit(0)).Block(
				jen.Return().List(jen.Id("cookies").Index(jen.Lit(0)), jen.Id("nil"))), jen.Return().List(jen.Id("nil"), jen.Id("errors").Dot("New").Call(jen.Lit("no cookies returned from request"))),
		),
		jen.Line(),
	)

	return ret
}

func usersTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())

	ret.Add(
		jen.Func().Id("TestV1Client_BuildGetUserRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodGet"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.Id("expectedID").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetUserRequest").Call(jen.Id("ctx"), jen.Id("expectedID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("actual").Dot("URL").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.Id("expectedID")))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_GetUser").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("User").Values(jen.Id("ID").Op(":").Lit(1)), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("req").Dot("URL").Dot("String").Call(), jen.Qual("strconv", "Itoa").Call(jen.Id("int").Call(jen.Id("expected").Dot("ID"))))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/users/%d"), jen.Id("expected").Dot("ID")), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected")))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("GetUser").Call(jen.Id("ctx"), jen.Id("expected").Dot("ID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildGetUsersRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodGet"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetUsersRequest").Call(jen.Id("ctx"), jen.Id("nil")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_GetUsers").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("UserList").Values(jen.Id("Users").Op(":").Index().Id("models").Dot("User").Values(jen.Values(jen.Id("ID").Op(":").Lit(1)))), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/users"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected")))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("GetUsers").Call(jen.Id("ctx"), jen.Id("nil")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildCreateUserRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("UserInput").Values(), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildCreateUserRequest").Call(jen.Id("ctx"), jen.Id("exampleInput")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_CreateUser").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("UserCreationResponse").Values(jen.Id("ID").Op(":").Lit(1)), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("UserInput").Values(), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/users"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Var().Id("x").Op("*").Id("models").Dot("UserInput"), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.Id("req").Dot("Body")).Dot("Decode").Call(jen.Op("&").Id("x"))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("exampleInput"), jen.Id("x")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected"))), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("CreateUser").Call(jen.Id("ctx"), jen.Id("exampleInput")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildArchiveUserRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodDelete"), jen.Id("expectedID").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildArchiveUserRequest").Call(jen.Id("ctx"), jen.Id("expectedID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual").Dot("URL")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("actual").Dot("URL").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.Id("expectedID")))), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_ArchiveUser").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/users/%d"), jen.Id("expected")), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodDelete")), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("err").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")).Dot("ArchiveUser").Call(jen.Id("ctx"), jen.Id("expected")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildLoginRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildLoginRequest").Call(jen.Lit("username"), jen.Lit("password"), jen.Lit("123456")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("req")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_Login").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/users/login"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Qual("net/http", "SetCookie").Call(jen.Id("res"), jen.Op("&").Qual("net/http", "Cookie").Values(jen.Id("Name").Op(":").Lit("hi"))), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("cookie"), jen.Id("err")).Op(":=").Id("c").Dot("Login").Call(jen.Id("ctx"), jen.Lit("username"), jen.Lit("password"), jen.Lit("123456")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("cookie")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err")))), jen.Id("T").Dot("Run").Call(jen.Lit("with timeout"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/users/login"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Qual("time", "Sleep").Call(jen.Lit(10).Op("*").Qual("time", "Hour")), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.Id("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Microsecond"), jen.List(jen.Id("cookie"), jen.Id("err")).Op(":=").Id("c").Dot("Login").Call(jen.Id("ctx"), jen.Lit("username"), jen.Lit("password"), jen.Lit("123456")), jen.Id("require").Dot("Nil").Call(jen.Id("t"), jen.Id("cookie")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")))), jen.Id("T").Dot("Run").Call(jen.Lit("with missing cookie"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/users/login"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("cookie"), jen.Id("err")).Op(":=").Id("c").Dot("Login").Call(jen.Id("ctx"), jen.Lit("username"), jen.Lit("password"), jen.Lit("123456")), jen.Id("require").Dot("Nil").Call(jen.Id("t"), jen.Id("cookie")), jen.Id("assert").Dot("Error").Call(jen.Id("t"), jen.Id("err")))),
		),
		jen.Line(),
	)

	return ret
}

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())
	ret.Add(jen.Var().Id("webhooksBasePath").Op("=").Lit("webhooks"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildGetWebhookRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("webhooksBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("id"), jen.Lit(10))), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Id("uri"), jen.Id("nil")),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("GetWebhook").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Id("webhook").Op("*").Id("models").Dot("Webhook"), jen.Id("err").Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetWebhookRequest").Call(jen.Id("ctx"), jen.Id("id")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("webhook")), jen.Return().List(jen.Id("webhook"), jen.Id("err")),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildGetWebhooksRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("filter").Op("*").Id("models").Dot("QueryFilter")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("filter").Dot("ToValues").Call(), jen.Id("webhooksBasePath")), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Id("uri"), jen.Id("nil")),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("GetWebhooks").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("filter").Op("*").Id("models").Dot("QueryFilter")).Params(
			jen.Id("webhooks").Op("*").Id("models").Dot("WebhookList"), jen.Id("err").Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetWebhooksRequest").Call(jen.Id("ctx"), jen.Id("filter")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("webhooks")), jen.Return().List(jen.Id("webhooks"), jen.Id("err")),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildCreateWebhookRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("body").Op("*").Id("models").Dot("WebhookCreationInput")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("webhooksBasePath")), jen.Return().Id("c").Dot("buildDataRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Id("uri"), jen.Id("body")),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("CreateWebhook").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("input").Op("*").Id("models").Dot("WebhookCreationInput")).Params(
			jen.Id("webhook").Op("*").Id("models").Dot("Webhook"), jen.Id("err").Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildCreateWebhookRequest").Call(jen.Id("ctx"), jen.Id("input")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.Id("err").Op("=").Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("webhook")), jen.Return().List(jen.Id("webhook"), jen.Id("err")),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildUpdateWebhookRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("updated").Op("*").Id("models").Dot("Webhook")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("webhooksBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("updated").Dot("ID"), jen.Lit(10))), jen.Return().Id("c").Dot("buildDataRequest").Call(jen.Qual("net/http", "MethodPut"), jen.Id("uri"), jen.Id("updated")),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("UpdateWebhook").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("updated").Op("*").Id("models").Dot("Webhook")).Params(
			jen.Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildUpdateWebhookRequest").Call(jen.Id("ctx"), jen.Id("updated")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err"))), jen.Return().Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("updated")),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildArchiveWebhookRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("webhooksBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("id"), jen.Lit(10))), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodDelete"), jen.Id("uri"), jen.Id("nil")),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("ArchiveWebhook").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildArchiveWebhookRequest").Call(jen.Id("ctx"), jen.Id("id")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err"))), jen.Return().Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Id("nil")),
		),
	)
	return ret
}

func webhooksTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())
	ret.ImportNames(map[string]string{
		"github.com/stretchr/testify/assert":  "assert",
		"github.com/stretchr/testify/mock":    "mock",
		"github.com/stretchr/testify/require": "require",
	})

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_BuildGetWebhookRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodGet"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.Id("expectedID").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetWebhookRequest").Call(jen.Id("ctx"), jen.Id("expectedID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("actual").Dot("URL").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.Id("expectedID")))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
	)

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_GetWebhook").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Webhook").Values(jen.Id("ID").Op(":").Lit(1), jen.Id("Name").Op(":").Lit("example")), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("req").Dot("URL").Dot("String").Call(), jen.Qual("strconv", "Itoa").Call(jen.Id("int").Call(jen.Id("expected").Dot("ID"))))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/api/v1/webhooks/%d"), jen.Id("expected").Dot("ID")), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected")))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("GetWebhook").Call(jen.Id("ctx"), jen.Id("expected").Dot("ID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
	)

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_BuildGetWebhooksRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodGet"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetWebhooksRequest").Call(jen.Id("ctx"), jen.Id("nil")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
	)

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_GetWebhooks").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("WebhookList").Values(jen.Id("Webhooks").Op(":").Index().Id("models").Dot("Webhook").Values(jen.Values(jen.Id("ID").Op(":").Lit(1), jen.Id("Name").Op(":").Lit("example")))), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/api/v1/webhooks"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected")))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("GetWebhooks").Call(jen.Id("ctx"), jen.Id("nil")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
	)

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_BuildCreateWebhookRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("WebhookCreationInput").Values(jen.Id("Name").Op(":").Lit("expected name")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildCreateWebhookRequest").Call(jen.Id("ctx"), jen.Id("exampleInput")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
	)

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_CreateWebhook").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Webhook").Values(jen.Id("ID").Op(":").Lit(1), jen.Id("Name").Op(":").Lit("example")), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("WebhookCreationInput").Values(jen.Id("Name").Op(":").Id("expected").Dot("Name")), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/api/v1/webhooks"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Var().Id("x").Op("*").Id("models").Dot("WebhookCreationInput"), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.Id("req").Dot("Body")).Dot("Decode").Call(jen.Op("&").Id("x"))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("exampleInput"), jen.Id("x")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected"))), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("CreateWebhook").Call(jen.Id("ctx"), jen.Id("exampleInput")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
	)

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_BuildUpdateWebhookRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPut"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("Webhook").Values(jen.Id("Name").Op(":").Lit("changed name")), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildUpdateWebhookRequest").Call(jen.Id("ctx"), jen.Id("exampleInput")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
	)

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_UpdateWebhook").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Webhook").Values(jen.Id("ID").Op(":").Lit(1), jen.Id("Name").Op(":").Lit("example")), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/api/v1/webhooks/%d"), jen.Id("expected").Dot("ID")), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPut")), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("err").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")).Dot("UpdateWebhook").Call(jen.Id("ctx"), jen.Id("expected")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")))),
		),
	)

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_BuildArchiveWebhookRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodDelete"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("expectedID").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildArchiveWebhookRequest").Call(jen.Id("ctx"), jen.Id("expectedID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual").Dot("URL")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("actual").Dot("URL").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.Id("expectedID")))), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
	)

	ret.Add(
		jen.Line(),
		jen.Func().Id("TestV1Client_ArchiveWebhook").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/api/v1/webhooks/%d"), jen.Id("expected")), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodDelete")), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("err").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")).Dot("ArchiveWebhook").Call(jen.Id("ctx"), jen.Id("expected")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")))),
		),
	)
	return ret
}

func itemsDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())
	ret.Add(jen.Var().Id("itemsBasePath").Op("=").Lit("items"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildGetItemRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("itemsBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("id"), jen.Lit(10))), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Id("uri"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("GetItem").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Id("item").Op("*").Id("models").Dot("Item"), jen.Id("err").Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetItemRequest").Call(jen.Id("ctx"), jen.Id("id")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.If(jen.Id("retrieveErr").Op(":=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("item")), jen.Id("retrieveErr").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("retrieveErr"))), jen.Return().List(jen.Id("item"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildGetItemsRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("filter").Op("*").Id("models").Dot("QueryFilter")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("filter").Dot("ToValues").Call(), jen.Id("itemsBasePath")), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Id("uri"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("GetItems").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("filter").Op("*").Id("models").Dot("QueryFilter")).Params(
			jen.Id("items").Op("*").Id("models").Dot("ItemList"), jen.Id("err").Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetItemsRequest").Call(jen.Id("ctx"), jen.Id("filter")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.If(jen.Id("retrieveErr").Op(":=").Id("c").Dot("retrieve").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("items")), jen.Id("retrieveErr").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Id("retrieveErr"))), jen.Return().List(jen.Id("items"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildCreateItemRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("body").Op("*").Id("models").Dot("ItemCreationInput")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("itemsBasePath")), jen.Return().Id("c").Dot("buildDataRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Id("uri"), jen.Id("body")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("CreateItem").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("input").Op("*").Id("models").Dot("ItemCreationInput")).Params(
			jen.Id("item").Op("*").Id("models").Dot("Item"), jen.Id("err").Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildCreateItemRequest").Call(jen.Id("ctx"), jen.Id("input")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err")))), jen.Id("err").Op("=").Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("item")), jen.Return().List(jen.Id("item"), jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildUpdateItemRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("updated").Op("*").Id("models").Dot("Item")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("itemsBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("updated").Dot("ID"), jen.Lit(10))), jen.Return().Id("c").Dot("buildDataRequest").Call(jen.Qual("net/http", "MethodPut"), jen.Id("uri"), jen.Id("updated")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("UpdateItem").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("updated").Op("*").Id("models").Dot("Item")).Params(
			jen.Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildUpdateItemRequest").Call(jen.Id("ctx"), jen.Id("updated")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err"))), jen.Return().Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("updated")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("BuildArchiveItemRequest").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Op("*").Qual("net/http", "Request"), jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(jen.Id("nil"), jen.Id("itemsBasePath"), jen.Qual("strconv", "FormatUint").Call(jen.Id("id"), jen.Lit(10))), jen.Return().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodDelete"), jen.Id("uri"), jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id("V1Client")).Id("ArchiveItem").Params(
			jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Id("uint64")).Params(
			jen.Id("error")).Block(
			jen.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("c").Dot("BuildArchiveItemRequest").Call(jen.Id("ctx"), jen.Id("id")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Id("err"))), jen.Return().Id("c").Dot("executeRequest").Call(jen.Id("ctx"), jen.Id("req"), jen.Id("nil")),
		),
	)
	return ret
}

func itemsTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.Add(jen.Null())
	ret.ImportNames(map[string]string{
		"github.com/stretchr/testify/assert":  "assert",
		"github.com/stretchr/testify/mock":    "mock",
		"github.com/stretchr/testify/require": "require",
	})

	ret.Add(
		jen.Func().Id("TestV1Client_BuildGetItemRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodGet"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.Id("expectedID").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetItemRequest").Call(jen.Id("ctx"), jen.Id("expectedID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("actual").Dot("URL").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.Id("expectedID")))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_GetItem").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Item").Values(jen.Id("ID").Op(":").Lit(1), jen.Id("Name").Op(":").Lit("example"), jen.Id("Details").Op(":").Lit("blah")), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("req").Dot("URL").Dot("String").Call(), jen.Qual("strconv", "Itoa").Call(jen.Id("int").Call(jen.Id("expected").Dot("ID"))))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/api/v1/items/%d"), jen.Id("expected").Dot("ID")), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected")))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("GetItem").Call(jen.Id("ctx"), jen.Id("expected").Dot("ID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildGetItemsRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodGet"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetItemsRequest").Call(jen.Id("ctx"), jen.Id("nil")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_GetItems").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("ItemList").Values(jen.Id("Items").Op(":").Index().Id("models").Dot("Item").Values(jen.Values(jen.Id("ID").Op(":").Lit(1), jen.Id("Name").Op(":").Lit("example"), jen.Id("Details").Op(":").Lit("blah")))), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/api/v1/items"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected")))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("GetItems").Call(jen.Id("ctx"), jen.Id("nil")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildCreateItemRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPost"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("ItemCreationInput").Values(jen.Id("Name").Op(":").Lit("expected name"), jen.Id("Details").Op(":").Lit("expected details")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildCreateItemRequest").Call(jen.Id("ctx"), jen.Id("exampleInput")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_CreateItem").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Item").Values(jen.Id("ID").Op(":").Lit(1), jen.Id("Name").Op(":").Lit("example"), jen.Id("Details").Op(":").Lit("blah")), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("ItemCreationInput").Values(jen.Id("Name").Op(":").Id("expected").Dot("Name"), jen.Id("Details").Op(":").Id("expected").Dot("Details")), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/api/v1/items"), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost")), jen.Var().Id("x").Op("*").Id("models").Dot("ItemCreationInput"), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.Id("req").Dot("Body")).Dot("Decode").Call(jen.Op("&").Id("x"))), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("exampleInput"), jen.Id("x")), jen.Id("require").Dot("NoError").Call(jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected"))), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("CreateItem").Call(jen.Id("ctx"), jen.Id("exampleInput")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("expected"), jen.Id("actual")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildUpdateItemRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodPut"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("Item").Values(jen.Id("Name").Op(":").Lit("changed name"), jen.Id("Details").Op(":").Lit("changed details")), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildUpdateItemRequest").Call(jen.Id("ctx"), jen.Id("exampleInput")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_UpdateItem").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Item").Values(jen.Id("ID").Op(":").Lit(1), jen.Id("Name").Op(":").Lit("example"), jen.Id("Details").Op(":").Lit("blah")), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/api/v1/items/%d"), jen.Id("expected").Dot("ID")), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPut")), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("err").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")).Dot("UpdateItem").Call(jen.Id("ctx"), jen.Id("expected")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_BuildArchiveItemRequest").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expectedMethod").Op(":=").Qual("net/http", "MethodDelete"), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Id("nil")), jen.Id("expectedID").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.Id("c").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")), jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildArchiveItemRequest").Call(jen.Id("ctx"), jen.Id("expectedID")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual")), jen.Id("require").Dot("NotNil").Call(jen.Id("t"), jen.Id("actual").Dot("URL")), jen.Id("assert").Dot("True").Call(jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(jen.Id("actual").Dot("URL").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.Id("expectedID")))), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("actual").Dot("Method"), jen.Id("expectedMethod"), jen.Lit("request should be a %s request"), jen.Id("expectedMethod")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_ArchiveItem").Params(
			jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(), jen.Id("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(
				jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Id("expected").Op(":=").Id("uint64").Call(jen.Lit(1)), jen.Id("ctx").Op(":=").Qual("context", "Background").Call(), jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
					jen.Id("res").Qual("net/http", "ResponseWriter"), jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
					jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(jen.Lit("/api/v1/items/%d"), jen.Id("expected")), jen.Lit("expected and actual path don't match")), jen.Id("assert").Dot("Equal").Call(jen.Id("t"), jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodDelete")), jen.Id("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))), jen.Id("err").Op(":=").Id("buildTestClient").Call(jen.Id("t"), jen.Id("ts")).Dot("ArchiveItem").Call(jen.Id("ctx"), jen.Id("expected")), jen.Id("assert").Dot("NoError").Call(jen.Id("t"), jen.Id("err"), jen.Lit("no error should be returned")))),
		),
		jen.Line(),
	)

	return ret
}
