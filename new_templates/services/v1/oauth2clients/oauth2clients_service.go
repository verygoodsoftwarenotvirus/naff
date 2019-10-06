package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientsServiceDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")
	utils.AddImports(ret)
	return ret
}
