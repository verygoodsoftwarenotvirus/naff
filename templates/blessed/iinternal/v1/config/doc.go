package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("config")

	ret.PackageComment("Package config provides configuration structs for every service\n")

	return ret
}
