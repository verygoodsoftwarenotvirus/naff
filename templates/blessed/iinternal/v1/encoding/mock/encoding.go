package mock

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func encodingDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder").Op("=").Parens(jen.Op("*").ID("EncoderDecoder")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("EncoderDecoder is a mock EncoderDecoder"),
		jen.Line(),
		jen.Type().ID("EncoderDecoder").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("EncodeResponse satisfies our EncoderDecoder interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("EncodeResponse").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("v").Interface()).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("v")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("DecodeRequest satisfies our EncoderDecoder interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("DecodeRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("v").Interface()).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("req"), jen.ID("v")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
