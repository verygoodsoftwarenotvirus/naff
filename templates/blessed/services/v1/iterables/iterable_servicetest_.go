package iterables

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableServiceTestDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(pkg, ret)

	ret.Add(utils.FakeSeedFunc())

	ret.Add(buildbuildTestServiceFuncDecl(pkg, typ)...)
	ret.Add(buildTestProvideServiceFuncDecl(pkg, typ)...)

	return ret
}

func buildbuildTestServiceFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	serviceValues := []jen.Code{
		jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
		jen.ID(fmt.Sprintf("%sCounter", uvn)).Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
		jen.ID(fmt.Sprintf("%sDatabase", uvn)).Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
	}

	if typ.BelongsToUser {
		serviceValues = append(serviceValues,
			jen.ID("userIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
		)
	}
	if typ.BelongsToStruct != nil {
		serviceValues = append(serviceValues,
			jen.IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
		)
	}

	serviceValues = append(serviceValues,
		jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
		jen.ID("encoderDecoder").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
		jen.ID("reporter").Op(":").ID("nil"),
	)

	lines := []jen.Code{
		jen.Func().ID("buildTestService").Params().Params(jen.Op("*").ID("Service")).Block(
			jen.Return().Op("&").ID("Service").Valuesln(serviceValues...),
		),
		jen.Line(),
	}

	return lines
}

func relevantIDFetcherParam(typ models.DataType) jen.Code {
	if typ.BelongsToUser || typ.BelongsToStruct != nil {
		return jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0))
	}
	return nil
}

func buildTestProvideServiceFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	cn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("TestProvide%sService", pn)).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expectation").Op(":=").Add(utils.FakeUint64Func()),
				jen.ID("uc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
					jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.Nil()),
				),
				jen.Line(),
				jen.ID("idm").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("idm").Dot("On").Call(jen.Lit(fmt.Sprintf("GetAll%sCount", pn)), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("s"), jen.Err()).Op(":=").ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					utils.CtxVar(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("idm"),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					relevantIDFetcherParam(typ),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Nil(),
				),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("s")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error providing unit counter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expectation").Op(":=").Add(utils.FakeUint64Func()),
				jen.ID("uc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
					jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				),
				jen.Line(),
				jen.ID("idm").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("idm").Dot("On").Call(jen.Lit(fmt.Sprintf("GetAll%sCount", pn)), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("s"), jen.Err()).Op(":=").ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					utils.CtxVar(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("idm"),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					relevantIDFetcherParam(typ),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Nil(),
				),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/require", "Nil").Call(jen.ID("t"), jen.ID("s")),
				jen.Qual("github.com/stretchr/testify/require", "Error").Call(jen.ID("t"), jen.Err()),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit(fmt.Sprintf("with error fetching %s count", cn)), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expectation").Op(":=").Add(utils.FakeUint64Func()),
				jen.ID("uc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
					jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.Nil()),
				),
				jen.Line(),
				jen.ID("idm").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("idm").Dot("On").Call(jen.Lit(fmt.Sprintf("GetAll%sCount", pn)), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("s"), jen.Err()).Op(":=").ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					utils.CtxVar(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("idm"),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					relevantIDFetcherParam(typ),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Nil(),
				),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/require", "Nil").Call(jen.ID("t"), jen.ID("s")),
				jen.Qual("github.com/stretchr/testify/require", "Error").Call(jen.ID("t"), jen.Err()),
			)),
		),
		jen.Line(),
	}

	return lines
}
