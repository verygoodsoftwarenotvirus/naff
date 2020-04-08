package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(jen.Null())
	ret.Add(jen.Const().Defs(
		jen.ID("oauth2ClientsBasePath").Equals().Lit("oauth2/clients"),
	))

	ret.Add(buildBuildGetOAuth2ClientRequest(proj)...)
	ret.Add(buildGetOAuth2Client(proj)...)
	ret.Add(buildBuildGetOAuth2ClientsRequest(proj)...)
	ret.Add(buildGetOAuth2Clients(proj)...)
	ret.Add(buildBuildCreateOAuth2ClientRequest(proj)...)
	ret.Add(buildCreateOAuth2Client(proj)...)
	ret.Add(buildBuildArchiveOAuth2ClientRequest(proj)...)
	ret.Add(buildArchiveOAuth2Client(proj)...)

	ret.Add()

	return ret
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
			utils.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildGetOAuth2ClientRequest builds an HTTP request for fetching an OAuth2 client"),
		jen.Line(),
		newClientMethod("BuildGetOAuth2ClientRequest").Params(
			utils.CtxParam(),
			jen.ID("id").Uint64(),
		).Params(
			jen.ParamPointer().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildGetOAuth2Client(proj *models.Project) []jen.Code {
	funcName := "GetOAuth2Client"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetOAuth2ClientRequest").Call(
			utils.CtxVar(),
			jen.ID("id"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
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
			utils.CtxVar(),
			jen.ID("req"),
			jen.AddressOf().ID("oauth2Client"),
		),
		jen.Return().List(
			jen.ID("oauth2Client"),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Comment("GetOAuth2Client gets an OAuth2 client"),
		jen.Line(),
		newClientMethod("GetOAuth2Client").Params(
			utils.CtxParam(),
			jen.ID("id").Uint64(),
		).Params(
			jen.ID("oauth2Client").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Err().Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildGetOAuth2ClientsRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetOAuth2ClientsRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.ID(utils.FilterVarName).Dot("ToValues").Call(),
			jen.ID("oauth2ClientsBasePath"),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildGetOAuth2ClientsRequest builds an HTTP request for fetching a list of OAuth2 clients"),
		jen.Line(),
		newClientMethod("BuildGetOAuth2ClientsRequest").Params(
			utils.CtxParam(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.ParamPointer().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildGetOAuth2Clients(proj *models.Project) []jen.Code {
	funcName := "GetOAuth2Clients"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetOAuth2ClientsRequest").Call(
			utils.CtxVar(),
			jen.ID(utils.FilterVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Var().ID("oauth2Clients").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList"),
		jen.Err().Equals().ID("c").Dot("retrieve").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.AddressOf().ID("oauth2Clients"),
		),
		jen.Return().List(
			jen.ID("oauth2Clients"),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Comment("GetOAuth2Clients gets a list of OAuth2 clients"),
		jen.Line(),
		newClientMethod("GetOAuth2Clients").Params(
			utils.CtxParam(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList"),
			jen.Error(),
		).Block(block...),
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
		jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.ID("req").Dot("AddCookie").Call(
			jen.ID("cookie"),
		),
		jen.Line(),
		jen.Return().List(
			jen.ID("req"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildCreateOAuth2ClientRequest builds an HTTP request for creating OAuth2 clients"),
		jen.Line(),
		newClientMethod("BuildCreateOAuth2ClientRequest").Paramsln(
			utils.CtxParam(),
			jen.ID("cookie").ParamPointer().Qual("net/http", "Cookie"),
			jen.ID("body").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput"),
		).Params(jen.ParamPointer().Qual("net/http", "Request"),
			jen.Error()).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildCreateOAuth2Client(proj *models.Project) []jen.Code {
	funcName := "CreateOAuth2Client"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Var().ID("oauth2Client").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
		jen.If(jen.ID("cookie").Op("==").ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("errors", "New").Call(
					jen.Lit("cookie required for request"),
				),
			),
		),
		jen.Line(),
		jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot("BuildCreateOAuth2ClientRequest").Call(
			utils.CtxVar(),
			jen.ID("cookie"),
			jen.ID("input"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.If(
			jen.ID("resErr").Assign().ID("c").Dot("executeUnauthenticatedDataRequest").Call(utils.CtxVar(), jen.ID("req"), jen.AddressOf().ID("oauth2Client")),
			jen.ID("resErr").DoesNotEqual().ID("nil"),
		).Block(
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
		jen.Comment("in order to receive a valid response"),
		jen.Line(),
		newClientMethod("CreateOAuth2Client").Paramsln(
			utils.CtxParam(),
			jen.ID("cookie").ParamPointer().Qual("net/http", "Cookie"),
			jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error(),
		).Block(block...,
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
			utils.CtxVar(),
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildArchiveOAuth2ClientRequest builds an HTTP request for archiving an oauth2 client"),
		jen.Line(),
		newClientMethod("BuildArchiveOAuth2ClientRequest").Params(utils.CtxParam(), jen.ID("id").Uint64()).Params(
			jen.ParamPointer().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildArchiveOAuth2Client(proj *models.Project) []jen.Code {
	funcName := "ArchiveOAuth2Client"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildArchiveOAuth2ClientRequest").Call(
			utils.CtxVar(),
			jen.ID("id"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("building request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("executeRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client"),
		jen.Line(),
		newClientMethod("ArchiveOAuth2Client").Params(utils.CtxParam(), jen.ID("id").Uint64()).Params(jen.Error()).Block(block...),
		jen.Line(),
	}

	return lines
}
