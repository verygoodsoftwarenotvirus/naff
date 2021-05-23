package fake

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)
	utils.AddImports(proj, code, false)

	code.Add(buildBuildFakeUser(proj)...)
	code.Add(buildBuildDatabaseCreationResponse(proj)...)
	code.Add(buildBuildFakeUserList(proj)...)
	code.Add(buildBuildFakeUserCreationInput(proj)...)
	code.Add(buildBuildFakeUserCreationInputFromUser(proj)...)
	code.Add(buildBuildFakeUserDatabaseCreationInputFromUser(proj)...)
	code.Add(buildBuildFakeUserLoginInputFromUser(proj)...)
	code.Add(buildBuildFakePasswordUpdateInput(proj)...)
	code.Add(buildBuildFakeTOTPSecretRefreshInput(proj)...)
	code.Add(buildBuildFakeTOTPSecretValidationInputForUser(proj)...)

	return code
}

func buildBuildFakeUser(proj *models.Project) []jen.Code {
	funcName := "BuildFakeUser"
	typeName := "User"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.Return(
				jen.AddressOf().Qual(proj.TypesPackage(), typeName).Valuesln(
					jen.ID("ID").MapAssign().Uint64().Call(utils.FakeUint32Func()),
					jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
					jen.Comment(`HashedPassword: ""`),
					jen.Comment("Salt:           []byte(fake.Word())"),
					jen.ID("TwoFactorSecret").MapAssign().Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(
						jen.Index().Byte().Call(jen.Qual(constants.FakeLibrary, "Password").Call(
							jen.False(),
							jen.True(),
							jen.True(),
							jen.False(),
							jen.False(),
							jen.Lit(32),
						)),
					),
					jen.ID("TwoFactorSecretVerifiedOn").MapAssign().Func().Params(jen.ID("i").Uint64()).PointerTo().Uint64().SingleLineBlock(
						jen.Return(jen.AddressOf().ID("i")),
					).Call(
						jen.Uint64().Call(
							jen.Uint32().Call(
								jen.Qual(constants.FakeLibrary, "Date").Call().Dot("Unix").Call(),
							),
						),
					),
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
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("user").PointerTo().Qual(proj.TypesPackage(), "User"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.Return(
				jen.AddressOf().Qual(proj.TypesPackage(), typeName).Valuesln(
					jen.ID("ID").MapAssign().ID("user").Dot("ID"),
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("TwoFactorSecret").MapAssign().ID("user").Dot("TwoFactorSecret"),
					jen.ID("PasswordLastChangedOn").MapAssign().ID("user").Dot("PasswordLastChangedOn"),
					jen.ID("IsAdmin").MapAssign().ID("user").Dot("IsAdmin"),
					jen.ID("CreatedOn").MapAssign().ID("user").Dot("CreatedOn"),
					jen.ID("LastUpdatedOn").MapAssign().ID("user").Dot("LastUpdatedOn"),
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
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.ID(utils.BuildFakeVarName("User1")).Assign().ID("BuildFakeUser").Call(),
			jen.ID(utils.BuildFakeVarName("User2")).Assign().ID("BuildFakeUser").Call(),
			jen.ID(utils.BuildFakeVarName("User3")).Assign().ID("BuildFakeUser").Call(),
			jen.Line(),
			jen.Return(
				jen.AddressOf().Qual(proj.TypesPackage(), typeName).Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.TypesPackage(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().One(),
						jen.ID("Limit").MapAssign().Lit(20),
					),
					jen.ID("Users").MapAssign().Index().Qual(proj.TypesPackage(), "User").Valuesln(
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
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.ID(utils.BuildFakeVarName("User")).Assign().ID("BuildFakeUser").Call(),
			jen.Return(
				jen.AddressOf().Qual(proj.TypesPackage(), typeName).Valuesln(
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
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("user").PointerTo().Qual(proj.TypesPackage(), "User"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.Return(
				jen.AddressOf().Qual(proj.TypesPackage(), typeName).Valuesln(
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
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("user").PointerTo().Qual(proj.TypesPackage(), "User"),
		).Params(
			jen.Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.Return(
				jen.Qual(proj.TypesPackage(), typeName).Valuesln(
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
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("user").PointerTo().Qual(proj.TypesPackage(), "User"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.Return(
				jen.AddressOf().Qual(proj.TypesPackage(), typeName).Valuesln(
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
					jen.ID("TOTPToken").MapAssign().Qual("fmt", "Sprintf").Call(jen.Lit(`0%s`), jen.Qual(constants.FakeLibrary, "Zip").Call()),
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
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.Return(
				jen.AddressOf().Qual(proj.TypesPackage(), typeName).Valuesln(
					jen.ID("NewPassword").MapAssign().Add(utils.FakePasswordFunc()),
					jen.ID("CurrentPassword").MapAssign().Add(utils.FakePasswordFunc()),
					jen.ID("TOTPToken").MapAssign().Qual("fmt", "Sprintf").Call(jen.Lit(`0%s`), jen.Qual(constants.FakeLibrary, "Zip").Call()),
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
		jen.Commentf("%s builds a faked %s.", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.Return(
				jen.AddressOf().Qual(proj.TypesPackage(), typeName).Valuesln(
					jen.ID("CurrentPassword").MapAssign().Add(utils.FakePasswordFunc()),
					jen.ID("TOTPToken").MapAssign().Qual("fmt", "Sprintf").Call(jen.Lit(`0%s`), jen.Qual(constants.FakeLibrary, "Zip").Call()),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeTOTPSecretValidationInputForUser(proj *models.Project) []jen.Code {
	funcName := "BuildFakeTOTPSecretValidationInputForUser"
	typeName := "TOTPSecretVerificationInput"

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s for a given user", funcName, typeName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID("user").PointerTo().Qual(proj.TypesPackage(), "User"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), typeName),
		).Body(
			jen.List(jen.ID("token"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.Err().DoesNotEqual().Nil()).Body(
				jen.Qual("log", "Panicf").Call(
					jen.Lit("error generating TOTP token for fake user: %v"),
					jen.Err(),
				),
			),
			jen.Line(),
			jen.Return(
				jen.AddressOf().Qual(proj.TypesPackage(), typeName).Valuesln(
					jen.ID(constants.UserIDFieldName).MapAssign().ID("user").Dot("ID"),
					jen.ID("TOTPToken").MapAssign().ID("token"),
				),
			),
		),
	}

	return lines
}
