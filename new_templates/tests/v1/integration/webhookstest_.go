package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func webhooksTestDotGo() *jen.File {
	ret := jen.NewFile("integration")
	utils.AddImports(ret)

	ret.Add(jen.Func().ID("checkWebhookEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").ID("models").Dot(
		"Webhook",
	)).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.ID("assert").Dot(
			"NotZero",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"ID",
		)),
		jen.ID("assert").Dot(
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"Name",
		), jen.ID("actual").Dot(
			"Name",
		)),
		jen.ID("assert").Dot(
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"ContentType",
		), jen.ID("actual").Dot(
			"ContentType",
		)),
		jen.ID("assert").Dot(
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"URL",
		), jen.ID("actual").Dot(
			"URL",
		)),
		jen.ID("assert").Dot(
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"Method",
		), jen.ID("actual").Dot(
			"Method",
		)),
		jen.ID("assert").Dot(
			"NotZero",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"CreatedOn",
		)),
	),
	)
	ret.Add(jen.Func().ID("buildDummyWebhookInput").Params().Params(jen.Op("*").ID("models").Dot(
		"WebhookCreationInput",
	)).Block(
		jen.ID("x").Op(":=").Op("&").ID("models").Dot(
			"WebhookCreationInput",
		).Valuesln(jen.ID("Name").Op(":").ID("fake").Dot(
			"Word",
		).Call(), jen.ID("URL").Op(":").ID("fake").Dot(
			"DomainName",
		).Call(), jen.ID("ContentType").Op(":").Lit("application/json"), jen.ID("Method").Op(":").Qual("net/http", "MethodPost")),
		jen.Return().ID("x"),
	),
	)
	ret.Add(jen.Func().ID("buildDummyWebhook").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("models").Dot(
		"Webhook",
	)).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.List(jen.ID("y"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
			"CreateWebhook",
		).Call(jen.Qual("context", "Background").Call(), jen.ID("buildDummyWebhookInput").Call()),
		utils.RequireNoError(jen.ID("t"), jen.ID("err")),
		jen.Return().ID("y"),
	),
	)
	ret.Add(jen.Func().ID("reverse").Params(jen.ID("s").ID("string")).Params(jen.ID("string")).Block(
		jen.ID("runes").Op(":=").Index().ID("rune").Call(jen.ID("s")),
		jen.For(jen.List(jen.ID("i"), jen.ID("j")).Op(":=").List(jen.Lit(0), jen.ID("len").Call(jen.ID("runes")).Op("-").Lit(1)), jen.ID("i").Op("<").ID("j"), jen.List(jen.ID("i"), jen.ID("j")).Op("=").List(jen.ID("i").Op("+").Lit(1), jen.ID("j").Op("-").Lit(1))).Block(
			jen.List(jen.ID("runes").Index(jen.ID("i")), jen.ID("runes").Index(jen.ID("j"))).Op("=").List(jen.ID("runes").Index(jen.ID("j")), jen.ID("runes").Index(jen.ID("i"))),
		),
		jen.Return().ID("string").Call(jen.ID("runes")),
	),
	)
	ret.Add(jen.Func().ID("TestWebhooks").Params(jen.ID("test").Op("*").Qual("testing", "T")).Block(
		jen.ID("test").Dot(
			"Parallel",
		).Call(),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be createable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("input").Op(":=").ID("buildDummyWebhookInput").Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Webhook",
				).Valuesln(jen.ID("Name").Op(":").ID("input").Dot(
					"Name",
				), jen.ID("URL").Op(":").ID("input").Dot(
					"URL",
				), jen.ID("ContentType").Op(":").ID("input").Dot(
					"ContentType",
				), jen.ID("Method").Op(":").ID("input").Dot(
					"Method",
				)),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateWebhook",
				).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
					"WebhookCreationInput",
				).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
					"Name",
				), jen.ID("ContentType").Op(":").ID("expected").Dot(
					"ContentType",
				), jen.ID("URL").Op(":").ID("expected").Dot(
					"URL",
				), jen.ID("Method").Op(":").ID("expected").Dot(
					"Method",
				))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("premade")),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"ArchiveWebhook",
				).Call(jen.ID("tctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetWebhook",
				).Call(jen.ID("tctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot(
					"NotZero",
				).Call(jen.ID("t"), jen.ID("actual").Dot(
					"ArchivedOn",
				)),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be able to be read in a list"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.Null().Var().ID("expected").Index().Op("*").ID("models").Dot(
					"Webhook",
				),
				jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
					jen.ID("expected").Op("=").ID("append").Call(jen.ID("expected"), jen.ID("buildDummyWebhook").Call(jen.ID("t"))),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetWebhooks",
				).Call(jen.ID("tctx"), jen.ID("nil")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				utils.AssertTrue(jen.ID("t"), jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot(
					"Webhooks",
				))),
				jen.For(jen.List(jen.ID("_"), jen.ID("webhook")).Op(":=").Range().ID("actual").Dot(
					"Webhooks",
				)).Block(
					jen.ID("err").Op("=").ID("todoClient").Dot(
						"ArchiveWebhook",
					).Call(jen.ID("tctx"), jen.ID("webhook").Dot(
						"ID",
					)),
					jen.ID("assert").Dot(
						"NoError",
					).Call(jen.ID("t"), jen.ID("err")),
				),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should return an error when trying to read something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetWebhook",
				).Call(jen.ID("tctx"), jen.ID("nonexistentID")),
				jen.ID("assert").Dot(
					"Error",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should be readable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("input").Op(":=").ID("buildDummyWebhookInput").Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Webhook",
				).Valuesln(jen.ID("Name").Op(":").ID("input").Dot(
					"Name",
				), jen.ID("URL").Op(":").ID("input").Dot(
					"URL",
				), jen.ID("ContentType").Op(":").ID("input").Dot(
					"ContentType",
				), jen.ID("Method").Op(":").ID("input").Dot(
					"Method",
				)),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateWebhook",
				).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
					"WebhookCreationInput",
				).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
					"Name",
				), jen.ID("ContentType").Op(":").ID("expected").Dot(
					"ContentType",
				), jen.ID("URL").Op(":").ID("expected").Dot(
					"URL",
				), jen.ID("Method").Op(":").ID("expected").Dot(
					"Method",
				))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetWebhook",
				).Call(jen.ID("tctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"ArchiveWebhook",
				).Call(jen.ID("tctx"), jen.ID("actual").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Updating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should return an error when trying to update something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("err").Op(":=").ID("todoClient").Dot(
					"UpdateWebhook",
				).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
					"Webhook",
				).Valuesln(jen.ID("ID").Op(":").ID("nonexistentID"))),
				jen.ID("assert").Dot(
					"Error",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should be updatable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("input").Op(":=").ID("buildDummyWebhookInput").Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Webhook",
				).Valuesln(jen.ID("Name").Op(":").ID("input").Dot(
					"Name",
				), jen.ID("URL").Op(":").ID("input").Dot(
					"URL",
				), jen.ID("ContentType").Op(":").ID("input").Dot(
					"ContentType",
				), jen.ID("Method").Op(":").ID("input").Dot(
					"Method",
				)),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateWebhook",
				).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
					"WebhookCreationInput",
				).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
					"Name",
				), jen.ID("ContentType").Op(":").ID("expected").Dot(
					"ContentType",
				), jen.ID("URL").Op(":").ID("expected").Dot(
					"URL",
				), jen.ID("Method").Op(":").ID("expected").Dot(
					"Method",
				))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.ID("premade").Dot(
					"Name",
				).Op("=").ID("reverse").Call(jen.ID("premade").Dot(
					"Name",
				)),
				jen.ID("expected").Dot(
					"Name",
				).Op("=").ID("premade").Dot(
					"Name",
				),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"UpdateWebhook",
				).Call(jen.ID("tctx"), jen.ID("premade")),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetWebhook",
				).Call(jen.ID("tctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot(
					"NotNil",
				).Call(jen.ID("t"), jen.ID("actual").Dot(
					"UpdatedOn",
				)),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"ArchiveWebhook",
				).Call(jen.ID("tctx"), jen.ID("actual").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be able to be deleted"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("input").Op(":=").ID("buildDummyWebhookInput").Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Webhook",
				).Valuesln(jen.ID("Name").Op(":").ID("input").Dot(
					"Name",
				), jen.ID("URL").Op(":").ID("input").Dot(
					"URL",
				), jen.ID("ContentType").Op(":").ID("input").Dot(
					"ContentType",
				), jen.ID("Method").Op(":").ID("input").Dot(
					"Method",
				)),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateWebhook",
				).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
					"WebhookCreationInput",
				).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
					"Name",
				), jen.ID("ContentType").Op(":").ID("expected").Dot(
					"ContentType",
				), jen.ID("URL").Op(":").ID("expected").Dot(
					"URL",
				), jen.ID("Method").Op(":").ID("expected").Dot(
					"Method",
				))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"ArchiveWebhook",
				).Call(jen.ID("tctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
		)),
	),
	)
	return ret
}
