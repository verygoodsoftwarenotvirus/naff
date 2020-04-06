package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func encodingTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("encoding")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Type().ID("example").Struct(
			jen.ID("Name").String().Tag(map[string]string{"json": "name", "xml": "name"}),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestServerEncoderDecoder_EncodeResponse").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectation").Assign().Lit("name"),
				jen.ID("ex").Assign().VarPointer().ID("example").Values(jen.ID("Name").MapAssign().ID("expectation")),
				jen.ID("ed").Assign().ID("ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Err().Assign().ID("ed").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("ex")),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("res").Dot("Body").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(jen.Lit("{%q:%q}\n"), jen.Lit("name"), jen.ID("ex").Dot("Name")), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"as XML",
				jen.ID("expectation").Assign().Lit("name"),
				jen.ID("ex").Assign().VarPointer().ID("example").Values(jen.ID("Name").MapAssign().ID("expectation")),
				jen.ID("ed").Assign().ID("ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("res").Dot("Header").Call().Dot("Set").Call(jen.ID("ContentTypeHeader"), jen.Lit("application/xml")),
				jen.Line(),
				jen.Err().Assign().ID("ed").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("ex")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("fmt", "Sprintf").Call(jen.Lit(`<example><name>%s</name></example>`), jen.ID("expectation")), jen.ID("res").Dot("Body").Dot("String").Call(), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestServerEncoderDecoder_DecodeRequest").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectation").Assign().Lit("name"),
				jen.ID("e").Assign().VarPointer().ID("example").Values(jen.ID("Name").MapAssign().ID("expectation")),
				jen.ID("ed").Assign().ID("ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.List(jen.ID("bs"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("e")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Qual("bytes", "NewReader").Call(jen.ID("bs"))),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.Var().ID("x").ID("example"),
				utils.AssertNoError(jen.ID("ed").Dot("DecodeRequest").Call(jen.ID("req"), jen.VarPointer().ID("x")), nil),
				utils.AssertEqual(jen.ID("x").Dot("Name"), jen.ID("e").Dot("Name"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"as XML",
				jen.ID("expectation").Assign().Lit("name"),
				jen.ID("e").Assign().VarPointer().ID("example").Values(jen.ID("Name").MapAssign().ID("expectation")),
				jen.ID("ed").Assign().ID("ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.List(jen.ID("bs"), jen.Err()).Assign().Qual("encoding/xml", "Marshal").Call(jen.ID("e")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Qual("bytes", "NewReader").Call(jen.ID("bs"))),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("req").Dot("Header").Dot("Set").Call(jen.ID("ContentTypeHeader"), jen.ID("XMLContentType")),
				jen.Line(),
				jen.Var().ID("x").ID("example"),
				utils.AssertNoError(jen.ID("ed").Dot("DecodeRequest").Call(jen.ID("req"), jen.VarPointer().ID("x")), nil),
				utils.AssertEqual(jen.ID("x").Dot("Name"), jen.ID("e").Dot("Name"), nil),
			),
		),
		jen.Line(),
	)
	return ret
}
