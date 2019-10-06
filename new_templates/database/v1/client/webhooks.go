package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("dbclient")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"WebhookDataManager",
	).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
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
