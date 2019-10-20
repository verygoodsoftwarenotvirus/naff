package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientsDotGo() *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("_").ID("models").Dot(
			"OAuth2ClientDataManager",
		).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachOAuth2ClientDatabaseIDToSpan is a consistent way to attach an oauth2 client's ID to a span"),
		jen.Line(),
		jen.Func().ID("attachOAuth2ClientDatabaseIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("oauth2ClientID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("oauth2client_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("oauth2ClientID"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachOAuth2ClientIDToSpan is a consistent way to attach an oauth2 client's Client ID to a span"),
		jen.Line(),
		jen.Func().ID("attachOAuth2ClientIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("clientID").ID("string")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("client_id"), jen.ID("clientID"))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Client gets an OAuth2 client from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetOAuth2Client")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("clientID")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("client_id").Op(":").ID("clientID"), jen.Lit("user_id").Op(":").ID("userID"))),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("GetOAuth2Client called")),
			jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("c").Dot("querier").Dot(
				"GetOAuth2Client",
			).Call(jen.ID("ctx"), jen.ID("clientID"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error fetching oauth2 client from the querier")),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Return().List(jen.ID("client"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2ClientByClientID fetches any OAuth2 client by client ID, regardless of ownership."),
		jen.Line(),
		jen.Comment("This is used by authenticating middleware to fetch client information it needs to validate."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetOAuth2ClientByClientID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetOAuth2ClientByClientID")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachOAuth2ClientIDToSpan").Call(jen.ID("span"), jen.ID("clientID")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("oauth2client_client_id"), jen.ID("clientID")),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("GetOAuth2ClientByClientID called")),
			jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("c").Dot("querier").Dot(
				"GetOAuth2ClientByClientID",
			).Call(jen.ID("ctx"), jen.ID("clientID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error fetching oauth2 client from the querier")),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Return().List(jen.ID("client"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2ClientCount gets the count of OAuth2 clients in the database that match the current filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot("QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetOAuth2ClientCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetOAuth2ClientCount called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"GetOAuth2ClientCount",
			).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2ClientCount gets the count of OAuth2 clients that match the current filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllOAuth2ClientCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetAllOAuth2ClientCount called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"GetAllOAuth2ClientCount",
			).Call(jen.ID("ctx")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2ClientsForUser returns all OAuth2 clients belonging to a given user"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllOAuth2ClientsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Op("*").ID("models").Dot(
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllOAuth2ClientsForUser")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetAllOAuth2ClientsForUser called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"GetAllOAuth2ClientsForUser",
			).Call(jen.ID("ctx"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2Clients returns all OAuth2 clients, irrespective of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Index().Op("*").ID("models").Dot(
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllOAuth2Clients")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetAllOAuth2Clients called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"GetAllOAuth2Clients",
			).Call(jen.ID("ctx")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Clients gets a list of OAuth2 clients"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot("QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("models").Dot(
			"OAuth2ClientList",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetOAuth2Clients")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetOAuth2Clients called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"GetOAuth2Clients",
			).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateOAuth2Client creates an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("CreateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot(
			"OAuth2ClientCreationInput",
		)).Params(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("CreateOAuth2Client")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("client_id").Op(":").ID("input").Dot(
					"ClientID",
				),
				jen.Lit("belongs_to").Op(":").ID("input").Dot("BelongsTo"))),
			jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("c").Dot("querier").Dot(
				"CreateOAuth2Client",
			).Call(jen.ID("ctx"), jen.ID("input")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"WithError",
				).Call(jen.ID("err")).Dot("Debug").Call(jen.Lit("error writing oauth2 client to the querier")),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("new oauth2 client created successfully")),
			jen.Return().List(jen.ID("client"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateOAuth2Client updates a OAuth2 client. Note that this function expects the input's"),
		jen.Line(),
		jen.Comment("ID field to be valid."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("models").Dot(
			"OAuth2Client",
		)).Params(jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("UpdateOAuth2Client")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("c").Dot("querier").Dot(
				"UpdateOAuth2Client",
			).Call(jen.ID("ctx"), jen.ID("updated")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ArchiveOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("ArchiveOAuth2Client")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("clientID")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("client_id").Op(":").ID("clientID"), jen.Lit("belongs_to").Op(":").ID("userID"))),
			jen.ID("err").Op(":=").ID("c").Dot("querier").Dot(
				"ArchiveOAuth2Client",
			).Call(jen.ID("ctx"), jen.ID("clientID"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"WithError",
				).Call(jen.ID("err")).Dot("Debug").Call(jen.Lit("error deleting oauth2 client to the querier")),
				jen.Return().ID("err"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("removed oauth2 client successfully")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)
	return ret
}
