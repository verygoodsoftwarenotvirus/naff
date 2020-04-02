package iterables

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(pkg, ret)

	sn := typ.Name.Singular()

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services"),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				jen.ID(fmt.Sprintf("Provide%sService", typ.Name.Plural())),
				jen.ID(fmt.Sprintf("Provide%sDataManager", sn)),
				jen.ID(fmt.Sprintf("Provide%sDataServer", sn)),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Provide%sDataManager turns a database into an %sDataManager", sn, sn),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sDataManager", sn)).Params(jen.ID("db").Qual(pkg.DatabaseV1Package(), "Database")).Params(jen.Qual(pkg.ModelsV1Package(), fmt.Sprintf("%sDataManager", sn))).Block(
			jen.Return().ID("db"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Provide%sDataServer is an arbitrary function for dependency injection's sake", sn),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sDataServer", sn)).Params(jen.ID("s").Op("*").ID("Service")).Params(jen.Qual(pkg.ModelsV1Package(), fmt.Sprintf("%sDataServer", sn))).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	return ret
}
