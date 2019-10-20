package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsTestDotGo() *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("TestClient_GetItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleItemID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Item",
				).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"ItemDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetItem"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleItemID"), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetItem",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleItemID"), jen.ID("exampleUserID")),
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
		jen.Func().ID("TestClient_GetItemCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"ItemDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetItemCount"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("models").Dot(
						"DefaultQueryFilter",
					).Call(), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetItemCount",
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
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"ItemDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetItemCount"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.Parens(jen.Op("*").ID("models").Dot(
						"QueryFilter",
					)).Call(jen.ID("nil")), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetItemCount",
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
		jen.Func().ID("TestClient_GetAllItemsCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"ItemDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetAllItemsCount"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				)).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetAllItemsCount",
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
		jen.Func().ID("TestClient_GetItems").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"ItemList",
				).Values(),
				jen.ID("mockDB").Dot(
					"ItemDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetItems"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("models").Dot(
						"DefaultQueryFilter",
					).Call(), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetItems",
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
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"ItemList",
				).Values(),
				jen.ID("mockDB").Dot(
					"ItemDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("GetItems"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.Parens(jen.Op("*").ID("models").Dot(
						"QueryFilter",
					)).Call(jen.ID("nil")), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"GetItems",
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
		jen.Func().ID("TestClient_CreateItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
					"ItemCreationInput",
				).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Item",
				).Values(),
				jen.ID("mockDB").Dot(
					"ItemDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("CreateItem"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleInput")).Dot(
					"Return",
				).Call(jen.ID("expected"), jen.ID("nil")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
					"CreateItem",
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
		jen.Func().ID("TestClient_UpdateItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
					"Item",
				).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),

				jen.Var().ID("expected").ID("error"),
				jen.ID("mockDB").Dot(
					"ItemDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("UpdateItem"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleInput")).Dot(
					"Return",
				).Call(jen.ID("expected")),
				jen.ID("err").Op(":=").ID("c").Dot(
					"UpdateItem",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_ArchiveItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleItemID").Op(":=").ID("uint64").Call(jen.Lit(123)),

				jen.Var().ID("expected").ID("error"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot(
					"ItemDataManager",
				).Dot(
					"On",
				).Call(jen.Lit("ArchiveItem"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleItemID"), jen.ID("exampleUserID")).Dot(
					"Return",
				).Call(jen.ID("expected")),
				jen.ID("err").Op(":=").ID("c").Dot(
					"ArchiveItem",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUserID"), jen.ID("exampleItemID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
		),
		jen.Line(),
	)
	return ret
}
