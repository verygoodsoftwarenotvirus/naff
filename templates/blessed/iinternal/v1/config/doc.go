package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("config")

	code.PackageComment("Package config provides configuration structs for every service\n")

	return code
}
