package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("OAuth2ClientKey is a ContextKey for use with contexts involving OAuth2 clients"),
			jen.ID("OAuth2ClientKey").ID("ContextKey").Equals().Lit("oauth2_client"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("OAuth2ClientDataManager handles OAuth2 clients"),
			jen.ID("OAuth2ClientDataManager").Interface(
				jen.ID("GetOAuth2Client").Params(constants.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.PointerTo().ID("OAuth2Client"), jen.Error()),
				jen.ID("GetOAuth2ClientByClientID").Params(constants.CtxParam(), jen.ID("clientID").String()).Params(jen.PointerTo().ID("OAuth2Client"), jen.Error()),
				jen.ID("GetAllOAuth2ClientCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()),
				jen.ID("GetOAuth2Clients").Params(constants.CtxParam(), jen.ID("userID").Uint64(), utils.QueryFilterParam(nil)).Params(jen.PointerTo().ID("OAuth2ClientList"), jen.Error()),
				jen.ID("CreateOAuth2Client").Params(constants.CtxParam(), jen.ID("input").PointerTo().ID("OAuth2ClientCreationInput")).Params(jen.PointerTo().ID("OAuth2Client"), jen.Error()),
				jen.ID("UpdateOAuth2Client").Params(constants.CtxParam(), jen.ID("updated").PointerTo().ID("OAuth2Client")).Params(jen.Error()),
				jen.ID("ArchiveOAuth2Client").Params(constants.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.Error()),
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
				jen.ID("ExtractOAuth2ClientFromRequest").Params(constants.CtxParam(), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().ID("OAuth2Client"), jen.Error()),
				jen.Line(),
				jen.Comment("wrappers for our implementation library"),
				jen.ID("HandleAuthorizeRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()),
				jen.ID("HandleTokenRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("OAuth2Client represents a user-authorized API client"),
			jen.ID("OAuth2Client").Struct(
				jen.ID("ID").Uint64().Tag(jsonTag("id")),
				jen.ID("Name").String().Tag(jsonTag("name")),
				jen.ID("ClientID").String().Tag(jsonTag("client_id")),
				jen.ID("ClientSecret").String().Tag(jsonTag("client_secret")),
				jen.ID("RedirectURI").String().Tag(jsonTag("redirect_uri")),
				jen.ID("Scopes").Index().String().Tag(jsonTag("scopes")),
				jen.ID("ImplicitAllowed").Bool().Tag(jsonTag("implicit_allowed")),
				jen.ID("CreatedOn").Uint64().Tag(jsonTag("created_on")),
				jen.ID("UpdatedOn").PointerTo().Uint64().Tag(jsonTag("updated_on")),
				jen.ID("ArchivedOn").PointerTo().Uint64().Tag(jsonTag("archived_on")),
				jen.ID(constants.UserOwnershipFieldName).Uint64().Tag(jsonTag("belongs_to_user")),
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
				jen.ID("Name").String().Tag(jsonTag("name")),
				jen.ID("ClientID").String().Tag(jsonTag("-")),
				jen.ID("ClientSecret").String().Tag(jsonTag("-")),
				jen.ID("RedirectURI").String().Tag(jsonTag("redirect_uri")),
				jen.ID(constants.UserOwnershipFieldName).Uint64().Tag(jsonTag("-")),
				jen.ID("Scopes").Index().String().Tag(jsonTag("scopes")),
			),
			jen.Line(),
			jen.Comment("OAuth2ClientUpdateInput is a struct for use when updating OAuth2 clients"),
			jen.ID("OAuth2ClientUpdateInput").Struct(
				jen.ID("RedirectURI").String().Tag(jsonTag("redirect_uri")),
				jen.ID("Scopes").Index().String().Tag(jsonTag("scopes")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Underscore().Qual("gopkg.in/oauth2.v3", "ClientInfo").Equals().Parens(jen.PointerTo().ID("OAuth2Client")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetID returns the client ID. NOTE: I believe this is implemented for the above interface spec (oauth2.ClientInfo)"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("OAuth2Client")).ID("GetID").Params().Params(jen.String()).Block(
			jen.Return().ID("c").Dot("ClientID"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetSecret returns the ClientSecret"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("OAuth2Client")).ID("GetSecret").Params().Params(jen.String()).Block(
			jen.Return().ID("c").Dot("ClientSecret"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetDomain returns the client's domain"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("OAuth2Client")).ID("GetDomain").Params().Params(jen.String()).Block(
			jen.Return().ID("c").Dot("RedirectURI"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserID returns the client's UserID"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("OAuth2Client")).ID("GetUserID").Params().Params(jen.String()).Block(
			jen.Return().Qual("strconv", "FormatUint").Call(jen.ID("c").Dot(constants.UserOwnershipFieldName), jen.Lit(10)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("HasScope returns whether or not the provided scope is included in the scope list"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("OAuth2Client")).ID("HasScope").Params(jen.ID("scope").String()).Params(jen.ID("found").Bool()).Block(
			jen.ID("scope").Equals().Qual("strings", "TrimSpace").Call(jen.ID("scope")),
			jen.If(jen.ID("scope").IsEqualTo().EmptyString()).Block(
				jen.Return().False(),
			),
			jen.If(jen.ID("c").DoesNotEqual().ID("nil").And().ID("c").Dot("Scopes").DoesNotEqual().ID("nil")).Block(
				jen.For(jen.List(jen.Underscore(), jen.ID("s")).Assign().Range().ID("c").Dot("Scopes")).Block(
					jen.If(jen.Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("s"))).IsEqualTo().Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("scope"))).Or().Qual("strings", "TrimSpace").Call(jen.ID("s")).IsEqualTo().Lit("*")).Block(
						jen.Return().True(),
					),
				),
			),
			jen.Return().False(),
		),
		jen.Line(),
	)

	return ret
}
