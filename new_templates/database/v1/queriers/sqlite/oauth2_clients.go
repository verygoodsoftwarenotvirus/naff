package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func oauth2ClientsDotGo() *jen.File {
	ret := jen.NewFile("sqlite")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("scopesSeparator").Op("=").Lit(`,`).Var().ID("oauth2ClientsTableName").Op("=").Lit("oauth2_clients"),
	)
	ret.Add(jen.Null().Var().ID("oauth2ClientsTableColumns").Op("=").Index().ID("string").Valuesln(jen.Lit("id"), jen.Lit("name"), jen.Lit("client_id"), jen.Lit("scopes"), jen.Lit("redirect_uri"), jen.Lit("client_secret"), jen.Lit("created_on"), jen.Lit("updated_on"), jen.Lit("archived_on"), jen.Lit("belongs_to")),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Null().Var().ID("getAllOAuth2ClientsQueryBuilder").Qual("sync", "Once").Var().ID("getAllOAuth2ClientsQuery").ID("string"),
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
	ret.Add(jen.Null().Var().ID("getAllOAuth2ClientCountQueryBuilder").Qual("sync", "Once").Var().ID("getAllOAuth2ClientCountQuery").ID("string"),
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
