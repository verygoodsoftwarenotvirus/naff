package iterables

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, ret)

	sn := typ.Name.Singular()
	ret.Add(
		jen.Func().IDf("TestProvide%sDataManager", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.IDf("Provide%sDataManager", sn).Call(jen.Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call()),
			)),
		),
	)

	ret.Add(
		jen.Func().IDf("TestProvide%sDataServer", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.IDf("Provide%sDataServer", sn).Call(jen.ID("buildTestService").Call()),
			)),
		),
	)

	return ret
}
