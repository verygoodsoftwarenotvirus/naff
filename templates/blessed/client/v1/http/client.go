package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func newClientMethod(name string) *jen.Statement {
	return jen.Func().Params(jen.ID("c").PointerTo().ID(v1)).ID(name)
}

func mainDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)
	ret.Add(jen.Line())

	// consts
	ret.Add(
		jen.Const().Defs(
			jen.ID("defaultTimeout").Equals().Lit(30).Times().Qual("time", "Second"),
			jen.ID("clientName").Equals().Lit("v1_client"),
		),
	)

	// vars
	ret.Add(
		jen.Var().Defs(
			jen.Comment("ErrNotFound is a handy error to return when we receive a 404 response."),
			jen.ID("ErrNotFound").Equals().Qual("fmt", "Errorf").Call(jen.Lit("%d: not found"), jen.Qual("net/http", "StatusNotFound")),
			jen.Line(),
			jen.Comment("ErrUnauthorized is a handy error to return when we receive a 401 response."),
			jen.ID("ErrUnauthorized").Equals().Qual("fmt", "Errorf").Call(jen.Lit("%d: not authorized"), jen.Qual("net/http", "StatusUnauthorized")),
			jen.Line(),
			jen.Comment("ErrInvalidTOTPToken is an error for when our TOTP validation request goes awry."),
			jen.ID("ErrInvalidTOTPToken").Equals().Qual("errors", "New").Call(jen.Lit("invalid TOTP token")),
		),
		jen.Line(),
	)

	// types
	ret.Add(
		jen.Commentf("%s is a client for interacting with v1 of our HTTP API.", v1),
		jen.Line(),
		jen.Type().ID(v1).Struct(
			jen.ID("plainClient").PointerTo().Qual("net/http", "Client"),
			jen.ID("authedClient").PointerTo().Qual("net/http", "Client"),
			jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
			jen.ID("Debug").Bool(),
			jen.ID("URL").PointerTo().Qual("net/url", "URL"),
			jen.ID("Scopes").Index().String(),
			jen.ID("tokenSource").Qual("golang.org/x/oauth2", "TokenSource"),
		),
		jen.Line(),
	)

	ret.Add(buildAuthenticatedClient()...)
	ret.Add(buildPlainClient()...)
	ret.Add(buildTokenSource()...)
	ret.Add(buildTokenEndpoint()...)
	ret.Add(buildNewClient()...)
	ret.Add(buildNewSimpleClient()...)
	ret.Add(buildBuildOAuthClient()...)
	ret.Add(buildCloseResponseBody()...)
	ret.Add(buildExportedBuildURL()...)
	ret.Add(buildUnexportedBuildURL()...)
	ret.Add(buildBuildVersionlessURL()...)
	ret.Add(buildBuildWebsocketURL()...)
	ret.Add(buildBuildHealthCheckRequest()...)
	ret.Add(buildIsUp()...)
	ret.Add(buildBuildDataRequest(proj)...)
	ret.Add(buildExecuteRequest(proj)...)
	ret.Add(buildExecuteRawRequest(proj)...)
	ret.Add(buildCheckExistence(proj)...)
	ret.Add(buildRetrieve(proj)...)
	ret.Add(buildExecuteUnauthenticatedDataRequest(proj)...)

	return ret
}

func buildAuthenticatedClient() []jen.Code {
	lines := []jen.Code{
		jen.Comment("AuthenticatedClient returns the authenticated *http.Client that we use to make most requests."),
		jen.Line(),
		newClientMethod("AuthenticatedClient").Params().Params(jen.PointerTo().Qual("net/http", "Client")).Block(
			jen.Return().ID("c").Dot("authedClient"),
		),
		jen.Line(),
	}

	return lines
}

func buildPlainClient() []jen.Code {
	lines := []jen.Code{
		jen.Comment("PlainClient returns the unauthenticated *http.Client that we use to make certain requests."),
		jen.Line(),
		newClientMethod("PlainClient").Params().Params(jen.PointerTo().Qual("net/http", "Client")).Block(
			jen.Return().ID("c").Dot("plainClient"),
		),
		jen.Line(),
	}

	return lines
}

