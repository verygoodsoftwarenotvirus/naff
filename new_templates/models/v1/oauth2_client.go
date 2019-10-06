package models

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("OAuth2ClientKey").ID("ContextKey").Op("=").Lit("oauth2_client"))
	ret.Add(jen.Null().Type().ID("OAuth2ClientDataManager").Interface(jen.ID("GetOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("OAuth2Client"), jen.ID("error")), jen.ID("GetOAuth2ClientByClientID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.Op("*").ID("OAuth2Client"), jen.ID("error")), jen.ID("GetAllOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")), jen.ID("GetOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")), jen.ID("GetOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("OAuth2ClientList"), jen.ID("error")), jen.ID("GetAllOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Index().Op("*").ID("OAuth2Client"), jen.ID("error")), jen.ID("GetAllOAuth2ClientsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Op("*").ID("OAuth2Client"), jen.ID("error")), jen.ID("CreateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("OAuth2ClientCreationInput")).Params(jen.Op("*").ID("OAuth2Client"), jen.ID("error")), jen.ID("UpdateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("OAuth2Client")).Params(jen.ID("error")), jen.ID("ArchiveOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error"))).Type().ID("OAuth2ClientDataServer").Interface(jen.ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")), jen.ID("OAuth2ClientInfoMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")), jen.ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("OAuth2Client"), jen.ID("error")), jen.ID("HandleAuthorizeRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")), jen.ID("HandleTokenRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error"))).Type().ID("OAuth2Client").Struct(
		jen.ID("ID").ID("uint64"),
		jen.ID("Name").ID("string"),
		jen.ID("ClientID").ID("string"),
		jen.ID("ClientSecret").ID("string"),
		jen.ID("RedirectURI").ID("string"),
		jen.ID("Scopes").Index().ID("string"),
		jen.ID("ImplicitAllowed").ID("bool"),
		jen.ID("CreatedOn").ID("uint64"),
		jen.ID("UpdatedOn").Op("*").ID("uint64"),
		jen.ID("ArchivedOn").Op("*").ID("uint64"),
		jen.ID("BelongsTo").ID("uint64"),
	).Type().ID("OAuth2ClientList").Struct(
		jen.ID("Pagination"),
		jen.ID("Clients").Index().ID("OAuth2Client"),
	).Type().ID("OAuth2ClientCreationInput").Struct(
		jen.ID("UserLoginInput"),
		jen.ID("Name").ID("string"),
		jen.ID("ClientID").ID("string"),
		jen.ID("ClientSecret").ID("string"),
		jen.ID("RedirectURI").ID("string"),
		jen.ID("BelongsTo").ID("uint64"),
		jen.ID("Scopes").Index().ID("string"),
	).Type().ID("OAuth2ClientUpdateInput").Struct(
		jen.ID("RedirectURI").ID("string"),
		jen.ID("Scopes").Index().ID("string"),
	),
	)
	ret.Add(jen.Null().Var().ID("_").ID("oauth2").Dot(
		"ClientInfo",
	).Op("=").Parens(jen.Op("*").ID("OAuth2Client")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
