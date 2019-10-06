package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func httpRoutesTestDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")
	utils.AddImports(ret)

	ret.Add(jen.Func().ID("Test_randString").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("actual").Op(":=").ID("randString").Call(),
			jen.ID("assert").Dot(
				"NotEmpty",
			).Call(jen.ID("t"), jen.ID("actual")),
		)),
	),
	)
	ret.Add(jen.Func().ID("buildRequest").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual("net/http", "Request")).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
		jen.ID("require").Dot(
			"NotNil",
		).Call(jen.ID("t"), jen.ID("req")),
		jen.ID("assert").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.Return().ID("req"),
	),
	)
	ret.Add(jen.Func().ID("Test_fetchUserID").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("expected"))),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"fetchUserID",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without context value present"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"fetchUserID",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestService_ListHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2Clients"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2ClientList",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ListHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with no rows returned"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2Clients"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2ClientList",
			)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ListHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2Clients"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2ClientList",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ListHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2Clients"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2ClientList",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ListHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestService_CreateHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("HashedPassword").Op(":").Lit("hashed_pass"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`blah`)), jen.ID("TwoFactorSecret").Op(":").Lit("SUPER SECRET")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientCreationInput",
			).Valuesln(jen.ID("UserLoginInput").Op(":").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456"))),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("a").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(),
			jen.ID("a").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"Password",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("a"),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("s").Dot(
				"oauth2ClientCounter",
			).Op("=").ID("uc"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusCreated"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with missing input"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error getting user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("HashedPassword").Op(":").Lit("hashed_pass"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`blah`)), jen.ID("TwoFactorSecret").Op(":").Lit("SUPER SECRET")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientCreationInput",
			).Valuesln(jen.ID("UserLoginInput").Op(":").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456"))),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid credentials"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("HashedPassword").Op(":").Lit("hashed_pass"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`blah`)), jen.ID("TwoFactorSecret").Op(":").Lit("SUPER SECRET")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientCreationInput",
			).Valuesln(jen.ID("UserLoginInput").Op(":").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456"))),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("a").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(),
			jen.ID("a").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"Password",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("false"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("a"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error validating password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("HashedPassword").Op(":").Lit("hashed_pass"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`blah`)), jen.ID("TwoFactorSecret").Op(":").Lit("SUPER SECRET")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientCreationInput",
			).Valuesln(jen.ID("UserLoginInput").Op(":").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456"))),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("a").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(),
			jen.ID("a").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"Password",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("a"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error creating oauth2 client"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("HashedPassword").Op(":").Lit("hashed_pass"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`blah`)), jen.ID("TwoFactorSecret").Op(":").Lit("SUPER SECRET")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientCreationInput",
			).Valuesln(jen.ID("UserLoginInput").Op(":").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456"))),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2Client",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("a").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(),
			jen.ID("a").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"Password",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("a"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("HashedPassword").Op(":").Lit("hashed_pass"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`blah`)), jen.ID("TwoFactorSecret").Op(":").Lit("SUPER SECRET")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientCreationInput",
			).Valuesln(jen.ID("UserLoginInput").Op(":").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456"))),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("a").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(),
			jen.ID("a").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"Password",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("a"),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("s").Dot(
				"oauth2ClientCounter",
			).Op("=").ID("uc"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("CreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusCreated"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestService_ReadHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("exampleOAuth2ClientID").Op(":=").ID("uint64").Call(jen.Lit(2)),
			jen.ID("s").Dot(
				"urlClientIDExtractor",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("exampleOAuth2ClientID"),
			),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleOAuth2ClientID"), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ReadHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with no rows found"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("exampleOAuth2ClientID").Op(":=").ID("uint64").Call(jen.Lit(2)),
			jen.ID("s").Dot(
				"urlClientIDExtractor",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("exampleOAuth2ClientID"),
			),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleOAuth2ClientID"), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(), jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ReadHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusNotFound"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching client from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("exampleOAuth2ClientID").Op(":=").ID("uint64").Call(jen.Lit(2)),
			jen.ID("s").Dot(
				"urlClientIDExtractor",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("exampleOAuth2ClientID"),
			),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleOAuth2ClientID"), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2Client",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ReadHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("exampleOAuth2ClientID").Op(":=").ID("uint64").Call(jen.Lit(2)),
			jen.ID("s").Dot(
				"urlClientIDExtractor",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("exampleOAuth2ClientID"),
			),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleOAuth2ClientID"), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ReadHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestService_ArchiveHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("exampleOAuth2ClientID").Op(":=").ID("uint64").Call(jen.Lit(2)),
			jen.ID("s").Dot(
				"urlClientIDExtractor",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("exampleOAuth2ClientID"),
			),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("ArchiveOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleOAuth2ClientID"), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("Decrement"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("s").Dot(
				"oauth2ClientCounter",
			).Op("=").ID("uc"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ArchiveHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with no rows found"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("exampleOAuth2ClientID").Op(":=").ID("uint64").Call(jen.Lit(2)),
			jen.ID("s").Dot(
				"urlClientIDExtractor",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("exampleOAuth2ClientID"),
			),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("ArchiveOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleOAuth2ClientID"), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ArchiveHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusNotFound"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error deleting record"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("exampleOAuth2ClientID").Op(":=").ID("uint64").Call(jen.Lit(2)),
			jen.ID("s").Dot(
				"urlClientIDExtractor",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("exampleOAuth2ClientID"),
			),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("ArchiveOAuth2Client"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleOAuth2ClientID"), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("userID"))),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"ArchiveHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),
	)
	return ret
}
