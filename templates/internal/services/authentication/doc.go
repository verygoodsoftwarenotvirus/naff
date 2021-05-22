package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageCommentf(`Package %s implements a user authentication layer for a web server, issuing
cookies, validating requests via middleware`, packageName)

	return code
}
