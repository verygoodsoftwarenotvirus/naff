package client

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)
	ret.Add(jen.Const().Defs(
		jen.ID("webhooksBasePath").Equals().Lit("webhooks"),
	))

	ret.Add(buildBuildGetWebhookRequest(proj)...)
	ret.Add(buildGetWebhook(proj)...)
	ret.Add(buildBuildGetWebhooksRequest(proj)...)
	ret.Add(buildGetWebhooks(proj)...)
	ret.Add(buildBuildCreateWebhookRequest(proj)...)
	ret.Add(buildCreateWebhook(proj)...)
	ret.Add(buildBuildUpdateWebhookRequest(proj)...)
	ret.Add(buildUpdateWebhook(proj)...)
	ret.Add(buildBuildArchiveWebhookRequest(proj)...)
	ret.Add(buildArchiveWebhook(proj)...)

	return ret
}

func buildBuildGetWebhookRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetWebhookRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.Nil(),
			jen.ID("webhooksBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
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
		jen.Commentf("%s builds an HTTP request for fetching a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ParamPointer().Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...,
		),
	}

	return lines
}

func buildGetWebhook(proj *models.Project) []jen.Code {
	funcName := "GetWebhook"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetWebhookRequest").Call(
			utils.CtxVar(),
			jen.ID("id"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(jen.Nil(),
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
			jen.VarPointer().ID("webhook"),
		),
		jen.Return().List(
			jen.ID("webhook"),
			jen.Err(),
		),
	}

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

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.ID(utils.FilterVarName).Dot("ToValues").Call(),
			jen.ID("webhooksBasePath"),
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
		jen.Commentf("%s builds an HTTP request for fetching webhooks", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter"),
		).Params(
			jen.ParamPointer().Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
	}

	return lines
}

func buildGetWebhooks(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath
	funcName := "GetWebhooks"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetWebhooksRequest").Call(
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
		jen.Err().Equals().ID("c").Dot("retrieve").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.VarPointer().ID("webhooks"),
		),
		jen.Return().List(
			jen.ID("webhooks"),
			jen.Err(),
		),
	}

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

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.Nil(),
			jen.ID("webhooksBasePath"),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for creating a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("body").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "WebhookCreationInput"),
		).Params(
			jen.ParamPointer().Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
	}

	return lines
}

func buildCreateWebhook(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath
	funcName := "CreateWebhook"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildCreateWebhookRequest").Call(
			utils.CtxVar(),
			jen.ID("input"),
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
		jen.Err().Equals().ID("c").Dot("executeRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.VarPointer().ID("webhook"),
		),
		jen.Return().List(
			jen.ID("webhook"),
			jen.Err(),
		),
	}

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

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.Nil(),
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
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("updated").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "Webhook"),
		).Params(
			jen.ParamPointer().Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
	}

	return lines
}

func buildUpdateWebhook(proj *models.Project) []jen.Code {
	funcName := "UpdateWebhook"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildUpdateWebhookRequest").Call(
			utils.CtxVar(),
			jen.ID("updated"),
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
			jen.ID("req"), jen.VarPointer().ID("updated"),
		),
	}

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

func buildBuildArchiveWebhookRequest(proj *models.Project) []jen.Code {
	funcName := "BuildArchiveWebhookRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.Nil(),
			jen.ID("webhooksBasePath"), jen.Qual("strconv", "FormatUint").Call(
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
		jen.Commentf("%s builds an HTTP request for updating a webhook", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ParamPointer().Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
	}

	return lines
}

func buildArchiveWebhook(proj *models.Project) []jen.Code {
	funcName := "ArchiveWebhook"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildArchiveWebhookRequest").Call(
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
