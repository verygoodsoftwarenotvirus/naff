package client

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(proj, ret)

	ret.Add(jen.Null())
	ret.Add(jen.Const().Defs(
		jen.ID("oauth2ClientsBasePath").Op("=").Lit("oauth2/clients"),
	))

	ret.Add(buildBuildGetOAuth2ClientRequest(proj)...)
	ret.Add(buildGetOAuth2Client(proj)...)
	ret.Add(buildBuildGetOAuth2ClientsRequest(proj)...)
	ret.Add(buildGetOAuth2Clients(proj)...)
	ret.Add(buildBuildCreateOAuth2ClientRequest(proj)...)
	ret.Add(buildCreateOAuth2Client(proj)...)
	ret.Add(buildBuildArchiveOAuth2ClientRequest()...)
	ret.Add(buildArchiveOAuth2Client()...)

	ret.Add()

	return ret
}

func buildBuildGetOAuth2ClientRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetOAuth2ClientRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
			jen.Nil(),
			jen.ID("oauth2ClientsBasePath"), jen.Qual("strconv", "FormatUint").Call(
				jen.ID("id"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	)

	lines := []jen.Code{
		jen.Comment("BuildGetOAuth2ClientRequest builds an HTTP request for fetching an OAuth2 client"),
		jen.Line(),
		newClientMethod("BuildGetOAuth2ClientRequest").Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildGetOAuth2Client(proj *models.Project) []jen.Code {
	funcName := "GetOAuth2Client"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildGetOAuth2ClientRequest").Call(
			utils.CtxVar(),
			jen.ID("id"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Op("=").ID("c").Dot("retrieve").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.Op("&").ID("oauth2Client"),
		),
		jen.Return().List(
			jen.ID("oauth2Client"),
			jen.Err(),
		),
	)

	lines := []jen.Code{
		jen.Comment("GetOAuth2Client gets an OAuth2 client"),
		jen.Line(),
		newClientMethod("GetOAuth2Client").Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("oauth2Client").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "OAuth2Client"),
			jen.Err().ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildGetOAuth2ClientsRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetOAuth2ClientsRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
			jen.ID(utils.FilterVarName).Dot("ToValues").Call(),
			jen.ID("oauth2ClientsBasePath"),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	)

	lines := []jen.Code{
		jen.Comment("BuildGetOAuth2ClientsRequest builds an HTTP request for fetching a list of OAuth2 clients"),
		jen.Line(),
		newClientMethod("BuildGetOAuth2ClientsRequest").Params(
			utils.CtxParam(),
			jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildGetOAuth2Clients(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath
	funcName := "GetOAuth2Clients"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildGetOAuth2ClientsRequest").Call(
			utils.CtxVar(),
			jen.ID(utils.FilterVarName),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Var().ID("oauth2Clients").Op("*").Qual(filepath.Join(outPath, "models/v1"), "OAuth2ClientList"),
		jen.Err().Op("=").ID("c").Dot("retrieve").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.Op("&").ID("oauth2Clients"),
		),
		jen.Return().List(
			jen.ID("oauth2Clients"),
			jen.Err(),
		),
	)

	lines := []jen.Code{
		jen.Comment("GetOAuth2Clients gets a list of OAuth2 clients"),
		jen.Line(),
		newClientMethod("GetOAuth2Clients").Params(
			utils.CtxParam(),
			jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(outPath, "models/v1"), "QueryFilter"),
		).Params(
			jen.Op("*").Qual(filepath.Join(outPath, "models/v1"), "OAuth2ClientList"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildCreateOAuth2ClientRequest(proj *models.Project) []jen.Code {
	funcName := "BuildCreateOAuth2ClientRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.Lit("oauth2"),
			jen.Lit("client"),
		),
		jen.Line(),
		jen.List(jen.ID("req"), jen.Err()).Op(":=").ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
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
	)

	lines := []jen.Code{
		jen.Comment("BuildCreateOAuth2ClientRequest builds an HTTP request for creating OAuth2 clients"),
		jen.Line(),
		newClientMethod("BuildCreateOAuth2ClientRequest").Paramsln(
			utils.CtxParam(),
			jen.ID("cookie").Op("*").Qual("net/http", "Cookie"),
			jen.ID("body").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "OAuth2ClientCreationInput"),
		).Params(jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error")).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildCreateOAuth2Client(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath
	funcName := "CreateOAuth2Client"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.Var().ID("oauth2Client").Op("*").Qual(filepath.Join(outPath, "models/v1"), "OAuth2Client"),
		jen.If(jen.ID("cookie").Op("==").ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("errors", "New").Call(
					jen.Lit("cookie required for request"),
				),
			),
		),
		jen.Line(),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildCreateOAuth2ClientRequest").Call(
			utils.CtxVar(),
			jen.ID("cookie"),
			jen.ID("input"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.List(
			jen.ID("res"),
			jen.Err(),
		).Op(":=").ID("c").Dot("executeRawRequest").Call(
			utils.CtxVar(),
			jen.ID("c").Dot("plainClient"),
			jen.ID("req"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("executing request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.If(jen.ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.ID("ErrNotFound"),
			),
		),
		jen.Line(),
		jen.If(jen.ID("resErr").Op(":=").ID("unmarshalBody").Call(
			jen.ID("res"),
			jen.Op("&").ID("oauth2Client"),
		),
			jen.ID("resErr").Op("!=").ID("nil"),
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
	)

	lines := []jen.Code{
		jen.Comment("CreateOAuth2Client creates an OAuth2 client. Note that cookie must not be nil"),
		jen.Line(),
		jen.Comment("in order to receive a valid response"),
		jen.Line(),
		newClientMethod("CreateOAuth2Client").Paramsln(
			utils.CtxParam(),
			jen.ID("cookie").Op("*").Qual("net/http", "Cookie"),
			jen.ID("input").Op("*").Qual(filepath.Join(outPath, "models/v1"), "OAuth2ClientCreationInput"),
		).Params(
			jen.Op("*").Qual(filepath.Join(outPath, "models/v1"), "OAuth2Client"),
			jen.ID("error"),
		).Block(block...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveOAuth2ClientRequest() []jen.Code {
	funcName := "BuildArchiveOAuth2ClientRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
			jen.Nil(),
			jen.ID("oauth2ClientsBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID("id"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequest").Call(
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.Nil(),
		),
	)

	lines := []jen.Code{
		jen.Comment("BuildArchiveOAuth2ClientRequest builds an HTTP request for archiving an oauth2 client"),
		jen.Line(),
		newClientMethod("BuildArchiveOAuth2ClientRequest").Params(utils.CtxParam(), jen.ID("id").ID("uint64")).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildArchiveOAuth2Client() []jen.Code {
	funcName := "ArchiveOAuth2Client"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildArchiveOAuth2ClientRequest").Call(
			utils.CtxVar(),
			jen.ID("id"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
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
	)

	lines := []jen.Code{
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client"),
		jen.Line(),
		newClientMethod("ArchiveOAuth2Client").Params(utils.CtxParam(), jen.ID("id").ID("uint64")).Params(jen.ID("error")).Block(block...),
		jen.Line(),
	}

	return lines
}
