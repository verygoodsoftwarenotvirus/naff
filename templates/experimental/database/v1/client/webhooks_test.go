package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func webhooksTestDotGo() *jen.File {
	ret := jen.NewFile("dbclient")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestClient_GetWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleID"), jen.ID("exampleUserID")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleID"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestClient_GetWebhookCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetWebhookCount"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetWebhookCount",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetWebhookCount"), jen.ID("mock").Dot(
				"Anything",
			), jen.Parens(jen.Op("*").ID("models").Dot(
				"QueryFilter",
			)).Call(jen.ID("nil")), jen.ID("exampleUserID")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetWebhookCount",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("nil"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestClient_GetAllWebhooksCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetAllWebhooksCount"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetAllWebhooksCount",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestClient_GetAllWebhooks").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetAllWebhooks"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetAllWebhooks",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestClient_GetWebhooks").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetWebhooks"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetWebhooks",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetWebhooks"), jen.ID("mock").Dot(
				"Anything",
			), jen.Parens(jen.Op("*").ID("models").Dot(
				"QueryFilter",
			)).Call(jen.ID("nil")), jen.ID("exampleUserID")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetWebhooks",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("nil"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestClient_CreateWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"CreateWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestClient_UpdateWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(),
			jen.Null().Var().ID("expected").ID("error"),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.ID("expected")),
			jen.ID("actual").Op(":=").ID("c").Dot(
				"UpdateWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestClient_ArchiveWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.Null().Var().ID("expected").ID("error"),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("ArchiveWebhook"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleID"), jen.ID("exampleUserID")).Dot(
				"Return",
			).Call(jen.ID("expected")),
			jen.ID("actual").Op(":=").ID("c").Dot(
				"ArchiveWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleID"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),

		jen.Line(),
	)
	return ret
}
