package iterables

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestProvideSomethingDataManager(proj, typ)...)
	code.Add(buildTestProvideSomethingDataServer(typ)...)

	return code
}

func buildTestProvideSomethingDataManager(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Func().IDf("TestProvide%sDataManager", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.IDf("Provide%sDataManager", sn).Call(jen.Qual(proj.DatabasePackage(), "BuildMockDatabase").Call()),
			)),
		),
	}

	return lines
}

func buildTestProvideSomethingDataServer(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Func().IDf("TestProvide%sDataServer", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.IDf("Provide%sDataServer", sn).Call(jen.ID("buildTestService").Call()),
			)),
		),
	}

	return lines
}
