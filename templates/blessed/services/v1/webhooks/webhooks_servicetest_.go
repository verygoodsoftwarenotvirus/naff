package webhooks

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksServiceTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params().Params(jen.Op("*").ID("Service")).Block(
			jen.Return().Op("&").ID("Service").Valuesln(
				jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.ID("webhookCounter").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("webhookDatabase").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "WebhookDataManager").Values(),
				jen.ID("userIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
				jen.ID("webhookIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
				jen.ID("encoderDecoder").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
				jen.ID("eventManager").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.ID("nil"), jen.ID("nil")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideWebhooksService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("uc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"), jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.ID("nil")),
				),
				jen.Line(),
				jen.ID("dm").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "WebhookDataManager").Values(),
				jen.ID("dm").Dot("On").Call(jen.Lit("GetAllWebhooksCount"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("ProvideWebhooksService").Callln(
					jen.Qual("context", "Background").Call(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("dm"),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.ID("nil"), jen.ID("nil")),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error providing counter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Var().ID("ucp").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName"),
					jen.ID("description").ID("string")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
					jen.ID("error")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("ProvideWebhooksService").Callln(
					jen.Qual("context", "Background").Call(), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "WebhookDataManager").Values(),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.ID("nil"), jen.ID("nil"))),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error setting count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("uc").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics/mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider").Op("=").Func().Paramsln(
					jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName"),
					jen.ID("description").ID("string"),
				).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"), jen.ID("error")).Block(
					jen.Return().List(jen.ID("uc"), jen.ID("nil")),
				),
				jen.Line(),
				jen.ID("dm").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "WebhookDataManager").Values(),
				jen.ID("dm").Dot("On").Call(jen.Lit("GetAllWebhooksCount"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("ProvideWebhooksService").Callln(
					jen.Qual("context", "Background").Call(), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("dm"),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).SingleLineBlock(jen.Return().Lit(0)),
					jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding/mock"), "EncoderDecoder").Values(), jen.ID("ucp"),
					jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.ID("nil"), jen.ID("nil")),
				),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
			)),
		),
		jen.Line(),
	)
	return ret
}
