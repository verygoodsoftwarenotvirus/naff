package events

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func subscriberTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestProvideSubscriber").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Enabled").Op(":").ID("true"), jen.ID("Provider").Op(":").ID("ProviderMemory"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").Lit("mem://whatever"), jen.ID("AckDeadline").Op(":").Qual("time", "Second")),
					jen.ID("topic").Op(":=").ID("mempubsub").Dot("NewTopic").Call(),
					jen.ID("subscription").Op(":=").ID("mempubsub").Dot("NewSubscription").Call(
						jen.ID("topic"),
						jen.ID("cfg").Dot("AckDeadline"),
					),
					jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideSubscriber").Call(
						jen.ID("logger"),
						jen.ID("subscription"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil subscription"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Values(),
					jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideSubscriber").Call(
						jen.ID("logger"),
						jen.ID("nil"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("topic").Op(":=").ID("mempubsub").Dot("NewTopic").Call(),
					jen.ID("subscription").Op(":=").ID("mempubsub").Dot("NewSubscription").Call(
						jen.ID("topic"),
						jen.Qual("time", "Second"),
					),
					jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideSubscriber").Call(
						jen.ID("logger"),
						jen.ID("subscription"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with disabled config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Enabled").Op(":").ID("false")),
					jen.ID("topic").Op(":=").ID("mempubsub").Dot("NewTopic").Call(),
					jen.ID("subscription").Op(":=").ID("mempubsub").Dot("NewSubscription").Call(
						jen.ID("topic"),
						jen.ID("cfg").Dot("AckDeadline"),
					),
					jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideSubscriber").Call(
						jen.ID("logger"),
						jen.ID("subscription"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_subscriber_HandleEvents").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Enabled").Op(":").ID("true"), jen.ID("Provider").Op(":").ID("ProviderMemory"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("SubscriptionIdentifier").Op(":").ID("t").Dot("Name").Call(), jen.ID("ConnectionURL").Op(":").Lit("mem://whatever"), jen.ID("AckDeadline").Op(":").Qual("time", "Second")),
					jen.ID("topic").Op(":=").ID("mempubsub").Dot("NewTopic").Call(),
					jen.ID("subscription").Op(":=").ID("mempubsub").Dot("NewSubscription").Call(
						jen.ID("topic"),
						jen.ID("cfg").Dot("AckDeadline"),
					),
					jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideSubscriber").Call(
						jen.ID("logger"),
						jen.ID("subscription"),
						jen.ID("cfg"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("s"),
					),
					jen.List(jen.ID("p"), jen.ID("err")).Op(":=").ID("ProvidePublisher").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("cfg"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("p"),
					),
					jen.ID("stopChan").Op(":=").ID("make").Call(
						jen.Chan().ID("bool"),
						jen.Lit(1),
					),
					jen.Var().Defs(
						jen.ID("calledHat").Qual("sync", "Mutex"),
						jen.ID("called").ID("bool"),
					),
					jen.ID("deadline").Op(":=").Qual("time", "After").Call(jen.Qual("time", "Second")),
					jen.ID("pollTicker").Op(":=").Qual("time", "NewTicker").Call(jen.Qual("time", "Second").Op("/").Lit(5)),
					jen.Go().Func().Params().Body(
						jen.For().Body(
							jen.Select().Body(
								jen.Case(jen.Op("<-").ID("deadline")).Body(
									jen.ID("stopChan").ReceiveFromChannel().ID("true")),
								jen.Case(jen.Op("<-").ID("pollTicker").Dot("C")).Body(
									jen.ID("require").Dot("NoError").Call(
										jen.ID("t"),
										jen.ID("topic").Dot("Send").Call(
											jen.ID("ctx"),
											jen.Op("&").Qual("gocloud.dev/pubsub", "Message").Valuesln(jen.ID("Body").Op(":").Index().ID("byte").Call(jen.Lit("{}"))),
										),
									)),
							))).Call(),
					jen.Go().ID("s").Dot("HandleEvents").Call(
						jen.Qual("time", "Second").Op("/").Lit(10),
						jen.ID("stopChan"),
						jen.Func().Params(jen.ID("body").Index().ID("byte")).Body(
							jen.ID("calledHat").Dot("Lock").Call(),
							jen.Defer().ID("calledHat").Dot("Unlock").Call(),
							jen.ID("called").Op("=").ID("true"),
						),
					),
					jen.Qual("time", "Sleep").Call(jen.Lit(2).Op("*").Qual("time", "Second")),
					jen.ID("calledHat").Dot("Lock").Call(),
					jen.Defer().ID("calledHat").Dot("Unlock").Call(),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("called"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
