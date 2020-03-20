package oauth2clients

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsServiceTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Service")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("manager").Op(":=").Qual("gopkg.in/oauth2.v3/manage", "NewDefaultManager").Call(),
			jen.List(jen.ID("tokenStore"), jen.Err()).Op(":=").Qual("gopkg.in/oauth2.v3/store", "NewMemoryTokenStore").Call(),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.ID("manager").Dot("MustTokenStorage").Call(jen.ID("tokenStore"), jen.Err()),
			jen.ID("server").Op(":=").Qual("gopkg.in/oauth2.v3/server", "NewDefaultServer").Call(jen.ID("manager")),
			jen.Line(),
			jen.ID("service").Op(":=").Op("&").ID("Service").Valuesln(
				jen.ID("database").Op(":").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.ID("encoderDecoder").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("authenticator").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth/mock"), "Authenticator").Values(),
				jen.ID("urlClientIDExtractor").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
				jen.ID("oauth2ClientCounter").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("tokenStore").Op(":").ID("tokenStore"),
				jen.ID("oauth2Handler").Op(":").ID("server"),
			),
			jen.Line(),
			jen.Return().ID("service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideOAuth2ClientsService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2Clients"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Values(), jen.Nil()),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2ClientCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.ID("uc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expected")).Dot("Return").Call(), jen.Line(),
				jen.Var().ID("ucp").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"), jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Op(":=").ID("ProvideOAuth2ClientsService").Callln(
					utils.CtxVar(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.ID("mockDB"),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth/mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("service")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error providing counter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2Clients"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Values(), jen.Nil()),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2ClientCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.ID("uc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expected")).Dot("Return").Call(), jen.Line(),
				jen.Var().ID("ucp").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName"),
					jen.ID("description").ID("string")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
					jen.ID("error"),
				).Block(
					jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Op(":=").ID("ProvideOAuth2ClientsService").Callln(
					utils.CtxVar(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("mockDB"),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth/mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("service")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching oauth2 clients"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2Clients"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Values(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2ClientCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("uc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expected")).Dot("Return").Call(), jen.Line(),
				jen.Var().ID("ucp").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName"),
					jen.ID("description").ID("string")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
					jen.ID("error"),
				).Block(
					jen.Return().List(jen.ID("uc"), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Op(":=").ID("ProvideOAuth2ClientsService").Callln(
					utils.CtxVar(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("mockDB"),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth/mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
				),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("service")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_clientStore_GetByID").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleID").Op(":=").Lit("blah"),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleID"),
				).Dot("Return").Call(jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Values(jen.ID("ClientID").Op(":").ID("exampleID")), jen.Nil()),
				jen.Line(),
				jen.ID("c").Op(":=").Op("&").ID("clientStore").Values(jen.ID("database").Op(":").ID("mockDB")),
				jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("c").Dot("GetByID").Call(jen.ID("exampleID")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("exampleID"), jen.ID("actual").Dot("GetID").Call()),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with no rows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleID").Op(":=").Lit("blah"),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleID"),
				).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.ID("c").Op(":=").Op("&").ID("clientStore").Values(jen.ID("database").Op(":").ID("mockDB")),
				jen.List(jen.ID("_"), jen.Err()).Op(":=").ID("c").Dot("GetByID").Call(jen.ID("exampleID")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleID").Op(":=").Lit("blah"),
				jen.Line(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleID"),
				).Dot("Return").Call(jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.ID("exampleID"))),
				jen.Line(),
				jen.ID("c").Op(":=").Op("&").ID("clientStore").Values(jen.ID("database").Op(":").ID("mockDB")),
				jen.List(jen.ID("_"), jen.Err()).Op(":=").ID("c").Dot("GetByID").Call(jen.ID("exampleID")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_HandleAuthorizeRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("moah").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
				jen.ID("moah").Dot("On").Callln(
					jen.Lit("HandleAuthorizeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Op("=").ID("moah"),
				jen.List(jen.ID("req"), jen.ID("res")).Op(":=").List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("s").Dot("HandleAuthorizeRequest").Call(jen.ID("res"), jen.ID("req"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_HandleTokenRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("moah").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
				jen.ID("moah").Dot("On").Callln(
					jen.Lit("HandleTokenRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Op("=").ID("moah"),
				jen.List(jen.ID("req"), jen.ID("res")).Op(":=").List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("s").Dot("HandleTokenRequest").Call(jen.ID("res"), jen.ID("req"))),
			)),
		),
		jen.Line(),
	)
	return ret
}
