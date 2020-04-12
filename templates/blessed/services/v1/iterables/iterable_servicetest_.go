package iterables

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableServiceTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, ret)

	ret.Add(buildbuildTestServiceFuncDecl(proj, typ)...)
	ret.Add(buildTestProvideServiceFuncDecl(proj, typ)...)

	return ret
}

func buildbuildTestServiceFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	serviceValues := []jen.Code{
		jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
		jen.ID(fmt.Sprintf("%sCounter", uvn)).MapAssign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
		jen.ID(fmt.Sprintf("%sDatabase", uvn)).MapAssign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
	}

	if typ.BelongsToUser {
		serviceValues = append(serviceValues,
			jen.ID("userIDFetcher").MapAssign().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
		)
	}
	if typ.BelongsToStruct != nil {
		serviceValues = append(serviceValues,
			jen.IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).MapAssign().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
		)
	}

	serviceValues = append(serviceValues,
		jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).MapAssign().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
		jen.ID("encoderDecoder").MapAssign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("reporter").MapAssign().ID("nil"),
	)

	lines := []jen.Code{
		jen.Func().ID("buildTestService").Params().Params(jen.PointerTo().ID("Service")).Block(
			jen.Return().AddressOf().ID("Service").Valuesln(serviceValues...),
		),
		jen.Line(),
	}

	return lines
}

func relevantIDFetcherParam(typ models.DataType) jen.Code {
	if typ.BelongsToUser || typ.BelongsToStruct != nil {
		return jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero())
	}
	return nil
}

func buildTestProvideServiceFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("TestProvide%sService", pn)).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
					jen.Error()).Block(
					jen.Return().List(jen.AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("s"), jen.Err()).Assign().ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					relevantIDFetcherParam(typ),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Nil(),
				),
				jen.Line(),
				utils.AssertNotNil(jen.ID("s"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error providing unit counter",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"), jen.Error()).Block(
					jen.Return().List(jen.Nil(), utils.ObligatoryError()),
				),
				jen.Line(),
				jen.List(jen.ID("s"), jen.Err()).Assign().ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					relevantIDFetcherParam(typ),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Nil(),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("s"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}
