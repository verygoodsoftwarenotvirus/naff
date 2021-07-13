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
		jen.Func().ID("buildTestService").Params().Params(jen.PointerTo().ID("service")).Body(
			jen.Return().AddressOf().ID("service").Valuesln(jen.ID("logger").Op(":").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				jen.IDf("%sCounter", uvn).Op(":").AddressOf().Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
				jen.IDf("%sDataManager", uvn).Op(":").AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.IDf("%sIDFetcher", uvn).Op(":").Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBody(jen.Return().Lit(0)),
				jen.ID("encoderDecoder").Op(":").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.ID("search").Op(":").AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
					}
					return jen.Null()
				}(),
				jen.ID("tracer").Op(":").Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.Lit("test")))),
		jen.Newline(),
	)

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Var().ID("ucp").Qual(proj.MetricsPackage(), "UnitCounterProvider").Equals().Func().Params(jen.List(jen.ID("counterName"),
			jen.ID("description")).String(),
		).Params(jen.Qual(proj.MetricsPackage(), "UnitCounter")).Body(
			jen.Return().AddressOf().Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
		),
		jen.Newline(),
		jen.ID("rpm").Assign().Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
	}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		tsn := dep.Name.Singular()
		trn := dep.Name.RouteName()

		bodyLines = append(bodyLines,
			jen.ID("rpm").Dot("On").Callln(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.Qual(constants.MockPkg, "IsType").Call(jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call()),
				jen.Qual(proj.ServicePackage(dep.Name.PackageName()), fmt.Sprintf("%sIDURIParamKey", tsn)),
				jen.Lit(trn),
			).Dot("Return").Call(jen.Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBody(jen.Return().Lit(0))),
		)
	}

	bodyLines = append(bodyLines,
		jen.ID("rpm").Dot("On").Callln(
			jen.Lit("BuildRouteParamIDFetcher"),
			jen.Qual(constants.MockPkg, "IsType").Call(jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call()),
			jen.IDf("%sIDURIParamKey", sn),
			jen.Lit(rn),
		).Dot("Return").Call(jen.Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBody(jen.Return().Lit(0))),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.List(jen.ID("s"),
			jen.ID("err")).Assign().ID("ProvideService").Callln(
			jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
			jen.ID("Config").Values(
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.ID("SearchIndexPath").Op(":").Lit("example/path")
					}
					return jen.Null()
				}(),
			),
			jen.AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
			jen.Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
			jen.ID("ucp"),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Func().Params(jen.ID("path").Qual(proj.InternalSearchPackage(), "IndexPath"),
						jen.ID("name").Qual(proj.InternalSearchPackage(), "IndexName"),
						jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
					).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(
							jen.AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
							jen.ID("nil"),
						),
					)
				}
				return jen.Null()
			}(),

			jen.ID("rpm"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "NotNil").Call(
			jen.ID("t"),
			jen.ID("s"),
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("rpm"),
		),
	)

	code.Add(
		jen.Func().IDf("TestProvide%sService", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					bodyLines...,
				),
			),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Newline().ID("T").Dot("Run").Call(
						jen.Lit("with error providing index"),
						jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
							jen.ID("t").Dot("Parallel").Call(),
							jen.Newline(),
							jen.Var().ID("ucp").Qual(proj.MetricsPackage(), "UnitCounterProvider").Equals().Func().Params(jen.List(jen.ID("counterName"),
								jen.ID("description")).String(),
							).Params(jen.Qual(proj.MetricsPackage(), "UnitCounter")).Body(
								jen.Return().AddressOf().Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
							),
							jen.Newline(),
							jen.List(jen.ID("s"),
								jen.ID("err")).Assign().ID("ProvideService").Callln(
								jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
								jen.ID("Config").Values(jen.ID("SearchIndexPath").Op(":").Lit("example/path")),
								jen.AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
								jen.Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
								jen.ID("ucp"),
								jen.Func().Params(jen.ID("path").Qual(proj.InternalSearchPackage(), "IndexPath"),
									jen.ID("name").Qual(proj.InternalSearchPackage(), "IndexName"),
									jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
								).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"),
									jen.ID("error")).Body(
									jen.Return().List(jen.ID("nil"),
										jen.Qual("errors", "New").Call(jen.Lit("blah")))),
								jen.Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
							),
							jen.Newline(),
							jen.Qual(constants.AssertionLibrary, "Nil").Call(
								jen.ID("t"),
								jen.ID("s"),
							),
							jen.Qual(constants.AssertionLibrary, "Error").Call(
								jen.ID("t"),
								jen.ID("err"),
							),
						),
					)
				}
				return jen.Null()
			}(),
		),
	)

	return code
}
