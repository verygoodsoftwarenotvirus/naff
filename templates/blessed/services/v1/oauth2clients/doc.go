package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")

	ret.PackageComment(`Package oauth2clients provides a series of HTTP handlers for managing
OAuth2 clients in a compatible database.`)

	return ret
}
