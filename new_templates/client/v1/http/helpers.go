package client

import jen "github.com/dave/jennifer/jen"

func helpersDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "models")
	ret.Add(jen.Line())

	ret.Add(
		jen.Comment("argIsNotPointer checks an argument and returns whether or not it is a pointer"),
		jen.Line(),
		jen.Func().Id("argIsNotPointer").Params(jen.Id("i").Interface()).Params(
			jen.Id("notAPointer").Id("bool"),
			jen.Id("err").Id("error"),
		).Block(
			jen.If(
				jen.Id("i").Op("==").Id("nil").
					Op("||").
					Qual("reflect", "TypeOf").Call(
					jen.Id("i"),
				).Dot("Kind").Call().Op("!=").Qual("reflect", "Ptr"),
			).Block(
				jen.Return().List(jen.Id("true"), jen.Qual("errors", "New").Call(
					jen.Lit("value is not a pointer"),
				),
				),
			),
			jen.Return().List(jen.Id("false"),
				jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("argIsNotNil checks an argument and returns whether or not it is nil"),
		jen.Line(),
		jen.Func().Id("argIsNotNil").Params(jen.Id("i").Interface()).Params(
			jen.Id("isNil").Id("bool"),
			jen.Id("err").Id("error"),
		).Block(
			jen.If(jen.Id("i").Op("==").Id("nil")).Block(
				jen.Return().List(
					jen.Id("true"),
					jen.Id("errors").Dot("New").Call(
						jen.Lit("value is nil"),
					),
				),
			),
			jen.Return().List(jen.Id("false"),
				jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("argIsNotPointerOrNil does what it says on the tin. This function"),
		jen.Line(),
		jen.Comment("is primarily useful for detecting if a destination value is valid"),
		jen.Line(),
		jen.Comment("before decoding an HTTP response, for instance"),
		jen.Line(),
		jen.Func().Id("argIsNotPointerOrNil").Params(jen.Id("i").Interface()).Params(jen.Id("error")).Block(
			jen.If(
				jen.List(
					jen.Id("nn"),
					jen.Id("err"),
				).Op(":=").Id("argIsNotNil").Call(
					jen.Id("i"),
				),
				jen.Id("nn").Op("||").Id("err").Op("!=").Id("nil"),
			).Block(jen.Return().Id("err")),
			jen.If(jen.List(
				jen.Id("np"),
				jen.Id("err"),
			).Op(":=").Id("argIsNotPointer").Call(
				jen.Id("i"),
			),
				jen.Id("np").
					Op("||").
					Id("err").Op("!=").Id("nil"),
			).Block(jen.Return().Id("err")),
			jen.Return().Id("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("unmarshalBody takes an HTTP response and JSON decodes its"),
		jen.Line(),
		jen.Comment("body into a destination value. `dest` must be a non-nil"),
		jen.Line(),
		jen.Comment("pointer to an object. Ideally, response is also not nil."),
		jen.Line(),
		jen.Comment("The error returned here should only ever be received in"),
		jen.Line(),
		jen.Comment("testing, and should never be encountered by an end-user."),
		jen.Line(),
		jen.Func().Id("unmarshalBody").Params(
			jen.Id("res").Op("*").Qual("net/http", "Response"),
			jen.Id("dest").Interface(),
		).Params(
			jen.Id("error"),
		).Block(
			jen.If(
				jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(
					jen.Id("dest"),
				),
				jen.Id("err").Op("!=").Id("nil"),
			).Block(
				jen.Return().Id("err"),
			),
			jen.List(
				jen.Id("bodyBytes"),
				jen.Id("err"),
			).Op(":=").Qual("io/ioutil", "ReadAll").
				Call(jen.Id("res").Dot("Body")),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Id("err"),
			),
			jen.If(jen.Id("res").Dot("StatusCode").Op(">=").Qual("net/http", "StatusBadRequest")).Block(
				jen.Id("apiErr").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "ErrorResponse").Values(),
				jen.If(jen.Id("err").Op("=").Qual("encoding/json", "Unmarshal").Call(
					jen.Id("bodyBytes"),
					jen.Op("&").Id("apiErr"),
				),
					jen.Id("err").Op("!=").Id("nil"),
				).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("unmarshaling error: %w"),
						jen.Id("err"),
					),
				),
				jen.Return().Id("apiErr"),
			),
			jen.If(
				jen.Id("err").Op("=").Qual("encoding/json", "Unmarshal").Call(
					jen.Id("bodyBytes"),
					jen.Op("&").Id("dest"),
				),
				jen.Id("err").Op("!=").Id("nil"),
			).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("unmarshaling body: %w"),
					jen.Id("err"),
				),
			), jen.Return().Id("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("createBodyFromStruct takes any value in and returns an io.Reader"),
		jen.Line(),
		jen.Comment("for placement within http.NewRequest's last argument."),
		jen.Line(),
		jen.Func().Id("createBodyFromStruct").Params(
			jen.Id("in").Interface(),
		).Params(
			jen.Qual("io", "Reader"),
			jen.Id("error"),
		).Block(
			jen.List(
				jen.Id("out"),
				jen.Id("err"),
			).Op(":=").Qual("encoding/json", "Marshal").Call(
				jen.Id("in"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"),
					jen.Id("err")),
			),
			jen.Return().List(jen.Qual("bytes", "NewReader").Call(
				jen.Id("out"),
			),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	return ret
}
