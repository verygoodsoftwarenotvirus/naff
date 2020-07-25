package fake

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)
	utils.AddImports(proj, code)

	code.Add(buildBuildFakeOAuth2Client(proj)...)
	code.Add(buildBuildFakeOAuth2ClientList(proj)...)
	code.Add(buildBuildFakeOAuth2ClientCreationInputFromClient(proj)...)

	return code
}

func buildBuildFakeOAuth2Client(proj *models.Project) []jen.Code {
	funcName := "BuildFakeOAuth2Client"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked OAuth2Client.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
		).Block(
			jen.Return(
				jen.AddressOf().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
					jen.ID("ClientID").MapAssign().Add(utils.FakeUUIDFunc()),
					jen.ID("ClientSecret").MapAssign().Add(utils.FakeUUIDFunc()),
					jen.ID("RedirectURI").MapAssign().Add(utils.FakeURLFunc()),
					jen.ID("Scopes").MapAssign().Index().String().Valuesln(
						utils.FakeStringFunc(),
						utils.FakeStringFunc(),
						utils.FakeStringFunc(),
					),
					jen.ID("ImplicitAllowed").MapAssign().False(),
					jen.ID(constants.UserOwnershipFieldName).MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("CreatedOn").MapAssign().Add(utils.FakeUnixTimeFunc()),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeOAuth2ClientList(proj *models.Project) []jen.Code {
	funcName := "BuildFakeOAuth2ClientList"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked OAuth2ClientList.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList"),
		).Block(
			jen.ID(utils.BuildFakeVarName("OAuth2Client1")).Assign().ID("BuildFakeOAuth2Client").Call(),
			jen.ID(utils.BuildFakeVarName("OAuth2Client2")).Assign().ID("BuildFakeOAuth2Client").Call(),
			jen.ID(utils.BuildFakeVarName("OAuth2Client3")).Assign().ID("BuildFakeOAuth2Client").Call(),
			jen.Line(),
			jen.Return(
				jen.AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().One(),
						jen.ID("Limit").MapAssign().Lit(20),
					),
					jen.ID("Clients").MapAssign().Index().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
						jen.PointerTo().ID("exampleOAuth2Client1"),
						jen.PointerTo().ID("exampleOAuth2Client2"),
						jen.PointerTo().ID("exampleOAuth2Client3"),
					),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeOAuth2ClientCreationInputFromClient(proj *models.Project) []jen.Code {
	funcName := "BuildFakeOAuth2ClientCreationInputFromClient"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked OAuth2ClientCreationInput.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("client").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput"),
		).Block(
			jen.Return(
				jen.AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").MapAssign().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
						jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
						jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
						jen.ID("TOTPToken").MapAssign().Qual("fmt", "Sprintf").Call(jen.Lit(`0%s`), jen.Qual(constants.FakeLibrary, "Zip").Call()),
					),
					jen.ID("Name").MapAssign().ID("client").Dot("Name"),
					jen.ID("Scopes").MapAssign().ID("client").Dot("Scopes"),
					jen.ID("ClientID").MapAssign().ID("client").Dot("ClientID"),
					jen.ID("ClientSecret").MapAssign().ID("client").Dot("ClientSecret"),
					jen.ID("RedirectURI").MapAssign().ID("client").Dot("RedirectURI"),
					jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("client").Dot(constants.UserOwnershipFieldName),
				),
			),
		),
	}

	return lines
}
