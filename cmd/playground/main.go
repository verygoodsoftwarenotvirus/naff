package main

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
)

func main() {
	s := wordsmith.FromSingularPascalCase("BlahBlahBlah")

	println(s.RouteName())
	println(s.PackageName())
	println(s.PluralRouteName())
}
