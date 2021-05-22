package requests

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)
	code.Add(jen.Const().Defs(
		jen.ID("webhooksBasePath").Equals().Lit("webhooks"),
	))

	code.Add(buildBuildGetWebhookRequest(proj)...)
	code.Add(buildGetWebhook(proj)...)
	code.Add(buildBuildGetWebhooksRequest(proj)...)
	code.Add(buildGetWebhooks(proj)...)
	code.Add(buildBuildCreateWebhookRequest(proj)...)
	code.Add(buildCreateWebhook(proj)...)
	code.Add(buildBuildUpdateWebhookRequest(proj)...)
	code.Add(buildUpdateWebhook(proj)...)
	code.Add(buildBuildArchiveWebhookRequest(proj)...)
	code.Add(buildArchiveWebhook(proj)...)

	return code
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
			constants.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for fetching a webhook.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...,
		),
	}

	return lines
}

func buildGetWebhook(proj *models.Project) []jen.Code {
	funcName := "GetWebhook"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetWebhookRequest").Call(
			constants.CtxVar(),
			jen.ID("id"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.Nil(),
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
			jen.AddressOf().ID("webhook"),
		),
		jen.Return().List(
			jen.ID("webhook"),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s retrieves a webhook.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
		).Params(
			jen.ID("webhook").PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Err().Error(),
		).Body(block...),
	}

	return lines
}

func buildBuildGetWebhooksRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetWebhooksRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Call(
			jen.ID(constants.FilterVarName).Dot("ToValues").Call(),
			jen.ID("webhooksBasePath"),
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
		jen.Commentf("%s builds an HTTP request for fetching webhooks.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
	}

	return lines
}

func buildGetWebhooks(proj *models.Project) []jen.Code {
	funcName := "GetWebhooks"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetWebhooksRequest").Call(
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
		jen.Err().Equals().ID("c").Dot("retrieve").Call(
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.AddressOf().ID("webhooks"),
		),
		jen.Return().List(
			jen.ID("webhooks"),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s gets a list of webhooks.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.ID("webhooks").PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"),
			jen.Err().Error(),
		).Body(block...),
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
			constants.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for creating a webhook.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("body").PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
	}

	return lines
}

func buildCreateWebhook(proj *models.Project) []jen.Code {
	funcName := "CreateWebhook"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildCreateWebhookRequest").Call(
			constants.CtxVar(),
			jen.ID("input"),
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
		jen.Err().Equals().ID("c").Dot("executeRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.AddressOf().ID("webhook"),
		),
		jen.Return().List(
			jen.ID("webhook"),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s creates a webhook.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput"),
		).Params(
			jen.ID("webhook").PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Err().Error(),
		).Body(block...),
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
			constants.CtxVar(),
			jen.Qual("net/http", "MethodPut"),
			jen.ID("uri"),
			jen.ID("updated"),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating a webhook.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
	}

	return lines
}

func buildUpdateWebhook(proj *models.Project) []jen.Code {
	funcName := "UpdateWebhook"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildUpdateWebhookRequest").Call(
			constants.CtxVar(),
			jen.ID("updated"),
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
			jen.ID(constants.RequestVarName), jen.AddressOf().ID("updated"),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s updates a webhook.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
		).Params(
			jen.Error(),
		).Body(block...),
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
			constants.CtxVar(),
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating a webhook.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
	}

	return lines
}

func buildArchiveWebhook(proj *models.Project) []jen.Code {
	funcName := "ArchiveWebhook"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildArchiveWebhookRequest").Call(
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
		jen.Commentf("%s archives a webhook.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
		).Params(
			jen.Error(),
		).Body(block...,
		),
	}

	return lines
}
