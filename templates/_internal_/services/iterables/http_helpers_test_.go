package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpHelpersTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()

	code.Add(
		jen.Type().IDf("%sServiceHTTPRoutesTestHelper", puvn).Struct(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("req").Op("*").Qual("net/http", "Request"),
			jen.ID("res").Op("*").ID("httptest").Dot("ResponseRecorder"),
			jen.ID("service").Op("*").ID("service"),
			jen.ID("exampleUser").Op("*").Qual(proj.TypesPackage(), "User"),
			jen.ID("exampleAccount").Op("*").Qual(proj.TypesPackage(), "Account"),
			jen.IDf("example%s", sn).Op("*").Qual(proj.TypesPackage(), sn),
			jen.ID("exampleCreationInput").Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)),
			jen.ID("exampleUpdateInput").Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn)),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("buildTestHelper").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").IDf("%sServiceHTTPRoutesTestHelper", puvn)).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Newline(),
			jen.ID("helper").Op(":=").Op("&").IDf("%sServiceHTTPRoutesTestHelper", puvn).Values(),
			jen.Newline(),
			jen.ID("helper").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("helper").Dot("service").Op("=").ID("buildTestService").Call(),
			jen.ID("helper").Dot("exampleUser").Op("=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
			jen.ID("helper").Dot("exampleAccount").Op("=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
			jen.ID("helper").Dot("exampleAccount").Dot("BelongsToUser").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
			jen.ID("helper").Dotf("example%s", sn).Op("=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
			jen.ID("helper").Dotf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("helper").Dot("exampleAccount").Dot("ID"),
			jen.ID("helper").Dot("exampleCreationInput").Op("=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.ID("helper").Dotf("example%s", sn)),
			jen.ID("helper").Dot("exampleUpdateInput").Op("=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.ID("helper").Dotf("example%s", sn)),
			jen.Newline(),
			jen.ID("helper").Dot("service").Dotf("%sIDFetcher", uvn).Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.Return().ID("helper").Dotf("example%s", sn).Dot("ID"),
			),
			jen.Newline(),
			jen.ID("sessionCtxData").Op(":=").Op("&").Qual(proj.TypesPackage(), "SessionContextData").Valuesln(
				jen.ID("Requester").Op(":").Qual(proj.TypesPackage(), "RequesterInfo").Valuesln(
					jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"),
					jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"),
					jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"),
					jen.ID("ServicePermissions").Op(":").Qual(proj.InternalAuthorizationPackage(), "NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("...")),
				),
				jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"),
				jen.ID("AccountPermissions").Op(":").Map(jen.ID("uint64")).Qual(proj.InternalAuthorizationPackage(), "AccountRolePermissionsChecker").Valuesln(jen.ID("helper").Dot("exampleAccount").Dot("ID").Op(":").Qual(proj.InternalAuthorizationPackage(), "NewAccountRolePermissionChecker").Call(jen.Qual(proj.InternalAuthorizationPackage(), "AccountMemberRole").Dot("String").Call()))),
			jen.Newline(),
			jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
				jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
			),
			jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"),
				jen.ID("error")).Body(
				jen.Return().List(jen.ID("sessionCtxData"),
					jen.ID("nil")),
			),
			jen.Newline(),
			jen.ID("req").Op(":=").Qual(proj.TestUtilsPackage(), "BuildTestRequest").Call(jen.ID("t")),
			jen.Newline(),
			jen.ID("helper").Dot("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
				jen.ID("req").Dot("Context").Call(),
				jen.Qual(proj.TypesPackage(), "SessionContextDataKey"),
				jen.ID("sessionCtxData"),
			)),
			jen.Newline(),
			jen.ID("helper").Dot("res").Op("=").ID("httptest").Dot("NewRecorder").Call(),
			jen.Newline(),
			jen.Return().ID("helper"),
		),
		jen.Newline(),
	)

	return code
}
