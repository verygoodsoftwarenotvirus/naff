package oauth2clients

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func implementationDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Comment("gopkg.in/oauth2.v3/server specific implementations"),
		jen.Line(),
		jen.Line(),
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "InternalErrorHandler").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")).Dot("OAuth2InternalErrorHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2InternalErrorHandler fulfills a role for the OAuth2 server-side provider"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("OAuth2InternalErrorHandler").Params(jen.ID("err").ID("error")).Params(jen.Op("*").Qual("gopkg.in/oauth2.v3/errors", "Response")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("OAuth2 Internal Error")),
			jen.Line(),
			jen.ID("res").Op(":=").Op("&").Qual("gopkg.in/oauth2.v3/errors", "Response").Valuesln(
				jen.ID("Error").Op(":").ID("err"),
				jen.ID("Description").Op(":").Lit("Internal error"),
				jen.ID("ErrorCode").Op(":").Qual("net/http", "StatusInternalServerError"),
				jen.ID("StatusCode").Op(":").Qual("net/http", "StatusInternalServerError"),
			),
			jen.Line(),
			jen.Return().ID("res"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "ResponseErrorHandler").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")).Dot("OAuth2ResponseErrorHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2ResponseErrorHandler fulfills a role for the OAuth2 server-side provider"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("OAuth2ResponseErrorHandler").Params(jen.ID("re").Op("*").Qual("gopkg.in/oauth2.v3/errors", "Response")).Block(
			jen.ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("error_code").Op(":").ID("re").Dot("ErrorCode"),
				jen.Lit("description").Op(":").ID("re").Dot("Description"),
				jen.Lit("uri").Op(":").ID("re").Dot("URI"),
				jen.Lit("status_code").Op(":").ID("re").Dot("StatusCode"),
				jen.Lit("header").Op(":").ID("re").Dot("Header"))).Dot("Error").Call(jen.ID("re").Dot("Error"), jen.Lit("OAuth2ResponseErrorHandler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "AuthorizeScopeHandler").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")).Dot("AuthorizeScopeHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("AuthorizeScopeHandler satisfies the oauth2server AuthorizeScopeHandler interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("AuthorizeScopeHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("scope").ID("string"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("AuthorizeScopeHandler")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("scope").Op("=").ID("determineScope").Call(jen.ID("req")),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("scope"), jen.ID("scope")).Dot("WithRequest").Call(jen.ID("req")),
			jen.Line(),
			jen.Comment("check for client and return if valid"),
			jen.Var().ID("client").Op("=").ID("s").Dot("fetchOAuth2ClientFromRequest").Call(jen.ID("req")),
			jen.If(jen.ID("client").Op("!=").ID("nil").Op("&&").ID("client").Dot("HasScope").Call(jen.ID("scope"))).Block(
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusOK")),
				jen.Return().List(jen.ID("scope"), jen.ID("nil")),
			),
			jen.Line(),
			jen.Comment("check to see if the client ID is present instead"),
			jen.If(jen.ID("clientID").Op(":=").ID("s").Dot("fetchOAuth2ClientIDFromRequest").Call(jen.ID("req")), jen.ID("clientID").Op("!=").Lit("")).Block(
				jen.Comment("fetch oauth2 client from database"),
				jen.List(jen.ID("client"), jen.ID("err")).Op("=").ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(jen.ID("ctx"), jen.ID("clientID")),
				jen.Line(),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error fetching OAuth2 Client")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
					jen.Return().List(jen.Lit(""), jen.ID("err")),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error fetching OAuth2 Client")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return().List(jen.Lit(""), jen.ID("err")),
				),
				jen.Line(),
				jen.Comment("authorization check"),
				jen.If(jen.Op("!").ID("client").Dot("HasScope").Call(jen.ID("scope"))).Block(
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return().List(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("not authorized for scope"))),
				),
				jen.Line(),
				jen.Return().List(jen.ID("scope"), jen.ID("nil")),
			),
			jen.Line(),
			jen.Comment("invalid credentials"),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
			jen.Return().List(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("no scope information found"))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "UserAuthorizationHandler").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")).Dot("UserAuthorizationHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UserAuthorizationHandler satisfies the oauth2server UserAuthorizationHandler interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UserAuthorizationHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("userID").ID("string"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UserAuthorizationHandler")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().ID("uid").ID("uint64"),
			jen.Line(),
			jen.Comment("check context for client"),
			jen.If(jen.List(jen.ID("client"), jen.ID("clientOk")).Op(":=").ID("ctx").Dot("Value").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientKey")).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")), jen.Op("!").ID("clientOk")).Block(
				jen.Comment("check for user instead"),
				jen.List(jen.ID("user"), jen.ID("userOk")).Op(":=").ID("ctx").Dot("Value").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserKey")).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User")),
				jen.If(jen.Op("!").ID("userOk")).Block(jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("no user attached to this request")),
					jen.Return().List(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("user not found"))),
				),
				jen.ID("uid").Op("=").ID("user").Dot("ID"),
			).Else().Block(
				jen.ID("uid").Op("=").ID("client").Dot("BelongsToUser"),
			),
			jen.Line(),
			jen.Return().List(jen.Qual("strconv", "FormatUint").Call(jen.ID("uid"), jen.Lit(10)), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "ClientAuthorizedHandler").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")).Dot("ClientAuthorizedHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ClientAuthorizedHandler satisfies the oauth2server ClientAuthorizedHandler interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ClientAuthorizedHandler").Params(jen.ID("clientID").ID("string"), jen.ID("grant").Qual("gopkg.in/oauth2.v3", "GrantType")).Params(jen.ID("allowed").ID("bool"), jen.ID("err").ID("error")).Block(
			jen.Comment("NOTE: it's a shame the interface we're implementing doesn't have this as its first argument"),
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.Qual("context", "Background").Call(), jen.Lit("ClientAuthorizedHandler")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("grant").Op(":").ID("grant"), jen.Lit("client_id").Op(":").ID("clientID"))),
			jen.Line(),
			jen.Comment("reject invalid grant type"),
			jen.If(jen.ID("grant").Op("==").Qual("gopkg.in/oauth2.v3", "PasswordCredentials")).Block(
				jen.Return().List(jen.ID("false"), jen.Qual("errors", "New").Call(jen.Lit("invalid grant type: password"))),
			),
			jen.Line(),
			jen.Comment("fetch client data"),
			jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(jen.ID("ctx"), jen.ID("clientID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("fetching oauth2 client from database")),
				jen.Return().List(jen.ID("false"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching oauth2 client from database: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.Comment("disallow implicit grants unless authorized"),
			jen.If(jen.ID("grant").Op("==").Qual("gopkg.in/oauth2.v3", "Implicit").Op("&&").Op("!").ID("client").Dot("ImplicitAllowed")).Block(
				jen.Return().List(jen.ID("false"), jen.Qual("errors", "New").Call(jen.Lit("client not authorized for implicit grants"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("true"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3/server", "ClientScopeHandler").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")).Dot("ClientScopeHandler"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ClientScopeHandler satisfies the oauth2server ClientScopeHandler interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ClientScopeHandler").Params(jen.List(jen.ID("clientID"), jen.ID("scope")).ID("string")).Params(jen.ID("authed").ID("bool"), jen.ID("err").ID("error")).Block(
			jen.Comment("NOTE: it's a shame the interface we're implementing doesn't have this as its first argument"),
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.Qual("context", "Background").Call(), jen.Lit("UserAuthorizationHandler")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("client_id").Op(":").ID("clientID"),
				jen.Lit("scope").Op(":").ID("scope")),
			),
			jen.Line(),
			jen.Comment("fetch client info"),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(jen.ID("ctx"), jen.ID("clientID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error fetching OAuth2 client for ClientScopeHandler")),
				jen.Return().List(jen.ID("false"), jen.ID("err")),
			),
			jen.Line(),
			jen.Comment("check for scope"),
			jen.If(jen.ID("c").Dot("HasScope").Call(jen.ID("scope"))).Block(
				jen.Return().List(jen.ID("true"), jen.ID("nil")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("false"), jen.Qual("errors", "New").Call(jen.Lit("unauthorized"))),
		),
		jen.Line(),
	)
	return ret
}
