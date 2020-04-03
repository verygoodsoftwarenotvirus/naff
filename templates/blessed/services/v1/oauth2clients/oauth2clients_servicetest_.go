package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsServiceTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(proj, ret)

	ret.Add(utils.FakeSeedFunc())

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Params(jen.PointerTo().ID("Service")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("manager").Assign().Qual("goproj.in/oauth2.v3/manage", "NewDefaultManager").Call(),
			jen.List(jen.ID("tokenStore"), jen.Err()).Assign().Qual("goproj.in/oauth2.v3/store", "NewMemoryTokenStore").Call(),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.ID("manager").Dot("MustTokenStorage").Call(jen.ID("tokenStore"), jen.Err()),
			jen.ID("server").Assign().Qual("goproj.in/oauth2.v3/server", "NewDefaultServer").Call(jen.ID("manager")),
			jen.Line(),
			jen.ID("service").Assign().VarPointer().ID("Service").Valuesln(
				jen.ID("database").MapAssign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.ID("encoderDecoder").MapAssign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("authenticator").MapAssign().VarPointer().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
				jen.ID("urlClientIDExtractor").MapAssign().Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
				jen.ID("oauth2ClientCounter").MapAssign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("tokenStore").MapAssign().ID("tokenStore"),
				jen.ID("oauth2Handler").MapAssign().ID("server"),
			),
			jen.Line(),
			jen.Return().ID("service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideOAuth2ClientsService").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2Clients"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(), jen.Nil()),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2ClientCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.ID("uc").Assign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expected")).Dot("Return").Call(), jen.Line(),
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Paramsln(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"), jen.Error()).Block(
					jen.Return().List(jen.ID("uc"), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideOAuth2ClientsService").Callln(
					utils.CtxVar(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.ID("mockDB"),
					jen.VarPointer().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertNotNil(jen.ID("service"), nil),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error providing counter"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2Clients"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(), jen.Nil()),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2ClientCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.ID("uc").Assign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expected")).Dot("Return").Call(), jen.Line(),
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Paramsln(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").ID("string")).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
					jen.Error(),
				).Block(
					jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideOAuth2ClientsService").Callln(
					utils.CtxVar(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("mockDB"),
					jen.VarPointer().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("service"), nil),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching oauth2 clients"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2Clients"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetAllOAuth2ClientCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("uc").Assign().VarPointer().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expected")).Dot("Return").Call(), jen.Line(),
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Paramsln(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").ID("string")).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
					jen.Error(),
				).Block(
					jen.Return().List(jen.ID("uc"), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideOAuth2ClientsService").Callln(
					utils.CtxVar(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("mockDB"),
					jen.VarPointer().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
				),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("service"), nil),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_clientStore_GetByID").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleID").Assign().Lit("blah"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleID"),
				).Dot("Return").Call(jen.VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(jen.ID("ClientID").MapAssign().ID("exampleID")), jen.Nil()),
				jen.Line(),
				jen.ID("c").Assign().VarPointer().ID("clientStore").Values(jen.ID("database").MapAssign().ID("mockDB")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetByID").Call(jen.ID("exampleID")),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleID"), jen.ID("actual").Dot("GetID").Call(), nil),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with no rows"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleID").Assign().Lit("blah"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.ID("c").Assign().VarPointer().ID("clientStore").Values(jen.ID("database").MapAssign().ID("mockDB")),
				jen.List(jen.ID("_"), jen.Err()).Assign().ID("c").Dot("GetByID").Call(jen.ID("exampleID")),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleID").Assign().Lit("blah"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.ID("exampleID"))),
				jen.Line(),
				jen.ID("c").Assign().VarPointer().ID("clientStore").Values(jen.ID("database").MapAssign().ID("mockDB")),
				jen.List(jen.ID("_"), jen.Err()).Assign().ID("c").Dot("GetByID").Call(jen.ID("exampleID")),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_HandleAuthorizeRequest").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("moah").Assign().VarPointer().ID("mockOauth2Handler").Values(),
				jen.ID("moah").Dot("On").Callln(
					jen.Lit("HandleAuthorizeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("moah"),
				jen.List(jen.ID("req"), jen.ID("res")).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
				jen.Line(),
				utils.AssertNoError(jen.ID("s").Dot("HandleAuthorizeRequest").Call(jen.ID("res"), jen.ID("req")), nil),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_HandleTokenRequest").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("moah").Assign().VarPointer().ID("mockOauth2Handler").Values(),
				jen.ID("moah").Dot("On").Callln(
					jen.Lit("HandleTokenRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("moah"),
				jen.List(jen.ID("req"), jen.ID("res")).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
				jen.Line(),
				utils.AssertNoError(jen.ID("s").Dot("HandleTokenRequest").Call(jen.ID("res"), jen.ID("req")), nil),
			)),
		),
		jen.Line(),
	)
	return ret
}
