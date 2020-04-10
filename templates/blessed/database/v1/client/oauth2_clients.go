package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(),
			"OAuth2ClientDataManager",
		).Equals().Parens(jen.PointerTo().ID("Client")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Client gets an OAuth2 client from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetOAuth2Client").Params(utils.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			utils.StartSpan(proj, true, "GetOAuth2Client"),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("clientID")),
			jen.Line(),
			jen.ID("logger").Assign().ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("client_id").MapAssign().ID("clientID"), jen.Lit("user_id").MapAssign().ID("userID"),
			)),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("GetOAuth2Client called")),
			jen.Line(),
			jen.List(jen.ID("client"), jen.Err()).Assign().ID("c").Dot("querier").Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("clientID"), jen.ID("userID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching oauth2 client from the querier")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2ClientByClientID fetches any OAuth2 client by client ID, regardless of ownership."),
		jen.Line(),
		jen.Comment("This is used by authenticating middleware to fetch client information it needs to validate."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetOAuth2ClientByClientID").Params(utils.CtxParam(), jen.ID("clientID").String()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetOAuth2ClientByClientID")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientIDToSpan").Call(jen.ID("span"), jen.ID("clientID")),
			jen.ID("logger").Assign().ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("oauth2client_client_id"), jen.ID("clientID")),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("GetOAuth2ClientByClientID called")),
			jen.Line(),
			jen.List(jen.ID("client"), jen.Err()).Assign().ID("c").Dot("querier").Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("clientID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching oauth2 client from the querier")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2ClientCount gets the count of OAuth2 clients that match the current filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetAllOAuth2ClientCount").Params(utils.CtxParam()).Params(jen.Uint64(), jen.Error()).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetAllOAuth2ClientCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetAllOAuth2ClientCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllOAuth2ClientCount").Call(utils.CtxVar()),
		),
		jen.Line(),
	)

	//ret.Add(
	//	jen.Comment("GetAllOAuth2Clients returns all OAuth2 clients, irrespective of ownership."),
	//	jen.Line(),
	//	jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetAllOAuth2Clients").Params(utils.CtxParam()).Params(
	//		jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
	//		jen.Error(),
	//	).Block(
	//		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetAllOAuth2Clients")),
	//		jen.Defer().ID("span").Dot("End").Call(),
	//		jen.Line(),
	//		jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetAllOAuth2Clients called")),
	//		jen.Line(),
	//		jen.Return().ID("c").Dot("querier").Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
	//	),
	//	jen.Line(),
	//)

	ret.Add(
		jen.Comment("GetOAuth2Clients gets a list of OAuth2 clients"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetOAuth2Clients").Params(
			utils.CtxParam(),
			jen.ID("userID").Uint64(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList"), jen.Error()).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetOAuth2Clients")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID(utils.FilterVarName)),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetOAuth2Clients called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.ID("userID"), jen.ID(utils.FilterVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateOAuth2Client creates an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("CreateOAuth2Client").Params(utils.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("CreateOAuth2Client")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Assign().ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("client_id").MapAssign().ID("input").Dot("ClientID"),
				jen.Lit("belongs_to_user").MapAssign().ID("input").Dot("BelongsToUser")),
			),
			jen.Line(),
			jen.List(jen.ID("client"), jen.Err()).Assign().ID("c").Dot("querier").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("input")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("WithError").Call(jen.Err()).Dot("Debug").Call(jen.Lit("error writing oauth2 client to the querier")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("new oauth2 client created successfully")),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateOAuth2Client updates a OAuth2 client. Note that this function expects the input's"),
		jen.Line(),
		jen.Comment("ID field to be valid."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("UpdateOAuth2Client").Params(utils.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.Error()).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("UpdateOAuth2Client")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("UpdateOAuth2Client").Call(utils.CtxVar(), jen.ID("updated")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("ArchiveOAuth2Client").Params(utils.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.Error()).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("ArchiveOAuth2Client")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("clientID")),
			jen.Line(),
			jen.ID("logger").Assign().ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("client_id").MapAssign().ID("clientID"),
				jen.Lit("belongs_to_user").MapAssign().ID("userID"),
			)),
			jen.Line(),
			jen.Err().Assign().ID("c").Dot("querier").Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("clientID"), jen.ID("userID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("WithError").Call(jen.Err()).Dot("Debug").Call(jen.Lit("error deleting oauth2 client to the querier")),
				jen.Return().Err(),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("removed oauth2 client successfully")),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)
	return ret
}
