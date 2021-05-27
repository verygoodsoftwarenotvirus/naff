package events

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestConfig_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderAWSSQS"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.ID("err").Op(":=").ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("t").Dot("Name").Call(), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.ID("err").Op(":=").ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvidePublishTopic").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_memory"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderMemory"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePublishTopic").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_aws"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderAWSSQS"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePublishTopic").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_kafka"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderKafka"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePublishTopic").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_rabbitmq"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderRabbitMQ"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePublishTopic").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_nats"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderNATS"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePublishTopic").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideSubscription").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_gcp"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderGoogleCloudPubSub"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideSubscription").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_aws"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderAWSSQS"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideSubscription").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_kafka"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderKafka"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideSubscription").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_rabbitmq"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderRabbitMQ"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideSubscription").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_nats"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderNATS"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideSubscription").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_azure"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderAzureServiceBus"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideSubscription").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
