package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(buildTestClient_GetOAuth2Client(proj)...)
	ret.Add(buildTestClient_GetOAuth2ClientByClientID(proj)...)
	ret.Add(buildTestClient_GetAllOAuth2ClientCount(proj)...)
	ret.Add(buildTestClient_GetAllOAuth2Clients(proj)...)
	ret.Add(buildTestClient_GetOAuth2Clients(proj)...)
	ret.Add(buildTestClient_CreateOAuth2Client(proj)...)
	ret.Add(buildTestClient_UpdateOAuth2Client(proj)...)
	ret.Add(buildTestClient_ArchiveOAuth2Client(proj)...)

	return ret
}

func buildTestClient_GetOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2Client").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleClientID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error returned from querier",
				utils.CreateCtx(),
				jen.ID("exampleClientID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
		),
	}

	return lines
}

func buildTestClient_GetOAuth2ClientByClientID(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2ClientByClientID").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleClientID").Assign().Lit("CLIENT_ID"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("exampleClientID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error returned from querier",
				utils.CreateCtx(),
				jen.ID("exampleClientID").Assign().Lit("CLIENT_ID"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("expected").Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID")).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("exampleClientID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
		),
	}

	return lines
}

func buildTestClient_GetOAuth2ClientCount(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2ClientCount").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				utils.CreateDefaultQueryFilter(proj),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientCount").Call(
					utils.CtxVar(),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		jen.Line(),
		utils.BuildSubTestWithoutContext(
			"with nil filter",
			utils.CreateCtx(),
			jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
			jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
			utils.CreateNilQueryFilter(proj),
			jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
				jen.Lit("GetOAuth2ClientCount"),
				jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				jen.ID("exampleUserID"),
				jen.ID(utils.FilterVarName),
			).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientCount").Call(
				utils.CtxVar(),
				jen.ID("exampleUserID"),
				jen.ID(utils.FilterVarName),
			),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			jen.Line(),
			jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
		),
		jen.Line(),
		utils.BuildSubTestWithoutContext(
			"with error returned from querier",
			utils.CreateCtx(),
			jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
			jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
			utils.CreateDefaultQueryFilter(proj),
			jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
				jen.Lit("GetOAuth2ClientCount"),
				jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				jen.ID("exampleUserID"),
				jen.ID(utils.FilterVarName),
			).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientCount").Call(
				utils.CtxVar(),
				jen.ID("exampleUserID"),
				jen.ID(utils.FilterVarName),
			),
			utils.AssertError(jen.Err(), nil),
			utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			jen.Line(),
			jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetAllOAuth2ClientCount(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetAllOAuth2ClientCount").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetAllOAuth2ClientCount"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllOAuth2ClientCount").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetAllOAuth2Clients(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetAllOAuth2Clients").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.Var().ID("expected").Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetAllOAuth2Clients"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetOAuth2Clients(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2Clients").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Values(),
				utils.CreateDefaultQueryFilter(proj),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Clients"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(
					utils.CtxVar(),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil filter",
				utils.CreateCtx(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Values(),
				utils.CreateNilQueryFilter(proj),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Clients"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(
					utils.CtxVar(),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error returned from querier",
				utils.CreateCtx(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList")).Call(jen.Nil()),
				utils.CreateDefaultQueryFilter(proj),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Clients"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(
					utils.CtxVar(),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
		),
	}

	return lines
}

func buildTestClient_CreateOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_CreateOAuth2Client").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(),
				jen.ID("exampleInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Values(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("CreateOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error returned from querier",
				utils.CreateCtx(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("expected").Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
				jen.ID("exampleInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Values(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("CreateOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
		),
	}

	return lines
}

func buildTestClient_UpdateOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_UpdateOAuth2Client").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("example").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(),
				jen.Var().ID("expected").ID("error"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("UpdateOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("example")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("UpdateOAuth2Client").Call(utils.CtxVar(), jen.ID("example")),
				utils.AssertNoError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_ArchiveOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_ArchiveOAuth2Client").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleClientID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Var().ID("expected").ID("error"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("ArchiveOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				utils.AssertNoError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error returned from querier",
				utils.CreateCtx(),
				jen.ID("exampleClientID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Qual("fmt", "Errorf").Call(jen.Lit("blah")),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("ArchiveOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				utils.AssertError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
		),
	}

	return lines
}
