package configgen

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func docDotGo() *jen.File {
	ret := jen.NewFile("main")
	ret.PackageComment(`
Command config_gen generates configuration files in the local repository, configured
via the precise mechanism that parses them to guard against invalid configuration
`)
	return ret
}
