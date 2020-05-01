package fake

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)
	utils.AddImports(proj, ret)

	ret.Add(buildBuildFakeWebhook(proj)...)
	ret.Add(buildBuildFakeWebhookList(proj)...)
	ret.Add(buildBuildFakeWebhookUpdateInputFromWebhook(proj)...)
	ret.Add(buildBuildFakeWebhookCreationInput(proj)...)
	ret.Add(buildBuildFakeWebhookCreationInputFromWebhook(proj)...)

	return ret
}

func buildBuildFakeWebhook(proj *models.Project) []jen.Code {
	funcName := "BuildFakeWebhook"
	typeName := "Webhook"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.AddressOf().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
					jen.ID("ContentType").MapAssign().Add(utils.FakeContentTypeFunc()),
					jen.ID("URL").MapAssign().Add(utils.FakeURLFunc()),
					jen.ID("Method").MapAssign().Add(utils.FakeHTTPMethodFunc()),
					jen.ID("Events").MapAssign().Index().String().Values(utils.FakeStringFunc()),
					jen.ID("DataTypes").MapAssign().Index().String().Values(utils.FakeStringFunc()),
					jen.ID("Topics").MapAssign().Index().String().Values(utils.FakeStringFunc()),
					jen.ID("CreatedOn").MapAssign().Add(utils.FakeUnixTimeFunc()),
					jen.ID("ArchivedOn").MapAssign().Nil(),
					jen.ID(constants.UserOwnershipFieldName).MapAssign().Add(utils.FakeUint64Func()),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeWebhookList(proj *models.Project) []jen.Code {
	funcName := "BuildFakeWebhookList"
	typeName := "WebhookList"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.ID(utils.BuildFakeVarName("Webhook1")).Assign().ID("BuildFakeWebhook").Call(),
			jen.ID(utils.BuildFakeVarName("Webhook2")).Assign().ID("BuildFakeWebhook").Call(),
			jen.ID(utils.BuildFakeVarName("Webhook3")).Assign().ID("BuildFakeWebhook").Call(),
			jen.Return(
				jen.AddressOf().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().One(),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().Lit(3),
					),
					jen.ID("Webhooks").MapAssign().Index().Qual(proj.ModelsV1Package(), "Webhook").Valuesln(
						jen.PointerTo().ID("exampleWebhook1"),
						jen.PointerTo().ID("exampleWebhook2"),
						jen.PointerTo().ID("exampleWebhook3"),
					),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeWebhookUpdateInputFromWebhook(proj *models.Project) []jen.Code {
	funcName := "BuildFakeWebhookUpdateInputFromWebhook"
	typeName := "WebhookUpdateInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("webhook").PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.AddressOf().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("Name").MapAssign().ID("webhook").Dot("Name"),
					jen.ID("ContentType").MapAssign().ID("webhook").Dot("ContentType"),
					jen.ID("URL").MapAssign().ID("webhook").Dot("URL"),
					jen.ID("Method").MapAssign().ID("webhook").Dot("Method"),
					jen.ID("Events").MapAssign().ID("webhook").Dot("Events"),
					jen.ID("DataTypes").MapAssign().ID("webhook").Dot("DataTypes"),
					jen.ID("Topics").MapAssign().ID("webhook").Dot("Topics"),
					jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("webhook").Dot(constants.UserOwnershipFieldName),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeWebhookCreationInput(proj *models.Project) []jen.Code {
	funcName := "BuildFakeWebhookCreationInput"
	typeName := "WebhookCreationInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.ID("webhook").Assign().ID("BuildFakeWebhook").Call(),
			jen.Return(
				jen.ID("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("webhook")),
			),
		),
	}

	return lines
}

func buildBuildFakeWebhookCreationInputFromWebhook(proj *models.Project) []jen.Code {
	funcName := "BuildFakeWebhookCreationInputFromWebhook"
	typeName := "WebhookCreationInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("webhook").PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.AddressOf().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("Name").MapAssign().ID("webhook").Dot("Name"),
					jen.ID("ContentType").MapAssign().ID("webhook").Dot("ContentType"),
					jen.ID("URL").MapAssign().ID("webhook").Dot("URL"),
					jen.ID("Method").MapAssign().ID("webhook").Dot("Method"),
					jen.ID("Events").MapAssign().ID("webhook").Dot("Events"),
					jen.ID("DataTypes").MapAssign().ID("webhook").Dot("DataTypes"),
					jen.ID("Topics").MapAssign().ID("webhook").Dot("Topics"),
					jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("webhook").Dot(constants.UserOwnershipFieldName),
				),
			),
		),
	}

	return lines
}
