package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual("net/http", "Handler").Equals().Parens(jen.PointerTo().ID("MockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("MockHTTPHandler").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("MockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_CreationInputMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusBadRequest"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UpdateInputMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("UpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("UpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusBadRequest"), nil),
			),
		),
		jen.Line(),
	)
	return ret
}
