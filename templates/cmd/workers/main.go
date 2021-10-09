package workers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("preWritesTopicName").Op("=").Lit("pre_writes"),
			jen.ID("dataChangesTopicName").Op("=").Lit("data_changes"),
			jen.ID("preUpdatesTopicName").Op("=").Lit("pre_updates"),
			jen.ID("preArchivesTopicName").Op("=").Lit("pre_archives"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("initializeLocalSecretManager").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("envVarKey").ID("string")).Params(jen.ID("secrets").Dot("SecretManager")).Body(
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("cfg").Op(":=").Op("&").ID("secrets").Dot("Config").Valuesln(jen.ID("Provider").Op(":").ID("secrets").Dot("ProviderLocal"), jen.ID("Key").Op(":").Qual("os", "Getenv").Call(jen.ID("envVarKey"))),
			jen.List(jen.ID("k"), jen.ID("err")).Op(":=").ID("secrets").Dot("ProvideSecretKeeper").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.List(jen.ID("sm"), jen.ID("err")).Op(":=").ID("secrets").Dot("ProvideSecretManager").Call(
				jen.ID("logger"),
				jen.ID("k"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.Return().ID("sm"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("configFilepathEnvVar").Op("=").Lit("CONFIGURATION_FILEPATH"),
			jen.ID("configStoreEnvVarKey").Op("=").Lit("TODO_WORKERS_LOCAL_CONFIG_STORE_KEY"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("main").Params().Body(
			jen.Var().Defs(
				jen.ID("addr").Op("=").Lit("worker_queue:6379"),
			),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("logger").Op(":=").ID("logging").Dot("ProvideLogger").Call(jen.ID("logging").Dot("Config").Valuesln(jen.ID("Provider").Op(":").ID("logging").Dot("ProviderZerolog"))),
			jen.ID("logger").Dot("Info").Call(jen.Lit("starting workers...")),
			jen.ID("configFilepath").Op(":=").Qual("os", "Getenv").Call(jen.ID("configFilepathEnvVar")),
			jen.If(jen.ID("configFilepath").Op("==").Lit("")).Body(
				jen.Qual("log", "Fatal").Call(jen.Lit("no config provided"))),
			jen.List(jen.ID("configBytes"), jen.ID("err")).Op(":=").Qual("os", "ReadFile").Call(jen.ID("configFilepath")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.ID("sm").Op(":=").ID("initializeLocalSecretManager").Call(
				jen.ID("ctx"),
				jen.ID("configStoreEnvVarKey"),
			),
			jen.Var().Defs(
				jen.ID("cfg").Op("*").ID("config").Dot("InstanceConfig"),
			),
			jen.If(jen.ID("err").Op("=").ID("sm").Dot("Decrypt").Call(
				jen.ID("ctx"),
				jen.ID("string").Call(jen.ID("configBytes")),
				jen.Op("&").ID("cfg"),
			), jen.ID("err").Op("!=").ID("nil").Op("||").ID("cfg").Op("==").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.ID("cfg").Dot("Observability").Dot("Tracing").Dot("Jaeger").Dot("ServiceName").Op("=").Lit("workers"),
			jen.List(jen.ID("flushFunc"), jen.ID("initializeTracerErr")).Op(":=").ID("cfg").Dot("Observability").Dot("Tracing").Dot("Initialize").Call(jen.ID("logger")),
			jen.If(jen.ID("initializeTracerErr").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("initializeTracerErr"),
					jen.Lit("initializing tracer"),
				)),
			jen.If(jen.ID("flushFunc").Op("!=").ID("nil")).Body(
				jen.Defer().ID("flushFunc").Call()),
			jen.ID("cfg").Dot("Database").Dot("RunMigrations").Op("=").ID("false"),
			jen.List(jen.ID("dataManager"), jen.ID("err")).Op(":=").ID("config").Dot("ProvideDatabaseClient").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("cfg"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.ID("pcfg").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/config", "Config").Valuesln(jen.ID("Provider").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/config", "ProviderRedis"), jen.ID("RedisConfig").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/config", "RedisConfig").Valuesln(jen.ID("QueueAddress").Op(":").ID("addr"))),
			jen.ID("consumerProvider").Op(":=").ID("consumers").Dot("ProvideRedisConsumerProvider").Call(
				jen.ID("logger"),
				jen.ID("addr"),
			),
			jen.List(jen.ID("publisherProvider"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/config", "ProvidePublisherProvider").Call(
				jen.ID("logger"),
				jen.ID("pcfg"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.ID("postWritesWorker").Op(":=").ID("workers").Dot("ProvideDataChangesWorker").Call(jen.ID("logger")),
			jen.List(jen.ID("postWritesConsumer"), jen.ID("err")).Op(":=").ID("consumerProvider").Dot("ProviderConsumer").Call(
				jen.ID("ctx"),
				jen.ID("dataChangesTopicName"),
				jen.ID("postWritesWorker").Dot("HandleMessage"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.Go().ID("postWritesConsumer").Dot("Consume").Call(
				jen.ID("nil"),
				jen.ID("nil"),
			),
			jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Valuesln(jen.ID("Timeout").Op(":").Lit(5).Op("*").Qual("time", "Second")),
			jen.List(jen.ID("postWritesPublisher"), jen.ID("err")).Op(":=").ID("publisherProvider").Dot("ProviderPublisher").Call(jen.ID("dataChangesTopicName")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.List(jen.ID("preWritesWorker"), jen.ID("err")).Op(":=").ID("workers").Dot("ProvidePreWritesWorker").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("client"),
				jen.ID("dataManager"),
				jen.ID("postWritesPublisher"),
				jen.Lit("http://elasticsearch:9200"),
				jen.ID("elasticsearch").Dot("NewIndexManager"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.List(jen.ID("preWritesConsumer"), jen.ID("err")).Op(":=").ID("consumerProvider").Dot("ProviderConsumer").Call(
				jen.ID("ctx"),
				jen.ID("preWritesTopicName"),
				jen.ID("preWritesWorker").Dot("HandleMessage"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.Go().ID("preWritesConsumer").Dot("Consume").Call(
				jen.ID("nil"),
				jen.ID("nil"),
			),
			jen.List(jen.ID("postUpdatesPublisher"), jen.ID("err")).Op(":=").ID("publisherProvider").Dot("ProviderPublisher").Call(jen.ID("dataChangesTopicName")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.List(jen.ID("preUpdatesWorker"), jen.ID("err")).Op(":=").ID("workers").Dot("ProvidePreUpdatesWorker").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("client"),
				jen.ID("dataManager"),
				jen.ID("postUpdatesPublisher"),
				jen.Lit("http://elasticsearch:9200"),
				jen.ID("elasticsearch").Dot("NewIndexManager"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.List(jen.ID("preUpdatesConsumer"), jen.ID("err")).Op(":=").ID("consumerProvider").Dot("ProviderConsumer").Call(
				jen.ID("ctx"),
				jen.ID("preUpdatesTopicName"),
				jen.ID("preUpdatesWorker").Dot("HandleMessage"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.Go().ID("preUpdatesConsumer").Dot("Consume").Call(
				jen.ID("nil"),
				jen.ID("nil"),
			),
			jen.List(jen.ID("postArchivesPublisher"), jen.ID("err")).Op(":=").ID("publisherProvider").Dot("ProviderPublisher").Call(jen.ID("dataChangesTopicName")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.List(jen.ID("preArchivesWorker"), jen.ID("err")).Op(":=").ID("workers").Dot("ProvidePreArchivesWorker").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("client"),
				jen.ID("dataManager"),
				jen.ID("postArchivesPublisher"),
				jen.Lit("http://elasticsearch:9200"),
				jen.ID("elasticsearch").Dot("NewIndexManager"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.List(jen.ID("preArchivesConsumer"), jen.ID("err")).Op(":=").ID("consumerProvider").Dot("ProviderConsumer").Call(
				jen.ID("ctx"),
				jen.ID("preArchivesTopicName"),
				jen.ID("preArchivesWorker").Dot("HandleMessage"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.Go().ID("preArchivesConsumer").Dot("Consume").Call(
				jen.ID("nil"),
				jen.ID("nil"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("working...")),
			jen.ID("sigChan").Op(":=").ID("make").Call(
				jen.Chan().Qual("os", "Signal"),
				jen.Lit(1),
			),
			jen.Qual("os/signal", "Notify").Call(
				jen.ID("sigChan"),
				jen.Qual("syscall", "SIGINT"),
				jen.Qual("syscall", "SIGTERM"),
			),
			jen.Op("<-").ID("sigChan"),
		),
		jen.Newline(),
	)

	return code
}
