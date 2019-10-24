package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("users")

	ret.PackageComment(`Package users provides a series of HTTP handlers for managing
users, passwords, and two factor secrets in a compatible database.`)

	return ret
}
