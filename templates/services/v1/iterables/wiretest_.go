package iterables

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code)

	sn := typ.Name.Singular()
	code.Add(
		jen.Func().IDf("TestProvide%sDataManager", sn).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
				jen.IDf("Provide%sDataManager", sn).Call(jen.Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call()),
			)),
		),
	)

	code.Add(
		jen.Func().IDf("TestProvide%sDataServer", sn).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
				jen.IDf("Provide%sDataServer", sn).Call(jen.ID("buildTestService").Call()),
			)),
		),
	)

	return code
}
