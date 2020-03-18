package oauth2clients

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("URIParamKey is used for referring to OAuth2 client IDs in router params"),
			jen.ID("URIParamKey").Op("=").Lit("oauth2ClientID"),
			jen.Line(),
			jen.ID("oauth2ClientIDURIParamKey").Op("=").Lit("client_id"),
			jen.ID("clientIDKey").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ContextKey").Op("=").Lit("client_id"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachUserIDToSpan provides a consistent way of attaching an user ID to a given span"),
		jen.Line(),
		jen.Func().ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("userID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachOAuth2ClientDatabaseIDToSpan provides a consistent way of attaching an oauth2 client ID to a given span"),
		jen.Line(),
		jen.Func().ID("attachOAuth2ClientDatabaseIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("clientID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("oauth2client_db_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("clientID"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachOAuth2ClientIDToSpan provides a consistent way of attaching a client ID to a given span"),
		jen.Line(),
		jen.Func().ID("attachOAuth2ClientIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("clientID").ID("string")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("client_id"), jen.ID("clientID"))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("randString produces a random string"),
		jen.Line(),
		jen.Comment("https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/"),
		jen.Line(),
		jen.Func().ID("randString").Params().Params(jen.ID("string")).Block(
			jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.Lit(32)),
			jen.If(jen.List(jen.ID("_"), jen.Err()).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.Comment("this is so that we don't end up with `=` in IDs"),
			jen.Return().Qual("encoding/base32", "StdEncoding").Dot("WithPadding").Call(jen.Qual("encoding/base32", "NoPadding")).Dot("EncodeToString").Call(jen.ID("b")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("fetchUserID grabs a userID out of the request context"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("fetchUserID").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.If(jen.List(jen.ID("id"), jen.ID("ok")).Op(":=").ID("req").Dot("Context").Call().Dot("Value").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserIDKey")).Assert(jen.ID("uint64")), jen.ID("ok")).Block(
				jen.Return().ID("id"),
			),
			jen.Return().Lit(0),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler is a handler that returns a list of OAuth2 clients"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ListHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("extract filter"),
				jen.ID("qf").Op(":=").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ExtractQueryFilter").Call(jen.ID("req")),
				jen.Line(),
				jen.Comment("determine user"),
				jen.ID("userID").Op(":=").ID("s").Dot("fetchUserID").Call(jen.ID("req")),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("fetch oauth2 clients"),
				jen.List(jen.ID("oauth2Clients"), jen.Err()).Op(":=").ID("s").Dot("database").Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.ID("qf"), jen.ID("userID")),
				jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Comment("just return an empty list if there are no results"),
					jen.ID("oauth2Clients").Op("=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientList").Valuesln(
						jen.ID("Clients").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Values(),
					),
				).Else().If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("encountered error getting list of oauth2 clients from database")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode response and peace"),
				jen.If(jen.Err().Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("oauth2Clients")), jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateHandler is our OAuth2 client creation route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CreateHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("fetch creation input from request context"),
				jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID(utils.ContextVarName).Dot("Value").Call(jen.ID("CreationMiddlewareCtxKey")).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientCreationInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot("logger").Dot("Info").Call(jen.Lit("valid input not attached to request")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("keep relevant data in mind"),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("username").Op(":").ID("input").Dot("Username"),
					jen.Lit("scopes").Op(":").ID("input").Dot("Scopes"),
					jen.Lit("redirect_uri").Op(":").ID("input").Dot("RedirectURI"))),
				jen.Line(),
				jen.Comment("retrieve user"),
				jen.List(jen.ID("user"), jen.Err()).Op(":=").ID("s").Dot("database").Dot("GetUserByUsername").Call(utils.CtxVar(), jen.ID("input").Dot("Username")),
				jen.If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("fetching user by username")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.ID("input").Dot("BelongsToUser").Op("=").ID("user").Dot("ID"),
				jen.Line(),
				jen.Comment("tag span since we have the info"),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("user").Dot("ID")),
				jen.Line(),
				jen.Comment("check credentials"),
				jen.List(jen.ID("valid"), jen.Err()).Op(":=").ID("s").Dot("authenticator").Dot("ValidateLogin").Callln(
					utils.CtxVar(), jen.ID("user").Dot("HashedPassword"),
					jen.ID("input").Dot("Password"),
					jen.ID("user").Dot("TwoFactorSecret"),
					jen.ID("input").Dot("TOTPToken"),
					jen.ID("user").Dot("Salt"),
				),
				jen.Line(),
				jen.If(jen.Op("!").ID("valid")).Block(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("invalid credentials provided")),
					utils.WriteXHeader("res", "StatusUnauthorized"),
					jen.Return(),
				).Else().If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("validating user credentials")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("set some data"),
				jen.ID("input").Dot("ClientID").Op("=").ID("randString").Call(),
				jen.ID("input").Dot("ClientSecret").Op("=").ID("randString").Call(),
				jen.ID("input").Dot("BelongsToUser").Op("=").ID("s").Dot("fetchUserID").Call(jen.ID("req")),
				jen.Line(),
				jen.Comment("create the client"),
				jen.List(jen.ID("client"), jen.Err()).Op(":=").ID("s").Dot("database").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("input")),
				jen.If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("creating oauth2Client in the database")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("notify interested parties"),
				jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("client").Dot("ID")),
				jen.ID("s").Dot("oauth2ClientCounter").Dot("Increment").Call(utils.CtxVar()),
				jen.Line(),
				utils.WriteXHeader("res", "StatusCreated"),
				jen.If(jen.Err().Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("client")), jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler is a route handler for retrieving an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ReadHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("determine subject of request"),
				jen.ID("userID").Op(":=").ID("s").Dot("fetchUserID").Call(jen.ID("req")),
				jen.ID("oauth2ClientID").Op(":=").ID("s").Dot("urlClientIDExtractor").Call(jen.ID("req")),
				jen.Line(),
				jen.Comment("keep the aforementioned in mind"),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("oauth2_client_id").Op(":").ID("oauth2ClientID"), jen.Lit("user_id").Op(":").ID("userID"))),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("oauth2ClientID")),
				jen.Line(),
				jen.Comment("fetch oauth2 client"),
				jen.List(jen.ID("x"), jen.Err()).Op(":=").ID("s").Dot("database").Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("oauth2ClientID"), jen.ID("userID")),
				jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("ReadHandler called on nonexistent client")),
					utils.WriteXHeader("res", "StatusNotFound"),
					jen.Return(),
				).Else().If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching oauth2Client from database")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode response and peace"),
				jen.If(jen.Err().Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler is a route handler for archiving an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ArchiveHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("determine subject matter"),
				jen.ID("userID").Op(":=").ID("s").Dot("fetchUserID").Call(jen.ID("req")),
				jen.ID("oauth2ClientID").Op(":=").ID("s").Dot("urlClientIDExtractor").Call(jen.ID("req")),
				jen.Line(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("oauth2_client_id").Op(":").ID("oauth2ClientID"), jen.Lit("user_id").Op(":").ID("userID"))),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("oauth2ClientID")),
				jen.Line(),
				jen.Comment("mark client as archived"),
				jen.Err().Op(":=").ID("s").Dot("database").Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("oauth2ClientID"), jen.ID("userID")),
				jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
					utils.WriteXHeader("res", "StatusNotFound"),
					jen.Return(),
				).Else().If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("encountered error deleting oauth2 client")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("notify relevant parties"),
				jen.ID("s").Dot("oauth2ClientCounter").Dot("Decrement").Call(utils.CtxVar()),
				utils.WriteXHeader("res", "StatusNoContent"),
			),
		),
		jen.Line(),
	)
	return ret
}
