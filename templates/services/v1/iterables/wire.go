package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code)

	code.Add(buildWireProviders(typ)...)
	code.Add(buildWireProvideSomethingDataManager(proj, typ)...)
	code.Add(buildWireProvideSomethingDataServer(proj, typ)...)

	return code
}

func buildWireProviders(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID(fmt.Sprintf("Provide%sService", typ.Name.Plural())),
				jen.ID(fmt.Sprintf("Provide%sDataManager", sn)),
				jen.ID(fmt.Sprintf("Provide%sDataServer", sn)),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.IDf("Provide%sServiceSearchIndex", pn)
					} else {
						return jen.Null()
					}
				}(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildWireProvideSomethingDataManager(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Commentf("Provide%sDataManager turns a database into an %sDataManager.", sn, sn),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sDataManager", sn)).Params(jen.ID("db").Qual(proj.DatabaseV1Package(), "DataManager")).Params(jen.Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataManager", sn))).Block(
			jen.Return().ID("db"),
		),
		jen.Line(),
	}

	return lines
}

func buildWireProvideSomethingDataServer(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Commentf("Provide%sDataServer is an arbitrary function for dependency injection's sake.", sn),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sDataServer", sn)).Params(jen.ID("s").PointerTo().ID("Service")).Params(jen.Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataServer", sn))).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	}

	return lines
}
