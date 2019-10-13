package users

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func httpRoutesTestDotGo() *jen.File {
	ret := jen.NewFile("users")
	ret.Add(jen.Null(),

		jen.Line(),
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

		jen.Line(),
	)
	ret.Add(jen.Func().ID("Test_randString").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("randString").Call(),
			jen.ID("assert").Dot(
				"NotEmpty",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_validateCredentialChangeRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
			jen.ID("examplePassword").Op(":=").Lit("password"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"HashedPassword",
			), jen.ID("examplePassword"), jen.ID("expected").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleTOTPToken"), jen.ID("expected").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("s").Dot(
				"validateCredentialChangeRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("examplePassword"), jen.ID("exampleTOTPToken")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("sc")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with no rows found in database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
			jen.ID("examplePassword").Op(":=").Lit("password"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("s").Dot(
				"validateCredentialChangeRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("examplePassword"), jen.ID("exampleTOTPToken")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusNotFound"), jen.ID("sc")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
			jen.ID("examplePassword").Op(":=").Lit("password"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("s").Dot(
				"validateCredentialChangeRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("examplePassword"), jen.ID("exampleTOTPToken")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("sc")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error validating login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
			jen.ID("examplePassword").Op(":=").Lit("password"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"HashedPassword",
			), jen.ID("examplePassword"), jen.ID("expected").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleTOTPToken"), jen.ID("expected").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("false"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("s").Dot(
				"validateCredentialChangeRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("examplePassword"), jen.ID("exampleTOTPToken")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("sc")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
			jen.ID("examplePassword").Op(":=").Lit("password"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expected").Dot(
				"HashedPassword",
			), jen.ID("examplePassword"), jen.ID("expected").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleTOTPToken"), jen.ID("expected").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("false"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("s").Dot(
				"validateCredentialChangeRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("examplePassword"), jen.ID("exampleTOTPToken")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("sc")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_List").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUsers"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"UserList",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
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
		).Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUsers"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"UserList",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
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
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUsers"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"UserList",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
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

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_Create").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"UserInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").ID("exampleInput").Dot(
				"Username",
			), jen.ID("HashedPassword").Op(":").Lit("blahblah")),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Password",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser").Dot(
				"HashedPassword",
			), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.ID("db").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("db").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("db"),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)),
			jen.ID("s").Dot(
				"userCounter",
			).Op("=").ID("mc"),
			jen.ID("r").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Valuesln(),
			jen.ID("r").Dot(
				"On",
			).Call(jen.Lit("Report"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("s").Dot(
				"reporter",
			).Op("=").ID("r"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserCreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"userCreationEnabled",
			).Op("=").ID("true"),
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
		).Call(jen.Lit("with user creation disabled"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("s").Dot(
				"userCreationEnabled",
			).Op("=").ID("false"),
			jen.ID("s").Dot(
				"CreateHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusForbidden"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with missing input"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("s").Dot(
				"userCreationEnabled",
			).Op("=").ID("true"),
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
		).Call(jen.Lit("with error hashing password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"UserInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").ID("exampleInput").Dot(
				"Username",
			), jen.ID("HashedPassword").Op(":").Lit("blahblah")),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Password",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser").Dot(
				"HashedPassword",
			), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserCreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"userCreationEnabled",
			).Op("=").ID("true"),
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
		).Call(jen.Lit("with error creating entry in database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"UserInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").ID("exampleInput").Dot(
				"Username",
			), jen.ID("HashedPassword").Op(":").Lit("blahblah")),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Password",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser").Dot(
				"HashedPassword",
			), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.ID("db").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("db").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("db"),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserCreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"userCreationEnabled",
			).Op("=").ID("true"),
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
		).Call(jen.Lit("with pre-existing entry in database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"UserInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").ID("exampleInput").Dot(
				"Username",
			), jen.ID("HashedPassword").Op(":").Lit("blahblah")),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Password",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser").Dot(
				"HashedPassword",
			), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.ID("db").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("db").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/client", "ErrUserExists")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("db"),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserCreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"userCreationEnabled",
			).Op("=").ID("true"),
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
		).Call(jen.Lit("with error encoding response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"UserInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").ID("exampleInput").Dot(
				"Username",
			), jen.ID("HashedPassword").Op(":").Lit("blahblah")),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"Password",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser").Dot(
				"HashedPassword",
			), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.ID("db").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("db").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("CreateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput")).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("db"),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Increment"), jen.ID("mock").Dot(
				"Anything",
			)),
			jen.ID("s").Dot(
				"userCounter",
			).Op("=").ID("mc"),
			jen.ID("r").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Valuesln(),
			jen.ID("r").Dot(
				"On",
			).Call(jen.Lit("Report"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("s").Dot(
				"reporter",
			).Op("=").ID("r"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserCreationMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"userCreationEnabled",
			).Op("=").ID("true"),
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

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_Read").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
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
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(), jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
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
		).Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
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
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
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

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_NewTOTPSecret").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"TOTPSecretRefreshInput",
			).Valuesln(),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("TOTPSecretRefreshMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"CurrentPassword",
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
			).Op("=").ID("auth"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("s").Dot(
				"NewTOTPSecretHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusAccepted"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without input attached to request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("s").Dot(
				"NewTOTPSecretHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with input attached but without user information"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"TOTPSecretRefreshInput",
			).Valuesln(),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("TOTPSecretRefreshMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"NewTOTPSecretHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error validating login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"TOTPSecretRefreshInput",
			).Valuesln(),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("TOTPSecretRefreshMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"CurrentPassword",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("false"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("s").Dot(
				"NewTOTPSecretHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error updating in database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"TOTPSecretRefreshInput",
			).Valuesln(),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("TOTPSecretRefreshMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"CurrentPassword",
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
			).Op("=").ID("auth"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("s").Dot(
				"NewTOTPSecretHandler",
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
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"TOTPSecretRefreshInput",
			).Valuesln(),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("TOTPSecretRefreshMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"CurrentPassword",
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
			).Op("=").ID("auth"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("s").Dot(
				"NewTOTPSecretHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusAccepted"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_UpdatePassword").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"PasswordUpdateInput",
			).Valuesln(jen.ID("NewPassword").Op(":").Lit("new_password"), jen.ID("CurrentPassword").Op(":").Lit("old_password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("PasswordChangeMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"CurrentPassword",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"NewPassword",
			)).Dot(
				"Return",
			).Call(jen.Lit("blah"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("s").Dot(
				"UpdatePasswordHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusAccepted"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without input attached to request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("s").Dot(
				"UpdatePasswordHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with input but without user info"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"PasswordUpdateInput",
			).Valuesln(jen.ID("NewPassword").Op(":").Lit("new_password"), jen.ID("CurrentPassword").Op(":").Lit("old_password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("PasswordChangeMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("s").Dot(
				"UpdatePasswordHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error validating login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"PasswordUpdateInput",
			).Valuesln(jen.ID("NewPassword").Op(":").Lit("new_password"), jen.ID("CurrentPassword").Op(":").Lit("old_password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("PasswordChangeMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"CurrentPassword",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("false"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.ID("s").Dot(
				"UpdatePasswordHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error hashing password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"PasswordUpdateInput",
			).Valuesln(jen.ID("NewPassword").Op(":").Lit("new_password"), jen.ID("CurrentPassword").Op(":").Lit("old_password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("PasswordChangeMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"CurrentPassword",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"NewPassword",
			)).Dot(
				"Return",
			).Call(jen.Lit("blah"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("s").Dot(
				"UpdatePasswordHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error updating user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(123)), jen.ID("HashedPassword").Op(":").Lit("not really lol"), jen.ID("Salt").Op(":").Index().ID("byte").Call(jen.Lit(`nah`)), jen.ID("TwoFactorSecret").Op(":").Lit("still no")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"PasswordUpdateInput",
			).Valuesln(jen.ID("NewPassword").Op(":").Lit("new_password"), jen.ID("CurrentPassword").Op(":").Lit("old_password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("PasswordChangeMiddlewareCtxKey"), jen.ID("exampleInput"))),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("exampleUser").Dot(
				"ID",
			))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Valuesln(),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"HashedPassword",
			), jen.ID("exampleInput").Dot(
				"CurrentPassword",
			), jen.ID("exampleUser").Dot(
				"TwoFactorSecret",
			), jen.ID("exampleInput").Dot(
				"TOTPToken",
			), jen.ID("exampleUser").Dot(
				"Salt",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("auth").Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleInput").Dot(
				"NewPassword",
			)).Dot(
				"Return",
			).Call(jen.Lit("blah"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Op("=").ID("auth"),
			jen.ID("s").Dot(
				"UpdatePasswordHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_Archive").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expectedUserID"),
			),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("ArchiveUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUserID")).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("r").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "Reporter").Valuesln(),
			jen.ID("r").Dot(
				"On",
			).Call(jen.Lit("Report"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("mc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Valuesln(),
			jen.ID("mc").Dot(
				"On",
			).Call(jen.Lit("Decrement"), jen.ID("mock").Dot(
				"Anything",
			)),
			jen.ID("s").Dot(
				"userCounter",
			).Op("=").ID("mc"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
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
		).Call(jen.Lit("with error updating database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("s").Dot(
				"userIDFetcher",
			).Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().ID("expectedUserID"),
			),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("ArchiveUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUserID")).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
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

		jen.Line(),
	)
	return ret
}
