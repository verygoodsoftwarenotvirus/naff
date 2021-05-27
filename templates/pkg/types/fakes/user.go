package fakes

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeUser builds a faked User."),
		jen.Line(),
		jen.Func().ID("BuildFakeUser").Params().Params(jen.Op("*").ID("types").Dot("User")).Body(
			jen.Return().Op("&").ID("types").Dot("User").Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Uint32").Call()), jen.ID("ExternalID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("Username").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("false"),
				jen.ID("false"),
				jen.Lit(32),
			), jen.ID("ServiceAccountStatus").Op(":").ID("types").Dot("GoodStandingAccountStatus"), jen.ID("TwoFactorSecret").Op(":").Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.Index().ID("byte").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("false"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("false"),
				jen.ID("false"),
				jen.Lit(32),
			))), jen.ID("TwoFactorSecretVerifiedOn").Op(":").Func().Params(jen.ID("i").ID("uint64")).Params(jen.Op("*").ID("uint64")).Body(
				jen.Return().Op("&").ID("i")).Call(jen.ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call()))), jen.ID("ServiceRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("ServiceUserRole").Dot("String").Call()), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserCreationResponseFromUser builds a faked UserCreationResponse."),
		jen.Line(),
		jen.Func().ID("BuildUserCreationResponseFromUser").Params(jen.ID("user").Op("*").ID("types").Dot("User")).Params(jen.Op("*").ID("types").Dot("UserCreationResponse")).Body(
			jen.Return().Op("&").ID("types").Dot("UserCreationResponse").Valuesln(jen.ID("CreatedUserID").Op(":").ID("user").Dot("ID"), jen.ID("Username").Op(":").ID("user").Dot("Username"), jen.ID("CreatedOn").Op(":").ID("user").Dot("CreatedOn"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeUserList builds a faked UserList."),
		jen.Line(),
		jen.Func().ID("BuildFakeUserList").Params().Params(jen.Op("*").ID("types").Dot("UserList")).Body(
			jen.Var().Defs(
				jen.ID("examples").Index().Op("*").ID("types").Dot("User"),
			),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("exampleQuantity"), jen.ID("i").Op("++")).Body(
				jen.ID("examples").Op("=").ID("append").Call(
					jen.ID("examples"),
					jen.ID("BuildFakeUser").Call(),
				)),
			jen.Return().Op("&").ID("types").Dot("UserList").Valuesln(jen.ID("Pagination").Op(":").ID("types").Dot("Pagination").Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("FilteredCount").Op(":").ID("exampleQuantity").Op("/").Lit(2), jen.ID("TotalCount").Op(":").ID("exampleQuantity")), jen.ID("Users").Op(":").ID("examples")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeUserCreationInput builds a faked UserRegistrationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeUserCreationInput").Params().Params(jen.Op("*").ID("types").Dot("UserRegistrationInput")).Body(
			jen.ID("exampleUser").Op(":=").ID("BuildFakeUser").Call(),
			jen.Return().Op("&").ID("types").Dot("UserRegistrationInput").Valuesln(jen.ID("Username").Op(":").ID("exampleUser").Dot("Username"), jen.ID("Password").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildTestUserCreationConfig builds a faked TestUserCreationConfig."),
		jen.Line(),
		jen.Func().ID("BuildTestUserCreationConfig").Params().Params(jen.Op("*").ID("types").Dot("TestUserCreationConfig")).Body(
			jen.ID("exampleUser").Op(":=").ID("BuildFakeUserCreationInput").Call(),
			jen.Return().Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(jen.ID("Username").Op(":").ID("exampleUser").Dot("Username"), jen.ID("Password").Op(":").ID("exampleUser").Dot("Password"), jen.ID("HashedPassword").Op(":").Lit("hashed passwords"), jen.ID("IsServiceAdmin").Op(":").ID("false")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeUserRegistrationInputFromUser builds a faked UserRegistrationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeUserRegistrationInputFromUser").Params(jen.ID("user").Op("*").ID("types").Dot("User")).Params(jen.Op("*").ID("types").Dot("UserRegistrationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("UserRegistrationInput").Valuesln(jen.ID("Username").Op(":").ID("user").Dot("Username"), jen.ID("Password").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeUserDataStoreCreationInputFromUser builds a faked UserDataStoreCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeUserDataStoreCreationInputFromUser").Params(jen.ID("user").Op("*").ID("types").Dot("User")).Params(jen.Op("*").ID("types").Dot("UserDataStoreCreationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("UserDataStoreCreationInput").Valuesln(jen.ID("Username").Op(":").ID("user").Dot("Username"), jen.ID("HashedPassword").Op(":").ID("user").Dot("HashedPassword"), jen.ID("TwoFactorSecret").Op(":").ID("user").Dot("TwoFactorSecret"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeUserReputationUpdateInputFromUser builds a faked UserReputationUpdateInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeUserReputationUpdateInputFromUser").Params(jen.ID("user").Op("*").ID("types").Dot("User")).Params(jen.Op("*").ID("types").Dot("UserReputationUpdateInput")).Body(
			jen.Return().Op("&").ID("types").Dot("UserReputationUpdateInput").Valuesln(jen.ID("TargetUserID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("NewReputation").Op(":").ID("user").Dot("ServiceAccountStatus"), jen.ID("Reason").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Sentence").Call(jen.Lit(10)))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeUserRegistrationInput builds a faked UserLoginInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeUserRegistrationInput").Params().Params(jen.Op("*").ID("types").Dot("UserRegistrationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("UserRegistrationInput").Valuesln(jen.ID("Username").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Username").Call(), jen.ID("Password").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeUserLoginInputFromUser builds a faked UserLoginInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeUserLoginInputFromUser").Params(jen.ID("user").Op("*").ID("types").Dot("User")).Params(jen.Op("*").ID("types").Dot("UserLoginInput")).Body(
			jen.Return().Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("user").Dot("Username"), jen.ID("Password").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			), jen.ID("TOTPToken").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit("0%s"),
				jen.Qual("github.com/brianvoe/gofakeit/v5", "Zip").Call(),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakePasswordUpdateInput builds a faked PasswordUpdateInput."),
		jen.Line(),
		jen.Func().ID("BuildFakePasswordUpdateInput").Params().Params(jen.Op("*").ID("types").Dot("PasswordUpdateInput")).Body(
			jen.Return().Op("&").ID("types").Dot("PasswordUpdateInput").Valuesln(jen.ID("NewPassword").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			), jen.ID("CurrentPassword").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			), jen.ID("TOTPToken").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit("0%s"),
				jen.Qual("github.com/brianvoe/gofakeit/v5", "Zip").Call(),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeTOTPSecretRefreshInput builds a faked TOTPSecretRefreshInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeTOTPSecretRefreshInput").Params().Params(jen.Op("*").ID("types").Dot("TOTPSecretRefreshInput")).Body(
			jen.Return().Op("&").ID("types").Dot("TOTPSecretRefreshInput").Valuesln(jen.ID("CurrentPassword").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			), jen.ID("TOTPToken").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit("0%s"),
				jen.Qual("github.com/brianvoe/gofakeit/v5", "Zip").Call(),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeTOTPSecretVerificationInput builds a faked TOTPSecretVerificationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeTOTPSecretVerificationInput").Params().Params(jen.Op("*").ID("types").Dot("TOTPSecretVerificationInput")).Body(
			jen.ID("user").Op(":=").ID("BuildFakeUser").Call(),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual("log", "Panicf").Call(
					jen.Lit("error generating TOTP token for fakes user: %v"),
					jen.ID("err"),
				)),
			jen.Return().Op("&").ID("types").Dot("TOTPSecretVerificationInput").Valuesln(jen.ID("UserID").Op(":").ID("user").Dot("ID"), jen.ID("TOTPToken").Op(":").ID("token")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeTOTPSecretVerificationInputForUser builds a faked TOTPSecretVerificationInput for a given user."),
		jen.Line(),
		jen.Func().ID("BuildFakeTOTPSecretVerificationInputForUser").Params(jen.ID("user").Op("*").ID("types").Dot("User")).Params(jen.Op("*").ID("types").Dot("TOTPSecretVerificationInput")).Body(
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual("log", "Panicf").Call(
					jen.Lit("error generating TOTP token for fakes user: %v"),
					jen.ID("err"),
				)),
			jen.Return().Op("&").ID("types").Dot("TOTPSecretVerificationInput").Valuesln(jen.ID("UserID").Op(":").ID("user").Dot("ID"), jen.ID("TOTPToken").Op(":").ID("token")),
		),
		jen.Line(),
	)

	return code
}