func buildTokenSource() []jen.Code {
	lines := []jen.Code{
		jen.Comment("TokenSource provides the client's token source."),
		jen.Line(),
		newClientMethod("TokenSource").Params().Params(jen.ID("oauth2").Dot("TokenSource")).Block(
			jen.Return().ID("c").Dot("tokenSource"),
		),
		jen.Line(),
	}

	return lines
}

func buildNewClient() []jen.Code {
	lines := []jen.Code{
		jen.Comment("NewClient builds a new API client for us."),
		jen.Line(),
		jen.Func().ID("NewClient").Paramsln(
			constants.CtxParam(),
			jen.Listln(
				jen.ID("clientID"),
				jen.ID("clientSecret"),
			).String(),
			jen.ID("address").PointerTo().Qual("net/url", "URL"),
			jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
			jen.ID("hclient").PointerTo().Qual("net/http", "Client"),
			jen.ID("scopes").Index().String(),
			jen.ID("debug").Bool(),
		).Params(
			jen.PointerTo().ID(v1),
			jen.Error(),
		).Block(
			jen.Var().ID("client").Equals().ID("hclient"),
			jen.If(jen.ID("client").IsEqualTo().ID("nil")).Block(
				jen.ID("client").Equals().AddressOf().Qual("net/http", "Client").Valuesln(
					jen.ID("Timeout").MapAssign().ID("defaultTimeout"),
				),
			),
			jen.If(jen.ID("client").Dot("Timeout").IsEqualTo().Zero()).Block(
				jen.ID("client").Dot("Timeout").Equals().ID("defaultTimeout"),
			),
			jen.Line(),
			jen.If(jen.ID("debug")).Block(
				jen.ID(constants.LoggerVarName).Dot("SetLevel").Call(
					jen.Qual(utils.LoggingPkg, "DebugLevel"),
				),
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(
					jen.Lit("log level set to debug!"),
				),
			),
			jen.Line(),
			jen.List(
				jen.ID("ac"),
				jen.ID("ts"),
			).Assign().ID("buildOAuthClient").Call(
				constants.CtxVar(),
				jen.ID("address"),
				jen.ID("clientID"),
				jen.ID("clientSecret"),
				jen.ID("scopes"),
				jen.ID("client").Dot("Timeout"),
			),
			jen.Line(),
			jen.ID("c").Assign().AddressOf().ID(v1).Valuesln(
				jen.ID("URL").MapAssign().ID("address"),
				jen.ID("plainClient").MapAssign().ID("client"),
				jen.ID(constants.LoggerVarName).MapAssign().Qual(constants.LoggerVarName, "WithName").Call(jen.ID("clientName")),
				jen.ID("Debug").MapAssign().ID("debug"),
				jen.ID("authedClient").MapAssign().ID("ac"),
				jen.ID("tokenSource").MapAssign().ID("ts"),
			),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Dot("WithValue").Call(
				jen.Lit("url"),
				jen.ID("address").Dot("String").Call(),
			).Dot("Debug").Call(
				jen.Lit("returning client"),
			),
			jen.Return().List(jen.ID("c"),
				jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildOAuthClient() []jen.Code {
	lines := []jen.Code{
		jen.Comment("buildOAuthClient takes care of all the OAuth2 noise and returns a nice pretty *http.Client for us to use."),
		jen.Line(),
		jen.Func().ID("buildOAuthClient").Paramsln(
			constants.CtxParam(),
			jen.ID("uri").PointerTo().Qual("net/url", "URL"),
			jen.Listln(
				jen.ID("clientID"),
				jen.ID("clientSecret"),
			).String(),
			jen.ID("scopes").Index().String(),
			jen.ID("timeout").Qual("time", "Duration"),
		).Params(
			jen.PointerTo().Qual("net/http", "Client"),
			jen.ID("oauth2").Dot("TokenSource"),
		).Block(
			jen.ID("conf").Assign().Qual("golang.org/x/oauth2/clientcredentials", "Config").Valuesln(
				jen.ID("ClientID").MapAssign().ID("clientID"),
				jen.ID("ClientSecret").MapAssign().ID("clientSecret"),
				jen.ID("Scopes").MapAssign().ID("scopes"),
				jen.ID("EndpointParams").MapAssign().Qual("net/url", "Values").Valuesln(
					jen.Lit("client_id").MapAssign().Index().String().Values(jen.ID("clientID")),
					jen.Lit("client_secret").MapAssign().Index().String().Values(jen.ID("clientSecret")),
				),
				jen.ID("TokenURL").MapAssign().ID("tokenEndpoint").Call(jen.ID("uri")).Dot("TokenURL"),
			),
			jen.Line(),
			jen.ID("ts").Assign().ID("oauth2").Dot("ReuseTokenSource").Call(
				jen.Nil(),
				jen.ID("conf").Dot("TokenSource").Call(
					constants.CtxVar(),
				),
			),
			jen.ID("client").Assign().AddressOf().Qual("net/http", "Client").Valuesln(
				jen.ID("Transport").MapAssign().AddressOf().ID("oauth2").Dot("Transport").Valuesln(
					jen.ID("Base").MapAssign().AddressOf().Qual("go.opencensus.io/plugin/ochttp", "Transport").Valuesln(
						jen.ID("Base").MapAssign().ID("newDefaultRoundTripper").Call(),
					),
					jen.ID("Source").MapAssign().ID("ts"),
				),
				jen.ID("Timeout").MapAssign().ID("timeout"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("client"),
				jen.ID("ts")),
		),
		jen.Line(),
	}

	return lines
}

func buildTokenEndpoint() []jen.Code {
	lines := []jen.Code{
		jen.Comment("tokenEndpoint provides the oauth2 Endpoint for a given host."),
		jen.Line(),
		jen.Func().ID("tokenEndpoint").Params(
			jen.ID("baseURL").PointerTo().Qual("net/url", "URL"),
		).Params(
			jen.ID("oauth2").Dot("Endpoint"),
		).Block(
			jen.List(
				jen.ID("tu"),
				jen.ID("au"),
			).Assign().List(jen.PointerTo().ID("baseURL"), jen.PointerTo().ID("baseURL")),
			jen.List(
				jen.ID("tu").Dot("Path"),
				jen.ID("au").Dot("Path"),
			).Equals().List(
				jen.Lit("oauth2/token"),
				jen.Lit("oauth2/authorize"),
			),
			jen.Line(),
			jen.Return().ID("oauth2").Dot("Endpoint").Valuesln(
				jen.ID("TokenURL").MapAssign().ID("tu").Dot("String").Call(),
				jen.ID("AuthURL").MapAssign().ID("au").Dot("String").Call(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildNewSimpleClient() []jen.Code {
	lines := []jen.Code{
		jen.Comment("NewSimpleClient is a client that is capable of much less than the normal client"),
		jen.Line(),
		jen.Comment("and has noops or empty values for most of its authentication and debug parts."),
		jen.Line(),
		jen.Comment("Its purpose at the time of this writing is merely so I can make users (which"),
		jen.Line(),
		jen.Comment("is a route that doesn't require authentication.)"),
		jen.Line(),
		jen.Func().ID("NewSimpleClient").Params(
			constants.CtxParam(),
			jen.ID("address").PointerTo().Qual("net/url", "URL"),
			jen.ID("debug").Bool(),
		).Params(
			jen.PointerTo().ID(v1),
			jen.Error(),
		).Block(
			jen.Return().ID("NewClient").Callln(
				constants.CtxVar(),
				jen.EmptyString(),
				jen.EmptyString(),
				jen.ID("address"),
				jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.AddressOf().Qual("net/http", "Client").Values(
					jen.ID("Timeout").MapAssign().Lit(5).Times().Qual("time", "Second"),
				),
				jen.Index().String().Values(jen.Lit("*")),
				jen.ID("debug"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildCloseResponseBody() []jen.Code {
	lines := []jen.Code{
		jen.Comment("closeResponseBody takes a given HTTP response and closes its body, logging if an error occurs."),
		jen.Line(),
		newClientMethod("closeResponseBody").Params(jen.ID(constants.ResponseVarName).PointerTo().Qual("net/http", "Response")).Block(
			jen.If(jen.ID(constants.ResponseVarName).DoesNotEqual().Nil()).Block(
				jen.If(jen.Err().Assign().ID(constants.ResponseVarName).Dot("Body").Dot("Close").Call(), jen.Err().DoesNotEqual().Nil()).Block(
					jen.ID("c").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("closing response body")),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildExecuteRawRequest(proj *models.Project) []jen.Code {
	block := []jen.Code{
		utils.StartSpan(proj, true, "executeRawRequest"),
		jen.Var().ID(constants.LoggerVarName).Equals().ID("c").Dot(constants.LoggerVarName),
		jen.If(jen.List(
			jen.ID("command"),
			jen.Err(),
		).Assign().Qual("github.com/moul/http2curl", "GetCurlCommand").Call(
			jen.ID(constants.RequestVarName),
		),
			jen.Err().IsEqualTo().ID("nil").And().ID("c").Dot("Debug"),
		).Block(
			jen.ID(constants.LoggerVarName).Equals().ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(
				jen.Lit("curl"),
				jen.ID("command").Dot("String").Call(),
			),
		),
		jen.Line(),
		jen.List(
			jen.ID(constants.ResponseVarName),
			jen.Err(),
		).Assign().ID("client").Dot("Do").Call(
			jen.ID(constants.RequestVarName).Dot("WithContext").Call(
				constants.CtxVar(),
			),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("executing request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.If(jen.ID("c").Dot("Debug")).Block(
			jen.List(
				jen.ID("bdump"),
				jen.Err(),
			).Assign().Qual("net/http/httputil", "DumpResponse").Call(
				jen.ID(constants.ResponseVarName),
				jen.True(),
			),
			jen.If(jen.Err().IsEqualTo().ID("nil").And().ID(constants.RequestVarName).Dot("Method").DoesNotEqual().Qual("net/http", "MethodGet")).Block(
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(
					jen.Lit("response_body"),
					jen.String().Call(
						jen.ID("bdump"),
					),
				),
			),
			jen.ID(constants.LoggerVarName).Dot("Debug").Call(
				jen.Lit("request executed"),
			),
		),
		jen.Line(),
		jen.Return().List(jen.ID(constants.ResponseVarName), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Comment("executeRawRequest takes a given *http.Request and executes it with the provided."),
		jen.Line(),
		jen.Comment("client, alongside some debugging logging."),
		jen.Line(),
		newClientMethod("executeRawRequest").Params(
			constants.CtxParam(),
			jen.ID("client").PointerTo().Qual("net/http", "Client"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params(jen.PointerTo().Qual("net/http", "Response"), jen.Error()).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildExportedBuildURL() []jen.Code {

	lines := []jen.Code{
		jen.Comment("BuildURL builds standard service URLs."),
		jen.Line(),
		newClientMethod("BuildURL").Params(
			jen.ID("qp").Qual("net/url", "Values"),
			jen.ID("parts").Spread().String(),
		).Params(jen.String()).Block(
			jen.Var().ID("u").PointerTo().Qual("net/url", "URL"),
			jen.If(jen.ID("qp").DoesNotEqual().ID("nil")).Block(
				jen.ID("u").Equals().ID("c").Dot("buildURL").Call(jen.ID("qp"), jen.ID("parts").Spread()),
			).Else().Block(
				jen.ID("u").Equals().ID("c").Dot("buildURL").Call(jen.Nil(), jen.ID("parts").Spread()),
			),
			jen.Line(),
			jen.If(jen.ID("u").DoesNotEqual().Nil()).Block(
				jen.Return(jen.ID("u").Dot("String").Call()),
			),
			jen.Return(jen.EmptyString()),
		),
		jen.Line(),
	}

	return lines
}

func buildUnexportedBuildURL() []jen.Code {
	lines := []jen.Code{
		jen.Comment("buildURL takes a given set of query parameters and URL parts, and returns."),
		jen.Line(),
		jen.Comment("a parsed URL object from them."),
		jen.Line(),
		newClientMethod("buildURL").Params(
			jen.ID("queryParams").Qual("net/url", "Values"),
			jen.ID("parts").Spread().String(),
		).Params(jen.PointerTo().Qual("net/url", "URL")).Block(
			jen.ID("tu").Assign().PointerTo().ID("c").Dot("URL"),
			jen.Line(),
			jen.ID("parts").Equals().ID("append").Call(
				jen.Index().String().Values(jen.Lit("api"), jen.Lit("v1")),
				jen.ID("parts").Spread(),
			),
			jen.List(
				jen.ID("u"),
				jen.Err(),
			).Assign().Qual("net/url", "Parse").Call(
				jen.Qual("strings", "Join").Call(
					jen.ID("parts"),
					jen.Lit("/"),
				),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("c").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("building URL")),
				jen.Return(jen.Nil()),
			),
			jen.Line(),
			jen.If(jen.ID("queryParams").DoesNotEqual().ID("nil")).Block(
				jen.ID("u").Dot("RawQuery").Equals().ID("queryParams").Dot("Encode").Call(),
			),
			jen.Line(),
			jen.Return().ID("tu").Dot("ResolveReference").Call(
				jen.ID("u"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildVersionlessURL() []jen.Code {
	lines := []jen.Code{
		jen.Comment("buildVersionlessURL builds a URL without the `/api/v1/` prefix. It should"),
		jen.Line(),
		jen.Comment("otherwise be identical to buildURL."),
		jen.Line(),
		newClientMethod("buildVersionlessURL").Params(
			jen.ID("qp").Qual("net/url", "Values"),
			jen.ID("parts").Spread().String(),
		).Params(jen.String()).Block(
			jen.ID("tu").Assign().PointerTo().ID("c").Dot("URL"),
			jen.Line(),
			jen.List(
				jen.ID("u"),
				jen.Err(),
			).Assign().Qual("net/url", "Parse").Call(
				jen.Qual("path", "Join").Call(
					jen.ID("parts").Spread(),
				),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("c").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("building URL")),
				jen.Return(jen.EmptyString()),
			),
			jen.Line(),
			jen.If(jen.ID("qp").DoesNotEqual().ID("nil")).Block(
				jen.ID("u").Dot("RawQuery").Equals().ID("qp").Dot("Encode").Call(),
			),
			jen.Line(),
			jen.Return().ID("tu").Dot("ResolveReference").Call(
				jen.ID("u"),
			).Dot("String").Call(),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildWebsocketURL() []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildWebsocketURL builds a standard URL and then converts its scheme to the websocket protocol."),
		jen.Line(),
		newClientMethod("BuildWebsocketURL").Params(
			jen.ID("parts").Spread().String(),
		).Params(jen.String()).Block(
			jen.ID("u").Assign().ID("c").Dot("buildURL").Call(
				jen.Nil(),
				jen.ID("parts").Spread(),
			),
			jen.ID("u").Dot("Scheme").Equals().Lit("ws"),
			jen.Line(),
			jen.Return().ID("u").Dot("String").Call(),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildHealthCheckRequest() []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildHealthCheckRequest builds a health check HTTP request."),
		jen.Line(),
		newClientMethod("BuildHealthCheckRequest").Params(constants.CtxParam()).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(
			jen.ID("u").Assign().PointerTo().ID("c").Dot("URL"),
			jen.ID("uri").Assign().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s://%s/_meta_/ready"),
				jen.ID("u").Dot("Scheme"),
				jen.ID("u").Dot("Host"),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequestWithContext").Call(
				constants.CtxVar(),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.Nil(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildIsUp() []jen.Code {
	lines := []jen.Code{
		jen.Comment("IsUp returns whether or not the service's health endpoint is returning 200s."),
		jen.Line(),
		newClientMethod("IsUp").Params(constants.CtxParam()).Params(jen.Bool()).Block(
			jen.List(
				jen.ID(constants.RequestVarName),
				jen.Err(),
			).Assign().ID("c").Dot("BuildHealthCheckRequest").Call(constants.CtxVar()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("c").Dot(constants.LoggerVarName).Dot("Error").Call(
					jen.Err(),
					jen.Lit("building request"),
				),
				jen.Return().False(),
			),
			jen.Line(),
			jen.List(
				jen.ID(constants.ResponseVarName),
				jen.Err(),
			).Assign().ID("c").Dot("plainClient").Dot("Do").Call(
				jen.ID(constants.RequestVarName),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("c").Dot(constants.LoggerVarName).Dot("Error").Call(
					jen.Err(),
					jen.Lit("health check"),
				),
				jen.Return().False(),
			),
			jen.ID("c").Dot("closeResponseBody").Call(jen.ID(constants.ResponseVarName)),
			jen.Line(),
			jen.Return().ID(constants.ResponseVarName).Dot("StatusCode").IsEqualTo().Qual("net/http", "StatusOK"),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildDataRequest(proj *models.Project) []jen.Code {
	block := []jen.Code{
		utils.StartSpan(proj, true, "buildDataRequest"),
		jen.List(
			jen.ID("body"),
			jen.Err(),
		).Assign().ID("createBodyFromStruct").Call(
			jen.ID("in"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(jen.Nil(), jen.Err()),
		),
		jen.Line(),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.ID("method"),
			jen.ID("uri"),
			jen.ID("body"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.ID(constants.RequestVarName).Dot("Header").Dot("Set").Call(
			jen.Lit("Content-type"),
			jen.Lit("application/json"),
		),
		jen.Return().List(jen.ID(constants.RequestVarName), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Comment("buildDataRequest builds an HTTP request for a given method, URL, and body data."),
		jen.Line(),
		newClientMethod("buildDataRequest").Params(
			constants.CtxParam(),
			jen.List(
				jen.ID("method"),
				jen.ID("uri"),
			).String(),
			jen.ID("in").Interface(),
		).Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildCheckExistence(proj *models.Project) []jen.Code {
	block := []jen.Code{
		utils.StartSpan(proj, true, "checkExistence"),
		jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().ID("c").Dot("executeRawRequest").Call(constants.CtxVar(), jen.ID("c").Dot("authedClient"), jen.ID(constants.RequestVarName)),
		jen.If(jen.Err().DoesNotEqual().Nil()).Block(
			jen.Return(jen.False(), jen.Err()),
		),
		jen.ID("c").Dot("closeResponseBody").Call(jen.ID(constants.ResponseVarName)),
		jen.Line(),
		jen.Return(jen.ID(constants.ResponseVarName).Dot("StatusCode").IsEqualTo().Qual("net/http", "StatusOK"), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Comment("checkExistence executes an HTTP request and loads the response content into a bool."),
		jen.Line(),
		newClientMethod("checkExistence").Params(
			constants.CtxParam(),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params(jen.Bool(), jen.Error()).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildRetrieve(proj *models.Project) []jen.Code {
	funcName := "retrieve"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.If(jen.Err().Assign().ID("argIsNotPointerOrNil").Call(
			jen.ID("obj"),
		),
			jen.Err().DoesNotEqual().ID("nil"),
		).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("struct to load must be a pointer: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.List(
			jen.ID(constants.ResponseVarName),
			jen.Err(),
		).Assign().ID("c").Dot("executeRawRequest").Call(
			constants.CtxVar(),
			jen.ID("c").Dot("authedClient"),
			jen.ID(constants.RequestVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("executing request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.If(jen.ID(constants.ResponseVarName).Dot("StatusCode").IsEqualTo().Qual("net/http", "StatusNotFound")).Block(
			jen.Return().ID("ErrNotFound"),
		),
		jen.Line(),
		jen.Return().ID("unmarshalBody").Call(constants.CtxVar(), jen.ID(constants.ResponseVarName), jen.AddressOf().ID("obj")),
	}

	lines := []jen.Code{
		jen.Commentf("%s executes an HTTP request and loads the response content into a struct. In the event of a 404,", funcName),
		jen.Line(),
		jen.Comment("the provided ErrNotFound is returned."),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			jen.ID("obj").Interface(),
		).Params(jen.Error()).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildExecuteRequest(proj *models.Project) []jen.Code {
	funcName := "executeRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.ResponseVarName),
			jen.Err(),
		).Assign().ID("c").Dot("executeRawRequest").Call(
			constants.CtxVar(),
			jen.ID("c").Dot("authedClient"),
			jen.ID(constants.RequestVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("executing request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Switch(jen.ID(constants.ResponseVarName).Dot("StatusCode")).Block(
			jen.Case(jen.Qual("net/http", "StatusNotFound")).Block(
				jen.Return().ID("ErrNotFound"),
			),
			jen.Case(jen.Qual("net/http", "StatusUnauthorized")).Block(
				jen.Return().ID("ErrUnauthorized"),
			),
		),
		jen.Line(),
		jen.If(jen.ID("out").DoesNotEqual().ID("nil")).Block(
			jen.If(
				jen.ID("resErr").Assign().ID("unmarshalBody").Call(constants.CtxVar(), jen.ID(constants.ResponseVarName), jen.ID("out")),
				jen.ID("resErr").DoesNotEqual().ID("nil"),
			).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("loading response from server: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Return().ID("nil"),
	}

	lines := []jen.Code{
		jen.Commentf("%s takes a given request and executes it with the auth client. It returns some errors", funcName),
		jen.Line(),
		jen.Comment("upon receiving certain status codes, but otherwise will return nil upon success."),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			jen.ID("out").Interface(),
		).Params(jen.Error()).Block(block...,
		),
		jen.Line(),
	}

	return lines
}

func buildExecuteUnauthenticatedDataRequest(proj *models.Project) []jen.Code {
	funcName := "executeUnauthenticatedDataRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Line(),
		jen.List(
			jen.ID(constants.ResponseVarName),
			jen.Err(),
		).Assign().ID("c").Dot("executeRawRequest").Call(
			constants.CtxVar(),
			jen.ID("c").Dot("plainClient"),
			jen.ID(constants.RequestVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("executing request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Switch(jen.ID(constants.ResponseVarName).Dot("StatusCode")).Block(
			jen.Case(jen.Qual("net/http", "StatusNotFound")).Block(
				jen.Return().ID("ErrNotFound"),
			),
			jen.Case(jen.Qual("net/http", "StatusUnauthorized")).Block(
				jen.Return().ID("ErrUnauthorized"),
			),
		),
		jen.Line(),
		jen.If(jen.ID("out").DoesNotEqual().ID("nil")).Block(
			jen.If(
				jen.ID("resErr").Assign().ID("unmarshalBody").Call(constants.CtxVar(), jen.ID(constants.ResponseVarName), jen.ID("out")),
				jen.ID("resErr").DoesNotEqual().ID("nil"),
			).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("loading response from server: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Return().ID("nil"),
	}

	lines := []jen.Code{
		jen.Commentf("%s takes a given request and loads the response into an interface value.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			jen.ID("out").Interface(),
		).Params(jen.Error()).Block(block...),
		jen.Line(),
	}

	return lines
}
