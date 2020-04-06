package fake

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)
	utils.AddImports(proj, ret)

	ret.Add(buildBuildFakeUser(proj)...)
	ret.Add(buildBuildDatabaseCreationResponse(proj)...)
	ret.Add(buildBuildFakeUserList(proj)...)
	ret.Add(buildBuildFakeUserCreationInput(proj)...)
	ret.Add(buildBuildFakeUserCreationInputFromUser(proj)...)
	ret.Add(buildBuildFakeUserDatabaseCreationInputFromUser(proj)...)
	ret.Add(buildBuildFakeUserLoginInputFromUser(proj)...)
	ret.Add(buildBuildFakePasswordUpdateInput(proj)...)
	ret.Add(buildBuildFakeTOTPSecretRefreshInput(proj)...)

	return ret
}

func buildBuildFakeUser(proj *models.Project) []jen.Code {
	funcName := "BuildFakeUser"
	typeName := "User"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.VarPointer().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("ID").MapAssign().Uint64().Call(utils.FakeUint32Func()),
					jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
					jen.Comment(`HashedPassword: ""`),
					jen.Comment("Salt:           []byte(fake.Word())"),
					jen.Comment(`TwoFactorSecret: ""`),
					jen.ID("IsAdmin").MapAssign().False(),
					jen.ID("CreatedOn").MapAssign().Add(utils.FakeUnixTimeFunc()),
				),
			),
		),
	}

	return lines
}

func buildBuildDatabaseCreationResponse(proj *models.Project) []jen.Code {
	funcName := "BuildDatabaseCreationResponse"
	typeName := "UserCreationResponse"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.VarPointer().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("ID").MapAssign().ID("user").Dot("ID"),
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("TwoFactorSecret").MapAssign().ID("user").Dot("TwoFactorSecret"),
					jen.ID("PasswordLastChangedOn").MapAssign().ID("user").Dot("PasswordLastChangedOn"),
					jen.ID("IsAdmin").MapAssign().ID("user").Dot("IsAdmin"),
					jen.ID("CreatedOn").MapAssign().ID("user").Dot("CreatedOn"),
					jen.ID("UpdatedOn").MapAssign().ID("user").Dot("UpdatedOn"),
					jen.ID("ArchivedOn").MapAssign().ID("user").Dot("ArchivedOn"),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeUserList(proj *models.Project) []jen.Code {
	funcName := "BuildFakeUserList"
	typeName := "UserList"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.ID("exampleUser1").Assign().ID("BuildFakeUser").Call(),
			jen.ID("exampleUser2").Assign().ID("BuildFakeUser").Call(),
			jen.ID("exampleUser3").Assign().ID("BuildFakeUser").Call(),
			jen.Line(),
			jen.Return(
				jen.VarPointer().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().One(),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().Lit(3),
					),
					jen.ID("Users").MapAssign().Index().Qual(proj.ModelsV1Package(), "User").Valuesln(
						jen.PointerTo().ID("exampleUser1"),
						jen.PointerTo().ID("exampleUser2"),
						jen.PointerTo().ID("exampleUser3"),
					),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeUserCreationInput(proj *models.Project) []jen.Code {
	funcName := "BuildFakeUserCreationInput"
	typeName := "UserCreationInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.ID("exampleUser").Assign().ID("BuildFakeUser").Call(),
			jen.Return(
				jen.VarPointer().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("Username").MapAssign().ID("exampleUser").Dot("Username"),
					jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeUserCreationInputFromUser(proj *models.Project) []jen.Code {
	funcName := "BuildFakeUserCreationInputFromUser"
	typeName := "UserCreationInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.VarPointer().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeUserDatabaseCreationInputFromUser(proj *models.Project) []jen.Code {
	funcName := "BuildFakeUserDatabaseCreationInputFromUser"
	typeName := "UserDatabaseCreationInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User"),
		).Params(
			jen.Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("HashedPassword").MapAssign().ID("user").Dot("HashedPassword"),
					jen.ID("TwoFactorSecret").MapAssign().ID("user").Dot("TwoFactorSecret"),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeUserLoginInputFromUser(proj *models.Project) []jen.Code {
	funcName := "BuildFakeUserLoginInputFromUser"
	typeName := "UserLoginInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.VarPointer().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
					jen.ID("TOTPToken").MapAssign().Qual("fmt", "Sprintf").Call(jen.Lit(`0%s`), jen.Qual(utils.FakeLibrary, "Zip").Call()),
				),
			),
		),
	}

	return lines
}

func buildBuildFakePasswordUpdateInput(proj *models.Project) []jen.Code {
	funcName := "BuildFakePasswordUpdateInput"
	typeName := "PasswordUpdateInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.VarPointer().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("NewPassword").MapAssign().Add(utils.FakePasswordFunc()),
					jen.ID("CurrentPassword").MapAssign().Add(utils.FakePasswordFunc()),
					jen.ID("TOTPToken").MapAssign().Qual("fmt", "Sprintf").Call(jen.Lit(`0%s`), jen.Qual(utils.FakeLibrary, "Zip").Call()),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeTOTPSecretRefreshInput(proj *models.Project) []jen.Code {
	funcName := "BuildFakeTOTPSecretRefreshInput"
	typeName := "TOTPSecretRefreshInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), typeName),
		).Block(
			jen.Return(
				jen.VarPointer().Qual(proj.ModelsV1Package(), typeName).Valuesln(
					jen.ID("CurrentPassword").MapAssign().Add(utils.FakePasswordFunc()),
					jen.ID("TOTPToken").MapAssign().Qual("fmt", "Sprintf").Call(jen.Lit(`0%s`), jen.Qual(utils.FakeLibrary, "Zip").Call()),
				),
			),
		),
	}

	return lines
}
