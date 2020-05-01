package main

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
)

func main() {
	s := wordsmith.FromSingularPascalCase("Postgres")

	println(s.LowercaseAbbreviation())
}
