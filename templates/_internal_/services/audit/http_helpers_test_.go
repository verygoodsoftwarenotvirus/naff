package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpHelpersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("auditServiceHTTPRoutesTestHelper").Struct(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("req").Op("*").Qual("net/http", "Request"),
			jen.ID("res").Op("*").ID("httptest").Dot("ResponseRecorder"),
			jen.ID("service").Op("*").ID("service"),
			jen.ID("exampleUser").Op("*").ID("types").Dot("User"),
			jen.ID("exampleAccount").Op("*").ID("types").Dot("Account"),
			jen.ID("exampleAuditLogEntry").Op("*").ID("types").Dot("AuditLogEntry"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestHelper").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("auditServiceHTTPRoutesTestHelper")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("helper").Op(":=").Op("&").ID("auditServiceHTTPRoutesTestHelper").Values(),
			jen.ID("helper").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("helper").Dot("service").Op("=").ID("buildTestService").Call(),
			jen.ID("helper").Dot("exampleUser").Op("=").ID("fakes").Dot("BuildFakeUser").Call(),
			jen.ID("helper").Dot("exampleAccount").Op("=").ID("fakes").Dot("BuildFakeAccount").Call(),
			jen.ID("helper").Dot("exampleAccount").Dot("BelongsToUser").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
			jen.ID("helper").Dot("exampleAuditLogEntry").Op("=").ID("fakes").Dot("BuildFakeAuditLogEntry").Call(),
			jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(
				jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(
					jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Valuesln(
					jen.ID("helper").Dot("exampleAccount").Dot("ID").Op(":").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()))),
			jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
				jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
				jen.ID("encoding").Dot("ContentTypeJSON"),
			),
			jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
				jen.Return().List(jen.ID("sessionCtxData"), jen.ID("nil"))),
			jen.ID("helper").Dot("service").Dot("auditLogEntryIDFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().ID("helper").Dot("exampleAuditLogEntry").Dot("ID")),
			jen.Var().ID("err").ID("error"),
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
