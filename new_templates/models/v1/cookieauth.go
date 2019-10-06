package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func cookieauthDotGo() *jen.File {
	ret := jen.NewFile("models")
	ret.Add(jen.Null().Type().ID("CookieAuth").Struct(
		jen.ID("UserID").ID("uint64"),
		jen.ID("Admin").ID("bool"),
		jen.ID("Username").ID("string"),
	),
	)
	return ret
}
