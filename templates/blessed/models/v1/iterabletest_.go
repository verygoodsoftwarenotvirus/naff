package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableTestDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)
	sn := typ.Name.Singular()

	ret.Add(
		jen.Func().IDf("Test%s_Update", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("i").Op(":=").Op("&").ID(sn).Values(),
				jen.Line(),
				jen.ID("expected").Op(":=").Op("&").IDf("%sUpdateInput", sn).Valuesln(
					jen.ID("Name").Op(":").Lit("expected name"),
					jen.ID("Details").Op(":").Lit("expected details"),
				),
				jen.Line(),
				jen.ID("i").Dot("Update").Call(jen.ID("expected")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Name"), jen.ID("i").Dot("Name")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Details"), jen.ID("i").Dot("Details")),
			)),
		),
		jen.Line(),
	)
	return ret
}
