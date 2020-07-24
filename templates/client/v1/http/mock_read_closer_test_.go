package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockReadCloserTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Underscore().Qual("io", "ReadCloser").Equals().Parens(jen.PointerTo().ID("ReadCloser")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadCloser is a mock io.ReadCloser for testprojects purposes."),
		jen.Line(),
		jen.Type().ID("ReadCloser").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("newMockReadCloser returns a new mock io.ReadCloser."),
		jen.Line(),
		jen.Func().ID("newMockReadCloser").Params().Params(jen.PointerTo().ID("ReadCloser")).Block(
			jen.Return().AddressOf().ID("ReadCloser").Values(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler implements the ReadHandler part of our ReadCloser."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("ReadCloser")).ID("Read").Params(jen.ID("b").Index().Byte()).Params(jen.ID("i").ID("int"), jen.Err().Error()).Block(
			jen.ID("retVals").Assign().ID("m").Dot("Called").Call(jen.ID("b")),
			jen.Return().List(jen.ID("retVals").Dot("Int").Call(jen.Zero()), jen.ID("retVals").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Close implements the Closer part of our ReadCloser."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("ReadCloser")).ID("Close").Params().Params(jen.Err().Error()).Block(
			jen.Return().ID("m").Dot("Called").Call().Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	return code
}
