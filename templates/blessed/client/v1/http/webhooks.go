package client

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(proj, ret)
	ret.Add(jen.Const().Defs(
		jen.ID("webhooksBasePath").Op("=").Lit("webhooks"),
	))

	ret.Add(buildBuildGetWebhookRequest()...)
	ret.Add(buildGetWebhook(proj)...)
	ret.Add(buildBuildGetWebhooksRequest(proj)...)
	ret.Add(buildGetWebhooks(proj)...)
	ret.Add(buildBuildCreateWebhookRequest(proj)...)
	ret.Add(buildCreateWebhook(proj)...)
	ret.Add(buildBuildUpdateWebhookRequest(proj)...)
	ret.Add(buildUpdateWebhook(proj)...)
	ret.Add(buildBuildArchiveWebhookRequest()...)
	ret.Add(buildArchiveWebhook()...)

	return ret
}

func buildBuildGetWebhookRequest() []jen.Code {
	funcName := "BuildGetWebhookRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
			jen.ID("nil"),
			jen.ID("webhooksBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID("id"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.ID("nil"),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for fetching a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...,
		),
	}

	return lines
}

func buildGetWebhook(proj *models.Project) []jen.Code {
	funcName := "GetWebhook"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildGetWebhookRequest").Call(
			utils.CtxVar(),
			jen.ID("id"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"),
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
			jen.Op("&").ID("webhook"),
		),
		jen.Return().List(
			jen.ID("webhook"),
			jen.Err(),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s retrieves a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("webhook").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "Webhook"),
			jen.Err().ID("error"),
		).Block(block...),
	}

	return lines
}

func buildBuildGetWebhooksRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetWebhooksRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
			jen.ID(utils.FilterVarName).Dot("ToValues").Call(),
			jen.ID("webhooksBasePath"),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.ID("nil"),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for fetching webhooks", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
	}

	return lines
}

func buildGetWebhooks(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath
	funcName := "GetWebhooks"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildGetWebhooksRequest").Call(
			utils.CtxVar(),
			jen.ID(utils.FilterVarName),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.ID("nil"),
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
			jen.Op("&").ID("webhooks"),
		),
		jen.Return().List(
			jen.ID("webhooks"),
			jen.Err(),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s gets a list of webhooks", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(outPath, "models/v1"), "QueryFilter"),
		).Params(
			jen.ID("webhooks").Op("*").Qual(filepath.Join(outPath, "models/v1"), "WebhookList"),
			jen.Err().ID("error"),
		).Block(block...),
	}

	return lines
}

func buildBuildCreateWebhookRequest(proj *models.Project) []jen.Code {
	funcName := "BuildCreateWebhookRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
			jen.ID("nil"),
			jen.ID("webhooksBasePath"),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for creating a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("body").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "WebhookCreationInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
	}

	return lines
}

func buildCreateWebhook(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath
	funcName := "CreateWebhook"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildCreateWebhookRequest").Call(
			utils.CtxVar(),
			jen.ID("input"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Op("=").ID("c").Dot("executeRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.Op("&").ID("webhook"),
		),
		jen.Return().List(
			jen.ID("webhook"),
			jen.Err(),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s creates a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("input").Op("*").Qual(filepath.Join(outPath, "models/v1"), "WebhookCreationInput"),
		).Params(
			jen.ID("webhook").Op("*").Qual(filepath.Join(outPath, "models/v1"), "Webhook"),
			jen.Err().ID("error"),
		).Block(block...),
	}

	return lines
}

func buildBuildUpdateWebhookRequest(proj *models.Project) []jen.Code {
	funcName := "BuildUpdateWebhookRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
			jen.ID("nil"),
			jen.ID("webhooksBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID("updated").Dot("ID"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPut"),
			jen.ID("uri"),
			jen.ID("updated"),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("updated").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "Webhook"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
	}

	return lines
}

func buildUpdateWebhook(proj *models.Project) []jen.Code {
	funcName := "UpdateWebhook"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildUpdateWebhookRequest").Call(
			utils.CtxVar(),
			jen.ID("updated"),
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
			jen.ID("req"), jen.Op("&").ID("updated"),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s updates a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("updated").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "Webhook"),
		).Params(
			jen.ID("error"),
		).Block(block...),
	}

	return lines
}

func buildBuildArchiveWebhookRequest() []jen.Code {
	funcName := "BuildArchiveWebhookRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
			jen.ID("nil"),
			jen.ID("webhooksBasePath"), jen.Qual("strconv", "FormatUint").Call(
				jen.ID("id"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequest").Call(
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.ID("nil"),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
	}

	return lines
}

func buildArchiveWebhook() []jen.Code {
	funcName := "ArchiveWebhook"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildArchiveWebhookRequest").Call(
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
			jen.ID("nil"),
		),
	)

	lines := []jen.Code{
		jen.Commentf("%s archives a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("error"),
		).Block(block...,
		),
	}

	return lines
}
