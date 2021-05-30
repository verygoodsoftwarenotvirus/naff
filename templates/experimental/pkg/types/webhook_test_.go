package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestWebhook_Update").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").Op("&").ID("WebhookUpdateInput").Valuesln(jen.ID("Name").Op(":").Lit("whatever"), jen.ID("ContentType").Op(":").Lit("application/xml"), jen.ID("URL").Op(":").Lit("https://blah.verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(jen.Lit("more_things")), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(jen.Lit("new_stuff")), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(jen.Lit("blah-blah"))),
					jen.ID("actual").Op(":=").Op("&").ID("Webhook").Valuesln(jen.ID("Name").Op(":").Lit("something_else"), jen.ID("ContentType").Op(":").Lit("application/json"), jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPost"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(jen.Lit("things")), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(jen.Lit("stuff")), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(jen.Lit("blah"))),
					jen.ID("expected").Op(":=").Op("&").ID("Webhook").Valuesln(jen.ID("Name").Op(":").ID("exampleInput").Dot("Name"), jen.ID("ContentType").Op(":").Lit("application/xml"), jen.ID("URL").Op(":").Lit("https://blah.verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(jen.Lit("more_things")), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(jen.Lit("new_stuff")), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(jen.Lit("blah-blah"))),
					jen.ID("actual").Dot("Update").Call(jen.ID("exampleInput")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestWebhookCreationInput_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("buildValidWebhookCreationInput").Op(":=").Func().Params().Params(jen.Op("*").ID("WebhookCreationInput")).Body(
				jen.Return().Op("&").ID("WebhookCreationInput").Valuesln(jen.ID("Name").Op(":").Lit("whatever"), jen.ID("ContentType").Op(":").Lit("application/xml"), jen.ID("URL").Op(":").Lit("https://blah.verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(jen.Lit("more_things")), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(jen.Lit("new_stuff")), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(jen.Lit("blah-blah")))),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("buildValidWebhookCreationInput").Call().Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("bad name"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("Name").Op("=").Lit(""),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("bad url"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("URL").Op("=").Qual("fmt", "Sprintf").Call(
						jen.Lit(`%s://verygoodsoftwarenotvirus.ru`),
						jen.ID("string").Call(jen.ID("byte").Call(jen.Lit(127))),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("bad method"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("Method").Op("=").Lit("balogna"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("bad content type"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("ContentType").Op("=").Lit("application/balogna"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("empty events"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("Events").Op("=").Index().ID("string").Valuesln(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("empty data types"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("DataTypes").Op("=").Index().ID("string").Valuesln(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestWebhookUpdateInput_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("buildValidWebhookCreationInput").Op(":=").Func().Params().Params(jen.Op("*").ID("WebhookUpdateInput")).Body(
				jen.Return().Op("&").ID("WebhookUpdateInput").Valuesln(jen.ID("Name").Op(":").Lit("whatever"), jen.ID("ContentType").Op(":").Lit("application/xml"), jen.ID("URL").Op(":").Lit("https://blah.verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(jen.Lit("more_things")), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(jen.Lit("new_stuff")), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(jen.Lit("blah-blah")))),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("buildValidWebhookCreationInput").Call().Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("bad name"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("Name").Op("=").Lit(""),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("bad url"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("URL").Op("=").Qual("fmt", "Sprintf").Call(
						jen.Lit(`%s://verygoodsoftwarenotvirus.ru`),
						jen.ID("string").Call(jen.ID("byte").Call(jen.Lit(127))),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("bad method"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("Method").Op("=").Lit("balogna"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("bad content type"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("ContentType").Op("=").Lit("application/balogna"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("empty events"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("Events").Op("=").Index().ID("string").Valuesln(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("empty data types"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("buildValidWebhookCreationInput").Call(),
					jen.ID("exampleInput").Dot("DataTypes").Op("=").Index().ID("string").Valuesln(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
