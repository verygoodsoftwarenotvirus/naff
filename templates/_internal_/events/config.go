package events

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ProviderGoogleCloudPubSub").Op("=").Lit("google_cloud_pubsub"),
			jen.ID("ProviderAWSSQS").Op("=").Lit("aws_sqs"),
			jen.ID("ProviderRabbitMQ").Op("=").Lit("rabbit_mq"),
			jen.ID("ProviderAzureServiceBus").Op("=").Lit("azure_service_bus"),
			jen.ID("ProviderKafka").Op("=").Lit("kafka"),
			jen.ID("ProviderNATS").Op("=").Lit("nats"),
			jen.ID("ProviderMemory").Op("=").Lit("memory"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errNilConfig").Op("=").Qual("errors", "New").Call(jen.Lit("nil config provided")),
			jen.ID("errNilSubscription").Op("=").Qual("errors", "New").Call(jen.Lit("nil subscription provided")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Config").Struct(
				jen.ID("Provider").ID("string"),
				jen.ID("Topic").ID("string"),
				jen.ID("SubscriptionIdentifier").ID("string"),
				jen.ID("ConnectionURL").ID("string"),
				jen.ID("GCPPubSub").ID("GCPPubSub"),
				jen.ID("AckDeadline").Qual("time", "Duration"),
				jen.ID("Enabled").ID("bool"),
			),
			jen.ID("GCPPubSub").Struct(jen.ID("ServiceAccountKeyFilepath").ID("string")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("Config")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates the Config struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("c"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("Provider"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						jen.ID("ProviderGoogleCloudPubSub"),
						jen.ID("ProviderAWSSQS"),
						jen.ID("ProviderRabbitMQ"),
						jen.ID("ProviderAzureServiceBus"),
						jen.ID("ProviderKafka"),
						jen.ID("ProviderNATS"),
						jen.ID("ProviderMemory"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errInvalidProvider").Op("=").Qual("errors", "New").Call(jen.Lit("invalid events subscription provider")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvidePublishTopic uses a configuration to provide a pubsub subscription."),
		jen.Line(),
		jen.Func().ID("ProvidePublishTopic").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("Config")).Params(jen.Op("*").Qual("gocloud.dev/pubsub", "Topic"), jen.ID("error")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				jen.Case(jen.ID("ProviderGoogleCloudPubSub")).Body(
					jen.Var().Defs(
						jen.ID("creds").Op("*").ID("google").Dot("Credentials"),
					), jen.If(jen.ID("cfg").Dot("GCPPubSub").Dot("ServiceAccountKeyFilepath").Op("!=").Lit("")).Body(
						jen.List(jen.ID("serviceAccountKeyBytes"), jen.ID("err")).Op(":=").Qual("os", "ReadFile").Call(jen.ID("cfg").Dot("GCPPubSub").Dot("ServiceAccountKeyFilepath")),
						jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
							jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
								jen.Lit("reading service account key file: %w"),
								jen.ID("err"),
							))),
						jen.If(jen.List(jen.ID("creds"), jen.ID("err")).Op("=").ID("google").Dot("CredentialsFromJSON").Call(
							jen.ID("ctx"),
							jen.ID("serviceAccountKeyBytes"),
						), jen.ID("err").Op("!=").ID("nil")).Body(
							jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
								jen.Lit("using service account key credentials: %w"),
								jen.ID("err"),
							))),
					).Else().Body(
						jen.Var().Defs(
							jen.ID("err").ID("error"),
						),
						jen.If(jen.List(jen.ID("creds"), jen.ID("err")).Op("=").Qual("gocloud.dev/gcp", "DefaultCredentials").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
							jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
								jen.Lit("constructing pub/sub credentials: %w"),
								jen.ID("err"),
							))),
					), jen.List(jen.ID("conn"), jen.ID("_"), jen.ID("err")).Op(":=").ID("gcppubsub").Dot("Dial").Call(
						jen.ID("ctx"),
						jen.ID("creds").Dot("TokenSource"),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("dialing connection to pub/sub %w"),
							jen.ID("err"),
						))), jen.List(jen.ID("pubClient"), jen.ID("err")).Op(":=").ID("gcppubsub").Dot("PublisherClient").Call(
						jen.ID("ctx"),
						jen.ID("conn"),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("establishing publisher client: %w"),
							jen.ID("err"),
						))), jen.Return().ID("gcppubsub").Dot("OpenTopicByPath").Call(
						jen.ID("pubClient"),
						jen.ID("cfg").Dot("Topic"),
						jen.ID("nil"),
					)),
				jen.Case(jen.ID("ProviderAWSSQS")).Body(
					jen.List(jen.ID("sess"), jen.ID("err")).Op(":=").ID("session").Dot("NewSession").Call(jen.ID("nil")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("establishing AWS session: %w"),
							jen.ID("err"),
						))), jen.ID("topic").Op(":=").ID("awssnssqs").Dot("OpenSQSTopic").Call(
						jen.ID("ctx"),
						jen.ID("sess"),
						jen.ID("cfg").Dot("Topic"),
						jen.ID("nil"),
					), jen.Return().List(jen.ID("topic"), jen.ID("nil"))),
				jen.Case(jen.ID("ProviderKafka")).Body(
					jen.ID("config").Op(":=").ID("kafkapubsub").Dot("MinimalConfig").Call(), jen.Return().ID("kafkapubsub").Dot("OpenTopic").Call(
						jen.Qual("strings", "Split").Call(
							jen.ID("cfg").Dot("ConnectionURL"),
							jen.Lit(","),
						),
						jen.ID("config"),
						jen.ID("cfg").Dot("Topic"),
						jen.ID("nil"),
					)),
				jen.Case(jen.ID("ProviderRabbitMQ")).Body(
					jen.List(jen.ID("rabbitConn"), jen.ID("err")).Op(":=").ID("amqp").Dot("Dial").Call(jen.ID("cfg").Dot("ConnectionURL")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("connecting to rabbitmq: %w"),
							jen.ID("err"),
						))), jen.ID("topic").Op(":=").ID("rabbitpubsub").Dot("OpenTopic").Call(
						jen.ID("rabbitConn"),
						jen.ID("cfg").Dot("Topic"),
						jen.ID("nil"),
					), jen.Return().List(jen.ID("topic"), jen.ID("nil"))),
				jen.Case(jen.ID("ProviderNATS")).Body(
					jen.List(jen.ID("natsConn"), jen.ID("err")).Op(":=").ID("nats").Dot("Connect").Call(jen.ID("cfg").Dot("ConnectionURL")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("connecting to nats: %w"),
							jen.ID("err"),
						))), jen.Return().ID("natspubsub").Dot("OpenTopic").Call(
						jen.ID("natsConn"),
						jen.ID("cfg").Dot("Topic"),
						jen.ID("nil"),
					)),
				jen.Case(jen.ID("ProviderAzureServiceBus")).Body(
					jen.List(jen.ID("busNamespace"), jen.ID("err")).Op(":=").ID("azuresb").Dot("NewNamespaceFromConnectionString").Call(jen.ID("cfg").Dot("ConnectionURL")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("establishing namespace for Azure Service Bus: %w"),
							jen.ID("err"),
						))), jen.List(jen.ID("busTopic"), jen.ID("err")).Op(":=").ID("azuresb").Dot("NewTopic").Call(
						jen.ID("busNamespace"),
						jen.ID("cfg").Dot("Topic"),
						jen.ID("nil"),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("establishing subscription for Azure Service Bus: %w"),
							jen.ID("err"),
						))), jen.Return().ID("azuresb").Dot("OpenTopic").Call(
						jen.ID("ctx"),
						jen.ID("busTopic"),
						jen.ID("nil"),
					)),
				jen.Case(jen.ID("ProviderMemory")).Body(
					jen.ID("topic").Op(":=").ID("mempubsub").Dot("NewTopic").Call(), jen.Return().List(jen.ID("topic"), jen.ID("nil"))),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.ID("errInvalidProvider"))),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideSubscription uses a configuration to provide a pub/sub subscription."),
		jen.Line(),
		jen.Func().ID("ProvideSubscription").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("Config")).Params(jen.Op("*").Qual("gocloud.dev/pubsub", "Subscription"), jen.ID("error")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				jen.Case(jen.ID("ProviderGoogleCloudPubSub")).Body(
					jen.Var().Defs(
						jen.ID("creds").Op("*").ID("google").Dot("Credentials"),
					), jen.If(jen.ID("cfg").Dot("GCPPubSub").Dot("ServiceAccountKeyFilepath").Op("!=").Lit("")).Body(
						jen.List(jen.ID("serviceAccountKeyBytes"), jen.ID("err")).Op(":=").Qual("os", "ReadFile").Call(jen.ID("cfg").Dot("GCPPubSub").Dot("ServiceAccountKeyFilepath")),
						jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
							jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
								jen.Lit("reading service account key file: %w"),
								jen.ID("err"),
							))),
						jen.If(jen.List(jen.ID("creds"), jen.ID("err")).Op("=").ID("google").Dot("CredentialsFromJSON").Call(
							jen.ID("ctx"),
							jen.ID("serviceAccountKeyBytes"),
						), jen.ID("err").Op("!=").ID("nil")).Body(
							jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
								jen.Lit("using service account key credentials: %w"),
								jen.ID("err"),
							))),
					).Else().Body(
						jen.Var().Defs(
							jen.ID("err").ID("error"),
						),
						jen.If(jen.List(jen.ID("creds"), jen.ID("err")).Op("=").Qual("gocloud.dev/gcp", "DefaultCredentials").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
							jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
								jen.Lit("constructing pub/sub credentials: %w"),
								jen.ID("err"),
							))),
					), jen.List(jen.ID("conn"), jen.ID("_"), jen.ID("err")).Op(":=").ID("gcppubsub").Dot("Dial").Call(
						jen.ID("ctx"),
						jen.ID("creds").Dot("TokenSource"),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("dialing connection to pub/sub %w"),
							jen.ID("err"),
						))), jen.List(jen.ID("subClient"), jen.ID("err")).Op(":=").ID("gcppubsub").Dot("SubscriberClient").Call(
						jen.ID("ctx"),
						jen.ID("conn"),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("dialing connection to pub/sub %w"),
							jen.ID("err"),
						))), jen.Return().ID("gcppubsub").Dot("OpenSubscriptionByPath").Call(
						jen.ID("subClient"),
						jen.ID("cfg").Dot("SubscriptionIdentifier"),
						jen.ID("nil"),
					)),
				jen.Case(jen.ID("ProviderAWSSQS")).Body(
					jen.List(jen.ID("sess"), jen.ID("err")).Op(":=").ID("session").Dot("NewSession").Call(jen.ID("nil")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("establishing AWS session: %w"),
							jen.ID("err"),
						))), jen.ID("subscription").Op(":=").ID("awssnssqs").Dot("OpenSubscription").Call(
						jen.ID("ctx"),
						jen.ID("sess"),
						jen.ID("cfg").Dot("SubscriptionIdentifier"),
						jen.ID("nil"),
					), jen.Return().List(jen.ID("subscription"), jen.ID("nil"))),
				jen.Case(jen.ID("ProviderKafka")).Body(
					jen.ID("addrs").Op(":=").Qual("strings", "Split").Call(
						jen.ID("cfg").Dot("ConnectionURL"),
						jen.Lit(","),
					), jen.ID("config").Op(":=").ID("kafkapubsub").Dot("MinimalConfig").Call(), jen.Return().ID("kafkapubsub").Dot("OpenSubscription").Call(
						jen.ID("addrs"),
						jen.ID("config"),
						jen.ID("cfg").Dot("Topic"),
						jen.Index().ID("string").Valuesln(jen.ID("cfg").Dot("SubscriptionIdentifier")),
						jen.ID("nil"),
					)),
				jen.Case(jen.ID("ProviderRabbitMQ")).Body(
					jen.List(jen.ID("rabbitConn"), jen.ID("err")).Op(":=").ID("amqp").Dot("Dial").Call(jen.ID("cfg").Dot("ConnectionURL")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("connecting to rabbitmq: %w"),
							jen.ID("err"),
						))), jen.ID("subscription").Op(":=").ID("rabbitpubsub").Dot("OpenSubscription").Call(
						jen.ID("rabbitConn"),
						jen.ID("cfg").Dot("SubscriptionIdentifier"),
						jen.ID("nil"),
					), jen.Return().List(jen.ID("subscription"), jen.ID("nil"))),
				jen.Case(jen.ID("ProviderNATS")).Body(
					jen.List(jen.ID("natsConn"), jen.ID("err")).Op(":=").ID("nats").Dot("Connect").Call(jen.ID("cfg").Dot("ConnectionURL")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("connecting to NATS: %w"),
							jen.ID("err"),
						))), jen.Return().ID("natspubsub").Dot("OpenSubscription").Call(
						jen.ID("natsConn"),
						jen.ID("cfg").Dot("SubscriptionIdentifier"),
						jen.ID("nil"),
					)),
				jen.Case(jen.ID("ProviderAzureServiceBus")).Body(
					jen.List(jen.ID("busNamespace"), jen.ID("err")).Op(":=").ID("azuresb").Dot("NewNamespaceFromConnectionString").Call(jen.ID("cfg").Dot("ConnectionURL")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("connecting to Azure Service Bus: %w"),
							jen.ID("err"),
						))), jen.List(jen.ID("busTopic"), jen.ID("err")).Op(":=").ID("azuresb").Dot("NewTopic").Call(
						jen.ID("busNamespace"),
						jen.ID("cfg").Dot("Topic"),
						jen.ID("nil"),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("establishing service bus subscription: %w"),
							jen.ID("err"),
						))), jen.List(jen.ID("busSub"), jen.ID("err")).Op(":=").ID("azuresb").Dot("NewSubscription").Call(
						jen.ID("busTopic"),
						jen.ID("cfg").Dot("SubscriptionIdentifier"),
						jen.ID("nil"),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("establishing service bus subscription: %w"),
							jen.ID("err"),
						))), jen.Return().ID("azuresb").Dot("OpenSubscription").Call(
						jen.ID("ctx"),
						jen.ID("busNamespace"),
						jen.ID("busTopic"),
						jen.ID("busSub"),
						jen.ID("nil"),
					)),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.ID("errInvalidProvider"))),
			)),
		jen.Line(),
	)

	return code
}
