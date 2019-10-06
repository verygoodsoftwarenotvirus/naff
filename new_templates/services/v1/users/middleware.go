package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareDotGo() *jen.File {
	ret := jen.NewFile("users")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("UserCreationMiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("user_creation_input").Var().ID("PasswordChangeMiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("user_password_change").Var().ID("TOTPSecretRefreshMiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("totp_refresh"),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
