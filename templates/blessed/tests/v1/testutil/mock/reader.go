package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func readerDotGo() *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("_").Qual("io", "ReadCloser").Op("=").Parens(jen.Op("*").ID("ReadCloser")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadCloser is a mock io.ReadCloser for testing purposes"),
		jen.Line(),
		jen.Type().ID("ReadCloser").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("NewMockReadCloser returns a new mock io.ReadCloser"),
		jen.Line(),
		jen.Func().ID("NewMockReadCloser").Params().Params(jen.Op("*").ID("ReadCloser")).Block(
			jen.Return().Op("&").ID("ReadCloser").Values(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler implements the ReadHandler part of our ReadCloser"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ReadCloser")).ID("Read").Params(jen.ID("b").Index().ID("byte")).Params(jen.ID("i").ID("int"), jen.ID("err").ID("error")).Block(
			jen.ID("retVals").Op(":=").ID("m").Dot("Called").Call(jen.ID("b")),
			jen.Return().List(jen.ID("retVals").Dot("Int").Call(jen.Lit(0)), jen.ID("retVals").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Close implements the Closer part of our ReadCloser"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ReadCloser")).ID("Close").Params().Params(jen.ID("err").ID("error")).Block(
			jen.Return().ID("m").Dot(
				"Called",
			).Call().Dot(
				"Error",
			).Call(jen.Lit(1)),
		),
		jen.Line(),
	)
	return ret
}
