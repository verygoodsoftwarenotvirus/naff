package iterables

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkgRoot string, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(pkgRoot, []models.DataType{typ}, ret)

	sn := typ.Name.Singular()

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(
				jen.ID(fmt.Sprintf("Provide%sService", typ.Name.Plural())),
				jen.ID(fmt.Sprintf("Provide%sDataManager", sn)),
				jen.ID(fmt.Sprintf("Provide%sDataServer", sn)),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Provide%sDataManager turns a database into an %sDataManager", sn, sn)),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sDataManager", sn)).Params(jen.ID("db").Qual(filepath.Join(pkgRoot, "database/v1"), "Database")).Params(jen.Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sDataManager", sn))).Block(
			jen.Return().ID("db"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Provide%sDataServer is an arbitrary function for dependency injection's sake", sn)),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sDataServer", sn)).Params(jen.ID("s").Op("*").ID("Service")).Params(jen.Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sDataServer", sn))).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	return ret
}
