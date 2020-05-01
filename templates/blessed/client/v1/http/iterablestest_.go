package client

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(buildTestV1Client_BuildSomethingExistsRequest(proj, typ)...)
	ret.Add(buildTestV1Client_SomethingExists(proj, typ)...)
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

func buildTestV1Client_BuildSomethingExistsRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	depVarDecls := typ.BuildDependentObjectsForHTTPClientBuildExistenceRequestMethodTest(proj)

	subtestLines := []jen.Code{
		utils.ExpectMethod("expectedMethod", "MethodHead"),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
	}
	subtestLines = append(subtestLines, depVarDecls...)

	subtestLines = append(subtestLines,
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("Build%sExistsRequest", ts)).Call(
			typ.BuildArgsForHTTPClientExistenceRequestBuildingMethodTest(proj)...,
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
				utils.FormatString("%d",
					jen.ID(utils.BuildFakeVarName(ts)).Dot("ID"),
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

func buildTestV1Client_SomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	// routes
	happyPathSubtestLines := typ.BuildDependentObjectsForHTTPClientExistenceMethodTest(proj)
	actualCallArgs := typ.BuildArgsForHTTPClientExistenceMethodTest(proj)
	//actualCallArgs := []jen.Code{utils.CtxVar(), jen.ID(utils.BuildFakeVarName(ts)).Dot("ID")}

	happyPathSubtestLines = append(happyPathSubtestLines,
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertTrue(
				jen.Qual("strings", "HasSuffix").Call(
					jen.ID(constants.RequestVarName).Dot("URL").Dot("String").Call(),
					jen.Qual("strconv", "Itoa").Call(
						jen.Int().Call(
							jen.ID(utils.BuildFakeVarName(ts)).Dot("ID"),
						),
					),
				),
				nil,
			),
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
				utils.FormatString(
					typ.BuildFormatStringForHTTPClientExistenceMethodTest(proj),
					typ.BuildFormatCallArgsForHTTPClientExistenceMethodTest(proj)...,
				),
				jen.Lit("expected and actual paths do not match"),
			),
			utils.AssertEqual(jen.ID(constants.RequestVarName).Dot("Method"), jen.Qual("net/http", "MethodHead"), nil),
			jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
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

	sadPathSubtestLines := typ.BuildDependentObjectsForHTTPClientExistenceMethodTest(proj)
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

	depVarDecls := typ.BuildDependentObjectsForHTTPClientBuildRetrievalRequestMethodTest(proj)

	subtestLines := []jen.Code{
		utils.ExpectMethod("expectedMethod", "MethodGet"),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
	}
	subtestLines = append(subtestLines, depVarDecls...)

	subtestLines = append(subtestLines,
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(
			typ.BuildArgsForHTTPClientRetrievalRequestBuilderMethodTest(proj)...,
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
				utils.FormatString("%d",
					jen.ID(utils.BuildFakeVarName(ts)).Dot("ID"),
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

	args := typ.BuildHTTPClientRetrievalTestCallArgs(proj)

	happyPathSubtestLines := append(
		typ.BuildDependentObjectsForHTTPClientRetrievalMethodTest(proj),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertTrue(
				jen.Qual("strings", "HasSuffix").Call(
					jen.ID(constants.RequestVarName).Dot("URL").Dot("String").Call(),
					jen.Qual("strconv", "Itoa").Call(
						jen.Int().Call(
							jen.ID(utils.BuildFakeVarName(ts)).Dot("ID"),
						),
					),
				),
				nil,
			),
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
				utils.FormatString(
					typ.BuildFormatStringForHTTPClientRetrievalMethodTest(proj),
					typ.BuildFormatCallArgsForHTTPClientRetrievalMethodTest(proj)...,
				),
				jen.Lit("expected and actual paths do not match"),
			),
			utils.AssertEqual(jen.ID(constants.RequestVarName).Dot("Method"), jen.Qual("net/http", "MethodGet"), nil),
			utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName(ts))), nil),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("Get%s", ts)).Call(args...),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(jen.ID(utils.BuildFakeVarName(ts)), jen.ID("actual"), nil),
	)

	invalidClientURLSubtestLines := append(
		typ.BuildDependentObjectsForHTTPClientRetrievalMethodTest(proj),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("Get%s", ts)).Call(args...),
		jen.Line(),
		utils.AssertNil(jen.ID("actual"), nil),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	)

	invalidResponseSubtestLines := append(
		typ.BuildDependentObjectsForHTTPClientRetrievalMethodTest(proj),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertTrue(
				jen.Qual("strings", "HasSuffix").Call(
					jen.ID(constants.RequestVarName).Dot("URL").Dot("String").Call(),
					jen.Qual("strconv", "Itoa").Call(
						jen.Int().Call(
							jen.ID(utils.BuildFakeVarName(ts)).Dot("ID"),
						),
					),
				),
				nil,
			),
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
				utils.FormatString(
					typ.BuildFormatStringForHTTPClientRetrievalMethodTest(proj),
					typ.BuildFormatCallArgsForHTTPClientRetrievalMethodTest(proj)...,
				),
				jen.Lit("expected and actual paths do not match"),
			),
			utils.AssertEqual(jen.ID(constants.RequestVarName).Dot("Method"), jen.Qual("net/http", "MethodGet"), nil),
			utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.Lit("BLAH")), nil),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("Get%s", ts)).Call(
			args...,
		),
		jen.Line(),
		utils.AssertNil(jen.ID("actual"), nil),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Get%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid client URL", invalidClientURLSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid response", invalidResponseSubtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildGetListOfSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural() // title plural

	subtestLines := append(
		typ.BuildDependentObjectsForHTTPClientBuildListRetrievalRequestMethodTest(proj),
		jen.ID(constants.FilterVarName).Assign().Call(jen.PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter")).Call(jen.Nil()),
		utils.ExpectMethod("expectedMethod", "MethodGet"),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(
			typ.BuildCallArgsForHTTPClientListRetrievalRequestBuildingMethodTest(proj)...,
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

	modelListRoute := typ.BuildFormatStringForHTTPClientListMethodTest(proj)

	var uriDec *jen.Statement
	urlFormatArgs := typ.BuildFormatCallArgsForHTTPClientListMethodTest(proj)
	if len(urlFormatArgs) > 0 {
		uriDec = utils.FormatString(
			modelListRoute,
			urlFormatArgs...,
		)
	} else {
		uriDec = jen.Lit(modelListRoute)
	}

	structDecls := typ.BuildDependentObjectsForHTTPClientListRetrievalTest(proj)
	happyPathSubtestLines := append(
		structDecls,
		jen.ID(constants.FilterVarName).Assign().Add(utils.NilQueryFilter(proj)),
		jen.Line(),
		jen.IDf("example%sList", ts).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", ts)).Call(),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
				uriDec,
				jen.Lit("expected and actual paths do not match"),
			),
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("Method"),
				jen.Qual("net/http", "MethodGet"),
				nil,
			),
			utils.RequireNoError(
				jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.IDf("example%sList", ts)),
				nil,
			),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).
			Op(":=").ID("c").Dot(fmt.Sprintf("Get%s", tp)).Call(
			typ.BuildCallArgsForHTTPClientListRetrievalMethodTest(proj)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(jen.IDf("example%sList", ts), jen.ID("actual"), nil),
	)

	invalidClientURLSubtestLines := append(
		structDecls,
		jen.ID(constants.FilterVarName).Assign().Add(utils.NilQueryFilter(proj)),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
		jen.List(jen.ID("actual"), jen.Err()).
			Op(":=").ID("c").Dot(fmt.Sprintf("Get%s", tp)).Call(
			typ.BuildCallArgsForHTTPClientListRetrievalMethodTest(proj)...,
		),
		jen.Line(),
		utils.AssertNil(jen.ID("actual"), nil),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	)

	invalidResponseSubtestLines := append(
		structDecls,
		jen.ID(constants.FilterVarName).Assign().Add(utils.NilQueryFilter(proj)),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
				uriDec,
				jen.Lit("expected and actual paths do not match"),
			),
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("Method"),
				jen.Qual("net/http", "MethodGet"),
				nil,
			),
			utils.RequireNoError(
				jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.Lit("BLAH")),
				nil,
			),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).
			Op(":=").ID("c").Dot(fmt.Sprintf("Get%s", tp)).Call(
			typ.BuildCallArgsForHTTPClientListRetrievalMethodTest(proj)...,
		),
		jen.Line(),
		utils.AssertNil(jen.ID("actual"), nil),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Get%s", tp)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid client URL", invalidClientURLSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid response", invalidResponseSubtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildCreateSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	cfs := []jen.Code{jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())}
	for _, field := range typ.Fields {
		cfs = append(cfs, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type, field.Pointer)()))
	}

	subtestLines := append(
		typ.BuildDependentObjectsForHTTPClientBuildCreationRequestMethodTest(proj),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", ts, ts)).Call(jen.ID(utils.BuildFakeVarName(ts))),
		jen.Line(),
		utils.ExpectMethod("expectedMethod", "MethodPost"),
		jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildCreate%sRequest", ts)).Call(
			typ.BuildHTTPClientCreationRequestBuildingMethodArgsForTest(proj)...,
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

	// routes
	modelListRoute := typ.BuildFormatStringForHTTPClientCreateMethodTest(proj)

	var uriDec *jen.Statement
	urlFormatArgs := typ.BuildFormatCallArgsForHTTPClientCreationMethodTest(proj)
	if len(urlFormatArgs) > 0 {
		uriDec = utils.FormatString(
			modelListRoute,
			urlFormatArgs...,
		)
	} else {
		uriDec = jen.Lit(modelListRoute)
	}

	happyPathSubtestLines := append(
		typ.BuildDependentObjectsForHTTPClientCreationMethodTest(proj),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", ts, ts)).Call(jen.ID(utils.BuildFakeVarName(ts))),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
				uriDec,
				jen.Lit("expected and actual paths do not match"),
			),
			utils.AssertEqual(jen.ID(constants.RequestVarName).Dot("Method"), jen.Qual("net/http", "MethodPost"), nil),
			jen.Line(),
			jen.Var().ID("x").PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", ts)),
			utils.RequireNoError(jen.Qual("encoding/json", "NewDecoder").Call(jen.ID(constants.RequestVarName).Dot("Body")).Dot("Decode").Call(jen.AddressOf().ID("x")), nil),
			jen.Line(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName("Input")).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().Zero()
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToUser {
					return jen.ID(utils.BuildFakeVarName("Input")).Dot(constants.UserOwnershipFieldName).Equals().Zero()
				}
				return jen.Null()
			}(),
			utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Input")), jen.ID("x"), nil),
			jen.Line(),
			utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName(ts))), nil),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("Create%s", ts)).Call(
			typ.BuildHTTPClientCreationMethodArgsForTest(proj)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(jen.ID(utils.BuildFakeVarName(ts)), jen.ID("actual"), nil),
	)

	invalidClientURLSubtestLines := append(
		typ.BuildDependentObjectsForHTTPClientCreationMethodTest(proj),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", ts, ts)).Call(jen.ID(utils.BuildFakeVarName(ts))),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("Create%s", ts)).Call(
			typ.BuildHTTPClientCreationMethodArgsForTest(proj)...,
		),
		jen.Line(),
		utils.AssertNil(jen.ID("actual"), nil),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Create%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid client URL", invalidClientURLSubtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildUpdateSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular() // title singular

	actualParams := typ.BuildCallArgsForHTTPClientUpdateRequestBuildingMethodTest(proj)

	subtestLines := typ.BuildDependentObjectsForHTTPClientBuildUpdateRequestMethodTest(proj)
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

	subtestLines := typ.BuildDependentObjectsForHTTPClientUpdateMethodTest(proj)

	happyPathSubtestLines := append(
		subtestLines,
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
				utils.FormatString(
					typ.BuildFormatStringForHTTPClientUpdateMethodTest(proj),
					typ.BuildFormatCallArgsForHTTPClientUpdateTest(proj)...,
				),
				jen.Lit("expected and actual paths do not match"),
			),
			utils.AssertEqual(jen.ID(constants.RequestVarName).Dot("Method"), jen.Qual("net/http", "MethodPut"), nil),
			utils.AssertNoError(
				jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName(ts))),
				nil,
			),
		),
		jen.Line(),
		jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot(fmt.Sprintf("Update%s", ts)).Call(
			typ.BuildCallArgsForHTTPClientUpdateMethodTest(proj)...,
		),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
	)

	withInvalidClientURLSubtestLines := append(
		subtestLines,
		jen.Line(),
		jen.Err().Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot(fmt.Sprintf("Update%s", ts)).Call(
			typ.BuildCallArgsForHTTPClientUpdateMethodTest(proj)...,
		),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Update%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid client URL", withInvalidClientURLSubtestLines...),
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

	subtestLines = append(subtestLines, typ.BuildDependentObjectsForHTTPClientBuildArchiveRequestMethodTest(proj)...)

	subtestLines = append(subtestLines,
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
			typ.BuildArgsForHTTPClientArchiveRequestBuildingMethodTest(proj)...,
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.RequireNotNil(jen.ID("actual").Dot("URL"), nil),
		utils.AssertTrue(
			jen.Qual("strings", "HasSuffix").Call(
				jen.ID("actual").Dot("URL").Dot("String").Call(),
				utils.FormatString("%d",
					jen.ID(utils.BuildFakeVarName(ts)).Dot("ID"),
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

	happyPathSubtestLines := append(
		typ.BuildDependentObjectsForHTTPClientArchiveMethodTest(proj),
		jen.Line(),
		utils.BuildTestServer(
			"ts",
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
				utils.FormatString(
					typ.BuildFormatStringForHTTPClientArchiveMethodTest(proj),
					typ.BuildArgsForHTTPClientArchiveMethodTestURLFormatCall(proj)...,
				),
				jen.Lit("expected and actual paths do not match"),
			),
			utils.AssertEqual(
				jen.ID(constants.RequestVarName).Dot("Method"),
				jen.Qual("net/http", "MethodDelete"),
				nil,
			),
			utils.WriteHeader("StatusOK"),
		),
		jen.Line(),
		jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot(fmt.Sprintf("Archive%s", ts)).Call(
			typ.BuildArgsForHTTPClientArchiveMethodTest(proj)...,
		),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
	)

	withInvalidClientURLSubtestLines := append(
		typ.BuildDependentObjectsForHTTPClientArchiveMethodTest(proj),
		jen.Line(),
		jen.Err().Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot(fmt.Sprintf("Archive%s", ts)).Call(
			typ.BuildArgsForHTTPClientArchiveMethodTest(proj)...,
		),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	)

	lines := []jen.Code{
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Archive%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid client URL", withInvalidClientURLSubtestLines...),
		),
		jen.Line(),
	}

	return lines
}
