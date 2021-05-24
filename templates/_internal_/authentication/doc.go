package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageCommentf(`Package %s provides functions and structures to facilitate salting and hashing passwords, as well as
validating TOTP tokens`, packageName)

	return code
}
