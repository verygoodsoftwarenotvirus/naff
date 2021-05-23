package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Null())
	code.Add(
		jen.Const().Defs(
			jen.ID("oauth2ClientsBasePath").Equals().Lit("oauth2/clients"),
		))

	code.Add(buildBuildGetOAuth2ClientRequest(proj)...)
	code.Add(buildGetOAuth2Client(proj)...)
	code.Add(buildBuildGetOAuth2ClientsRequest(proj)...)
	code.Add(buildGetOAuth2Clients(proj)...)
	code.Add(buildBuildCreateOAuth2ClientRequest(proj)...)
	code.Add(buildCreateOAuth2Client(proj)...)
	code.Add(buildBuildArchiveOAuth2ClientRequest(proj)...)
	code.Add(buildArchiveOAuth2Client(proj)...)

	code.Add()

	return code
}

func buildBuildGetOAuth2ClientRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetOAuth2ClientRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.Nil(),
			jen.ID("oauth2ClientsBasePath"), jen.Qual("strconv", "FormatUint").Call(
				jen.ID("id"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildGetOAuth2ClientRequest builds an HTTP request for fetching an OAuth2 client."),
		jen.Line(),
		newClientMethod("BuildGetOAuth2ClientRequest").Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildGetOAuth2Client(proj *models.Project) []jen.Code {
	funcName := "GetOAuth2Client"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetOAuth2ClientRequest").Call(
			constants.CtxVar(),
			jen.ID("id"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Equals().ID("c").Dot("retrieve").Call(
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.AddressOf().ID("oauth2Client"),
		),
		jen.Return().List(
			jen.ID("oauth2Client"),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Comment("GetOAuth2Client gets an OAuth2 client."),
		jen.Line(),
		newClientMethod("GetOAuth2Client").Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
		).Params(
			jen.ID("oauth2Client").PointerTo().Qual(proj.TypesPackage(), "OAuth2Client"),
			jen.Err().Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildGetOAuth2ClientsRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetOAuth2ClientsRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.ID(constants.FilterVarName).Dot("ToValues").Call(),
			jen.ID("oauth2ClientsBasePath"),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildGetOAuth2ClientsRequest builds an HTTP request for fetching a list of OAuth2 clients."),
		jen.Line(),
		newClientMethod("BuildGetOAuth2ClientsRequest").Params(
			constants.CtxParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.TypesPackage(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildGetOAuth2Clients(proj *models.Project) []jen.Code {
	funcName := "GetOAuth2Clients"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetOAuth2ClientsRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.FilterVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Var().ID("oauth2Clients").PointerTo().Qual(proj.TypesPackage(), "OAuth2ClientList"),
		jen.Err().Equals().ID("c").Dot("retrieve").Call(
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.AddressOf().ID("oauth2Clients"),
		),
		jen.Return().List(
			jen.ID("oauth2Clients"),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Comment("GetOAuth2Clients gets a list of OAuth2 clients."),
		jen.Line(),
		newClientMethod("GetOAuth2Clients").Params(
			constants.CtxParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.TypesPackage(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2ClientList"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildCreateOAuth2ClientRequest(proj *models.Project) []jen.Code {
	funcName := "BuildCreateOAuth2ClientRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.Lit("oauth2"),
			jen.Lit("client"),
		),
		jen.Line(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("c").Dot("buildDataRequest").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.ID(constants.RequestVarName).Dot("AddCookie").Call(
			jen.ID("cookie"),
		),
		jen.Line(),
		jen.Return().List(
			jen.ID(constants.RequestVarName),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildCreateOAuth2ClientRequest builds an HTTP request for creating OAuth2 clients."),
		jen.Line(),
		newClientMethod("BuildCreateOAuth2ClientRequest").Paramsln(
			constants.CtxParam(),
			jen.ID("cookie").PointerTo().Qual("net/http", "Cookie"),
			jen.ID("body").PointerTo().Qual(proj.TypesPackage(), "OAuth2ClientCreationInput"),
		).Params(jen.PointerTo().Qual("net/http", "Request"),
			jen.Error()).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildCreateOAuth2Client(proj *models.Project) []jen.Code {
	funcName := "CreateOAuth2Client"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Var().ID("oauth2Client").PointerTo().Qual(proj.TypesPackage(), "OAuth2Client"),
		jen.If(jen.ID("cookie").IsEqualTo().ID("nil")).Body(
			jen.Return().List(
				jen.Nil(),
				utils.Error("cookie required for request"),
			),
		),
		jen.Line(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("c").Dot("BuildCreateOAuth2ClientRequest").Call(
			constants.CtxVar(),
			jen.ID("cookie"),
			jen.ID("input"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.If(
			jen.ID("resErr").Assign().ID("c").Dot("executeUnauthenticatedDataRequest").Call(constants.CtxVar(), jen.ID(constants.RequestVarName), jen.AddressOf().ID("oauth2Client")),
			jen.ID("resErr").DoesNotEqual().ID("nil"),
		).Body(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("loading response from server: %w"),
					jen.ID("resErr"),
				),
			),
		),
		jen.Line(),
		jen.Return().List(
			jen.ID("oauth2Client"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("CreateOAuth2Client creates an OAuth2 client. Note that cookie must not be nil"),
		jen.Line(),
		jen.Comment("in order to receive a valid response."),
		jen.Line(),
		newClientMethod("CreateOAuth2Client").Paramsln(
			constants.CtxParam(),
			jen.ID("cookie").PointerTo().Qual("net/http", "Cookie"),
			jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "OAuth2ClientCreationInput"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client"),
			jen.Error(),
		).Body(block...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveOAuth2ClientRequest(proj *models.Project) []jen.Code {
	funcName := "BuildArchiveOAuth2ClientRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.Nil(),
			jen.ID("oauth2ClientsBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID("id"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildArchiveOAuth2ClientRequest builds an HTTP request for archiving an oauth2 client."),
		jen.Line(),
		newClientMethod("BuildArchiveOAuth2ClientRequest").Params(constants.CtxParam(), jen.ID("id").Uint64()).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildArchiveOAuth2Client(proj *models.Project) []jen.Code {
	funcName := "ArchiveOAuth2Client"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildArchiveOAuth2ClientRequest").Call(
			constants.CtxVar(),
			jen.ID("id"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("building request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("executeRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client."),
		jen.Line(),
		newClientMethod("ArchiveOAuth2Client").Params(constants.CtxParam(), jen.ID("id").Uint64()).Params(jen.Error()).Body(block...),
		jen.Line(),
	}

	return lines
}
