package frontend

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("buildRequest").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().Qual("net/http", "Request")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodGet"),
				jen.Lit("https://verygoodsoftwarenotvirus.ru"),
				jen.Nil(),
			),
			jen.Line(),
			utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
			utils.AssertNoError(jen.Err(), nil),
			jen.Return().ID(constants.RequestVarName),
		),
		jen.Line(),
	)

	ret.Add(buildTestService_StaticDir(proj)...)

	ret.Add(
		jen.Func().ID("TestService_Routes").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				utils.AssertNotNil(jen.Parens(jen.AddressOf().ID("Service").Values()).Dot("Routes").Call(), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_buildStaticFileServer").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("s").Assign().AddressOf().ID("Service").Valuesln(
					jen.ID("config").MapAssign().Qual(proj.InternalConfigV1Package(), "FrontendSettings").Valuesln(
						jen.ID("CacheStaticFiles").MapAssign().True(),
					),
				),
				jen.List(jen.ID("cwd"), jen.Err()).Assign().Qual("os", "Getwd").Call(),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("buildStaticFileServer").Call(jen.ID("cwd")),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
		),
		jen.Line(),
	)
	return ret
}

func buildTestService_StaticDir(proj *models.Project) []jen.Code {
	block := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		utils.BuildSubTestWithoutContext(
			"happy path",
			jen.ID("s").Assign().AddressOf().ID("Service").Values(jen.ID(constants.LoggerVarName).MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			jen.List(jen.ID("cwd"), jen.Err()).Assign().Qual("os", "Getwd").Call(),
			utils.RequireNoError(jen.Err(), nil),
			jen.Line(),
			jen.List(jen.ID("hf"), jen.Err()).Assign().ID("s").Dot("StaticDir").Call(jen.ID("cwd")),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertNotNil(jen.ID("hf"), nil),
			jen.Line(),
			jen.List(jen.ID(constants.RequestVarName), jen.ID(constants.ResponseVarName)).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
			jen.ID(constants.RequestVarName).Dot("URL").Dot("Path").Equals().Lit("/http_routes_test.go"),
			jen.ID("hf").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			jen.Line(),
			utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		),
		jen.Line(),
		utils.BuildSubTestWithoutContext(
			"with frontend routing path",
			jen.ID("s").Assign().AddressOf().ID("Service").Values(jen.ID(constants.LoggerVarName).MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.ID(utils.BuildFakeVarName("Dir")).Assign().Lit("."),
			jen.Line(),
			jen.List(jen.ID("hf"), jen.Err()).Assign().ID("s").Dot("StaticDir").Call(jen.ID(utils.BuildFakeVarName("Dir"))),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertNotNil(jen.ID("hf"), nil),
			jen.Line(),
			jen.List(jen.ID(constants.RequestVarName), jen.ID(constants.ResponseVarName)).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
			jen.ID(constants.RequestVarName).Dot("URL").Dot("Path").Equals().Lit("/login"),
			jen.ID("hf").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			jen.Line(),
			utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		),
		jen.Line(),
	}

	for _, typ := range proj.DataTypes {
		tpcn := typ.Name.PluralCommonName()

		block = append(block,
			jen.Line(), jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with frontend %s routing path", tpcn),
				jen.ID("s").Assign().AddressOf().ID("Service").Values(jen.ID(constants.LoggerVarName).MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID(utils.BuildFakeVarName("Dir")).Assign().Lit("."),
				jen.Line(),
				jen.List(jen.ID("hf"), jen.Err()).Assign().ID("s").Dot("StaticDir").Call(jen.ID(utils.BuildFakeVarName("Dir"))),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertNotNil(jen.ID("hf"), nil),
				jen.Line(),
				jen.List(jen.ID(constants.RequestVarName), jen.ID(constants.ResponseVarName)).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path").Equals().Lit(fmt.Sprintf("/%s/123", typ.Name.PluralRouteName())),
				jen.ID("hf").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
			),
		)
	}

	lines := []jen.Code{
		jen.Func().ID("TestService_StaticDir").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(block...), jen.Line(),
	}

	return lines
}
