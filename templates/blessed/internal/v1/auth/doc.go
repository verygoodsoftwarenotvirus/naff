package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("auth")

	ret.PackageComment(`Package auth provides functions and structures to facilitate salting and hashing passwords, as well as
validating TOTP tokens`)

	return ret
}
