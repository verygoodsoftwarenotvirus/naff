package client

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("client")

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

func buildVarDeclarationsOfDependentIDs(pkg *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.IDf("%sID", pt.Name.UnexportedVarName()).Op(":=").Add(utils.FakeUint64Func()))
	}
	lines = append(lines, jen.IDf("%sID", typ.Name.UnexportedVarName()).Op(":=").Add(utils.FakeUint64Func()))

	return lines
}

func buildVarDeclarationsOfDependentStructs(pkg *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lines = append(lines,
			jen.ID(pt.Name.UnexportedVarName()).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models", "v1"), pt.Name.Singular()).Valuesln(
				jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				func() jen.Code {
					if pt.BelongsToStruct != nil {
						return jen.IDf("BelongsTo%s", pt.BelongsToStruct.Singular()).Op(":").ID(pt.BelongsToStruct.UnexportedVarName()).Dot("ID")
					} else {
						return nil
					}
				}(),
			),
		)
	}

	lines = append(lines,
		jen.ID(typ.Name.UnexportedVarName()).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models", "v1"), typ.Name.Singular()).Valuesln(
			jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").ID(typ.BelongsToStruct.UnexportedVarName()).Dot("ID")
				} else {
					return nil
				}
			}(),
		),
	)

	return lines
}

func buildCreationVarDeclarationsOfDependentStructs(pkg *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		values := []jen.Code{
			jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
		}

		for _, field := range pt.Fields {
			if field.ValidForCreationInput {
				values = append(values, jen.ID(field.Name.Singular()).Op(":").Add(utils.FakeFuncForType(field.Type)()))
			}
		}

		values = append(values,
			func() jen.Code {
				if pt.BelongsToStruct != nil {
					return jen.IDf("BelongsTo%s", pt.BelongsToStruct.Singular()).Op(":").ID(pt.BelongsToStruct.UnexportedVarName()).Dot("ID")
				} else {
					return nil
				}
			}(),
		)

		lines = append(lines,
			jen.ID(pt.Name.UnexportedVarName()).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models", "v1"), pt.Name.Singular()).Valuesln(values...),
		)
	}

	values := []jen.Code{
		jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
	}

	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			values = append(values, jen.ID(field.Name.Singular()).Op(":").Add(utils.FakeFuncForType(field.Type)()))
		}
	}

	values = append(values,
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").ID(typ.BelongsToStruct.UnexportedVarName()).Dot("ID")
			} else {
				return nil
			}
		}(),
	)

	lines = append(lines,
		jen.ID(typ.Name.UnexportedVarName()).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models", "v1"), typ.Name.Singular()).Valuesln(values...),
	)

	return lines
}

