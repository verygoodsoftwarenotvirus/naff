package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("OAuth2ClientKey is a ContextKey for use with contexts involving OAuth2 clients"),
			jen.ID("OAuth2ClientKey").ID("ContextKey").Op("=").Lit("oauth2_client"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("OAuth2ClientDataManager handles OAuth2 clients"),
			jen.ID("OAuth2ClientDataManager").Interface(
				jen.ID("GetOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("OAuth2Client"), jen.ID("error")),
				jen.ID("GetOAuth2ClientByClientID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.Op("*").ID("OAuth2Client"), jen.ID("error")),
				jen.ID("GetAllOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("OAuth2ClientList"), jen.ID("error")),
				jen.ID("GetAllOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Index().Op("*").ID("OAuth2Client"), jen.ID("error")),
				jen.ID("GetAllOAuth2ClientsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Op("*").ID("OAuth2Client"), jen.ID("error")),
				jen.ID("CreateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("OAuth2ClientCreationInput")).Params(jen.Op("*").ID("OAuth2Client"), jen.ID("error")),
				jen.ID("UpdateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("OAuth2Client")).Params(jen.ID("error")),
				jen.ID("ArchiveOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("OAuth2ClientDataServer describes a structure capable of serving traffic related to oauth2 clients"),
			jen.ID("OAuth2ClientDataServer").Interface(
				jen.ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.Comment("There is deliberately no update function"),
				jen.ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.Line(),
				jen.ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("OAuth2ClientInfoMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("OAuth2Client"), jen.ID("error")),
				jen.Line(),
				jen.Comment("wrappers for our implementation library"),
				jen.ID("HandleAuthorizeRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")),
				jen.ID("HandleTokenRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("OAuth2Client represents a user-authorized API client"),
			jen.ID("OAuth2Client").Struct(
				jen.ID("ID").ID("uint64").Tag(jsonTag("id")),
				jen.ID("Name").ID("string").Tag(jsonTag("name")),
				jen.ID("ClientID").ID("string").Tag(jsonTag("client_id")),
				jen.ID("ClientSecret").ID("string").Tag(jsonTag("client_secret")),
				jen.ID("RedirectURI").ID("string").Tag(jsonTag("redirect_uri")),
				jen.ID("Scopes").Index().ID("string").Tag(jsonTag("scopes")),
				jen.ID("ImplicitAllowed").ID("bool").Tag(jsonTag("implicit_allowed")),
				jen.ID("CreatedOn").ID("uint64").Tag(jsonTag("created_on")),
				jen.ID("UpdatedOn").Op("*").ID("uint64").Tag(jsonTag("updated_on")),
				jen.ID("ArchivedOn").Op("*").ID("uint64").Tag(jsonTag("archived_on")),
				jen.ID("BelongsToUser").ID("uint64").Tag(jsonTag("belongs_to_user")),
			),
			jen.Line(),
			jen.Comment("OAuth2ClientList is a response struct containing a list of OAuth2Clients"),
			jen.ID("OAuth2ClientList").Struct(
				jen.ID("Pagination"),
				jen.ID("Clients").Index().ID("OAuth2Client").Tag(jsonTag("clients")),
			),
			jen.Line(),
			jen.Comment("OAuth2ClientCreationInput is a struct for use when creating OAuth2 clients."),
			jen.ID("OAuth2ClientCreationInput").Struct(
				jen.ID("UserLoginInput"),
				jen.ID("Name").ID("string").Tag(jsonTag("name")),
				jen.ID("ClientID").ID("string").Tag(jsonTag("-")),
				jen.ID("ClientSecret").ID("string").Tag(jsonTag("-")),
				jen.ID("RedirectURI").ID("string").Tag(jsonTag("redirect_uri")),
				jen.ID("BelongsToUser").ID("uint64").Tag(jsonTag("-")),
				jen.ID("Scopes").Index().ID("string").Tag(jsonTag("scopes")),
			),
			jen.Line(),
			jen.Comment("OAuth2ClientUpdateInput is a struct for use when updating OAuth2 clients"),
			jen.ID("OAuth2ClientUpdateInput").Struct(
				jen.ID("RedirectURI").ID("string").Tag(jsonTag("redirect_uri")),
				jen.ID("Scopes").Index().ID("string").Tag(jsonTag("scopes")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("gopkg.in/oauth2.v3", "ClientInfo").Op("=").Parens(jen.Op("*").ID("OAuth2Client")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetID returns the client ID. NOTE: I believe this is implemented for the above interface spec (oauth2.ClientInfo)"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("OAuth2Client")).ID("GetID").Params().Params(jen.ID("string")).Block(
			jen.Return().ID("c").Dot("ClientID"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetSecret returns the ClientSecret"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("OAuth2Client")).ID("GetSecret").Params().Params(jen.ID("string")).Block(
			jen.Return().ID("c").Dot("ClientSecret"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetDomain returns the client's domain"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("OAuth2Client")).ID("GetDomain").Params().Params(jen.ID("string")).Block(
			jen.Return().ID("c").Dot("RedirectURI"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserID returns the client's UserID"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("OAuth2Client")).ID("GetUserID").Params().Params(jen.ID("string")).Block(
			jen.Return().Qual("strconv", "FormatUint").Call(jen.ID("c").Dot("BelongsToUser"), jen.Lit(10)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("HasScope returns whether or not the provided scope is included in the scope list"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("OAuth2Client")).ID("HasScope").Params(jen.ID("scope").ID("string")).Params(jen.ID("found").ID("bool")).Block(
			jen.ID("scope").Op("=").Qual("strings", "TrimSpace").Call(jen.ID("scope")),
			jen.If(jen.ID("scope").Op("==").Lit("")).Block(
				jen.Return().ID("false"),
			),
			jen.If(jen.ID("c").Op("!=").ID("nil").Op("&&").ID("c").Dot("Scopes").Op("!=").ID("nil")).Block(
				jen.For(jen.List(jen.ID("_"), jen.ID("s")).Op(":=").Range().ID("c").Dot("Scopes")).Block(
					jen.If(jen.Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("s"))).Op("==").Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("scope"))).Op("||").Qual("strings", "TrimSpace").Call(jen.ID("s")).Op("==").Lit("*")).Block(
						jen.Return().ID("true"),
					),
				),
			),
			jen.Return().ID("false"),
		),
		jen.Line(),
	)
	return ret
}
