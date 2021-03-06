package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildMiddlewareConstantDefs()...)
	code.Add(buildCreationInputMiddleware(proj)...)
	code.Add(buildExtractOAuth2ClientFromRequest(proj)...)
	code.Add(buildDetermineScope()...)
	code.Add(buildOAuth2TokenAuthenticationMiddleware(proj)...)
	code.Add(buildOAuth2ClientInfoMiddleware(proj)...)
	code.Add(buildServiceFetchOAuth2ClientFromRequest(proj)...)
	code.Add(buildServiceFetchOAuth2ClientIDFromRequest()...)

	return code
}

func buildMiddlewareConstantDefs() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("scopesSeparator").Equals().Lit(","),
			jen.ID("apiPathPrefix").Equals().Lit("/api/v1/"),
		),
		jen.Line(),
	}

	return lines
}

func buildCreationInputMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreationInputMiddleware is a middleware for attaching OAuth2 client info to a request."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("CreationInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")),
				jen.Line(),
				jen.Comment("decode value from request."),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("creationMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildExtractOAuth2ClientFromRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ExtractOAuth2ClientFromRequest extracts OAuth2 client data from a request."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ExtractOAuth2ClientFromRequest").Params(constants.CtxParam(), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("ExtractOAuth2ClientFromRequest")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("validate bearer token."),
			jen.List(jen.ID("token"), jen.Err()).Assign().ID("s").Dot("oauth2Handler").Dot("ValidationBearerToken").Call(jen.ID(constants.RequestVarName)),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("validating bearer token: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Comment("fetch client ID."),
			jen.ID("clientID").Assign().ID("token").Dot("GetClientID").Call(),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("client_id"), jen.ID("clientID")),
			jen.Line(),
			jen.Comment("fetch client by client ID."),
			jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(constants.CtxVar(), jen.ID("clientID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 Client")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Comment("determine the scope."),
			jen.ID("scope").Assign().ID("determineScope").Call(jen.ID(constants.RequestVarName)),
			jen.ID("hasScope").Assign().ID("c").Dot("HasScope").Call(jen.ID("scope")),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("scope"), jen.ID("scope")).Dot("WithValue").Call(jen.Lit("scopes"), jen.Qual("strings", "Join").Call(jen.ID("c").Dot("Scopes"), jen.ID("scopesSeparator"))),
			jen.Line(),
			jen.If(jen.Not().ID("hasScope")).Body(
				jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("rejecting client for invalid scope")),
				jen.Return().List(jen.Nil(), utils.Error("client not authorized for scope")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("c"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildDetermineScope() []jen.Code {
	lines := []jen.Code{
		jen.Comment("determineScope determines the scope of a request by its URL."),
		jen.Line(),
		jen.Func().ID("determineScope").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.String()).Body(
			jen.If(jen.Qual("strings", "HasPrefix").Call(jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
				jen.ID("apiPathPrefix"))).Body(
				jen.ID("x").Assign().Qual("strings", "TrimPrefix").Call(jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"), jen.ID("apiPathPrefix")),
				jen.If(jen.ID("y").Assign().Qual("strings", "Split").Call(jen.ID("x"), jen.Lit("/")), jen.Len(jen.ID("y")).GreaterThan().Zero()).Body(
					jen.ID("x").Equals().ID("y").Index(jen.Zero()),
				),
				jen.Return().ID("x"),
			),
			jen.Line(),
			jen.Return().EmptyString(),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2TokenAuthenticationMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("OAuth2TokenAuthenticationMiddleware authenticates Oauth tokens."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("OAuth2TokenAuthenticationMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("OAuth2TokenAuthenticationMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("ExtractOAuth2ClientFromRequest").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error authenticated token-authed request")),
					jen.Qual("net/http", "Error").Call(jen.ID(constants.ResponseVarName), jen.Lit("invalid token"), jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return(),
				),
				jen.Line(),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("c").Dot(constants.UserOwnershipFieldName)),
				jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("c").Dot("ClientID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientDatabaseIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("c").Dot("ID")),
				jen.Line(),
				jen.Comment("attach the client object to the request."),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"), jen.ID("c")),
				jen.Line(),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2ClientInfoMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("OAuth2ClientInfoMiddleware fetches clientOAuth2Client info from requests and attaches it explicitly to a request."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("OAuth2ClientInfoMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("OAuth2ClientInfoMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.If(jen.ID("v").Assign().ID(constants.RequestVarName).Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.ID("oauth2ClientIDURIParamKey")), jen.ID("v").DoesNotEqual().EmptyString()).Body(
					jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("oauth2_client_id"), jen.ID("v")),
					jen.Line(),
					jen.List(jen.ID("client"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(constants.CtxVar(), jen.ID("v")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
						jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 client")),
						jen.Qual("net/http", "Error").Call(jen.ID(constants.ResponseVarName), jen.Lit("invalid request"), jen.Qual("net/http", "StatusUnauthorized")),
						jen.Return(),
					),
					jen.Line(),
					jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("client").Dot(constants.UserOwnershipFieldName)),
					jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("client").Dot("ClientID")),
					jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientDatabaseIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("client").Dot("ID")),
					jen.Line(),
					constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"), jen.ID("client")),
					jen.Line(),
					jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar()),
				),
				jen.Line(),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceFetchOAuth2ClientFromRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("fetchOAuth2ClientFromRequest").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Body(
			jen.List(jen.ID("client"), jen.ID("ok")).Assign().ID(constants.RequestVarName).Dot("Context").Call().Dot("Value").Call(jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")),
			jen.Underscore().Equals().ID("ok").Comment("we don't really care, but the linters do"),
			jen.Return().ID("client"),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceFetchOAuth2ClientIDFromRequest() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("fetchOAuth2ClientIDFromRequest").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.String()).Body(
			jen.List(jen.ID("clientID"), jen.ID("ok")).Assign().ID(constants.RequestVarName).Dot("Context").Call().Dot("Value").Call(jen.ID("clientIDKey")).Assert(jen.String()),
			jen.Underscore().Equals().ID("ok").Comment("we don't really care, but the linters do"),
			jen.Return().ID("clientID"),
		),
		jen.Line(),
	}

	return lines
}
