package dbclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersTestDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)


	ret.Add(jen.Func().ID("TestClient_GetUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleID")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleID")),
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
	)
	ret.Add(jen.Func().ID("TestClient_GetUserByUsername").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUsername").Op(":=").Lit("username"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUsername")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetUserByUsername",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUsername")),
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
	)
	ret.Add(jen.Func().ID("TestClient_GetUserCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserCount"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call()).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetUserCount",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call()),
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
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserCount"), jen.ID("mock").Dot(
				"Anything",
			), jen.Parens(jen.Op("*").ID("models").Dot(
				"QueryFilter",
			)).Call(jen.ID("nil"))).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetUserCount",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
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
	)
	ret.Add(jen.Func().ID("TestClient_GetUsers").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"UserList",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUsers"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call()).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetUsers",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call()),
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
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"UserList",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUsers"), jen.ID("mock").Dot(
				"Anything",
			), jen.Parens(jen.Op("*").ID("models").Dot(
				"QueryFilter",
			)).Call(jen.ID("nil"))).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetUsers",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
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
	)
	ret.Add(jen.Func().ID("TestClient_CreateUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"UserInput",
			).Valuesln(),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"CreateUser",
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
	)
	ret.Add(jen.Func().ID("TestClient_UpdateUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(),
			jen.Null().Var().ID("expected").ID("error"),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("err").Op(":=").ID("c").Dot(
				"UpdateUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestClient_ArchiveUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.Null().Var().ID("expected").ID("error"),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("ArchiveUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("err").Op(":=").ID("c").Dot(
				"ArchiveUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("mockDB").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),
	)
	return ret
}
