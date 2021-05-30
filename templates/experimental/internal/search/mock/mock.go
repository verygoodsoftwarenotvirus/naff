package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("search").Dot("IndexManager").Op("=").Parens(jen.Op("*").ID("IndexManager")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("IndexManager").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Index implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("IndexManager")).ID("Index").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("id").ID("uint64"), jen.ID("value").Interface()).Params(jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("id"),
				jen.ID("value"),
			),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Search implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("IndexManager")).ID("Search").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("userID").ID("uint64")).Params(jen.ID("ids").Index().ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SearchForAdmin implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("IndexManager")).ID("SearchForAdmin").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string")).Params(jen.ID("ids").Index().ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("query"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Delete implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("IndexManager")).ID("Delete").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("id").ID("uint64")).Params(jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("id"),
			),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	return code
}
