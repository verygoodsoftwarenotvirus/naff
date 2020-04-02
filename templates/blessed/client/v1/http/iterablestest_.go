package client

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(buildTestV1Client_BuildItemExistsRequest(proj, typ)...)
	ret.Add(buildTestV1Client_ItemExists(proj, typ)...)
	ret.Add(buildTestV1Client_BuildGetSomethingRequest(proj, typ)...)
	ret.Add(buildTestV1Client_GetSomething(proj, typ)...)
	ret.Add(buildTestV1Client_BuildGetListOfSomethingRequest(proj, typ)...)
	ret.Add(buildTestV1Client_GetListOfSomething(proj, typ)...)
	ret.Add(buildTestV1Client_BuildCreateSomethingRequest(proj, typ)...)
	ret.Add(buildTestV1Client_CreateSomething(proj, typ)...)
	ret.Add(buildTestV1Client_BuildUpdateSomethingRequest(proj, typ)...)
	ret.Add(buildTestV1Client_UpdateSomething(proj, typ)...)
	ret.Add(buildTestV1Client_BuildArchiveSomethingRequest(proj, typ)...)
	ret.Add(buildTestV1Client_ArchiveSomething(proj, typ)...)

	return ret
}

func buildVarDeclarationsOfDependentStructs(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		pts := pt.Name.Singular()
		lines = append(lines, jen.IDf("example%s", pts).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
	}
	lines = append(lines, jen.IDf("example%s", typ.Name.Singular()).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	return lines
}

func buildCreationVarDeclarationsOfDependentStructs(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		values := []jen.Code{
			jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
		}

		for _, field := range pt.Fields {
			if field.ValidForCreationInput {
				values = append(values, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type)()))
			}
		}

		values = append(values,
			func() jen.Code {
				if pt.BelongsToStruct != nil {
					return jen.IDf("BelongsTo%s", pt.BelongsToStruct.Singular()).MapAssign().ID(pt.BelongsToStruct.UnexportedVarName()).Dot("ID")
				} else {
					return nil
				}
			}(),
		)

		lines = append(lines,
			jen.ID(pt.Name.UnexportedVarName()).Assign().VarPointer().Qual(proj.ModelsV1Package(), pt.Name.Singular()).Valuesln(values...),
		)
	}

	values := []jen.Code{
		jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
	}

	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			values = append(values, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type)()))
		}
	}

	values = append(values,
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().ID(typ.BelongsToStruct.UnexportedVarName()).Dot("ID")
			} else {
				return nil
			}
		}(),
	)

	lines = append(lines,
		jen.ID(typ.Name.UnexportedVarName()).Assign().VarPointer().Qual(proj.ModelsV1Package(), typ.Name.Singular()).Valuesln(values...),
	)

	return lines
}

func buildFormatStringForSingleInstanceRoute(proj *models.Project, typ models.DataType) (path string) {
	modelRoute := "/api/v1/"
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func buildFormatStringForListRoute(proj *models.Project, typ models.DataType) (path string) {
	modelRoute := "/api/v1/"
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s", typ.Name.PluralRouteName())

	return modelRoute
}

func buildFormatStringForSingleInstanceCreationRoute(proj *models.Project, typ models.DataType) (path string) {
	modelRoute := "/api/v1/"
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += typ.Name.PluralRouteName()

	return modelRoute
}

func buildFormatCallArgsForSingleInstanceRoute(proj *models.Project, typ models.DataType) (args []jen.Code) {
	callArgs := []jen.Code{}
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		callArgs = append(callArgs, jen.IDf("example%s", pt.Name.Singular()).Dot("ID"))
	}
	callArgs = append(callArgs, jen.IDf("example%s", typ.Name.Singular()).Dot("ID"))

	return callArgs
}

func buildFormatCallArgsForListRoute(proj *models.Project, typ models.DataType) (args []jen.Code) {
	callArgs := []jen.Code{}
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		callArgs = append(callArgs, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
	}

	return callArgs
}

func buildFormatCallArgsForSingleInstanceCreationRoute(proj *models.Project, typ models.DataType) (args []jen.Code) {
	callArgs := []jen.Code{}
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		callArgs = append(callArgs, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
	}

	return callArgs
}

