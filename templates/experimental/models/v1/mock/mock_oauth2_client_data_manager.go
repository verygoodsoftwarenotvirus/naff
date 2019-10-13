package mock

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mockOauth2ClientDataManagerDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"OAuth2ClientDataManager",
	).Op("=").Parens(jen.Op("*").ID("OAuth2ClientDataManager")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("OAuth2ClientDataManager").Struct(jen.ID("mock").Dot(
		"Mock",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetOAuth2Client is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("models").Dot(
		"OAuth2Client",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("clientID"), jen.ID("userID")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetOAuth2ClientByClientID is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetOAuth2ClientByClientID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("identifier").ID("string")).Params(jen.Op("*").ID("models").Dot(
		"OAuth2Client",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("identifier")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetOAuth2ClientCount is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllOAuth2ClientCount is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetAllOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllOAuth2Clients is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetAllOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Index().Op("*").ID("models").Dot(
		"OAuth2Client",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Index().Op("*").ID("models").Dot(
			"OAuth2Client",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllOAuth2ClientsForUser is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetAllOAuth2ClientsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Op("*").ID("models").Dot(
		"OAuth2Client",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("userID")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Index().Op("*").ID("models").Dot(
			"OAuth2Client",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetOAuth2Clients is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("models").Dot(
		"OAuth2ClientList",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"OAuth2ClientList",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreateOAuth2Client is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("CreateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot(
		"OAuth2ClientCreationInput",
	)).Params(jen.Op("*").ID("models").Dot(
		"OAuth2Client",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("input")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UpdateOAuth2Client is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("UpdateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("models").Dot(
		"OAuth2Client",
	)).Params(jen.ID("error")).Block(
		jen.Return().ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("updated")).Dot(
			"Error",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ArchiveOAuth2Client is a mock function").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("ArchiveOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
		jen.Return().ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("clientID"), jen.ID("userID")).Dot(
			"Error",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	return ret
}
