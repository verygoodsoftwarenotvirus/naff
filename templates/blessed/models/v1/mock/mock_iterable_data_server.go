package mock

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockIterableDataServerDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkg, ret)

	sn := typ.Name.Singular()

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sDataServer", sn)).Op("=").Parens(jen.Op("*").IDf("%sDataServer", sn)).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("%sDataServer is a mocked models.%sDataServer for testing", sn, sn),
		jen.Line(),
		jen.Type().IDf("%sDataServer", sn).Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreationInputMiddleware implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataServer", sn)).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateInputMiddleware implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataServer", sn)).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataServer", sn)).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataServer", sn)).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ExistenceHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataServer", sn)).ID("ExistenceHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataServer", sn)).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataServer", sn)).ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataServer", sn)).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)
	return ret
}
