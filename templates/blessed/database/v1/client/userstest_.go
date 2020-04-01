package client

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("TestClient_GetUser").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleID")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUser").Call(utils.CtxVar(), jen.ID("exampleID")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetUserByUsername").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleUsername").Assign().Lit("username"),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("GetUserByUsername"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleUsername")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUserByUsername").Call(utils.CtxVar(), jen.ID("exampleUsername")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetUserCount").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				utils.CreateDefaultQueryFilter(pkg),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUserCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUserCount").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				utils.CreateNilQueryFilter(pkg),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUserCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUserCount").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetUsers").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserList").Values(),
				utils.CreateDefaultQueryFilter(pkg),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUsers"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserList").Values(),
				utils.CreateNilQueryFilter(pkg),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUsers"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_CreateUser").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleInput").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserInput").Values(),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("CreateUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(utils.CtxVar(), jen.ID("exampleInput")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_UpdateUser").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleInput").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("UpdateUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("UpdateUser").Call(utils.CtxVar(), jen.ID("exampleInput")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_ArchiveUser").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleInput").Assign().Add(utils.FakeUint64Func()),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("ArchiveUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("exampleInput")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)
	return ret
}
