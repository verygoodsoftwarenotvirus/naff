package v1

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("main")
	ret.HeaderComment("+build wireinject")

	utils.AddImports(proj, ret)

	newsmanImp := "gitlab.com/verygoodsoftwarenotvirus/newsman"
	loggingImp := "gitlab.com/verygoodsoftwarenotvirus/logging/v1"

	internalConfigImp := fmt.Sprintf("%s/internal/v1/config", proj.OutputPath)
	internalMetricsImp := fmt.Sprintf("%s/internal/v1/metrics", proj.OutputPath)
	internalEncodingImp := fmt.Sprintf("%s/internal/v1/encoding", proj.OutputPath)
	databaseClientImp := fmt.Sprintf("%s/database/v1", proj.OutputPath)
	internalAuthImp := fmt.Sprintf("%s/internal/v1/auth", proj.OutputPath)
	authServiceImp := fmt.Sprintf("%s/services/v1/auth", proj.OutputPath)
	usersServiceImp := fmt.Sprintf("%s/services/v1/users", proj.OutputPath)
	frontendServiceImp := fmt.Sprintf("%s/services/v1/frontend", proj.OutputPath)
	webhooksServiceImp := fmt.Sprintf("%s/services/v1/webhooks", proj.OutputPath)
	oauth2ClientsServiceImp := fmt.Sprintf("%s/services/v1/oauth2clients", proj.OutputPath)
	httpServerImp := fmt.Sprintf("%s/server/v1/http", proj.OutputPath)
	serverImp := fmt.Sprintf("%s/server/v1", proj.OutputPath)

	// if proj.EnableNewsman {
	ret.Add(
		jen.Comment("ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day."),
		jen.Line(),
		jen.Func().ID("ProvideReporter").Params(jen.ID("n").PointerTo().Qual(newsmanImp, "Newsman")).Params(jen.Qual(newsmanImp, "Reporter")).Block(
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
			jen.Qual(proj.InternalSearchV1Package("bleve"), "Providers"),
			jen.Qual(serverImp, "Providers"),
			jen.Qual(internalEncodingImp, "Providers"),
			jen.Qual(httpServerImp, "Providers"),
			jen.Comment("metrics"),
			jen.Qual(internalMetricsImp, "Providers"),
			jen.Comment("external libs"),
		}

		//if proj.EnableNewsman {
		args = append(args,
			jen.Qual(newsmanImp, "NewNewsman"),
		)
		//}

		args = append(args,
			jen.ID("ProvideReporter"),
			jen.Comment("services"),
			jen.Qual(authServiceImp, "Providers"),
			jen.Qual(usersServiceImp, "Providers"),
		)

		for _, typ := range proj.DataTypes {
			args = append(args,
				jen.Qual(fmt.Sprintf("%s/services/v1/%s", proj.OutputPath, typ.Name.PackageName()), "Providers"),
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
		jen.Comment("BuildServer builds a server."),
		jen.Line(),
		jen.Func().ID("BuildServer").Paramsln(
			constants.CtxParam(),
			jen.ID("cfg").PointerTo().Qual(internalConfigImp, "ServerConfig"),
			jen.ID(constants.LoggerVarName).Qual(loggingImp, "Logger"),
			jen.ID("database").Qual(databaseClientImp, "DataManager"),
			jen.ID("db").PointerTo().Qual("database/sql", "DB"),
		).Params(jen.PointerTo().Qual(serverImp, "Server"), jen.Error()).Block(
			jen.Qual("github.com/google/wire", "Build").Callln(
				buildWireBuildCallArgs()...,
			),
			jen.Return().List(jen.Nil(), jen.Nil()),
		),
		jen.Line(),
	)

	return ret
}
