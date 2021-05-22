package apiclients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageCommentf(`Package %s provides a series of HTTP handlers for managing
OAuth2 clients in a compatible database.`, packageName)

	return code
}
