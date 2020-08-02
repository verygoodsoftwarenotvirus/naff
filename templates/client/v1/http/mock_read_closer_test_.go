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

	code.Add(buildMockReadCloserInterfaceAssurance()...)
	code.Add(buildMockReadCloserDecl()...)
	code.Add(buildNewMockReadCloser()...)
	code.Add(buildMockReadCloserReadHandler()...)
	code.Add(buildMockReadCloserClose()...)

	return code
}

func buildMockReadCloserInterfaceAssurance() []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual("io", "ReadCloser").Equals().Parens(jen.PointerTo().ID("ReadCloser")).Call(jen.Nil()),
		jen.Line(),
	}

	return lines
}

func buildMockReadCloserDecl() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ReadCloser is a mock io.ReadCloser for testing purposes."),
		jen.Line(),
		jen.Type().ID("ReadCloser").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildNewMockReadCloser() []jen.Code {
	lines := []jen.Code{
		jen.Comment("newMockReadCloser returns a new mock io.ReadCloser."),
		jen.Line(),
		jen.Func().ID("newMockReadCloser").Params().Params(jen.PointerTo().ID("ReadCloser")).Block(
			jen.Return().AddressOf().ID("ReadCloser").Values(),
		),
		jen.Line(),
	}

	return lines
}

func buildMockReadCloserReadHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ReadHandler implements the ReadHandler part of our ReadCloser."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("ReadCloser")).ID("Read").Params(jen.ID("b").Index().Byte()).Params(jen.ID("i").ID("int"), jen.Err().Error()).Block(
			jen.ID("retVals").Assign().ID("m").Dot("Called").Call(jen.ID("b")),
			jen.Return().List(jen.ID("retVals").Dot("Int").Call(jen.Zero()), jen.ID("retVals").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildMockReadCloserClose() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Close implements the Closer part of our ReadCloser."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("ReadCloser")).ID("Close").Params().Params(jen.Err().Error()).Block(
			jen.Return().ID("m").Dot("Called").Call().Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
