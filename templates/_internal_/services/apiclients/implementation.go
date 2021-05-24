package apiclients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func implementationDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("gopkg.in/oauth2.v3/server specific implementations"),
		jen.Line(),
	)

	code.Add(buildImplementationOAuth2InternalErrorHandler()...)
	code.Add(buildImplementationOAuth2ResponseErrorHandler()...)
	code.Add(buildImplementationAuthorizeScopeHandler(proj)...)
	code.Add(buildImplementationUserAuthorizationHandler(proj)...)
	code.Add(buildImplementationClientAuthorizedHandler(proj)...)
	code.Add(buildImplementationClientScopeHandler(proj)...)

	return code
}

func buildImplementationOAuth2InternalErrorHandler() []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual("gopkg.in/oauth2.v3/server", "InternalErrorHandler").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()).Dot("OAuth2InternalErrorHandler"),
		jen.Line(),
		jen.Comment("OAuth2InternalErrorHandler fulfills a role for the OAuth2 server-side provider"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("OAuth2InternalErrorHandler").Params(jen.Err().Error()).Params(jen.PointerTo().Qual("gopkg.in/oauth2.v3/errors", "Response")).Body(
			jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("OAuth2 Internal Error")),
			jen.Line(),
			jen.ID(constants.ResponseVarName).Assign().AddressOf().Qual("gopkg.in/oauth2.v3/errors", "Response").Valuesln(
				jen.ID("Error").MapAssign().Err(),
				jen.ID("Description").MapAssign().Lit("Internal error"),
				jen.ID("ErrorCode").MapAssign().Qual("net/http", "StatusInternalServerError"),
				jen.ID("StatusCode").MapAssign().Qual("net/http", "StatusInternalServerError"),
			),
			jen.Line(),
			jen.Return().ID(constants.ResponseVarName),
		),
		jen.Line(),
	}

	return lines
}

func buildImplementationOAuth2ResponseErrorHandler() []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual("gopkg.in/oauth2.v3/server", "ResponseErrorHandler").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()).Dot("OAuth2ResponseErrorHandler"),
		jen.Line(),
		jen.Comment("OAuth2ResponseErrorHandler fulfills a role for the OAuth2 server-side provider"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("OAuth2ResponseErrorHandler").Params(jen.ID("re").PointerTo().Qual("gopkg.in/oauth2.v3/errors", "Response")).Body(
			jen.ID("s").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("error_code").MapAssign().ID("re").Dot("ErrorCode"),
				jen.Lit("description").MapAssign().ID("re").Dot("Description"),
				jen.Lit("uri").MapAssign().ID("re").Dot("URI"),
				jen.Lit("status_code").MapAssign().ID("re").Dot("StatusCode"),
				jen.Lit("header").MapAssign().ID("re").Dot("Header"))).Dot("Error").Call(jen.ID("re").Dot("Error"), jen.Lit("OAuth2ResponseErrorHandler")),
		),
		jen.Line(),
	}

	return lines
}

func buildImplementationAuthorizeScopeHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual("gopkg.in/oauth2.v3/server", "AuthorizeScopeHandler").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()).Dot("AuthorizeScopeHandler"),
		jen.Line(),
		jen.Comment("AuthorizeScopeHandler satisfies the oauth2server AuthorizeScopeHandler interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("AuthorizeScopeHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.ID("scope").String(), jen.Err().Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("AuthorizeScopeHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.ID("scope").Equals().ID("determineScope").Call(jen.ID(constants.RequestVarName)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("scope"), jen.ID("scope")),
			jen.Line(),
			jen.Comment("check for client and return if valid."),
			jen.Var().ID("client").Equals().ID("s").Dot("fetchOAuth2ClientFromRequest").Call(jen.ID(constants.RequestVarName)),
			jen.If(jen.ID("client").DoesNotEqual().ID("nil").And().ID("client").Dot("HasScope").Call(jen.ID("scope"))).Body(
				jen.ID(constants.ResponseVarName).Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusOK")),
				jen.Return().List(jen.ID("scope"), jen.Nil()),
			),
			jen.Line(),
			jen.Comment("check to see if the client ID is present instead."),
			jen.If(jen.ID("clientID").Assign().ID("s").Dot("fetchOAuth2ClientIDFromRequest").Call(jen.ID(constants.RequestVarName)), jen.ID("clientID").DoesNotEqual().EmptyString()).Body(
				jen.Comment("fetch oauth2 client from database."),
				jen.List(jen.ID("client"), jen.Err()).Equals().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(constants.CtxVar(), jen.ID("clientID")),
				jen.Line(),
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Body(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 Client")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusNotFound"),
					jen.Return().List(jen.EmptyString(), jen.Err()),
				).Else().If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 Client")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return().List(jen.EmptyString(), jen.Err()),
				),
				jen.Line(),
				jen.Comment("authorization check."),
				jen.If(jen.Not().ID("client").Dot("HasScope").Call(jen.ID("scope"))).Body(
					utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
					jen.Return().List(jen.EmptyString(), utils.Error("not authorized for scope")),
				),
				jen.Line(),
				jen.Return().List(jen.ID("scope"), jen.Nil()),
			),
			jen.Line(),
			jen.Comment("invalid credentials."),
			utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
			jen.Return().List(jen.EmptyString(), utils.Error("no scope information found")),
		),
		jen.Line(),
	}

	return lines
}

func buildImplementationUserAuthorizationHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual("gopkg.in/oauth2.v3/server", "UserAuthorizationHandler").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()).Dot("UserAuthorizationHandler"),
		jen.Line(),
		jen.Comment("UserAuthorizationHandler satisfies the oauth2server UserAuthorizationHandler interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UserAuthorizationHandler").Params(
			jen.Underscore().Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params(
			jen.ID(constants.UserIDVarName).String(),
			jen.Err().Error(),
		).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("UserAuthorizationHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Var().ID("uid").Uint64(),
			jen.Line(),
			jen.Comment("check context for client."),
			jen.If(jen.List(jen.ID("client"), jen.ID("clientOk")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.Qual(proj.TypesPackage(), "OAuth2ClientKey")).Assert(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client")), jen.Not().ID("clientOk")).Body(
				jen.Comment("check for user instead."),
				jen.List(jen.ID("si"), jen.ID("userOk")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.Qual(proj.TypesPackage(), "SessionInfoKey")).Assert(jen.PointerTo().Qual(proj.TypesPackage(), "SessionInfo")),
				jen.If(jen.Not().ID("userOk").Or().ID("si").IsEqualTo().Nil()).Body(jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no user attached to this request")),
					jen.Return().List(jen.EmptyString(), utils.Error("user not found")),
				),
				jen.ID("uid").Equals().ID("si").Dot("UserID"),
			).Else().Body(
				jen.ID("uid").Equals().ID("client").Dot(constants.UserOwnershipFieldName),
			),
			jen.Line(),
			jen.Return().List(jen.Qual("strconv", "FormatUint").Call(jen.ID("uid"), jen.Lit(10)), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildImplementationClientAuthorizedHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual("gopkg.in/oauth2.v3/server", "ClientAuthorizedHandler").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()).Dot("ClientAuthorizedHandler"),
		jen.Line(),
		jen.Comment("ClientAuthorizedHandler satisfies the oauth2server ClientAuthorizedHandler interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ClientAuthorizedHandler").Params(jen.ID("clientID").String(), jen.ID("grant").Qual("gopkg.in/oauth2.v3", "GrantType")).Params(jen.ID("allowed").Bool(), jen.Err().Error()).Body(
			jen.Comment("NOTE: it's a shame the interface we're implementing doesn't have this as its first argument"),
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.InlineCtx(), jen.Lit("ClientAuthorizedHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("grant").MapAssign().ID("grant"), jen.Lit("client_id").MapAssign().ID("clientID"))),
			jen.Line(),
			jen.Comment("reject invalid grant type."),
			jen.If(jen.ID("grant").IsEqualTo().Qual("gopkg.in/oauth2.v3", "PasswordCredentials")).Body(
				jen.Return().List(jen.False(), utils.Error("invalid grant type: password")),
			),
			jen.Line(),
			jen.Comment("fetch client data."),
			jen.List(jen.ID("client"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(constants.CtxVar(), jen.ID("clientID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("fetching oauth2 client from database")),
				jen.Return().List(jen.False(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching oauth2 client from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Comment("disallow implicit grants unless authorized."),
			jen.If(jen.ID("grant").IsEqualTo().Qual("gopkg.in/oauth2.v3", "Implicit").And().Not().ID("client").Dot("ImplicitAllowed")).Body(
				jen.Return().List(jen.False(), utils.Error("client not authorized for implicit grants")),
			),
			jen.Line(),
			jen.Return().List(jen.True(), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildImplementationClientScopeHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual("gopkg.in/oauth2.v3/server", "ClientScopeHandler").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()).Dot("ClientScopeHandler"),
		jen.Line(),
		jen.Comment("ClientScopeHandler satisfies the oauth2server ClientScopeHandler interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ClientScopeHandler").Params(jen.List(jen.ID("clientID"), jen.ID("scope")).String()).Params(jen.ID("authed").Bool(), jen.Err().Error()).Body(
			jen.Comment("NOTE: it's a shame the interface we're implementing doesn't have this as its first argument"),
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.InlineCtx(), jen.Lit("UserAuthorizationHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("client_id").MapAssign().ID("clientID"),
				jen.Lit("scope").MapAssign().ID("scope")),
			),
			jen.Line(),
			jen.Comment("fetch client info."),
			jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(constants.CtxVar(), jen.ID("clientID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 client for ClientScopeHandler")),
				jen.Return().List(jen.False(), jen.Err()),
			),
			jen.Line(),
			jen.Comment("check for scope."),
			jen.If(jen.ID("c").Dot("HasScope").Call(jen.ID("scope"))).Body(
				jen.Return().List(jen.True(), jen.Nil()),
			),
			jen.Line(),
			jen.Return().List(jen.False(), utils.Error("unauthorized")),
		),
		jen.Line(),
	}

	return lines
}
