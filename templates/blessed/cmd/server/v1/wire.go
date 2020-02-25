package v1

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("main")
	ret.HeaderComment("+build wireinject")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	newsmanImp := "gitlab.com/verygoodsoftwarenotvirus/newsman"
	loggingImp := "gitlab.com/verygoodsoftwarenotvirus/logging/v1"

	internalConfigImp := fmt.Sprintf("%s/internal/v1/config", pkg.OutputPath)
	internalMetricsImp := fmt.Sprintf("%s/internal/v1/metrics", pkg.OutputPath)
	internalEncodingImp := fmt.Sprintf("%s/internal/v1/encoding", pkg.OutputPath)
	databaseClientImp := fmt.Sprintf("%s/database/v1", pkg.OutputPath)
	internalAuthImp := fmt.Sprintf("%s/internal/v1/auth", pkg.OutputPath)
	authServiceImp := fmt.Sprintf("%s/services/v1/auth", pkg.OutputPath)
	usersServiceImp := fmt.Sprintf("%s/services/v1/users", pkg.OutputPath)
	frontendServiceImp := fmt.Sprintf("%s/services/v1/frontend", pkg.OutputPath)
	webhooksServiceImp := fmt.Sprintf("%s/services/v1/webhooks", pkg.OutputPath)
	oauth2ClientsServiceImp := fmt.Sprintf("%s/services/v1/oauth2clients", pkg.OutputPath)
	httpServerImp := fmt.Sprintf("%s/server/v1/http", pkg.OutputPath)
	serverImp := fmt.Sprintf("%s/server/v1", pkg.OutputPath)

	// if pkg.EnableNewsman {
	ret.Add(
		jen.Comment("ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day"),
		jen.Line(),
		jen.Func().ID("ProvideReporter").Params(jen.ID("n").Op("*").Qual(newsmanImp, "Newsman")).Params(jen.Qual(newsmanImp, "Reporter")).Block(
			jen.Return().ID("n"),
		),
		jen.Line(),
	)
	// }

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
		}

		// if pkg.EnableNewsman {
		args = append(args,
			jen.Qual(newsmanImp, "NewNewsman"),
		)
		// }

		args = append(args,
			jen.ID("ProvideReporter"),
			jen.Comment("services"),
			jen.Qual(authServiceImp, "Providers"),
			jen.Qual(usersServiceImp, "Providers"),
		)

		for _, typ := range pkg.DataTypes {
			args = append(args,
				jen.Qual(fmt.Sprintf("%s/services/v1/%s", pkg.OutputPath, typ.Name.PackageName()), "Providers"),
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
			utils.CtxParam(),
			jen.ID("cfg").Op("*").Qual(internalConfigImp, "ServerConfig"),
			jen.ID("logger").Qual(loggingImp, "Logger"),
			jen.ID("database").Qual(databaseClientImp, "Database")).
			Params(jen.Op("*").Qual(serverImp, "Server"), jen.ID("error")).Block(
			jen.Qual("github.com/google/wire", "Build").Callln(
				buildWireBuildCallArgs()...,
			),
			jen.Return().List(jen.ID("nil"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return ret
}
