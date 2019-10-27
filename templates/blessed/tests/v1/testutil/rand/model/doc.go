package model

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("randmodel")

	ret.PackageComment("Package randmodel contains functions for generating randominstances of models for testing or demonstration purposes\n")

	return ret
}
