package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("UserKey is the non-string type we use for referencing a user in a context"),
			jen.ID("UserKey").ID("ContextKey").Equals().Lit("user"),
			jen.Comment("UserIDKey is the non-string type we use for referencing a user ID in a context"),
			jen.ID("UserIDKey").ID("ContextKey").Equals().Lit("user_id"),
			jen.Comment("UserIsAdminKey is the non-string type we use for referencing a user's admin status in a context"),
			jen.ID("UserIsAdminKey").ID("ContextKey").Equals().Lit("is_admin"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("User represents a user"),
			jen.ID("User").Struct(
				jen.ID("ID").Uint64().Tag(jsonTag("id")),
				jen.ID("Username").String().Tag(jsonTag("username")),
				jen.ID("HashedPassword").String().Tag(jsonTag("-")),
				jen.ID("Salt").Index().Byte().Tag(jsonTag("-")),
				jen.ID("TwoFactorSecret").String().Tag(jsonTag("-")),
				jen.ID("PasswordLastChangedOn").PointerTo().Uint64().Tag(jsonTag("password_last_changed_on")),
				jen.ID("IsAdmin").Bool().Tag(jsonTag("is_admin")),
				jen.ID("CreatedOn").Uint64().Tag(jsonTag("created_on")),
				jen.ID("UpdatedOn").PointerTo().Uint64().Tag(jsonTag("updated_on")),
				jen.ID("ArchivedOn").PointerTo().Uint64().Tag(jsonTag("archived_on")),
			),
			jen.Line(),
			jen.Comment("UserList represents a list of users"),
			jen.ID("UserList").Struct(
				jen.ID("Pagination"),
				jen.ID("Users").Index().ID("User").Tag(jsonTag("users")),
			),
			jen.Line(),
			jen.Comment("UserLoginInput represents the payload used to log in a user"),
			jen.ID("UserLoginInput").Struct(
				jen.ID("Username").String().Tag(jsonTag("username")),
				jen.ID("Password").String().Tag(jsonTag("password")),
				jen.ID("TOTPToken").String().Tag(jsonTag("totp_token")),
			),
			jen.Line(),
			jen.Comment("UserCreationInput represents the input required from users to register an account"),
			jen.ID("UserCreationInput").Struct(
				jen.ID("Username").String().Tag(jsonTag("username")),
				jen.ID("Password").String().Tag(jsonTag("password")),
			),
			jen.Line(),
			jen.Comment("UserDatabaseCreationInput is used by the user creation route to communicate with the database"),
			jen.ID("UserDatabaseCreationInput").Struct(
				jen.ID("Username").String(),
				jen.ID("HashedPassword").String(),
				jen.ID("TwoFactorSecret").String(),
			),
			jen.Line(),
			jen.Comment("UserCreationResponse is a response structure for Users that doesn't contain password fields, but does contain the two factor secret"),
			jen.ID("UserCreationResponse").Struct(
				jen.ID("ID").Uint64().Tag(jsonTag("id")),
				jen.ID("Username").String().Tag(jsonTag("username")),
				jen.ID("TwoFactorSecret").String().Tag(jsonTag("two_factor_secret")),
				jen.ID("PasswordLastChangedOn").PointerTo().Uint64().Tag(jsonTag("password_last_changed_on")),
				jen.ID("IsAdmin").Bool().Tag(jsonTag("is_admin")),
				jen.ID("CreatedOn").Uint64().Tag(jsonTag("created_on")),
				jen.ID("UpdatedOn").PointerTo().Uint64().Tag(jsonTag("updated_on")),
				jen.ID("ArchivedOn").PointerTo().Uint64().Tag(jsonTag("archived_on")),
				jen.ID("TwoFactorQRCode").String().Tag(jsonTag("qr_code")),
			),
			jen.Line(),
			jen.Comment("PasswordUpdateInput represents input a user would provide when updating their password"),
			jen.ID("PasswordUpdateInput").Struct(
				jen.ID("NewPassword").String().Tag(jsonTag("new_password")),
				jen.ID("CurrentPassword").String().Tag(jsonTag("current_password")),
				jen.ID("TOTPToken").String().Tag(jsonTag("totp_token")),
			),
			jen.Line(),
			jen.Comment("TOTPSecretRefreshInput represents input a user would provide when updating their 2FA secret"),
			jen.ID("TOTPSecretRefreshInput").Struct(
				jen.ID("CurrentPassword").String().Tag(jsonTag("current_password")),
				jen.ID("TOTPToken").String().Tag(jsonTag("totp_token")),
			),
			jen.Line(),
			jen.Comment("TOTPSecretRefreshResponse represents the response we provide to a user when updating their 2FA secret"),
			jen.ID("TOTPSecretRefreshResponse").Struct(
				jen.ID("TwoFactorSecret").String().Tag(jsonTag("two_factor_secret")),
			),
			jen.Line(),
			jen.Comment("UserDataManager describes a structure which can manage users in permanent storage"),
			jen.ID("UserDataManager").Interface(
				jen.ID("GetUser").Params(utils.CtxParam(), jen.ID("userID").Uint64()).Params(jen.PointerTo().ID("User"), jen.Error()),
				jen.ID("GetUserByUsername").Params(utils.CtxParam(), jen.ID("username").String()).Params(jen.PointerTo().ID("User"), jen.Error()),
				jen.ID("GetAllUserCount").Params(utils.CtxParam()).Params(jen.Uint64(), jen.Error()),
				jen.ID("GetUsers").Params(utils.CtxParam(), utils.QueryFilterParam()).Params(jen.PointerTo().ID("UserList"), jen.Error()),
				jen.ID("CreateUser").Params(utils.CtxParam(), jen.ID("input").ID("UserDatabaseCreationInput")).Params(jen.PointerTo().ID("User"), jen.Error()),
				jen.ID("UpdateUser").Params(utils.CtxParam(), jen.ID("updated").PointerTo().ID("User")).Params(jen.Error()),
				jen.ID("ArchiveUser").Params(utils.CtxParam(), jen.ID("userID").Uint64()).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("UserDataServer describes a structure capable of serving traffic related to users"),
			jen.ID("UserDataServer").Interface(
				jen.ID("UserInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("PasswordUpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("TOTPSecretRefreshInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.Line(),
				jen.ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("NewTOTPSecretHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("UpdatePasswordHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Update accepts a User as input and merges those values if they're set"),
		jen.Line(),
		jen.Func().Params(jen.ID("u").PointerTo().ID("User")).ID("Update").Params(jen.ID("input").PointerTo().ID("User")).Block(
			jen.If(jen.ID("input").Dot("Username").DoesNotEqual().Lit("").And().ID("input").Dot("Username").DoesNotEqual().ID("u").Dot("Username")).Block(
				jen.ID("u").Dot("Username").Equals().ID("input").Dot("Username"),
			),
			jen.Line(),
			jen.If(jen.ID("input").Dot("HashedPassword").DoesNotEqual().Lit("").And().ID("input").Dot("HashedPassword").DoesNotEqual().ID("u").Dot("HashedPassword")).Block(
				jen.ID("u").Dot("HashedPassword").Equals().ID("input").Dot("HashedPassword"),
			),
			jen.Line(),
			jen.If(jen.ID("input").Dot("TwoFactorSecret").DoesNotEqual().Lit("").And().ID("input").Dot("TwoFactorSecret").DoesNotEqual().ID("u").Dot("TwoFactorSecret")).Block(
				jen.ID("u").Dot("TwoFactorSecret").Equals().ID("input").Dot("TwoFactorSecret"),
			),
		),
		jen.Line(),
	)
	return ret
}
