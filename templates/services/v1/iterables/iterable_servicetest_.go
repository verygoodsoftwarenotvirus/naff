package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableServiceTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code)

	code.Add(buildbuildTestServiceFuncDecl(proj, typ)...)
	code.Add(buildTestProvideServiceFuncDecl(proj, typ)...)

	return code
}

func buildbuildTestServiceFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	serviceValues := []jen.Code{
		jen.ID(constants.LoggerVarName).MapAssign().Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
		jen.ID(fmt.Sprintf("%sCounter", uvn)).MapAssign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		serviceValues = append(serviceValues,
			jen.ID(fmt.Sprintf("%sDataManager", ot.Name.UnexportedVarName())).MapAssign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", ot.Name.Singular())).Values(),
		)
	}
	serviceValues = append(serviceValues, jen.ID(fmt.Sprintf("%sDataManager", uvn)).MapAssign().AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values())

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		serviceValues = append(serviceValues,
			jen.IDf("%sIDFetcher", ot.Name.UnexportedVarName()).MapAssign().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
		)
	}
	serviceValues = append(serviceValues, jen.IDf("%sIDFetcher", uvn).MapAssign().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()))

	if typ.OwnedByAUserAtSomeLevel(proj) {
		serviceValues = append(serviceValues,
			jen.ID("userIDFetcher").MapAssign().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
		)
	}

	serviceValues = append(serviceValues,
		jen.ID("encoderDecoder").MapAssign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("reporter").MapAssign().ID("nil"),
	)

	if typ.SearchEnabled {
		serviceValues = append(serviceValues,
			jen.ID("search").MapAssign().AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values(),
		)
	}

	lines := []jen.Code{
		jen.Func().ID("buildTestService").Params().Params(jen.PointerTo().ID("Service")).Block(
			jen.Return().AddressOf().ID("Service").Valuesln(serviceValues...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideServiceFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	provideServiceLines := []jen.Code{
		jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		provideServiceLines = append(provideServiceLines, jen.AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", ot.Name.Singular())).Values())
	}
	provideServiceLines = append(provideServiceLines, jen.AddressOf().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataManager", sn)).Values())

	for range proj.FindOwnerTypeChain(typ) {
		provideServiceLines = append(provideServiceLines, jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()))
	}
	provideServiceLines = append(provideServiceLines, jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()))
	if typ.OwnedByAUserAtSomeLevel(proj) {
		provideServiceLines = append(provideServiceLines, jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()))
	}

	provideServiceLines = append(provideServiceLines,
		jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		jen.ID("ucp"),
		jen.Nil(),
	)

	if typ.SearchEnabled {
		provideServiceLines = append(provideServiceLines,
			jen.AddressOf().Qual(proj.InternalSearchV1Package("mock"), "IndexManager").Values(),
		)
	}

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
					provideServiceLines...,
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
					jen.Return().List(jen.Nil(), constants.ObligatoryError()),
				),
				jen.Line(),
				jen.List(jen.ID("s"), jen.Err()).Assign().ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					provideServiceLines...,
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
