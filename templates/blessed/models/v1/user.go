package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("UserKey is the non-string type we use for referencing a user in a context"),
			jen.ID("UserKey").ID("ContextKey").Op("=").Lit("user"),
			jen.Comment("UserIDKey is the non-string type we use for referencing a user ID in a context"),
			jen.ID("UserIDKey").ID("ContextKey").Op("=").Lit("user_id"),
			jen.Comment("UserIsAdminKey is the non-string type we use for referencing a user's admin status in a context"),
			jen.ID("UserIsAdminKey").ID("ContextKey").Op("=").Lit("is_admin"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("User represents a user"),
			jen.ID("User").Struct(
				jen.ID("ID").ID("uint64").Tag(jsonTag("id")),
				jen.ID("Username").ID("string").Tag(jsonTag("username")),
				jen.ID("HashedPassword").ID("string").Tag(jsonTag("-")),
				jen.ID("Salt").Index().ID("byte").Tag(jsonTag("-")),
				jen.ID("TwoFactorSecret").ID("string").Tag(jsonTag("-")),
				jen.ID("PasswordLastChangedOn").Op("*").ID("uint64").Tag(jsonTag("password_last_changed_on")),
				jen.ID("IsAdmin").ID("bool").Tag(jsonTag("is_admin")),
				jen.ID("CreatedOn").ID("uint64").Tag(jsonTag("created_on")),
				jen.ID("UpdatedOn").Op("*").ID("uint64").Tag(jsonTag("updated_on")),
				jen.ID("ArchivedOn").Op("*").ID("uint64").Tag(jsonTag("archived_on")),
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
				jen.ID("Username").ID("string").Tag(jsonTag("username")),
				jen.ID("Password").ID("string").Tag(jsonTag("password")),
				jen.ID("TOTPToken").ID("string").Tag(jsonTag("totp_token")),
			),
			jen.Line(),
			jen.Comment("UserInput represents the input required to modify/create users"),
			jen.ID("UserInput").Struct(
				jen.ID("Username").ID("string").Tag(jsonTag("username")),
				jen.ID("Password").ID("string").Tag(jsonTag("password")),
				jen.ID("TwoFactorSecret").ID("string").Tag(jsonTag("-")),
			),
			jen.Line(),
			jen.Comment("UserCreationResponse is a response structure for Users that doesn't contain password fields, but does contain the two factor secret"),
			jen.ID("UserCreationResponse").Struct(
				jen.ID("ID").ID("uint64").Tag(jsonTag("id")),
				jen.ID("Username").ID("string").Tag(jsonTag("username")),
				jen.ID("TwoFactorSecret").ID("string").Tag(jsonTag("two_factor_secret")),
				jen.ID("PasswordLastChangedOn").Op("*").ID("uint64").Tag(jsonTag("password_last_changed_on")),
				jen.ID("IsAdmin").ID("bool").Tag(jsonTag("is_admin")),
				jen.ID("CreatedOn").ID("uint64").Tag(jsonTag("created_on")),
				jen.ID("UpdatedOn").Op("*").ID("uint64").Tag(jsonTag("updated_on")),
				jen.ID("ArchivedOn").Op("*").ID("uint64").Tag(jsonTag("archived_on")),
				jen.ID("TwoFactorQRCode").ID("string").Tag(jsonTag("qr_code")),
			),
			jen.Line(),
			jen.Comment("PasswordUpdateInput represents input a user would provide when updating their password"),
			jen.ID("PasswordUpdateInput").Struct(
				jen.ID("NewPassword").ID("string").Tag(jsonTag("new_password")),
				jen.ID("CurrentPassword").ID("string").Tag(jsonTag("current_password")),
				jen.ID("TOTPToken").ID("string").Tag(jsonTag("totp_token")),
			),
			jen.Line(),
			jen.Comment("TOTPSecretRefreshInput represents input a user would provide when updating their 2FA secret"),
			jen.ID("TOTPSecretRefreshInput").Struct(
				jen.ID("CurrentPassword").ID("string").Tag(jsonTag("current_password")),
				jen.ID("TOTPToken").ID("string").Tag(jsonTag("totp_token")),
			),
			jen.Line(),
			jen.Comment("TOTPSecretRefreshResponse represents the response we provide to a user when updating their 2FA secret"),
			jen.ID("TOTPSecretRefreshResponse").Struct(
				jen.ID("TwoFactorSecret").ID("string").Tag(jsonTag("two_factor_secret")),
			),
			jen.Line(),
			jen.Comment("UserDataManager describes a structure which can manage users in permanent storage"),
			jen.ID("UserDataManager").Interface(
				jen.ID("GetUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("User"), jen.ID("error")),
				jen.ID("GetUserByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.Op("*").ID("User"), jen.ID("error")),
				jen.ID("GetUserCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("UserList"), jen.ID("error")),
				jen.ID("CreateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("UserInput")).Params(jen.Op("*").ID("User"), jen.ID("error")),
				jen.ID("UpdateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("User")).Params(jen.ID("error")),
				jen.ID("ArchiveUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error")),
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
		jen.Func().Params(jen.ID("u").Op("*").ID("User")).ID("Update").Params(jen.ID("input").Op("*").ID("User")).Block(
			jen.If(jen.ID("input").Dot("Username").Op("!=").Lit("").Op("&&").ID("input").Dot("Username").Op("!=").ID("u").Dot("Username")).Block(
				jen.ID("u").Dot("Username").Op("=").ID("input").Dot("Username"),
			),
			jen.Line(),
			jen.If(jen.ID("input").Dot("HashedPassword").Op("!=").Lit("").Op("&&").ID("input").Dot("HashedPassword").Op("!=").ID("u").Dot("HashedPassword")).Block(
				jen.ID("u").Dot("HashedPassword").Op("=").ID("input").Dot("HashedPassword"),
			),
			jen.Line(),
			jen.If(jen.ID("input").Dot("TwoFactorSecret").Op("!=").Lit("").Op("&&").ID("input").Dot("TwoFactorSecret").Op("!=").ID("u").Dot("TwoFactorSecret")).Block(
				jen.ID("u").Dot("TwoFactorSecret").Op("=").ID("input").Dot("TwoFactorSecret"),
			),
		),
		jen.Line(),
	)
	return ret
}
