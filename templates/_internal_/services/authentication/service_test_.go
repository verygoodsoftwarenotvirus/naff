package authentication

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("service")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("encoderDecoder").Op(":=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
				jen.ID("logger"),
				jen.ID("encoding").Dot("ContentTypeJSON"),
			),
			jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideService").Call(
				jen.ID("logger"),
				jen.Op("&").ID("Config").Valuesln(jen.ID("Cookies").Op(":").ID("CookieConfig").Valuesln(jen.ID("Name").Op(":").ID("DefaultCookieName"), jen.ID("SigningKey").Op(":").Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!")), jen.ID("PASETO").Op(":").ID("PASETOConfig").Valuesln(jen.ID("Issuer").Op(":").Lit("test"), jen.ID("LocalModeKey").Op(":").Index().ID("byte").Call(jen.Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!")), jen.ID("Lifetime").Op(":").Qual("time", "Hour"))),
				jen.Op("&").ID("authentication").Dot("MockAuthenticator").Values(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Values(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Values(),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Values(),
				jen.ID("scs").Dot("New").Call(),
				jen.ID("encoderDecoder"),
				jen.ID("chi").Dot("NewRouteParamManager").Call(),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().ID("s").Assert(jen.Op("*").ID("service")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logger"),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideService").Call(
						jen.ID("logger"),
						jen.Op("&").ID("Config").Valuesln(jen.ID("Cookies").Op(":").ID("CookieConfig").Valuesln(jen.ID("Name").Op(":").ID("DefaultCookieName"), jen.ID("SigningKey").Op(":").Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"))),
						jen.Op("&").ID("authentication").Dot("MockAuthenticator").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Values(),
						jen.ID("scs").Dot("New").Call(),
						jen.ID("encoderDecoder"),
						jen.ID("chi").Dot("NewRouteParamManager").Call(),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid cookie key"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logger"),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideService").Call(
						jen.ID("logger"),
						jen.Op("&").ID("Config").Valuesln(jen.ID("Cookies").Op(":").ID("CookieConfig").Valuesln(jen.ID("Name").Op(":").ID("DefaultCookieName"), jen.ID("SigningKey").Op(":").Lit("BLAHBLAHBLAH"))),
						jen.Op("&").ID("authentication").Dot("MockAuthenticator").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Values(),
						jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Values(),
						jen.ID("scs").Dot("New").Call(),
						jen.ID("encoderDecoder"),
						jen.ID("chi").Dot("NewRouteParamManager").Call(),
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
