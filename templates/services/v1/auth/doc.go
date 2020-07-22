package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageComment(`Package auth implements a user authentication layer for a web server, issuing
cookies, validating requests via middleware`)

	return code
}
