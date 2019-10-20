package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func webhooksTestDotGo() *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("TestClient_GetWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Webhook",
				).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetWebhook"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleID"), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetWebhook",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetWebhookCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetWebhookCount"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("models").Dot(
						"DefaultQueryFilter",
					).Call(), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetWebhookCount",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetWebhookCount"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.Parens(jen.Op("*").ID("models").Dot(
						"QueryFilter",
					)).Call(jen.ID("nil")), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetWebhookCount",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("nil"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetAllWebhooksCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetAllWebhooksCount"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				)).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetAllWebhooksCount",
				).Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetAllWebhooks").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"WebhookList",
				).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetAllWebhooks"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				)).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetAllWebhooks",
				).Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetWebhooks").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"WebhookList",
				).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetWebhooks"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("models").Dot(
						"DefaultQueryFilter",
					).Call(), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetWebhooks",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"WebhookList",
				).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetWebhooks"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.Parens(jen.Op("*").ID("models").Dot(
						"QueryFilter",
					)).Call(jen.ID("nil")), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetWebhooks",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("nil"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_CreateWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
					"WebhookCreationInput",
				).Values(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Webhook",
				).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("CreateWebhook"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleInput")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"CreateWebhook",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_UpdateWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
					"Webhook",
				).Values(),

				jen.Var().ID("expected").ID("error"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("UpdateWebhook"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleInput")).Dot(
					"Return",
				).Call(jen.ID("expected")),
				jen.ID("actual").Op(":=").ID("c").Dot(
					"UpdateWebhook",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_ArchiveWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),

				jen.Var().ID("expected").ID("error"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"WebhookDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("ArchiveWebhook"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleID"), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected")),
				jen.ID("actual").Op(":=").ID("c").Dot(
					"ArchiveWebhook",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("mockDB").Dot(
					"AssertExpectations",
				).Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)
	return ret
}
