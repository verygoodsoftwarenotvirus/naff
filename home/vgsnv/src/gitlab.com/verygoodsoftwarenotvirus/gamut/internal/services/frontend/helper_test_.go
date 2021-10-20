package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helperTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("serviceHTTPRoutesTestHelper").Struct(
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("req").Op("*").Qual("net/http", "Request"),
				jen.ID("res").Op("*").ID("httptest").Dot("ResponseRecorder"),
				jen.ID("service").Op("*").ID("service"),
				jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"),
				jen.ID("exampleUser").Op("*").ID("types").Dot("User"),
				jen.ID("exampleAccount").Op("*").ID("types").Dot("Account"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("buildTestHelper").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("serviceHTTPRoutesTestHelper")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("helper").Op(":=").Op("&").ID("serviceHTTPRoutesTestHelper").Valuesln(),
			jen.ID("helper").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("helper").Dot("exampleUser").Op("=").ID("fakes").Dot("BuildFakeUser").Call(),
			jen.ID("helper").Dot("exampleAccount").Op("=").ID("fakes").Dot("BuildFakeAccount").Call(),
			jen.ID("helper").Dot("exampleAccount").Dot("BelongsToUser").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
			jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(),
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("authService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/gamut/pkg/types/mock", "AuthService").Valuesln(),
			jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/gamut/internal/routing/mock", "NewRouteParamManager").Call(),
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.ID("mock").Dot("IsType").Call(jen.ID("logging").Dot("NewNoopLogger").Call()),
				jen.ID("mock").Dot("AnythingOfType").Call(jen.Lit("string")),
				jen.ID("mock").Dot("AnythingOfType").Call(jen.Lit("string")),
			).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().Lit(0))),
			jen.Var().Defs(
				jen.ID("ok").ID("bool"),
			),
			jen.List(jen.ID("helper").Dot("service"), jen.ID("ok")).Op("=").ID("ProvideService").Call(
				jen.ID("cfg"),
				jen.ID("logger"),
				jen.ID("authService"),
			).Assert(jen.Op("*").ID("service")),
			jen.ID("require").Dot("True").Call(
				jen.ID("t"),
				jen.ID("ok"),
			),
			jen.ID("helper").Dot("sessionCtxData").Op("=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").Map(jen.ID("string")).ID("authorization").Dot("AccountRolePermissionsChecker").Valuesln(jen.ID("helper").Dot("exampleAccount").Dot("ID").Op(":").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()))),
			jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
				jen.Return().List(jen.ID("helper").Dot("sessionCtxData"), jen.ID("nil"))),
			jen.ID("req").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/gamut/tests/utils", "BuildTestRequest").Call(jen.ID("t")),
			jen.ID("helper").Dot("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
				jen.ID("req").Dot("Context").Call(),
				jen.ID("types").Dot("SessionContextDataKey"),
				jen.ID("helper").Dot("sessionCtxData"),
			)),
			jen.ID("helper").Dot("res").Op("=").ID("httptest").Dot("NewRecorder").Call(),
			jen.Return().ID("helper"),
		),
		jen.Newline(),
	)

	return code
}
