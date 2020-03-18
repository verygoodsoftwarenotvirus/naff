package iterables

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(pkg, ret)

	ret.Add(buildTestServiceListFuncDecl(pkg, typ)...)
	ret.Add(buildTestServiceCreateFuncDecl(pkg, typ)...)
	ret.Add(buildTestServiceReadFuncDecl(pkg, typ)...)
	ret.Add(buildTestServiceUpdateFuncDecl(pkg, typ)...)
	ret.Add(buildTestServiceArchiveFuncDecl(pkg, typ)...)

	return ret
}

func buildOwnerVarName(typ models.DataType) string {
	if typ.BelongsToUser {
		return "requestingUser"
	}
	if typ.BelongsToStruct != nil {
		return fmt.Sprintf("requesting%s", typ.BelongsToStruct.Singular())
	}

	return ""
}

func buildRelevantOwnerVar(pkg *models.Project, typ models.DataType) jen.Code {
	if typ.BelongsToUser {
		return jen.ID(buildOwnerVarName(typ)).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Add(utils.FakeUint64Func()))
	}
	if typ.BelongsToStruct != nil {
		return jen.ID(buildOwnerVarName(typ)).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), typ.BelongsToStruct.Singular()).Values(jen.ID("ID").Op(":").Add(utils.FakeUint64Func()))
	}

	return nil
}

func buildDBCallOwnerVar(typ models.DataType) jen.Code {
	if typ.BelongsToUser || typ.BelongsToStruct != nil {
		return jen.ID(buildOwnerVarName(typ)).Dot("ID")
	}

	return nil
}

func buildRelevantIDFetcher(typ models.DataType) jen.Code {
	if typ.BelongsToUser {
		return jen.ID("s").Dot("userIDFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().ID(buildOwnerVarName(typ)).Dot("ID"),
		)
	}
	if typ.BelongsToStruct != nil {
		return jen.ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().ID(buildOwnerVarName(typ)).Dot("ID"),
		)
	}

	return nil
}

func buildTestServiceListFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_List", pn)).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)).Valuesln(
					jen.ID(pn).Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
						),
					),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", pn),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with no rows returned"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", pn),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn))).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with error fetching %s from database", pcn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", pn),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn))).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusInternalServerError")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)).Valuesln(
					jen.ID(pn).Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", pn),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceCreateFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	buildCreationInputFromExpectedLines := func() []jen.Code {
		lines := []jen.Code{}

		for _, field := range typ.Fields {
			if field.ValidForCreationInput {
				sn := field.Name.Singular()
				lines = append(lines, jen.ID(sn).Op(":").ID("expected").Dot(sn))
			}
		}

		return lines
	}

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_Create", pn)).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				jen.ID("mc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Op("=").ID("mc"),
				jen.Line(),
				jen.ID("r").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Op("=").ID("r"),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Create%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dotf("%sCreationInput", sn).Valuesln(
					buildCreationInputFromExpectedLines()...,
				),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusCreated")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without input attached"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusBadRequest")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with error creating %s", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Create%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dotf("%sCreationInput", sn).Valuesln(
					buildCreationInputFromExpectedLines()...,
				),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusInternalServerError")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				jen.ID("mc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"), jen.Qual("github.com/stretchr/testify/mock", "Anything")),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Op("=").ID("mc"),
				jen.Line(),
				jen.ID("r").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Op("=").ID("r"),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Create%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dotf("%sCreationInput", sn).Valuesln(
					buildCreationInputFromExpectedLines()...,
				),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusCreated")),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceReadFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{

		jen.Func().ID(fmt.Sprintf("Test%sService_Read", pn)).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", sn),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with no such %s in database", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusNotFound")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with error fetching %s from database", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", sn),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusInternalServerError")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceUpdateFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	buildUpdateInputFromExpectedLines := func() []jen.Code {
		lines := []jen.Code{}

		for _, field := range typ.Fields {
			if field.ValidForUpdateInput {
				sn := field.Name.Singular()
				lines = append(lines, jen.ID(sn).Op(":").ID("expected").Dot(sn))
			}
		}

		return lines
	}

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_Update", pn)).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				jen.ID("mc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Op("=").ID("mc"),
				jen.Line(),
				jen.ID("r").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Op("=").ID("r"),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.ID("id").Dot("On").Call(jen.Litf("Update%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dotf("%sUpdateInput", sn).Valuesln(
					buildUpdateInputFromExpectedLines()...,
				),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without update input"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusBadRequest")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with no rows fetching %s", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dotf("%sUpdateInput", sn).Valuesln(
					buildUpdateInputFromExpectedLines()...,
				),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusNotFound")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with error fetching %s", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dotf("%sUpdateInput", sn).Valuesln(
					buildUpdateInputFromExpectedLines()...,
				),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusInternalServerError")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with error updating %s", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				jen.ID("mc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"), jen.Qual("github.com/stretchr/testify/mock", "Anything")),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Op("=").ID("mc"),
				jen.Line(),
				jen.ID("r").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Op("=").ID("r"),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.ID("id").Dot("On").Call(jen.Litf("Update%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dotf("%sUpdateInput", sn).Valuesln(
					buildUpdateInputFromExpectedLines()...,
				),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusInternalServerError")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				jen.ID("mc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Increment"), jen.Qual("github.com/stretchr/testify/mock", "Anything")),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Op("=").ID("mc"),
				jen.Line(),
				jen.ID("r").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Op("=").ID("r"),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.ID("id").Dot("On").Call(jen.Litf("Update%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("EncodeResponse"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dotf("%sUpdateInput", sn).Valuesln(
					buildUpdateInputFromExpectedLines()...,
				),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("exampleInput"))),
				jen.Line(),
				jen.ID("s").Dot("UpdateHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestServiceArchiveFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("Test%sService_Archive", pn)).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				jen.ID("r").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Values(),
				jen.ID("r").Dot("On").Call(jen.Lit("Report"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("reporter").Op("=").ID("r"),
				jen.Line(),
				jen.ID("mc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("mc").Dot("On").Call(jen.Lit("Decrement")).Dot("Return").Call(),
				jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Op("=").ID("mc"),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Archive%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("EncodeResponse"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ArchiveHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusNoContent")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with no %s in database", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Archive%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ),
				).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ArchiveHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusNotFound")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				buildRelevantOwnerVar(pkg, typ),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Add(utils.FakeUint64Func()),
				),
				jen.Line(),
				buildRelevantIDFetcher(typ),
				jen.Line(),
				jen.ID("s").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Return().ID("expected").Dot("ID"),
				),
				jen.Line(),
				jen.ID("id").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("id").Dot("On").Call(jen.Litf("Archive%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("expected").Dot("ID"),
					buildDBCallOwnerVar(typ)).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Op("=").ID("id"),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("s").Dot("ArchiveHandler").Call().Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusInternalServerError")),
			)),
		),
		jen.Line(),
	}

	return lines
}
