package accounts

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
			jen.Return().Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("NewNoopLogger").Call(), jen.ID("accountCounter").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(), jen.ID("accountDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountDataManager").Valuesln(), jen.ID("accountMembershipDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(), jen.ID("accountIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().Lit(0)), jen.ID("encoderDecoder").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("test")))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideAccountsService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("ucp").ID("metrics").Dot("UnitCounterProvider").Op("=").Func().Params(jen.List(jen.ID("counterName"), jen.ID("description")).ID("string")).Params(jen.ID("metrics").Dot("UnitCounter")).Body(
					jen.Return().Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln()),
			),
			jen.ID("l").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.ID("mock").Dot("IsType").Call(jen.ID("l")),
				jen.ID("AccountIDURIParamKey"),
				jen.Lit("account"),
			).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().Lit(0))),
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.ID("mock").Dot("IsType").Call(jen.ID("l")),
				jen.ID("UserIDURIParamKey"),
				jen.Lit("user"),
			).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().Lit(0))),
			jen.ID("s").Op(":=").ID("ProvideService").Call(
				jen.ID("logging").Dot("NewNoopLogger").Call(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountDataManager").Valuesln(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
				jen.ID("ucp"),
				jen.ID("rpm"),
			),
			jen.ID("assert").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("s"),
			),
			jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
				jen.ID("t"),
				jen.ID("rpm"),
			),
		),
		jen.Line(),
	)

	return code
}
