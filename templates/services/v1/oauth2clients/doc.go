package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageComment(`Package oauth2clients provides a series of HTTP handlers for managing
OAuth2 clients in a compatible database.`)

	return code
}
