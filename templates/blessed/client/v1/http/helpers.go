package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(jen.Line())

	ret.Add(utils.Comments("argIsNotPointer checks an argument and returns whether or not it is a pointer")...)
	ret.Add(
		jen.Func().ID("argIsNotPointer").Params(jen.ID("i").Interface()).Params(
			jen.ID("notAPointer").Bool(),
			jen.Err().Error(),
		).Block(
			jen.If(
				jen.ID("i").Op("==").ID("nil").
					Op("||").
					Qual("reflect", "TypeOf").Call(
					jen.ID("i"),
				).Dot("Kind").Call().DoesNotEqual().Qual("reflect", "Ptr"),
			).Block(
				jen.Return().List(
					jen.ID("true"),
					jen.Qual("errors", "New").Call(jen.Lit("value is not a pointer")),
				),
			),
			jen.Return().List(jen.ID("false"),
				jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(utils.Comments("argIsNotNil checks an argument and returns whether or not it is nil")...)
	ret.Add(
		jen.Func().ID("argIsNotNil").Params(jen.ID("i").Interface()).Params(
			jen.ID("isNil").Bool(),
			jen.Err().Error(),
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
				jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(utils.Comments(
		"argIsNotPointerOrNil does what it says on the tin. This function",
		"is primarily useful for detecting if a destination value is valid",
		"before decoding an HTTP response, for instance",
	)...)
	ret.Add(
		jen.Func().ID("argIsNotPointerOrNil").Params(jen.ID("i").Interface()).Params(jen.Error()).Block(
			jen.If(
				jen.List(jen.ID("nn"), jen.Err()).Assign().ID("argIsNotNil").Call(jen.ID("i")),
				jen.ID("nn").Or().ID("err").DoesNotEqual().ID("nil"),
			).Block(jen.Return().ID("err")),
			jen.Line(),
			jen.If(
				jen.List(jen.ID("np"), jen.Err()).Assign().ID("argIsNotPointer").Call(jen.ID("i")),
				jen.ID("np").Or().ID("err").DoesNotEqual().ID("nil"),
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
			utils.CtxParam(),
			jen.ID("res").ParamPointer().Qual("net/http", "Response"),
			jen.ID("dest").Interface(),
		).Params(
			jen.Error(),
		).Block(
			utils.StartSpan(proj, false, "unmarshalBody"),
			jen.If(jen.Err().Assign().ID("argIsNotPointerOrNil").Call(jen.ID("dest")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().ID("err"),
			),
			jen.Line(),
			jen.List(
				jen.ID("bodyBytes"),
				jen.Err(),
			).Assign().Qual("io/ioutil", "ReadAll").
				Call(jen.ID("res").Dot("Body")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().ID("err"),
			),
			jen.Line(),
			jen.If(jen.ID("res").Dot("StatusCode").Op(">=").Qual("net/http", "StatusBadRequest")).Block(
				jen.ID("apiErr").Assign().VarPointer().Qual(proj.ModelsV1Package(), "ErrorResponse").Values(),
				jen.If(jen.Err().Equals().Qual("encoding/json", "Unmarshal").Call(
					jen.ID("bodyBytes"),
					jen.VarPointer().ID("apiErr"),
				),
					jen.Err().DoesNotEqual().ID("nil"),
				).Block(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("unmarshaling error: %w"),
						jen.Err(),
					),
				),
				jen.Return().ID("apiErr"),
			),
			jen.Line(),
			jen.If(
				jen.Err().Equals().Qual("encoding/json", "Unmarshal").Call(
					jen.ID("bodyBytes"),
					jen.VarPointer().ID("dest"),
				),
				jen.Err().DoesNotEqual().ID("nil"),
			).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("unmarshaling body: %w"),
					jen.Err(),
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
			jen.Error(),
		).Block(
			jen.List(
				jen.ID("out"),
				jen.Err(),
			).Assign().Qual("encoding/json", "Marshal").Call(
				jen.ID("in"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(),
					jen.Err()),
			),
			jen.Return().List(jen.Qual("bytes", "NewReader").Call(
				jen.ID("out"),
			),
				jen.Nil(),
			),
		),
		jen.Line(),
	)

	return ret
}
