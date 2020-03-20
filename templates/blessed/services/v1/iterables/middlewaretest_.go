package iterables

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().ID("_").Qual("net/http", "Handler").Op("=").Parens(jen.Op("*").ID("mockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockHTTPHandler").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_CreationInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("actual").Op(":=").ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("actual").Op(":=").ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusBadRequest")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UpdateInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("actual").Op(":=").ID("s").Dot("UpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("actual").Op(":=").ID("s").Dot("UpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusBadRequest")),
			)),
		),
		jen.Line(),
	)
	return ret
}
