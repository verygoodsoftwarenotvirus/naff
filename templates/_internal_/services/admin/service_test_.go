package admin

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
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
			jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.ID("mock").Dot("IsType").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
				jen.ID("UserIDURIParamKey"),
				jen.Lit("user"),
			).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().Lit(0))),
			jen.ID("s").Op(":=").ID("ProvideService").Call(
				jen.ID("logger"),
				jen.Op("&").ID("authentication").Dot("Config").Valuesln(
					jen.ID("Cookies").Op(":").ID("authentication").Dot("CookieConfig").Valuesln(
						jen.ID("SigningKey").Op(":").Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"))),
				jen.Op("&").ID("authentication").Dot("MockAuthenticator").Values(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AdminUserDataManager").Values(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
				jen.ID("scs").Dot("New").Call(),
				jen.ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
					jen.ID("logger"),
					jen.ID("encoding").Dot("ContentTypeJSON"),
				),
				jen.ID("rpm"),
			),
			jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
				jen.ID("t"),
				jen.ID("rpm"),
			),
			jen.List(jen.ID("srv"), jen.ID("ok")).Op(":=").ID("s").Assert(jen.Op("*").ID("service")),
			jen.ID("require").Dot("True").Call(
				jen.ID("t"),
				jen.ID("ok"),
			),
			jen.Return().ID("srv"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideAdminService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.ID("mock").Dot("IsType").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
						jen.ID("UserIDURIParamKey"),
						jen.Lit("user"),
					).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().Lit(0))),
					jen.ID("s").Op(":=").ID("ProvideService").Call(
						jen.ID("logger"),
						jen.Op("&").ID("authentication").Dot("Config").Valuesln(
							jen.ID("Cookies").Op(":").ID("authentication").Dot("CookieConfig").Valuesln(
								jen.ID("SigningKey").Op(":").Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"))),
						jen.Op("&").ID("authentication").Dot("MockAuthenticator").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AdminUserDataManager").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
						jen.ID("scs").Dot("New").Call(),
						jen.ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
							jen.ID("logger"),
							jen.ID("encoding").Dot("ContentTypeJSON"),
						),
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
