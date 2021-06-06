package iterables

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	rn := typ.Name.RouteName()
	uvn := typ.Name.UnexportedVarName()

	code.Add(
		jen.Func().ID("buildTestService").Params().Params(jen.Op("*").ID("service")).Body(
			jen.Return().Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("NewNoopLogger").Call(),
				jen.IDf("%sCounter", uvn).Op(":").Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
				jen.IDf("%sDataManager", uvn).Op(":").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.IDf("%sIDFetcher", uvn).Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBody(jen.Return().Lit(0)),
				jen.ID("encoderDecoder").Op(":").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
				jen.ID("search").Op(":").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
				jen.ID("tracer").Op(":").Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.Lit("test")))),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestProvide%sService", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.Var().ID("ucp").Qual(proj.MetricsPackage(), "UnitCounterProvider").Op("=").Func().Params(jen.List(jen.ID("counterName"),
						jen.ID("description")).ID("string"),
					).Params(jen.Qual(proj.MetricsPackage(), "UnitCounter")).Body(
						jen.Return().Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call()),
						jen.IDf("%sIDURIParamKey", sn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBody(jen.Return().Lit(0))),
					jen.Newline(),
					jen.List(jen.ID("s"),
						jen.ID("err")).Op(":=").ID("ProvideService").Callln(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("Config").Values(jen.ID("SearchIndexPath").Op(":").Lit("example/path")),
						jen.Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
						jen.Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
						jen.ID("ucp"),
						jen.Func().Params(jen.ID("path").Qual(proj.InternalSearchPackage(), "IndexPath"),
							jen.ID("name").Qual(proj.InternalSearchPackage(), "IndexName"),
							jen.ID("logger").ID("logging").Dot("Logger"),
						).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
							jen.Return().List(jen.Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
								jen.ID("nil"))),
						jen.ID("rpm"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("rpm"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error providing index"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.Var().ID("ucp").Qual(proj.MetricsPackage(), "UnitCounterProvider").Op("=").Func().Params(jen.List(jen.ID("counterName"),
						jen.ID("description")).ID("string"),
					).Params(jen.Qual(proj.MetricsPackage(), "UnitCounter")).Body(
						jen.Return().Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
					),
					jen.Newline(),
					jen.List(jen.ID("s"),
						jen.ID("err")).Op(":=").ID("ProvideService").Callln(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("Config").Values(jen.ID("SearchIndexPath").Op(":").Lit("example/path")),
						jen.Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
						jen.Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
						jen.ID("ucp"),
						jen.Func().Params(jen.ID("path").Qual(proj.InternalSearchPackage(), "IndexPath"),
							jen.ID("name").Qual(proj.InternalSearchPackage(), "IndexName"),
							jen.ID("logger").ID("logging").Dot("Logger"),
						).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"),
							jen.ID("error")).Body(
							jen.Return().List(jen.ID("nil"),
								jen.Qual("errors", "New").Call(jen.Lit("blah")))),
						jen.Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
