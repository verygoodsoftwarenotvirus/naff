package testutil

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func docDotGo() *jen.File {
	ret := jen.NewFile("testutil")
	return ret
}
