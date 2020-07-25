package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(jen.Line())

	code.Add(
		jen.Comment("argIsNotPointer checks an argument and returns whether or not it is a pointer."),
		jen.Line(),
		jen.Func().ID("argIsNotPointer").Params(jen.ID("i").Interface()).Params(
			jen.ID("notAPointer").Bool(),
			jen.Err().Error(),
		).Block(
			jen.If(
				jen.ID("i").IsEqualTo().ID("nil").Or().
					Qual("reflect", "TypeOf").Call(
					jen.ID("i"),
				).Dot("Kind").Call().DoesNotEqual().Qual("reflect", "Ptr"),
			).Block(
				jen.Return().List(jen.True(), utils.Error("value is not a pointer")),
			),
			jen.Return().List(jen.False(),
				jen.Nil()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("argIsNotNil checks an argument and returns whether or not it is nil."),
		jen.Line(),
		jen.Func().ID("argIsNotNil").Params(jen.ID("i").Interface()).Params(
			jen.ID("isNil").Bool(),
			jen.Err().Error(),
		).Block(
			jen.If(jen.ID("i").IsEqualTo().ID("nil")).Block(
				jen.Return().List(
					jen.True(),
					utils.Error("value is nil"),
				),
			),
			jen.Return().List(jen.False(),
				jen.Nil()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("argIsNotPointerOrNil does what it says on the tin. This function"),
		jen.Line(),
		jen.Comment("is primarily useful for detecting if a destination value is valid"),
		jen.Line(),
		jen.Comment("before decoding an HTTP response, for instance."),
		jen.Line(),
		jen.Func().ID("argIsNotPointerOrNil").Params(jen.ID("i").Interface()).Params(jen.Error()).Block(
			jen.If(
				jen.List(jen.ID("nn"), jen.Err()).Assign().ID("argIsNotNil").Call(jen.ID("i")),
				jen.ID("nn").Or().Err().DoesNotEqual().ID("nil"),
			).Block(jen.Return().Err()),
			jen.Line(),
			jen.If(
				jen.List(jen.ID("np"), jen.Err()).Assign().ID("argIsNotPointer").Call(jen.ID("i")),
				jen.ID("np").Or().Err().DoesNotEqual().ID("nil"),
			).Block(jen.Return().Err()),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
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
		jen.Func().ID("unmarshalBody").Params(
			constants.CtxParam(),
			jen.ID(constants.ResponseVarName).PointerTo().Qual("net/http", "Response"),
			jen.ID("dest").Interface(),
		).Params(
			jen.Error(),
		).Block(
			utils.StartSpan(proj, false, "unmarshalBody"),
			jen.If(jen.Err().Assign().ID("argIsNotPointerOrNil").Call(jen.ID("dest")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().Err(),
			),
			jen.Line(),
			jen.List(
				jen.ID("bodyBytes"),
				jen.Err(),
			).Assign().Qual("io/ioutil", "ReadAll").
				Call(jen.ID(constants.ResponseVarName).Dot("Body")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().Err(),
			),
			jen.Line(),
			jen.If(jen.ID(constants.ResponseVarName).Dot("StatusCode").Op(">=").Qual("net/http", "StatusBadRequest")).Block(
				jen.ID("apiErr").Assign().AddressOf().Qual(proj.ModelsV1Package(), "ErrorResponse").Values(),
				jen.If(jen.Err().Equals().Qual("encoding/json", "Unmarshal").Call(
					jen.ID("bodyBytes"),
					jen.AddressOf().ID("apiErr"),
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
					jen.AddressOf().ID("dest"),
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

	code.Add(
		jen.Comment("createBodyFromStruct takes any value in and returns an io.Reader"),
		jen.Line(),
		jen.Comment("for placement within http.NewRequest's last argument."),
		jen.Line(),
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

	return code
}
