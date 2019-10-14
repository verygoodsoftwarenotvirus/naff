package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func cookieauthDotGo() *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)

	ret.Add(jen.Null().Type().ID("CookieAuth").Struct(jen.ID("UserID").ID("uint64"), jen.ID("Admin").ID("bool"), jen.ID("Username").ID("string")),

		jen.Line(),
	)
	return ret
}
