package fakes

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func fakeDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.Qual("github.com/brianvoe/gofakeit/v5", "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("exampleQuantity").Op("=").Lit(3),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeSQLQuery builds a fake SQL query and arg pair."),
		jen.Line(),
		jen.Func().ID("BuildFakeSQLQuery").Params().Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("s").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s %s WHERE things = ? AND stuff = ?"),
				jen.Qual("github.com/brianvoe/gofakeit/v5", "RandomString").Call(jen.Index().ID("string").Valuesln(jen.Lit("SELECT * FROM"), jen.Lit("INSERT INTO"), jen.Lit("UPDATE"))),
				jen.Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
					jen.ID("true"),
					jen.ID("true"),
					jen.ID("true"),
					jen.ID("false"),
					jen.ID("false"),
					jen.Lit(32),
				),
			),
			jen.Return().List(jen.ID("s"), jen.Index().Interface().Valuesln(jen.Lit("things"), jen.Lit("stuff"))),
		),
		jen.Line(),
	)

	return code
}
