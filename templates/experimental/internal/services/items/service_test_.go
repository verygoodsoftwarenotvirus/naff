package items

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildTestService").Params().Params(jen.Op("*").ID("service")).Body(
			jen.Return().Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("NewNoopLogger").Call(), jen.ID("itemCounter").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(), jen.ID("itemDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(), jen.ID("itemIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().Lit(0)), jen.ID("encoderDecoder").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(), jen.ID("search").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("test")))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideItemsService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Var().Defs(
						jen.ID("ucp").ID("metrics").Dot("UnitCounterProvider").Op("=").Func().Params(jen.List(jen.ID("counterName"), jen.ID("description")).ID("string")).Params(jen.ID("metrics").Dot("UnitCounter")).Body(
							jen.Return().Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln()),
					),
					jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.ID("mock").Dot("IsType").Call(jen.ID("logging").Dot("NewNoopLogger").Call()),
						jen.ID("ItemIDURIParamKey"),
						jen.Lit("item"),
					).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().Lit(0))),
					jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideService").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("Config").Valuesln(jen.ID("SearchIndexPath").Op(":").Lit("example/path")),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
						jen.ID("ucp"),
						jen.Func().Params(jen.ID("path").ID("search").Dot("IndexPath"), jen.ID("name").ID("search").Dot("IndexName"), jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
							jen.Return().List(jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(), jen.ID("nil"))),
						jen.ID("rpm"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("rpm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error providing index"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Var().Defs(
						jen.ID("ucp").ID("metrics").Dot("UnitCounterProvider").Op("=").Func().Params(jen.List(jen.ID("counterName"), jen.ID("description")).ID("string")).Params(jen.ID("metrics").Dot("UnitCounter")).Body(
							jen.Return().Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln()),
					),
					jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideService").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("Config").Valuesln(jen.ID("SearchIndexPath").Op(":").Lit("example/path")),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
						jen.ID("ucp"),
						jen.Func().Params(jen.ID("path").ID("search").Dot("IndexPath"), jen.ID("name").ID("search").Dot("IndexName"), jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
							jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					),
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
		jen.Line(),
	)

	return code
}
