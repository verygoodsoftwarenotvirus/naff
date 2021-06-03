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

	utils.AddImports(proj, code, false)

	code.Add(buildbuildTestServiceFuncDecl(proj, typ)...)
	code.Add(buildTestProvideServiceFuncDecl(proj, typ)...)

	return code
}

func buildbuildTestServiceFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	serviceValues := []jen.Code{
		jen.ID(constants.LoggerVarName).MapAssign().Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
		jen.ID(fmt.Sprintf("%sCounter", uvn)).MapAssign().AddressOf().Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		serviceValues = append(serviceValues,
			jen.ID(fmt.Sprintf("%sDataManager", ot.Name.UnexportedVarName())).MapAssign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", ot.Name.Singular())).Values(),
		)
	}
	serviceValues = append(serviceValues, jen.ID(fmt.Sprintf("%sDataManager", uvn)).MapAssign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values())

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
		jen.ID("encoderDecoder").MapAssign().AddressOf().Qual(proj.EncodingPackage("mock"), "EncoderDecoder").Values(),
		jen.ID("reporter").MapAssign().ID("nil"),
	)

	if typ.SearchEnabled {
		serviceValues = append(serviceValues,
			jen.ID("search").MapAssign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
		)
	}

	lines := []jen.Code{
		jen.Func().ID("buildTestService").Params().Params(jen.PointerTo().ID("Service")).Body(
			jen.Return().AddressOf().ID("Service").Valuesln(serviceValues...),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestProvideServiceFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	provideServiceLines := []jen.Code{
		jen.Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		provideServiceLines = append(provideServiceLines, jen.AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", ot.Name.Singular())).Values())
	}
	provideServiceLines = append(provideServiceLines, jen.AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values())

	for range proj.FindOwnerTypeChain(typ) {
		provideServiceLines = append(provideServiceLines, jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()))
	}
	provideServiceLines = append(provideServiceLines, jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()))
	if typ.OwnedByAUserAtSomeLevel(proj) {
		provideServiceLines = append(provideServiceLines, jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()))
	}

	provideServiceLines = append(provideServiceLines,
		jen.AddressOf().Qual(proj.EncodingPackage("mock"), "EncoderDecoder").Values(),
		jen.ID("ucp"),
		jen.Nil(),
	)

	if typ.SearchEnabled {
		provideServiceLines = append(provideServiceLines,
			jen.AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
		)
	}

	lines := []jen.Code{
		jen.Func().ID(fmt.Sprintf("TestProvide%sService", pn)).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.Var().ID("ucp").Qual(proj.MetricsPackage(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.MetricsPackage(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.MetricsPackage(), "UnitCounter"),
					jen.Error()).Body(
					jen.Return().List(jen.AddressOf().Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(), jen.Nil()),
				),
				jen.Newline(),
				jen.List(jen.ID("s"), jen.Err()).Assign().ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					provideServiceLines...,
				),
				jen.Newline(),
				utils.AssertNotNil(jen.ID("s"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Newline(),
			utils.BuildSubTestWithoutContext(
				"with error providing unit counter",
				jen.Var().ID("ucp").Qual(proj.MetricsPackage(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.MetricsPackage(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.MetricsPackage(), "UnitCounter"), jen.Error()).Body(
					jen.Return().List(jen.Nil(), constants.ObligatoryError()),
				),
				jen.Newline(),
				jen.List(jen.ID("s"), jen.Err()).Assign().ID(fmt.Sprintf("Provide%sService", pn)).Callln(
					provideServiceLines...,
				),
				jen.Newline(),
				utils.AssertNil(jen.ID("s"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Newline(),
	}

	return lines
}
