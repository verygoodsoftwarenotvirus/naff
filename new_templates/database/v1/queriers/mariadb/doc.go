package mariadb

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func docDotGo() *jen.File {
	ret := jen.NewFile("mariadb")
	return ret
}
