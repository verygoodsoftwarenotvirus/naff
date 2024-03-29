package workers

import (
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("preWritesTopicName").Equals().Lit("pre_writes"),
			jen.ID("dataChangesTopicName").Equals().Lit("data_changes"),
			jen.ID("preUpdatesTopicName").Equals().Lit("pre_updates"),
			jen.ID("preArchivesTopicName").Equals().Lit("pre_archives"),
			jen.Newline(),
			jen.ID("configFilepathEnvVar").Equals().Lit("CONFIGURATION_FILEPATH"),
			jen.ID("configStoreEnvVarKey").Equals().Litf("%s_WORKERS_LOCAL_CONFIG_STORE_KEY", strings.ToUpper(proj.Name.Singular())),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("initializeLocalSecretManager").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("envVarKey").String()).Params(jen.Qual(proj.InternalSecretsPackage(), "SecretManager")).Body(
			jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
			jen.Newline(),
			jen.ID("cfg").Assign().Op("&").Qual(proj.InternalSecretsPackage(), "Config").Valuesln(
				jen.ID("Provider").MapAssign().Qual(proj.InternalSecretsPackage(), "ProviderLocal"),
				jen.ID("Key").MapAssign().Qual("os", "Getenv").Call(jen.ID("envVarKey")),
			),
			jen.Newline(),
			jen.List(jen.ID("k"), jen.ID("err")).Assign().Qual(proj.InternalSecretsPackage(), "ProvideSecretKeeper").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.List(jen.ID("sm"), jen.ID("err")).Assign().Qual(proj.InternalSecretsPackage(), "ProvideSecretManager").Call(
				jen.ID("logger"),
				jen.ID("k"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.Return().ID("sm"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("main").Params().Body(
			jen.ID("ctx").Assign().Qual("context", "Background").Call(),
			jen.Newline(),
			jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "ProvideLogger").Call(jen.Qual(proj.InternalLoggingPackage(), "Config").Valuesln(
				jen.ID("Provider").MapAssign().Qual(proj.InternalLoggingPackage(), "ProviderZerolog"))),
			jen.Newline(),
			jen.ID("logger").Dot("Info").Call(jen.Lit("starting workers...")),
			jen.Newline(),
			jen.Comment("find and validate our configuration filepath."),
			jen.ID("configFilepath").Assign().Qual("os", "Getenv").Call(jen.ID("configFilepathEnvVar")),
			jen.If(jen.ID("configFilepath").IsEqualTo().Lit("")).Body(
				jen.Qual("log", "Fatal").Call(jen.Lit("no config provided")),
			),
			jen.Newline(),
			jen.List(jen.ID("configBytes"), jen.ID("err")).Assign().Qual("os", "ReadFile").Call(jen.ID("configFilepath")),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.ID("sm").Assign().ID("initializeLocalSecretManager").Call(
				jen.ID("ctx"),
				jen.ID("configStoreEnvVarKey"),
			),
			jen.Newline(),
			jen.Var().ID("cfg").Op("*").Qual(proj.InternalConfigPackage(), "InstanceConfig"),
			jen.If(jen.ID("err").Equals().ID("sm").Dot("Decrypt").Call(jen.ID("ctx"), jen.String().Call(jen.ID("configBytes")), jen.Op("&").ID("cfg")), jen.ID("err").DoesNotEqual().Nil().Op("||").ID("cfg").IsEqualTo().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.ID("cfg").Dot("Observability").Dot("Tracing").Dot("Jaeger").Dot("ServiceName").Equals().Lit("workers"),
			jen.Newline(),
			jen.List(jen.ID("flushFunc"), jen.ID("initializeTracerErr")).Assign().ID("cfg").Dot("Observability").Dot("Tracing").Dot("Initialize").Call(jen.ID("logger")),
			jen.If(jen.ID("initializeTracerErr").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("initializeTracerErr"),
					jen.Lit("initializing tracer"),
				),
			),
			jen.Newline(),
			jen.Comment("if tracing is disabled, this will be nil"),
			jen.If(jen.ID("flushFunc").DoesNotEqual().Nil()).Body(
				jen.Defer().ID("flushFunc").Call(),
			),
			jen.Newline(),
			jen.ID("cfg").Dot("Database").Dot("RunMigrations").Equals().False(),
			jen.Newline(),
			jen.List(jen.ID("dataManager"), jen.ID("err")).Assign().Qual(proj.InternalConfigPackage(), "ProvideDatabaseClient").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("cfg"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.Newline(),
			jen.ID("consumerProvider").Assign().Qual(proj.InternalMessageQueueConsumersPackage(), "ProvideRedisConsumerProvider").Call(
				jen.ID("logger"),
				jen.String().Call(jen.ID("cfg").Dot("Events").Dot("RedisConfig").Dot("QueueAddress")),
			),
			jen.Newline(),
			jen.List(jen.ID("publisherProvider"), jen.ID("err")).Assign().Qual(proj.InternalMessageQueueConfigPackage(), "ProvidePublisherProvider").Call(
				jen.ID("logger"),
				jen.AddressOf().ID("cfg").Dot("Events"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.Comment("post-writes worker"),
			jen.Newline(),
			jen.ID("postWritesWorker").Assign().Qual(proj.InternalWorkersPackage(), "ProvideDataChangesWorker").Call(jen.ID("logger")),
			jen.List(jen.ID("postWritesConsumer"), jen.ID("err")).Assign().ID("consumerProvider").Dot("ProviderConsumer").Call(
				jen.ID("ctx"),
				jen.ID("dataChangesTopicName"),
				jen.ID("postWritesWorker").Dot("HandleMessage"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.Go().ID("postWritesConsumer").Dot("Consume").Call(
				jen.Nil(),
				jen.Nil(),
			),
			jen.Newline(),
			jen.Comment("pre-writes worker"),
			jen.Newline(),
			jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Valuesln(
				jen.ID("Timeout").MapAssign().Lit(5).Op("*").Qual("time", "Second")),
			jen.Newline(),
			jen.List(jen.ID("postWritesPublisher"), jen.ID("err")).Assign().ID("publisherProvider").Dot("ProviderPublisher").Call(jen.ID("dataChangesTopicName")),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.List(jen.ID("preWritesWorker"), jen.ID("err")).Assign().Qual(proj.InternalWorkersPackage(), "ProvidePreWritesWorker").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("client"),
				jen.ID("dataManager"),
				jen.ID("postWritesPublisher"),
				jen.ID("cfg").Dot("Search").Dot("Address"),
				jen.Qual(proj.InternalSearchPackage("elasticsearch"), "NewIndexManager"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.Newline(),
			jen.List(jen.ID("preWritesConsumer"), jen.ID("err")).Assign().ID("consumerProvider").Dot("ProviderConsumer").Call(
				jen.ID("ctx"),
				jen.ID("preWritesTopicName"),
				jen.ID("preWritesWorker").Dot("HandleMessage"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.Go().ID("preWritesConsumer").Dot("Consume").Call(
				jen.Nil(),
				jen.Nil(),
			),
			jen.Newline(),
			jen.Comment("pre-updates worker"),
			jen.Newline(),
			jen.List(jen.ID("postUpdatesPublisher"), jen.ID("err")).Assign().ID("publisherProvider").Dot("ProviderPublisher").Call(jen.ID("dataChangesTopicName")),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.Newline(),
			jen.List(jen.ID("preUpdatesWorker"), jen.ID("err")).Assign().Qual(proj.InternalWorkersPackage(), "ProvidePreUpdatesWorker").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("client"),
				jen.ID("dataManager"),
				jen.ID("postUpdatesPublisher"),
				jen.ID("cfg").Dot("Search").Dot("Address"),
				jen.Qual(proj.InternalSearchPackage("elasticsearch"), "NewIndexManager"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.Newline(),
			jen.List(jen.ID("preUpdatesConsumer"), jen.ID("err")).Assign().ID("consumerProvider").Dot("ProviderConsumer").Call(
				jen.ID("ctx"),
				jen.ID("preUpdatesTopicName"),
				jen.ID("preUpdatesWorker").Dot("HandleMessage"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.Go().ID("preUpdatesConsumer").Dot("Consume").Call(
				jen.Nil(),
				jen.Nil(),
			),
			jen.Newline(),
			jen.Comment("pre-archives worker"),
			jen.Newline(),
			jen.List(jen.ID("postArchivesPublisher"), jen.ID("err")).Assign().ID("publisherProvider").Dot("ProviderPublisher").Call(jen.ID("dataChangesTopicName")),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.List(jen.ID("preArchivesWorker"), jen.ID("err")).Assign().Qual(proj.InternalWorkersPackage(), "ProvidePreArchivesWorker").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("client"),
				jen.ID("dataManager"),
				jen.ID("postArchivesPublisher"),
				jen.ID("cfg").Dot("Search").Dot("Address"),
				jen.Qual(proj.InternalSearchPackage("elasticsearch"), "NewIndexManager"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err"))),
			jen.Newline(),
			jen.List(jen.ID("preArchivesConsumer"), jen.ID("err")).Assign().ID("consumerProvider").Dot("ProviderConsumer").Call(
				jen.ID("ctx"),
				jen.ID("preArchivesTopicName"),
				jen.ID("preArchivesWorker").Dot("HandleMessage"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.Go().ID("preArchivesConsumer").Dot("Consume").Call(
				jen.Nil(),
				jen.Nil(),
			),
			jen.Newline(),
			jen.ID("logger").Dot("Info").Call(jen.Lit("working...")),
			jen.Newline(),
			jen.Comment("wait for signal to exit"),
			jen.ID("sigChan").Assign().ID("make").Call(
				jen.Chan().Qual("os", "Signal"),
				jen.One(),
			),
			jen.Qual("os/signal", "Notify").Call(
				jen.ID("sigChan"),
				jen.Qual("syscall", "SIGINT"),
				jen.Qual("syscall", "SIGTERM"),
			),
			jen.Newline(),
			jen.Op("<-").ID("sigChan"),
		),
		jen.Newline(),
	)

	return code
}
