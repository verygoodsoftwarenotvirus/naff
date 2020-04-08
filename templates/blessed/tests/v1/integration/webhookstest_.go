package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("checkWebhookEquality").Params(jen.ID("t").ParamPointer().Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			utils.AssertNotZero(jen.ID("actual").Dot("ID"), nil),
			utils.AssertEqual(jen.ID("expected").Dot("Name"), jen.ID("actual").Dot("Name"), nil),
			utils.AssertEqual(jen.ID("expected").Dot("ContentType"), jen.ID("actual").Dot("ContentType"), nil),
			utils.AssertEqual(jen.ID("expected").Dot("URL"), jen.ID("actual").Dot("URL"), nil),
			utils.AssertEqual(jen.ID("expected").Dot("Method"), jen.ID("actual").Dot("Method"), nil),
			utils.AssertNotZero(jen.ID("actual").Dot("CreatedOn"), nil),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyWebhookInput").Params().Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput")).Block(
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "WebhookCreationInput").Valuesln(
				jen.ID("Name").MapAssign().Qual(utils.FakeLibrary, "Word").Call(),
				jen.ID("URL").MapAssign().Qual(utils.FakeLibrary, "DomainName").Call(),
				jen.ID("ContentType").MapAssign().Lit("application/json"),
				jen.ID("Method").MapAssign().Qual("net/http", "MethodPost"),
			),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyWebhook").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			utils.CreateCtx(),
			jen.Line(),
			jen.List(jen.ID("y"), jen.Err()).Assign().ID("todoClient").Dot("CreateWebhook").Call(utils.CtxVar(), jen.ID("buildDummyWebhookInput").Call()),
			utils.RequireNoError(jen.Err(), nil),
			jen.Return().ID("y"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("reverse").Params(jen.ID("s").String()).Params(jen.String()).Block(
			jen.ID("runes").Assign().Index().ID("rune").Call(jen.ID("s")),
			jen.For(jen.List(jen.ID("i"), jen.ID("j")).Assign().List(jen.Zero(), jen.ID("len").Call(jen.ID("runes")).Op("-").One()), jen.ID("i").Op("<").ID("j"), jen.List(jen.ID("i"), jen.ID("j")).Equals().List(jen.ID("i").Op("+").One(), jen.ID("j").Op("-").One())).Block(
				jen.List(jen.ID("runes").Index(jen.ID("i")), jen.ID("runes").Index(jen.ID("j"))).Equals().List(jen.ID("runes").Index(jen.ID("j")), jen.ID("runes").Index(jen.ID("i"))),
			),
			jen.Return().String().Call(jen.ID("runes")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestWebhooks").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be createable",
					jen.Line(),
					jen.Comment("Create webhook"),
					jen.ID("input").Assign().ID("buildDummyWebhookInput").Call(),
					jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "Webhook").Valuesln(
						jen.ID("Name").MapAssign().ID("input").Dot("Name"),
						jen.ID("URL").MapAssign().ID("input").Dot("URL"),
						jen.ID("ContentType").MapAssign().ID("input").Dot("ContentType"),
						jen.ID("Method").MapAssign().ID("input").Dot("Method")),
					jen.List(jen.ID("premade"), jen.Err()).Assign().ID("todoClient").Dot("CreateWebhook").Call(utils.CtxVar(), jen.AddressOf().Qual(proj.ModelsV1Package(), "WebhookCreationInput").Valuesln(
						jen.ID("Name").MapAssign().ID("expected").Dot("Name"),
						jen.ID("ContentType").MapAssign().ID("expected").Dot("ContentType"),
						jen.ID("URL").MapAssign().ID("expected").Dot("URL"),
						jen.ID("Method").MapAssign().ID("expected").Dot("Method"))),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert webhook equality"),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("premade")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().ID("todoClient").Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("premade").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
					jen.Line(), // REPEATME NOTICEME
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					utils.AssertNotZero(jen.ID("actual").Dot("ArchivedOn"), nil),
				),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be able to be read in a list",
					jen.Line(),
					jen.Comment("Create webhooks"),
					jen.Var().ID("expected").Index().PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
					jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
						jen.ID("expected").Equals().ID("append").Call(jen.ID("expected"), jen.ID("buildDummyWebhook").Call(jen.ID("t"))),
					),
					jen.Line(),
					jen.Comment("Assert webhook list equality"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetWebhooks").Call(utils.CtxVar(), jen.Nil()),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					utils.AssertTrue(jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Webhooks")), nil),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.For(jen.List(jen.Underscore(), jen.ID("webhook")).Assign().Range().ID("actual").Dot("Webhooks")).Block(
						jen.Err().Equals().ID("todoClient").Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("webhook").Dot("ID")),
						utils.AssertNoError(jen.Err(), nil),
					),
				),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"it should return an error when trying to read something that doesn't exist",
					jen.Line(),
					jen.Comment("Fetch webhook"),
					jen.List(jen.Underscore(), jen.Err()).Assign().ID("todoClient").Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("nonexistentID")),
					utils.AssertError(jen.Err(), nil),
				),
				jen.Line(),
				utils.BuildSubTest(
					"it should be readable",
					jen.Line(),
					jen.Comment("Create webhook"),
					jen.ID("input").Assign().ID("buildDummyWebhookInput").Call(),
					jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "Webhook").Valuesln(
						jen.ID("Name").MapAssign().ID("input").Dot("Name"),
						jen.ID("URL").MapAssign().ID("input").Dot("URL"),
						jen.ID("ContentType").MapAssign().ID("input").Dot("ContentType"),
						jen.ID("Method").MapAssign().ID("input").Dot("Method")),
					jen.List(jen.ID("premade"), jen.Err()).Assign().ID("todoClient").Dot("CreateWebhook").Call(
						utils.CtxVar(),
						jen.AddressOf().Qual(proj.ModelsV1Package(), "WebhookCreationInput").Valuesln(
							jen.ID("Name").MapAssign().ID("expected").Dot("Name"),
							jen.ID("ContentType").MapAssign().ID("expected").Dot("ContentType"),
							jen.ID("URL").MapAssign().ID("expected").Dot("URL"),
							jen.ID("Method").MapAssign().ID("expected").Dot("Method"),
						),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Fetch webhook"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert webhook equality"),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().ID("todoClient").Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("actual").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("Updating"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"it should return an error when trying to update something that doesn't exist",
					jen.Line(),
					jen.Err().Assign().ID("todoClient").Dot("UpdateWebhook").Call(utils.CtxVar(), jen.AddressOf().Qual(proj.ModelsV1Package(), "Webhook").Values(jen.ID("ID").MapAssign().ID("nonexistentID"))),
					utils.AssertError(jen.Err(), nil),
				),
				jen.Line(),
				utils.BuildSubTest(
					"it should be updatable",
					jen.Line(),
					jen.Comment("Create webhook"),
					jen.ID("input").Assign().ID("buildDummyWebhookInput").Call(),
					jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "Webhook").Valuesln(
						jen.ID("Name").MapAssign().ID("input").Dot("Name"),
						jen.ID("URL").MapAssign().ID("input").Dot("URL"),
						jen.ID("ContentType").MapAssign().ID("input").Dot("ContentType"),
						jen.ID("Method").MapAssign().ID("input").Dot("Method")),
					jen.List(jen.ID("premade"), jen.Err()).Assign().ID("todoClient").Dot("CreateWebhook").Call(
						utils.CtxVar(),
						jen.AddressOf().Qual(proj.ModelsV1Package(), "WebhookCreationInput").Valuesln(
							jen.ID("Name").MapAssign().ID("expected").Dot("Name"),
							jen.ID("ContentType").MapAssign().ID("expected").Dot("ContentType"),
							jen.ID("URL").MapAssign().ID("expected").Dot("URL"),
							jen.ID("Method").MapAssign().ID("expected").Dot("Method"),
						),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Change webhook"),
					jen.ID("premade").Dot("Name").Equals().ID("reverse").Call(jen.ID("premade").Dot("Name")),
					jen.ID("expected").Dot("Name").Equals().ID("premade").Dot("Name"),
					jen.Err().Equals().ID("todoClient").Dot("UpdateWebhook").Call(utils.CtxVar(), jen.ID("premade")),
					utils.AssertNoError(jen.Err(), nil),
					jen.Line(),
					jen.Comment("Fetch webhook"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert webhook equality"),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					utils.AssertNotNil(jen.ID("actual").Dot("UpdatedOn"), nil),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().ID("todoClient").Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("actual").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be able to be deleted",
					jen.Line(),
					jen.Comment("Create webhook"),
					jen.ID("input").Assign().ID("buildDummyWebhookInput").Call(),
					jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "Webhook").Valuesln(
						jen.ID("Name").MapAssign().ID("input").Dot("Name"),
						jen.ID("URL").MapAssign().ID("input").Dot("URL"),
						jen.ID("ContentType").MapAssign().ID("input").Dot("ContentType"),
						jen.ID("Method").MapAssign().ID("input").Dot("Method"),
					),
					jen.List(jen.ID("premade"), jen.Err()).Assign().ID("todoClient").Dot("CreateWebhook").Call(utils.CtxVar(), jen.AddressOf().Qual(proj.ModelsV1Package(), "WebhookCreationInput").Valuesln(
						jen.ID("Name").MapAssign().ID("expected").Dot("Name"),
						jen.ID("ContentType").MapAssign().ID("expected").Dot("ContentType"),
						jen.ID("URL").MapAssign().ID("expected").Dot("URL"),
						jen.ID("Method").MapAssign().ID("expected").Dot("Method")),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().ID("todoClient").Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("premade").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
			jen.Line(),
		),
	)

	return ret
}
