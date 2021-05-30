package admin

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpHelpersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("adminServiceHTTPRoutesTestHelper").Struct(
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("service").Op("*").ID("service"),
				jen.ID("exampleUser").Op("*").ID("types").Dot("User"),
				jen.ID("exampleAccount").Op("*").ID("types").Dot("Account"),
				jen.ID("exampleInput").Op("*").ID("types").Dot("UserReputationUpdateInput"),
				jen.ID("req").Op("*").Qual("net/http", "Request"),
				jen.ID("res").Op("*").ID("httptest").Dot("ResponseRecorder"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("helper").Op("*").ID("adminServiceHTTPRoutesTestHelper")).ID("neuterAdminUser").Params().Body(
			jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("=").Index().ID("string").Valuesln(jen.ID("authorization").Dot("ServiceUserRole").Dot("String").Call()),
			jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
				jen.Return().List(jen.Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Valuesln(jen.ID("helper").Dot("exampleAccount").Dot("ID").Op(":").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()))), jen.ID("nil"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestHelper").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("adminServiceHTTPRoutesTestHelper")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("helper").Op(":=").Op("&").ID("adminServiceHTTPRoutesTestHelper").Valuesln(),
			jen.ID("helper").Dot("service").Op("=").ID("buildTestService").Call(jen.ID("t")),
			jen.Var().Defs(
				jen.ID("err").ID("error"),
			),
			jen.List(jen.ID("helper").Dot("ctx"), jen.ID("err")).Op("=").ID("helper").Dot("service").Dot("sessionManager").Dot("Load").Call(
				jen.Qual("context", "Background").Call(),
				jen.Lit(""),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("helper").Dot("exampleUser").Op("=").ID("fakes").Dot("BuildFakeUser").Call(),
			jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("=").Index().ID("string").Valuesln(jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
			jen.ID("helper").Dot("exampleAccount").Op("=").ID("fakes").Dot("BuildFakeAccount").Call(),
			jen.ID("helper").Dot("exampleAccount").Dot("BelongsToUser").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
			jen.ID("helper").Dot("exampleInput").Op("=").ID("fakes").Dot("BuildFakeUserReputationUpdateInput").Call(),
			jen.ID("helper").Dot("res").Op("=").ID("httptest").Dot("NewRecorder").Call(),
			jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("helper").Dot("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.Lit("https://blah.com"),
				jen.ID("nil"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("helper").Dot("req"),
			),
			jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Valuesln(jen.ID("helper").Dot("exampleAccount").Dot("ID").Op(":").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()))),
			jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
				jen.ID("logging").Dot("NewNoopLogger").Call(),
				jen.ID("encoding").Dot("ContentTypeJSON"),
			),
			jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
				jen.Return().List(jen.ID("sessionCtxData"), jen.ID("nil"))),
			jen.ID("helper").Dot("service").Dot("userIDFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().ID("helper").Dot("exampleUser").Dot("ID")),
			jen.Return().ID("helper"),
		),
		jen.Line(),
	)

	return code
}
