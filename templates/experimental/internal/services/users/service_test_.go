package users

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
			jen.ID("expectedUserCount").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
			jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
				jen.Lit("GetAllUsersCount"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
			).Dot("Return").Call(
				jen.ID("expectedUserCount"),
				jen.ID("nil"),
			),
			jen.ID("s").Op(":=").ID("ProvideUsersService").Call(
				jen.Op("&").ID("authentication").Dot("Config").Valuesln(),
				jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountDataManager").Valuesln(),
				jen.Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
				jen.Func().Params(jen.List(jen.ID("counterName"), jen.ID("description")).ID("string")).Params(jen.ID("metrics").Dot("UnitCounter")).Body(
					jen.Return().ID("uc")),
				jen.Op("&").ID("images").Dot("MockImageUploadProcessor").Valuesln(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/uploads/mock", "UploadManager").Valuesln(),
				jen.ID("chi").Dot("NewRouteParamManager").Call(),
			),
			jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
				jen.ID("t"),
				jen.ID("mockDB"),
				jen.ID("uc"),
			),
			jen.Return().ID("s").Assert(jen.Op("*").ID("service")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideUsersService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.ID("mock").Dot("IsType").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
						jen.ID("UserIDURIParamKey"),
						jen.Lit("user"),
					).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().Lit(0))),
					jen.ID("s").Op(":=").ID("ProvideUsersService").Call(
						jen.Op("&").ID("authentication").Dot("Config").Valuesln(),
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountDataManager").Valuesln(),
						jen.Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
						jen.Func().Params(jen.List(jen.ID("counterName"), jen.ID("description")).ID("string")).Params(jen.ID("metrics").Dot("UnitCounter")).Body(
							jen.Return().Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln()),
						jen.Op("&").ID("images").Dot("MockImageUploadProcessor").Valuesln(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/uploads/mock", "UploadManager").Valuesln(),
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
			),
		),
		jen.Line(),
	)

	return code
}
