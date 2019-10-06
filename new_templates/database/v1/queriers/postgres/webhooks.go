package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("postgres")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("eventsSeparator").Op("=").Lit(`,`).Var().ID("typesSeparator").Op("=").Lit(`,`).Var().ID("topicsSeparator").Op("=").Lit(`,`).Var().ID("webhooksTableName").Op("=").Lit("webhooks"),
	)
	ret.Add(jen.Null().Var().ID("webhooksTableColumns").Op("=").Index().ID("string").Valuesln(jen.Lit("id"), jen.Lit("name"), jen.Lit("content_type"), jen.Lit("url"), jen.Lit("method"), jen.Lit("events"), jen.Lit("data_types"), jen.Lit("topics"), jen.Lit("created_on"), jen.Lit("updated_on"), jen.Lit("archived_on"), jen.Lit("belongs_to")),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Null().Var().ID("getAllWebhooksCountQueryBuilder").Qual("sync", "Once").Var().ID("getAllWebhooksCountQuery").ID("string"),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Null().Var().ID("getAllWebhooksQueryBuilder").Qual("sync", "Once").Var().ID("getAllWebhooksQuery").ID("string"),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
