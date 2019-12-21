package users

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Var().ID("_").Qual("net/http", "Handler").Op("=").Parens(jen.Op("*").ID("MockHTTPHandler")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("MockHTTPHandler").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("MockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UserInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(
					jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("s").Dot("UserInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(
					jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("s").Dot("UserInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusBadRequest")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_PasswordUpdateInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(
					jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("s").Dot("PasswordUpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(
					jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUserCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("uint64").Call(jen.Lit(123)), jen.ID("nil")),
				jen.ID("s").Dot("database").Op("=").ID("mockDB"),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("s").Dot("PasswordUpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusBadRequest")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_TOTPSecretRefreshInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(
					jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("nil")),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("s").Dot("TOTPSecretRefreshInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusOK")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(
					jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("s").Dot("TOTPSecretRefreshInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusBadRequest")),
			)),
		),
		jen.Line(),
	)
	return ret
}
