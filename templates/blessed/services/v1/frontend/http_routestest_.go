package frontend

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("buildRequest").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Params(jen.ParamPointer().Qual("net/http", "Request")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodGet"),
				jen.Lit("https://verygoodsoftwarenotvirus.ru"),
				jen.Nil(),
			),
			jen.Line(),
			utils.RequireNotNil(jen.ID("req"), nil),
			utils.AssertNoError(jen.Err(), nil),
			jen.Return().ID("req"),
		),
		jen.Line(),
	)

	ret.Add(buildTestService_StaticDir(proj)...)

	ret.Add(
		jen.Func().ID("TestService_Routes").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				utils.AssertNotNil(jen.Parens(jen.VarPointer().ID("Service").Values()).Dot("Routes").Call(), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_buildStaticFileServer").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("s").Assign().VarPointer().ID("Service").Valuesln(
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
			jen.ID("s").Assign().VarPointer().ID("Service").Values(jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			jen.List(jen.ID("cwd"), jen.Err()).Assign().Qual("os", "Getwd").Call(),
			utils.RequireNoError(jen.Err(), nil),
			jen.Line(),
			jen.List(jen.ID("hf"), jen.Err()).Assign().ID("s").Dot("StaticDir").Call(jen.ID("cwd")),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertNotNil(jen.ID("hf"), nil),
			jen.Line(),
			jen.List(jen.ID("req"), jen.ID("res")).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
			jen.ID("req").Dot("URL").Dot("Path").Equals().Lit("/http_routes_test.go"),
			jen.ID("hf").Call(jen.ID("res"), jen.ID("req")),
			jen.Line(),
			utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
		),
		jen.Line(),
		utils.BuildSubTestWithoutContext(
			"with frontend routing path",
			jen.ID("s").Assign().VarPointer().ID("Service").Values(jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.ID("exampleDir").Assign().Lit("."),
			jen.Line(),
			jen.List(jen.ID("hf"), jen.Err()).Assign().ID("s").Dot("StaticDir").Call(jen.ID("exampleDir")),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertNotNil(jen.ID("hf"), nil),
			jen.Line(),
			jen.List(jen.ID("req"), jen.ID("res")).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
			jen.ID("req").Dot("URL").Dot("Path").Equals().Lit("/login"),
			jen.ID("hf").Call(jen.ID("res"), jen.ID("req")),
			jen.Line(),
			utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
		),
		jen.Line(),
	}

	for _, typ := range proj.DataTypes {
		tpcn := typ.Name.PluralCommonName()

		block = append(block,
			jen.Line(), jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with frontend %s routing path", tpcn),
				jen.ID("s").Assign().VarPointer().ID("Service").Values(jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("exampleDir").Assign().Lit("."),
				jen.Line(),
				jen.List(jen.ID("hf"), jen.Err()).Assign().ID("s").Dot("StaticDir").Call(jen.ID("exampleDir")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertNotNil(jen.ID("hf"), nil),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("res")).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
				jen.ID("req").Dot("URL").Dot("Path").Equals().Lit(fmt.Sprintf("/%s/123", typ.Name.PluralRouteName())),
				jen.ID("hf").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
		)
	}

	lines := []jen.Code{
		jen.Func().ID("TestService_StaticDir").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(block...), jen.Line(),
	}

	return lines
}
