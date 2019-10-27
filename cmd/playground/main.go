package main

import (
	"bytes"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"log"
)

func main() {

	ret := jen.NewFile("postgres")

	ret.Add(
		jen.Var().ID("x").Op("=").Qual("strings", "NewReplacer").PairedCallln(
			jen.Lit("$"), jen.RawString(`\$`),
			jen.Lit("("), jen.Lit("\\("),
		),
		jen.Line(),
	)

	var b bytes.Buffer
	if err := ret.Render(&b); err != nil {
		log.Fatal(err)
	}
	x := b.String()
	log.Println(x)
}
