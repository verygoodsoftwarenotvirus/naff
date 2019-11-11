package v1

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("main")
	ret.HeaderComment("+build wireinject")

	utils.AddImports(pkgRoot, types, ret)

	newsmanImp := "gitlab.com/verygoodsoftwarenotvirus/newsman"
	loggingImp := "gitlab.com/verygoodsoftwarenotvirus/logging/v1"

	internalConfigImp := fmt.Sprintf("%s/internal/v1/config", pkgRoot)
	internalMetricsImp := fmt.Sprintf("%s/internal/v1/metrics", pkgRoot)
	internalEncodingImp := fmt.Sprintf("%s/internal/v1/encoding", pkgRoot)
	databaseClientImp := fmt.Sprintf("%s/database/v1", pkgRoot)
	internalAuthImp := fmt.Sprintf("%s/internal/v1/auth", pkgRoot)
	authServiceImp := fmt.Sprintf("%s/services/v1/auth", pkgRoot)
	usersServiceImp := fmt.Sprintf("%s/services/v1/users", pkgRoot)
	frontendServiceImp := fmt.Sprintf("%s/services/v1/frontend", pkgRoot)
	webhooksServiceImp := fmt.Sprintf("%s/services/v1/webhooks", pkgRoot)
	oauth2ClientsServiceImp := fmt.Sprintf("%s/services/v1/oauth2clients", pkgRoot)
	httpServerImp := fmt.Sprintf("%s/server/v1/http", pkgRoot)
	serverImp := fmt.Sprintf("%s/server/v1", pkgRoot)

	ret.Add(
		jen.Comment("ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day"),
		jen.Line(),
		jen.Func().ID("ProvideReporter").Params(jen.ID("n").Op("*").Qual(newsmanImp, "Newsman")).Params(jen.Qual(newsmanImp, "Reporter")).Block(
			jen.Return().ID("n"),
		),
		jen.Line(),
	)

	buildWireBuildCallArgs := func() []jen.Code {
		args := []jen.Code{
			jen.Qual(internalConfigImp, "Providers"),
			jen.Qual(internalAuthImp, "Providers"),
			jen.Comment("server things"),
			jen.Qual(serverImp, "Providers"),
			jen.Qual(internalEncodingImp, "Providers"),
			jen.Qual(httpServerImp, "Providers"),
			jen.Comment("metrics"),
			jen.Qual(internalMetricsImp, "Providers"),
			jen.Comment("external libs"),
			jen.Qual(newsmanImp, "NewNewsman"),
			jen.ID("ProvideReporter"),
			jen.Comment("services"),
			jen.Qual(authServiceImp, "Providers"),
			jen.Qual(usersServiceImp, "Providers"),
		}

		for _, typ := range types {
			args = append(args,
				jen.Qual(fmt.Sprintf("%s/services/v1/%s", pkgRoot, typ.Name.PluralRouteName()), "Providers"),
			)
		}

		args = append(args,
			jen.Qual(frontendServiceImp, "Providers"),
			jen.Qual(webhooksServiceImp, "Providers"),
			jen.Qual(oauth2ClientsServiceImp, "Providers"),
		)

		return args
	}

	ret.Add(
		jen.Comment("BuildServer builds a server"),
		jen.Line(),
		jen.Func().ID("BuildServer").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("cfg").Op("*").Qual(internalConfigImp, "ServerConfig"),
			jen.ID("logger").Qual(loggingImp, "Logger"),
			jen.ID("database").Qual(databaseClientImp, "Database")).
			Params(jen.Op("*").Qual(serverImp, "Server"), jen.ID("error")).Block(
			jen.ID("wire").Dot("Build").Callln(
				buildWireBuildCallArgs()...,
			),
			jen.Return().List(jen.ID("nil"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return ret
}
