package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("auth")

	code.PackageComment(`Package auth provides functions and structures to facilitate salting and hashing passwords, as well as
validating TOTP tokens`)

	return code
}
