package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersServiceTestDotGo() *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Service")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("expectedUserCount").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
			jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
				jen.Lit("GetUserCount"),
				jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				jen.Parens(jen.Op("*").ID("models").Dot("QueryFilter")).Call(jen.ID("nil")),
			).Dot("Return").Call(jen.ID("expectedUserCount"), jen.ID("nil")),
			jen.Line(),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
			jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.Qual("github.com/stretchr/testify/mock", "Anything")),
			jen.Var().ID("ucp").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider").Op("=").Func().Paramsln(
				jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName"),
				jen.ID("description").ID("string"),
			).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"), jen.ID("error")).SingleLineBlock(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.Line(),
			jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Callln(
				jen.Qual("context", "Background").Call(),
				jen.ID("config").Dot("AuthSettings").Values(),
				jen.ID("noop").Dot("ProvideNoopLogger").Call(),
				jen.ID("mockDB"),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(),
				jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
				jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
				jen.ID("ucp"), jen.ID("newsman").Dot("NewNewsman").Call(jen.ID("nil"), jen.ID("nil")),
			),
			jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Line(),
			jen.Return().ID("service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideUsersService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("mockUserCount").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUserCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("mockUserCount"), jen.ID("nil")),
				jen.Line(),
				jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("mockUserCount")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName"),
					jen.ID("description").ID("string")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
					jen.ID("error"),
				).Block(
					jen.Return().List(jen.ID("uc"), jen.ID("nil")),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Callln(
					jen.Qual("context", "Background").Call(),
					jen.ID("config").Dot("AuthSettings").Values(),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
					jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.ID("nil"),
				),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.ID("service")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil userIDFetcher"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("mockUserCount").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUserCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("mockUserCount"), jen.ID("nil")),
				jen.Line(),
				jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("mockUserCount")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName"),
					jen.ID("description").ID("string")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
					jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.ID("nil")),
				),
				jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Callln(
					jen.Qual("context", "Background").Call(), jen.ID("config").Dot(
						"AuthSettings",
					).Values(), jen.ID("noop").Dot(
						"ProvideNoopLogger",
					).Call(), jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(), jen.ID("nil"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(), jen.ID("ucp"), jen.ID("nil")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("service")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error initializing counter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("mockUserCount").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUserCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("mockUserCount"), jen.ID("nil")),
				jen.Line(),
				jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("mockUserCount")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName"),
					jen.ID("description").ID("string")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
					jen.ID("error"),
				).Block(
					jen.Return().List(jen.ID("uc"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Callln(
					jen.Qual("context", "Background").Call(),
					jen.ID("config").Dot("AuthSettings").Values(),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
					jen.ID("mockDB"),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(), jen.ID("ucp"), jen.ID("nil"),
				),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("service")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error getting user count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("mockUserCount").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUserCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("mockUserCount"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
				jen.Var().ID("ucp").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName"),
					jen.ID("description").ID("string")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
					jen.ID("error"),
				).Block(
					jen.Return().List(jen.ID("uc"), jen.ID("nil")),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Callln(
					jen.Qual("context", "Background").Call(),
					jen.ID("config").Dot("AuthSettings").Values(),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
					jen.ID("mockDB"),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(), jen.ID("ucp"), jen.ID("nil"),
				),
				jen.Line(),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("service")),
			)),
		),
		jen.Line(),
	)
	return ret
}