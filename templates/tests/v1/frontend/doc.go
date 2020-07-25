package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("frontend")

	code.PackageComment(`Package frontend is a series of selenium tests which validate certain aspects of our
frontend, to guard against failed contributions to the frontend code.`)

	return code
}
