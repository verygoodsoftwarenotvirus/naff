package httpserver

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(pkg, ret)

	buildServerLines := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("DebugMode").Op(":").ID("true"),
			jen.ID("db").Op(":").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
			jen.ID("config").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "ServerConfig").Values(),
			jen.ID("encoder").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
			jen.ID("httpServer").Op(":").ID("provideHTTPServer").Call(),
			jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("frontendService").Op(":").Qual(filepath.Join(pkg.OutputPath, "services/v1/frontend"), "ProvideFrontendService").Callln(
				jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "FrontendSettings").Values(),
			),
			jen.ID("webhooksService").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "WebhookDataServer").Values(),
			jen.ID("usersService").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "UserDataServer").Values(),
			jen.ID("authService").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "services/v1/auth"), "Service").Values(),
		}
		for _, typ := range pkg.DataTypes {
			tpuvn := typ.Name.PluralUnexportedVarName()
			tsn := typ.Name.Singular()
			lines = append(lines,
				jen.IDf("%sService", tpuvn).Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), fmt.Sprintf("%sDataServer", tsn)).Values(),
			)
		}

		lines = append(lines,
			jen.ID("oauth2ClientsService").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "OAuth2ClientDataServer").Values(),
		)

		return lines
	}

	ret.Add(
		jen.Func().ID("buildTestServer").Params().Params(jen.Op("*").ID("Server")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Server").Valuesln(
				buildServerLines()...,
			),
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	buildProvideServerArgs := func() []jen.Code {
		args := []jen.Code{
			jen.Qual("context", "Background").Call(),
			jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "ServerConfig").Valuesln(
				jen.ID("Auth").Op(":").Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "AuthSettings").Valuesln(
					jen.ID("CookieSecret").Op(":").Lit("THISISAVERYLONGSTRINGFORTESTPURPOSES"),
				),
			),
			jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "services/v1/auth"), "Service").Values(),
			jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "services/v1/frontend"), "Service").Values(),
		}

		for _, typ := range pkg.DataTypes {
			pn := typ.Name.PackageName()
			args = append(args, jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "services/v1", pn), "Service").Values())
		}

		args = append(args,
			jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "services/v1/users"), "Service").Values(),
			jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "services/v1/oauth2clients"), "Service").Values(),
			jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "services/v1/webhooks"), "Service").Values(),
			jen.ID("mockDB"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
		)

		// if pkg.EnableNewsman {
		args = append(args,
			jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.ID("nil"), jen.ID("nil")),
		)
		// }

		return args
	}

	ret.Add(
		jen.Func().ID("TestProvideServer").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(jen.Lit("GetAllWebhooks"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Values(), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("ProvideServer").Callln(
					buildProvideServerArgs()...,
				),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
			)),
		),
		jen.Line(),
	)
	return ret
}
