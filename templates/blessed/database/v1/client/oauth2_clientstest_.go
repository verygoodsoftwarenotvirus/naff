package client

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsTestDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Func().ID("TestClient_GetOAuth2Client").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleClientID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error returned from querier"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleClientID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Parens(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client")).Call(jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetOAuth2ClientByClientID").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleClientID").Op(":=").Lit("CLIENT_ID"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client").Values(),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2ClientByClientID").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error returned from querier"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleClientID").Op(":=").Lit("CLIENT_ID"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Parens(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client")).Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID")).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2ClientByClientID").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetOAuth2ClientCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2ClientCount"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2ClientCount").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2ClientCount"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual("github.com/stretchr/testify/mock", "AnythingOfType").Call(jen.Lit("*models.QueryFilter")), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2ClientCount").Call(jen.Qual("context", "Background").Call(), jen.ID("nil"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error returned from querier"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2ClientCount"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2ClientCount").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetAllOAuth2ClientCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetAllOAuth2ClientCount"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAllOAuth2ClientCount").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetAllOAuth2Clients").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.Var().ID("expected").Index().Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetAllOAuth2Clients"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAllOAuth2Clients").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetOAuth2Clients").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientList").Values(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2Clients"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2Clients").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientList").Values(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2Clients"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Parens(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter")).Call(jen.ID("nil")), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2Clients").Call(jen.Qual("context", "Background").Call(), jen.ID("nil"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error returned from querier"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Parens(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientList")).Call(jen.ID("nil")),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetOAuth2Clients"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2Clients").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_CreateOAuth2Client").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client").Values(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientCreationInput").Values(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("CreateOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error returned from querier"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Parens(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client")).Call(jen.ID("nil")),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientCreationInput").Values(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("CreateOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.ID("errors").Dot("New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_UpdateOAuth2Client").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("example").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client").Values(),
				jen.Var().ID("expected").ID("error"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("UpdateOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("example")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("c").Dot("UpdateOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("example")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_ArchiveOAuth2Client").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleClientID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Var().ID("expected").ID("error"),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("ArchiveOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("c").Dot("ArchiveOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error returned from querier"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleClientID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Qual("fmt", "Errorf").Call(jen.Lit("blah")),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("ArchiveOAuth2Client"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleClientID"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("c").Dot("ArchiveOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)
	return ret
}
