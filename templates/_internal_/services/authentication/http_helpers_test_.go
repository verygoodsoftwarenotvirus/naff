package authentication

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpHelpersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("attachCookieToRequestForTest").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("s").Op("*").ID("service"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("user").Op("*").ID("types").Dot("User")).Params(jen.Qual("context", "Context"), jen.Op("*").Qual("net/http", "Request"), jen.ID("string")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
			jen.List(jen.ID("ctx"), jen.ID("sessionErr")).Op(":=").ID("s").Dot("sessionManager").Dot("Load").Call(
				jen.ID("req").Dot("Context").Call(),
				jen.Lit(""),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("sessionErr"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("s").Dot("sessionManager").Dot("RenewToken").Call(jen.ID("ctx")),
			),
			jen.ID("s").Dot("sessionManager").Dot("Put").Call(
				jen.ID("ctx"),
				jen.ID("userIDContextKey"),
				jen.ID("user").Dot("ID"),
			),
			jen.ID("s").Dot("sessionManager").Dot("Put").Call(
				jen.ID("ctx"),
				jen.ID("accountIDContextKey"),
				jen.ID("exampleAccount").Dot("ID"),
			),
			jen.List(jen.ID("token"), jen.ID("_"), jen.ID("err")).Op(":=").ID("s").Dot("sessionManager").Dot("Commit").Call(jen.ID("ctx")),
			jen.ID("assert").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("token"),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot("buildCookie").Call(
				jen.ID("token"),
				jen.Qual("time", "Now").Call().Dot("Add").Call(jen.ID("s").Dot("config").Dot("Cookies").Dot("Lifetime")),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("req").Dot("AddCookie").Call(jen.ID("c")),
			jen.Return().List(jen.ID("ctx"), jen.ID("req").Dot("WithContext").Call(jen.ID("ctx")), jen.ID("token")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("authServiceHTTPRoutesTestHelper").Struct(
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("req").Op("*").Qual("net/http", "Request"),
				jen.ID("res").Op("*").ID("httptest").Dot("ResponseRecorder"),
				jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"),
				jen.ID("service").Op("*").ID("service"),
				jen.ID("exampleUser").Op("*").ID("types").Dot("User"),
				jen.ID("exampleAccount").Op("*").ID("types").Dot("Account"),
				jen.ID("exampleAPIClient").Op("*").ID("types").Dot("APIClient"),
				jen.ID("examplePerms").Map(jen.ID("uint64")).Op("*").ID("types").Dot("UserAccountMembershipInfo"),
				jen.ID("examplePermCheckers").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker"),
				jen.ID("exampleLoginInput").Op("*").ID("types").Dot("UserLoginInput"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("helper").Op("*").ID("authServiceHTTPRoutesTestHelper")).ID("setContextFetcher").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
			jen.ID("helper").Dot("sessionCtxData").Op("=").ID("sessionCtxData"),
			jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
				jen.Return().List(jen.ID("sessionCtxData"), jen.ID("nil"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestHelper").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("authServiceHTTPRoutesTestHelper")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("helper").Op(":=").Op("&").ID("authServiceHTTPRoutesTestHelper").Values(),
			jen.ID("helper").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("helper").Dot("service").Op("=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("helper").Dot("exampleUser").Op("=").ID("fakes").Dot("BuildFakeUser").Call(),
			jen.ID("helper").Dot("exampleAccount").Op("=").ID("fakes").Dot("BuildFakeAccount").Call(),
			jen.ID("helper").Dot("exampleAccount").Dot("BelongsToUser").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
			jen.ID("helper").Dot("exampleAPIClient").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
			jen.ID("helper").Dot("exampleAPIClient").Dot("BelongsToUser").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
			jen.ID("helper").Dot("exampleLoginInput").Op("=").ID("fakes").Dot("BuildFakeUserLoginInputFromUser").Call(jen.ID("helper").Dot("exampleUser")),
			jen.ID("helper").Dot("examplePerms").Op("=").Map(jen.ID("uint64")).Op("*").ID("types").Dot("UserAccountMembershipInfo").Valuesln(jen.ID("helper").Dot("exampleAccount").Dot("ID").Op(":").Valuesln(jen.ID("AccountName").Op(":").ID("helper").Dot("exampleAccount").Dot("Name"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()))),
			jen.ID("helper").Dot("examplePermCheckers").Op("=").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Valuesln(jen.ID("helper").Dot("exampleAccount").Dot("ID").Op(":").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call())),
			jen.ID("helper").Dot("setContextFetcher").Call(jen.ID("t")),
			jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
				jen.ID("logging").Dot("NewNoopLogger").Call(),
				jen.ID("encoding").Dot("ContentTypeJSON"),
			),
			jen.Var().Defs(
				jen.ID("err").ID("error"),
			),
			jen.ID("helper").Dot("res").Op("=").ID("httptest").Dot("NewRecorder").Call(),
			jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("helper").Dot("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
				jen.ID("nil"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("helper").Dot("req"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().ID("helper"),
		),
		jen.Line(),
	)

	return code
}
