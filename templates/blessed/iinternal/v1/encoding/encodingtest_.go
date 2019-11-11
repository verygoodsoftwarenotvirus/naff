package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func encodingTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("encoding")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Type().ID("example").Struct(
			jen.ID("Name").ID("string").Tag(map[string]string{"json": "name", "xml": "name"}),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestServerEncoderDecoder_EncodeResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectation").Op(":=").Lit("name"),
				jen.ID("ex").Op(":=").Op("&").ID("example").Values(jen.ID("Name").Op(":").ID("expectation")),
				jen.ID("ed").Op(":=").ID("ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("err").Op(":=").ID("ed").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("ex")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("res").Dot("Body").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("{%q:%q}\n"), jen.Lit("name"), jen.ID("ex").Dot("Name"))),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("as XML"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectation").Op(":=").Lit("name"),
				jen.ID("ex").Op(":=").Op("&").ID("example").Values(jen.ID("Name").Op(":").ID("expectation")),
				jen.ID("ed").Op(":=").ID("ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("res").Dot("Header").Call().Dot("Set").Call(jen.ID("ContentTypeHeader"), jen.Lit("application/xml")),
				jen.Line(),
				jen.ID("err").Op(":=").ID("ed").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("ex")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("fmt", "Sprintf").Call(jen.Lit(`<example><name>%s</name></example>`), jen.ID("expectation")), jen.ID("res").Dot("Body").Dot("String").Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestServerEncoderDecoder_DecodeRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectation").Op(":=").Lit("name"),
				jen.ID("e").Op(":=").Op("&").ID("example").Values(jen.ID("Name").Op(":").ID("expectation")),
				jen.ID("ed").Op(":=").ID("ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.List(jen.ID("bs"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("e")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Qual("bytes", "NewReader").Call(jen.ID("bs"))),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.Var().ID("x").ID("example"),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("ed").Dot("DecodeRequest").Call(jen.ID("req"), jen.Op("&").ID("x"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("x").Dot("Name"), jen.ID("e").Dot("Name")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("as XML"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectation").Op(":=").Lit("name"),
				jen.ID("e").Op(":=").Op("&").ID("example").Values(jen.ID("Name").Op(":").ID("expectation")),
				jen.ID("ed").Op(":=").ID("ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.List(jen.ID("bs"), jen.ID("err")).Op(":=").Qual("encoding/xml", "Marshal").Call(jen.ID("e")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Qual("bytes", "NewReader").Call(jen.ID("bs"))),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("req").Dot("Header").Dot("Set").Call(jen.ID("ContentTypeHeader"), jen.ID("XMLContentType")),
				jen.Line(),
				jen.Var().ID("x").ID("example"),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("ed").Dot("DecodeRequest").Call(jen.ID("req"), jen.Op("&").ID("x"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("x").Dot("Name"), jen.ID("e").Dot("Name")),
			)),
		),
		jen.Line(),
	)
	return ret
}