func buildFormatStringForSingleInstanceRoute(pkg *models.Project, typ models.DataType) (path string) {
	modelRoute := "/api/v1/"
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func buildFormatStringForListRoute(pkg *models.Project, typ models.DataType) (path string) {
	modelRoute := "/api/v1/"
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s", typ.Name.PluralRouteName())

	return modelRoute
}

func buildFormatStringForSingleInstanceCreationRoute(pkg *models.Project, typ models.DataType) (path string) {
	modelRoute := "/api/v1/"
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += typ.Name.PluralRouteName()

	return modelRoute
}

func buildFormatCallArgsForSingleInstanceRoute(pkg *models.Project, typ models.DataType) (args []jen.Code) {
	callArgs := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		callArgs = append(callArgs, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
	}
	callArgs = append(callArgs, jen.ID(typ.Name.UnexportedVarName()).Dot("ID"))

	return callArgs
}

func buildFormatCallArgsForSingleInstanceRouteWithoutStructs(pkg *models.Project, typ models.DataType) (args []jen.Code) {
	callArgs := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		callArgs = append(callArgs, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	callArgs = append(callArgs, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	return callArgs
}

func buildFormatCallArgsForListRoute(pkg *models.Project, typ models.DataType) (args []jen.Code) {
	callArgs := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		callArgs = append(callArgs, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
	}

	return callArgs
}

func buildFormatCallArgsForSingleInstanceCreationRoute(pkg *models.Project, typ models.DataType) (args []jen.Code) {
	callArgs := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		callArgs = append(callArgs, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
	}

	return callArgs
}

func buildFormatCallArgsForSingleInstanceRouteThatIncludesItsOwnType(pkg *models.Project, typ models.DataType) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := pkg.FindOwnerTypeChain(typ)

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

	depVarDecls := buildVarDeclarationsOfDependentIDs(proj, typ)

	subtestLines := []jen.Code{
		utils.ExpectMethod("expectedMethod", "MethodHead"),
		jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
	}
	subtestLines = append(subtestLines, depVarDecls...)

	subtestLines = append(subtestLines,
		jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("c").Dot(fmt.Sprintf("Build%sExistsRequest", ts)).Call(
			buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, true)...,
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
					jen.IDf("%sID", typ.Name.UnexportedVarName()),
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
	uvn := typ.Name.UnexportedVarName()

	// routes
	var subtestLines []jen.Code
	actualCallArgs := []jen.Code{utils.CtxVar()}

	//for _, pt := range proj.FindOwnerTypeChain(typ) {
	//	actualCallArgs = append(actualCallArgs, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
	//
	//	if pt.BelongsToStruct != nil {
	//		subtestLines = append(subtestLines,
	//			jen.ID(pt.Name.UnexportedVarName()).Op(":=").Op("&").Qual(filepath.Join(proj.OutputPath, "models/v1"), pt.Name.Singular()).Valuesln(
	//				jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
	//				jen.IDf("BelongsTo%s", pt.BelongsToStruct.Singular()).Op(":").ID(pt.BelongsToStruct.UnexportedVarName()).Dot("ID"),
	//			),
	//		)
	//	} else {
	//		subtestLines = append(subtestLines,
	//			jen.IDf(pt.Name.UnexportedVarName()).Op(":=").Op("&").Qual(filepath.Join(proj.OutputPath, "models/v1"), pt.Name.Singular()).Values(jen.ID("ID").Op(":").Add(utils.FakeUint64Func())),
	//		)
	//	}
	//}
	actualCallArgs = append(actualCallArgs, jen.IDf("%sID", uvn))

	subtestLines = append(subtestLines,
		jen.IDf("%sID", uvn).Op(":=").Add(utils.FakeUint64Func()),
		jen.ID("expected").Op(":=").True(),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertTrue(
				jen.Qual("strings", "HasSuffix").Call(
					jen.ID("req").Dot("URL").Dot("String").Call(),
					jen.Qual("strconv", "Itoa").Call(
						jen.ID("int").Call(
							jen.IDf("%sID", uvn),
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
						buildFormatCallArgsForSingleInstanceRouteWithoutStructs(proj, typ)...,
					)...,
				),
				jen.Lit("expected and actual paths don't match"),
			),
			utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodHead"), nil),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
		),
		jen.Line(),
		jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("c").Dot(fmt.Sprintf("%sExists", ts)).Call(
			actualCallArgs...,
		),
		jen.Line(),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_%sExists", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", subtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildGetSomethingRequest(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	depVarDecls := buildVarDeclarationsOfDependentIDs(pkg, typ)

	subtestLines := []jen.Code{
		utils.ExpectMethod("expectedMethod", "MethodGet"),
		jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
	}
	subtestLines = append(subtestLines, depVarDecls...)

	subtestLines = append(subtestLines,
		jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(
			buildParamsForMethodThatHandlesAnInstanceWithIDs(pkg, typ, true)...,
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
					jen.IDf("%sID", typ.Name.UnexportedVarName()),
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

func buildTestV1Client_GetSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular
	uvn := typ.Name.UnexportedVarName()

	// routes
	var subtestLines []jen.Code
	actualCallArgs := []jen.Code{utils.CtxVar()}

	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		actualCallArgs = append(actualCallArgs, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))

		if pt.BelongsToStruct != nil {
			subtestLines = append(subtestLines,
				jen.ID(pt.Name.UnexportedVarName()).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), pt.Name.Singular()).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
					jen.IDf("BelongsTo%s", pt.BelongsToStruct.Singular()).Op(":").ID(pt.BelongsToStruct.UnexportedVarName()).Dot("ID"),
				),
			)
		} else {
			subtestLines = append(subtestLines,
				jen.ID(pt.Name.UnexportedVarName()).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), pt.Name.Singular()).Values(jen.ID("ID").Op(":").Add(utils.FakeUint64Func())),
			)
		}
	}
	actualCallArgs = append(actualCallArgs, jen.ID(uvn).Dot("ID"))

	subtestLines = append(subtestLines,
		jen.ID(uvn).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), ts).Valuesln(
			jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").ID(typ.BelongsToStruct.UnexportedVarName()).Dot("ID")
				} else {
					return nil
				}
			}(),
		),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertTrue(
				jen.Qual("strings", "HasSuffix").Call(
					jen.ID("req").Dot("URL").Dot("String").Call(),
					jen.Qual("strconv", "Itoa").Call(
						jen.ID("int").Call(
							jen.ID(uvn).Dot("ID"),
						),
					),
				),
				nil,
			),
			utils.AssertEqual(
				jen.ID("req").Dot("URL").Dot("Path"),
				jen.Qual("fmt", "Sprintf").Call(
					append(
						[]jen.Code{jen.Lit(buildFormatStringForSingleInstanceRoute(pkg, typ))},
						buildFormatCallArgsForSingleInstanceRoute(pkg, typ)...,
					)...,
				),
				jen.Lit("expected and actual paths don't match"),
			),
			utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodGet"), nil),
			utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID(uvn)), nil),
		),
		jen.Line(),
		jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("c").Dot(fmt.Sprintf("Get%s", ts)).Call(
			actualCallArgs...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(jen.ID(uvn), jen.ID("actual"), nil),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Get%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", subtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildGetListOfSomethingRequest(pkg *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural() // title plural

	structDecls := buildVarDeclarationsOfDependentStructs(pkg, typ)
	subtestLines := structDecls[:len(structDecls)-1]
	subtestLines = append(subtestLines,
		jen.ID(utils.FilterVarName).Op(":=").Call(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Call(jen.Nil()),
		utils.ExpectMethod("expectedMethod", "MethodGet"),
		jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(
			buildParamsForMethodThatRetrievesAListOfADataTypeFromStructs(pkg, typ, true)...,
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

func buildTestV1Client_GetListOfSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural()   // title plural
	ts := typ.Name.Singular() // title singular
	puvn := typ.Name.PluralUnexportedVarName()

	modelListRoute := buildFormatStringForListRoute(pkg, typ)

	var uriDec *jen.Statement
	urlFormatArgs := buildFormatCallArgsForListRoute(pkg, typ)
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

	structDecls := buildVarDeclarationsOfDependentStructs(pkg, typ)
	subtestLines := structDecls[:len(structDecls)-1]
	subtestLines = append(subtestLines,
		jen.ID(utils.FilterVarName).Op(":=").Add(utils.NilQueryFilter(pkg)),
		jen.Line(),
		jen.ID(puvn).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", ts)).Valuesln(
			jen.ID(tp).Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), ts).Valuesln(
				jen.Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
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
		jen.ID("c").Op(":=").ID("buildTestClient").Call(
			jen.ID("t"),
			jen.ID("ts"),
		),
		jen.List(jen.ID("actual"), jen.Err()).
			Op(":=").ID("c").Dot(fmt.Sprintf("Get%s", tp)).Call(
			buildParamsForMethodThatFetchesAListOfDataTypesFromStructs(pkg, typ, true)...,
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

func buildTestV1Client_BuildCreateSomethingRequest(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	cfs := []jen.Code{jen.ID("ID").Op(":").Add(utils.FakeUint64Func())}
	for _, field := range typ.Fields {
		cfs = append(cfs, jen.ID(field.Name.Singular()).Op(":").Add(utils.FakeFuncForType(field.Type)()))
	}

	structDecls := buildVarDeclarationsOfDependentStructs(pkg, typ)
	subtestLines := structDecls[:len(structDecls)-1]
	subtestLines = append(subtestLines,
		jen.Line(),
		utils.ExpectMethod("expectedMethod", "MethodPost"),
		jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("input").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", ts)).Valuesln(
			cfs[1:]...,
		),
		jen.ID("c").Op(":=").ID("buildTestClient").Call(
			jen.ID("t"),
			jen.ID("ts"),
		),
		jen.List(
			jen.ID("actual"),
			jen.Err(),
		).Op(":=").ID("c").Dot(fmt.Sprintf("BuildCreate%sRequest", ts)).Call(
			buildParamsForMethodThatCreatesADataTypeFromStructs(pkg, typ, true)...,
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

func buildTestV1Client_CreateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular
	uvn := typ.Name.UnexportedVarName()

	// routes
	modelListRoute := buildFormatStringForSingleInstanceCreationRoute(pkg, typ)

	var createCreationInputLines []jen.Code
	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			createCreationInputLines = append(createCreationInputLines, jen.ID(field.Name.Singular()).Op(":").ID(uvn).Dot(field.Name.Singular()))
		}
	}

	var uriDec *jen.Statement
	urlFormatArgs := buildFormatCallArgsForSingleInstanceCreationRoute(pkg, typ)
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

	subtestLines := buildCreationVarDeclarationsOfDependentStructs(pkg, typ)
	subtestLines = append(subtestLines,
		jen.ID("input").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", ts)).Valuesln(
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
			jen.Var().ID("x").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", ts)),
			utils.RequireNoError(jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot("Body")).Dot("Decode").Call(jen.Op("&").ID("x")), nil),
			utils.AssertEqual(jen.ID("input"), jen.ID("x"), nil),
			jen.Line(),
			utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID(uvn)), nil),
		),
		jen.Line(),
		jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("c").Dot(fmt.Sprintf("Create%s", ts)).Call(
			buildParamsForMethodThatCreatesADataTypeFromStructs(pkg, typ, true)...,
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

func buildTestV1Client_BuildUpdateSomethingRequest(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	actualParams := buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(pkg, typ)

	subtestLines := buildVarDeclarationsOfDependentStructs(pkg, typ)
	subtestLines = append(subtestLines,
		utils.ExpectMethod("expectedMethod", "MethodPut"),
		jen.Line(),
		jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("c").Dot(fmt.Sprintf("BuildUpdate%sRequest", ts)).Call(
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

func buildTestV1Client_UpdateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular
	uvn := typ.Name.UnexportedVarName()

	subtestLines := buildVarDeclarationsOfDependentStructs(pkg, typ)
	subtestLines = append(subtestLines,
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID("req").Dot("URL").Dot("Path"),
				jen.Qual("fmt", "Sprintf").Call(
					append(
						[]jen.Code{jen.Lit(buildFormatStringForSingleInstanceRoute(pkg, typ))},
						buildFormatCallArgsForSingleInstanceRouteThatIncludesItsOwnType(pkg, typ)...,
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
		jen.Err().Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot(fmt.Sprintf("Update%s", ts)).Call(
			buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(pkg, typ)...,
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

func buildTestV1Client_BuildArchiveSomethingRequest(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	subtestLines := []jen.Code{
		utils.ExpectMethod("expectedMethod", "MethodDelete"),
		jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
	}

	subtestLines = append(subtestLines, buildVarDeclarationsOfDependentStructs(pkg, typ)...)

	subtestLines = append(subtestLines,
		jen.Line(),
		jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructs(pkg, typ)...,
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

func buildTestV1Client_ArchiveSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	subtestLines := buildVarDeclarationsOfDependentStructs(pkg, typ)

	subtestLines = append(subtestLines,
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID("req").Dot("URL").Dot("Path"),
				jen.Qual("fmt", "Sprintf").Call(
					append(
						[]jen.Code{jen.Lit(buildFormatStringForSingleInstanceRoute(pkg, typ))},
						buildParamsForMethodThatHandlesAnInstanceWithStructs(pkg, typ)[1:]...,
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
		jen.Err().Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot(fmt.Sprintf("Archive%s", ts)).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructs(pkg, typ)...,
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
