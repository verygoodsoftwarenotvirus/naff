package auth

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Func().ID("TestService_CookieAuthenticationMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("Username").Op(":").Lit("username")),
				jen.Line(),
				jen.ID("md").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "UserDataManager").Values(),
				jen.ID("md").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("exampleUser"), jen.ID("nil")),
				jen.ID("s").Dot("userDB").Op("=").ID("md"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot("buildAuthCookie").Call(jen.ID("exampleUser")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("cookie")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Op(":=").ID("s").Dot("CookieAuthenticationMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("Username").Op(":").Lit("username")),
				jen.ID("md").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "UserDataManager").Values(),
				jen.ID("md").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User")).Call(jen.ID("nil")), jen.ID("nil")),
				jen.ID("s").Dot("userDB").Op("=").ID("md"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot("buildAuthCookie").Call(jen.ID("exampleUser")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("cookie")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Op(":=").ID("s").Dot("CookieAuthenticationMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusUnauthorized")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without user attached"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Op(":=").ID("s").Dot("CookieAuthenticationMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_AuthenticationMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("exampleClient").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"),
					jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"),
					jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot("ID"),
				),
				jen.Line(),
				jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("oauth2ClientsService").Op("=").ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClient").Dot("BelongsTo")).Dot("Return").Call(jen.ID("exampleUser"), jen.ID("nil")),
				jen.ID("s").Dot("userDB").Op("=").ID("mockDB"),
				jen.Line(),
				jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("true")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path without allowing cookies"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("exampleClient").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"),
					jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"),
					jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot("ID"),
				),
				jen.Line(),
				jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("oauth2ClientsService").Op("=").ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClient").Dot("BelongsTo")).Dot("Return").Call(jen.ID("exampleUser"), jen.ID("nil")),
				jen.ID("s").Dot("userDB").Op("=").ID("mockDB"),
				jen.Line(),
				jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("false")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching client but able to use cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(1), jen.ID("Username").Op(":").Lit("username")),
				jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("oauth2ClientsService").Op("=").ID("ocv"),
				jen.Line(),
				jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot("buildAuthCookie").Call(jen.ID("exampleUser")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("c")),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleUser").Dot("ID")).Dot("Return").Call(jen.ID("exampleUser"), jen.ID("nil")),
				jen.ID("s").Dot("userDB").Op("=").ID("mockDB"),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("true")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("able to use cookies but error fetching user info"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(1), jen.ID("Username").Op(":").Lit("username")),
				jen.ID("exampleClient").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"),
					jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"),
					jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot("ID"),
				),
				jen.Line(),
				jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("oauth2ClientsService").Op("=").ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClient").Dot("BelongsTo")).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User")).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("userDB").Op("=").ID("mockDB"),
				jen.Line(),
				jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot("buildAuthCookie").Call(jen.ID("exampleUser")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("c")),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("true")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("no cookies allowed, with error fetching user info"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("exampleClient").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"),
					jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"),
					jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot("ID"),
				),
				jen.Line(),
				jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("oauth2ClientsService").Op("=").ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClient").Dot("BelongsTo")).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User")).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("userDB").Op("=").ID("mockDB"),
				jen.Line(),
				jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("false")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching client but able to use cookie but unable to decode cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("oauth2ClientsService").Op("=").ID("ocv"),
				jen.Line(),
				jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot("buildAuthCookie").Call(jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(1), jen.ID("Username").Op(":").Lit("username"))),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("c")),
				jen.Line(),
				jen.ID("cb").Op(":=").Op("&").ID("mockCookieEncoderDecoder").Values(),
				jen.ID("cb").Dot("On").Call(jen.Lit("Decode"), jen.ID("CookieName"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("cookieManager").Op("=").ID("cb"),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("true")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid authentication"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Call(jen.ID("nil")), jen.ID("nil")),
				jen.ID("s").Dot("oauth2ClientsService").Op("=").ID("ocv"),
				jen.Line(),
				jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("false")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("res").Dot("Code"), jen.Qual("net/http", "StatusUnauthorized")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("nightmare path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("exampleClient").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"),
					jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"),
					jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot("ID"),
				),
				jen.Line(),
				jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("oauth2ClientsService").Op("=").ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClient").Dot("BelongsTo")).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User")).Call(jen.ID("nil")), jen.ID("nil")),
				jen.ID("s").Dot("userDB").Op("=").ID("mockDB"),
				jen.Line(),
				jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("false")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_parseLoginInputFromForm").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("Password").Op(":").Lit("password"),
					jen.ID("TOTPToken").Op(":").Lit("123456"),
				),
				jen.Line(),
				jen.ID("req").Dot("Form").Op("=").Map(jen.ID("string")).Index().ID("string").Valuesln(
					jen.ID("UsernameFormKey").Op(":").Values(jen.ID("expected").Dot("Username")),
					jen.ID("PasswordFormKey").Op(":").Values(jen.ID("expected").Dot("Password")),
					jen.ID("TOTPTokenFormKey").Op(":").Values(jen.ID("expected").Dot("TOTPToken")),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("parseLoginInputFromForm").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("returns nil with error parsing form"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("req").Dot("URL").Dot("RawQuery").Op("=").Lit("%gh&%ij"),
				jen.ID("req").Dot("Form").Op("=").ID("nil"),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("parseLoginInputFromForm").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UserLoginInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("Password").Op(":").Lit("password"),
					jen.ID("TOTPToken").Op(":").Lit("1233456"),
				),
				jen.Line(),
				jen.Var().ID("b").Qual("bytes", "Buffer"),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("b")).Dot("Encode").Call(jen.ID("exampleInput"))),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Op("&").ID("b")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Op(":=").ID("s").Dot("UserLoginInputMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("Password").Op(":").Lit("password"),
					jen.ID("TOTPToken").Op(":").Lit("1233456"),
				),
				jen.Line(),
				jen.Var().ID("b").Qual("bytes", "Buffer"),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("b")).Dot("Encode").Call(jen.ID("exampleInput"))),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Op("&").ID("b")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("h").Op(":=").ID("s").Dot("UserLoginInputMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request but valid value attached to form"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("Password").Op(":").Lit("password"),
					jen.ID("TOTPToken").Op(":").Lit("1233456"),
				),
				jen.ID("form").Op(":=").Qual("net/url", "Values").Valuesln(
					jen.ID("UsernameFormKey").Op(":").Values(jen.ID("exampleInput").Dot("Username")),
					jen.ID("PasswordFormKey").Op(":").Values(jen.ID("exampleInput").Dot("Password")),
					jen.ID("TOTPTokenFormKey").Op(":").Values(jen.ID("exampleInput").Dot("TOTPToken")),
				),
				jen.Line(),
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
				),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Dot("Header").Dot("Set").Call(jen.Lit("Content-type"), jen.Lit("application/x-www-form-urlencoded")),
				jen.Line(),
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ed").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Op("=").ID("ed"),
				jen.Line(),
				jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Op(":=").ID("s").Dot("UserLoginInputMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_AdminMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserKey"),
						jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("IsAdmin").Op(":").ID("true")),
					),
				),
				jen.Line(),
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Op(":=").ID("s").Dot("AdminMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without user attached"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.Line(),
				jen.ID("h").Op(":=").ID("s").Dot("AdminMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with non-admin user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("require").Dot("NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserKey"),
						jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("IsAdmin").Op(":").ID("false")),
					),
				),
				jen.Line(),
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Values(),
				jen.Line(),
				jen.ID("h").Op(":=").ID("s").Dot("AdminMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code")),
			)),
		),
		jen.Line(),
	)
	return ret
}
