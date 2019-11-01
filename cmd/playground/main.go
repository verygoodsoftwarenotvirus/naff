package main

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
)

func main() {
	y := wordsmith.FromSingularPascalCase("MariaDB")
	fmt.Println(y.SingularPackageName())
}
