package iterables

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, ret)

	ret.Add(buildTestServiceListFuncDecl(proj, typ)...)
	ret.Add(buildTestServiceCreateFuncDecl(proj, typ)...)
	ret.Add(buildTestServiceExistenceFuncDecl(proj, typ)...)
	ret.Add(buildTestServiceReadFuncDecl(proj, typ)...)
	ret.Add(buildTestServiceUpdateFuncDecl(proj, typ)...)
	ret.Add(buildTestServiceArchiveFuncDecl(proj, typ)...)

	return ret
}

func buildOwnerVarName(typ models.DataType) string {
	if typ.BelongsToUser {
		return "requestingUser"
	} else if typ.BelongsToStruct != nil {
		return fmt.Sprintf("requesting%s", typ.BelongsToStruct.Singular())
	}

	return ""
}

func buildRelevantOwnerVar(proj *models.Project, typ models.DataType) jen.Code {
	g := &jen.Group{}

	if typ.BelongsToUser {
		g.Add(jen.ID(buildOwnerVarName(typ)).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	} else if typ.BelongsToStruct != nil {
		g.Add(jen.ID(buildOwnerVarName(typ)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.BelongsToStruct.Singular())).Call())
	}

	return g
}

func includeUserFetcher(typ models.DataType) jen.Code {
	if typ.BelongsToUser {
		return jen.ID("s").Dot("userIDFetcher").Equals().ID("userIDFetcher")
	}
	return jen.Null()
}

func includeOwnerFetcher(typ models.DataType) jen.Code {
	if typ.BelongsToStruct != nil {
		btsuvn := typ.BelongsToStruct.UnexportedVarName()
		return jen.ID("s").Dotf("%sIDFetcher", btsuvn).Equals().IDf("%sIDFetcher", btsuvn)
	}
	return jen.Null()
}

func buildDBCallOwnerVar(typ models.DataType) jen.Code {
	if typ.BelongsToUser || typ.BelongsToStruct != nil {
		return jen.ID(buildOwnerVarName(typ)).Dot("ID")
	}

	return nil
}

func buildRelevantIDFetcher(typ models.DataType) jen.Code {
	if typ.BelongsToUser {
		return jen.ID("userIDFetcher").Assign().Func().Params(jen.Underscore().ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(buildOwnerVarName(typ)).Dot("ID"),
		)
	} else if typ.BelongsToStruct != nil {
		return jen.IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Assign().Func().Params(jen.Underscore().ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
			jen.Return().ID(buildOwnerVarName(typ)).Dot("ID"),
		)
	}

	return nil
}

func buildTestServiceListFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	uvn := typ.Name.UnexportedVarName()

	dataManagerVarName := fmt.Sprintf("%sdm", typ.Name.LowercaseAbbreviation())

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_ListHandler", pn)).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			buildRelevantOwnerVar(proj, typ),
			buildRelevantIDFetcher(typ),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%sList", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", pn),
					jen.Qual(utils.MockPkg, "Anything"),
					buildDBCallOwnerVar(typ),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.IDf("example%sList", sn), jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with no rows returned",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", pn),
					jen.Qual(utils.MockPkg, "Anything"),
					buildDBCallOwnerVar(typ),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn))).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with error fetching %s from database", pcn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", pn),
					jen.Qual(utils.MockPkg, "Anything"),
					buildDBCallOwnerVar(typ),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn))).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error encoding response",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%sList", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", pn),
					jen.Qual(utils.MockPkg, "Anything"),
					buildDBCallOwnerVar(typ),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.IDf("example%sList", sn), jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
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

	dataManagerVarName := fmt.Sprintf("%sdm", typ.Name.LowercaseAbbreviation())

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_CreateHandler", pn)).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			buildRelevantOwnerVar(proj, typ),
			buildRelevantIDFetcher(typ),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
				jen.Line(),
				jen.ID("mc").Assign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"),
					jen.Qual(utils.MockPkg, "Anything"),
				),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Equals().ID("mc"),
				jen.Line(),
				jen.ID("r").Assign().VarPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Equals().ID("r"),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Create%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusCreated"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without input attached",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with error creating %s", scn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Create%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error encoding response",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
				jen.Line(),
				jen.ID("mc").Assign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"), jen.Qual(utils.MockPkg, "Anything")),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Equals().ID("mc"),
				jen.Line(),
				jen.ID("r").Assign().VarPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Equals().ID("r"),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Create%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusCreated"), jen.ID("res").Dot("Code"), nil),
			),
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

	dataManagerVarName := fmt.Sprintf("%sdm", typ.Name.LowercaseAbbreviation())

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_ExistenceHandler", pn)).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			buildRelevantOwnerVar(proj, typ),
			buildRelevantIDFetcher(typ),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("%sExists", sn),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.True(), jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ExistenceHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with no such %s in database", scn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("%sExists", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.False(), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ExistenceHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with error fetching %s from database", scn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("%sExists", sn),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.False(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ExistenceHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error encoding response",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("%sExists", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.True(), jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ExistenceHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
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

	dataManagerVarName := fmt.Sprintf("%sdm", typ.Name.LowercaseAbbreviation())

	lines := []jen.Code{

		jen.Func().ID(fmt.Sprintf("Test%sService_ReadHandler", pn)).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			buildRelevantOwnerVar(proj, typ),
			buildRelevantIDFetcher(typ),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", sn),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with no such %s in database", scn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with error fetching %s from database", scn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", sn),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error encoding response",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
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

	dataManagerVarName := fmt.Sprintf("%sdm", typ.Name.LowercaseAbbreviation())

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_UpdateHandler", pn)).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			buildRelevantOwnerVar(proj, typ),
			buildRelevantIDFetcher(typ),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mc").Assign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"),
					jen.Qual(utils.MockPkg, "Anything"),
				),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Equals().ID("mc"),
				jen.Line(),
				jen.ID("r").Assign().VarPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Equals().ID("r"),
				jen.Line(),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Update%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without update input",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				//includeOwnerFetcher(typ),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with no rows fetching %s", scn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with error fetching %s", scn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with error updating %s", scn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mc").Assign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"), jen.Qual(utils.MockPkg, "Anything")),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Equals().ID("mc"),
				jen.Line(),
				jen.ID("r").Assign().VarPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Equals().ID("r"),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Update%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error encoding response",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mc").Assign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"), jen.Qual(utils.MockPkg, "Anything")),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Equals().ID("mc"),
				jen.Line(),
				jen.ID("r").Assign().VarPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Equals().ID("r"),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Update%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceArchiveFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	dataManagerVarName := fmt.Sprintf("%sdm", typ.Name.LowercaseAbbreviation())

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_ArchiveHandler", pn)).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			buildRelevantOwnerVar(proj, typ),
			buildRelevantIDFetcher(typ),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID("r").Assign().VarPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Equals().ID("r"),
				jen.Line(),
				jen.ID("mc").Assign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Decrement"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Equals().ID("mc"),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Archive%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("EncodeResponse"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ArchiveHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("with no %s in database", scn),
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Archive%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ArchiveHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error reading from database",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				includeUserFetcher(typ),
				includeOwnerFetcher(typ),
				jen.Line(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return().IDf("example%s", sn).Dot("ID"),
				),
				jen.Line(),
				jen.ID(dataManagerVarName).Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID(dataManagerVarName).Dot("On").Call(jen.Litf("Archive%s", sn), jen.Qual(utils.MockPkg, "Anything"),
					jen.IDf("example%s", sn).Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Equals().ID(dataManagerVarName),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("s").Dot("ArchiveHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}
