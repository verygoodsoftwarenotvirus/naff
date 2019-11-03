package client

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(jen.Line())

	ret.Add(utils.Comments("argIsNotPointer checks an argument and returns whether or not it is a pointer")...)
	ret.Add(
		jen.Func().ID("argIsNotPointer").Params(jen.ID("i").Interface()).Params(
			jen.ID("notAPointer").ID("bool"),
			jen.ID("err").ID("error"),
		).Block(
			jen.If(
				jen.ID("i").Op("==").ID("nil").
					Op("||").
					Qual("reflect", "TypeOf").Call(
					jen.ID("i"),
				).Dot("Kind").Call().Op("!=").Qual("reflect", "Ptr"),
			).Block(
				jen.Return().List(
					jen.ID("true"),
					jen.Qual("errors", "New").Call(jen.Lit("value is not a pointer")),
				),
			),
			jen.Return().List(jen.ID("false"),
				jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(utils.Comments("argIsNotNil checks an argument and returns whether or not it is nil")...)
	ret.Add(
		jen.Func().ID("argIsNotNil").Params(jen.ID("i").Interface()).Params(
			jen.ID("isNil").ID("bool"),
			jen.ID("err").ID("error"),
		).Block(
			jen.If(jen.ID("i").Op("==").ID("nil")).Block(
				jen.Return().List(
					jen.ID("true"),
					jen.Qual("errors", "New").Call(
						jen.Lit("value is nil"),
					),
				),
			),
			jen.Return().List(jen.ID("false"),
				jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(utils.Comments(
		"argIsNotPointerOrNil does what it says on the tin. This function",
		"is primarily useful for detecting if a destination value is valid",
		"before decoding an HTTP response, for instance",
	)...)
	ret.Add(
		jen.Func().ID("argIsNotPointerOrNil").Params(jen.ID("i").Interface()).Params(jen.ID("error")).Block(
			jen.If(
				jen.List(jen.ID("nn"), jen.ID("err")).Op(":=").ID("argIsNotNil").Call(jen.ID("i")),
				jen.ID("nn").Op("||").ID("err").Op("!=").ID("nil"),
			).Block(jen.Return().ID("err")),
			jen.Line(),
			jen.If(
				jen.List(jen.ID("np"), jen.ID("err")).Op(":=").ID("argIsNotPointer").Call(jen.ID("i")),
				jen.ID("np").Op("||").ID("err").Op("!=").ID("nil"),
			).Block(jen.Return().ID("err")),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	ret.Add(utils.Comments(
		"unmarshalBody takes an HTTP response and JSON decodes its",
		"body into a destination value. `dest` must be a non-nil",
		"pointer to an object. Ideally, response is also not nil.",
		"The error returned here should only ever be received in",
		"testing, and should never be encountered by an end-user.",
	)...)
	ret.Add(
		jen.Func().ID("unmarshalBody").Params(
			jen.ID("res").Op("*").Qual("net/http", "Response"),
			jen.ID("dest").Interface(),
		).Params(
			jen.ID("error"),
		).Block(
			jen.If(jen.ID("err").Op(":=").ID("argIsNotPointerOrNil").Call(jen.ID("dest")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().ID("err"),
			),
			jen.Line(),
			jen.List(
				jen.ID("bodyBytes"),
				jen.ID("err"),
			).Op(":=").Qual("io/ioutil", "ReadAll").
				Call(jen.ID("res").Dot("Body")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().ID("err"),
			),
			jen.Line(),
			jen.If(jen.ID("res").Dot("StatusCode").Op(">=").Qual("net/http", "StatusBadRequest")).Block(
				jen.ID("apiErr").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "ErrorResponse").Values(),
				jen.If(jen.ID("err").Op("=").Qual("encoding/json", "Unmarshal").Call(
					jen.ID("bodyBytes"),
					jen.Op("&").ID("apiErr"),
				),
					jen.ID("err").Op("!=").ID("nil"),
				).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("unmarshaling error: %w"),
						jen.ID("err"),
					),
				),
				jen.Return().ID("apiErr"),
			),
			jen.Line(),
			jen.If(
				jen.ID("err").Op("=").Qual("encoding/json", "Unmarshal").Call(
					jen.ID("bodyBytes"),
					jen.Op("&").ID("dest"),
				),
				jen.ID("err").Op("!=").ID("nil"),
			).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("unmarshaling body: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	ret.Add(utils.Comments(
		"createBodyFromStruct takes any value in and returns an io.Reader",
		"for placement within http.NewRequest's last argument.",
	)...)
	ret.Add(
		jen.Func().ID("createBodyFromStruct").Params(
			jen.ID("in").Interface(),
		).Params(
			jen.Qual("io", "Reader"),
			jen.ID("error"),
		).Block(
			jen.List(
				jen.ID("out"),
				jen.ID("err"),
			).Op(":=").Qual("encoding/json", "Marshal").Call(
				jen.ID("in"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.ID("err")),
			),
			jen.Return().List(jen.Qual("bytes", "NewReader").Call(
				jen.ID("out"),
			),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	return ret
}
