package oauth2clients

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func implementationDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Comment("gopkg.in/oauth2.v3/server specific implementations"),
		jen.Line(),
		jen.Line(),
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "InternalErrorHandler").Equals().Parens(jen.Op("*").ID("Service")).Call(jen.Nil()).Dot("OAuth2InternalErrorHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2InternalErrorHandler fulfills a role for the OAuth2 server-side provider"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("OAuth2InternalErrorHandler").Params(jen.Err().ID("error")).Params(jen.ParamPointer().Qual("gopkg.in/oauth2.v3/errors", "Response")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("OAuth2 Internal Error")),
			jen.Line(),
			jen.ID("res").Assign().VarPointer().Qual("gopkg.in/oauth2.v3/errors", "Response").Valuesln(
				jen.ID("Error").MapAssign().ID("err"),
				jen.ID("Description").MapAssign().Lit("Internal error"),
				jen.ID("ErrorCode").MapAssign().Qual("net/http", "StatusInternalServerError"),
				jen.ID("StatusCode").MapAssign().Qual("net/http", "StatusInternalServerError"),
			),
			jen.Line(),
			jen.Return().ID("res"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "ResponseErrorHandler").Equals().Parens(jen.Op("*").ID("Service")).Call(jen.Nil()).Dot("OAuth2ResponseErrorHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2ResponseErrorHandler fulfills a role for the OAuth2 server-side provider"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("OAuth2ResponseErrorHandler").Params(jen.ID("re").ParamPointer().Qual("gopkg.in/oauth2.v3/errors", "Response")).Block(
			jen.ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("error_code").MapAssign().ID("re").Dot("ErrorCode"),
				jen.Lit("description").MapAssign().ID("re").Dot("Description"),
				jen.Lit("uri").MapAssign().ID("re").Dot("URI"),
				jen.Lit("status_code").MapAssign().ID("re").Dot("StatusCode"),
				jen.Lit("header").MapAssign().ID("re").Dot("Header"))).Dot("Error").Call(jen.ID("re").Dot("Error"), jen.Lit("OAuth2ResponseErrorHandler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "AuthorizeScopeHandler").Equals().Parens(jen.Op("*").ID("Service")).Call(jen.Nil()).Dot("AuthorizeScopeHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("AuthorizeScopeHandler satisfies the oauth2server AuthorizeScopeHandler interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("AuthorizeScopeHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.ID("scope").ID("string"), jen.Err().ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("AuthorizeScopeHandler")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("scope").Equals().ID("determineScope").Call(jen.ID("req")),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("scope"), jen.ID("scope")).Dot("WithRequest").Call(jen.ID("req")),
			jen.Line(),
			jen.Comment("check for client and return if valid"),
			jen.Var().ID("client").Equals().ID("s").Dot("fetchOAuth2ClientFromRequest").Call(jen.ID("req")),
			jen.If(jen.ID("client").DoesNotEqual().ID("nil").Op("&&").ID("client").Dot("HasScope").Call(jen.ID("scope"))).Block(
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusOK")),
				jen.Return().List(jen.ID("scope"), jen.Nil()),
			),
			jen.Line(),
			jen.Comment("check to see if the client ID is present instead"),
			jen.If(jen.ID("clientID").Assign().ID("s").Dot("fetchOAuth2ClientIDFromRequest").Call(jen.ID("req")), jen.ID("clientID").DoesNotEqual().Lit("")).Block(
				jen.Comment("fetch oauth2 client from database"),
				jen.List(jen.ID("client"), jen.Err()).Equals().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("clientID")),
				jen.Line(),
				jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 Client")),
					utils.WriteXHeader("res", "StatusNotFound"),
					jen.Return().List(jen.Lit(""), jen.Err()),
				).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 Client")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return().List(jen.Lit(""), jen.Err()),
				),
				jen.Line(),
				jen.Comment("authorization check"),
				jen.If(jen.Op("!").ID("client").Dot("HasScope").Call(jen.ID("scope"))).Block(
					utils.WriteXHeader("res", "StatusUnauthorized"),
					jen.Return().List(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("not authorized for scope"))),
				),
				jen.Line(),
				jen.Return().List(jen.ID("scope"), jen.Nil()),
			),
			jen.Line(),
			jen.Comment("invalid credentials"),
			utils.WriteXHeader("res", "StatusBadRequest"),
			jen.Return().List(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("no scope information found"))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "UserAuthorizationHandler").Equals().Parens(jen.Op("*").ID("Service")).Call(jen.Nil()).Dot("UserAuthorizationHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UserAuthorizationHandler satisfies the oauth2server UserAuthorizationHandler interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UserAuthorizationHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.ID("userID").ID("string"), jen.Err().ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UserAuthorizationHandler")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().ID("uid").ID("uint64"),
			jen.Line(),
			jen.Comment("check context for client"),
			jen.If(jen.List(jen.ID("client"), jen.ID("clientOk")).Assign().ID(utils.ContextVarName).Dot("Value").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientKey")).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")), jen.Op("!").ID("clientOk")).Block(
				jen.Comment("check for user instead"),
				jen.List(jen.ID("user"), jen.ID("userOk")).Assign().ID(utils.ContextVarName).Dot("Value").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserKey")).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User")),
				jen.If(jen.Op("!").ID("userOk")).Block(jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("no user attached to this request")),
					jen.Return().List(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("user not found"))),
				),
				jen.ID("uid").Equals().ID("user").Dot("ID"),
			).Else().Block(
				jen.ID("uid").Equals().ID("client").Dot("BelongsToUser"),
			),
			jen.Line(),
			jen.Return().List(jen.Qual("strconv", "FormatUint").Call(jen.ID("uid"), jen.Lit(10)), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "ClientAuthorizedHandler").Equals().Parens(jen.Op("*").ID("Service")).Call(jen.Nil()).Dot("ClientAuthorizedHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ClientAuthorizedHandler satisfies the oauth2server ClientAuthorizedHandler interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ClientAuthorizedHandler").Params(jen.ID("clientID").ID("string"), jen.ID("grant").Qual("gopkg.in/oauth2.v3", "GrantType")).Params(jen.ID("allowed").ID("bool"), jen.Err().ID("error")).Block(
			jen.Comment("NOTE: it's a shame the interface we're implementing doesn't have this as its first argument"),
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(utils.InlineCtx(), jen.Lit("ClientAuthorizedHandler")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("grant").MapAssign().ID("grant"), jen.Lit("client_id").MapAssign().ID("clientID"))),
			jen.Line(),
			jen.Comment("reject invalid grant type"),
			jen.If(jen.ID("grant").Op("==").Qual("gopkg.in/oauth2.v3", "PasswordCredentials")).Block(
				jen.Return().List(jen.ID("false"), jen.Qual("errors", "New").Call(jen.Lit("invalid grant type: password"))),
			),
			jen.Line(),
			jen.Comment("fetch client data"),
			jen.List(jen.ID("client"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("clientID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("fetching oauth2 client from database")),
				jen.Return().List(jen.ID("false"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching oauth2 client from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Comment("disallow implicit grants unless authorized"),
			jen.If(jen.ID("grant").Op("==").Qual("gopkg.in/oauth2.v3", "Implicit").Op("&&").Op("!").ID("client").Dot("ImplicitAllowed")).Block(
				jen.Return().List(jen.ID("false"), jen.Qual("errors", "New").Call(jen.Lit("client not authorized for implicit grants"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("true"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "ClientScopeHandler").Equals().Parens(jen.Op("*").ID("Service")).Call(jen.Nil()).Dot("ClientScopeHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ClientScopeHandler satisfies the oauth2server ClientScopeHandler interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ClientScopeHandler").Params(jen.List(jen.ID("clientID"), jen.ID("scope")).ID("string")).Params(jen.ID("authed").ID("bool"), jen.Err().ID("error")).Block(
			jen.Comment("NOTE: it's a shame the interface we're implementing doesn't have this as its first argument"),
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(utils.InlineCtx(), jen.Lit("UserAuthorizationHandler")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("client_id").MapAssign().ID("clientID"),
				jen.Lit("scope").MapAssign().ID("scope")),
			),
			jen.Line(),
			jen.Comment("fetch client info"),
			jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("clientID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 client for ClientScopeHandler")),
				jen.Return().List(jen.ID("false"), jen.Err()),
			),
			jen.Line(),
			jen.Comment("check for scope"),
			jen.If(jen.ID("c").Dot("HasScope").Call(jen.ID("scope"))).Block(
				jen.Return().List(jen.ID("true"), jen.Nil()),
			),
			jen.Line(),
			jen.Return().List(jen.ID("false"), jen.Qual("errors", "New").Call(jen.Lit("unauthorized"))),
		),
		jen.Line(),
	)
	return ret
}
