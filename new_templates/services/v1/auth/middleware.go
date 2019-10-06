package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func middlewareDotGo() *jen.File {
	ret := jen.NewFile("auth")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("UserLoginInputMiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("user_login_input").Var().ID("UsernameFormKey").Op("=").Lit("username").Var().ID("PasswordFormKey").Op("=").Lit("password").Var().ID("TOTPTokenFormKey").Op("=").Lit("totp_token"),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
