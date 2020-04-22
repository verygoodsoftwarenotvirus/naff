package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("checkWebhookEquality").Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Block(
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
		jen.Func().ID("reverse").Params(jen.ID("s").String()).Params(jen.String()).Block(
			jen.ID("runes").Assign().Index().ID("rune").Call(jen.ID("s")),
			jen.For(jen.List(jen.ID("i"), jen.ID("j")).Assign().List(jen.Zero(), jen.Len(jen.ID("runes")).Minus().One()), jen.ID("i").LessThan().ID("j"), jen.List(jen.ID("i"), jen.ID("j")).Equals().List(jen.ID("i").Plus().One(), jen.ID("j").Minus().One())).Block(
				jen.List(jen.ID("runes").Index(jen.ID("i")), jen.ID("runes").Index(jen.ID("j"))).Equals().List(jen.ID("runes").Index(jen.ID("j")), jen.ID("runes").Index(jen.ID("i"))),
			),
			jen.Return().String().Call(jen.ID("runes")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestWebhooks").Params(jen.ID("test").PointerTo().Qual("testing", "T")).Block(
			jen.ID("test").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be createable",
					utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create webhook"),
					utils.BuildFakeVar(proj, "Webhook"),
					utils.BuildFakeVarWithCustomName(
						proj,
						"exampleWebhookInput",
						"BuildFakeWebhookCreationInputFromWebhook",
						jen.ID(utils.BuildFakeVarName("Webhook")),
					),
					jen.List(jen.ID("premade"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateWebhook").Call(
						constants.CtxVar(),
						jen.ID("exampleWebhookInput"),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert webhook equality"),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("premade")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveWebhook").Call(
						constants.CtxVar(),
						jen.ID("premade").Dot("ID"),
					),
					utils.AssertNoError(jen.Err(), nil),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("GetWebhook").Call(
						constants.CtxVar(),
						jen.ID("premade").Dot("ID"),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("actual")),
					utils.AssertNotZero(jen.ID("actual").Dot("ArchivedOn"), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be able to be read in a list",
					utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create webhooks"),
					jen.Var().ID("expected").Index().PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
					jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().Lit(5), jen.ID("i").Op("++")).Block(
						utils.BuildFakeVar(proj, "Webhook"),
						utils.BuildFakeVarWithCustomName(
							proj,
							"exampleWebhookInput",
							"BuildFakeWebhookCreationInputFromWebhook",
							jen.ID(utils.BuildFakeVarName("Webhook")),
						),
						jen.List(jen.ID("createdWebhook"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateWebhook").Call(
							constants.CtxVar(),
							jen.ID("exampleWebhookInput"),
						),
						jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("createdWebhook"), jen.Err()),
						jen.Line(),
						utils.AppendItemsToList(jen.ID("expected"), jen.ID("createdWebhook")),
					),
					jen.Line(),
					jen.Comment("Assert webhook list equality"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("GetWebhooks").Call(constants.CtxVar(), jen.Nil()),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					utils.AssertTrue(jen.Len(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Webhooks")), nil),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.For(jen.List(jen.Underscore(), jen.ID("webhook")).Assign().Range().ID("actual").Dot("Webhooks")).Block(
						jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveWebhook").Call(constants.CtxVar(), jen.ID("webhook").Dot("ID")),
						utils.AssertNoError(jen.Err(), nil),
					),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"it should return an error when trying to read something that doesn't exist",
					utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Fetch webhook"),
					jen.List(jen.Underscore(), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("GetWebhook").Call(constants.CtxVar(), jen.ID("nonexistentID")),
					utils.AssertError(jen.Err(), nil),
				),
				jen.Line(),
				utils.BuildSubTest(
					"it should be readable",
					utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create webhook"), utils.BuildFakeVar(proj, "Webhook"),
					utils.BuildFakeVarWithCustomName(
						proj,
						"exampleWebhookInput",
						"BuildFakeWebhookCreationInputFromWebhook",
						jen.ID(utils.BuildFakeVarName("Webhook")),
					),
					jen.List(jen.ID("premade"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateWebhook").Call(
						constants.CtxVar(),
						jen.ID("exampleWebhookInput"),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Fetch webhook"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("GetWebhook").Call(constants.CtxVar(), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert webhook equality"),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveWebhook").Call(constants.CtxVar(), jen.ID("actual").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Updating"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"it should return an error when trying to update something that doesn't exist",
					utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					utils.BuildFakeVar(proj, "Webhook"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID").Equals().ID("nonexistentID"),
					jen.Line(),
					jen.Err().Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("UpdateWebhook").Call(
						constants.CtxVar(),
						jen.ID(utils.BuildFakeVarName("Webhook")),
					),
					utils.AssertError(jen.Err(), nil),
				),
				jen.Line(),
				utils.BuildSubTest(
					"it should be updatable",
					utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create webhook"),
					utils.BuildFakeVar(proj, "Webhook"),
					utils.BuildFakeVarWithCustomName(
						proj,
						"exampleWebhookInput",
						"BuildFakeWebhookCreationInputFromWebhook",
						jen.ID(utils.BuildFakeVarName("Webhook")),
					),
					jen.List(jen.ID("premade"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateWebhook").Call(
						constants.CtxVar(),
						jen.ID("exampleWebhookInput"),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Change webhook"),
					jen.ID("premade").Dot("Name").Equals().ID("reverse").Call(jen.ID("premade").Dot("Name")),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Name").Equals().ID("premade").Dot("Name"),
					jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("UpdateWebhook").Call(constants.CtxVar(), jen.ID("premade")),
					utils.AssertNoError(jen.Err(), nil),
					jen.Line(),
					jen.Comment("Fetch webhook"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("GetWebhook").Call(constants.CtxVar(), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert webhook equality"),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("actual")),
					utils.AssertNotNil(jen.ID("actual").Dot("UpdatedOn"), nil),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveWebhook").Call(constants.CtxVar(), jen.ID("actual").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be able to be deleted",
					utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create webhook"),
					utils.BuildFakeVar(proj, "Webhook"),
					utils.BuildFakeVarWithCustomName(
						proj,
						"exampleWebhookInput",
						"BuildFakeWebhookCreationInputFromWebhook",
						jen.ID(utils.BuildFakeVarName("Webhook")),
					),
					jen.List(jen.ID("premade"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateWebhook").Call(
						constants.CtxVar(),
						jen.ID("exampleWebhookInput"),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveWebhook").Call(constants.CtxVar(), jen.ID("premade").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
		),
	)

	return ret
}
