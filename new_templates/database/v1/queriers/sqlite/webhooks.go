package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("sqlite")
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
	ret.Add(jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildWebhookCreationTimeQuery").Params(jen.ID("webhookID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
		jen.Null().Var().ID("err").ID("error"),
		jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
			"sqlBuilder",
		).Dot(
			"Select",
		).Call(jen.Lit("created_on")).Dot(
			"From",
		).Call(jen.ID("webhooksTableName")).Dot(
			"Where",
		).Call(jen.ID("squirrel").Dot(
			"Eq",
		).Valuesln(jen.Lit("id").Op(":").ID("webhookID"))).Dot(
			"ToSql",
		).Call(),
		jen.ID("s").Dot(
			"logQueryBuildingError",
		).Call(jen.ID("err")),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	),
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
