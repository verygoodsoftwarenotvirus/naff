package iterables

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpHelpersTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("itemsServiceHTTPRoutesTestHelper").Struct(
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("req").Op("*").Qual("net/http", "Request"),
				jen.ID("res").Op("*").ID("httptest").Dot("ResponseRecorder"),
				jen.ID("service").Op("*").ID("service"),
				jen.ID("exampleUser").Op("*").ID("types").Dot("User"),
				jen.ID("exampleAccount").Op("*").ID("types").Dot("Account"),
				jen.ID("exampleItem").Op("*").ID("types").Dot("Item"),
				jen.ID("exampleCreationInput").Op("*").ID("types").Dot("ItemCreationInput"),
				jen.ID("exampleUpdateInput").Op("*").ID("types").Dot("ItemUpdateInput"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("buildTestHelper").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("itemsServiceHTTPRoutesTestHelper")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("helper").Op(":=").Op("&").ID("itemsServiceHTTPRoutesTestHelper").Values(),
			jen.ID("helper").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("helper").Dot("service").Op("=").ID("buildTestService").Call(),
			jen.ID("helper").Dot("exampleUser").Op("=").ID("fakes").Dot("BuildFakeUser").Call(),
			jen.ID("helper").Dot("exampleAccount").Op("=").ID("fakes").Dot("BuildFakeAccount").Call(),
			jen.ID("helper").Dot("exampleAccount").Dot("BelongsToAccount").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
			jen.ID("helper").Dot("exampleItem").Op("=").ID("fakes").Dot("BuildFakeItem").Call(),
			jen.ID("helper").Dot("exampleItem").Dot("BelongsToAccount").Op("=").ID("helper").Dot("exampleAccount").Dot("ID"),
			jen.ID("helper").Dot("exampleCreationInput").Op("=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("helper").Dot("exampleItem")),
			jen.ID("helper").Dot("exampleUpdateInput").Op("=").ID("fakes").Dot("BuildFakeItemUpdateInputFromItem").Call(jen.ID("helper").Dot("exampleItem")),
			jen.ID("helper").Dot("service").Dot("itemIDFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().ID("helper").Dot("exampleItem").Dot("ID")),
			jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Valuesln(jen.ID("helper").Dot("exampleAccount").Dot("ID").Op(":").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()))),
			jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
				jen.ID("logging").Dot("NewNoopLogger").Call(),
				jen.ID("encoding").Dot("ContentTypeJSON"),
			),
			jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
				jen.Return().List(jen.ID("sessionCtxData"), jen.ID("nil"))),
			jen.ID("req").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BuildTestRequest").Call(jen.ID("t")),
			jen.ID("helper").Dot("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
				jen.ID("req").Dot("Context").Call(),
				jen.ID("types").Dot("SessionContextDataKey"),
				jen.ID("sessionCtxData"),
			)),
			jen.ID("helper").Dot("res").Op("=").ID("httptest").Dot("NewRecorder").Call(),
			jen.Return().ID("helper"),
		),
		jen.Newline(),
	)

	return code
}
