package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func userDotGo() *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("UserKey").ID("ContextKey").Op("=").Lit("user").Var().ID("UserIDKey").ID("ContextKey").Op("=").Lit("user_id").Var().ID("UserIsAdminKey").ID("ContextKey").Op("=").Lit("is_admin"),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("User").Struct(jen.ID("ID").ID("uint64"), jen.ID("Username").ID("string"), jen.ID("HashedPassword").ID("string"), jen.ID("Salt").Index().ID("byte"), jen.ID("TwoFactorSecret").ID("string"), jen.ID("PasswordLastChangedOn").Op("*").ID("uint64"), jen.ID("IsAdmin").ID("bool"), jen.ID("CreatedOn").ID("uint64"), jen.ID("UpdatedOn").Op("*").ID("uint64"), jen.ID("ArchivedOn").Op("*").ID("uint64")).Type().ID("UserList").Struct(jen.ID("Pagination"), jen.ID("Users").Index().ID("User")).Type().ID("UserLoginInput").Struct(jen.ID("Username").ID("string"), jen.ID("Password").ID("string"), jen.ID("TOTPToken").ID("string")).Type().ID("UserInput").Struct(jen.ID("Username").ID("string"), jen.ID("Password").ID("string"), jen.ID("TwoFactorSecret").ID("string")).Type().ID("UserCreationResponse").Struct(jen.ID("ID").ID("uint64"), jen.ID("Username").ID("string"), jen.ID("TwoFactorSecret").ID("string"), jen.ID("PasswordLastChangedOn").Op("*").ID("uint64"), jen.ID("IsAdmin").ID("bool"), jen.ID("CreatedOn").ID("uint64"), jen.ID("UpdatedOn").Op("*").ID("uint64"), jen.ID("ArchivedOn").Op("*").ID("uint64"), jen.ID("TwoFactorQRCode").ID("string")).Type().ID("PasswordUpdateInput").Struct(jen.ID("NewPassword").ID("string"), jen.ID("CurrentPassword").ID("string"), jen.ID("TOTPToken").ID("string")).Type().ID("TOTPSecretRefreshInput").Struct(jen.ID("CurrentPassword").ID("string"), jen.ID("TOTPToken").ID("string")).Type().ID("TOTPSecretRefreshResponse").Struct(jen.ID("TwoFactorSecret").ID("string")).Type().ID("UserDataManager").Interface(jen.ID("GetUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("User"), jen.ID("error")), jen.ID("GetUserByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.Op("*").ID("User"), jen.ID("error")), jen.ID("GetUserCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.ID("uint64"), jen.ID("error")), jen.ID("GetUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("UserList"), jen.ID("error")), jen.ID("CreateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("UserInput")).Params(jen.Op("*").ID("User"), jen.ID("error")), jen.ID("UpdateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("User")).Params(jen.ID("error")), jen.ID("ArchiveUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error"))).Type().ID("UserDataServer").Interface(jen.ID("UserInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")), jen.ID("PasswordUpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")), jen.ID("TOTPSecretRefreshInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")), jen.ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("NewTOTPSecretHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("UpdatePasswordHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc"))),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Update accepts a User as input and merges those values if they're set"),
		jen.Line(),
		jen.Func().Params(jen.ID("u").Op("*").ID("User")).ID("Update").Params(jen.ID("input").Op("*").ID("User")).Block(
			jen.If(jen.ID("input").Dot(
				"Username",
			).Op("!=").Lit("").Op("&&").ID("input").Dot(
				"Username",
			).Op("!=").ID("u").Dot(
				"Username",
			)).Block(
				jen.ID("u").Dot(
					"Username",
				).Op("=").ID("input").Dot(
					"Username",
				),
			),
			jen.If(jen.ID("input").Dot(
				"HashedPassword",
			).Op("!=").Lit("").Op("&&").ID("input").Dot(
				"HashedPassword",
			).Op("!=").ID("u").Dot(
				"HashedPassword",
			)).Block(
				jen.ID("u").Dot(
					"HashedPassword",
				).Op("=").ID("input").Dot(
					"HashedPassword",
				),
			),
			jen.If(jen.ID("input").Dot(
				"TwoFactorSecret",
			).Op("!=").Lit("").Op("&&").ID("input").Dot(
				"TwoFactorSecret",
			).Op("!=").ID("u").Dot(
				"TwoFactorSecret",
			)).Block(
				jen.ID("u").Dot(
					"TwoFactorSecret",
				).Op("=").ID("input").Dot(
					"TwoFactorSecret",
				),
			),
		),
		jen.Line(),
	)
	return ret
}
