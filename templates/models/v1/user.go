package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("models")

	utils.AddImports(proj, code)

	code.Add(
		jen.Type().Defs(
			jen.Comment("User represents a user."),
			jen.ID("User").Struct(
				jen.ID("Salt").Index().Byte().Tag(jsonTag("-")),
				jen.ID("Username").String().Tag(jsonTag("username")),
				jen.ID("HashedPassword").String().Tag(jsonTag("-")),
				jen.ID("TwoFactorSecret").String().Tag(jsonTag("-")),
				jen.ID("ID").Uint64().Tag(jsonTag("id")),
				jen.ID("PasswordLastChangedOn").PointerTo().Uint64().Tag(jsonTag("passwordLastChangedOn")),
				jen.ID("TwoFactorSecretVerifiedOn").PointerTo().Uint64().Tag(jsonTag("-")),
				jen.ID("CreatedOn").Uint64().Tag(jsonTag("createdOn")),
				jen.ID("LastUpdatedOn").PointerTo().Uint64().Tag(jsonTag("lastUpdatedOn")),
				jen.ID("ArchivedOn").PointerTo().Uint64().Tag(jsonTag("archivedOn")),
				jen.ID("IsAdmin").Bool().Tag(jsonTag("isAdmin")),
				jen.ID("RequiresPasswordChange").Bool().Tag(jsonTag("requiresPasswordChange")),
			),
			jen.Line(),
			jen.Comment("UserList represents a list of users."),
			jen.ID("UserList").Struct(
				jen.ID("Pagination"),
				jen.ID("Users").Index().ID("User").Tag(jsonTag("users")),
			),
			jen.Line(),
			jen.Comment("UserLoginInput represents the payload used to log in a user."),
			jen.ID("UserLoginInput").Struct(
				jen.ID("Username").String().Tag(jsonTag("username")),
				jen.ID("Password").String().Tag(jsonTag("password")),
				jen.ID("TOTPToken").String().Tag(jsonTag("totpToken")),
			),
			jen.Line(),
			jen.Comment("UserCreationInput represents the input required from users to register an account."),
			jen.ID("UserCreationInput").Struct(
				jen.ID("Username").String().Tag(jsonTag("username")),
				jen.ID("Password").String().Tag(jsonTag("password")),
			),
			jen.Line(),
			jen.Comment("UserDatabaseCreationInput is used by the user creation route to communicate with the database."),
			jen.ID("UserDatabaseCreationInput").Struct(
				jen.ID("Salt").Index().Byte().Tag(jsonTag("-")),
				jen.ID("Username").String().Tag(jsonTag("-")),
				jen.ID("HashedPassword").String().Tag(jsonTag("-")),
				jen.ID("TwoFactorSecret").String().Tag(jsonTag("-")),
			),
			jen.Line(),
			jen.Comment("UserCreationResponse is a response structure for Users that doesn't contain password fields, but does contain the two factor secret."),
			jen.ID("UserCreationResponse").Struct(
				jen.ID("ID").Uint64().Tag(jsonTag("id")),
				jen.ID("Username").String().Tag(jsonTag("username")),
				jen.ID("TwoFactorSecret").String().Tag(jsonTag("twoFactorSecret")),
				jen.ID("PasswordLastChangedOn").PointerTo().Uint64().Tag(jsonTag("passwordLastChangedOn")),
				jen.ID("IsAdmin").Bool().Tag(jsonTag("isAdmin")),
				jen.ID("CreatedOn").Uint64().Tag(jsonTag("createdOn")),
				jen.ID("LastUpdatedOn").PointerTo().Uint64().Tag(jsonTag("lastUpdatedOn")),
				jen.ID("ArchivedOn").PointerTo().Uint64().Tag(jsonTag("archivedOn")),
				jen.ID("TwoFactorQRCode").String().Tag(jsonTag("qrCode")),
			),
			jen.Line(),
			jen.Comment("PasswordUpdateInput represents input a user would provide when updating their password."),
			jen.ID("PasswordUpdateInput").Struct(
				jen.ID("NewPassword").String().Tag(jsonTag("newPassword")),
				jen.ID("CurrentPassword").String().Tag(jsonTag("currentPassword")),
				jen.ID("TOTPToken").String().Tag(jsonTag("totpToken")),
			),
			jen.Line(),
			jen.Comment("TOTPSecretRefreshInput represents input a user would provide when updating their 2FA secret."),
			jen.ID("TOTPSecretRefreshInput").Struct(
				jen.ID("CurrentPassword").String().Tag(jsonTag("currentPassword")),
				jen.ID("TOTPToken").String().Tag(jsonTag("totpToken")),
			),
			jen.Line(),
			jen.Comment("TOTPSecretVerificationInput represents input a user would provide when validating their 2FA secret."),
			jen.ID("TOTPSecretVerificationInput").Struct(
				jen.ID("UserID").Uint64().Tag(jsonTag("userID")),
				jen.ID("TOTPToken").String().Tag(jsonTag("totpToken")),
			),
			jen.Line(),
			jen.Comment("TOTPSecretRefreshResponse represents the response we provide to a user when updating their 2FA secret."),
			jen.ID("TOTPSecretRefreshResponse").Struct(
				jen.ID("TwoFactorSecret").String().Tag(jsonTag("twoFactorSecret")),
			),
			jen.Line(),
			jen.Comment("UserDataManager describes a structure which can manage users in permanent storage."),
			jen.ID("UserDataManager").Interface(
				jen.ID("GetUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().ID("User"), jen.Error()),
				jen.ID("GetUserWithUnverifiedTwoFactorSecret").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().ID("User"), jen.Error()),
				jen.ID("VerifyUserTwoFactorSecret").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.Error()),
				jen.ID("GetUserByUsername").Params(constants.CtxParam(), jen.ID("username").String()).Params(jen.PointerTo().ID("User"), jen.Error()),
				jen.ID("GetAllUsersCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()),
				jen.ID("GetUsers").Params(constants.CtxParam(), utils.QueryFilterParam(nil)).Params(jen.PointerTo().ID("UserList"), jen.Error()),
				jen.ID("CreateUser").Params(constants.CtxParam(), jen.ID("input").ID("UserDatabaseCreationInput")).Params(jen.PointerTo().ID("User"), jen.Error()),
				jen.ID("UpdateUser").Params(constants.CtxParam(), jen.ID("updated").PointerTo().ID("User")).Params(jen.Error()),
				jen.ID("UpdateUserPassword").Params(constants.CtxParam(), constants.UserIDParam(), jen.ID("newHash").String()).Params(jen.Error()),
				jen.ID("ArchiveUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("UserDataServer describes a structure capable of serving traffic related to users."),
			jen.ID("UserDataServer").Interface(
				jen.ID("UserInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("PasswordUpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("TOTPSecretRefreshInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("TOTPSecretVerificationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.Line(),
				jen.ID("ListHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("CreateHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("ReadHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("NewTOTPSecretHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("TOTPSecretVerificationHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("UpdatePasswordHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("ArchiveHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Update accepts a User as input and merges those values if they're set."),
		jen.Line(),
		jen.Func().Params(jen.ID("u").PointerTo().ID("User")).ID("Update").Params(jen.ID("input").PointerTo().ID("User")).Block(
			jen.If(jen.ID("input").Dot("Username").DoesNotEqual().EmptyString().And().ID("input").Dot("Username").DoesNotEqual().ID("u").Dot("Username")).Block(
				jen.ID("u").Dot("Username").Equals().ID("input").Dot("Username"),
			),
			jen.Line(),
			jen.If(jen.ID("input").Dot("HashedPassword").DoesNotEqual().EmptyString().And().ID("input").Dot("HashedPassword").DoesNotEqual().ID("u").Dot("HashedPassword")).Block(
				jen.ID("u").Dot("HashedPassword").Equals().ID("input").Dot("HashedPassword"),
			),
			jen.Line(),
			jen.If(jen.ID("input").Dot("TwoFactorSecret").DoesNotEqual().EmptyString().And().ID("input").Dot("TwoFactorSecret").DoesNotEqual().ID("u").Dot("TwoFactorSecret")).Block(
				jen.ID("u").Dot("TwoFactorSecret").Equals().ID("input").Dot("TwoFactorSecret"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ToSessionInfo accepts a User as input and merges those values if they're set."),
		jen.Line(),
		jen.Func().Params(jen.ID("u").PointerTo().ID("User")).ID("ToSessionInfo").Params().PointerTo().ID("SessionInfo").Block(
			jen.Return(jen.AddressOf().ID("SessionInfo").Valuesln(
				jen.ID(constants.UserIDFieldName).MapAssign().ID("u").Dot("ID"),
				jen.ID("UserIsAdmin").MapAssign().ID("u").Dot("IsAdmin"),
			)),
		),
		jen.Line(),
	)

	return code
}
