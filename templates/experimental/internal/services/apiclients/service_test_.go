package apiclients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("service")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Return().Op("&").ID("service").Valuesln(jen.ID("apiClientDataManager").Op(":").ID("database").Dot("BuildMockDatabase").Call(), jen.ID("logger").Op(":").ID("logging").Dot("NewNonOperationalLogger").Call(), jen.ID("encoderDecoder").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(), jen.ID("authenticator").Op(":").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(), jen.ID("sessionContextDataFetcher").Op(":").ID("authentication").Dot("FetchContextFromRequest"), jen.ID("urlClientIDExtractor").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().Lit(0)), jen.ID("apiClientCounter").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(), jen.ID("secretGenerator").Op(":").Op("&").ID("random").Dot("MockGenerator").Valuesln(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")), jen.ID("cfg").Op(":").Op("&").ID("config").Valuesln()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideAPIClientsService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("mockAPIClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.ID("mock").Dot("IsType").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
						jen.ID("APIClientIDURIParamKey"),
						jen.Lit("api client"),
					).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().Lit(0))),
					jen.ID("s").Op(":=").ID("ProvideAPIClientsService").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("mockAPIClientDataManager"),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
						jen.Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
						jen.Func().Params(jen.List(jen.ID("counterName"), jen.ID("description")).ID("string")).Params(jen.ID("metrics").Dot("UnitCounter")).Body(
							jen.Return().ID("nil")),
						jen.ID("rpm"),
						jen.Op("&").ID("config").Valuesln(),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIClientDataManager"),
						jen.ID("rpm"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