func buildFormatCallArgsForSingleInstanceRouteThatIncludesItsOwnType(proj *models.Project, typ models.DataType) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)

	for i, pt := range owners {
		if typ.BelongsToStruct != nil && i == len(owners)-1 {
			callArgs = append(callArgs, jen.ID(typ.Name.UnexportedVarName()).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		} else {
			callArgs = append(callArgs, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
		}
	}
	callArgs = append(callArgs, jen.ID(typ.Name.UnexportedVarName()).Dot("ID"))

	return callArgs
}

func buildTestV1Client_BuildItemExistsRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	depVarDecls := buildVarDeclarationsOfDependentStructs(proj, typ)

	subtestLines := []jen.Code{
		utils.ExpectMethod("expectedMethod", "MethodHead"),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
	}
	subtestLines = append(subtestLines, depVarDecls...)

	subtestLines = append(subtestLines,
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("Build%sExistsRequest", ts)).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructs(proj, typ)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(
			jen.Err(),
			jen.Lit("no error should be returned"),
		),
		utils.AssertTrue(
			jen.Qual("strings", "HasSuffix").Call(
				jen.ID("actual").Dot("URL").Dot("String").Call(),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%d"),
					jen.IDf("example%s", ts).Dot("ID"),
				),
			),
			nil,
		),
		utils.AssertEqual(
			jen.ID("actual").Dot("Method"),
			jen.ID("expectedMethod"),
			jen.Lit("request should be a %s request"),
			jen.ID("expectedMethod"),
		),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Build%sExistsRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", subtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_ItemExists(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	// routes
	happyPathSubtestLines := buildVarDeclarationsOfDependentStructs(proj, typ)
	actualCallArgs := []jen.Code{utils.CtxVar(), jen.IDf("example%s", ts).Dot("ID")}

	happyPathSubtestLines = append(happyPathSubtestLines,
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertTrue(
				jen.Qual("strings", "HasSuffix").Call(
					jen.ID("req").Dot("URL").Dot("String").Call(),
					jen.Qual("strconv", "Itoa").Call(
						jen.ID("int").Call(
							jen.IDf("example%s", ts).Dot("ID"),
						),
					),
				),
				nil,
			),
			utils.AssertEqual(
				jen.ID("req").Dot("URL").Dot("Path"),
				jen.Qual("fmt", "Sprintf").Call(
					append(
						[]jen.Code{jen.Lit(buildFormatStringForSingleInstanceRoute(proj, typ))},
						buildFormatCallArgsForSingleInstanceRoute(proj, typ)...,
					)...,
				),
				jen.Lit("expected and actual paths don't match"),
			),
			utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodHead"), nil),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("%sExists", ts)).Call(
			actualCallArgs...,
		),
		jen.Line(),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertTrue(jen.ID("actual"), nil),
	)

	sadPathSubtestLines := buildVarDeclarationsOfDependentStructs(proj, typ)
	sadPathSubtestLines = append(sadPathSubtestLines,
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("%sExists", ts)).Call(
			actualCallArgs...,
		),
		jen.Line(),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
		utils.AssertFalse(jen.ID("actual"), nil),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_%sExists", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with erroneous response", sadPathSubtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildGetSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	depVarDecls := buildVarDeclarationsOfDependentStructs(proj, typ)

	subtestLines := []jen.Code{
		utils.ExpectMethod("expectedMethod", "MethodGet"),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
	}
	subtestLines = append(subtestLines, depVarDecls...)

	subtestLines = append(subtestLines,
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructs(proj, typ)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(
			jen.Err(),
			jen.Lit("no error should be returned"),
		),
		utils.AssertTrue(
			jen.Qual("strings", "HasSuffix").Call(
				jen.ID("actual").Dot("URL").Dot("String").Call(),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%d"),
					jen.IDf("example%s", ts).Dot("ID"),
				),
			),
			nil,
		),
		utils.AssertEqual(
			jen.ID("actual").Dot("Method"),
			jen.ID("expectedMethod"),
			jen.Lit("request should be a %s request"),
			jen.ID("expectedMethod"),
		),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildGet%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", subtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_GetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	// routes
	var happyPathSubtestLines []jen.Code
	happyPathSubtestLines = append(happyPathSubtestLines, buildVarDeclarationsOfDependentStructs(proj, typ)...)
	happyPathSubtestLines = append(happyPathSubtestLines,
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertTrue(
				jen.Qual("strings", "HasSuffix").Call(
					jen.ID("req").Dot("URL").Dot("String").Call(),
					jen.Qual("strconv", "Itoa").Call(
						jen.ID("int").Call(
							jen.IDf("example%s", ts).Dot("ID"),
						),
					),
				),
				nil,
			),
			utils.AssertEqual(
				jen.ID("req").Dot("URL").Dot("Path"),
				jen.Qual("fmt", "Sprintf").Call(
					append(
						[]jen.Code{jen.Lit(buildFormatStringForSingleInstanceRoute(proj, typ))},
						buildFormatCallArgsForSingleInstanceRoute(proj, typ)...,
					)...,
				),
				jen.Lit("expected and actual paths don't match"),
			),
			utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodGet"), nil),
			utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.IDf("example%s", ts)), nil),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("Get%s", ts)).Call(
			typ.BuildGetSomethingArgs(proj)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(jen.IDf("example%s", ts), jen.ID("actual"), nil),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Get%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildGetListOfSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural() // title plural

	structDecls := buildVarDeclarationsOfDependentStructs(proj, typ)
	subtestLines := structDecls[:len(structDecls)-1]
	subtestLines = append(subtestLines,
		jen.ID(utils.FilterVarName).Assign().Call(jen.Op("*").Qual(proj.ModelsV1Package(), "QueryFilter")).Call(jen.Nil()),
		utils.ExpectMethod("expectedMethod", "MethodGet"),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(
			buildParamsForMethodThatRetrievesAListOfADataTypeFromStructs(proj, typ, true)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(
			jen.ID("actual").Dot("Method"),
			jen.ID("expectedMethod"),
			jen.Lit("request should be a %s request"),
			jen.ID("expectedMethod"),
		),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildGet%sRequest", tp)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", subtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_GetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural()   // title plural
	ts := typ.Name.Singular() // title singular
	puvn := typ.Name.PluralUnexportedVarName()

	modelListRoute := buildFormatStringForListRoute(proj, typ)

	var uriDec *jen.Statement
	urlFormatArgs := buildFormatCallArgsForListRoute(proj, typ)
	if len(urlFormatArgs) > 0 {
		uriDec = jen.Qual("fmt", "Sprintf").Call(
			append(
				[]jen.Code{jen.Lit(modelListRoute)},
				urlFormatArgs...,
			)...,
		)
	} else {
		uriDec = jen.Lit(modelListRoute)
	}

	structDecls := buildVarDeclarationsOfDependentStructs(proj, typ)
	subtestLines := structDecls[:len(structDecls)-1]
	subtestLines = append(subtestLines,
		jen.ID(utils.FilterVarName).Assign().Add(utils.NilQueryFilter(proj)),
		jen.Line(),
		jen.ID(puvn).Assign().VarPointer().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", ts)).Valuesln(
			jen.ID(tp).MapAssign().Index().Qual(proj.ModelsV1Package(), ts).Valuesln(
				jen.Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
				),
			),
		),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID("req").Dot("URL").Dot("Path"),
				uriDec,
				jen.Lit("expected and actual paths don't match"),
			),
			utils.AssertEqual(
				jen.ID("req").Dot("Method"),
				jen.Qual("net/http", "MethodGet"),
				nil,
			),
			utils.RequireNoError(
				jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID(puvn)),
				nil,
			),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(
			jen.ID("t"),
			jen.ID("ts"),
		),
		jen.List(jen.ID("actual"), jen.Err()).
			Op(":=").ID("c").Dot(fmt.Sprintf("Get%s", tp)).Call(
			buildParamsForMethodThatFetchesAListOfDataTypesFromStructs(proj, typ, true)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(
			jen.Err(),
			jen.Lit("no error should be returned"),
		),
		utils.AssertEqual(jen.ID(puvn), jen.ID("actual"), nil),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Get%s", tp)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", subtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildCreateSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	cfs := []jen.Code{jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())}
	for _, field := range typ.Fields {
		cfs = append(cfs, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type)()))
	}

	structDecls := buildVarDeclarationsOfDependentStructs(proj, typ)
	subtestLines := structDecls[:len(structDecls)-1]
	subtestLines = append(subtestLines,
		jen.Line(),
		utils.ExpectMethod("expectedMethod", "MethodPost"),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("input").Assign().VarPointer().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", ts)).Valuesln(
			cfs[1:]...,
		),
		jen.ID("c").Assign().ID("buildTestClient").Call(
			jen.ID("t"),
			jen.ID("ts"),
		),
		jen.List(
			jen.ID("actual"),
			jen.Err(),
		).Assign().ID("c").Dot(fmt.Sprintf("BuildCreate%sRequest", ts)).Call(
			buildParamsForMethodThatCreatesADataTypeFromStructs(proj, typ, true)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(
			jen.Err(),
			jen.Lit("no error should be returned"),
		),
		utils.AssertEqual(
			jen.ID("actual").Dot("Method"),
			jen.ID("expectedMethod"),
			jen.Lit("request should be a %s request"),
			jen.ID("expectedMethod"),
		),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildCreate%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", subtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_CreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular
	uvn := typ.Name.UnexportedVarName()

	// routes
	modelListRoute := buildFormatStringForSingleInstanceCreationRoute(proj, typ)

	var createCreationInputLines []jen.Code
	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			createCreationInputLines = append(createCreationInputLines, jen.ID(field.Name.Singular()).MapAssign().ID(uvn).Dot(field.Name.Singular()))
		}
	}

	var uriDec *jen.Statement
	urlFormatArgs := buildFormatCallArgsForSingleInstanceCreationRoute(proj, typ)
	if len(urlFormatArgs) > 0 {
		uriDec = jen.Qual("fmt", "Sprintf").Call(
			append(
				[]jen.Code{jen.Lit(modelListRoute)},
				urlFormatArgs...,
			)...,
		)
	} else {
		uriDec = jen.Lit(modelListRoute)
	}

	subtestLines := buildCreationVarDeclarationsOfDependentStructs(proj, typ)
	subtestLines = append(subtestLines,
		jen.ID("input").Assign().VarPointer().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", ts)).Valuesln(
			createCreationInputLines...,
		),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID("req").Dot("URL").Dot("Path"),
				uriDec,
				jen.Lit("expected and actual paths don't match"),
			),
			utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodPost"), nil),
			jen.Line(),
			jen.Var().ID("x").Op("*").Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", ts)),
			utils.RequireNoError(jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot("Body")).Dot("Decode").Call(jen.VarPointer().ID("x")), nil),
			utils.AssertEqual(jen.ID("input"), jen.ID("x"), nil),
			jen.Line(),
			utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID(uvn)), nil),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("Create%s", ts)).Call(
			buildParamsForMethodThatCreatesADataTypeFromStructs(proj, typ, true)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(jen.ID(uvn), jen.ID("actual"), nil),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Create%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				subtestLines...,
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildUpdateSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	actualParams := buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj, typ)

	subtestLines := buildVarDeclarationsOfDependentStructs(proj, typ)
	subtestLines = append(subtestLines,
		utils.ExpectMethod("expectedMethod", "MethodPut"),
		jen.Line(),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildUpdate%sRequest", ts)).Call(
			actualParams...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(
			jen.ID("actual").Dot("Method"),
			jen.ID("expectedMethod"),
			jen.Lit("request should be a %s request"),
			jen.ID("expectedMethod"),
		),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildUpdate%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", subtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_UpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular
	uvn := typ.Name.UnexportedVarName()

	subtestLines := buildVarDeclarationsOfDependentStructs(proj, typ)
	subtestLines = append(subtestLines,
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID("req").Dot("URL").Dot("Path"),
				jen.Qual("fmt", "Sprintf").Call(
					append(
						[]jen.Code{jen.Lit(buildFormatStringForSingleInstanceRoute(proj, typ))},
						buildFormatCallArgsForSingleInstanceRouteThatIncludesItsOwnType(proj, typ)...,
					)...,
				),
				jen.Lit("expected and actual paths don't match"),
			),
			utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodPut"), nil),
			utils.AssertNoError(
				jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID(uvn)),
				nil,
			),
		),
		jen.Line(),
		jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot(fmt.Sprintf("Update%s", ts)).Call(
			buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj, typ)...,
		),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Update%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", subtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildArchiveSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	subtestLines := []jen.Code{
		utils.ExpectMethod("expectedMethod", "MethodDelete"),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
	}

	subtestLines = append(subtestLines, buildVarDeclarationsOfDependentStructs(proj, typ)...)

	subtestLines = append(subtestLines,
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructs(proj, typ)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.RequireNotNil(jen.ID("actual").Dot("URL"), nil),
		utils.AssertTrue(
			jen.Qual("strings", "HasSuffix").Call(
				jen.ID("actual").Dot("URL").Dot("String").Call(),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%d"),
					jen.ID(typ.Name.UnexportedVarName()).Dot("ID"),
				),
			),
			nil,
		),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(
			jen.ID("actual").Dot("Method"),
			jen.ID("expectedMethod"),
			jen.Lit("request should be a %s request"),
			jen.ID("expectedMethod"),
		),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildArchive%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				subtestLines...,
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_ArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	subtestLines := buildVarDeclarationsOfDependentStructs(proj, typ)

	subtestLines = append(subtestLines,
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID("req").Dot("URL").Dot("Path"),
				jen.Qual("fmt", "Sprintf").Call(
					append(
						[]jen.Code{jen.Lit(buildFormatStringForSingleInstanceRoute(proj, typ))},
						buildParamsForMethodThatHandlesAnInstanceWithStructs(proj, typ)[1:]...,
					)...,
				),
				jen.Lit("expected and actual paths don't match"),
			),
			utils.AssertEqual(
				jen.ID("req").Dot("Method"),
				jen.Qual("net/http", "MethodDelete"),
				nil,
			),
			utils.WriteHeader("StatusOK"),
		),
		jen.Line(),
		jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot(fmt.Sprintf("Archive%s", ts)).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructs(proj, typ)...,
		),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Archive%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				subtestLines...,
			),
		),
		jen.Line(),
	}

	return lines
}
