package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Underscore().Qual(proj.TypesPackage(),
			"OAuth2ClientDataManager",
		).Equals().Parens(jen.PointerTo().ID("Client")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildGetOAuth2Client(proj)...)
	code.Add(buildGetOAuth2ClientByClientID(proj)...)
	code.Add(buildGetAllOAuth2ClientCount(proj)...)
	code.Add(buildGetOAuth2ClientsForUser(proj)...)
	code.Add(buildCreateOAuth2Client(proj)...)
	code.Add(buildUpdateOAuth2Client(proj)...)
	code.Add(buildArchiveOAuth2Client(proj)...)

	return code
}

func buildGetOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetOAuth2Client gets an OAuth2 client from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetOAuth2Client").Params(constants.CtxParam(), jen.List(jen.ID("clientID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client"), jen.Error()).Body(
			utils.StartSpan(proj, true, "GetOAuth2Client"),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.Qual(proj.InternalTracingPackage(), "AttachOAuth2ClientDatabaseIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("clientID")),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("client_id").MapAssign().ID("clientID"), jen.Lit("user_id").MapAssign().ID(constants.UserIDVarName),
			)),
			jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("GetOAuth2Client called")),
			jen.Line(),
			jen.List(jen.ID("client"), jen.Err()).Assign().ID("c").Dot("querier").Dot("GetOAuth2Client").Call(constants.CtxVar(), jen.ID("clientID"), jen.ID(constants.UserIDVarName)),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching oauth2 client from the querier")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildGetOAuth2ClientByClientID(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetOAuth2ClientByClientID fetches any OAuth2 client by client ID, regardless of ownership."),
		jen.Line(),
		jen.Comment("This is used by authenticating middleware to fetch client information it needs to validate."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetOAuth2ClientByClientID").Params(constants.CtxParam(), jen.ID("clientID").String()).Params(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client"), jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetOAuth2ClientByClientID")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachOAuth2ClientIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("clientID")),
			jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("oauth2client_client_id"), jen.ID("clientID")),
			jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("GetOAuth2ClientByClientID called")),
			jen.Line(),
			jen.List(jen.ID("client"), jen.Err()).Assign().ID("c").Dot("querier").Dot("GetOAuth2ClientByClientID").Call(constants.CtxVar(), jen.ID("clientID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching oauth2 client from the querier")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllOAuth2ClientCount(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetAllOAuth2ClientCount gets the count of OAuth2 clients that match the current filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetAllOAuth2ClientCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetAllOAuth2ClientCount")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("GetAllOAuth2ClientCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllOAuth2ClientCount").Call(constants.CtxVar()),
		),
		jen.Line(),
	}

	return lines
}

func buildGetOAuth2ClientsForUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetOAuth2ClientsForUser gets a list of OAuth2 clients."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetOAuth2ClientsForUser").Params(
			constants.CtxParam(),
			constants.UserIDParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.TypesPackage(), "QueryFilter"),
		).Params(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2ClientList"), jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetOAuth2ClientsForUser")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.FilterVarName)),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)).Dot("Debug").Call(jen.Lit("GetOAuth2ClientsForUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetOAuth2ClientsForUser").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName), jen.ID(constants.FilterVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateOAuth2Client creates an OAuth2 client."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("CreateOAuth2Client").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "OAuth2ClientCreationInput")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client"), jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("CreateOAuth2Client")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("oauth2client_client_id").MapAssign().ID("input").Dot("ClientID"),
				jen.Lit("belongs_to_user").MapAssign().ID("input").Dot(constants.UserOwnershipFieldName)),
			),
			jen.Line(),
			jen.List(jen.ID("client"), jen.Err()).Assign().ID("c").Dot("querier").Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID("input")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("WithError").Call(jen.Err()).Dot("Debug").Call(jen.Lit("error writing oauth2 client to the querier")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("client_id"), jen.ID("client").Dot("ID")).Dot("Debug").Call(jen.Lit("new oauth2 client created successfully")),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("UpdateOAuth2Client updates a OAuth2 client. Note that this function expects the input's"),
		jen.Line(),
		jen.Comment("ID field to be valid."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("UpdateOAuth2Client").Params(constants.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.TypesPackage(), "OAuth2Client")).Params(jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("UpdateOAuth2Client")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("UpdateOAuth2Client").Call(constants.CtxVar(), jen.ID("updated")),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("ArchiveOAuth2Client").Params(constants.CtxParam(), jen.List(jen.ID("clientID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("ArchiveOAuth2Client")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.Qual(proj.InternalTracingPackage(), "AttachOAuth2ClientDatabaseIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("clientID")),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("client_id").MapAssign().ID("clientID"),
				jen.Lit("belongs_to_user").MapAssign().ID(constants.UserIDVarName),
			)),
			jen.Line(),
			jen.Err().Assign().ID("c").Dot("querier").Dot("ArchiveOAuth2Client").Call(constants.CtxVar(), jen.ID("clientID"), jen.ID(constants.UserIDVarName)),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("WithError").Call(jen.Err()).Dot("Debug").Call(jen.Lit("error deleting oauth2 client to the querier")),
				jen.Return().Err(),
			),
			jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("removed oauth2 client successfully")),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	}

	return lines
}
