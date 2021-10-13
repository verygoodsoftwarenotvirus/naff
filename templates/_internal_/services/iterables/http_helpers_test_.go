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

	structLines := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
		jen.ID("req").PointerTo().Qual("net/http", "Request"),
		jen.ID("res").PointerTo().ID("httptest").Dot("ResponseRecorder"),
		jen.ID("service").PointerTo().ID("service"),
		jen.ID("exampleUser").PointerTo().Qual(proj.TypesPackage(), "User"),
		jen.ID("exampleAccount").PointerTo().Qual(proj.TypesPackage(), "Account"),
	}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		structLines = append(structLines, jen.IDf("example%s", dep.Name.Singular()).PointerTo().Qual(proj.TypesPackage(), dep.Name.Singular()))
	}

	structLines = append(structLines,
		jen.IDf("example%s", sn).PointerTo().Qual(proj.TypesPackage(), sn),
		jen.ID("exampleCreationInput").PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationRequestInput", sn)),
		jen.ID("exampleUpdateInput").PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn)),
	)

	code.Add(
		jen.Type().IDf("%sServiceHTTPRoutesTestHelper", puvn).Struct(
			structLines...,
		),
		jen.Newline(),
	)

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Helper").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().AddressOf().IDf("%sServiceHTTPRoutesTestHelper", puvn).Values(),
		jen.Newline(),
		jen.ID("helper").Dot("ctx").Equals().Qual("context", "Background").Call(),
		jen.ID("helper").Dot("service").Equals().ID("buildTestService").Call(),
		jen.ID("helper").Dot("exampleUser").Equals().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
		jen.ID("helper").Dot("exampleAccount").Equals().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
		jen.ID("helper").Dot("exampleAccount").Dot("BelongsToUser").Equals().ID("helper").Dot("exampleUser").Dot("ID"),
	}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		tsn := dep.Name.Singular()

		bodyLines = append(bodyLines,
			jen.ID("helper").Dotf("example%s", tsn).Equals().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", tsn)).Call(),
		)

		if dep.BelongsToStruct != nil {
			btssn := dep.BelongsToStruct.Singular()
			bodyLines = append(bodyLines, jen.ID("helper").Dotf("example%s", tsn).Dotf("BelongsTo%s", btssn).Equals().ID("helper").Dotf("example%s", btssn).Dot("ID"))
		}

		if dep.BelongsToAccount {
			bodyLines = append(bodyLines, jen.ID("helper").Dotf("example%s", tsn).Dot("BelongsToAccount").Equals().ID("helper").Dot("exampleAccount").Dot("ID"))
		}
	}

	bodyLines = append(bodyLines,
		jen.ID("helper").Dotf("example%s", sn).Equals().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
	)

	if typ.BelongsToStruct != nil {
		bodyLines = append(bodyLines, jen.ID("helper").Dotf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID("helper").Dotf("example%s", typ.BelongsToStruct.Singular()).Dot("ID"))
	}

	if typ.BelongsToAccount {
		bodyLines = append(bodyLines, jen.ID("helper").Dotf("example%s", sn).Dot("BelongsToAccount").Equals().ID("helper").Dot("exampleAccount").Dot("ID"))
	}

	bodyLines = append(bodyLines,
		jen.ID("helper").Dot("exampleCreationInput").Equals().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationRequestInputFrom%s", sn, sn)).Call(jen.ID("helper").Dotf("example%s", sn)),
		jen.ID("helper").Dot("exampleUpdateInput").Equals().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInputFrom%s", sn, sn)).Call(jen.ID("helper").Dotf("example%s", sn)),
		jen.Newline(),
	)

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		tsn := dep.Name.Singular()
		tuvn := dep.Name.UnexportedVarName()

		bodyLines = append(bodyLines,
			jen.ID("helper").Dot("service").Dotf("%sIDFetcher", tuvn).Equals().Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()).Body(
				jen.Return().ID("helper").Dotf("example%s", tsn).Dot("ID"),
			),
			jen.Newline(),
		)
	}

	bodyLines = append(bodyLines,
		jen.ID("helper").Dot("service").Dotf("%sIDFetcher", uvn).Equals().Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()).Body(
			jen.Return().ID("helper").Dotf("example%s", sn).Dot("ID"),
		),
		jen.Newline(),
	)

	bodyLines = append(bodyLines,
		jen.ID("sessionCtxData").Assign().AddressOf().Qual(proj.TypesPackage(), "SessionContextData").Valuesln(
			jen.ID("Requester").MapAssign().Qual(proj.TypesPackage(), "RequesterInfo").Valuesln(
				jen.ID("UserID").MapAssign().ID("helper").Dot("exampleUser").Dot("ID"),
				jen.ID("Reputation").MapAssign().ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"),
				jen.ID("ReputationExplanation").MapAssign().ID("helper").Dot("exampleUser").Dot("ReputationExplanation"),
				jen.ID("ServicePermissions").MapAssign().Qual(proj.InternalAuthorizationPackage(), "NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("...")),
			),
			jen.ID("ActiveAccountID").MapAssign().ID("helper").Dot("exampleAccount").Dot("ID"),
			jen.ID("AccountPermissions").MapAssign().Map(jen.String()).Qual(proj.InternalAuthorizationPackage(), "AccountRolePermissionsChecker").Valuesln(jen.ID("helper").Dot("exampleAccount").Dot("ID").MapAssign().Qual(proj.InternalAuthorizationPackage(), "NewAccountRolePermissionChecker").Call(jen.Qual(proj.InternalAuthorizationPackage(), "AccountMemberRole").Dot("String").Call()))),
		jen.Newline(),
		jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
			jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
			jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
		),
		jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Equals().Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"),
			jen.ID("error")).Body(
			jen.Return().List(jen.ID("sessionCtxData"),
				jen.Nil()),
		),
		jen.Newline(),
		jen.ID("req").Assign().Qual(proj.TestUtilsPackage(), "BuildTestRequest").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("helper").Dot("req").Equals().ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
			jen.ID("req").Dot("Context").Call(),
			jen.Qual(proj.TypesPackage(), "SessionContextDataKey"),
			jen.ID("sessionCtxData"),
		)),
		jen.Newline(),
		jen.ID("helper").Dot("res").Equals().ID("httptest").Dot("NewRecorder").Call(),
		jen.Newline(),
		jen.Return().ID("helper"),
	)

	code.Add(
		jen.Func().ID("buildTestHelper").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().IDf("%sServiceHTTPRoutesTestHelper", puvn)).Body(
			bodyLines...,
		),
		jen.Newline(),
	)

	return code
}
