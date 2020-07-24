package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code)

	code.Add(buildTestServiceListFuncDecl(proj, typ)...)
	if typ.SearchEnabled {
		code.Add(buildTestServiceSearchFuncDecl(proj, typ)...)
	}
	code.Add(buildTestServiceCreateFuncDecl(proj, typ)...)
	code.Add(buildTestServiceExistenceFuncDecl(proj, typ)...)
	code.Add(buildTestServiceReadFuncDecl(proj, typ)...)
	code.Add(buildTestServiceUpdateFuncDecl(proj, typ)...)
	code.Add(buildTestServiceArchiveFuncDecl(proj, typ)...)

	return code
}

func includeOwnerFetchers(proj *models.Project, typ models.DataType) []jen.Code {
	out := []jen.Code{jen.Line()}

	if typ.OwnedByAUserAtSomeLevel(proj) {
		out = append(out,
			jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher"),
		)
	}
	for _, ot := range proj.FindOwnerTypeChain(typ) {
		btsuvn := ot.Name.UnexportedVarName()
		out = append(out, jen.ID("s").Dotf("%sIDFetcher", btsuvn).Equals().IDf("%sIDFetcher", btsuvn))
	}

	out = append(out, jen.Line())

	return out
}

func buildRelevantIDFetchers(proj *models.Project, typ models.DataType) []jen.Code {
	out := []jen.Code{}
	if typ.OwnedByAUserAtSomeLevel(proj) {
		out = append(out,
			utils.BuildFakeVar(proj, "User"),
			jen.ID("userIDFetcher").Assign().Func().Params(jen.Underscore().PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
				jen.Return().ID(utils.BuildFakeVarName("User")).Dot("ID"),
			),
			jen.Line(),
		)
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		out = append(out,
			utils.BuildFakeVar(proj, ot.Name.Singular()),
			func() jen.Code {
				if ot.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName(ot.Name.Singular())).Dotf("BelongsTo%s", ot.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(ot.BelongsToStruct.Singular())).Dot("ID")
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if ot.BelongsToUser {
					return jen.ID(utils.BuildFakeVarName(ot.Name.Singular())).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
				}
				return jen.Null()
			}(),
			jen.IDf("%sIDFetcher", ot.Name.UnexportedVarName()).Assign().Func().Params(jen.Underscore().PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
				jen.Return().ID(utils.BuildFakeVarName(ot.Name.Singular())).Dot("ID"),
			),
			jen.Line(),
		)
	}

	return out
}

func setupDataManagersForCreation(proj *models.Project, typ models.DataType, actualCallArgs, returnValues []jen.Code, indexToReturnFalse int, returnErr bool) (out []jen.Code) {
	owners := proj.FindOwnerTypeChain(typ)
	for i, ot := range owners {
		sn := ot.Name.Singular()
		uvn := ot.Name.UnexportedVarName()
		dataManagerVarName := fmt.Sprintf("%sDataManager", ot.Name.UnexportedVarName())

		if i > indexToReturnFalse && indexToReturnFalse > -1 {
			return
		}

		returnVals := []jen.Code{
			jen.True(),
			jen.Nil(),
		}
		if i == indexToReturnFalse && !returnErr {
			returnVals[0] = jen.False()
		} else if i == indexToReturnFalse && returnErr {
			returnVals[1] = constants.ObligatoryError()
		}

		callArgs := append(
			[]jen.Code{jen.Litf("%sExists", sn), jen.Qual(constants.MockPkg, "Anything")},
			ot.BuildRequisiteFakeVarCallArgsForServiceExistenceHandlerTest(proj)...,
		)

		out = append(out,
			jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
			jen.ID(dataManagerVarName).Dot("On").Call(callArgs...).Dot("Return").Call(returnVals...),
			jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
			jen.Line(),
		)
	}

	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	dataManagerVarName := fmt.Sprintf("%sDataManager", typ.Name.UnexportedVarName())

	if indexToReturnFalse == len(owners)-1 && len(owners) > 0 {
		return
	}

	out = append(out,
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(actualCallArgs...).Dot("Return").Call(returnValues...),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
	)

	return
}

func determineMockExpecters(proj *models.Project, typ models.DataType, indexToStopAt int) []string {
	out := []string{}

	for i, ot := range append(proj.FindOwnerTypeChain(typ), typ) {
		if i > indexToStopAt && indexToStopAt > -1 {
			continue
		}
		dataManagerVarName := fmt.Sprintf("%sDataManager", ot.Name.UnexportedVarName())
		out = append(out, dataManagerVarName)
	}

	return out
}

func buildTestServiceListFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	uvn := typ.Name.UnexportedVarName()

	dataManagerVarName := fmt.Sprintf("%sDataManager", typ.Name.UnexportedVarName())

	getSomethingExpectedArgs := []jen.Code{
		jen.Litf("Get%s", pn),
		jen.Qual(constants.MockPkg, "Anything"), // ctx
	}
	getSomethingExpectedArgs = append(getSomethingExpectedArgs, typ.BuildCallArgsForDBClientListRetrievalMethodTest(proj)...)
	getSomethingExpectedArgs = append(getSomethingExpectedArgs, jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.QueryFilter")))

	firstSubtestLines := append([]jen.Code{jen.ID("s").Assign().ID("buildTestService").Call()}, includeOwnerFetchers(proj, typ)...)
	firstSubtestLines = append(firstSubtestLines,
		jen.IDf("example%sList", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(getSomethingExpectedArgs...).Dot("Return").Call(jen.IDf("example%sList", sn), jen.Nil()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%sList", sn)),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ListHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName, "ed"),
	)

	secondSubtestLines := append([]jen.Code{jen.ID("s").Assign().ID("buildTestService").Call()}, includeOwnerFetchers(proj, typ)...)
	secondSubtestLines = append(secondSubtestLines, jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(getSomethingExpectedArgs...).Dot("Return").Call(
			jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn))).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%sList", sn)),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ListHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName, "ed"),
	)

	thirdSubtestLines := append([]jen.Code{jen.ID("s").Assign().ID("buildTestService").Call()}, includeOwnerFetchers(proj, typ)...)
	thirdSubtestLines = append(thirdSubtestLines,
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(getSomethingExpectedArgs...).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn))).Call(jen.Nil()), constants.ObligatoryError()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ListHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName),
	)

	fourthSubtestLines := append([]jen.Code{jen.ID("s").Assign().ID("buildTestService").Call()}, includeOwnerFetchers(proj, typ)...)
	fourthSubtestLines = append(fourthSubtestLines,
		jen.IDf("example%sList", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(getSomethingExpectedArgs...).Dot("Return").Call(jen.IDf("example%sList", sn), jen.Nil()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%sList", sn)),
		).Dot("Return").Call(constants.ObligatoryError()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ListHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName, "ed"),
	)

	block := append([]jen.Code{jen.ID("T").Dot("Parallel").Call(), jen.Line()}, buildRelevantIDFetchers(proj, typ)...)
	block = append(block,
		jen.Line(), utils.BuildSubTestWithoutContext("happy path", firstSubtestLines...),
		jen.Line(), utils.BuildSubTestWithoutContext("with no rows returned", secondSubtestLines...),
		jen.Line(), utils.BuildSubTestWithoutContext(fmt.Sprintf("with error fetching %s from database", pcn), thirdSubtestLines...),
		jen.Line(), utils.BuildSubTestWithoutContext("with error encoding response", fourthSubtestLines...),
	)

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_ListHandler", pn)).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceSearchFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()

	happyPathSubtest := []jen.Code{
		jen.ID("s").Assign().ID("buildTestService").Call(),
		jen.Line(),
		jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher"),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName("Query")).Assign().Lit("whatever"),
		jen.ID(utils.BuildFakeVarName("Limit")).Assign().Uint8().Call(jen.Lit(123)),
		jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call().Dot(pn),
		jen.Var().IDf("example%sIDs", sn).Index().Uint64(),
		jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn)).Block(
			jen.IDf("example%sIDs", sn).Equals().Append(jen.IDf("example%sIDs", sn), jen.ID("x").Dot("ID")),
		),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Dot("On").Call(
					jen.Lit("Search"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Query")),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("Return").Call(
					jen.IDf("example%sIDs", sn),
					jen.Nil(),
				)
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("s").Dot("search").Equals().ID("si")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.IDf("%sDataManager", uvn).Dot("On").Call(
			jen.Litf("Get%sWithIDs", pn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
			jen.ID(utils.BuildFakeVarName("Limit")),
			jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sIDs", sn))),
		).Dot("Return").Call(
			jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))),
			jen.Nil(),
		),
		jen.ID("s").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("[]models.%s", sn)),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().Qual("net/http/httptest", "NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d"),
				jen.ID(utils.BuildFakeVarName("Query")),
				jen.ID(utils.BuildFakeVarName("Limit")),
			),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("SearchHandler").Call().Call(
			jen.ID(constants.ResponseVarName),
			jen.ID(constants.RequestVarName),
		),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("si", fmt.Sprintf("%sDataManager", uvn), "ed"),
	}

	searchErrSubtest := []jen.Code{
		jen.ID("s").Assign().ID("buildTestService").Call(),
		jen.Line(),
		jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher"),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName("Query")).Assign().Lit("whatever"),
		jen.ID(utils.BuildFakeVarName("Limit")).Assign().Uint8().Call(jen.Lit(123)),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Dot("On").Call(
					jen.Lit("Search"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Query")),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("Return").Call(
					jen.Index().Uint64().Values(),
					constants.ObligatoryError(),
				)
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("s").Dot("search").Equals().ID("si")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().Qual("net/http/httptest", "NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d"),
				jen.ID(utils.BuildFakeVarName("Query")),
				jen.ID(utils.BuildFakeVarName("Limit")),
			),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("SearchHandler").Call().Call(
			jen.ID(constants.ResponseVarName),
			jen.ID(constants.RequestVarName),
		),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("si"),
	}

	noRowsSubtest := []jen.Code{
		jen.ID("s").Assign().ID("buildTestService").Call(),
		jen.Line(),
		jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher"),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName("Query")).Assign().Lit("whatever"),
		jen.ID(utils.BuildFakeVarName("Limit")).Assign().Uint8().Call(jen.Lit(123)),
		jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call().Dot(pn),
		jen.Var().IDf("example%sIDs", sn).Index().Uint64(),
		jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn)).Block(
			jen.IDf("example%sIDs", sn).Equals().Append(jen.IDf("example%sIDs", sn), jen.ID("x").Dot("ID")),
		),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Dot("On").Call(
					jen.Lit("Search"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Query")),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("Return").Call(
					jen.IDf("example%sIDs", sn),
					jen.Nil(),
				)
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("s").Dot("search").Equals().ID("si")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.IDf("%sDataManager", uvn).Dot("On").Call(
			jen.Litf("Get%sWithIDs", pn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
			jen.ID(utils.BuildFakeVarName("Limit")),
			jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sIDs", sn))),
		).Dot("Return").Call(
			jen.Index().Qual(proj.ModelsV1Package(), sn).Values(),
			jen.Qual("database/sql", "ErrNoRows"),
		),
		jen.ID("s").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("[]models.%s", sn)),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().Qual("net/http/httptest", "NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d"),
				jen.ID(utils.BuildFakeVarName("Query")),
				jen.ID(utils.BuildFakeVarName("Limit")),
			),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("SearchHandler").Call().Call(
			jen.ID(constants.ResponseVarName),
			jen.ID(constants.RequestVarName),
		),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("si", fmt.Sprintf("%sDataManager", uvn), "ed"),
	}

	fetchErrSubtest := []jen.Code{
		jen.ID("s").Assign().ID("buildTestService").Call(),
		jen.Line(),
		jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher"),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName("Query")).Assign().Lit("whatever"),
		jen.ID(utils.BuildFakeVarName("Limit")).Assign().Uint8().Call(jen.Lit(123)),
		jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call().Dot(pn),
		jen.Var().IDf("example%sIDs", sn).Index().Uint64(),
		jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn)).Block(
			jen.IDf("example%sIDs", sn).Equals().Append(jen.IDf("example%sIDs", sn), jen.ID("x").Dot("ID")),
		),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Dot("On").Call(
					jen.Lit("Search"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Query")),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("Return").Call(
					jen.IDf("example%sIDs", sn),
					jen.Nil(),
				)
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("s").Dot("search").Equals().ID("si")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.IDf("%sDataManager", uvn).Dot("On").Call(
			jen.Litf("Get%sWithIDs", pn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
			jen.ID(utils.BuildFakeVarName("Limit")),
			jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sIDs", sn))),
		).Dot("Return").Call(
			jen.Index().Qual(proj.ModelsV1Package(), sn).Values(),
			constants.ObligatoryError(),
		),
		jen.ID("s").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().Qual("net/http/httptest", "NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d"),
				jen.ID(utils.BuildFakeVarName("Query")),
				jen.ID(utils.BuildFakeVarName("Limit")),
			),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("SearchHandler").Call().Call(
			jen.ID(constants.ResponseVarName),
			jen.ID(constants.RequestVarName),
		),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("si", fmt.Sprintf("%sDataManager", uvn)),
	}

	encodeErrSubtest := []jen.Code{
		jen.ID("s").Assign().ID("buildTestService").Call(),
		jen.Line(),
		jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher"),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName("Query")).Assign().Lit("whatever"),
		jen.ID(utils.BuildFakeVarName("Limit")).Assign().Uint8().Call(jen.Lit(123)),
		jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call().Dot(pn),
		jen.Var().IDf("example%sIDs", sn).Index().Uint64(),
		jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn)).Block(
			jen.IDf("example%sIDs", sn).Equals().Append(jen.IDf("example%sIDs", sn), jen.ID("x").Dot("ID")),
		),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Dot("On").Call(
					jen.Lit("Search"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Query")),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("Return").Call(
					jen.IDf("example%sIDs", sn),
					jen.Nil(),
				)
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("s").Dot("search").Equals().ID("si")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.IDf("%sDataManager", uvn).Dot("On").Call(
			jen.Litf("Get%sWithIDs", pn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
			jen.ID(utils.BuildFakeVarName("Limit")),
			jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sIDs", sn))),
		).Dot("Return").Call(
			jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))),
			jen.Nil(),
		),
		jen.ID("s").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("[]models.%s", sn)),
		).Dot("Return").Call(constants.ObligatoryError()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().Qual("net/http/httptest", "NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d"),
				jen.ID(utils.BuildFakeVarName("Query")),
				jen.ID(utils.BuildFakeVarName("Limit")),
			),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("SearchHandler").Call().Call(
			jen.ID(constants.ResponseVarName),
			jen.ID(constants.RequestVarName),
		),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("si", fmt.Sprintf("%sDataManager", uvn), "ed"),
	}

	block := append([]jen.Code{jen.ID("T").Dot("Parallel").Call(), jen.Line()}, buildRelevantIDFetchers(proj, typ)...)
	block = append(block,
		jen.Line(), utils.BuildSubTestWithoutContext("happy path", happyPathSubtest...),
		jen.Line(), utils.BuildSubTestWithoutContext("with error conducting search", searchErrSubtest...),
		jen.Line(), utils.BuildSubTestWithoutContext("with now rows returned", noRowsSubtest...),
		jen.Line(), utils.BuildSubTestWithoutContext("with error fetching from database", fetchErrSubtest...),
		jen.Line(), utils.BuildSubTestWithoutContext("with error encoding response", encodeErrSubtest...),
	)

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_SearchHandler", pn)).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceCreateFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	happyPathSubtest := []jen.Code{
		jen.ID("s").Assign().ID("buildTestService").Call(),
		jen.Line(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName).Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.ID(utils.BuildFakeVarName(sn))),
		jen.Line(),
		jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.IDf("%sDataManager", uvn).Dot("On").Call(
			jen.Litf("Create%s", sn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%sCreationInput", sn)),
		).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
		jen.ID("s").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
		jen.Line(),
		jen.ID("mc").Assign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
		jen.ID("mc").Dot("On").Call(jen.Lit("Increment"),
			jen.Qual(constants.MockPkg, "Anything"),
		),
		jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Equals().ID("mc"),
		jen.Line(),
		jen.ID("r").Assign().AddressOf().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
		jen.ID("r").Dot("On").Call(
			jen.Lit("Report"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("newsman.Event")),
		).Dot("Return").Call(),
		jen.ID("s").Dot("reporter").Equals().ID("r"),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Dot("On").Call(
					jen.Lit("Index"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
					jen.ID(utils.BuildFakeVarName(sn)),
				).Dot("Return").Call(jen.Nil())
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("s").Dot("search").Equals().ID("si")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%s", sn)),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().Qual("net/http/httptest", "NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Call(
			jen.Qual("context", "WithValue").Call(
				jen.ID(constants.RequestVarName).Dot("Context").Call(),
				jen.ID("createMiddlewareCtxKey"),
				jen.ID(utils.BuildFakeVarName("Input")),
			),
		),
		jen.Line(),
		jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusCreated"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return utils.AssertExpectationsFor(fmt.Sprintf("%sDataManager", uvn), "mc", "r", "si", "ed")
			}
			return utils.AssertExpectationsFor(fmt.Sprintf("%sDataManager", uvn), "mc", "r", "ed")
		}(),
	}

	noInputSubtest := []jen.Code{
		jen.ID("s").Assign().ID("buildTestService").Call(),
		jen.Line(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().Qual("net/http/httptest", "NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
	}

	creationErrSubtest := []jen.Code{
		jen.ID("s").Assign().ID("buildTestService").Call(),
		jen.Line(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName).Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.ID(utils.BuildFakeVarName(sn))),
		jen.Line(),
		jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.IDf("%sDataManager", uvn).Dot("On").Call(
			jen.Litf("Create%s", sn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%sCreationInput", sn)),
		).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Parens(jen.Nil()), constants.ObligatoryError()),
		jen.ID("s").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().Qual("net/http/httptest", "NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Call(
			jen.Qual("context", "WithValue").Call(
				jen.ID(constants.RequestVarName).Dot("Context").Call(),
				jen.ID("createMiddlewareCtxKey"),
				jen.ID(utils.BuildFakeVarName("Input")),
			),
		),
		jen.Line(),
		jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(fmt.Sprintf("%sDataManager", uvn)),
	}

	encodeErrSubtest := []jen.Code{
		jen.ID("s").Assign().ID("buildTestService").Call(),
		jen.Line(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName).Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.ID(utils.BuildFakeVarName(sn))),
		jen.Line(),
		jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.IDf("%sDataManager", uvn).Dot("On").Call(
			jen.Litf("Create%s", sn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%sCreationInput", sn)),
		).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
		jen.ID("s").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
		jen.Line(),
		jen.ID("mc").Assign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
		jen.ID("mc").Dot("On").Call(jen.Lit("Increment"),
			jen.Qual(constants.MockPkg, "Anything"),
		),
		jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Equals().ID("mc"),
		jen.Line(),
		jen.ID("r").Assign().AddressOf().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
		jen.ID("r").Dot("On").Call(
			jen.Lit("Report"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("newsman.Event")),
		).Dot("Return").Call(),
		jen.ID("s").Dot("reporter").Equals().ID("r"),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Dot("On").Call(
					jen.Lit("Index"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
					jen.ID(utils.BuildFakeVarName(sn)),
				).Dot("Return").Call(jen.Nil())
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("s").Dot("search").Equals().ID("si")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%s", sn)),
		).Dot("Return").Call(constants.ObligatoryError()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().Qual("net/http/httptest", "NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Call(
			jen.Qual("context", "WithValue").Call(
				jen.ID(constants.RequestVarName).Dot("Context").Call(),
				jen.ID("createMiddlewareCtxKey"),
				jen.ID(utils.BuildFakeVarName("Input")),
			),
		),
		jen.Line(),
		jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusCreated"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return utils.AssertExpectationsFor(fmt.Sprintf("%sDataManager", uvn), "mc", "r", "si", "ed")
			}
			return utils.AssertExpectationsFor(fmt.Sprintf("%sDataManager", uvn), "mc", "r", "ed")
		}(),
	}

	block := append(
		[]jen.Code{
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
		},
		buildRelevantIDFetchers(proj, typ)...,
	)
	block = append(block, jen.Line())
	block = append(
		block,
		jen.ID("T").Dot("Run").Call(
			jen.Lit("happy path"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
				happyPathSubtest...,
			),
		),
		jen.Line(),
	)

	block = append(
		block,
		jen.ID("T").Dot("Run").Call(
			jen.Lit("without input attached"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
				noInputSubtest...,
			),
		),
		jen.Line(),
	)

	block = append(
		block,
		jen.ID("T").Dot("Run").Call(
			jen.Litf("with error creating %s", scn),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
				creationErrSubtest...,
			),
		),
		jen.Line(),
	)

	block = append(
		block,
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error encoding response"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
				encodeErrSubtest...,
			),
		),
		jen.Line(),
	)

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_CreateHandler", pn)).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceExistenceFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	dataManagerVarName := fmt.Sprintf("%sDataManager", typ.Name.UnexportedVarName())

	expectedCallArgs := append(
		[]jen.Code{
			jen.Litf("%sExists", sn),
			jen.Qual(constants.MockPkg, "Anything"),
		},
		typ.BuildRequisiteFakeVarCallArgsForServiceExistenceHandlerTest(proj)...,
	)

	firstSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)

	firstSubtestLines = append(firstSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(expectedCallArgs...).Dot("Return").Call(jen.True(), jen.Nil()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ExistenceHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName),
	)
	firstSubtest := utils.BuildSubTestWithoutContext("happy path", firstSubtestLines...)

	secondSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)

	secondSubtestLines = append(secondSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(), func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(expectedCallArgs...).Dot("Return").Call(jen.False(), jen.Qual("database/sql", "ErrNoRows")),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ExistenceHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName),
	)
	secondSubtest := utils.BuildSubTestWithoutContext(fmt.Sprintf("with no such %s in database", scn), secondSubtestLines...)

	thirdSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)

	thirdSubtestLines = append(thirdSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(), func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(expectedCallArgs...).Dot("Return").Call(jen.False(), constants.ObligatoryError()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ExistenceHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName),
	)
	thirdSubtest := utils.BuildSubTestWithoutContext(fmt.Sprintf("with error fetching %s from database", scn), thirdSubtestLines...)

	block := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
	}
	block = append(block, buildRelevantIDFetchers(proj, typ)...)

	block = append(block,
		jen.Line(),
		firstSubtest,
		jen.Line(),
		secondSubtest,
		jen.Line(),
		thirdSubtest,
	)

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_ExistenceHandler", pn)).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceReadFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	dataManagerVarName := fmt.Sprintf("%sDataManager", typ.Name.UnexportedVarName())

	getSomethingExpectationArgs := []jen.Code{
		jen.Litf("Get%s", sn),
		jen.Qual(constants.MockPkg, "Anything"),
	}
	getSomethingExpectationArgs = append(getSomethingExpectationArgs, typ.BuildRequisiteFakeVarCallArgsForServiceReadHandlerTest(proj)...)

	firstSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)

	firstSubtestLines = append(firstSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(getSomethingExpectationArgs...).Dot("Return").Call(jen.ID(utils.BuildFakeVarName(sn)), jen.Nil()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%s", sn)),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName, "ed"),
	)
	firstSubtest := utils.BuildSubTestWithoutContext("happy path", firstSubtestLines...)

	secondSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	secondSubtestLines = append(secondSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(getSomethingExpectationArgs...).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName),
	)
	secondSubtest := utils.BuildSubTestWithoutContext(fmt.Sprintf("with no such %s in database", scn), secondSubtestLines...)

	thirdSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	thirdSubtestLines = append(thirdSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(getSomethingExpectationArgs...).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Call(jen.Nil()), constants.ObligatoryError()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName),
	)
	thirdSubtest := utils.BuildSubTestWithoutContext(fmt.Sprintf("with error fetching %s from database", scn), thirdSubtestLines...)

	fourthSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	fourthSubtestLines = append(fourthSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(getSomethingExpectationArgs...).Dot("Return").Call(jen.ID(utils.BuildFakeVarName(sn)), jen.Nil()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%s", sn)),
		).Dot("Return").Call(constants.ObligatoryError()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName, "ed"),
	)
	fourthSubtest := utils.BuildSubTestWithoutContext("with error encoding response", fourthSubtestLines...)

	block := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
	}

	block = append(block, buildRelevantIDFetchers(proj, typ)...)
	block = append(block,
		jen.Line(), firstSubtest,
		jen.Line(), secondSubtest,
		jen.Line(), thirdSubtest,
		jen.Line(), fourthSubtest,
	)

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_ReadHandler", pn)).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceUpdateFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	dataManagerVarName := fmt.Sprintf("%sDataManager", typ.Name.UnexportedVarName())

	expectedDBRetrievalCallArgs := []jen.Code{
		jen.Litf("Get%s", sn),
		jen.Qual(constants.MockPkg, "Anything"),
	}
	expectedDBRetrievalCallArgs = append(expectedDBRetrievalCallArgs, typ.BuildRequisiteFakeVarCallArgsForServiceUpdateHandlerTest(proj)...)

	firstSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	firstSubtestLines = append(firstSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.ID(utils.BuildFakeVarName(sn))),
		jen.Line(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(expectedDBRetrievalCallArgs...).Dot("Return").Call(jen.ID(utils.BuildFakeVarName(sn)), jen.Nil()),
		jen.ID(dataManagerVarName).Dot("On").Call(
			jen.Litf("Update%s", sn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%s", sn)),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID("r").Assign().AddressOf().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
		jen.ID("r").Dot("On").Call(
			jen.Lit("Report"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("newsman.Event")),
		).Dot("Return").Call(),
		jen.ID("s").Dot("reporter").Equals().ID("r"),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Dot("On").Call(
					jen.Lit("Index"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
					jen.ID(utils.BuildFakeVarName(sn)),
				).Dot("Return").Call(jen.Nil())
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("s").Dot("search").Equals().ID("si")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%s", sn)),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("updateMiddlewareCtxKey"), jen.ID(utils.BuildFakeVarName("Input")))),
		jen.Line(),
		jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("r", dataManagerVarName, "ed"),
	)
	firstSubtest := utils.BuildSubTestWithoutContext("happy path", firstSubtestLines...)

	secondSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	secondSubtestLines = append(secondSubtestLines,
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
	)
	secondSubtest := utils.BuildSubTestWithoutContext("without update input", secondSubtestLines...)

	thirdSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	thirdSubtestLines = append(thirdSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.ID(utils.BuildFakeVarName(sn))),
		jen.Line(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(expectedDBRetrievalCallArgs...).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("updateMiddlewareCtxKey"), jen.ID(utils.BuildFakeVarName("Input")))),
		jen.Line(),
		jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName),
	)
	thirdSubtest := utils.BuildSubTestWithoutContext(fmt.Sprintf("with no rows fetching %s", scn), thirdSubtestLines...)

	fourthSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	fourthSubtestLines = append(fourthSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.ID(utils.BuildFakeVarName(sn))),
		jen.Line(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(expectedDBRetrievalCallArgs...).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Call(jen.Nil()), constants.ObligatoryError()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("updateMiddlewareCtxKey"), jen.ID(utils.BuildFakeVarName("Input")))),
		jen.Line(),
		jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName),
	)
	fourthSubtest := utils.BuildSubTestWithoutContext(fmt.Sprintf("with error fetching %s", scn), fourthSubtestLines...)

	fifthSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	fifthSubtestLines = append(fifthSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.ID(utils.BuildFakeVarName(sn))),
		jen.Line(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(expectedDBRetrievalCallArgs...).Dot("Return").Call(jen.ID(utils.BuildFakeVarName(sn)), jen.Nil()),
		jen.ID(dataManagerVarName).Dot("On").Call(
			jen.Litf("Update%s", sn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%s", sn)),
		).Dot("Return").Call(constants.ObligatoryError()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("updateMiddlewareCtxKey"), jen.ID(utils.BuildFakeVarName("Input")))),
		jen.Line(),
		jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(dataManagerVarName),
	)
	fifthSubtest := utils.BuildSubTestWithoutContext(fmt.Sprintf("with error updating %s", scn), fifthSubtestLines...)

	sixthSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	sixthSubtestLines = append(sixthSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.ID(utils.BuildFakeVarName(sn))),
		jen.Line(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		),
		jen.Line(),
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(expectedDBRetrievalCallArgs...).Dot("Return").Call(jen.ID(utils.BuildFakeVarName(sn)), jen.Nil()),
		jen.ID(dataManagerVarName).Dot("On").Call(
			jen.Litf("Update%s", sn),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%s", sn)),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
		jen.ID("r").Assign().AddressOf().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
		jen.ID("r").Dot("On").Call(
			jen.Lit("Report"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("newsman.Event")),
		).Dot("Return").Call(),
		jen.ID("s").Dot("reporter").Equals().ID("r"),
		jen.Line(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("si").Dot("On").Call(
					jen.Lit("Index"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
					jen.ID(utils.BuildFakeVarName(sn)),
				).Dot("Return").Call(jen.Nil())
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("s").Dot("search").Equals().ID("si")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ed").Dot("On").Call(
			jen.Lit("EncodeResponse"),
			jen.Qual(constants.MockPkg, "Anything"),
			jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Litf("*models.%s", sn)),
		).Dot("Return").Call(constants.ObligatoryError()),
		jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("updateMiddlewareCtxKey"), jen.ID(utils.BuildFakeVarName("Input")))),
		jen.Line(),
		jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("r", dataManagerVarName, "ed"),
	)
	sixthSubtest := utils.BuildSubTestWithoutContext("with error encoding response", sixthSubtestLines...)

	block := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
	}

	block = append(block, buildRelevantIDFetchers(proj, typ)...)
	block = append(block,
		jen.Line(),
		firstSubtest, jen.Line(),
		secondSubtest, jen.Line(),
		thirdSubtest, jen.Line(),
		fourthSubtest, jen.Line(),
		fifthSubtest, jen.Line(),
		sixthSubtest,
	)

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_UpdateHandler", pn)).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func setupDataManagersForDeletion(proj *models.Project, typ models.DataType, actualCallArgs, returnValues []jen.Code, indexToReturnFalse int, returnErr, returnFalse bool) (out []jen.Code) {
	owners := proj.FindOwnerTypeChain(typ)
	for i, ot := range owners {
		sn := ot.Name.Singular()
		uvn := ot.Name.UnexportedVarName()
		dataManagerVarName := fmt.Sprintf("%sDataManager", ot.Name.UnexportedVarName())

		if i > indexToReturnFalse && indexToReturnFalse > -1 {
			return
		}

		returnVals := []jen.Code{
			jen.True(),
			jen.Nil(),
		}
		if i == indexToReturnFalse && returnFalse {
			returnVals[0] = jen.False()
		}
		if i == indexToReturnFalse && returnErr {
			returnVals[1] = constants.ObligatoryError()
		}

		callArgs := append(
			[]jen.Code{jen.Litf("%sExists", sn), jen.Qual(constants.MockPkg, "Anything")},
			ot.BuildRequisiteFakeVarCallArgsForServiceExistenceHandlerTest(proj)...,
		)

		out = append(out,
			jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
			jen.ID(dataManagerVarName).Dot("On").Call(callArgs...).Dot("Return").Call(returnVals...),
			jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
			jen.Line(),
		)
	}

	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	dataManagerVarName := fmt.Sprintf("%sDataManager", typ.Name.UnexportedVarName())

	if indexToReturnFalse == len(owners)-1 && len(owners) > 0 {
		return
	}

	out = append(out,
		jen.ID(dataManagerVarName).Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
		jen.ID(dataManagerVarName).Dot("On").Call(actualCallArgs...).Dot("Return").Call(returnValues...),
		jen.ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Equals().ID(dataManagerVarName),
		jen.Line(),
	)

	return
}

func buildTestServiceArchiveFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	subtests := []jen.Code{}

	expectedCallArgs := []jen.Code{
		jen.Litf("Archive%s", sn),
		jen.Qual(constants.MockPkg, "Anything"),
	}
	expectedCallArgs = append(expectedCallArgs, typ.BuildRequisiteFakeVarCallArgsForServiceArchiveHandlerTest()...)

	buildHappyPathSubtestVariant := func(name string, index int, returnFalse, returnErr, includeSelf, includeAfterEffects, includeSearch bool) jen.Code {
		lines := append(
			[]jen.Code{
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
			},
			includeOwnerFetchers(proj, typ)...,
		)
		lines = append(lines, jen.Line())

		if includeSelf {
			lines = append(lines,
				jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if typ.BelongsToUser {
						return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
					}
					return jen.Null()
				}(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().ID(utils.BuildFakeVarName(sn)).Dot("ID"),
				),
				jen.Line(),
			)
		}

		returnValues := []jen.Code{
			jen.Nil(),
		}

		var expectedStatus string
		if index < 0 {
			expectedStatus = "StatusNoContent"
		} else {
			expectedStatus = "StatusNotFound"
		}
		if returnErr {
			expectedStatus = "StatusInternalServerError"
		}

		var elems []string
		if index == -1 && includeAfterEffects {
			elems = []string{
				"mc",
				"r",
			}
		}

		lines = append(lines, setupDataManagersForDeletion(proj, typ, expectedCallArgs, returnValues, index, returnErr, returnFalse)...)

		if includeAfterEffects {
			lines = append(lines,
				jen.ID("r").Assign().AddressOf().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(
					jen.Lit("Report"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("newsman.Event")),
				).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Equals().ID("r"),
				jen.Line(),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.ID("si").Assign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values()
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.ID("si").Dot("On").Call(
							jen.Lit("Delete"),
							jen.Qual(constants.MockPkg, "Anything"),
							jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
						).Dot("Return").Call(
							jen.Nil(),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.ID("s").Dot("search").Equals().ID("si")
					}
					return jen.Null()
				}(),
				jen.Line(),
				jen.ID("mc").Assign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Decrement"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Equals().ID("mc"),
				jen.Line(),
			)
		}

		lines = append(lines,
			jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
			jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodGet"),
				jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
				jen.Nil(),
			),
			utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
			utils.RequireNoError(jen.Err(), nil),
			jen.Line(),
			jen.ID("s").Dot("ArchiveHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			jen.Line(),
			utils.AssertEqual(jen.Qual("net/http", expectedStatus), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
			jen.Line(),
			utils.AssertExpectationsFor(append(determineMockExpecters(proj, typ, index), elems...)...),
		)

		return utils.BuildSubTestWithoutContext(name, lines...)
	}

	subtests = append(subtests, buildHappyPathSubtestVariant(
		"happy path",
		-1,
		false,
		false,
		true,
		true,
		true,
	))

	for i, ot := range proj.FindOwnerTypeChain(typ) {
		subtests = append(subtests, buildHappyPathSubtestVariant(
			fmt.Sprintf("with nonexistent %s", ot.Name.SingularCommonName()),
			i,
			true,
			false,
			false,
			false,
			false,
		))
		subtests = append(subtests, buildHappyPathSubtestVariant(
			fmt.Sprintf("with error checking %s existence", ot.Name.SingularCommonName()),
			i,
			false,
			true,
			false,
			false,
			false,
		))
	}

	secondSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	secondSubtestLines = append(secondSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return(jen.IDf(utils.BuildFakeVarName(sn)).Dot("ID")),
		),
		jen.Line(), jen.Line(),
	)

	secondSubtestLines = append(secondSubtestLines, setupDataManagersForDeletion(proj, typ, expectedCallArgs, []jen.Code{jen.Qual("database/sql", "ErrNoRows")}, -1, false, false)...)
	secondSubtestLines = append(secondSubtestLines,
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ArchiveHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(append(determineMockExpecters(proj, typ, -1))...),
	)
	subtests = append(subtests, utils.BuildSubTestWithoutContext(fmt.Sprintf("with no %s in database", scn), secondSubtestLines...))

	thirdSubtestLines := append(
		[]jen.Code{
			jen.ID("s").Assign().ID("buildTestService").Call(),
			jen.Line(),
		},
		includeOwnerFetchers(proj, typ)...,
	)
	thirdSubtestLines = append(thirdSubtestLines,
		jen.Line(),
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return(jen.IDf(utils.BuildFakeVarName(sn)).Dot("ID")),
		),
		jen.Line(), jen.Line(),
	)

	thirdSubtestLines = append(thirdSubtestLines, setupDataManagersForDeletion(proj, typ, expectedCallArgs, []jen.Code{constants.ObligatoryError()}, -1, false, false)...)
	thirdSubtestLines = append(thirdSubtestLines,
		jen.Line(),
		jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
			jen.Nil(),
		),
		utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.ID("s").Dot("ArchiveHandler").Call().Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		jen.Line(),
		utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
		jen.Line(),
		utils.AssertExpectationsFor(append(determineMockExpecters(proj, typ, -1))...),
	)
	subtests = append(subtests, utils.BuildSubTestWithoutContext("with error writing to database", thirdSubtestLines...))

	block := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
	}

	block = append(block, buildRelevantIDFetchers(proj, typ)...)

	for _, st := range subtests {
		block = append(block, jen.Line(), st)
	}

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_ArchiveHandler", pn)).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			block...,
		),
		jen.Line(),
	}

	return lines
}
